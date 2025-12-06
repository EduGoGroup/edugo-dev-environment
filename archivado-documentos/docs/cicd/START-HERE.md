# 🎯 COMIENZA AQUÍ - edugo-dev-environment

⚠️ **UBICACIÓN Y CONTEXTO DE TRABAJO:**

```
┌─────────────────────────────────────────────────────────────────────┐
│ 📍 Estás en: docs/cicd/ (dentro del repo edugo-dev-environment)   │
│ 📂 Ruta: /repos-separados/edugo-dev-environment/docs/cicd/        │
│ ⚠️ NO uses archivos de la raíz del repo (son viejos)              │
│ ✅ SOLO usa archivos dentro de docs/cicd/                         │
└─────────────────────────────────────────────────────────────────────┘
```

**Última actualización:** 20 Nov 2025, 18:30 hrs

---

## 🗺️ MAPA DE UBICACIÓN

```
edugo-dev-environment/ (repositorio de infraestructura)
│
├── [scripts/, docker/, configs/, etc.]          ← Infraestructura
├── README.md                                    ← README del proyecto
├── docker-compose.yml                           ← Compose principal
│
└── docs/
    └── cicd/                                    ← 👉 ESTÁS AQUÍ
        ├── START-HERE.md                        ← Este archivo
        ├── INDEX.md                             ← Navegación completa
        ├── PROMPTS.md                           ← Prompts para cada fase
        ├── README.md                            ← Plan de CI/CD
        ├── docs/                                ← Documentación
        ├── sprints/                             ← ⭐ Planes de sprint
        │   ├── SPRINT-1-TASKS.md
        │   ├── SPRINT-3-TASKS.md (si existe)
        │   └── SPRINT-4-TASKS.md
        ├── tracking/                            ← Estado y seguimiento
        │   ├── SPRINT-STATUS.md                 ← Estado actual
        │   └── REGLAS.md
        └── assets/
```

---

## 🎯 ¿QUÉ QUIERES HACER?

### 🔍 Opción 1: Ver Estado Actual del Proyecto
```bash
cat docs/cicd/tracking/SPRINT-STATUS.md | head -40
```

Lee el archivo para saber:
- Sprint activo
- Fase actual (1, 2, o 3)
- Progreso
- Próxima tarea

### ▶️ Opción 2: Continuar Trabajo desde donde quedó

**Prompt a usar:**
```
Continúa el trabajo de CI/CD en edugo-dev-environment desde donde quedó.
```

Ver detalles en: [PROMPTS.md](PROMPTS.md#continuar-desde-donde-quedó)

### 🆕 Opción 3: Iniciar Nuevo Sprint

**Prompt a usar:**
```
Ejecuta FASE 1 del SPRINT-X en edugo-dev-environment.
```

Reemplaza X con: 1, 3, o 4  
Ver detalles en: [PROMPTS.md](PROMPTS.md#fase-1)

### 📚 Opción 4: Entender el Sistema Completo

**Lee en orden:**
1. [INDEX.md](INDEX.md) - Navegación general (5 min)
2. [docs/QUICK-START.md](docs/QUICK-START.md) - Guía rápida (10 min)
3. [tracking/REGLAS.md](tracking/REGLAS.md) - Reglas detalladas (15 min)

---

## 📍 NAVEGACIÓN RÁPIDA

| Quiero... | Abrir... |
|-----------|----------|
| 🗺️ Navegar el proyecto | [INDEX.md](INDEX.md) |
| 🎯 Prompts para ejecutar | [PROMPTS.md](PROMPTS.md) ⭐ |
| 📊 Estado actual | [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md) |
| 📜 Reglas de ejecución | [tracking/REGLAS.md](tracking/REGLAS.md) |
| 📖 Guía rápida | [docs/QUICK-START.md](docs/QUICK-START.md) |
| 🎯 Ver tareas del sprint | [sprints/](sprints/) |
| 📈 Ver progreso | [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md) |

---

## 🤖 PARA CLAUDE CODE (INSTRUCCIONES CRÍTICAS)

### ⚠️ Antes de Hacer CUALQUIER COSA:

1. **Lee SIEMPRE:** `docs/cicd/INDEX.md`
2. **Verifica ubicación:**
   ```bash
   pwd
   # Debe contener: /edugo-dev-environment/docs/cicd
   ```
3. **Lee estado:** `docs/cicd/tracking/SPRINT-STATUS.md`
4. **Identifica:**
   - Sprint activo
   - Fase actual
   - Próxima tarea

### ⚠️ NO Uses Archivos Fuera de docs/cicd/

**Archivos PROHIBIDOS:**
- ❌ Cualquier archivo en raíz del repo (viejos)
- ❌ `/docs/isolated/*` (otra carpeta)
- ❌ Archivos de otros proyectos

**Archivos PERMITIDOS:**
- ✅ `docs/cicd/sprints/SPRINT-X-TASKS.md`
- ✅ `docs/cicd/tracking/*`
- ✅ `docs/cicd/docs/*`

### ⚠️ Cómo Verificar que Estás en el Archivo Correcto:

```bash
# Al abrir un archivo de sprint, verifica:
readlink -f docs/cicd/sprints/SPRINT-1-TASKS.md
# Debe mostrar: .../edugo-dev-environment/docs/cicd/sprints/SPRINT-1-TASKS.md

# Si muestra otra ruta, estás en el lugar equivocado
```

---

## 🔗 Enlaces Importantes

- **Plan general del proyecto:** [README.md](README.md)
- **Navegación completa:** [INDEX.md](INDEX.md)
- **Prompts para ejecutar:** [PROMPTS.md](PROMPTS.md) ⭐
- **Estado en tiempo real:** [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md)
- **Reglas de ejecución:** [tracking/REGLAS.md](tracking/REGLAS.md)

---

## 📊 COMANDOS RÁPIDOS

### Ver estado actual:
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment
cat docs/cicd/tracking/SPRINT-STATUS.md | head -40
```

### Ver próxima tarea:
```bash
grep "⏳\|🔄" docs/cicd/tracking/SPRINT-STATUS.md | head -1
```

### Ver sprints disponibles:
```bash
ls -1 docs/cicd/sprints/
```

### Ver logs de sesiones anteriores:
```bash
ls -lt docs/cicd/tracking/logs/ | head -5
```

---

**Generado por:** Claude Code  
**Fecha:** 20 de Noviembre, 2025  
**Versión:** 2.0 (con sistema de prompts)
