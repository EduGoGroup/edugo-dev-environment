# Guía de Docker Compose Profiles

Documentación detallada de los profiles disponibles en el ambiente de desarrollo.

## ¿Qué son los Profiles?

Los profiles de Docker Compose permiten iniciar subconjuntos de servicios según necesidad, optimizando recursos y tiempo de inicio.

## Profiles Disponibles

### 1. `full` (Default)
**Servicios:** Todos  
**Uso:** Desarrollo completo del sistema

```bash
./scripts/setup.sh
# o explícitamente
./scripts/setup.sh --profile full
```

**Incluye:**
- PostgreSQL
- MongoDB  
- RabbitMQ
- API Mobile
- API Administración
- Worker

**Recursos:** ~2GB RAM, puertos 5432, 27017, 5672, 15672, 8081, 8082

---

### 2. `db-only`
**Servicios:** Solo infraestructura  
**Uso:** Testing de migraciones, desarrollo de schemas

```bash
./scripts/setup.sh --profile db-only
```

**Incluye:**
- PostgreSQL (5432)
- MongoDB (27017)
- RabbitMQ (5672, 15672)

**Recursos:** ~500MB RAM

---

### 3. `api-only`
**Servicios:** DBs + APIs (sin Worker)  
**Uso:** Desarrollo de endpoints, testing de APIs

```bash
./scripts/setup.sh --profile api-only
```

**Incluye:**
- PostgreSQL + MongoDB + RabbitMQ
- API Mobile (8081)
- API Administración (8082)

**Recursos:** ~1.5GB RAM

---

### 4. `mobile-only`
**Servicios:** DBs + API Mobile  
**Uso:** Desarrollo exclusivo de app móvil

```bash
./scripts/setup.sh --profile mobile-only
```

**Incluye:**
- PostgreSQL + MongoDB + RabbitMQ
- API Mobile (8081)

**Recursos:** ~1GB RAM

---

### 5. `admin-only`
**Servicios:** DBs + API Admin  
**Uso:** Desarrollo de panel administrativo

```bash
./scripts/setup.sh --profile admin-only
```

**Incluye:**
- PostgreSQL + MongoDB
- API Administración (8082)

**Recursos:** ~800MB RAM

---

### 6. `worker-only`
**Servicios:** DBs + Worker  
**Uso:** Testing de procesamiento asíncrono

```bash
./scripts/setup.sh --profile worker-only
```

**Incluye:**
- PostgreSQL + MongoDB + RabbitMQ
- Worker

**Recursos:** ~800MB RAM

---

## Casos de Uso

### Desarrollo de Features

```bash
# Trabajando en API Mobile
./scripts/setup.sh --profile mobile-only --seed

# Trabajando en Worker
./scripts/setup.sh --profile worker-only
```

### Testing

```bash
# Testing de migraciones
./scripts/setup.sh --profile db-only

# Testing end-to-end
./scripts/setup.sh --profile full --seed
```

### CI/CD

```bash
# Tests de integración
docker-compose --profile api-only up -d
```

## Gestión de Perfiles

### Iniciar
```bash
./scripts/setup.sh --profile <nombre>
```

### Detener
```bash
./scripts/stop.sh --profile <nombre>
```

### Ver servicios activos
```bash
cd docker
docker-compose --profile <nombre> ps
```

### Logs
```bash
cd docker
docker-compose --profile <nombre> logs -f
```

## Tips

1. **Desarrollo enfocado:** Usa el profile mínimo necesario
2. **Recursos:** `db-only` usa 75% menos memoria que `full`
3. **Inicio rápido:** Profiles específicos inician 50% más rápido
4. **Testing:** Combina profiles con `--seed` para datos de prueba

## Troubleshooting

**Error: "no configuration file provided"**
```bash
cd docker
docker-compose --profile <nombre> up -d
```

**Servicios no inician:**
```bash
# Verificar logs
docker-compose --profile <nombre> logs

# Reiniciar
./scripts/stop.sh --profile <nombre>
./scripts/setup.sh --profile <nombre>
```

**Puerto ocupado:**
Edita `docker/.env` para cambiar puertos.

---

*Última actualización: Noviembre 2025*
