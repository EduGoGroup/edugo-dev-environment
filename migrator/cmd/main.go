package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	mongoMigrations "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	postgresSeeds "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds"
)

func main() {
	fmt.Println("=== EduGo Migrator ===")
	fmt.Println("Iniciando proceso de migraciones...")

	forceMigration := os.Getenv("FORCE_MIGRATION") == "true"
	applyMockData := os.Getenv("APPLY_MOCK_DATA") != "false" // default: true
	postgresOnly := os.Getenv("POSTGRES_ONLY") == "true"
	mongoOnly := os.Getenv("MONGO_ONLY") == "true"

	if forceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO - Se eliminarán y recrearán todas las bases de datos")
	}
	if !applyMockData {
		fmt.Println("ℹ️  Mock data deshabilitado (APPLY_MOCK_DATA=false)")
	}
	fmt.Println()

	// Ejecutar migraciones según flags
	if !mongoOnly {
		fmt.Println("\n--- PostgreSQL Migrations ---")
		if err := runPostgresMigrations(forceMigration, applyMockData); err != nil {
			log.Fatalf("❌ Error ejecutando migraciones de PostgreSQL: %v\n", err)
		}
	}

	if !postgresOnly {
		fmt.Println("\n--- MongoDB Migrations ---")
		if err := runMongoMigrations(forceMigration); err != nil {
			log.Fatalf("❌ Error ejecutando migraciones de MongoDB: %v\n", err)
		}
	}

	fmt.Println("\n✅ Todas las migraciones se ejecutaron correctamente")
}

// buildPostgresConnStr construye el connection string de PostgreSQL.
// Prioriza POSTGRES_URI si está definido (para Neon, Supabase, etc.).
// Si no, construye el string con las variables individuales.
func buildPostgresConnStr() (connStr string, user string) {
	if uri := os.Getenv("POSTGRES_URI"); uri != "" {
		// Extraer user para logs (best effort)
		user = os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "(from URI)"
		}
		return uri, user
	}

	setEnvIfEmpty("POSTGRES_HOST", "localhost")
	setEnvIfEmpty("POSTGRES_PORT", "5432")
	setEnvIfEmpty("POSTGRES_DB", "edugo")
	setEnvIfEmpty("POSTGRES_USER", "edugo")
	setEnvIfEmpty("POSTGRES_PASSWORD", "edugo123")
	setEnvIfEmpty("POSTGRES_SSLMODE", "disable")

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user = os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("POSTGRES_SSLMODE")

	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	return connStr, user
}

// buildMongoURI construye el URI de MongoDB.
// Prioriza MONGO_URI si está definido (para Atlas con mongodb+srv://, etc.).
// Si no, construye el URI con las variables individuales.
func buildMongoURI() (uri string, dbName string) {
	dbName = os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "edugo"
	}

	if uri = os.Getenv("MONGO_URI"); uri != "" {
		return uri, dbName
	}

	setEnvIfEmpty("MONGO_HOST", "localhost")
	setEnvIfEmpty("MONGO_PORT", "27017")
	setEnvIfEmpty("MONGO_USER", "edugo")
	setEnvIfEmpty("MONGO_PASSWORD", "edugo123")

	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")

	uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)
	return uri, dbName
}

