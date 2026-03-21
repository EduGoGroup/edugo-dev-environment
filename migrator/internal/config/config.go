package config

import (
	"fmt"
	"net/url"
	"os"
)

// PostgresConfig contiene la configuración de conexión a PostgreSQL.
type PostgresConfig struct {
	ConnStr string // Connection string listo para database/sql
	User    string // Usuario (para logs y operaciones DROP/GRANT)
}

// MongoConfig contiene la configuración de conexión a MongoDB.
type MongoConfig struct {
	URI    string
	DBName string
}

// Config contiene toda la configuración del migrator leída desde variables de entorno.
type Config struct {
	// Flags de control de ejecución
	ForceMigration bool
	ApplyMockData  bool
	PostgresOnly   bool
	MongoOnly      bool
	StatusOnly     bool

	// Configuración de bases de datos
	Postgres PostgresConfig
	Mongo    MongoConfig
}

// Load carga la configuración completa desde variables de entorno.
// Es el único lugar donde se leen variables de entorno en el migrator.
func Load() Config {
	return Config{
		ForceMigration: os.Getenv("FORCE_MIGRATION") == "true",
		ApplyMockData:  os.Getenv("APPLY_MOCK_DATA") != "false", // default: true
		PostgresOnly:   os.Getenv("POSTGRES_ONLY") == "true",
		MongoOnly:      os.Getenv("MONGO_ONLY") == "true",
		StatusOnly:     os.Getenv("STATUS_ONLY") == "true",
		Postgres:       loadPostgresConfig(),
		Mongo:          loadMongoConfig(),
	}
}

func loadPostgresConfig() PostgresConfig {
	if uri := os.Getenv("POSTGRES_URI"); uri != "" {
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			if parsed, err := url.Parse(uri); err == nil && parsed.User != nil {
				user = parsed.User.Username()
			}
		}
		if user == "" {
			user = "postgres"
		}
		return PostgresConfig{ConnStr: uri, User: user}
	}

	host := envOrDefault("POSTGRES_HOST", "localhost")
	port := envOrDefault("POSTGRES_PORT", "5432")
	user := envOrDefault("POSTGRES_USER", "edugo")
	password := envOrDefault("POSTGRES_PASSWORD", "edugo123")
	dbname := envOrDefault("POSTGRES_DB", "edugo")
	sslmode := envOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	return PostgresConfig{ConnStr: connStr, User: user}
}

func loadMongoConfig() MongoConfig {
	dbName := envOrDefault("MONGO_DB_NAME", "edugo")

	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return MongoConfig{URI: uri, DBName: dbName}
	}

	host := envOrDefault("MONGO_HOST", "localhost")
	port := envOrDefault("MONGO_PORT", "27017")
	user := envOrDefault("MONGO_USER", "edugo")
	password := envOrDefault("MONGO_PASSWORD", "edugo123")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)
	return MongoConfig{URI: uri, DBName: dbName}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
