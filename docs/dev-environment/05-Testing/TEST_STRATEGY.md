# Test Strategy - spec-05
## Tests de Scripts
```bash
# Test setup.sh
./scripts/setup.sh --profile db-only
docker ps | grep postgres

# Test seeds
./scripts/seed-data.sh
# Verificar datos
```
