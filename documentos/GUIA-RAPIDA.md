# Guía Rápida - EduGo Dev Environment

## Prerequisitos

1. **Docker Desktop** instalado y corriendo
2. **Acceso a GitHub Container Registry** (token con scope `read:packages`)

### Instalar Docker Desktop

```bash
# macOS
brew install --cask docker

# O descarga desde: https://www.docker.com/products/docker-desktop
```

### Verificar Docker

```bash
docker --version
docker info  # Debe responder sin errores
```

---

## Setup Rápido (3 pasos)

### Paso 1: Clonar

```bash
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment
```

### Paso 2: Ejecutar Setup

```bash
./scripts/setup.sh
```

El script:
1. Verifica Docker
2. Solicita tu GitHub Personal Access Token
3. Descarga imágenes Docker
4. Crea archivo `.env`
5. Levanta todos los servicios
6. **Espera automáticamente** a que PostgreSQL, MongoDB y RabbitMQ estén saludables
7. Ejecuta migraciones

#### Opciones de Setup

```bash
# Setup básico
./scripts/setup.sh

# Setup con datos de prueba
./scripts/setup.sh --seed
./scripts/setup.sh -s

# Setup solo bases de datos
./scripts/setup.sh --profile db-only

# Setup con timeout personalizado para health checks
./scripts/setup.sh --timeout 180
```

| Opción | Descripción |
|--------|-------------|
| `-s, --seed` | Cargar datos de prueba después de iniciar |
| `-p, --profile <nombre>` | Usar perfil específico (full, db-only, api-only, etc.) |
| `-t, --timeout <segundos>` | Timeout para health checks (default: 120) |
| `-h, --help` | Mostrar ayuda |

### Paso 3: Verificar

```bash
# Ver servicios corriendo
cd docker && docker-compose ps

# Probar API Mobile
curl http://localhost:8081/health

# Probar API Admin
curl http://localhost:8082/health
```

---

## URLs de Servicios

| Servicio | URL |
|----------|-----|
| API Mobile | http://localhost:8081 |
| API Mobile Swagger | http://localhost:8081/swagger/index.html |
| API Admin | http://localhost:8082 |
| API Admin Swagger | http://localhost:8082/swagger/index.html |
| RabbitMQ UI | http://localhost:15672 (edugo/edugo123) |

---

## Usuarios de Prueba

**Contraseña para TODOS:** `edugo2024`

| Email | Rol |
|-------|-----|
| admin@edugo.test | Administrador |
| teacher.math@edugo.test | Profesor |
| teacher.science@edugo.test | Profesor |
| student1@edugo.test | Estudiante |
| student2@edugo.test | Estudiante |
| student3@edugo.test | Estudiante |
| guardian1@edugo.test | Tutor |
| guardian2@edugo.test | Tutor |

### Ejemplo de Login

```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "student1@edugo.test", "password": "edugo2024"}'
```

---

## Comandos con Makefile (Recomendado)

```bash
make help        # Ver todos los comandos disponibles

make up          # Iniciar servicios
make down        # Detener servicios
make restart     # Reiniciar servicios
make status      # Ver estado

make logs        # Ver todos los logs
make logs-api    # Ver logs de API Mobile
make logs-admin  # Ver logs de API Admin

make diagnose    # Ejecutar diagnóstico completo
make health      # Verificar health de APIs

make psql        # Conectar a PostgreSQL
make mongo       # Conectar a MongoDB

make reset       # Reset completo (borra datos)
```

## Comandos Docker Compose (Alternativa)

```bash
# Iniciar
cd docker && docker-compose up -d

# Detener
docker-compose stop

# Ver logs
docker-compose logs -f api-mobile

# Reset completo
docker-compose down -v && cd .. && ./scripts/setup.sh
```

---

## Conectar Frontend

### React / Next.js

```javascript
const API_BASE = 'http://localhost:8081';

// Login
const response = await fetch(`${API_BASE}/v1/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
});
const { token, user } = await response.json();

// Usar token
const courses = await fetch(`${API_BASE}/api/v1/courses`, {
  headers: { Authorization: `Bearer ${token}` }
});
```

### Vue.js / Axios

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8081'
});

// Login
const { data } = await api.post('/v1/auth/login', { email, password });
localStorage.setItem('token', data.token);

// Configurar interceptor
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});
```

---

## Acceso Directo a Bases de Datos

### PostgreSQL

```bash
docker exec -it edugo-postgres psql -U edugo -d edugo

# Comandos útiles
\dt           # Listar tablas
\d users      # Describir tabla
SELECT * FROM users LIMIT 5;
\q            # Salir
```

### MongoDB

```bash
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 edugo --authSource admin

# Comandos útiles
show collections
db.documents.find().limit(5)
exit
```

---

## Scripts Disponibles

| Script | Propósito |
|--------|-----------|
| `./scripts/setup.sh` | Setup inicial completo |
| `./scripts/validate.sh` | Validar docker-compose |
| `./scripts/stop.sh` | Detener servicios |
| `./scripts/cleanup.sh` | Limpiar ambiente (interactivo) |
| `./scripts/update-images.sh` | Actualizar imágenes Docker |
| `./scripts/seed-data.sh` | Cargar datos de prueba |

### Perfiles de Setup

```bash
./scripts/setup.sh                      # Todo (default)
./scripts/setup.sh --profile db-only    # Solo bases de datos
./scripts/setup.sh --profile api-only   # DBs + APIs
./scripts/setup.sh -s                   # Con seed data
```

---

## Troubleshooting Rápido

### "Cannot connect to Docker daemon"
```bash
open -a Docker  # Abrir Docker Desktop
```

### "Port already in use"
```bash
lsof -ti:8081 | xargs kill -9  # Liberar puerto
```

### "pull access denied"
```bash
docker login ghcr.io
# Username: tu-usuario-github
# Password: tu-personal-access-token (ghp_...)
```

### "API responde 500"
```bash
docker-compose logs api-mobile  # Ver error real
```

---

**Ver también:** [FAQ.md](./FAQ.md) para más soluciones.
