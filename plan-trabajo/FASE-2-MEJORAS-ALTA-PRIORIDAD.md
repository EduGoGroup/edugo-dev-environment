# FASE 2: Mejoras de Alta Prioridad

**Prioridad:** Alta  
**Estimación:** 4-6 horas  
**Riesgo:** Medio  

---

## Objetivo

Implementar mejoras críticas que impactan directamente la experiencia de desarrollo.

---

## 2.1 Validación de Health Checks en Setup

### Problema Actual

`setup.sh` espera 10 segundos fijos en lugar de verificar health checks reales.

```bash
# Código actual problemático:
sleep 10
```

### Solución Propuesta

Implementar espera activa con verificación de health checks.

### Pasos de Implementación

#### Paso 2.1.1: Crear función de espera de health

**Archivo:** `scripts/setup.sh`

Agregar función:

```bash
wait_for_healthy() {
  local service=$1
  local max_attempts=${2:-30}
  local attempt=1
  
  echo "Esperando que $service esté saludable..."
  
  while [ $attempt -le $max_attempts ]; do
    local health=$(docker inspect --format='{{.State.Health.Status}}' "edugo-$service" 2>/dev/null)
    
    if [ "$health" = "healthy" ]; then
      echo "$service está saludable"
      return 0
    fi
    
    echo "Intento $attempt/$max_attempts: $service estado=$health"
    sleep 2
    attempt=$((attempt + 1))
  done
  
  echo "ERROR: $service no alcanzó estado saludable"
  return 1
}
```

#### Paso 2.1.2: Reemplazar sleep por verificaciones

Cambiar:

```bash
sleep 10
```

Por:

```bash
wait_for_healthy "postgres" 30 || exit 1
wait_for_healthy "mongodb" 30 || exit 1
wait_for_healthy "rabbitmq" 30 || exit 1
```

#### Paso 2.1.3: Agregar timeout total

```bash
TOTAL_TIMEOUT=120
START_TIME=$(date +%s)

# En wait_for_healthy, verificar timeout global
CURRENT_TIME=$(date +%s)
if [ $((CURRENT_TIME - START_TIME)) -gt $TOTAL_TIMEOUT ]; then
  echo "ERROR: Timeout global alcanzado"
  exit 1
fi
```

### Validación

- [ ] `setup.sh` espera a que postgres esté healthy
- [ ] `setup.sh` espera a que mongodb esté healthy  
- [ ] `setup.sh` espera a que rabbitmq esté healthy
- [ ] Timeout global funciona correctamente
- [ ] Mensajes de progreso son claros

### Commit Sugerido

```
feat(scripts): implementar espera activa de health checks en setup

- Agregar función wait_for_healthy() con reintentos
- Reemplazar sleep 10 por verificaciones reales
- Agregar timeout global de 120 segundos
- Mejorar mensajes de progreso
```

---

## 2.2 Seed Data Automático

### Problema Actual

`seed-data.sh` existe pero no se ejecuta por defecto en setup.

### Solución Propuesta

Integrar seed como paso opcional automático en setup.

### Pasos de Implementación

#### Paso 2.2.1: Agregar flag --seed a setup.sh

**Archivo:** `scripts/setup.sh`

```bash
# En la sección de parsing de argumentos
LOAD_SEEDS=false

while [[ $# -gt 0 ]]; do
  case $1 in
    -s|--seed)
      LOAD_SEEDS=true
      shift
      ;;
    # ... otros casos
  esac
done
```

#### Paso 2.2.2: Ejecutar seeds después de migraciones

```bash
# Después de que migrator termine
if [ "$LOAD_SEEDS" = true ]; then
  echo "Cargando datos de prueba..."
  ./scripts/seed-data.sh
fi
```

#### Paso 2.2.3: Actualizar documentación de ayuda

```bash
show_help() {
  echo "Uso: setup.sh [opciones]"
  echo ""
  echo "Opciones:"
  echo "  -s, --seed       Cargar datos de prueba después del setup"
  echo "  --profile NAME   Usar perfil específico (full, db-only, api-only)"
  echo "  -h, --help       Mostrar esta ayuda"
}
```

#### Paso 2.2.4: Agregar al Makefile

**Archivo:** `Makefile`

```makefile
setup-with-seeds: ## Setup completo con datos de prueba
	./scripts/setup.sh --seed

seed: ## Cargar datos de prueba
	./scripts/seed-data.sh
```

### Validación

- [ ] `./scripts/setup.sh --seed` carga datos de prueba
- [ ] `./scripts/setup.sh -s` carga datos de prueba
- [ ] Sin flag, no carga seeds
- [ ] `make setup-with-seeds` funciona
- [ ] `make seed` funciona independientemente

### Commit Sugerido

```
feat(scripts): integrar seed data como opción en setup

- Agregar flags -s/--seed para cargar datos de prueba
- Ejecutar seeds después de migraciones exitosas
- Agregar comandos make setup-with-seeds y make seed
- Actualizar documentación de ayuda
```

---

## Resumen de Commits de Fase 2

1. `feat(scripts): implementar espera activa de health checks en setup`
2. `feat(scripts): integrar seed data como opción en setup`

---

## Dependencias

- Fase 1 debe estar completada (paquetes actualizados)
- Docker debe estar funcionando correctamente

---

## Flujo de Trabajo Git

### 1. Crear rama desde dev

```bash
git checkout dev
git pull origin dev
git checkout -b fase-2-mejoras-alta-prioridad
```

### 2. Realizar los cambios

Ejecutar los pasos de implementación descritos arriba, haciendo commits atómicos.

### 3. Crear PR hacia dev

```bash
git push origin fase-2-mejoras-alta-prioridad
# Crear PR en GitHub hacia dev
```

---

## Documentación a Actualizar

Al completar esta fase, actualizar los siguientes documentos:

| Documento | Cambio Requerido |
|-----------|------------------|
| `documentos/GUIA-RAPIDA.md` | Documentar nuevos flags de setup.sh (-s/--seed) |
| `documentos/FAQ.md` | Agregar preguntas sobre health checks y seeds |
| `documentos/DEPRECADO-MEJORAS.md` | Marcar mejoras 1 y 2 como completadas |

### Checklist de Cierre

- [ ] Health checks implementados y funcionando
- [ ] Seed data integrado con flag --seed
- [ ] Makefile actualizado con nuevos comandos
- [ ] `documentos/GUIA-RAPIDA.md` actualizado
- [ ] `documentos/FAQ.md` actualizado
- [ ] `documentos/DEPRECADO-MEJORAS.md` actualizado
- [ ] PR creado hacia `dev`
- [ ] PR revisado y aprobado
- [ ] PR mergeado a `dev`
