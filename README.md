# EduGo - Ambiente de Desarrollo Backend

**Para Desarrolladores Frontend** 🎨

Este repositorio te permite levantar **todo el backend de EduGo** en tu Mac con un solo comando.

---

## 🎯 ¿Qué Obtienes?

Después de seguir esta guía (5-10 minutos), tendrás corriendo:

- ✅ **API Mobile** en http://localhost:8081
- ✅ **API Administración** en http://localhost:8082  
- ✅ **Worker** procesando PDFs en background
- ✅ **PostgreSQL** con datos de prueba
- ✅ **MongoDB** para almacenar documentos
- ✅ **RabbitMQ** para mensajería

**Todo funcional y listo para conectar tu app frontend.**

---

## ⚡ Inicio Rápido (3 Pasos)

### Paso 1: Instalar Docker Desktop

```bash
# macOS
brew install --cask docker

# O descarga desde: https://www.docker.com/products/docker-desktop
```

**Abrir Docker Desktop** y esperar a que inicie (ver ícono en barra de menú).

---

### Paso 2: Clonar y Configurar

```bash
# Clonar este repositorio
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# Ejecutar setup (pedirá tu GitHub token)
./scripts/setup.sh
```

**Cuando pida credenciales:**
- **Usuario:** tu-usuario-github  
- **Token:** [Crear token aquí](https://github.com/settings/tokens) con scope `read:packages`

---

### Paso 3: ¡Listo! 🎉

Verifica que todo esté corriendo:

```bash
# Ver servicios
cd docker && docker-compose ps

# Probar API Mobile
curl http://localhost:8081/health

# Probar API Admin
curl http://localhost:8082/health
```

**Resultado esperado:**
```json
{
  "status": "ok",
  "database": "connected",
  "mongodb": "connected",
  "rabbitmq": "connected"
}
```

---

## 📡 Endpoints de las APIs

### API Mobile (Puerto 8081)

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/v1/auth/login` | POST | Login de usuarios |
| `/api/v1/auth/register` | POST | Registro de usuarios |
| `/api/v1/users` | GET | Listar usuarios |
| `/api/v1/courses` | GET | Listar cursos |
| `/api/v1/documents` | POST | Subir PDF |

**Documentación completa:** http://localhost:8081/swagger

### API Administración (Puerto 8082)

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/v1/admin/users` | GET | Gestionar usuarios |
| `/api/v1/admin/institutions` | GET | Gestionar instituciones |
| `/api/v1/admin/reports` | GET | Reportes |

**Documentación completa:** http://localhost:8082/swagger

---

## 🧪 Datos de Prueba

El ambiente viene con datos de prueba pre-cargados:

### Usuarios de Prueba

| Email | Password | Rol |
|-------|----------|-----|
| `admin@edugo.com` | `admin123` | Administrador |
| `profesor@edugo.com` | `profesor123` | Profesor |
| `estudiante@edugo.com` | `estudiante123` | Estudiante |

### Ejemplo: Login desde tu Frontend

```javascript
// React / Vue / Angular
const response = await fetch('http://localhost:8081/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'estudiante@edugo.com',
    password: 'estudiante123'
  })
});

const { token, user } = await response.json();
console.log('Token JWT:', token);
console.log('Usuario:', user);
```

---

## 🔄 Comandos Útiles

### Iniciar el Backend

```bash
cd docker
docker-compose up -d
```

### Detener el Backend

```bash
docker-compose stop
```

### Ver Logs (Debugging)

```bash
# Todos los servicios
docker-compose logs -f

# Solo API Mobile
docker-compose logs -f api-mobile

# Solo Worker
docker-compose logs -f worker
```

### Reiniciar un Servicio

```bash
docker-compose restart api-mobile
```

### Reset Completo (Borra datos)

```bash
docker-compose down -v
./scripts/setup.sh
```

---

## 🔌 Conectar Tu Frontend

### React / Next.js

```javascript
// lib/api.js
const API_BASE_URL = 'http://localhost:8081/api/v1';

export async function loginUser(email, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  return response.json();
}

export async function getCourses(token) {
  const response = await fetch(`${API_BASE_URL}/courses`, {
    headers: { 
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}
```

### Vue.js

```javascript
// services/api.js
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8081/api/v1'
});

export default {
  login(email, password) {
    return api.post('/auth/login', { email, password });
  },
  
  getCourses(token) {
    return api.get('/courses', {
      headers: { Authorization: `Bearer ${token}` }
    });
  }
};
```

### Angular

```typescript
// services/api.service.ts
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private apiUrl = 'http://localhost:8081/api/v1';

  constructor(private http: HttpClient) {}

  login(email: string, password: string) {
    return this.http.post(`${this.apiUrl}/auth/login`, { email, password });
  }

  getCourses(token: string) {
    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    return this.http.get(`${this.apiUrl}/courses`, { headers });
  }
}
```

---

## 🐛 Problemas Comunes

### "Cannot connect to Docker daemon"

**Solución:** Abre Docker Desktop y espera a que inicie.

```bash
open -a Docker
```

### "Port 8081 already in use"

**Solución:** Algo está usando el puerto. Detenerlo:

```bash
lsof -ti:8081 | xargs kill -9
```

### "pull access denied"

**Solución:** Autenticarte en GitHub Container Registry:

```bash
docker login ghcr.io
# Usuario: tu-github-username
# Password: tu-personal-access-token
```

### "API responde 500"

**Solución:** Ver logs del servicio:

```bash
cd docker
docker-compose logs api-mobile
```

### Más Problemas?

Ver [Troubleshooting Completo](#-troubleshooting-detallado) al final de este README.

---

## 🎨 Ejemplo Completo: App de Login

```javascript
// App.jsx (React)
import { useState } from 'react';

function App() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [user, setUser] = useState(null);
  const [error, setError] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    
    try {
      const response = await fetch('http://localhost:8081/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      
      if (!response.ok) throw new Error('Login falló');
      
      const data = await response.json();
      setUser(data.user);
      localStorage.setItem('token', data.token);
    } catch (err) {
      setError(err.message);
    }
  };

  if (user) {
    return (
      <div>
        <h1>Bienvenido, {user.firstName}!</h1>
        <p>Email: {user.email}</p>
        <p>Rol: {user.role}</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleLogin}>
      <h1>Login EduGo</h1>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      
      <input 
        type="email" 
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
      />
      
      <input 
        type="password" 
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      
      <button type="submit">Iniciar Sesión</button>
      
      <p>Usuario de prueba: estudiante@edugo.com / estudiante123</p>
    </form>
  );
}

export default App;
```

---

## 📱 Testing con Postman

1. **Importar colección:**
   - Archivo: `docs/postman/EduGo-APIs.postman_collection.json` (si existe)
   - O crear requests manualmente

2. **Request de ejemplo:**

```
POST http://localhost:8081/api/v1/auth/login
Content-Type: application/json

{
  "email": "estudiante@edugo.com",
  "password": "estudiante123"
}
```

3. **Guardar el token** y usarlo en otros requests:

```
GET http://localhost:8081/api/v1/courses
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 🗄️ Acceso Directo a Bases de Datos

### PostgreSQL (si necesitas consultar directo)

```bash
docker exec -it postgres psql -U edugo -d edugo
```

**Comandos útiles:**
```sql
-- Ver todas las tablas
\dt

-- Ver usuarios
SELECT * FROM users;

-- Ver cursos
SELECT * FROM courses;

-- Salir
\q
```

### MongoDB (para documentos/PDFs)

```bash
docker exec -it mongodb mongosh -u edugo -p edugo123 edugo
```

**Comandos útiles:**
```javascript
// Ver colecciones
show collections

// Ver documentos procesados
db.documents.find().pretty()

// Salir
exit
```

### RabbitMQ UI (ver mensajes en cola)

Abrir en navegador: http://localhost:15672

**Credenciales:**
- Usuario: `edugo`
- Password: `edugo123`

---

## 📚 Documentación Adicional

Si necesitas más detalles:

| Documento | Para Qué |
|-----------|----------|
| [Quick Start](docker/QUICK_START.md) | Setup rápido con opciones |
| [Ejemplo End-to-End](docs/EXAMPLE.md) | Tutorial completo paso a paso |
| [Guía Completa](docker/README.md) | Documentación técnica detallada |
| [Scripts](scripts/README.md) | Referencia de scripts disponibles |

---

## 🏗️ Arquitectura (Para los Curiosos)

```
Tu App Frontend (React/Vue/Angular)
        │
        │ HTTP REST
        ↓
┌─────────────────────────────────────┐
│   Docker Compose (Este Repo)       │
│                                     │
│  API Mobile (:8081) ─┐              │
│                      │              │
│  API Admin  (:8082) ─┼─→ PostgreSQL │
│                      │              │
│  Worker      ────────┼─→ MongoDB    │
│                      │              │
│  RabbitMQ     ───────┘              │
└─────────────────────────────────────┘
```

**Flujo típico:**
1. Tu frontend hace login en API Mobile (puerto 8081)
2. Recibe un token JWT
3. Usa el token para obtener cursos, usuarios, etc.
4. Cuando subes un PDF, el Worker lo procesa automáticamente

---

## 🛑 Detener y Limpiar

### Detener (Mantiene Datos)

```bash
cd docker
docker-compose stop
```

Próxima vez que hagas `docker-compose up -d`, todo sigue donde lo dejaste.

### Reset Completo (Borra Todo)

```bash
cd docker
docker-compose down -v

# Re-inicializar
cd ..
./scripts/setup.sh
```

---

## ❓ FAQ Rápido

**Q: ¿Necesito saber Go/Backend para usar esto?**  
A: No. Solo ejecuta el setup y usa las APIs desde tu frontend.

**Q: ¿Puedo cambiar los puertos?**  
A: Sí, edita `docker/.env` y cambia `API_MOBILE_PORT=8081` a otro puerto.

**Q: ¿Los datos se pierden al detener Docker?**  
A: No, se mantienen en volúmenes. Solo se borran con `docker-compose down -v`.

**Q: ¿Puedo trabajar offline?**  
A: Después del primer setup, sí. Las imágenes quedan en tu Mac.

**Q: ¿Cómo actualizo las APIs a la última versión?**  
A: `./scripts/update-images.sh` y reinicia con `docker-compose up -d`.

---

## 🤝 Soporte

### Algo no funciona?

1. **Ver logs:** `docker-compose logs -f`
2. **Reiniciar:** `docker-compose restart api-mobile`
3. **Reset completo:** `docker-compose down -v && ./scripts/setup.sh`

### Errores comunes resueltos

Ver sección [Troubleshooting Detallado](#-troubleshooting-detallado) abajo.

### Reportar un bug

[Crear issue](https://github.com/EduGoGroup/edugo-dev-environment/issues/new)

---

## 🔧 Troubleshooting Detallado

### Error: "Cannot connect to Docker daemon"

**Causa:** Docker Desktop no está corriendo.

**Solución:**
```bash
open -a Docker
# Esperar a que el ícono aparezca en la barra de menú
docker ps  # Verificar que funciona
```

---

### Error: "Port already in use"

**Causa:** Otro servicio usa el puerto 5432, 8081, etc.

**Solución:**
```bash
# Ver qué usa el puerto
lsof -ti:5432  # PostgreSQL
lsof -ti:8081  # API Mobile
lsof -ti:8082  # API Admin

# Matar el proceso
lsof -ti:8081 | xargs kill -9

# O cambiar puerto en docker/.env
echo "API_MOBILE_PORT=8083" >> docker/.env
```

---

### Error: "pull access denied for ghcr.io/edugogroup/..."

**Causa:** No estás autenticado en GitHub Container Registry.

**Solución:**
```bash
# Crear token en: https://github.com/settings/tokens
# Scope: read:packages

# Login
docker login ghcr.io
Username: tu-usuario-github
Password: ghp_tu_token_aqui

# Re-ejecutar setup
./scripts/setup.sh
```

---

### Error: API responde "dial tcp: connection refused"

**Causa:** PostgreSQL o MongoDB no están corriendo.

**Solución:**
```bash
# Ver estado
cd docker
docker-compose ps

# Si postgres/mongodb están "Exited", reiniciar
docker-compose up -d postgres mongodb

# Ver logs
docker-compose logs postgres mongodb
```

---

### Error: "relation 'users' does not exist"

**Causa:** Migraciones no se ejecutaron.

**Solución:**
```bash
# Ejecutar migraciones manualmente
docker-compose run --rm migrator

# Verificar tablas
docker exec -it postgres psql -U edugo -d edugo -c "\dt"
```

---

### Error: Worker no procesa PDFs

**Causa:** RabbitMQ no está conectado o falta OPENAI_API_KEY.

**Solución:**
```bash
# Ver logs del worker
docker-compose logs -f worker

# Verificar RabbitMQ
docker-compose ps rabbitmq
open http://localhost:15672  # UI

# Verificar variables de entorno
grep OPENAI_API_KEY docker/.env

# Si falta, agregarla
echo "OPENAI_API_KEY=sk-tu-key-aqui" >> docker/.env
docker-compose restart worker
```

---

### Error: "no space left on device"

**Causa:** Docker llenó tu disco.

**Solución:**
```bash
# Ver uso de Docker
docker system df

# Limpiar imágenes viejas
docker image prune -a

# Limpiar volúmenes sin usar (⚠️ cuidado)
docker volume prune

# Limpieza completa
docker system prune -a --volumes
```

---

### API responde 500 Internal Server Error

**Causa:** Error en el código del backend o BD no disponible.

**Solución:**
```bash
# Ver logs detallados
docker-compose logs -f api-mobile

# Verificar conectividad a BD
docker-compose exec api-mobile ping postgres
docker-compose exec api-mobile ping mongodb

# Reiniciar API
docker-compose restart api-mobile

# Si persiste, ver variables
docker-compose exec api-mobile env | grep DATABASE
```

---

## 🎓 Para Saber Más

### ¿Por Qué NO Hay CI/CD en Este Repo?

Este proyecto es **configuración Docker**, no código fuente. La validación se hace localmente en segundos, no necesita CI/CD.

**Más detalles:** Ver sección completa en la [documentación técnica](docs/cicd/README.md).

### ¿Quieres Contribuir?

1. Fork el repo
2. Crea branch: `git checkout -b feature/mi-mejora`
3. Haz cambios
4. Valida: `./scripts/validate.sh`
5. Push y crea PR

### Scripts Disponibles

- `./scripts/setup.sh` - Setup inicial
- `./scripts/validate.sh` - Validar docker-compose
- `./scripts/update-images.sh` - Actualizar imágenes
- `./scripts/cleanup.sh` - Limpiar ambiente
- `./scripts/stop.sh` - Detener servicios

**Documentación completa:** [scripts/README.md](scripts/README.md)

---

## 📝 Credenciales por Defecto

### ⚠️ SOLO PARA DESARROLLO LOCAL

**PostgreSQL:**
```
Host: localhost
Port: 5432
User: edugo
Password: edugo123
Database: edugo
```

**MongoDB:**
```
Host: localhost
Port: 27017
User: edugo
Password: edugo123
Database: edugo
```

**RabbitMQ:**
```
Host: localhost
AMQP Port: 5672
Management UI: 15672
User: edugo
Password: edugo123
```

**Usuarios de Prueba:**
```
admin@edugo.com / admin123
profesor@edugo.com / profesor123
estudiante@edugo.com / estudiante123
```

---

## 📦 ¿Qué Contiene Este Repo?

```
edugo-dev-environment/
├── docker/
│   ├── docker-compose.yml      ← Configuración principal
│   └── .env                    ← Variables de entorno
├── scripts/
│   ├── setup.sh                ← Setup automático
│   ├── validate.sh             ← Validar configuración
│   └── ...
├── docs/
│   └── EXAMPLE.md              ← Tutorial completo
└── README.md                   ← Este archivo
```

---

## 🚀 ¡Comienza Ahora!

```bash
# 1. Clonar
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# 2. Setup
./scripts/setup.sh

# 3. Verificar
curl http://localhost:8081/health

# 4. ¡A programar tu frontend! 🎨
```

---

**Última actualización:** 22 de Noviembre, 2025  
**Versión:** 2.0.0  
**Mantenedor:** Equipo EduGo

**¿Dudas?** Abre un [issue](https://github.com/EduGoGroup/edugo-dev-environment/issues) o consulta la [documentación completa](docs/EXAMPLE.md).
