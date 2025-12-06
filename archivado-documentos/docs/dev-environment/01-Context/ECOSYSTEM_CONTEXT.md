# ECOSYSTEM CONTEXT - Dev Environment

## Posición en EduGo

**Rol:** Orquestador de Infraestructura - Hace funcionar toda la pila localmente  
**Interacción:** Contiene todos los otros servicios en Docker

---

## Mapa de Ecosistema

```
┌───────────────────────────────────────────────────┐
│    Dev Environment (Docker Compose)               │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ PostgreSQL Container                         │ │
│ │ - Port 5432                                  │ │
│ │ - Volumen: postgres_data                     │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ MongoDB Container                            │ │
│ │ - Port 27017                                 │ │
│ │ - Volumen: mongo_data                        │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ RabbitMQ Container                           │ │
│ │ - Port 5672, 15672 (Management)              │ │
│ │ - Volumen: rabbitmq_data                     │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ API Mobile Container (8080)                  │ │
│ ├─ Depende: PostgreSQL, MongoDB, RabbitMQ    │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ API Admin Container (8081)                   │ │
│ ├─ Depende: PostgreSQL                        │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ Worker Container (background)                │ │
│ ├─ Depende: RabbitMQ, MongoDB, PostgreSQL    │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ ┌──────────────────────────────────────────────┐ │
│ │ Redis Container (6379) - Optional            │ │
│ └──────────────────────────────────────────────┘ │
│                                                   │
│ Network: edugo-network (bridge)                  │
└───────────────────────────────────────────────────┘
```

---

## Flujos de Comunicación Internos

```
┌─────────────┐
│ API Mobile  │
│ :8080       │
└──────┬──────┘
       │
       ├─ TCP ──────────────► PostgreSQL:5432
       ├─ TCP ──────────────► MongoDB:27017
       └─ TCP ──────────────► RabbitMQ:5672
                                    │
                    ┌───────────────┴────────────────┐
                    │                                │
             ┌──────▼──────┐                  ┌──────▼──────┐
             │   Worker    │                  │ API Mobile  │
             │ (consumes)  │                  │ (consumes)  │
             └─────────────┘                  └─────────────┘


┌─────────────┐
│ API Admin   │
│ :8081       │
└──────┬──────┘
       │
       └─ TCP ──────────────► PostgreSQL:5432
```

---

## Datos Compartidos Entre Servicios

```
PostgreSQL (educgo_mobile, edugo_admin):

Tablas compartidas:
├─ users
├─ schools
├─ academic_units
├─ teachers
├─ students
├─ memberships
└─ enrollments

API Mobile escribe:        API Admin escribe:
├─ evaluations            ├─ schools
├─ questions              ├─ academic_units
├─ question_options       ├─ teachers
├─ answer_drafts          ├─ students
└─ evaluation_assignments └─ memberships


MongoDB (edugo_assessments):

Collections:
├─ evaluation_results (escrito por: API Mobile)
├─ evaluation_audit (escrito por: API Mobile)
├─ material_assessment (escrito por: Worker)
├─ material_summary (escrito por: Worker)
└─ material_event (escrito por: Worker)


RabbitMQ:

Messages:
├─ assessment.requests (API Mobile → Worker)
├─ assessment.responses (Worker → API Mobile)
├─ evaluation.events (API Mobile → audit)
└─ admin.events (API Admin → audit)
```

---

## Orden de Startup

```
1. PostgreSQL levanta (health check: SELECT 1)
   ↓
2. MongoDB levanta (health check: db.ping())
   ↓
3. RabbitMQ levanta (health check: rabbitmq-diagnostics)
   ↓
4. APIs esperan que dependencias estén ready (depends_on condition)
   ↓
5. API Mobile inicia cuando PostgreSQL + MongoDB + RabbitMQ = healthy
   ↓
6. API Admin inicia cuando PostgreSQL = healthy
   ↓
7. Worker inicia cuando RabbitMQ + MongoDB + PostgreSQL = healthy
   ↓
8. Seed data se carga (si dev environment)
   ↓
9. Sistema listo para testing
```

