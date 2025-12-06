# Reglas de Gestión de Sprints y Pull Requests - EduGo Dev Environment

Este documento establece las reglas obligatorias para la ejecución de sprints, gestión de ramas, creación de PRs y su ciclo de vida completo.

> **Nota:** Este proyecto (`edugo-dev-environment`) es especial porque puede modificar código propio Y de otros repositorios del ecosistema EduGo.

---

## 1. Pre-Sprint: Preparación del Entorno

### 1.1 Reglas de Ramas

| Regla | Descripción |
|-------|-------------|
| **Origen de ramas** | Toda rama feature debe partir de `dev` (o `main` si `dev` no existe) |
| **Commits en dev/main** | Prohibido hacer commits directamente, excepto documentación (nunca código) |
| **Estado de dev** | `dev` debe estar siempre funcional; no enviar PRs con fallos |
| **Rol de main** | `main` es producción; solo recibe merges desde `dev` cuando se solicite explícitamente |
| **Releases** | Solo crear releases e imágenes Docker desde `main` (excepto casos especiales documentados) |

### 1.2 Proyectos que Este Repositorio Puede Modificar

Este proyecto (`edugo-dev-environment`) tiene acceso para modificar:

| Proyecto | Rol | Flujo de PR |
|----------|-----|-------------|
| `edugo-dev-environment` | Propio | feature → dev → main |
| `edugo-api-mobile` | Dependencia | feature → dev → main (+ release si se solicita) |
| `edugo-api-administracion` | Dependencia | feature → dev → main (+ release si se solicita) |
| `edugo-worker` | Dependencia | feature → dev → main (+ release si se solicita) |
| `edugo-shared` | Dependencia | feature → dev → main + **release obligatorio** |
| `edugo-infrastructure` | Dependencia | feature → dev → main + **release obligatorio** |

### 1.3 Reglas para Proyectos de Dependencia (edugo-shared, edugo-infrastructure)

1. **Flujo obligatorio:** El código debe pasar por el flujo completo de PR (feature → dev → main)
2. **Releases por módulo:** Los releases se crean desde `main` por módulo, nunca como versión única del proyecto
3. **Versionamiento:** Usar formato `0.x.y` (estamos en desarrollo)
4. **Consumo de dependencias:** Siempre usar `go get` con versión específica; nunca referenciar código local
5. **GitHub Releases:** Solo mediante GitHub Releases (no tags simples) se puede acceder via `go get`

### 1.4 Proyectos API (api-mobile, api-administracion, worker)

- Objetivo: Llegar hasta `dev` con PRs exitosos
- PR de `dev` a `main`: Solo cuando se solicite explícitamente
- Releases manuales: Solo cuando se indique explícitamente

---

## 2. Pre-Sprint: Evaluación del Estado Actual

### 2.1 Si estás en rama `dev` o `main`

```
┌─────────────────────────────────────────────────────────────┐
│ VERIFICAR: ¿Hay código modificado sin commitear?            │
├─────────────────────────────────────────────────────────────┤
│ SI → Descartar cambios (ejecución incompleta de sprint     │
│      anterior; el código no debe estar en dev sin PR)       │
│ NO → Continuar con sincronización                           │
└─────────────────────────────────────────────────────────────┘
```

**Excepción:** Si solo hay documentación creada, incluirla en el sprint actual para su guardado.

### 2.2 Si estás en una rama feature

Ejecutar la siguiente secuencia de evaluación:

```
┌──────────────────────────────────────────────────────────────┐
│ PASO 1: Identificar el sprint asociado a esta rama           │
├──────────────────────────────────────────────────────────────┤
│ PASO 2: Verificar estado del PR                              │
│                                                              │
│   ¿El PR ya fue mergeado a dev (o main)?                     │
│   ├─ SI → Eliminar rama local y remota                       │
│   │       Cambiar a dev y sincronizar con remoto             │
│   │       FIN                                                │
│   └─ NO → Continuar a Paso 3                                 │
│                                                              │
│ PASO 3: Verificar si existe PR abierto                       │
│   ├─ SI → Ir a sección "3.3 Gestión de PR Abierto"           │
│   └─ NO → Continuar a Paso 4                                 │
│                                                              │
│ PASO 4: Evaluar tracking de tareas                           │
│   • Revisar tareas ejecutadas vs planificadas                │
│   • Si no existe tracking → Crearlo comparando plan vs código│
│   • Si existe tracking → Validar coherencia con código       │
│   • Actualizar tracking según hallazgos                      │
│   • Continuar ejecución del sprint                           │
│   • Leer documentación del sprint para entrar en contexto    │
└──────────────────────────────────────────────────────────────┘
```

---

## 3. Crear Pull Request

### 3.1 Pre-requisitos antes de crear el PR

Ejecutar **en este orden** (algunos pasos pueden paralelizarse):

