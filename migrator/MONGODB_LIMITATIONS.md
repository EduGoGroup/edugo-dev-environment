# Limitaciones de MongoDB en el Migrator

## 🔍 Problema Actual

Las migraciones de MongoDB **no se ejecutan automáticamente** en el migrator debido a limitaciones técnicas.

### Causa Raíz

Las migraciones de MongoDB en `edugo-infrastructure/mongodb/` están escritas como scripts JavaScript que requieren **mongosh** (MongoDB Shell):

```javascript
// mongodb/migrations/001_create_material_assessment.up.js
db.material_assessments.createIndex({ material_id: 1, created_at: -1 })
```

Estos scripts se ejecutan con:
```bash
mongosh < migration.js
```

### Por Qué No Funciona en Docker

1. **mongosh es un binario x64** - No funciona en Apple Silicon (ARM)
2. **Incompatibilidad con Alpine Linux** - Requiere dependencias glibc que Alpine no tiene
3. **Imagen pesada** - Instalar la imagen completa de MongoDB (800MB+) solo por mongosh no es óptimo

## ✅ PostgreSQL Funciona Perfectamente

- ✅ Migraciones SQL ejecutadas automáticamente
- ✅ 11 migraciones aplicadas correctamente
- ✅ Todas las tablas creadas

## 🔧 Soluciones Posibles

### Opción 1: Migrar MongoDB manualmente (ACTUAL)

```bash
# Entrar al contenedor de MongoDB
docker compose exec mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin

# Ejecutar migraciones manualmente
use edugo
db.material_assessments.createIndex(...)
```

### Opción 2: Modificar edugo-infrastructure (RECOMENDADO)

Cambiar las migraciones de MongoDB para usar el driver de Go en vez de scripts JavaScript:

```go
// En vez de archivos .js, usar código Go
func migration001Up(db *mongo.Database) error {
    _, err := db.Collection("material_assessments").Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
            Keys: bson.D{
                {Key: "material_id", Value: 1},
                {Key: "created_at", Value: -1},
            },
        },
    )
    return err
}
```

### Opción 3: Usar imagen base de MongoDB

Cambiar el Dockerfile para usar `FROM mongo:7.0` en vez de Alpine, pero la imagen sería mucho más pesada (800MB vs 400MB).

## 📊 Estado Actual

| Base de Datos | Migraciones Automáticas | Estado |
|---------------|------------------------|--------|
| PostgreSQL | ✅ Funciona | 11 migraciones aplicadas |
| MongoDB | ❌ Manual | Requiere mongosh |

## 💡 Recomendación

**Para desarrollo local:**
- PostgreSQL: ✅ Automático con migrator
- MongoDB: Ejecutar migraciones manualmente cuando sea necesario

**Para producción:**
- Modificar `edugo-infrastructure` para usar drivers de Go en vez de mongosh
- Esto permitirá migraciones automáticas para ambas bases de datos

## 🎯 Impacto

**Los servicios funcionan correctamente** porque:
- MongoDB está corriendo y accesible
- Las colecciones se crean dinámicamente cuando se insertan datos
- Los índices no son críticos para el funcionamiento básico
- Solo afecta el rendimiento en queries grandes
