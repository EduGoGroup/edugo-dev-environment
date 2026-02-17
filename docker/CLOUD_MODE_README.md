# 🚀 EduGo - Modo Cloud

## ✅ ¿Qué es el Modo Cloud?

El modo cloud levanta **solo las APIs** conectándose a servicios en la nube:

| Servicio | Ubicación | Configuración |
|----------|-----------|---------------|
| PostgreSQL | ☁️ **Neon** | Ya configurado |
| MongoDB | ☁️ **Atlas** | Ya configurado |
| Redis | ☁️ **Upstash** | Ya configurado |
| RabbitMQ | 🐳 Local (opcional) | Docker local |

## 🎯 Ventajas del Modo Cloud

- ✅ **Inicio rápido**: No esperar a que levanten PostgreSQL, MongoDB, Redis
- ✅ **Recursos ligeros**: Solo levanta las APIs necesarias
- ✅ **Datos persistentes**: Bases de datos siempre disponibles
- ✅ **Colaboración**: Todo el equipo comparte las mismas bases de datos
- ✅ **Desarrollo remoto**: Funciona desde cualquier lugar

## 📋 Uso

### Opción 1: Levantar todo (APIs + RabbitMQ)

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --profile full up -d
```

### Opción 2: Solo API Mobile

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --profile apps up -d
```

### Opción 3: Solo API Admin

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --profile admin up -d
```

### Opción 4: Solo Worker

```bash
cd docker
docker-compose -f docker-compose.cloud.yml --profile worker up -d
```

### Opción 5: APIs sin RabbitMQ (más ligero)

```bash
cd docker
docker-compose -f docker-compose.cloud.yml up -d api-mobile api-administracion worker
```

## 🔧 Configuración en Zed Editor

### Para Desarrollo Local (sin Docker)

Cada proyecto ahora tiene una nueva configuración de debug:

**API Mobile:**
- `Go: Debug main (CLOUD MODE - Neon/Atlas/Upstash)`

**API Administración:**
- `Go: Debug main (CLOUD MODE - Neon/Atlas/Upstash)`

**Worker:**
- `Go: Debug main (CLOUD MODE - Neon/Atlas/Upstash)`

### Cómo Usarlas

1. Abre el proyecto en Zed
2. Ve a la paleta de comandos (Cmd+Shift+P)
3. Busca "Debug: Select Configuration"
4. Selecciona la opción **"CLOUD MODE"**
5. Inicia el debug normalmente

## 📊 Comparación de Modos

| Aspecto | Modo Docker (tradicional) | Modo Cloud (nuevo) |
|---------|---------------------------|-------------------|
| **PostgreSQL** | 🐳 Contenedor local | ☁️ Neon |
| **MongoDB** | 🐳 Contenedor local | ☁️ Atlas |
| **Redis** | 🐳 Contenedor local | ☁️ Upstash |
| **RabbitMQ** | 🐳 Contenedor local | 🐳 Contenedor local (opcional) |
| **Tiempo inicio** | ~30-60 segundos | ~5-10 segundos |
| **Memoria RAM** | ~2-3 GB | ~500 MB |
| **Persistencia** | Se pierde con `down -v` | Siempre persistente |

## 🛠️ Comandos Útiles

### Ver logs de las APIs

```bash
# API Mobile
docker logs -f edugo-api-mobile-cloud

# API Admin
docker logs -f edugo-api-administracion-cloud

# Worker
docker logs -f edugo-worker-cloud
```

### Detener todo

```bash
cd docker
docker-compose -f docker-compose.cloud.yml down
```

### Detener y eliminar volúmenes (RabbitMQ)

```bash
cd docker
docker-compose -f docker-compose.cloud.yml down -v
```

## 🔄 Cambiar entre Modos

### De Cloud a Docker Tradicional

```bash
cd docker

# Detener modo cloud
docker-compose -f docker-compose.cloud.yml down

# Iniciar modo tradicional
docker-compose up -d
```

### De Docker Tradicional a Cloud

```bash
cd docker

# Detener modo tradicional
docker-compose down

# Iniciar modo cloud
docker-compose -f docker-compose.cloud.yml --profile full up -d
```

## ⚠️ Notas Importantes

1. **RabbitMQ es opcional**: Si tus APIs no usan mensajería, no necesitas levantarlo
2. **Datos compartidos**: Todos los desarrolladores comparten las mismas bases de datos cloud
3. **Límites gratuitos**: Revisa los límites en `CLOUD_SETUP.md`
4. **Variables de entorno**: Se cargan desde `.env` o `.env.neon`

## 🆘 Troubleshooting

### Error de conexión a PostgreSQL

Verifica que el host sea correcto:
```
POSTGRES_HOST=ep-green-frost-ado4abbi-pooler.c-2.us-east-1.aws.neon.tech
POSTGRES_SSLMODE=require
```

### Error de conexión a MongoDB

Verifica la URI de MongoDB Atlas:
```
MONGODB_URI=mongodb+srv://medinatello_db_user:6NQjJDaOkN4nvldT@edugo.alxme5j.mongodb.net/?appName=Edugo
```

### API no se conecta a RabbitMQ

Si no usas RabbitMQ, configura:
```
BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=false
```

## 📚 Documentación Relacionada

- `CLOUD_SETUP.md` - Guía completa de configuración cloud
- `docker-compose.yml` - Modo tradicional (con contenedores locales)
- `docker-compose.cloud.yml` - Modo cloud (este archivo)
- `.env.neon` - Variables de entorno para cloud

## 💡 Recomendaciones

**Para Desarrolladores Frontend:**
```bash
# Solo levanta las APIs que necesites
docker-compose -f docker-compose.cloud.yml up -d api-mobile api-administracion
```

**Para Desarrolladores Backend:**
```bash
# Usa el modo de debug de Zed con "CLOUD MODE"
# No necesitas Docker, todo corre localmente conectándose a cloud
```

**Para Testing Completo:**
```bash
# Levanta todo incluyendo RabbitMQ
docker-compose -f docker-compose.cloud.yml --profile full up -d
```
