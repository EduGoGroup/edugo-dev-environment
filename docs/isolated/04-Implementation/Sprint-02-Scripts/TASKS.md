# Tareas Sprint 02
## TASK-02-001: Scripts Automatizados
```bash
# scripts/setup.sh
#!/bin/bash
PROFILE=${1:-full}
docker compose --profile $PROFILE up -d

# scripts/seed-data.sh
psql < seeds/postgresql/*.sql
mongosh < seeds/mongodb/*.js

# scripts/stop.sh
docker compose down
```
**Tiempo:** 3h
