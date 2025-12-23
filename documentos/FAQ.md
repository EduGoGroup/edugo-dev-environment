# FAQ - Preguntas Frecuentes

## Preguntas Generales

### ¿Necesito saber Go/Backend para usar esto?

**No.** Solo ejecuta el setup y usa las APIs desde tu frontend. El backend ya está compilado y empaquetado en imágenes Docker.

### ¿Los datos se pierden al detener Docker?

**No.** Los datos persisten en volúmenes Docker. Solo se borran si ejecutas:
```bash
docker-compose down -v  # La flag -v elimina volúmenes
```

### ¿Puedo trabajar offline?

**Sí**, después del primer setup. Las imágenes Docker quedan almacenadas localmente en tu Mac.

### ¿Cómo actualizo las APIs a la última versión?

```bash
./scripts/update-images.sh
docker-compose up -d  # Reiniciar con nuevas imágenes
```

### ¿Puedo cambiar los puertos?

**Sí.** Edita `docker/.env`:
```bash
API_MOBILE_PORT=8083  # Cambiar puerto de API Mobile
API_ADMIN_PORT=8084   # Cambiar puerto de API Admin
```

### ¿Cómo cargo datos de prueba?

```bash
# Opción 1: Durante el setup
./scripts/setup.sh --seed

# Opción 2: Después del setup
./scripts/seed-data.sh

# Opción 3: Usando make
make seed
```

### ¿Cuánto tarda el setup en verificar que todo está listo?

El script espera automáticamente a que PostgreSQL, MongoDB y RabbitMQ estén saludables. Por defecto espera hasta 120 segundos. Puedes cambiar esto:

```bash
./scripts/setup.sh --timeout 180  # Esperar hasta 3 minutos
```

---

## Troubleshooting

### Error: "Cannot connect to Docker daemon"

**Causa:** Docker Desktop no está corriendo.

**Solución:**
```bash
open -a Docker  # macOS
# Esperar a que el ícono aparezca en la barra de menú
docker ps       # Verificar
```

### Error: "Port already in use"

**Causa:** Otro proceso usa el puerto.

**Solución:**
```bash
# Ver qué usa el puerto
lsof -ti:8081

# Matar el proceso
lsof -ti:8081 | xargs kill -9

# O cambiar puerto en docker/.env
```

### Error: "pull access denied for ghcr.io/edugogroup/..."

**Causa:** No estás autenticado en GitHub Container Registry.

**Solución:**
```bash
# 1. Crear token en: https://github.com/settings/tokens
#    Scope requerido: read:packages

# 2. Login
docker login ghcr.io
Username: tu-usuario-github
Password: ghp_tu_token_aqui
```

### Error: "dial tcp: connection refused" en API

**Causa:** PostgreSQL o MongoDB no están corriendo o saludables.

**Solución:**
```bash
# Ver estado
docker-compose ps

# Si postgres/mongodb están "Exited", reiniciar
docker-compose up -d postgres mongodb

# Ver logs de error
docker-compose logs postgres mongodb
```

### Error: "relation 'users' does not exist"

**Causa:** Migraciones no se ejecutaron.

**Solución:**
```bash
# Ejecutar migrator manualmente
docker-compose up migrator

# Verificar tablas
docker exec -it edugo-postgres psql -U edugo -d edugo -c "\dt"
```

### Worker no procesa PDFs

**Causa:** Falta OPENAI_API_KEY o RabbitMQ no conecta.

**Solución:**
```bash
# 1. Verificar logs
docker-compose logs worker

# 2. Verificar API key en .env
grep OPENAI_API_KEY docker/.env

# 3. Si falta, agregar
echo "OPENAI_API_KEY=sk-tu-key-real" >> docker/.env

# 4. Reiniciar worker
docker-compose restart worker
```

### Error: "no space left on device"

**Causa:** Docker llenó el disco.

**Solución:**
```bash
# Ver uso
docker system df

# Limpiar imágenes sin usar
docker image prune -a

# Limpiar todo (cuidado)
docker system prune -a --volumes
```

### API responde 500 Internal Server Error

**Causa:** Error interno o base de datos no disponible.

**Solución:**
```bash
# 1. Ver logs detallados
docker-compose logs api-mobile

# 2. Verificar que DBs están healthy
docker-compose ps

# 3. Reiniciar todo
docker-compose restart
```

### Los contenedores se reinician constantemente

**Causa:** Health checks fallan o error en configuración.

**Solución:**
```bash
# Ver logs de arranque
docker-compose logs --tail=50 api-mobile

# Verificar variables de entorno
docker-compose exec api-mobile env | grep DATABASE
```

---

## Preguntas de Desarrollo

### ¿Cómo genero un hash de contraseña para nuevos usuarios?

```bash
./scripts/generate-password.sh mi-password
```

### ¿Cómo veo las colas de RabbitMQ?

Abre http://localhost:15672
- Usuario: `edugo`
- Contraseña: `edugo123`

### ¿Cómo ejecuto solo las bases de datos?

```bash
./scripts/setup.sh --profile db-only
# O manualmente:
docker-compose up -d postgres mongodb rabbitmq
```

### ¿Cómo levanto las APIs opcionales (Admin, Worker)?

Las APIs adicionales usan profiles en Docker Compose:

```bash
# Solo API Mobile (default)
cd docker && docker-compose -f docker-compose-apps.yml up -d

# API Mobile + API Admin
docker-compose -f docker-compose-apps.yml --profile with-admin up -d

# API Mobile + Worker
docker-compose -f docker-compose-apps.yml --profile with-worker up -d

# Todos los servicios
docker-compose -f docker-compose-apps.yml --profile full up -d
```

**Nota:** API Admin y Worker requieren configuración adicional (ver docker-compose-apps.yml).

### ¿Cómo cargo datos de prueba adicionales?

```bash
./scripts/seed-data.sh
```

### ¿Cómo veo los documentos en MongoDB?

```bash
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 edugo --authSource admin

# Ver colecciones
show collections

# Ver documentos
db.documents.find().pretty()
```

---

## Preguntas de Configuración

### ¿Cuáles son las credenciales por defecto?

| Servicio | Usuario | Contraseña |
|----------|---------|------------|
| PostgreSQL | edugo | edugo123 |
| MongoDB | edugo | edugo123 |
| RabbitMQ | edugo | edugo123 |
| Usuarios de prueba | * | edugo2024 |

### ¿Dónde está el archivo de configuración?

```
docker/.env          # Configuración activa
docker/.env.example  # Plantilla con todos los valores
```

### ¿Qué versiones de imágenes se usan?

Por defecto `latest`. Para usar versiones específicas, edita `.env`:
```bash
API_MOBILE_VERSION=v1.2.3
API_ADMIN_VERSION=v1.2.3
WORKER_VERSION=v1.2.3
```

### ¿Cómo deshabilito S3?

S3 ya está marcado como opcional por defecto:
```bash
BOOTSTRAP_OPTIONAL_RESOURCES_S3=true  # Ya está en .env
```

---

## Reportar Problemas

Si ninguna solución funciona:

1. Captura los logs: `docker-compose logs > logs.txt`
2. Captura el estado: `docker-compose ps > status.txt`
3. Crea un issue en GitHub con ambos archivos

---

**Ver también:** [GUIA-RAPIDA.md](./GUIA-RAPIDA.md) para instrucciones de setup.
