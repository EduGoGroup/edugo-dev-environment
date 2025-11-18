# EXECUTION PLAN - Dev Environment

## Información del Proyecto

**Proyecto:** EduGo Dev Environment  
**Objetivo:** Infraestructura Docker Compose para desarrollo local completo  
**Duración:** 3 Sprints (6 semanas)  
**Equipo:** 1 DevOps engineer  
**Repositorio:** https://github.com/EduGoGroup/edugo-dev-environment

---

## Fase 1: Docker Compose Stack (Sprint 1)

### 1.1 Servicios Base
- [ ] Crear docker-compose.yml v3.8
- [ ] PostgreSQL 15 Alpine
- [ ] MongoDB 7.0
- [ ] RabbitMQ 3.12 Management
- [ ] Redis 7 (opcional)

### 1.2 Configuración
- [ ] Volúmenes para persistencia
- [ ] Networks (bridge)
- [ ] Environment variables
- [ ] Health checks para cada servicio

### 1.3 Scripts Básicos
- [ ] setup.sh - Iniciar ambiente
- [ ] teardown.sh - Detener servicios
- [ ] clean.sh - Limpiar volúmenes

### 1.4 Documentación
- [ ] README.md
- [ ] SETUP.md (guía de inicio)
- [ ] SERVICES.md (detalle de cada servicio)

**Entregables:**
- [ ] docker-compose.yml funcional
- [ ] Todos los servicios levantable con: docker-compose up
- [ ] Documentación de uso

---

## Fase 2: Integración de APIs + Seed Data (Sprint 2)

### 2.1 APIs en Docker Compose
- [ ] Agregar api-mobile (puerto 8080)
- [ ] Agregar api-admin (puerto 8081)
- [ ] Agregar worker (sin puerto, backend)
- [ ] Configurar environment vars
- [ ] Configurar depends_on

### 2.2 Seed Data
- [ ] docker/postgres/init.sql
  - [ ] Crear todas las tablas
  - [ ] Crear índices
  - [ ] Insertar datos de ejemplo
  
- [ ] docker/postgres/seed.sql
  - [ ] Escuela de prueba
  - [ ] Facultad y departamento
  - [ ] Docentes y estudiantes
  - [ ] Usuarios de prueba
  
- [ ] docker/mongo/init.js
  - [ ] Crear colecciones
  - [ ] Crear índices
  - [ ] Insertar documentos ejemplo
  
- [ ] docker/rabbitmq/rabbitmq.conf
  - [ ] Declarar exchanges
  - [ ] Declarar queues
  - [ ] Configurar bindings

### 2.3 Scripts de Seed
- [ ] seed-data.sh - Cargar datos de prueba
- [ ] reset-data.sh - Resetear a estado inicial
- [ ] export-data.sh - Exportar datos

### 2.4 Health Checks
- [ ] health-check.sh
  - [ ] Verificar PostgreSQL
  - [ ] Verificar MongoDB
  - [ ] Verificar RabbitMQ
  - [ ] Verificar APIs (health endpoints)
  - [ ] Reporte de estado general

**Entregables:**
- [ ] APIs levantadas y funcionando
- [ ] Datos de prueba cargados
- [ ] Health checks implementados
- [ ] Scripts de operación

---

## Fase 3: Perfil Dev + Documentación Completa (Sprint 3)

### 3.1 Docker Compose Override para Desarrollo
- [ ] docker-compose.override.yml
- [ ] Volúmenes de código (hot reload)
- [ ] Comandos para desarrollo (go run)
- [ ] Logs en stdout

### 3.2 Dockerfiles Dev
- [ ] Dockerfile.dev para cada API
- [ ] Support para air (go hot reload)
- [ ] Support para debugger (dlv)

### 3.3 Scripts Adicionales
- [ ] logs.sh - Ver logs de servicios
- [ ] shell.sh - Entrar a contenedor
- [ ] restart.sh - Reiniciar servicios
- [ ] rebuild.sh - Reconstruir imágenes

### 3.4 Documentación Completa
- [ ] TROUBLESHOOTING.md
  - [ ] Puertos en uso
  - [ ] Volúmenes dañados
  - [ ] Errores de conexión
  - [ ] Performance issues
  
