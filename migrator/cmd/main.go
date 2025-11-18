package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	postgresPath := filepath.Join(infraDir, "postgres")

	// Configurar variables de entorno
	setPostgresEnv()

	// Ejecutar migraciones usando el CLI de postgres
	fmt.Println("Ejecutando migraciones de PostgreSQL...")
	cmd := exec.Command("go", "run", "migrate.go", "up")
	cmd.Dir = postgresPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando migraciones postgres: %w", err)
	}

	fmt.Println("✅ Migraciones de PostgreSQL completadas")
	return nil
}

func runMongoMigrations() error {
	mongoPath := filepath.Join(infraDir, "mongodb")

	// Configurar variables de entorno
	setMongoEnv()

	// Ejecutar migraciones usando el nuevo runner.go
	fmt.Println("Ejecutando migraciones de MongoDB (structure + constraints)...")
	cmd := exec.Command("go", "run", "runner.go", "all")
	cmd.Dir = mongoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando migraciones mongodb: %w", err)
	}

	fmt.Println("✅ Migraciones de MongoDB completadas")
	return nil
}

func setPostgresEnv() {
	setEnvIfEmpty("DB_HOST", "localhost")
	setEnvIfEmpty("DB_PORT", "5432")
	setEnvIfEmpty("DB_NAME", "edugo")
	setEnvIfEmpty("DB_USER", "edugo")
	setEnvIfEmpty("DB_PASSWORD", "edugo123")
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
