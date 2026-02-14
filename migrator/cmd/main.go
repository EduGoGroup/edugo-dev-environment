package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoMigrations "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
)

func main() {
	fmt.Println("=== EduGo Migrator ===")
	fmt.Println("Iniciando proceso de migraciones...")

	// Verificar si se solicita migración forzada
	forceMigration := os.Getenv("FORCE_MIGRATION") == "true"
	if forceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO - Se eliminarán y recrearán todas las bases de datos")
	}
	fmt.Println()

	// 1. Ejecutar migraciones de PostgreSQL
	fmt.Println("\n--- PostgreSQL Migrations ---")
	if err := runPostgresMigrations(forceMigration); err != nil {
		log.Fatalf("❌ Error ejecutando migraciones de PostgreSQL: %v\n", err)
	}

	// 2. Ejecutar migraciones de MongoDB
	fmt.Println("\n--- MongoDB Migrations ---")
	if err := runMongoMigrations(forceMigration); err != nil {
		log.Fatalf("❌ Error ejecutando migraciones de MongoDB: %v\n", err)
	}

	fmt.Println("\n✅ Todas las migraciones se ejecutaron correctamente")
}

func runPostgresMigrations(force bool) error {
	// Configurar variables de entorno
	setPostgresEnv()

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Conectar a PostgreSQL
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

	fmt.Printf("✓ Conectado a PostgreSQL: %s@%s:%s/%s\n", user, host, port, dbname)

	// Si force=true, eliminar schema y recrear
	if force {
		fmt.Println("🔥 Eliminando schema público de PostgreSQL...")
		if err := dropPostgresSchema(db, user); err != nil {
			return fmt.Errorf("error eliminando schema postgres: %w", err)
		}
		fmt.Println("✅ Schema eliminado exitosamente")
	} else {
		// Verificar si ya existen tablas (idempotencia)
		if hasPostgresTables(db) {
			fmt.Println("✅ PostgreSQL ya tiene tablas - migraciones omitidas (idempotente)")
			return nil
		}
	}

	// Aplicar todas las migraciones usando el paquete importado
	fmt.Println("📦 Aplicando migraciones de estructura y constraints...")
	if err := postgresMigrations.ApplyAll(db); err != nil {
		return fmt.Errorf("error aplicando migraciones: %w", err)
	}

	// Aplicar seeds (datos esenciales del sistema: roles, permisos, etc.)
	fmt.Println("📦 Aplicando datos iniciales (seeds)...")
	if err := postgresMigrations.ApplySeeds(db); err != nil {
		return fmt.Errorf("error aplicando seeds: %w", err)
	}

	// Aplicar datos de prueba/testing
	fmt.Println("📦 Aplicando datos de prueba (testing)...")
	if err := postgresMigrations.ApplyMockData(db); err != nil {
		return fmt.Errorf("error aplicando datos de prueba: %w", err)
	}

	fmt.Println("✅ Migraciones de PostgreSQL completadas")
	return nil
}

func runMongoMigrations(force bool) error {
	// Configurar variables de entorno
	setMongoEnv()

	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	dbname := os.Getenv("MONGO_DB_NAME")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Conectar a MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("no se pudo conectar a MongoDB: %w", err)
	}
	defer func() {
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("⚠️  Error desconectando cliente MongoDB: %v", disconnectErr)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %w", err)
	}

	fmt.Printf("✓ Conectado a MongoDB: %s@%s:%s/%s\n", user, host, port, dbname)

	db := client.Database(dbname)

	// Si force=true, eliminar base de datos completa
	if force {
		fmt.Println("🔥 Eliminando base de datos MongoDB...")
		if err := db.Drop(ctx); err != nil {
			return fmt.Errorf("error eliminando database mongodb: %w", err)
		}
		fmt.Println("✅ Base de datos MongoDB eliminada exitosamente")
	} else {
		// Verificar si ya existen colecciones (idempotencia)
		if hasMongoCollections(ctx, db) {
			fmt.Println("✅ MongoDB ya tiene colecciones - migraciones omitidas (idempotente)")
			return nil
		}
	}

	// Aplicar todas las migraciones usando el paquete importado
	fmt.Println("📦 Aplicando migraciones de estructura y constraints...")
	if err := mongoMigrations.ApplyAll(ctx, db); err != nil {
		return fmt.Errorf("error aplicando migraciones: %w", err)
	}

	fmt.Println("✅ Migraciones de MongoDB completadas")
	return nil
}

func setPostgresEnv() {
	setEnvIfEmpty("POSTGRES_HOST", "localhost")
	setEnvIfEmpty("POSTGRES_PORT", "5432")
	setEnvIfEmpty("POSTGRES_DB", "edugo")
	setEnvIfEmpty("POSTGRES_USER", "edugo")
	setEnvIfEmpty("POSTGRES_PASSWORD", "edugo123")
}

func setMongoEnv() {
	setEnvIfEmpty("MONGO_HOST", "localhost")
	setEnvIfEmpty("MONGO_PORT", "27017")
	setEnvIfEmpty("MONGO_USER", "edugo")
	setEnvIfEmpty("MONGO_PASSWORD", "edugo123")
	setEnvIfEmpty("MONGO_DB_NAME", "edugo")
}

func setEnvIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		if err := os.Setenv(key, defaultValue); err != nil {
			log.Printf("⚠️  Error configurando variable de entorno %s: %v", key, err)
		}
	}
}

// hasPostgresTables verifica si PostgreSQL ya tiene tablas creadas
func hasPostgresTables(db *sql.DB) bool {
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'users'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("⚠️  Error verificando tablas: %v\n", err)
		return false
	}

	return exists
}

// hasMongoCollections verifica si MongoDB ya tiene colecciones creadas
func hasMongoCollections(ctx context.Context, db *mongo.Database) bool {
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		fmt.Printf("⚠️  Error listando colecciones: %v\n", err)
		return false
	}

	return len(collections) > 0
}

// dropPostgresSchema elimina y recrea el schema público en PostgreSQL
func dropPostgresSchema(db *sql.DB, user string) error {
	// Eliminar schema ui_config (Dynamic UI) antes del schema público
	_, err := db.Exec("DROP SCHEMA IF EXISTS ui_config CASCADE")
	if err != nil {
		return fmt.Errorf("error eliminando schema ui_config: %w", err)
	}

	// Eliminar schema público CASCADE (elimina todas las tablas, funciones, etc.)
	_, err = db.Exec("DROP SCHEMA public CASCADE")
	if err != nil {
		return fmt.Errorf("error eliminando schema: %w", err)
	}

	// Recrear schema público
	_, err = db.Exec("CREATE SCHEMA public")
	if err != nil {
		return fmt.Errorf("error creando schema: %w", err)
	}

	// Otorgar permisos
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