- [ ] DEV_GUIDE.md
  - [ ] Desarrollo local
  - [ ] Hot reload setup
  - [ ] Debugging
  - [ ] Common tasks
  
- [ ] ARCHITECTURE.md
  - [ ] Diagrama de servicios
  - [ ] Flujos de datos
  - [ ] Puertos y endpoints

### 3.5 CI/CD Integration
- [ ] GitHub Actions para validar docker-compose
- [ ] Test que stack pueda levantarse
- [ ] Test que health checks pasen

**Entregables:**
- [ ] Perfil development completo
- [ ] Hot reload funcionando
- [ ] Documentación exhaustiva
- [ ] Troubleshooting guide

---

## Fase 4: Optimizaciones y Testing (Post-Sprints)

### 4.1 Performance
- [ ] Memory limits en compose
- [ ] CPU limits si es necesario
- [ ] Caché de layers Docker

### 4.2 Testing
- [ ] Test que docker-compose.yml es válido
- [ ] Test que todos servicios levanten
- [ ] Test que health checks pasen
- [ ] Test que APIs responden

### 4.3 Documentación de Operaciones
- [ ] Backup y restore de datos
- [ ] Upgrade de versiones
- [ ] Scaling (si se necesita)

---

## Estructura Final

```
dev-environment/
├── docker-compose.yml           # Stack producción
├── docker-compose.override.yml  # Overrides locales (dev)
├── docker/
│   ├── postgres/
│   │   ├── Dockerfile          # Base de datos
│   │   ├── init.sql            # Schema
│   │   └── seed.sql            # Datos
│   ├── mongo/
│   │   ├── Dockerfile
│   │   ├── init.js
│   │   └── seed.js
│   ├── rabbitmq/
│   │   └── rabbitmq.conf       # Config
│   └── redis/
│       └── Dockerfile
├── scripts/
│   ├── setup.sh                # Iniciar
│   ├── teardown.sh             # Detener
│   ├── clean.sh                # Limpiar
│   ├── seed-data.sh            # Cargar datos
│   ├── health-check.sh         # Verificar salud
│   ├── logs.sh                 # Ver logs
│   ├── shell.sh                # Entrar contenedor
│   └── rebuild.sh              # Reconstruir
├── docs/
│   ├── SETUP.md                # Instalación
│   ├── TROUBLESHOOTING.md      # Problemas
│   ├── DEV_GUIDE.md            # Desarrollo
│   ├── SERVICES.md             # Detalles
│   └── ARCHITECTURE.md         # Diagrama
└── .gitignore
```

---

## Comandos Principales de Usuario

```bash
# Iniciar ambiente completo
./scripts/setup.sh

# Ver logs
./scripts/logs.sh
./scripts/logs.sh postgres    # Log específico

# Recargar datos
./scripts/clean.sh
./scripts/seed-data.sh

# Verificar salud
./scripts/health-check.sh

# Detener
./scripts/teardown.sh
```

---

## Criterios de Aceptación

- [ ] docker-compose up -d levanta todo
- [ ] Todos servicios saludables (health check pasan)
- [ ] Datos de prueba cargados
- [ ] APIs responden en puertos correctos
- [ ] Logs visibles
- [ ] Docker Compose file válido
- [ ] Documentación 100%
- [ ] Scripts funcionan en Linux y macOS

---

## Requisitos Previos para Usuarios

```bash
# Verificar requisitos antes de usar
- Docker Desktop 20.10+ (o Docker + Docker Compose 2.0+)
- 8GB RAM mínimo
- 20GB almacenamiento libre
- Puertos 5432, 27017, 5672, 15672, 8080, 8081 libres
```

---

## Riesgos y Mitigaciones

| Riesgo | Mitigation |
|--------|-----------|
| Puertos en uso | Usar docker-compose up -p "prefix" |
| Out of memory | Limitar memoria en services |
| Volúmenes corruptos | Script de clean.sh |
| Diferencias Mac/Linux | Usar alpine images |
| Datos no sincronizados | Seed scripts idempotentes |

---

## Próximos Pasos Después de Sprint 3

1. Publicar README en GitHub
2. Agregar badges (build status, etc)
3. Crear issues para feedback
4. Documentar troubleshooting común
5. Crear video tutorial (opcional)

---

**Próxima revisión:** Después de Sprint 1  
**Última actualización:** 15 de Noviembre, 2025
