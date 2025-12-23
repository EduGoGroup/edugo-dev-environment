# FASE 3: Mejoras de Media Prioridad

**Prioridad:** Media  
**Estimación:** 1-2 días  
**Riesgo:** Bajo-Medio  

---

## Objetivo

Implementar mejoras que facilitan el uso y la documentación del proyecto.

---

## 3.1 Documentar Variables de Entorno

### Problema Actual

Variables repetidas entre servicios con nombres diferentes:
- API Mobile usa `DATABASE_POSTGRES_*`
- API Admin usa `EDUGO_ADMIN_DATABASE_POSTGRES_*`
- Worker usa `EDUGO_WORKER_DATABASE_POSTGRES_*`

### Solución Propuesta

Crear documentación clara de todas las variables y su mapeo.

### Pasos de Implementación

#### Paso 3.2.1: Crear documento de variables

**Archivo:** `documentos/VARIABLES-ENTORNO.md`

```markdown
# Variables de Entorno

## Resumen de Prefijos por Servicio

| Servicio | Prefijo | Ejemplo |
|----------|---------|---------|
| API Mobile | `DATABASE_*` | DATABASE_POSTGRES_HOST |
| API Admin | `EDUGO_ADMIN_*` | EDUGO_ADMIN_DATABASE_POSTGRES_HOST |
| Worker | `EDUGO_WORKER_*` | EDUGO_WORKER_DATABASE_POSTGRES_HOST |

## PostgreSQL

| Variable Base | API Mobile | API Admin | Worker |
|---------------|------------|-----------|--------|
| Host | DATABASE_POSTGRES_HOST | EDUGO_ADMIN_DATABASE_POSTGRES_HOST | EDUGO_WORKER_DATABASE_POSTGRES_HOST |
| Port | DATABASE_POSTGRES_PORT | EDUGO_ADMIN_DATABASE_POSTGRES_PORT | EDUGO_WORKER_DATABASE_POSTGRES_PORT |
| User | DATABASE_POSTGRES_USER | EDUGO_ADMIN_DATABASE_POSTGRES_USER | EDUGO_WORKER_DATABASE_POSTGRES_USER |
| Password | DATABASE_POSTGRES_PASSWORD | EDUGO_ADMIN_DATABASE_POSTGRES_PASSWORD | EDUGO_WORKER_DATABASE_POSTGRES_PASSWORD |
| Database | DATABASE_POSTGRES_DB | EDUGO_ADMIN_DATABASE_POSTGRES_DB | EDUGO_WORKER_DATABASE_POSTGRES_DB |

## MongoDB

[Similar estructura...]

## RabbitMQ

[Similar estructura...]

## Variables Globales

| Variable | Descripción | Default |
|----------|-------------|---------|
| OPENAI_API_KEY | API Key de OpenAI para Worker | (requerido) |
| AUTH_JWT_SECRET | Secreto para tokens JWT | edugo-secret-2024 |
```

#### Paso 3.2.2: Actualizar .env.example

**Archivo:** `docker/.env.example`

Agregar comentarios explicativos:

```bash
# ===================================
# BASES DE DATOS - Configuración base
# ===================================
# Estas variables son usadas por la infraestructura Docker
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=edugo
POSTGRES_PASSWORD=edugo123
POSTGRES_DB=edugo

# ===================================
# API MOBILE
# ===================================
# Prefijo: DATABASE_*
DATABASE_POSTGRES_HOST=${POSTGRES_HOST}
# ... etc
```

### Validación

- [ ] Documento VARIABLES-ENTORNO.md creado
- [ ] .env.example tiene comentarios claros
- [ ] README principal referencia el documento

### Commit Sugerido

```
docs: documentar variables de entorno y prefijos por servicio

- Crear VARIABLES-ENTORNO.md con mapeo completo
- Mejorar comentarios en .env.example
- Referenciar desde README
```

---

## 3.2 Seeds Más Completos

### Problema Actual

Los seeds actuales son mínimos (solo 2 archivos SQL y 1 JS).

### Solución Propuesta

Expandir datos de prueba para cubrir más casos de uso.

### Pasos de Implementación

#### Paso 3.3.1: Expandir seeds de PostgreSQL

**Archivo:** `seeds/postgresql/02_more_users.sql`

```sql
-- Más usuarios de prueba con diferentes roles
INSERT INTO users (id, email, password_hash, role, ...) VALUES
  -- Profesores adicionales
  ('uuid-teacher-3', 'teacher.history@edugo.test', '...', 'teacher', ...),
  ('uuid-teacher-4', 'teacher.english@edugo.test', '...', 'teacher', ...),
  
  -- Estudiantes adicionales
  ('uuid-student-4', 'student4@edugo.test', '...', 'student', ...),
  ('uuid-student-5', 'student5@edugo.test', '...', 'student', ...),
  -- ... hasta 10 estudiantes
;
```

