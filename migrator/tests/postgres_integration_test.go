package tests

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/ontop"
	"github.com/EduGoGroup/edugo-dev-environment/migrator/internal/provision"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresIntegrationSuite es la suite de tests de integración para PostgreSQL
type PostgresIntegrationSuite struct {
	suite.Suite
	db        *sql.DB
	connStr   string
	container testcontainers.Container
	ctx       context.Context
}

// SetupSuite se ejecuta una vez antes de todos los tests de la suite
func (s *PostgresIntegrationSuite) SetupSuite() {
	s.ctx = context.Background()

	// Crear contenedor PostgreSQL para tests
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		// El entrypoint de postgres arranca, inicializa y reinicia: el puerto
		// escucha antes de que acepte conexiones (57P03). Se espera el segundo
		// "ready", como hacen los tests de integración de infra.
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err, "Failed to start PostgreSQL container")
	s.container = container

	// Obtener host y puerto del contenedor
	host, err := container.Host(s.ctx)
	s.Require().NoError(err)

	port, err := container.MappedPort(s.ctx, "5432")
	s.Require().NoError(err)

	// Conectar a PostgreSQL
	connStr := "host=" + host + " port=" + port.Port() + " user=testuser password=testpass dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	s.Require().NoError(err, "Failed to connect to PostgreSQL")

	s.db = db
	s.connStr = connStr

	// Aplicar migraciones completas con system seeds + datos de prueba (fixture base) vía Migrate(Force=true).
	_, err = migrations.Migrate(s.db, migrations.MigrateOptions{
		Force:        true,
		PlaygroundV2: "base",
		DBUser:       "testuser",
	})
	s.Require().NoError(err, "Failed to apply migrations")
}

// SetupTest se ejecuta antes de cada test individual
func (s *PostgresIntegrationSuite) SetupTest() {
	// Los datos ya fueron aplicados en SetupSuite
	// Este método está disponible para lógica de setup específica de cada test
}

// TearDownTest se ejecuta después de cada test individual
func (s *PostgresIntegrationSuite) TearDownTest() {
	// Los datos de prueba se aplican solo una vez y se reutilizan entre tests
	// Si un test necesita limpiar datos específicos, debe hacerlo explícitamente
	// No hacemos TRUNCATE aquí para mantener los datos de prueba disponibles
}

// TearDownSuite se ejecuta una vez después de todos los tests
func (s *PostgresIntegrationSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "Failed to close PostgreSQL connection")
	}
	if s.container != nil {
		err := s.container.Terminate(s.ctx)
		s.NoError(err, "Failed to terminate PostgreSQL container")
	}
}

// TestExampleQuery es un ejemplo de test que usa la base de datos
func (s *PostgresIntegrationSuite) TestExampleQuery() {
	// Verificar que las tablas fueron creadas en los schemas del dominio
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema IN ('auth', 'iam', 'academic', 'content', 'assessment', 'ui_config')").Scan(&count)
	s.NoError(err, "Failed to query tables")
	s.Greater(count, 0, "Expected at least one table to exist in domain schemas")
}

// TestMockDataExists es un ejemplo de test que verifica los datos de prueba
func (s *PostgresIntegrationSuite) TestMockDataExists() {
	// Verificar que hay datos de prueba en la tabla auth.users
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM auth.users").Scan(&count)
	s.NoError(err, "Failed to query users")
	s.Greater(count, 0, "Expected mock data in auth.users table")
}

// TestOnTop cubre el flag -on-top (plan 053 T5.2) sobre la BD migrada con base:
// sin los usuarios del escenario falla listándolos; con ellos (insertados como
// lo haría la API: UUID nuevo, hash inutilizable) siembra y encola su sync, y
// el lookup del alta resuelve esos mismos ids.
func (s *PostgresIntegrationSuite) TestOnTop() {
	const name = "onboarding"
	users, err := playground_v2.ScenarioUsers(name)
	s.Require().NoError(err)

	err = ontop.Apply(s.db, name)
	s.Require().Error(err, "sin usuarios debía fallar")
	for _, u := range users {
		s.Contains(err.Error(), u.Email)
	}

	s.Require().Error(ontop.Apply(s.db, "base"), "base no es on-top")

	ids := make([]string, 0, len(users))
	for _, u := range users {
		id := uuid.NewString()
		ids = append(ids, id)
		_, err := s.db.Exec(`INSERT INTO auth.users (id, email, first_name, last_name, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, true, now(), now())`, id, u.Email, u.FirstName, u.LastName)
		s.Require().NoError(err, "insertando %s", u.Email)
	}

	lookup := provision.SQLLookup{DB: s.db}
	for i, u := range users {
		got, found, err := lookup.UserIDByEmail(s.ctx, u.Email)
		s.Require().NoError(err)
		s.True(found)
		s.Equal(ids[i], got)
	}
	_, found, err := lookup.UserIDByEmail(s.ctx, "nadie@edugo.local")
	s.Require().NoError(err)
	s.False(found)

	s.Require().NoError(ontop.Apply(s.db, name))
	var outbox int
	s.Require().NoError(s.db.QueryRow(
		`SELECT COUNT(*) FROM auth.identity_sync_outbox WHERE user_id::text = ANY(string_to_array($1, ','))`,
		strings.Join(ids, ",")).Scan(&outbox))
	s.Equal(len(users), outbox, "una fila de outbox por usuario del escenario")

	// Re-ejecutable.
	s.Require().NoError(ontop.Apply(s.db, name))
}

// TestOnTopRejectsUnmigratedDB: -on-top no migra; sobre una BD vacía se niega.
func (s *PostgresIntegrationSuite) TestOnTopRejectsUnmigratedDB() {
	_, err := s.db.Exec(`CREATE DATABASE ontop_empty`)
	s.Require().NoError(err)
	empty, err := sql.Open("postgres", strings.Replace(s.connStr, "dbname=testdb", "dbname=ontop_empty", 1))
	s.Require().NoError(err)
	defer empty.Close()

	err = ontop.Apply(empty, "onboarding")
	s.Require().Error(err)
	s.Contains(err.Error(), "no está migrada")
}

// TestPostgresIntegration ejecuta la suite de tests
// Bug corregido en postgres/v0.16.4:
// - v0.16.0: Separó tipos en archivo independiente (000_create_types.sql)
// - v0.16.1: Renombró a 000_base_types.sql para orden alfabético correcto
// - v0.16.2: Habilitó extensión uuid-ossp
// - v0.16.3: Actualizó datos de prueba para sistema RBAC
// - v0.16.4: Corrigió formato de UUIDs
func TestPostgresIntegration(t *testing.T) {
	suite.Run(t, new(PostgresIntegrationSuite))
}