func runPostgresMigrations(force bool, applyMock bool) error {
	connStr, user := buildPostgresConnStr()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a PostgreSQL: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("⚠️  Error cerrando conexión PostgreSQL: %v", closeErr)
		}
	}()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("no se pudo hacer ping a PostgreSQL: %w", err)
	}

	fmt.Printf("✓ Conectado a PostgreSQL (user=%s)\n", user)

	if force {
		fmt.Println("🔥 Eliminando schemas de PostgreSQL...")
		pgUser := os.Getenv("POSTGRES_USER")
		if pgUser == "" {
			pgUser = user
		}
		if err := dropPostgresSchema(db, pgUser); err != nil {
			return fmt.Errorf("error eliminando schemas postgres: %w", err)
		}
		fmt.Println("✅ Schemas eliminados exitosamente")
	} else {
		if hasPostgresTables(db) {
			fmt.Println("✅ PostgreSQL ya tiene tablas - migraciones omitidas (idempotente)")
			return nil
		}
	}

	fmt.Println("📦 Aplicando migraciones de estructura...")
	if err := postgresMigrations.ApplyAll(db); err != nil {
		return fmt.Errorf("error aplicando migraciones: %w", err)
	}

	fmt.Println("📦 Aplicando datos de producción (seeds)...")
	if err := postgresSeeds.ApplyProduction(db); err != nil {
		return fmt.Errorf("error aplicando seeds de producción: %w", err)
	}

	if applyMock {
		fmt.Println("📦 Aplicando datos de desarrollo...")
		if err := postgresSeeds.ApplyDevelopment(db); err != nil {
			return fmt.Errorf("error aplicando datos de desarrollo: %w", err)
		}
	} else {
		fmt.Println("⏭️  Datos de desarrollo deshabilitados")
	}

	fmt.Println("✅ Migraciones de PostgreSQL completadas")
	return nil
}

func runMongoMigrations(force bool) error {
	mongoURI, dbName := buildMongoURI()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("no se pudo conectar a MongoDB: %w", err)
	}
	defer func() {
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("⚠️  Error desconectando cliente MongoDB: %v", disconnectErr)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %w", err)
	}

	fmt.Printf("✓ Conectado a MongoDB (db=%s)\n", dbName)

	db := client.Database(dbName)

	if force {
		fmt.Println("🔥 Eliminando base de datos MongoDB...")
		if err := db.Drop(ctx); err != nil {
			return fmt.Errorf("error eliminando database mongodb: %w", err)
		}
		fmt.Println("✅ Base de datos MongoDB eliminada exitosamente")
	} else {
		if hasMongoCollections(ctx, db) {
			fmt.Println("✅ MongoDB ya tiene colecciones - migraciones omitidas (idempotente)")
			return nil
		}
	}

	fmt.Println("📦 Aplicando migraciones de estructura y constraints...")
	if err := mongoMigrations.ApplyAll(ctx, db); err != nil {
		return fmt.Errorf("error aplicando migraciones: %w", err)
	}

	fmt.Println("✅ Migraciones de MongoDB completadas")
	return nil
}

func setEnvIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		if err := os.Setenv(key, defaultValue); err != nil {
			log.Printf("⚠️  Error configurando variable de entorno %s: %v", key, err)
		}
	}
}

func hasPostgresTables(db *sql.DB) bool {
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'auth'
		AND table_name = 'users'
	)`
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("⚠️  Error verificando tablas: %v\n", err)
		return false
	}
	return exists
}

func hasMongoCollections(ctx context.Context, db *mongo.Database) bool {
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		fmt.Printf("⚠️  Error listando colecciones: %v\n", err)
		return false
	}

	return len(collections) > 0
}

func dropPostgresSchema(db *sql.DB, user string) error {
	schemas := []string{"ui_config", "assessment", "content", "academic", "iam", "auth", "public"}
	for _, schema := range schemas {
		_, err := db.Exec("DROP SCHEMA IF EXISTS " + schema + " CASCADE")
		if err != nil {
			return fmt.Errorf("error eliminando schema %s: %w", schema, err)
		}
	}
	// Recreate public schema
	_, err := db.Exec("CREATE SCHEMA public")
	if err != nil {
		return fmt.Errorf("error creando schema: %w", err)
	}
	_, err = db.Exec("GRANT ALL ON SCHEMA public TO " + user)
	if err != nil {
		return fmt.Errorf("error otorgando permisos al usuario: %w", err)
	}
	_, err = db.Exec("GRANT ALL ON SCHEMA public TO public")
	if err != nil {
		return fmt.Errorf("error otorgando permisos públicos: %w", err)
	}
	return nil
}
