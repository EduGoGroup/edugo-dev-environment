//go:build integration

// Package guardian_invitation_flow valida el plan 024 · F4 · S2: el colegio
// invita a un representante ELIGIENDO al alumno (school_invitations.student_id),
// y al APROBAR el join-request (doble-gate school→unit) se materializa el
// vínculo academic.guardian_relations en estado "active" — cerrando la arista
// "vínculo ausente" del bug 0045.
//
// Por qué un Setup propio (multi-API, igual que guardian_ward_grades_flow): el
// flujo cruza IDENTITY (login + switch-context del admin y del representante) y
// ACADEMIC (crear invitación, redimir, aprobar, leer notas del acudido). Se
// levantan AMBAS APIs in-process sobre la MISMA gorm.DB y el MISMO
// AUTH_JWT_SECRET, sembradas con migrations + playground_v2 `base`.
//
// El Setup expone también el *sql.DB del container para que las aserciones lean
// directamente academic.guardian_relations y para limpiar lo insertado (defer).
package guardian_invitation_flow_test

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
	sharedJWTSecret = "guardian-invitation-flow-secret-32chars"
	sharedJWTIssuer = "edugo-test-guardian-invitation-flow"

	testDBName     = "test_guardian_invitation_flow"
	testDBUser     = "test"
	testDBPassword = "test"
	testDBImage    = "postgres:16-alpine"

	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// Variables globales compartidas por todos los subtests del paquete.
var (
	testDB    *gorm.DB
	testSQLDB *sql.DB // acceso directo al pool para aserciones/limpieza en BD.

	identityServer *httptest.Server
	academicServer *httptest.Server
)

func TestMain(m *testing.M) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		fmt.Fprintln(os.Stderr, "guardian_invitation_flow: ENABLE_INTEGRATION_TESTS!=true — skipping")
		os.Exit(0)
	}

	ctx := context.Background()

	// 1. Postgres efímero.
	container, gdb, sqlDB, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("guardian_invitation_flow: postgres: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
		_ = container.Terminate(ctx)
	}()
	testDB = gdb
	testSQLDB = sqlDB

	// 2. migrations.Migrate(Force=true, PlaygroundV2="base") — DDL L0..L4 + el
	//    mundo `base`: 2 escuelas, usuarios @edugo.test, los vínculos guardian
	//    sembrados (mendoza↔sofia, castro↔carlos/diego) y los catálogos de tipos
	//    de invitación (incl. key="guardian").
	if _, err := infraMigrations.Migrate(sqlDB, infraMigrations.MigrateOptions{
		Force:        true,
		PlaygroundV2: "base",
		DBUser:       testDBUser,
	}); err != nil {
		log.Fatalf("guardian_invitation_flow: migrate: %v", err)
	}

	// 2.1 base.Apply usa upsert idempotente y NO trunca tablas, así que las filas
	//     L0 de system sobreviven. Re-aplicar system (idempotente vía OnConflict
	//     DoNothing por id) blinda el catálogo, igual que guardian_ward_grades_flow.
	if err := system.ApplySystem(sqlDB, ""); err != nil {
		log.Fatalf("guardian_invitation_flow: re-apply system: %v", err)
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
