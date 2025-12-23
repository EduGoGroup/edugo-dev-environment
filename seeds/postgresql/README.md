# Seeds de PostgreSQL

Datos de prueba para desarrollo local.

## Archivos

| Archivo | Descripción | Dependencias |
|---------|-------------|--------------|
| `01_schools.sql` | 5 escuelas de ejemplo | Ninguna |
| `02_users.sql` | 2 admins, 5 docentes, 10 estudiantes | Ninguna |
| `03_materials.sql` | 12 materiales educativos | `01_schools.sql`, `02_users.sql` |
| `04_memberships.sql` | Asignación de usuarios a escuelas | `01_schools.sql`, `02_users.sql` |

## Orden de Ejecución

Los archivos se ejecutan en orden alfabético. El prefijo numérico asegura el orden correcto de dependencias.

## Credenciales de Prueba

Todos los usuarios tienen la misma contraseña: `password123`

| Email | Rol | Escuela |
|-------|-----|---------|
| admin@edugo.com | Admin | Liceo Técnico Santiago |
| teacher.fisica@edugo.com | Docente | Liceo Técnico Santiago |
| teacher.matematicas@edugo.com | Docente | Liceo Técnico Santiago |
| teacher.historia@edugo.com | Docente | Colegio Valparaíso |
| student1@edugo.com | Estudiante | Liceo Técnico Santiago |

## Uso

```bash
# Desde el directorio raíz del proyecto
make seed-data

# O directamente
./scripts/seed-data.sh
```

## Notas

- Los seeds usan `ON CONFLICT DO NOTHING` para ser idempotentes
- Los UUIDs son fijos para facilitar referencias cruzadas
- El hash de contraseña es bcrypt con cost 10
