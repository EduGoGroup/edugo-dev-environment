# Seeds de MongoDB

Datos de prueba para desarrollo local.

## Archivos

| Archivo | Descripción | Colección |
|---------|-------------|-----------|
| `01_material_summaries.js` | Resúmenes de materiales generados por IA | `material_summaries` |
| `02_assessments.js` | Evaluaciones con preguntas de ejemplo | `assessments` |

## Orden de Ejecución

Los archivos se ejecutan en orden alfabético. El prefijo numérico asegura el orden correcto.

## Contenido

### material_summaries

Cada documento incluye:
- `material_id`: Referencia al material en PostgreSQL
- `summary`: Resumen generado por IA
- `key_points`: Lista de puntos clave
- `topics`: Temas cubiertos
- `difficulty_level`: basic, intermediate, advanced
- `estimated_reading_time_minutes`: Tiempo estimado de lectura

### assessments

Cada documento incluye:
- `material_id`: Referencia al material en PostgreSQL
- `title`: Título de la evaluación
- `questions`: Array de preguntas (multiple_choice, true_false, short_answer)
- `total_points`: Puntos totales
- `passing_score`: Puntaje mínimo para aprobar
- `time_limit_minutes`: Tiempo límite

## Uso

```bash
# Desde el directorio raíz del proyecto
make seed-data

# O directamente
./scripts/seed-data.sh
```

## Notas

- Los seeds se relacionan con materiales de PostgreSQL por `material_id`
- Use los mismos UUIDs que en `seeds/postgresql/03_materials.sql`
