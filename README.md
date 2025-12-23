# EduGo - Ambiente de Desarrollo

Repositorio para levantar el backend completo de EduGo localmente con Docker.

## Inicio Rápido

### Prerequisitos

- [Docker Desktop](https://www.docker.com/products/docker-desktop) instalado y corriendo
- [GitHub Personal Access Token](https://github.com/settings/tokens) con scope `read:packages`

### Setup (3 comandos)

```bash
# 1. Clonar
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# 2. Setup automático
./scripts/setup.sh

# 3. Verificar
curl http://localhost:8081/health
```

---

## Servicios Disponibles

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **API Mobile** | http://localhost:8081 | API principal |
| **API Admin** | http://localhost:8082 | API de administración |
| **Swagger Mobile** | http://localhost:8081/swagger/index.html | Documentación API |
| **Swagger Admin** | http://localhost:8082/swagger/index.html | Documentación API |
| **RabbitMQ UI** | http://localhost:15672 | Gestión de colas (edugo/edugo123) |

---

## Usuarios de Prueba

**Contraseña para todos:** `edugo2024`

| Email | Rol |
|-------|-----|
| admin@edugo.test | Administrador |
| teacher.math@edugo.test | Profesor |
| student1@edugo.test | Estudiante |

**Ejemplo de login:**
```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "student1@edugo.test", "password": "edugo2024"}'
```

---

## Comandos Útiles

```bash
make help        # Ver todos los comandos
make up          # Iniciar servicios
make down        # Detener servicios
make logs        # Ver logs
make diagnose    # Diagnóstico del ambiente
make health      # Verificar APIs
make psql        # Conectar a PostgreSQL
make mongo       # Conectar a MongoDB
```

---

## Estructura del Proyecto

```
edugo-dev-environment/
├── Makefile             # Comandos simplificados (make help)
├── docker/              # Configuración Docker Compose
├── scripts/             # Scripts de utilidad
├── migrator/            # Herramienta de migraciones
├── seeds/               # Datos de prueba
├── documentos/          # Documentación detallada
└── README.md            # Este archivo
```

---

## Documentación

Documentación completa en [`documentos/`](./documentos/):

- [Arquitectura](./documentos/ARQUITECTURA.md) - Diagrama y componentes
- [Servicios](./documentos/SERVICIOS.md) - Configuración de cada servicio
- [Guía Rápida](./documentos/GUIA-RAPIDA.md) - Tutorial paso a paso
- [FAQ](./documentos/FAQ.md) - Preguntas frecuentes y troubleshooting
- [Deprecado/Mejoras](./documentos/DEPRECADO-MEJORAS.md) - Deuda técnica

---

## Troubleshooting

### Docker no está corriendo
```bash
open -a Docker  # macOS
```

### Error de autenticación en ghcr.io
```bash
docker login ghcr.io
# Usuario: tu-github-username
# Password: ghp_tu_token
```

### Puerto en uso
```bash
lsof -ti:8081 | xargs kill -9
```

Más soluciones en [FAQ.md](./documentos/FAQ.md).

---

## Credenciales por Defecto

| Servicio | Usuario | Contraseña |
|----------|---------|------------|
| PostgreSQL | edugo | edugo123 |
| MongoDB | edugo | edugo123 |
| RabbitMQ | edugo | edugo123 |

---

**Última actualización:** Diciembre 2025  
**Mantenedor:** Equipo EduGo