| # | Paso | Paralelizable |
|---|------|---------------|
| 1 | Revisar: planificado vs tracking vs código generado | No |
| 2 | Actualizar tracking de la tarea | No |
| 3 | Hacer commit de los cambios | No |
| 4 | Compilar la aplicación (todas las salidas) | Sí |
| 5 | Ejecutar linters según tecnología | Sí |
| 6 | Ejecutar pruebas unitarias | Sí |
| 7 | Ejecutar pruebas de integración | Sí (después de unitarias) |

**Nota:** Los pasos 4, 5 y 6 pueden ejecutarse en paralelo para optimizar tiempo.

### 3.2 Destino del PR

| Rama actual | Destino del PR |
|-------------|----------------|
| `feature/*` | `dev` |
| `dev` | `main` (solo para shared/infrastructure o cuando se solicite) |

### 3.3 Gestión de PR Abierto

#### 3.3.1 Monitoreo de Pipelines y Code Review

- **Intervalo de revisión:** Cada 1 minuto
- **Tiempo máximo de espera:** 10 minutos
- **Si excede 10 minutos:** Detener monitoreo, notificar al usuario con:
  - Estado actual de los pipelines
  - Instrucciones de cómo proceder manualmente
  - Razón probable del retraso

#### 3.3.2 Resolución de Fallos en Pipelines

```
┌─────────────────────────────────────────────────────────────┐
│ REGLA DE 3 INTENTOS                                         │
├─────────────────────────────────────────────────────────────┤
│ • Máximo 3 intentos para resolver cada error                │
│ • Documentar cada intento en el tracking                    │
│ • Cada push reinicia el conteo de tiempo (10 min)           │
│ • Si falla tras 3 intentos:                                 │
│   1. Detener el proceso                                     │
│   2. Notificar al usuario con informe detallado             │
│   3. Incluir: error, intentos realizados, posibles causas   │
└─────────────────────────────────────────────────────────────┘
```

#### 3.3.3 Gestión de Comentarios de Copilot/Code Review

Evaluar cada comentario según la siguiente matriz de decisión:

| Tipo de Comentario | Puntos Fibonacci | Acción |
|--------------------|------------------|--------|
| Traducción/redacción | - | Ignorar |
| Valioso y aplicable | ≤ 3 | Aplicar de inmediato |
| Mejora moderada | 4-8 | Crear tarea en siguiente sprint |
| Mejora significativa | > 8 | Crear documento de deuda técnica y DETENER si es bloqueante |
| De otro sprint (menor) | ≤ 3 | Corregir ahora |
| De otro sprint (medio) | 4-8 | Buscar sprint planificado, responder en comentario o agregar al siguiente |
| De otro sprint (mayor) | > 8 | Crear deuda técnica; detener si es bloqueante |

**Documento de deuda técnica debe incluir:**
- Descripción detallada del problema
- Razón de la deuda técnica
- Por qué no se detectó en el análisis
- Pasos para resolverlo

**Documentación obligatoria:**
- Documentar todos los comentarios no aplicados con justificación
- Cada cambio de documentación o tracking debe commitearse y pushearse (para mantener historial)

#### 3.3.4 Criterios para Hacer Merge

| Criterio | Obligatorio |
|----------|-------------|
| Todos los pipelines exitosos | ✅ Sí |
| Comentarios de Copilot resueltos o documentados | ✅ Sí |
| Todos los tests pasando | ✅ Sí |

**Excepciones para merge con pipeline fallido:**
- Error administrativo en los pipelines
- Límites de servicios de GitHub (suscripción agotada)

**Sobre tests:**
- Ningún test que falle puede ignorarse
- Si un test falla constantemente → Analizar si:
  - El test está mal diseñado
  - El test detectó fallas reales en el código

---

## 4. Post-Merge

### 4.1 Si el PR era hacia `dev`

```
┌─────────────────────────────────────────────────────────────┐
│ CHECKLIST POST-MERGE A DEV                                  │
├─────────────────────────────────────────────────────────────┤
│ □ Eliminar rama feature local: git branch -d feature/xxx    │
│ □ Eliminar rama feature remota: git push origin --delete    │
│ □ Cambiar a dev: git checkout dev                           │
│ □ Sincronizar con remoto: git pull origin dev               │
│ □ Validar que todas las tareas del sprint están en dev      │
│   • Si falta alguna → Investigar razón                      │
│   • Si es necesario → Crear tarea para siguiente sprint     │
│ □ Comenzar análisis del siguiente sprint                    │
│ □ Repetir ciclo desde sección 2                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Si el PR era hacia `main`

Aplica para: `edugo-shared`, `edugo-infrastructure`, o cuando se solicite explícitamente.

```
┌─────────────────────────────────────────────────────────────┐
│ CHECKLIST POST-MERGE A MAIN                                 │
├─────────────────────────────────────────────────────────────┤
│ □ Aplicar mismas reglas de 4.1 (cambiar a main)             │
│ □ Monitorear pipelines post-merge                           │
│   • Si alguno falla:                                        │
│     1. Investigar la causa                                  │
│     2. Crear hotfix en dev                                  │
│     3. Hacer PR de dev a main                               │
│     4. Seguir proceso completo                              │
│   • Nota: Pipelines post-merge no afectan código pero       │
│     pueden enmascarar otros errores                         │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 Creación de Release

