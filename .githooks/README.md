# Git Hooks - EduGo Dev Environment

Este directorio contiene hooks opcionales de Git para mejorar la calidad del código y prevenir errores comunes.

---

## 🪝 Hooks Disponibles

### pre-commit

Valida configuración antes de permitir commits.

**¿Qué Valida?**

- ✅ Sintaxis de archivos docker-compose (si fueron modificados)
- ✅ Que archivos `.env` NO se commiteen accidentalmente
- ✅ Que archivos de credenciales NO se commiteen
- ✅ Que scripts `.sh` tengan permisos de ejecución

---

## 🚀 Activar Hooks

### Opción A: Configurar Git (Recomendado)

```bash
# Desde la raíz del proyecto
git config core.hooksPath .githooks
```

**Ventajas:**
- ✅ Automático para todos los commits
- ✅ Fácil de activar/desactivar
- ✅ No modifica `.git/hooks`

### Opción B: Symlink

```bash
# Crear enlace simbólico
ln -s ../../.githooks/pre-commit .git/hooks/pre-commit
```

**Ventajas:**
- ✅ Compatible con versiones antiguas de Git
- ❌ Requiere hacerlo en cada clone

---

## 🛑 Desactivar Hooks

### Desactivar Temporalmente (Un Commit)

Si necesitas hacer un commit urgente sin validación:

```bash
git commit --no-verify
```

**⚠️ ADVERTENCIA:** Solo usa `--no-verify` si estás seguro de lo que haces.

### Desactivar Permanentemente

```bash
git config --unset core.hooksPath
```

---

## 🧪 Probar Hooks

### Probar pre-commit hook

```bash
# 1. Modificar un docker-compose con error
echo "invalid yaml:" >> docker/docker-compose.yml

# 2. Intentar commit
git add docker/docker-compose.yml
git commit -m "test"

# Esperado: Commit bloqueado con mensaje de error

# 3. Revertir cambio
git checkout docker/docker-compose.yml
```

### Probar protección de .env

```bash
# 1. Intentar agregar .env
touch docker/.env
git add docker/.env
git commit -m "test"

# Esperado: Commit bloqueado con advertencia

# 2. Limpiar
git reset HEAD docker/.env
rm docker/.env
```

---

## 📋 Detalles del Pre-commit Hook

### Flujo de Validación

```
┌─────────────────────────┐
│ git commit              │
└──────────┬──────────────┘
           │
           ↓
┌─────────────────────────┐
│ pre-commit hook         │
│ se ejecuta              │
└──────────┬──────────────┘
           │
           ↓
    ┌──────┴───────┐
    │              │
    ↓              ↓
┌────────┐    ┌────────┐
│ Valida │    │ Valida │
│ YAML   │    │ .env   │
└────┬───┘    └───┬────┘
     │            │
     ↓            ↓
  ┌─────────────────┐
  │ ¿Todo OK?       │
  └────┬───────┬────┘
       │       │
    SÍ │       │ NO
       ↓       ↓
┌─────────┐ ┌──────────┐
│ Commit  │ │ Bloquear │
│ OK ✅   │ │ commit ❌│
└─────────┘ └──────────┘
```

### Archivos Protegidos

El hook previene el commit de:

| Archivo/Patrón | Razón |
|----------------|-------|
| `docker/.env` | Contiene credenciales |
| `.env.local` | Configuración local |
| `.env.production` | Credenciales de producción |
| `credentials.json` | Credenciales de servicios |
| `serviceAccount.json` | Credenciales de Google Cloud |

### Validaciones Automáticas

| Validación | Cuándo | Acción |
|------------|--------|--------|
| Sintaxis docker-compose | Si `docker/*.yml` modificado | Ejecuta `./scripts/validate.sh` |
| Archivos .env | Si `.env` en staged | Bloquea commit |
| Permisos de scripts | Si `scripts/*.sh` modificado | Agrega permisos `+x` |

---

## 🔧 Personalizar Hooks

### Agregar Validaciones Personalizadas

Edita `.githooks/pre-commit`:

```bash
# Tu validación personalizada
echo "🔍 Validando código personalizado..."

if ! tu-comando-validacion; then
    echo "❌ Validación falló"
    exit 1
fi
```

### Agregar Más Archivos Protegidos

Edita la sección `CREDENTIAL_FILES`:

```bash
CREDENTIAL_FILES=(
    ".env"
    ".env.local"
    "tu-archivo-secreto.key"  # Agregar aquí
)
```

---

## ❓ FAQ

### ¿Por qué el hook bloquea mi commit?

**Respuesta:** El hook detectó un problema:
- Sintaxis YAML inválida en docker-compose
- Intentando commitear archivos con credenciales
- Script sin permisos de ejecución

**Solución:** Corrige el problema o usa `--no-verify` si estás seguro.

### ¿Cómo saber qué validó el hook?

**Respuesta:** El hook imprime mensajes detallados:

```bash
🔍 Ejecutando validaciones pre-commit...

📝 Archivos docker-compose modificados, validando...
✅ docker-compose válido

✅ Todas las validaciones pasaron
```

### ¿Puedo usar estos hooks en otros proyectos?

**Respuesta:** Sí, copia el directorio `.githooks` y configura:

```bash
git config core.hooksPath .githooks
```

### ¿Los hooks se comparten con el equipo?

**Respuesta:** Los hooks en `.githooks/` se comparten vía Git.
Cada desarrollador debe activarlos con:

```bash
git config core.hooksPath .githooks
```

---

## 🆘 Troubleshooting

### Hook no se ejecuta

**Problema:** Hago commit pero el hook no corre.

**Solución:**
```bash
# Verificar configuración
git config core.hooksPath

# Debería mostrar: .githooks

# Si no, activar:
git config core.hooksPath .githooks
```

### Hook da error de permisos

**Problema:** `Permission denied: .githooks/pre-commit`

**Solución:**
```bash
chmod +x .githooks/pre-commit
```

### Quiero deshabilitar temporalmente

**Problema:** Necesito hacer commit rápido sin validación.

**Solución:**
```bash
git commit --no-verify -m "mensaje"
```

---

## 📝 Mejores Prácticas

1. **Activar hooks al clonar** - Primera acción después de clone
2. **No usar --no-verify habitualmente** - Solo en emergencias
3. **Revisar mensajes del hook** - Entender por qué falla
4. **Mantener hooks actualizados** - Pull regularmente
5. **Reportar problemas** - Si el hook falla incorrectamente

---

## 🔗 Referencias

- [Git Hooks Documentation](https://git-scm.com/docs/githooks)
- [Core.hooksPath Config](https://git-scm.com/docs/git-config#Documentation/git-config.txt-corehooksPath)

---

**Última actualización:** 22 de Noviembre, 2025  
**Versión:** 1.0  
**Mantenedor:** Equipo EduGo
