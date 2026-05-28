//go:build integration

// Package roleflow provee setup compartido para los integration tests
// per-rol del paquete `test/integration/<role>_flow/`.
//
// Patrón (Pass 2 — single-path Grants):
//
//  1. Levanta UN postgres efímero por process (TestMain).
//  2. Aplica migrations + demo seed (los 21 usuarios canónicos quedan
//     persistidos con sus user_roles ya mirroreados a iam.role_grants).
//  3. Levanta UN solo identity server contra esa BD; ya no existe
//     "path legacy": el server siempre devuelve `Grants{Allow, Deny}`.
//
// Las suites per-rol llaman `Setup(m)` desde su `TestMain` y luego usan
// `Login(t, env.Server, email, password)` que retorna `LoginResponse`
// con `ActiveContext.Grants.Allow/Deny` (patrones path-based literales).
//
// La aserción central es `AssertGrantsContains` — verifica que el set
// de patterns esperado del seed L4 está literalmente presente en
// `Grants.Allow` (sin matching glob: queremos paridad exacta con la
// matriz `roles_permissions.go`).
package roleflow

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	infraMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"

	identityBuilder "github.com/edugo/edugo-api-identity/cmd/builder"
	identityConfig "github.com/edugo/edugo-api-identity/cmd/config"
)

const (
	sharedJWTSecret = "roleflow-cross-role-jwt-secret-32chars"
	sharedJWTIssuer = "edugo-test-roleflow"

	testDBName     = "test_roleflow"
	testDBUser     = "test"
	testDBPassword = "test"
	testDBImage    = "postgres:16-alpine"

	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour

	// DemoPassword es el password plano de todos los usuarios sembrados
	// en `seeds/demo/development.go` (defaultPasswordHash hashea exactamente
	// este string vía bcrypt). Documentado en seeds/README.md.
	DemoPassword = "12345678"
)

// Env mantiene el handle al identity server más el container postgres.
// Es global por process porque levantar testcontainers por suite tendría
// costo prohibitivo; el seed demo es read-only para los tests.
type Env struct {
	container *tcpostgres.PostgresContainer
	sqlDB     *sql.DB
	DB        *gorm.DB

	Server *httptest.Server
}

var defaultEnv *Env

// Setup arranca el entorno compartido. Llamar desde TestMain de cada
// paquete <role>_flow_test. Si ENABLE_INTEGRATION_TESTS != "true" sale
// con código 0 (skip). En éxito asigna defaultEnv y retorna m.Run().
func Setup(m *testing.M) int {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		fmt.Fprintln(os.Stderr, "roleflow: ENABLE_INTEGRATION_TESTS!=true — skipping")
		return 0
	}

	ctx := context.Background()
	env, err := bootstrap(ctx)
	if err != nil {
		log.Fatalf("roleflow: bootstrap: %v", err)
	}
	defer env.teardown(ctx)

	defaultEnv = env
	return m.Run()
}

// Get devuelve el env activo. Falla si no se llamó Setup primero.
func Get() *Env {
	if defaultEnv == nil {
		panic("roleflow: defaultEnv not initialised — Setup(m) must be called from TestMain")
	}
	return defaultEnv
}

