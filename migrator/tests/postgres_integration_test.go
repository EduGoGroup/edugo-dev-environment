package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresIntegrationSuite es la suite de tests de integración para PostgreSQL
type PostgresIntegrationSuite struct {
	suite.Suite
	db        *sql.DB
	container testcontainers.Container
	ctx       context.Context
}

// SetupSuite se ejecuta una vez antes de todos los tests de la suite
func (s *PostgresIntegrationSuite) SetupSuite() {
	s.ctx = context.Background()

	// Crear contenedor PostgreSQL para tests
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
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

	// Aplicar todas las migraciones usando el paquete de infrastructure
	err = migrations.ApplyAll(s.db)
	s.Require().NoError(err, "Failed to apply migrations")
}

// SetupTest se ejecuta antes de cada test individual
func (s *PostgresIntegrationSuite) SetupTest() {
	// Aplicar datos de prueba (mock data) para cada test
	err := migrations.ApplyMockData(s.db)
	s.Require().NoError(err, "Failed to apply mock data")
}

// TearDownTest se ejecuta después de cada test individual
func (s *PostgresIntegrationSuite) TearDownTest() {
	// Limpiar datos de prueba entre tests
	_, err := s.db.Exec("TRUNCATE users, schools, academic_units, memberships, materials, assessment, assessment_attempt, assessment_attempt_answer CASCADE")
	s.NoError(err, "Failed to truncate tables")
}

// TearDownSuite se ejecuta una vez después de todos los tests
func (s *PostgresIntegrationSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.container != nil {
		s.container.Terminate(s.ctx)
	}
}

// TestExampleQuery es un ejemplo de test que usa la base de datos
func (s *PostgresIntegrationSuite) TestExampleQuery() {
	// Ejemplo: Verificar que las tablas fueron creadas
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
	s.NoError(err, "Failed to query tables")
	s.Greater(count, 0, "Expected at least one table to exist")
}

// TestMockDataExists es un ejemplo de test que verifica los datos de prueba
func (s *PostgresIntegrationSuite) TestMockDataExists() {
	// Ejemplo: Verificar que hay datos de prueba en la tabla users
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	s.NoError(err, "Failed to query users")
	s.Greater(count, 0, "Expected mock data in users table")
}

// TestPostgresIntegration ejecuta la suite de tests
// TEMPORALMENTE DESHABILITADO: Bug en postgres/v0.15.0
// El archivo structure/000_create_functions.sql usa el tipo permission_scope
// que se crea después en structure/013_create_permissions.sql
// TODO: Habilitar cuando se corrija el orden de migraciones en edugo-infrastructure
func TestPostgresIntegration(t *testing.T) {
	t.Skip("DESHABILITADO: Bug en orden de migraciones de postgres/v0.15.0 - tipo permission_scope se usa antes de crearse")
	suite.Run(t, new(PostgresIntegrationSuite))
}
