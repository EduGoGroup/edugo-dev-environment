//go:build integration

// Package guardian_onboarding_policy_flow valida el plan 024 · F4 · S3: la
// politica de representante de la escuela (academic.school_guardian_policy) altera
// la admision de un ALUMNO de dos formas:
//
//   - Tramo B (gates_activation=true): al APROBAR la admision del alumno su
//     membership nace 'pending' y solo se ACTIVA ('active') cuando un representante
//     suyo queda aprobado (gating_approver='any').
//   - Tramo A2 (invitation_mode='on_enrollment'): al aprobar la admision del alumno
//     se AUTO-CREA una invitacion de tipo 'guardian' (codigo opaco, sin email)
//     apuntando al alumno (school_invitations.student_id = alumno).
//
// El comportamiento por defecto (sin fila de politica, o 'manual'/gates=false) es
// el de hoy: la membership nace 'active' y no se genera ninguna invitacion.
//
// Por que un Setup propio (multi-API, igual que guardian_invitation_flow): el flujo
// cruza IDENTITY (signup + login + switch-context) y ACADEMIC (crear invitacion,
// redimir, aprobar, leer estado). Se levantan AMBAS APIs in-process sobre la MISMA
// gorm.DB y el MISMO AUTH_JWT_SECRET, sembradas con migrations + playground_v2 `base`.
//
// El Setup expone tambien el *sql.DB del container para que las aserciones lean
// directamente academic.memberships / school_invitations / guardian_relations, para
// insertar la fila de politica test-local del Test 1 (DML de prueba, NO edita el
// seed) y para limpiar lo insertado (defer).
package guardian_onboarding_policy_flow_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	infraMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"

	academicBuilder "github.com/edugo/edugo-api-academic/cmd/builder"
	academicConfig "github.com/edugo/edugo-api-academic/cmd/config"
	identityBuilder "github.com/edugo/edugo-api-identity/cmd/builder"
	identityConfig "github.com/edugo/edugo-api-identity/cmd/config"
)

// JWT compartido por identity y academic: un token emitido por identity debe ser
// aceptado por academic. Si los secretos divergen, el primer GET con
// Authorization rompe en 401.
const (
	sharedJWTSecret = "guardian-onboarding-policy-flow-secret-32"
	sharedJWTIssuer = "edugo-test-guardian-onboarding-policy-flow"

	testDBName     = "test_guardian_onboarding_policy_flow"
	testDBUser     = "test"
	testDBPassword = "test"
	testDBImage    = "postgres:16-alpine"

	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// Variables globales compartidas por todos los subtests del paquete.
var (
	testDB    *gorm.DB
	testSQLDB *sql.DB // acceso directo al pool para aserciones/limpieza/DML en BD.

	identityServer *httptest.Server
	academicServer *httptest.Server
)

func TestMain(m *testing.M) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		fmt.Fprintln(os.Stderr, "guardian_onboarding_policy_flow: ENABLE_INTEGRATION_TESTS!=true — skipping")
		os.Exit(0)
	}

	ctx := context.Background()

	// 1. Postgres efimero.
	container, gdb, sqlDB, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("guardian_onboarding_policy_flow: postgres: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
		_ = container.Terminate(ctx)
	}()
	testDB = gdb
	testSQLDB = sqlDB

	// 2. migrations.Migrate(Force=true, PlaygroundV2="base") — DDL L0..L4 + el
	//    mundo `base`: 2 escuelas, usuarios @edugo.test, catalogos de invitation
	//    types (incl. key="student"/"guardian") y la fila school_guardian_policy de
	//    S3 (on_enrollment + gates + any).
	if _, err := infraMigrations.Migrate(sqlDB, infraMigrations.MigrateOptions{
		Force:        true,
		PlaygroundV2: "base",
		DBUser:       testDBUser,
	}); err != nil {
		log.Fatalf("guardian_onboarding_policy_flow: migrate: %v", err)
	}

	// 2.1 base.Apply usa upsert idempotente y NO trunca el catalogo de system, asi
	//     que las filas L0 sobreviven. Re-aplicar system (idempotente via OnConflict
	//     DoNothing por id) blinda el catalogo, igual que guardian_invitation_flow.
	if err := system.ApplySystem(sqlDB, ""); err != nil {
		log.Fatalf("guardian_onboarding_policy_flow: re-apply system: %v", err)
	}

	// 3. Levantar identity + academic sobre la MISMA gorm.DB y el MISMO JWT.
	identityServer = startIdentityServer(testDB)
	defer identityServer.Close()

	academicServer = startAcademicServer(testDB)
	defer academicServer.Close()

	os.Exit(m.Run())
}

