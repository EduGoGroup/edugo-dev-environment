package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	infraRepoURL = "https://github.com/EduGoGroup/edugo-infrastructure.git"
	infraDir     = ".infrastructure"
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

	// 1. Clonar/actualizar el repositorio de infraestructura
	if err := ensureInfrastructureRepo(); err != nil {
		log.Fatalf("Error obteniendo repositorio de infraestructura: %v", err)
	}

	// 2. Ejecutar migraciones de PostgreSQL
	fmt.Println("\n--- PostgreSQL Migrations ---")
	if err := runPostgresMigrations(forceMigration); err != nil {
		fmt.Printf("⚠️  Error ejecutando migraciones de PostgreSQL: %v\n", err)
		fmt.Println("Continuando con MongoDB...")
	}

	// 3. Ejecutar migraciones de MongoDB
	fmt.Println("\n--- MongoDB Migrations ---")
	if err := runMongoMigrations(forceMigration); err != nil {
		log.Fatalf("Error ejecutando migraciones de MongoDB: %v", err)
	}

	fmt.Println("\n✅ Todas las migraciones se ejecutaron correctamente")
}

func ensureInfrastructureRepo() error {
	fmt.Println("📦 Obteniendo repositorio de infraestructura...")

	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		// Clonar el repositorio si no existe
		fmt.Println("Clonando edugo-infrastructure...")
		cmd := exec.Command("git", "clone", "--branch", "main", infraRepoURL, infraDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error clonando repositorio: %w", err)
		}
	} else {
		// Actualizar si ya existe
		fmt.Println("Actualizando edugo-infrastructure...")
		cmd := exec.Command("git", "pull", "origin", "main")
		cmd.Dir = infraDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error actualizando repositorio: %w", err)
		}
	}

	// Debug: Ver archivos de migraciones disponibles
	structurePath := filepath.Join(infraDir, "postgres", "migrations", "structure")
	files, err := os.ReadDir(structurePath)
	if err == nil {
		fmt.Printf("📁 Archivos en structure: %d\n", len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".sql") {
				fmt.Printf("   - %s\n", f.Name())
			}
		}
	}

	fmt.Println("✅ Repositorio de infraestructura listo")
	return nil
}

func runPostgresMigrations(force bool) error {
	// Configurar variables de entorno
	setPostgresEnv()

	// Si force=true, eliminar schema y recrear
	if force {
		fmt.Println("🔥 Eliminando schema público de PostgreSQL...")
		if err := dropPostgresSchema(); err != nil {
			return fmt.Errorf("error eliminando schema postgres: %w", err)
		}
		fmt.Println("✅ Schema eliminado exitosamente")
	} else {
		// Verificar si ya existen tablas (idempotencia)
		if hasPostgresTables() {
			fmt.Println("✅ PostgreSQL ya tiene tablas - migraciones omitidas (idempotente)")
			return nil
		}
	}

	postgresPath := filepath.Join(infraDir, "postgres")

	// Primero ejecutar runner.go que incluye structure, constraints, seeds y testing
	fmt.Println("Ejecutando runner de PostgreSQL (estructura, constraints, seeds y testing)...")

	// Cambiar al directorio migrations dentro de postgres
	migrationsPath := filepath.Join(postgresPath, "migrations")

	cmd := exec.Command("go", "run", "../cmd/runner/runner.go")
	cmd.Dir = migrationsPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando runner postgres: %w", err)
	}

	fmt.Println("✅ Migraciones de PostgreSQL completadas")
	return nil
}

func runMongoMigrations(force bool) error {
	// Configurar variables de entorno
	setMongoEnv()

	// Si force=true, eliminar base de datos completa
	if force {
		fmt.Println("🔥 Eliminando base de datos MongoDB...")
		if err := dropMongoDatabase(); err != nil {
			return fmt.Errorf("error eliminando database mongodb: %w", err)
		}
		fmt.Println("✅ Base de datos MongoDB eliminada exitosamente")
	} else {
		// Verificar si ya existen colecciones (idempotencia)
		if hasMongoCollections() {
			fmt.Println("✅ MongoDB ya tiene colecciones - migraciones omitidas (idempotente)")
			return nil
		}
	}

	mongoPath := filepath.Join(infraDir, "mongodb")

	// Ejecutar runner.go que incluye structure y constraints
	fmt.Println("Ejecutando runner de MongoDB (estructura y constraints)...")

	cmd := exec.Command("go", "run", "./migrations/cmd/runner.go", "all")
	cmd.Dir = mongoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando runner mongodb: %w", err)
	}

	fmt.Println("✅ Migraciones de MongoDB completadas")
	return nil
}

func setPostgresEnv() {
	// Variables para runner.go (usa POSTGRES_*)
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
		os.Setenv(key, defaultValue)
	}
}

// hasPostgresTables verifica si PostgreSQL ya tiene tablas creadas
func hasPostgresTables() bool {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("⚠️  No se pudo conectar a PostgreSQL para verificar: %v\n", err)
		return false
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("⚠️  No se pudo hacer ping a PostgreSQL: %v\n", err)
		return false
	}

	// Verificar si existe la tabla 'users' (tabla principal)
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'users'
	)`

	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		fmt.Printf("⚠️  Error verificando tablas: %v\n", err)
		return false
	}

	return exists
}

// hasMongoCollections verifica si MongoDB ya tiene colecciones creadas
func hasMongoCollections() bool {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	dbname := os.Getenv("MONGO_DB_NAME")

	// Conectar a MongoDB usando el driver
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("⚠️  No se pudo conectar a MongoDB: %v\n", err)
		return false
	}
	defer client.Disconnect(context.Background())

	// Verificar conexión
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("⚠️  No se pudo hacer ping a MongoDB: %v\n", err)
		return false
	}

	// Obtener lista de colecciones
	db := client.Database(dbname)
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		fmt.Printf("⚠️  Error listando colecciones: %v\n", err)
		return false
	}

	// Si hay al menos 1 colección, consideramos que ya está migrado
	return len(collections) > 0
}

// dropPostgresSchema elimina y recrea el schema público en PostgreSQL
func dropPostgresSchema() error {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a PostgreSQL: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("no se pudo hacer ping a PostgreSQL: %w", err)
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

// dropMongoDatabase elimina completamente la base de datos de MongoDB
func dropMongoDatabase() error {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	dbname := os.Getenv("MONGO_DB_NAME")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("no se pudo conectar a MongoDB: %w", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %w", err)
	}

	// Eliminar la base de datos completa
	db := client.Database(dbname)
	if err := db.Drop(ctx); err != nil {
		return fmt.Errorf("error eliminando base de datos: %w", err)
	}

	return nil
}
