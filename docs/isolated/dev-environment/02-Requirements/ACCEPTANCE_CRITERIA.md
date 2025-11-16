# Acceptance Criteria - spec-05
## AC-001: Profile db-only Funciona
```bash
./scripts/setup.sh --profile db-only
docker ps | grep postgres
```
## AC-002: Seeds Insertan Datos
```bash
./scripts/seed-data.sh
psql -c "SELECT COUNT(*) FROM schools;"  # >3
```
