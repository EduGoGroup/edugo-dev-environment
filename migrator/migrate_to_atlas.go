package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoMigrations "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"
)

func main() {
	fmt.Println("=== EduGo Migrator - MongoDB Atlas Edition ===")
	fmt.Println("Migrando base de datos MongoDB a Atlas...")
	fmt.Println()

	// Credenciales de MongoDB Atlas
	mongoURI := "mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/?appName=Edugo"
	dbName := "edugo"

	// Verificar si se solicita migración forzada
	forceMigration := os.Getenv("FORCE_MIGRATION") == "true"
	if forceMigration {
		fmt.Println("⚠️  MODO FORZADO ACTIVADO - Se eliminará y recreará toda la base de datos")
		fmt.Println()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Conectar a MongoDB Atlas
	fmt.Println("📡 Conectando a MongoDB Atlas...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Error conectando a MongoDB Atlas: %v\n", err)
	}
	defer func() {
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("⚠️  Error desconectando: %v", disconnectErr)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Error haciendo ping a MongoDB Atlas: %v\n", err)
	}

	fmt.Printf("✓ Conectado a MongoDB Atlas: %s\n", dbName)
	fmt.Println()

	db := client.Database(dbName)

	// Si force=true, eliminar base de datos completa
	if forceMigration {
		fmt.Println("🔥 Eliminando base de datos MongoDB...")
		if err := db.Drop(ctx); err != nil {
			log.Fatalf("❌ Error eliminando base de datos: %v\n", err)
		}
		fmt.Println("✅ Base de datos eliminada exitosamente")
		fmt.Println()
	} else {
		// Verificar si ya existen colecciones (idempotencia)
		if hasMongoCollections(ctx, db) {
			fmt.Println("✅ MongoDB ya tiene colecciones - migraciones omitidas (idempotente)")
			fmt.Println("💡 Si deseas recrear la base de datos, ejecuta con: FORCE_MIGRATION=true")
			return
		}
	}

	// Aplicar todas las migraciones
	fmt.Println("📦 Aplicando migraciones de estructura y constraints...")
	if err := mongoMigrations.ApplyAll(ctx, db); err != nil {
		log.Fatalf("❌ Error aplicando migraciones: %v\n", err)
	}
	fmt.Println("✅ Migraciones de estructura completadas")
	fmt.Println()

	// Aplicar datos de prueba/testing (opcional)
	applyMockData := os.Getenv("APPLY_MOCK_DATA") != "false" // Por defecto true
	if applyMockData {
		fmt.Println("📦 Aplicando datos de prueba (testing)...")
		if err := mongoMigrations.ApplyMockData(ctx, db); err != nil {
			log.Fatalf("❌ Error aplicando datos de prueba: %v\n", err)
		}
		fmt.Println("✅ Datos de prueba aplicados")
	} else {
		fmt.Println("⏭️  Saltando datos de prueba (APPLY_MOCK_DATA=false)")
	}
	fmt.Println()

	// Verificar colecciones creadas
	fmt.Println("📂 Verificando colecciones creadas...")
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Error listando colecciones: %v\n", err)
	}
	fmt.Printf("✅ Colecciones creadas (%d): %v\n", len(collections), collections)
	fmt.Println()

	fmt.Println("✅ ¡Migración a MongoDB Atlas completada exitosamente!")
	fmt.Println()
	fmt.Println("📋 URI de conexión para tus aplicaciones:")
	fmt.Printf("   mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/%s?appName=Edugo\n", dbName)
}

// hasMongoCollections verifica si MongoDB ya tiene colecciones creadas
func hasMongoCollections(ctx context.Context, db *mongo.Database) bool {
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("⚠️  Error listando colecciones: %v\n", err)
		return false
	}

	return len(collections) > 0
}
