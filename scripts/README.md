# Scripts de Utilidad - EduGo Dev Environment

Este directorio contiene scripts de utilidad para gestionar el ambiente de desarrollo local de EduGo.

---

## 📋 Scripts Disponibles

| Script | Propósito | Cuándo Usar |
|--------|-----------|-------------|
| **setup.sh** | Inicializar ambiente completo | Primera vez o reset completo |
| **validate.sh** | Validar docker-compose files | Antes de hacer cambios o commits |
| **seed-data.sh** | Cargar datos de prueba | Desarrollo y testing |
| **stop.sh** | Detener servicios | Finalizar trabajo |
| **cleanup.sh** | Limpiar ambiente | Liberar espacio en disco |
| **update-images.sh** | Actualizar imágenes Docker | Obtener últimas versiones |

---

## 🔍 validate.sh

Valida la sintaxis de todos los archivos `docker-compose*.yml` sin levantar contenedores.

### Uso

```bash
./scripts/validate.sh
```

### ¿Qué Valida?

- ✅ Sintaxis YAML correcta
- ✅ Servicios definidos
- ✅ Volúmenes definidos
- ✅ Puertos expuestos
- ✅ Variables de entorno requeridas
- ⚠️ Existencia de .env

### Archivos Validados

- `docker/docker-compose.yml` (principal)
- `docker/docker-compose.full.yml` (completo)
- `docker/docker-compose.local.yml` (local)

### Salida Esperada

```
🔍 Validando configuración de Docker Compose...

✅ docker-compose instalado
   Versión: Docker Compose version v2.x.x

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📄 Validando: docker-compose.yml
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Archivo encontrado
📝 Validando sintaxis YAML...
✅ Sintaxis YAML válida

🔍 Servicios encontrados:
  ✓ postgres
  ✓ mongodb
  ✓ rabbitmq
  ✓ api-mobile
  ✓ api-admin
  ✓ worker

💾 Volúmenes encontrados:
  ✓ postgres_data
  ✓ mongodb_data
  ✓ rabbitmq_data

🌐 Puertos expuestos:
  ✓ 5432:5432
  ✓ 15672:15672
  ✓ 27017:27017
  ✓ 5672:5672
  ✓ 8081:8081
  ✓ 8082:8082

✅ Validación completada exitosamente

Próximo paso:
  cd docker && docker-compose up -d
```

### Errores Comunes

**Error: docker-compose no instalado**
```bash
❌ ERROR: docker-compose no está instalado
```
**Solución:** Instalar Docker Compose

**Error: sintaxis YAML inválida**
```bash
❌ ERROR: Sintaxis YAML inválida
```
**Solución:** Revisar indentación y sintaxis en el archivo correspondiente

---

## 🚀 setup.sh

Inicializa el ambiente de desarrollo completo.

### Uso

```bash
# Setup completo (todos los servicios)
./scripts/setup.sh

# Solo bases de datos
./scripts/setup.sh --profile db-only

# APIs con datos de prueba
./scripts/setup.sh --profile api-only --seed
```

### Opciones

| Opción | Descripción | Default |
|--------|-------------|---------|
| `-p, --profile` | Perfil de Docker Compose | `full` |
| `-s, --seed` | Cargar datos de prueba | `false` |
| `-h, --help` | Mostrar ayuda | - |

### Perfiles Disponibles

- `full` - Todos los servicios (default)
- `db-only` - Solo bases de datos
- `api-only` - Bases de datos + APIs
- `mobile-only` - Bases de datos + API Mobile
- `admin-only` - Bases de datos + API Administración
- `worker-only` - Bases de datos + Worker

### Qué Hace

1. Verifica que Docker Desktop esté corriendo
2. Solicita autenticación en GitHub Container Registry
3. Descarga las imágenes Docker necesarias
4. Crea archivo `.env` desde `.env.example`
5. Levanta los servicios según el perfil seleccionado
6. Ejecuta migraciones automáticas
7. (Opcional) Carga datos de prueba

---

## 🌱 seed-data.sh

Carga datos de prueba en las bases de datos.