**Archivo:** `seeds/postgresql/03_courses_with_content.sql`

```sql
-- Cursos con contenido completo
INSERT INTO courses (...) VALUES (...);
INSERT INTO units (...) VALUES (...);
INSERT INTO materials (...) VALUES (...);
```

#### Paso 3.3.2: Expandir seeds de MongoDB

**Archivo:** `seeds/mongodb/02_documents.js`

```javascript
// Documentos de ejemplo
db.documents.insertMany([
  {
    _id: ObjectId("..."),
    title: "Introducción a Matemáticas",
    content: "...",
    processed: true,
    // ...
  }
]);
```

#### Paso 3.3.3: Actualizar seed-data.sh

**Archivo:** `scripts/seed-data.sh`

```bash
# Ejecutar todos los seeds en orden
for sql_file in seeds/postgresql/*.sql; do
  echo "Ejecutando: $sql_file"
  docker exec -i edugo-postgres psql -U edugo -d edugo < "$sql_file"
done

for js_file in seeds/mongodb/*.js; do
  echo "Ejecutando: $js_file"
  docker exec -i edugo-mongodb mongosh -u edugo -p edugo123 edugo --authSource admin < "$js_file"
done
```

### Validación

- [ ] Nuevos archivos de seed creados
- [ ] seed-data.sh ejecuta todos los archivos
- [ ] Datos verificables en bases de datos

### Commit Sugerido

```
feat(seeds): expandir datos de prueba

- Agregar más usuarios con diferentes roles
- Agregar cursos con contenido completo
- Agregar documentos de ejemplo en MongoDB
- Actualizar seed-data.sh para ejecutar todos los archivos
```

---

## 3.3 Corregir Profiles en docker-compose-apps.yml

### Problema Actual

Usa `profiles: ["with-admin"]` pero no hay documentación ni funciona bien.

### Solución Propuesta

Corregir o eliminar profiles no funcionales.

### Pasos de Implementación

#### Paso 3.4.1: Evaluar uso actual

Revisar si alguien usa los profiles:
- Si se usan: documentar y corregir
- Si no se usan: eliminar

#### Paso 3.4.2: Opción A - Corregir profiles

**Archivo:** `docker/docker-compose-apps.yml`

```yaml
services:
  api-mobile:
    # Sin profile - siempre disponible
    ...
    
  api-administracion:
    profiles: ["with-admin", "full"]
    ...
    
  worker:
    profiles: ["with-worker", "full"]
    ...
```

#### Paso 3.4.3: Opción B - Eliminar profiles

Simplificar el archivo eliminando profiles si no aportan valor.

### Validación

- [ ] docker-compose config valida sin errores
- [ ] Profiles funcionan según documentación
- [ ] O profiles eliminados si no se usan

### Commit Sugerido

```
fix(docker): corregir/eliminar profiles no funcionales en docker-compose-apps.yml
```

---

## Resumen de Commits de Fase 3

1. `docs: documentar variables de entorno y prefijos por servicio`
2. `feat(seeds): expandir datos de prueba`
3. `fix(docker): corregir profiles en docker-compose-apps.yml`

---

## Dependencias

- Fases 1 y 2 completadas (opcional pero recomendado)

---

## Flujo de Trabajo Git

### 1. Crear rama desde dev

```bash
git checkout dev
git pull origin dev
git checkout -b fase-3-mejoras-media-prioridad
```

### 2. Realizar los cambios

Ejecutar los pasos de implementación descritos arriba, haciendo commits atómicos por cada mejora (3.1, 3.2, 3.3).

### 3. Crear PR hacia dev

```bash
git push origin fase-3-mejoras-media-prioridad
# Crear PR en GitHub hacia dev
```

---

## Documentación a Actualizar

Al completar esta fase, actualizar los siguientes documentos:

| Documento | Cambio Requerido |
|-----------|------------------|
| `documentos/GUIA-RAPIDA.md` | Referenciar nuevo doc de variables |
| `documentos/FAQ.md` | Agregar sección sobre profiles |
| `documentos/DEPRECADO-MEJORAS.md` | Marcar mejoras como completadas |

### Checklist de Cierre

- [x] VARIABLES-ENTORNO.md creado
- [x] Seeds expandidos y funcionando
- [x] Profiles de docker-compose corregidos
- [x] `documentos/GUIA-RAPIDA.md` actualizado
- [x] `documentos/FAQ.md` actualizado
- [x] `documentos/DEPRECADO-MEJORAS.md` actualizado
- [ ] PR creado hacia `dev`
- [ ] PR revisado y aprobado
- [ ] PR mergeado a `dev`
