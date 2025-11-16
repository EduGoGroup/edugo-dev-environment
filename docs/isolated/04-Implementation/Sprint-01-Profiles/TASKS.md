# Tareas Sprint 01
## TASK-01-001: Docker Compose Profiles
```yaml
# docker-compose.yml
services:
  postgres:
    profiles: [full, db-only]
  mongodb:
    profiles: [full, db-only]
  rabbitmq:
    profiles: [full, db-only]
  api-mobile:
    profiles: [full, api-only, mobile-only]
  api-admin:
    profiles: [full, api-only, admin-only]
  worker:
    profiles: [full, worker-only]
```
**Tiempo:** 3h
