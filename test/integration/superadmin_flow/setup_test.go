//go:build integration

// Package superadmin_flow exercita el flujo end-to-end del super_admin
// "global" (rol L0 super_admin, school_id=NULL, sin membership) atravesando
// los tres AppServer reales del backend: identity, academic, platform.
//
// El test reproduce los 4 bugs detectados en sesión 2026-05-12 y los
// blinda con asserts cross-API:
//
//  1. Seed L4: super_admin recibe context:browse_schools + context:browse_units.
//  2. academic.routes_school/routes_unit listan con RequireAnyPermission(read|browse).
//  3. KMP SchoolModels usa el campo `schools` (no `data`) — el shape del DTO
//     del backend lo respalda.
//  4. switch_context.go no exige membership cuando el rol matched es global,
//     e inyecta school_id/name/academic_unit_id/name en el contexto retornado.
//
// Arquitectura (A2 + B3 + C1 + D):
//
//   - A2 cross-API único: 3 httptest.Server en el mismo proceso compartiendo
//     JWT_SECRET vía constante.
//   - B3 fixture e2e nueva: global_user_no_membership + scenario
//     super_admin_global_flow.
//   - C1 testcontainers por suite: 1 postgres:16-alpine efímero levantado en
//     TestMain.
//   - D CI: build tag `integration` + env ENABLE_INTEGRATION_TESTS=true.
package superadmin_flow_test

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
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/scenarios"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"

	academicBuilder "github.com/edugo/edugo-api-academic/cmd/builder"
	academicConfig "github.com/edugo/edugo-api-academic/cmd/config"
	identityBuilder "github.com/edugo/edugo-api-identity/cmd/builder"
	identityConfig "github.com/edugo/edugo-api-identity/cmd/config"
	platformBuilder "github.com/edugo/edugo-api-platform/cmd/builder"
	platformConfig "github.com/edugo/edugo-api-platform/cmd/config"
)

// JWT compartido por las 3 APIs. El test obliga a que un token emitido por
// identity sea aceptado por academic y platform — si los secretos divergen,
// la cadena rompe en el primer GET con Authorization.
const (
	sharedJWTSecret = "superadmin-flow-cross-api-secret-32chars"
	sharedJWTIssuer = "edugo-test-superadmin-flow"

	testDBName     = "test_superadmin_flow"
	testDBUser     = "test"
	testDBPassword = "test"
	testDBImage    = "postgres:16-alpine"

	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// Variables globales compartidas por todos los subtests.
var (
	testDB *gorm.DB

	identityServer *httptest.Server
	academicServer *httptest.Server
	platformServer *httptest.Server

	// Constantes exportadas por el scenario super_admin_global_flow,
	// resueltas en TestMain. Los tests las usan para evitar hardcode de
	// emails/passwords/UUIDs.
	globalUserEmail    string
	globalUserPassword string
)

func TestMain(m *testing.M) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		fmt.Fprintln(os.Stderr, "superadmin_flow: ENABLE_INTEGRATION_TESTS!=true — skipping")
		os.Exit(0)
	}

	ctx := context.Background()

	// 1. Postgres efímero.
	container, gdb, sqlDB, err := startPostgres(ctx)
	if err != nil {
		log.Fatalf("superadmin_flow: postgres: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
		_ = container.Terminate(ctx)
	}()
	testDB = gdb

	// 2. migrations.Migrate(Force=true, SeedDemo=true) — DDL + L0..L4 +
	//    demo (necesario para tener academic_units listables).
	if _, err := infraMigrations.Migrate(sqlDB, infraMigrations.MigrateOptions{
		Force:    true,
		SeedDemo: true,
		DBUser:   testDBUser,
	}); err != nil {
		log.Fatalf("superadmin_flow: migrate: %v", err)
	}

	// 2.1 demo.ApplyDemo trunca auth.users y academic.schools (entre
	// otras tablas) destruyendo las filas L0 sembradas por system.
	// Re-aplicamos system para restaurarlas — es idempotente vía
	// OnConflict DoNothing por id.
	if err := system.ApplySystem(sqlDB, ""); err != nil {
		log.Fatalf("superadmin_flow: re-apply system: %v", err)
	}

	// 3. Aplicar scenario super_admin_global_flow (crea el global user
	//    sin membership + escuela baseline).
	reg := framework.NewRegistry()
	if err := scenarios.RegisterAll(reg); err != nil {
		log.Fatalf("superadmin_flow: register scenarios: %v", err)
	}
	composer := framework.NewComposer(reg, framework.NewNopLogger())
	applyCtx, err := composer.Apply(testDB, "super_admin_global_flow")
	if err != nil {
		log.Fatalf("superadmin_flow: apply scenario: %v", err)
	}
	globalUserEmail = applyCtx.Constants["E2EFixtureGlobalUserEmail"]
	globalUserPassword = applyCtx.Constants["E2EFixtureGlobalUserPassword"]
	if globalUserEmail == "" || globalUserPassword == "" {
		log.Fatalf("superadmin_flow: missing exported constants email=%q password=%q",
			globalUserEmail, globalUserPassword)
	}

	// 4. Levantar los 3 AppServer en httptest.Server distintos,
	//    compartiendo gorm.DB y JWT secret.
	identityServer = startIdentityServer(testDB)
	defer identityServer.Close()

	academicServer = startAcademicServer(testDB)
	defer academicServer.Close()

	platformServer = startPlatformServer(testDB)
	defer platformServer.Close()

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
	return httptest.NewServer(c.SetupRouter(cfg, log))
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
	return httptest.NewServer(c.SetupRouter(cfg, log))
}

func startPlatformServer(db *gorm.DB) *httptest.Server {
	cfg := &platformConfig.Config{
		Environment: "development",
		Server: platformConfig.ServerConfig{
			Port: 0,
			Host: "127.0.0.1",
		},
		Auth: platformConfig.AuthConfig{
			JWT: platformConfig.JWTConfig{
				Secret:               sharedJWTSecret,
				Issuer:               sharedJWTIssuer,
				AccessTokenDuration:  accessTokenDuration,
				RefreshTokenDuration: refreshTokenDuration,
			},
		},
		Logging: platformConfig.LoggingConfig{Level: "error", Format: "json"},
		CORS: platformConfig.CORSConfig{
			AllowedOrigins: "*",
			AllowedMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowedHeaders: "Origin,Content-Type,Accept,Authorization",
		},
	}

	blacklist := auth.NewInMemoryBlacklist(context.Background())
	log := newNoOpLogger()
	auditLog := audit.NewNoopAuditLogger()
	c := platformBuilder.NewContainer(db, log, cfg, blacklist, auditLog, nil)
	return httptest.NewServer(c.SetupRouter(cfg, log))
}

// noOpLogger silencia los logs de los AppServer en tests.
type noOpLogger struct{}

func newNoOpLogger() logger.Logger                  { return &noOpLogger{} }
func (l *noOpLogger) Debug(_ string, _ ...any)     {}
func (l *noOpLogger) Info(_ string, _ ...any)      {}
func (l *noOpLogger) Warn(_ string, _ ...any)      {}
func (l *noOpLogger) Error(_ string, _ ...any)     {}
func (l *noOpLogger) Fatal(_ string, _ ...any)     {}
func (l *noOpLogger) With(_ ...any) logger.Logger  { return l }
func (l *noOpLogger) Sync() error                  { return nil }