func bootstrap(ctx context.Context) (*Env, error) {
	container, gdb, sqlDB, err := startPostgres(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	env := &Env{container: container, sqlDB: sqlDB, DB: gdb}

	// migrations.Migrate(Force=true, SeedDemo=true) — incluye L0..L4 +
	// los 21 usuarios canónicos del seed `demo/development.go` con sus
	// user_roles asignados, mirroreados 1:1 a iam.role_grants.
	if _, err := infraMigrations.Migrate(sqlDB, infraMigrations.MigrateOptions{
		Force:    true,
		SeedDemo: true,
		DBUser:   testDBUser,
	}); err != nil {
		env.teardown(ctx)
		return nil, fmt.Errorf("migrate: %w", err)
	}

	// ApplyDemo trunca tablas y luego sembra, posiblemente destruyendo
	// las filas L0 que sembraba `system`. Re-aplicar system es idempotente
	// vía OnConflict DoNothing y restaura el catálogo.
	if err := system.ApplySystem(sqlDB, ""); err != nil {
		env.teardown(ctx)
		return nil, fmt.Errorf("re-apply system: %w", err)
	}

	env.Server = startIdentityServer(gdb)
	return env, nil
}

func (e *Env) teardown(ctx context.Context) {
	if e.Server != nil {
		e.Server.Close()
	}
	if e.sqlDB != nil {
		_ = e.sqlDB.Close()
	}
	if e.container != nil {
		_ = e.container.Terminate(ctx)
	}
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

// Grants es el sub-set tipado de `active_context.grants` del payload
// `POST /auth/login`. Match 1:1 con `dto.GrantsDTO` y con `auth.Grants`.
type Grants struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// LoginResponse es un sub-set tipado del payload `POST /auth/login`
// que cubre los campos usados por los tests per-rol.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Schools      []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"schools"`
	ActiveContext *struct {
		RoleID     string `json:"role_id"`
		RoleName   string `json:"role_name"`
		SchoolID   string `json:"school_id"`
		SchoolName string `json:"school_name"`
		Grants     Grants `json:"grants"`
	} `json:"active_context"`
}

// Login ejecuta POST /api/v1/auth/login contra `server`. Falla el test
// si el status no es 200 o el body no parsea.
func Login(t *testing.T, server *httptest.Server, email, password string) LoginResponse {
	t.Helper()
	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost,
		server.URL+"/api/v1/auth/login",
		bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, resp.StatusCode,
		"login %s: expected 200, got %d body=%s", email, resp.StatusCode, string(raw))

	var out LoginResponse
	require.NoError(t, json.Unmarshal(raw, &out), "login: parse body=%s", string(raw))
	require.NotEmpty(t, out.AccessToken, "login: access_token empty")
	require.NotNil(t, out.ActiveContext, "login: active_context nil")
	return out
}

// GetJSON ejecuta GET con bearer y retorna (status, body bytes).
func GetJSON(t *testing.T, server *httptest.Server, path, bearer string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	require.NoError(t, err)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, body
}

// AssertGrantsContains valida que cada permiso esperado está cubierto
// por algún pattern en `grants.Allow` y NO está negado por `grants.Deny`,
// usando el mismo matcher que el runtime (`auth.PermissionMatches`).
// Pass 3 wildcard-first: el seed emite patterns como `academic.*`, así
// que la aserción ya no compara literales — pregunta semántica "este
// rol puede X" y deja al matcher decidir.
func AssertGrantsContains(t *testing.T, grants Grants, expectedAllow ...string) {
	t.Helper()
	missing := make([]string, 0)
	for _, want := range expectedAllow {
		if !GrantsAllow(grants, want) {
			missing = append(missing, want)
		}
	}
	require.Emptyf(t, missing,
		"grants no cubre permisos esperados: %v\n  got allow=%v deny=%v",
		missing, grants.Allow, grants.Deny)
}

// GrantsAllow aplica la semántica deny>allow del matcher: un permiso
// queda permitido si algún pattern de Allow lo matchea y ninguno de
// Deny lo cubre.
func GrantsAllow(grants Grants, request string) bool {
	for _, d := range grants.Deny {
		if auth.PermissionMatches(d, request) {
			return false
		}
	}
	for _, a := range grants.Allow {
		if auth.PermissionMatches(a, request) {
			return true
		}
	}
	return false
}

// noOpLogger silencia los logs de las APIs en tests.
type noOpLogger struct{}

func newNoOpLogger() logger.Logger                 { return &noOpLogger{} }
func (l *noOpLogger) Debug(_ string, _ ...any)    {}
func (l *noOpLogger) Info(_ string, _ ...any)     {}
func (l *noOpLogger) Warn(_ string, _ ...any)     {}
func (l *noOpLogger) Error(_ string, _ ...any)    {}
func (l *noOpLogger) Fatal(_ string, _ ...any)    {}
func (l *noOpLogger) With(_ ...any) logger.Logger { return l }
func (l *noOpLogger) Sync() error                 { return nil }
