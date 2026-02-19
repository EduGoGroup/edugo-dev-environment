# EduGo - Modo Cloud

## Qu茅 es el Modo Cloud

El modo cloud levanta **solo las APIs** conect谩ndose a servicios en la nube:

| Servicio | Ubicaci贸n | Configuraci贸n |
|----------|-----------|---------------|
| PostgreSQL | Neon | Via `.env.cloud` |
| MongoDB | Atlas | Via `.env.cloud` |
| Redis | Upstash | Via `.env.cloud` |
| RabbitMQ | Docker local | Opcional |

## Requisitos

1. Configurar `docker/.env.cloud` con las credenciales de tus servicios cloud:
   ```bash
   cp .env.cloud.example .env.cloud
   # Editar .env.cloud con valores reales
   ```

2. Ejecutar migraciones en cloud:
   ```bash
   make db-migrate-cloud
   ```

## Uso

### Levantar todo (APIs + RabbitMQ)

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud --profile full up -d
```

### Solo API Mobile

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud --profile apps up -d
```

### Solo API Admin

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud --profile admin up -d
```

### Solo Worker

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud --profile worker up -d
```

## Comparaci贸n de Modos

| Aspecto | Modo Docker | Modo Cloud |
|---------|-------------|------------|
| PostgreSQL | Contenedor local | Neon |
| MongoDB | Contenedor local | Atlas |
| Redis | Contenedor local | Upstash |
| RabbitMQ | Contenedor local | Contenedor local (opcional) |
| Tiempo inicio | ~30-60 segundos | ~5-10 segundos |
| Persistencia | Se pierde con `down -v` | Siempre persistente |

## Comandos 脷tiles

```bash
# Ver logs
docker logs -f edugo-api-mobile-cloud
docker logs -f edugo-api-administracion-cloud
docker logs -f edugo-worker-cloud

# Detener todo
cd docker
docker-compose -f docker-compose.cloud.yml down

# Detener y eliminar vol煤menes (RabbitMQ)
cd docker
docker-compose -f docker-compose.cloud.yml down -v
```

## Cambiar entre Modos

```bash
# De Cloud a Docker Tradicional
cd docker
docker-compose -f docker-compose.cloud.yml down
docker-compose up -d

# De Docker Tradicional a Cloud
cd docker
docker-compose down
docker-compose -f docker-compose.cloud.yml --env-file .env.cloud --profile full up -d
```

## Troubleshooting

### Error de conexi贸n a PostgreSQL
Verifica que `POSTGRES_HOST` y `POSTGRES_SSLMODE=require` est茅n configurados en `.env.cloud`.

### Error de conexi贸n a MongoDB
Verifica que `MONGODB_URI` est茅 configurado correctamente en `.env.cloud`.

### API no se conecta a RabbitMQ
Si no usas RabbitMQ, configura `BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=false`.

## Documentaci贸n Relacionada

- `CLOUD_SETUP.md` - Gu铆a completa de configuraci贸n cloud
- `docker-compose.yml` - Modo tradicional (contenedores locales)
- `docker-compose.cloud.yml` - Modo cloud
- `.env.cloud.example` - Plantilla de variables de entorno