### Uso

```bash
./scripts/seed-data.sh
```

### Qué Hace

- Carga usuarios de prueba
- Carga instituciones de ejemplo
- Carga datos de configuración
- Verifica que los datos se cargaron correctamente

### Pre-requisitos

- Servicios deben estar corriendo
- Migraciones deben estar ejecutadas

---

## 🛑 stop.sh

Detiene los servicios de desarrollo.

### Uso

```bash
# Detener servicios (mantiene datos)
./scripts/stop.sh

# Detener perfil específico
./scripts/stop.sh --profile db-only

# Detener y eliminar volúmenes (⚠️ borra datos)
./scripts/stop.sh --volumes
```

### Opciones

| Opción | Descripción | Efecto |
|--------|-------------|--------|
| `--profile` | Perfil a detener | Solo detiene servicios del perfil |
| `--volumes` | Eliminar volúmenes | ⚠️ Borra todos los datos |

---

## 🧹 cleanup.sh

Limpia el ambiente de desarrollo y libera espacio en disco.

### Uso

```bash
./scripts/cleanup.sh
```

### Qué Hace (Interactivo)

El script pregunta antes de cada acción:

1. **Detener contenedores** - Para servicios corriendo
2. **Eliminar volúmenes** - ⚠️ Borra datos de BD
3. **Limpiar imágenes no usadas** - Libera espacio
4. **Eliminar imágenes de EduGo** - Fuerza re-descarga

### Casos de Uso

- **Espacio en disco lleno** - Ejecutar limpieza completa
- **Reset completo** - Eliminar todo y empezar de cero
- **Problemas con imágenes** - Eliminar y re-descargar

---

## 🔄 update-images.sh

Actualiza las imágenes Docker a sus últimas versiones.

### Uso

```bash
./scripts/update-images.sh
```

### Qué Hace

1. Descarga últimas versiones de:
   - `ghcr.io/edugogroup/edugo-api-mobile:latest`
   - `ghcr.io/edugogroup/edugo-api-administracion:latest`
   - `ghcr.io/edugogroup/edugo-worker:latest`

2. Reinicia servicios con nuevas imágenes

### Cuándo Usar

- Después de un release de APIs
- Para obtener últimas funcionalidades
- Para probar cambios recientes

---

## 🔗 Flujo de Trabajo Típico

### Primera Vez

```bash
# 1. Validar configuración
./scripts/validate.sh

# 2. Setup completo
./scripts/setup.sh

# 3. (Opcional) Cargar datos de prueba
./scripts/seed-data.sh
```

### Desarrollo Diario

```bash
# Iniciar servicios
cd docker && docker-compose up -d

# Trabajar...

# Detener al finalizar
./scripts/stop.sh
```

### Actualizar Versiones

```bash
# Actualizar imágenes
./scripts/update-images.sh

# Verificar nuevas versiones
docker-compose ps
```

### Limpiar Ambiente

```bash
# Limpieza completa
./scripts/cleanup.sh

# Re-inicializar
./scripts/setup.sh
```

---

## ⚠️ Notas Importantes

- **validate.sh** NO requiere servicios corriendo (solo valida sintaxis)
- **setup.sh** requiere autenticación en `ghcr.io`
- **cleanup.sh** es DESTRUCTIVO - elimina datos si se confirma
- **stop.sh** con `--volumes` BORRA datos permanentemente
- Todos los scripts deben ejecutarse desde la raíz del proyecto

---

## 🐛 Troubleshooting

### Script no ejecuta

**Error:**
```
Permission denied
```

**Solución:**
```bash
chmod +x scripts/*.sh
```

### Docker no encontrado

**Error:**
```
docker: command not found
```

**Solución:**
```bash
# macOS
open -a Docker

# Esperar a que inicie
```

### Error de autenticación en ghcr.io

**Error:**
```
Error response from daemon: unauthorized
```

**Solución:**
```bash
# Re-ejecutar setup que solicita credenciales
./scripts/setup.sh
```

---

## 📝 Licencia

Privado - EduGo © 2025