---

## Modificaciones en Dev Environment para Desarrollo

**docker-compose.override.yml** agrega:

```yaml
services:
  api-mobile:
    # Volumen: código local
    volumes:
      - ../repos-separados/edugo-api-mobile:/app
    # Hot reload
    command: air -c .air.toml
    # Más logs
    environment:
      - LOG_LEVEL=debug

  api-admin:
    volumes:
      - ../repos-separados/edugo-api-administracion:/app
    command: air -c .air.toml
    environment:
      - LOG_LEVEL=debug

  worker:
    volumes:
      - ../repos-separados/edugo-worker:/app
    command: air -c .air.toml
    environment:
      - LOG_LEVEL=debug
```

**Beneficio:**
- Cambios en código = recarga automática
- No requiere rebuild
- Desarrollo rápido

---

## Integración de Datos de Prueba

```
Setup completo:
1. docker-compose up -d
2. ./scripts/seed-data.sh
   ├─ Carga schema.sql en PostgreSQL
   ├─ Inserta escuelas en PostgreSQL
   ├─ Inserta académicas en PostgreSQL
   ├─ Inserta usuarios en PostgreSQL
   ├─ Declara exchanges en RabbitMQ
   ├─ Declara queues en RabbitMQ
   └─ Inserta documentos en MongoDB

3. ./scripts/health-check.sh
   ├─ Verifica PostgreSQL accesible
   ├─ Verifica MongoDB accesible
   ├─ Verifica RabbitMQ accesible
   ├─ Verifica API Mobile :8080
   ├─ Verifica API Admin :8081
   └─ Reporte final

4. Sistema listo
```

---

## Limpieza y Reset

```
Escenario: Datos se corrompieron, necesito reset

1. ./scripts/clean.sh
   ├─ docker-compose down -v
   │  └─ Elimina contenedores y volúmenes
   
2. ./scripts/setup.sh
   ├─ docker-compose build (si necesario)
   ├─ docker-compose up -d
   └─ Espera health checks

3. ./scripts/seed-data.sh
   └─ Carga datos frescos

4. Sistema listo
```

---

## Monitoreo y Troubleshooting

```
Ver estado de servicios:
docker-compose ps

Ver logs de servicio:
./scripts/logs.sh postgres
./scripts/logs.sh api-mobile

Entrar a contenedor:
./scripts/shell.sh api-mobile

Verificar conectividad:
./scripts/health-check.sh

Debug de puerto:
lsof -i :8080
```

---

## CI/CD Integration

Dev Environment puede ser usado en CI:

```yaml
# GitHub Actions
- name: Start dev environment
  run: docker-compose up -d

- name: Wait for services
  run: ./scripts/health-check.sh

- name: Run tests
  run: go test ./...

- name: Cleanup
  run: docker-compose down -v
```

---

## Diferencias: Local vs Staging vs Prod

```
LOCAL (Dev Environment):
├─ All services in Docker Compose
├─ Seed data included
├─ Logs to stdout
├─ Hot reload enabled
└─ Easier debugging

STAGING (Kubernetes):
├─ Separate PostgreSQL (RDS)
├─ Separate MongoDB (Atlas)
├─ Separate RabbitMQ (managed)
├─ APIs scaled horizontally
└─ Closer to production

PRODUCTION (Kubernetes):
├─ All managed services
├─ HA setup (replicas)
├─ Monitoring & alerting
├─ Backup & recovery
└─ Zero downtime deployments
```

---

## Checklist de Integración

- [ ] docker-compose.yml válido
- [ ] Todos servicios levantan correctamente
- [ ] Health checks implementados
- [ ] Seed data cargable
- [ ] APIs responden en puertos correctos
- [ ] Logs visibles
- [ ] Scripts funcionan en Linux y macOS
- [ ] Documentación completa
- [ ] README tiene instrucciones claras
