package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	postgresMigrations "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
)

func main() {
	fmt.Println("=== EduGo Migrator - Neon Edition ===")
	fmt.Println("Migrando base de datos PostgreSQL a Neon...")
	fmt.Println()

	// Credenciales de Neon
	host := "ep-green-frost-ado4abbi-pooler.c-2.us-east-1.aws.neon.tech"
	port := "5432"
	user := "neondb_owner"
	password := "npg_sC2u9pTVwQJI"
	dbname := "edugo"
	sslmode := "require"

	// Verificar si se solicita migración forzada
	forceMigration := os.Getenv("FORCE_MIGRATION") == "true"
	if forceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO - Se eliminarán y recrearán todas las tablas")
		fmt.Println()
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Conectar a PostgreSQL en Neon
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Error conectando a Neon: %v\n", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("⚠️  Error cerrando conexión: %v", closeErr)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Error haciendo ping a Neon: %v\n", err)
	}

	fmt.Printf("✓ Conectado a Neon: %s@%s/%s\n", user, host, dbname)
	fmt.Println()

	// Si force=true, eliminar schema y recrear
	if forceMigration {
		fmt.Println("🔥 Eliminando schema público...")
		if err := dropPostgresSchema(db, user); err != nil {
			log.Fatalf("❌ Error eliminando schema: %v\n", err)
		}
		fmt.Println("✅ Schema eliminado exitosamente")
		fmt.Println()
	} else {
		// Verificar si ya existen tablas (idempotencia)
		if hasPostgresTables(db) {
			fmt.Println("✅ La base de datos ya tiene tablas - migraciones omitidas (idempotente)")
			fmt.Println("💡 Si deseas recrear la base de datos, ejecuta con: FORCE_MIGRATION=true")
			return
		}
	}

	// Aplicar todas las migraciones
	fmt.Println("📦 Aplicando migraciones de estructura y constraints...")
	if err := postgresMigrations.ApplyAll(db); err != nil {
		log.Fatalf("❌ Error aplicando migraciones: %v\n", err)
	}
	fmt.Println("✅ Migraciones de estructura completadas")
	fmt.Println()

	// Aplicar seeds (datos esenciales)
	fmt.Println("📦 Aplicando datos iniciales (seeds)...")
	if err := postgresMigrations.ApplySeeds(db); err != nil {
		log.Fatalf("❌ Error aplicando seeds: %v\n", err)
	}
	fmt.Println("✅ Datos iniciales aplicados")
	fmt.Println()

	// Aplicar datos de prueba/testing (opcional)
	applyMockData := os.Getenv("APPLY_MOCK_DATA") != "false" // Por defecto true
	if applyMockData {
		fmt.Println("📦 Aplicando datos de prueba (testing)...")
		if err := postgresMigrations.ApplyMockData(db); err != nil {
			log.Fatalf("❌ Error aplicando datos de prueba: %v\n", err)
		}
		fmt.Println("✅ Datos de prueba aplicados")
	} else {
		fmt.Println("⏭️  Saltando datos de prueba (APPLY_MOCK_DATA=false)")
	}
	fmt.Println()

	fmt.Println("✅ ¡Migración a Neon completada exitosamente!")
	fmt.Println()
	fmt.Println("📋 Cadena de conexión para tus aplicaciones:")
	fmt.Printf("   host=%s port=%s user=%s password=%s dbname=%s sslmode=%s\n",
		host, port, user, password, dbname, sslmode)
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

// dropPostgresSchema elimina y recrea el schema público en PostgreSQL
func dropPostgresSchema(db *sql.DB, user string) error {
	// Eliminar schema ui_config (Dynamic UI) antes del schema público
	_, err := db.Exec("DROP SCHEMA IF EXISTS ui_config CASCADE")
	if err != nil {
		return fmt.Errorf("error eliminando schema ui_config: %w", err)
	}

	// Eliminar schema público CASCADE
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
