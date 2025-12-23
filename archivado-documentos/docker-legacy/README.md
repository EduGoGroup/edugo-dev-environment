# Archivos Docker Legacy

Estos archivos fueron archivados porque:
- Se consolidaron en `docker/docker-compose.yml` principal con profiles
- Son documentos temporales de validación que ya no son necesarios

## Docker Compose Files Archivados

| Archivo | Razón |
|---------|-------|
| `docker-compose-apps.yml` | Consolidado en docker-compose.yml con profiles |
| `docker-compose-infrastructure.yml` | Consolidado en docker-compose.yml (servicios sin profile) |
| `docker-compose-mock.yml` | No se usa actualmente |

## Documentos Archivados

| Archivo | Razón |
|---------|-------|
| `ACTUALIZAR_BASE_DATOS.md` | Documento de proceso específico |
| `PLAN_PRUEBAS_DOCKER_COMPOSE.md` | Plan de pruebas temporal |
| `QUICK_START.md` | Redundante con `documentos/GUIA-RAPIDA.md` |
| `RESULTADO_VALIDACION.md` | Reporte de validación puntual |

## Usar el nuevo docker-compose.yml

```bash
# Solo infraestructura (postgres, mongodb, rabbitmq)
docker-compose up -d

# Con API Mobile
docker-compose --profile apps up -d

# Con todas las apps
docker-compose --profile full up -d
```

---

**Archivado:** Diciembre 2025
