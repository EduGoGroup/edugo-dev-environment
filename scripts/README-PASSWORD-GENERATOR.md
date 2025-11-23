# 🔐 Generador de Hashes Bcrypt

Este script te permite generar hashes bcrypt para contraseñas de prueba.

## 📋 Uso

```bash
./scripts/generate-password.sh <tu-password>
```

## 📝 Ejemplos

### Generar hash para una contraseña

```bash
./scripts/generate-password.sh mipassword123
```

**Salida:**
```
🔐 Generando hash bcrypt...

Password: mipassword123
Hash:     $2a$10$xYz123abcDEF456ghiJKL789mnoPQR012stuVWX345yzABC678def

✅ Hash generado exitosamente

💡 Puedes usar este hash en:
   - Migraciones SQL (columna password_hash)
   - Tests de integración
   - Datos de prueba
```

### Generar hash para usuario de prueba

```bash
./scripts/generate-password.sh edugo2024
```

## 🎯 Casos de Uso

### 1. Crear usuarios de prueba en PostgreSQL

```sql
INSERT INTO users (email, password_hash, role) VALUES
('test@edugo.com', '$2a$10$...hash-generado...', 'student');
```

### 2. Testing de autenticación

```javascript
// En tus tests
const testUser = {
  email: 'test@example.com',
  password: 'mipassword123' // Password original
};

// El hash $2a$10$... ya está en la BD
const response = await api.post('/login', testUser);
```

### 3. Seeds de datos

```bash
# Generar hash
./scripts/generate-password.sh test123

# Copiar el hash a tu archivo de seeds
# seeds/postgresql/02_users.sql
```

## 🔧 Cómo Funciona

1. Toma el password como argumento
2. Usa el módulo `golang.org/x/crypto/bcrypt` del migrator
3. Genera un hash con costo por defecto (10)
4. Muestra el resultado en consola

## ⚠️ Notas Importantes

- **Costo 10**: El hash usa bcrypt con costo 10 (balance seguridad/rendimiento)
- **No reproducible**: Cada ejecución genera un hash diferente (bcrypt usa salt)
- **Mismo password**: Ambos hashes son válidos para verificar el mismo password
- **Solo desarrollo**: Estos hashes son para entornos de desarrollo/testing

## 🚀 Ejemplo Completo

Crear un usuario de prueba completo:

```bash
# 1. Generar hash
./scripts/generate-password.sh mydevpassword

# Resultado: $2a$10$ABC123...

# 2. Insertar en PostgreSQL
docker exec -it edugo-postgres psql -U edugo -d edugo

# 3. Ejecutar SQL
INSERT INTO users (id, email, password_hash, role, first_name, last_name) 
VALUES (
  gen_random_uuid(),
  'developer@test.com',
  '$2a$10$ABC123...',  -- El hash generado
  'student',
  'Dev',
  'Tester'
);

# 4. Probar login desde tu frontend
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"developer@test.com","password":"mydevpassword"}'
```

## 📚 Recursos

- [bcrypt en Go](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [Documentación de bcrypt](https://en.wikipedia.org/wiki/Bcrypt)
