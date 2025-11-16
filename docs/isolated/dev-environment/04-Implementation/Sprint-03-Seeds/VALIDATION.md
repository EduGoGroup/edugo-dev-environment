# Val Sprint 03
```bash
./scripts/seed-data.sh
psql -c "SELECT COUNT(*) FROM schools;"  # Esperado: 3
mongosh --eval "db.material_summary.countDocuments()"  # Esperado: 5
```
