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

	// Aplicar migraciones de estructura y constraints (sin mock data)
	// ApplyAll ya no incluye seeds ni mock data, se aplican separadamente
	err = migrations.ApplyAll(s.db)
	s.Require().NoError(err, "Failed to apply migrations")

	// Aplicar seeds del sistema (roles/permisos base) requeridos por mock data
	err = migrations.ApplySeeds(s.db)
	s.Require().NoError(err, "Failed to apply seeds")

	// Aplicar datos de prueba una sola vez al inicio
	err = migrations.ApplyMockData(s.db)
	s.Require().NoError(err, "Failed to apply mock data")
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
// Bug corregido en postgres/v0.16.4:
// - v0.16.0: Separó tipos en archivo independiente (000_create_types.sql)
// - v0.16.1: Renombró a 000_base_types.sql para orden alfabético correcto
// - v0.16.2: Habilitó extensión uuid-ossp
// - v0.16.3: Actualizó datos de prueba para sistema RBAC
// - v0.16.4: Corrigió formato de UUIDs
func TestPostgresIntegration(t *testing.T) {
	suite.Run(t, new(PostgresIntegrationSuite))
}
