package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	fmt.Println()

	// 1. Clonar/actualizar el repositorio de infraestructura
	if err := ensureInfrastructureRepo(); err != nil {
		log.Fatalf("Error obteniendo repositorio de infraestructura: %v", err)
	}

	// 2. Ejecutar migraciones de PostgreSQL
	fmt.Println("\n--- PostgreSQL Migrations ---")
	if err := runPostgresMigrations(); err != nil {
		fmt.Printf("⚠️  Error ejecutando migraciones de PostgreSQL: %v\n", err)
		fmt.Println("Continuando con MongoDB...")
	}

	// 3. Ejecutar migraciones de MongoDB
	fmt.Println("\n--- MongoDB Migrations ---")
	if err := runMongoMigrations(); err != nil {
		log.Fatalf("Error ejecutando migraciones de MongoDB: %v", err)
	}

	fmt.Println("\n✅ Todas las migraciones se ejecutaron correctamente")
}

func ensureInfrastructureRepo() error {
	fmt.Println("📦 Obteniendo repositorio de infraestructura...")

	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		// Clonar el repositorio si no existe
		fmt.Println("Clonando edugo-infrastructure...")
		cmd := exec.Command("git", "clone", infraRepoURL, infraDir)
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

	fmt.Println("✅ Repositorio de infraestructura listo")
	return nil
}

func runPostgresMigrations() error {
	// Configurar variables de entorno
	setPostgresEnv()

	// Verificar si ya existen tablas (idempotencia)
	if hasPostgresTables() {
		fmt.Println("✅ PostgreSQL ya tiene tablas - migraciones omitidas (idempotente)")
		return nil
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

func runMongoMigrations() error {
	// Configurar variables de entorno
	setMongoEnv()

	// Verificar si ya existen colecciones (idempotencia)
	if hasMongoCollections() {
		fmt.Println("✅ MongoDB ya tiene colecciones - migraciones omitidas (idempotente)")
		return nil
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