func startPostgres(ctx context.Context) (*tcpostgres.PostgresContainer, *gorm.DB, *sql.DB, error) {
	container, err := tcpostgres.Run(ctx,
		testDBImage,
		tcpostgres.WithDatabase(testDBName),
		tcpostgres.WithUsername(testDBUser),
		tcpostgres.WithPassword(testDBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("start postgres container: %w", err)
	}
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("get DSN: %w", err)
	}
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("gorm open: %w", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("gorm.DB(): %w", err)
	}
	return container, gdb, sqlDB, nil
}

func startIdentityServer(db *gorm.DB) *httptest.Server {
	cfg := &identityConfig.Config{
		Environment: "development",
		Server: identityConfig.ServerConfig{
			Port: 0,
			Host: "127.0.0.1",
		},
		Auth: identityConfig.AuthConfig{
			JWT: identityConfig.JWTConfig{
				Secret:               sharedJWTSecret,
				Issuer:               sharedJWTIssuer,
				AccessTokenDuration:  accessTokenDuration,
				RefreshTokenDuration: refreshTokenDuration,
			},
		},
		Logging: identityConfig.LoggingConfig{Level: "error", Format: "text"},
		CORS: identityConfig.CORSConfig{
			AllowedOrigins: "*",
			AllowedMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowedHeaders: "Content-Type,Authorization,X-Request-ID",
		},
	}

	blacklist := auth.NewInMemoryBlacklist(context.Background())
	log := newNoOpLogger()
	c := identityBuilder.NewContainer(db, log, cfg, blacklist)
	return httptest.NewServer(c.SetupRouter(cfg, log, "dev", "dev"))
}

func startAcademicServer(db *gorm.DB) *httptest.Server {
	cfg := &academicConfig.Config{
		Environment: "development",
		Server: academicConfig.ServerConfig{
			Port: 0,
			Host: "127.0.0.1",
		},
		Auth: academicConfig.AuthConfig{
			JWT: academicConfig.JWTConfig{
				Secret:               sharedJWTSecret,
				Issuer:               sharedJWTIssuer,
				AccessTokenDuration:  accessTokenDuration,
				RefreshTokenDuration: refreshTokenDuration,
			},
		},
		Logging: academicConfig.LoggingConfig{Level: "error", Format: "text"},
		CORS: academicConfig.CORSConfig{
			AllowedOrigins: "*",
			AllowedMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowedHeaders: "Content-Type,Authorization,X-Request-ID",
		},
	}

	blacklist := auth.NewInMemoryBlacklist(context.Background())
	log := newNoOpLogger()
	auditLog := audit.NewNoopAuditLogger()
	c := academicBuilder.NewContainer(db, log, cfg, blacklist, auditLog)
	return httptest.NewServer(c.SetupRouter(cfg, log, "dev", "dev"))
}

// noOpLogger silencia los logs de los AppServer en tests.
type noOpLogger struct{}

func newNoOpLogger() logger.Logger                { return &noOpLogger{} }
func (l *noOpLogger) Debug(_ string, _ ...any)    {}
func (l *noOpLogger) Info(_ string, _ ...any)     {}
func (l *noOpLogger) Warn(_ string, _ ...any)     {}
func (l *noOpLogger) Error(_ string, _ ...any)    {}
func (l *noOpLogger) Fatal(_ string, _ ...any)    {}
func (l *noOpLogger) With(_ ...any) logger.Logger { return l }
func (l *noOpLogger) Sync() error                 { return nil }
