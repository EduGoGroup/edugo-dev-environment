# 📋 Plan de Acción Post-Merge

**PR Relacionado**: https://github.com/EduGoGroup/edugo-infrastructure/pull/32  
**Rama**: `fix/implement-real-demo-user-passwords`

---

## 🔄 Flujo Completo

### ✅ FASE 1: COMPLETADA
- [x] Implementar passwords reales en infraestructura
- [x] Crear tablas refresh_tokens y login_attempts
- [x] Crear PR en edugo-infrastructure
- [x] Probar funcionamiento completo desde cero
- [x] Crear herramientas (script generador de passwords)
- [x] Actualizar documentación en dev-environment

---

### ⏳ FASE 2: PENDIENTE (Después del merge del PR)

#### 1. Crear Release en edugo-infrastructure

**Ejecutar**:
```bash
cd edugo-infrastructure
git checkout main
git pull origin main

# Crear tag
git tag -a v0.10.0 -m "Release v0.10.0

## Nuevas Features
- Passwords reales para usuarios demo (password: edugo2024)
- Tabla refresh_tokens para gestión de sesiones JWT
- Tabla login_attempts para rate limiting

## Fixes
- Login completo funcional sin errores
- Token storage implementado
"

# Push del tag
git push origin v0.10.0
```

**En GitHub**:
1. Ir a: https://github.com/EduGoGroup/edugo-infrastructure/releases/new
2. Tag: `v0.10.0`
3. Título: `v0.10.0 - Auth Tables & Real Demo Passwords`
4. Descripción: [Copiar del tag]
5. Publicar release

---

#### 2. Actualizar dev-environment

**A. Actualizar go.mod del migrator**:
```bash
cd edugo-dev-environment/migrator

# Actualizar a nueva versión
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.10.0
go mod tidy
```

**B. Revertir cambio temporal en cmd/main.go**:

Eliminar la parte que hace checkout a la rama fix y dejar solo:
```go
// Actualizar si ya existe
fmt.Println("Actualizando edugo-infrastructure...")
cmd := exec.Command("git", "pull", "origin", "main")
cmd.Dir = infraDir
// ...
```

**C. Reconstruir y probar**:
```bash
cd docker
docker-compose down -v
docker-compose build migrator
docker-compose up -d

# Esperar y probar
curl -X POST http://localhost:8081/v1/auth/login \
  -d '{"email":"admin@edugo.test","password":"edugo2024"}'
```

**D. Commit final**:
```bash
git add migrator/
git commit -m "chore: actualizar infraestructura a v0.10.0 y revertir cambio temporal

- Actualizar go.mod a edugo-infrastructure v0.10.0
- Revertir checkout temporal a rama fix
- Migrator ahora usa versión oficial del módulo

Relacionado: https://github.com/EduGoGroup/edugo-infrastructure/pull/32"
```

---

#### 3. Actualizar edugo-api-mobile

```bash
cd edugo-api-mobile

# Actualizar infraestructura
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.10.0
go mod tidy

# Probar
go test ./...
go run ./cmd/api

# Commit
git commit -am "chore: actualizar infraestructura a v0.10.0"
```

---

#### 4. Actualizar edugo-api-admin

```bash
cd edugo-api-admin

# Actualizar infraestructura
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0
go mod tidy

# Probar
go test ./...
go run ./cmd/api

# Commit
git commit -am "chore: actualizar infraestructura a v0.10.0"
```

---

#### 5. Actualizar edugo-worker

```bash
cd edugo-worker

# Actualizar infraestructura (si la usa)
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.10.0
go mod tidy

# Probar
go test ./...
go run ./cmd/worker

# Commit
git commit -am "chore: actualizar infraestructura a v0.10.0"
```

---

## 📝 Checklist Post-Release

### edugo-infrastructure
- [ ] PR mergeado a main
- [ ] Tag v0.10.0 creado
- [ ] Release publicado en GitHub
- [ ] Módulos Go disponibles

### edugo-dev-environment
- [ ] go.mod actualizado a v0.10.0
- [ ] Cambio temporal revertido
- [ ] Probado desde cero
- [ ] Commit y push

### edugo-api-mobile
- [ ] go.mod actualizado a v0.10.0
- [ ] Tests pasando
- [ ] Commit y push
- [ ] Nueva imagen Docker publicada

### edugo-api-admin
- [ ] go.mod actualizado a v0.10.0
- [ ] Tests pasando
- [ ] Commit y push
- [ ] Nueva imagen Docker publicada

### edugo-worker
- [ ] go.mod actualizado (si aplica)
- [ ] Tests pasando
- [ ] Commit y push
- [ ] Nueva imagen Docker publicada

---

## 🎯 Resultado Final Esperado

### Para Desarrolladores Frontend
```javascript
// Login funcionando al 100%
const response = await fetch('http://localhost:8081/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'student1@edugo.test',
    password: 'edugo2024'
  })
});

const { access_token, refresh_token, user } = await response.json();
// ✅ Funciona sin errores
// ✅ Tokens válidos
// ✅ Usuario disponible
```

### Para DevOps/Backend
- ✅ Migraciones automáticas completas
- ✅ Base de datos con todas las tablas necesarias
- ✅ Ambiente de desarrollo sin configuración manual
- ✅ Usuarios de prueba funcionales
- ✅ Rate limiting activo
- ✅ Auditoría de login implementada

---

## ⚠️ Notas Importantes

1. **NO hacer push del cambio temporal**: El commit 0c48cd3 es temporal y se revertirá
2. **Esperar el release antes de actualizar APIs**: Asegurarse que v0.10.0 esté publicado
3. **Probar cada proyecto individualmente**: Después de actualizar go.mod
4. **Coordinar imágenes Docker**: Las APIs necesitarán rebuild y republish

---

## 📞 Contacto

Si hay algún problema durante el proceso:
- Revisar logs del migrator: `docker-compose logs migrator`
- Verificar tablas: `docker exec edugo-postgres psql -U edugo -d edugo -c "\\dt"`
- Probar login: `curl -X POST http://localhost:8081/v1/auth/login -d '...'`

---

**Creado**: 2025-11-22  
**Estado actual**: ✅ PR Creado, esperando merge  
**Próximo paso**: Mergear PR #32 y crear release v0.10.0