| Proyecto | Cuándo crear release |
|----------|---------------------|
| `edugo-shared` | **Obligatorio** después de merge a main (necesario para `go get`) |
| `edugo-infrastructure` | **Obligatorio** después de merge a main (necesario para `go get`) |
| `edugo-api-mobile` | Solo cuando se solicite explícitamente |
| `edugo-api-administracion` | Solo cuando se solicite explícitamente |
| `edugo-worker` | Solo cuando se solicite explícitamente |
| `edugo-dev-environment` | Solo cuando se solicite explícitamente |

**Proceso de release:**
1. Ejecutar pipeline "manual release" desde `main`
2. Monitorear solo el inicio (que arranque correctamente)
3. No es necesario esperar la finalización (puede tardar por creación de imagen Docker)

---

## 5. Flujo para Modificar Proyectos Externos

Cuando `edugo-dev-environment` necesita modificar otro proyecto:

```
┌─────────────────────────────────────────────────────────────┐
│ FLUJO PARA MODIFICAR PROYECTO EXTERNO                       │
├─────────────────────────────────────────────────────────────┤
│ 1. Clonar/actualizar el proyecto destino                    │
│ 2. Crear rama feature desde dev del proyecto                │
│ 3. Realizar cambios necesarios                              │
│ 4. Ejecutar tests del proyecto                              │
│ 5. Commit y push                                            │
│ 6. Crear PR a dev                                           │
│ 7. Monitorear pipelines (regla de 10 min / 3 intentos)      │
│ 8. Merge a dev                                              │
│ 9. Si se solicita → PR de dev a main                        │
│ 10. Si se solicita → Crear release                          │
│ 11. Actualizar edugo-dev-environment si es necesario        │
│     (ej: nueva versión de imagen Docker)                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Resumen Visual del Flujo

```
                    ┌─────────────┐
                    │   INICIO    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ ¿En qué     │
                    │ rama estás? │
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────▼─────┐   ┌──────▼──────┐   ┌─────▼─────┐
    │    dev    │   │  feature/*  │   │   main    │
    └─────┬─────┘   └──────┬──────┘   └─────┬─────┘
          │                │                │
    ┌─────▼─────┐   ┌──────▼──────┐   ┌─────▼─────┐
    │ Descartar │   │  Evaluar    │   │  (Raro)   │
    │ cambios   │   │  estado PR  │   │  Evaluar  │
    │ si hay    │   │  y tracking │   │           │
    └─────┬─────┘   └──────┬──────┘   └───────────┘
          │                │
          └───────┬────────┘
                  │
           ┌──────▼──────┐
           │  Ejecutar   │
           │   Sprint    │
           └──────┬──────┘
                  │
           ┌──────▼──────┐
           │  Validar    │
           │  (compile,  │
           │  lint,test) │
           └──────┬──────┘
                  │
           ┌──────▼──────┐
           │  Crear PR   │
           └──────┬──────┘
                  │
           ┌──────▼──────┐
           │  Monitorear │
           │  (10 min    │
           │  máximo)    │
           └──────┬──────┘
                  │
           ┌──────▼──────┐
           │ ¿Éxito?     │
           └──────┬──────┘
                  │
        ┌─────────┼─────────┐
        │                   │
  ┌─────▼─────┐       ┌─────▼─────┐
  │    SÍ     │       │    NO     │
  │  Merge    │       │  Resolver │
  └─────┬─────┘       │  (máx 3   │
        │             │  intentos)│
        │             └─────┬─────┘
        │                   │
  ┌─────▼─────┐             │
  │  Post-    │◄────────────┘
  │  Merge    │  (si resuelto)
  └─────┬─────┘
        │
  ┌─────▼─────┐
  │ Siguiente │
  │  Sprint   │
  └───────────┘
```

---

## 7. Constantes y Límites

| Parámetro | Valor |
|-----------|-------|
| Tiempo máximo monitoreo pipeline | 10 minutos |
| Intervalo de verificación | 1 minuto |
| Máximo intentos por error | 3 |
| Umbral Fibonacci para aplicar ahora | ≤ 3 puntos |
| Umbral Fibonacci para siguiente sprint | 4-8 puntos |
| Umbral Fibonacci para deuda técnica | > 8 puntos |

---

**Última actualización:** 1 de Diciembre, 2025
**Versión:** 1.0 (adaptado para edugo-dev-environment)
