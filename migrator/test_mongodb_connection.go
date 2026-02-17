package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("=== Test de Conexión a MongoDB (Atlas) ===")
	fmt.Println()

	// URI de MongoDB Atlas
	mongoURI := "mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/?appName=Edugo"
	dbName := "edugo"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Conectar a MongoDB
	fmt.Println("📡 Conectando a MongoDB Atlas...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Error conectando a MongoDB: %v\n", err)
	}
	defer func() {
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("⚠️  Error desconectando: %v", disconnectErr)
		}
	}()

	// 2. Test PING
	fmt.Println("📡 Probando conexión con PING...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Error en PING: %v\n", err)
	}
	fmt.Println("✅ PING exitoso")
	fmt.Println()

	// 3. Obtener información del servidor
	fmt.Println("📊 Obteniendo información del servidor...")
	var serverStatus bson.M
	err = client.Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&serverStatus)
	if err != nil {
		fmt.Printf("⚠️  Advertencia: No se pudo obtener serverStatus: %v\n\n", err)
	} else {
		if version, ok := serverStatus["version"].(string); ok {
			fmt.Printf("✅ Versión de MongoDB: %s\n", version)
		}
		if host, ok := serverStatus["host"].(string); ok {
			fmt.Printf("✅ Host: %s\n", host)
		}
		fmt.Println()
	}

	// 4. Listar bases de datos
	fmt.Println("📚 Listando bases de datos...")
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Error listando bases de datos: %v\n", err)
	}
	fmt.Printf("✅ Bases de datos encontradas: %v\n", databases)
	fmt.Println()

	// 5. Acceder a la base de datos edugo
	db := client.Database(dbName)
	fmt.Printf("📦 Accediendo a base de datos: %s\n", dbName)

	// 6. Listar colecciones
	fmt.Println("📂 Listando colecciones...")
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Error listando colecciones: %v\n", err)
	}

	if len(collections) == 0 {
		fmt.Println("⚠️  No hay colecciones (base de datos vacía)")
		fmt.Println("💡 Ejecuta las migraciones para crear las colecciones")
	} else {
		fmt.Printf("✅ Colecciones encontradas (%d): %v\n", len(collections), collections)
	}
	fmt.Println()

	// 7. Test de escritura/lectura
	fmt.Println("📝 Probando escritura/lectura en colección de prueba...")
	testCollection := db.Collection("_test_connection")

	// Insertar documento de prueba
	testDoc := bson.M{
		"message":   "Test desde EduGo",
		"timestamp": time.Now(),
	}
	insertResult, err := testCollection.InsertOne(ctx, testDoc)
	if err != nil {
		log.Fatalf("❌ Error insertando documento: %v\n", err)
	}
	fmt.Printf("✅ Documento insertado con ID: %v\n", insertResult.InsertedID)

	// Leer documento
	var retrievedDoc bson.M
	err = testCollection.FindOne(ctx, bson.M{"_id": insertResult.InsertedID}).Decode(&retrievedDoc)
	if err != nil {
		log.Fatalf("❌ Error leyendo documento: %v\n", err)
	}
	fmt.Printf("✅ Documento recuperado: %s\n", retrievedDoc["message"])

	// Eliminar documento de prueba
	_, err = testCollection.DeleteOne(ctx, bson.M{"_id": insertResult.InsertedID})
	if err != nil {
		log.Fatalf("❌ Error eliminando documento: %v\n", err)
	}
	fmt.Println("✅ Documento de prueba eliminado")
	fmt.Println()

	// 8. Eliminar colección de prueba
	err = testCollection.Drop(ctx)
	if err != nil {
		log.Fatalf("❌ Error eliminando colección de prueba: %v\n", err)
	}
	fmt.Println("✅ Colección de prueba eliminada")
	fmt.Println()

	fmt.Println("✅ ¡Todas las pruebas de MongoDB pasaron exitosamente!")
	fmt.Println()
	fmt.Println("📋 Configuración de MongoDB para tu aplicación:")
	fmt.Println("   Cluster: edugo.alxme5j.mongodb.net")
	fmt.Println("   Base de datos: edugo")
	fmt.Println("   Usuario: medinatello_db_user")
	fmt.Println("   URI: mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/?appName=Edugo")
}
