package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("=== Test de Conexión a Redis (Upstash) ===")
	fmt.Println()

	// Credenciales de Upstash
	redisURL := "redis://default:AaCrAAIncDJmMTFjYjJiOGU1M2U0YmM3YWIxMDQyZTA2ZjdlZDgxZXAyNDExMzE@living-wildcat-41131.upstash.io:6379"

	// Parsear URL de Redis
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("❌ Error parseando URL de Redis: %v\n", err)
	}

	// Habilitar TLS
	opt.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Crear cliente Redis
	client := redis.NewClient(opt)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Test PING
	fmt.Println("📡 Probando conexión con PING...")
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Error en PING: %v\n", err)
	}
	fmt.Printf("✅ PING exitoso: %s\n\n", pong)

	// 2. Test SET
	fmt.Println("📝 Probando SET...")
	testKey := "edugo:test:connection"
	testValue := fmt.Sprintf("Test desde EduGo - %s", time.Now().Format(time.RFC3339))
	err = client.Set(ctx, testKey, testValue, 5*time.Minute).Err()
	if err != nil {
		log.Fatalf("❌ Error en SET: %v\n", err)
	}
	fmt.Printf("✅ SET exitoso: %s = %s\n\n", testKey, testValue)

	// 3. Test GET
	fmt.Println("📖 Probando GET...")
	retrievedValue, err := client.Get(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("❌ Error en GET: %v\n", err)
	}
	fmt.Printf("✅ GET exitoso: %s\n\n", retrievedValue)

	// 4. Test TTL
	fmt.Println("⏰ Probando TTL...")
	ttl, err := client.TTL(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("❌ Error en TTL: %v\n", err)
	}
	fmt.Printf("✅ TTL: %.0f segundos\n\n", ttl.Seconds())

	// 5. Test INFO
	fmt.Println("📊 Obteniendo información del servidor...")
	info, err := client.Info(ctx, "server").Result()
	if err != nil {
		fmt.Printf("⚠️  Advertencia: No se pudo obtener INFO (normal en Upstash)\n\n")
	} else {
		fmt.Printf("✅ Información del servidor obtenida (%d bytes)\n\n", len(info))
	}

	// 6. Limpiar
	fmt.Println("🧹 Limpiando clave de prueba...")
	err = client.Del(ctx, testKey).Err()
	if err != nil {
		log.Fatalf("❌ Error en DEL: %v\n", err)
	}
	fmt.Printf("✅ Clave eliminada: %s\n\n", testKey)

	fmt.Println("✅ ¡Todas las pruebas de Redis pasaron exitosamente!")
	fmt.Println()
	fmt.Println("📋 Configuración de Redis para tu aplicación:")
	fmt.Println("   Host: living-wildcat-41131.upstash.io")
	fmt.Println("   Port: 6379")
	fmt.Println("   Password: AaCrAAIncDJmMTFjYjJiOGU1M2U0YmM3YWIxMDQyZTA2ZjdlZDgxZXAyNDExMzE")
	fmt.Println("   TLS: Habilitado")
}
