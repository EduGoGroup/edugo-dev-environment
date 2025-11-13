# Seeds de PostgreSQL

Datos de prueba para desarrollo y testing.

## Estructura

Los archivos se ejecutan en orden alfabético:

1. `01_schools.sql` - Escuelas
2. `02_users.sql` - Usuarios (admins, teachers, students)
3. `03_subjects.sql` - Materias/asignaturas
4. `04_materials.sql` - Materiales educativos

## Carga

```bash
# Automática con setup
./scripts/setup.sh --seed

# Manual
./scripts/seed-data.sh
```

## Nota

Estos son datos de ejemplo. En producción usar datos reales.
