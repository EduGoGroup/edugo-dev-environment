# Migrator - EduGo Infrastructure v0.9.0

Este módulo contiene ejemplos de integración con la última versión de `edugo-infrastructure` (v0.9.0) utilizando tests de integración con testcontainers.

## 📦 Dependencias

El proyecto utiliza las siguientes dependencias de `edugo-infrastructure`:

```go
require (
    github.com/EduGoGroup/edugo-infrastructure/postgres v0.9.0
    github.com/EduGoGroup/edugo-infrastructure/mongodb v0.9.0
)
```

## 🧪 Tests de Integración

### PostgreSQL

El archivo `tests/postgres_integration_test.go` demuestra cómo:

1. **SetupSuite**: Configurar un contenedor PostgreSQL con testcontainers
2. **Aplicar migraciones**: Usar `migrations.ApplyAll(db)` del paquete infrastructure
3. **SetupTest**: Aplicar datos de prueba con `migrations.ApplyMockData(db)`
4. **TearDownTest**: Limpiar datos entre tests
5. **TearDownSuite**: Cerrar conexiones y detener contenedores

**Ejemplo de uso:**

```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"

func (s *Suite) SetupSuite() {
    // Conectar a PostgreSQL (testcontainer o local)
    db := conectarPostgres()
    
    // Aplicar migraciones
    migrations.ApplyAll(db)
}

func (s *Suite) SetupTest() {
    // Datos de prueba
    migrations.ApplyMockData(db)
}
```

### MongoDB

El archivo `tests/mongodb_integration_test.go` demuestra cómo:

1. **SetupSuite**: Configurar un contenedor MongoDB con testcontainers
2. **Aplicar migraciones**: Usar `migrations.ApplyAll(ctx, db)` del paquete infrastructure
3. **SetupTest**: Aplicar datos de prueba con `migrations.ApplyMockData(ctx, db)`
4. **TearDownTest**: Limpiar colecciones entre tests
5. **TearDownSuite**: Desconectar y detener contenedores

**Ejemplo de uso:**

```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrations"

func (s *Suite) SetupSuite() {
    // Conectar a MongoDB
    db := conectarMongo()
    
    // Aplicar migraciones
    migrations.ApplyAll(ctx, db)
}

func (s *Suite) SetupTest() {
    // Datos de prueba
    migrations.ApplyMockData(ctx, db)
}
```

## 🚀 Ejecutar Tests

```bash
# Compilar tests sin ejecutarlos
go test -c ./tests

# Ejecutar tests de PostgreSQL
go test -v ./tests -run TestPostgresIntegration

# Ejecutar tests de MongoDB
go test -v ./tests -run TestMongoDBIntegration

# Ejecutar todos los tests
go test -v ./tests
```

## 📝 Notas

- Los tests utilizan **testcontainers-go** para crear contenedores temporales
- Docker debe estar corriendo para ejecutar los tests
- Los contenedores se eliminan automáticamente después de los tests
- Cada test se ejecuta con datos frescos gracias a `SetupTest` y `TearDownTest`

## 🔗 Referencias

- [edugo-infrastructure](https://github.com/EduGoGroup/edugo-infrastructure)
- [postgres/USAGE_EXAMPLES.md](https://github.com/EduGoGroup/edugo-infrastructure/tree/main/postgres)
- [mongodb/USAGE_EXAMPLES.md](https://github.com/EduGoGroup/edugo-infrastructure/tree/main/mongodb)
- [testcontainers-go](https://golang.testcontainers.org/)

## ✅ Verificación

Para verificar que todo está configurado correctamente:

```bash
# Descargar dependencias
go mod download

# Limpiar y actualizar dependencias
go mod tidy

# Compilar proyecto
go build ./...

# Compilar tests
go test -c ./tests
```

Todos los comandos deben ejecutarse sin errores.
