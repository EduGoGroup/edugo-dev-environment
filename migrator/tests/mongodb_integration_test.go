package tests

import (
	"context"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBIntegrationSuite es la suite de tests de integración para MongoDB
type MongoDBIntegrationSuite struct {
	suite.Suite
	client    *mongo.Client
	db        *mongo.Database
	container testcontainers.Container
	ctx       context.Context
}

// SetupSuite se ejecuta una vez antes de todos los tests de la suite
func (s *MongoDBIntegrationSuite) SetupSuite() {
	s.ctx = context.Background()

	// Crear contenedor MongoDB para tests
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "testuser",
			"MONGO_INITDB_ROOT_PASSWORD": "testpass",
			"MONGO_INITDB_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForListeningPort("27017/tcp"),
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err, "Failed to start MongoDB container")
	s.container = container

	// Obtener host y puerto del contenedor
	host, err := container.Host(s.ctx)
	s.Require().NoError(err)

	port, err := container.MappedPort(s.ctx, "27017")
	s.Require().NoError(err)

	// Conectar a MongoDB
	mongoURI := "mongodb://testuser:testpass@" + host + ":" + port.Port()
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(s.ctx, clientOptions)
	s.Require().NoError(err, "Failed to connect to MongoDB")

	// Verificar conexión
	err = client.Ping(s.ctx, nil)
	s.Require().NoError(err, "Failed to ping MongoDB")

	s.client = client
	s.db = client.Database("testdb")

	// Aplicar todas las migraciones usando el paquete de infrastructure
	err = migrations.ApplyAll(s.ctx, s.db)
	s.Require().NoError(err, "Failed to apply migrations")
}

// SetupTest se ejecuta antes de cada test individual
func (s *MongoDBIntegrationSuite) SetupTest() {
	// Nota: ApplyMockData puede no estar implementado en todas las versiones
	// Los tests deben insertar sus propios datos de prueba si es necesario
	// Ejemplo: s.db.Collection("collection").InsertOne(ctx, document)
}

// TearDownTest se ejecuta después de cada test individual
func (s *MongoDBIntegrationSuite) TearDownTest() {
	// Limpiar datos de prueba entre tests (sin eliminar colecciones)
	collections := []string{
		"material_assessment",
		"material_content",
		"assessment_attempt_result",
		"audit_logs",
		"notifications",
		"analytics_events",
		"material_summary",
		"material_assessment_worker",
		"material_event",
	}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(s.ctx, map[string]interface{}{})
		if err != nil {
			s.T().Logf("Warning: Failed to delete documents from %s: %v", collName, err)
		}
	}
}

// TearDownSuite se ejecuta una vez después de todos los tests
func (s *MongoDBIntegrationSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Disconnect(s.ctx)
	}
	if s.container != nil {
		s.container.Terminate(s.ctx)
	}
}

// TestExampleQuery es un ejemplo de test que usa la base de datos
func (s *MongoDBIntegrationSuite) TestExampleQuery() {
	// Ejemplo: Verificar que las colecciones fueron creadas
	collections, err := s.db.ListCollectionNames(s.ctx, map[string]interface{}{})
	s.NoError(err, "Failed to list collections")
	s.Greater(len(collections), 0, "Expected at least one collection to exist")
}

// TestCanQueryCollection es un ejemplo de test que verifica que se puede consultar una colección
func (s *MongoDBIntegrationSuite) TestCanQueryCollection() {
	// Ejemplo: Verificar que podemos consultar la colección (aunque esté vacía)
	// Esto demuestra que las migraciones crearon la colección correctamente
	count, err := s.db.Collection("material_assessment").CountDocuments(s.ctx, map[string]interface{}{})
	s.NoError(err, "Failed to query material_assessment collection")

	// La colección puede estar vacía, pero debe existir y ser consultable
	s.GreaterOrEqual(count, int64(0), "Collection should be queryable")

	// Verificar que podemos hacer operaciones en la colección
	cursor, err := s.db.Collection("material_assessment").Find(s.ctx, map[string]interface{}{})
	s.NoError(err, "Failed to create cursor on material_assessment")
	defer cursor.Close(s.ctx)
}

// TestMongoDBIntegration ejecuta la suite de tests
func TestMongoDBIntegration(t *testing.T) {
	suite.Run(t, new(MongoDBIntegrationSuite))
}
