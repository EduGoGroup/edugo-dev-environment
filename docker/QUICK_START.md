# 🚀 EduGo - Inicio Rápido

## 🎨 Modo Mock - Lo Más Rápido (30 segundos)

**Para frontend developers que solo necesitan las APIs funcionando:**

```bash
cd edugo-dev-environment/docker

# Levantar APIs en modo mock (sin bases de datos)
docker-compose -f docker-compose-mock.yml up -d

# Verificar
curl http://localhost:8081/health
curl http://localhost:8082/health
```

**¡Listo!** APIs corriendo sin PostgreSQL, MongoDB ni RabbitMQ.

- API Mobile: http://localhost:8081
- API Admin: http://localhost:8082
- Login: `admin@edugo.test` / `edugo2024`

> ⚠️ Los datos son mock (en memoria) y se reinician con cada restart.

---

## Instalación Completa (5 minutos)

```bash
# 1. Clonar el repositorio (si no lo has hecho)
git clone <repo-url>
cd edugo-dev-environment/docker

# 2. Copiar variables de entorno (opcional, ya hay un .env funcional)
cp .env.example .env

# 3. Levantar todo
docker-compose --profile full up -d

# 4. Esperar ~30 segundos para que todo inicie

# 5. Validar que funciona
curl http://localhost:8081/health | jq
```

**¡Listo!** API Mobile corriendo en http://localhost:8081/swagger/index.html

---

## Solo Infraestructura (bases de datos)

```bash
cd edugo-dev-environment/docker
docker-compose -f docker-compose-infrastructure.yml up -d
```

**Conexiones**:
- PostgreSQL: `postgresql://edugo:edugo123@localhost:5432/edugo`
- MongoDB: `mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin`
- RabbitMQ: `amqp://edugo:edugo123@localhost:5672/`

---

## Solo API Mobile

```bash
# 1. Levantar infraestructura primero
docker-compose -f docker-compose-infrastructure.yml up -d

# 2. Crear red si no existe
docker network create edugo-network 2>/dev/null || true

# 3. Levantar API
docker-compose -f docker-compose-apps.yml up -d api-mobile

# 4. Validar
open http://localhost:8081/swagger/index.html
```

---

## Comandos Útiles

```bash
# Ver logs en tiempo real
docker-compose logs -f api-mobile

# Ver estado de servicios
docker-compose ps

# Reiniciar un servicio
docker-compose restart api-mobile

# Detener todo
docker-compose down

# Detener y eliminar volúmenes (CUIDADO: borra datos)
docker-compose down -v

# === MODO MOCK ===
# Levantar modo mock
docker-compose -f docker-compose-mock.yml up -d

# Ver logs modo mock
docker-compose -f docker-compose-mock.yml logs -f

# Detener modo mock
docker-compose -f docker-compose-mock.yml down
```

---

## ⚠️ Notas Importantes

1. **API Admin y Worker**: Actualmente requieren archivos `config.yaml` para funcionar. Ver [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md) para soluciones.

2. **S3**: Las credenciales de S3 en `.env` son placeholders. API Mobile funciona porque tiene `BOOTSTRAP_OPTIONAL_RESOURCES_S3=true`.

3. **OpenAI**: Worker necesita una API key real de OpenAI. Agregar en `.env`:
   ```
   OPENAI_API_KEY=sk-proj-tu-api-key-real-aqui
   ```

4. **Puertos en uso**: Si hay error de puerto ocupado:
   ```bash
   lsof -ti:8081 | xargs kill -9
   lsof -ti:8082 | xargs kill -9
   ```

---

## 📚 Más Información

- [README.md](./README.md) - Guía completa de uso
- [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md) - Reporte de validación detallado
- [../docs/dev-environment/](../docs/dev-environment/) - Documentación completa del proyecto

---

## 🆘 ¿Problemas?

1. Ver logs: `docker-compose logs -f [servicio]`
2. Consultar [README.md](./README.md) sección Troubleshooting
3. Revisar [RESULTADO_VALIDACION.md](./RESULTADO_VALIDACION.md) para problemas conocidos
