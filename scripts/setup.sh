#!/bin/bash

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables por defecto
PROFILE="full"
SEED_DATA=false
HEALTH_TIMEOUT=120  # Timeout global en segundos
HEALTH_INTERVAL=3   # Intervalo entre verificaciones

# Función para mostrar ayuda
show_help() {
    echo "🚀 EduGo - Setup de Ambiente de Desarrollo"
    echo ""
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -p, --profile <profile>   Perfil de Docker Compose a usar (default: full)"
    echo "  -s, --seed                Cargar datos de prueba después de iniciar"
    echo "  -t, --timeout <segundos>  Timeout para health checks (default: 120)"
    echo "  -h, --help                Mostrar esta ayuda"
    echo ""
    echo "Perfiles disponibles:"
    echo "  full         - Todos los servicios (PostgreSQL + MongoDB + RabbitMQ + APIs + Worker)"
    echo "  db-only      - Solo bases de datos (PostgreSQL + MongoDB + RabbitMQ)"
    echo "  api-only     - Bases de datos + APIs (sin Worker)"
    echo "  mobile-only  - Bases de datos + API Mobile"
    echo "  admin-only   - Bases de datos + API Administración"
    echo "  worker-only  - Bases de datos + Worker"
    echo ""
    echo "Ejemplos:"
    echo "  $0                          # Inicia todo (profile: full)"
    echo "  $0 --profile db-only        # Solo bases de datos"
    echo "  $0 --profile api-only -s    # APIs + DBs con datos de prueba"
    echo "  $0 -s                       # Todo con datos de prueba"
}

# Función para esperar que un servicio esté saludable
wait_for_healthy() {
    local container_name=$1
    local service_name=$2
    local max_wait=$HEALTH_TIMEOUT
    local elapsed=0

    echo -e "${BLUE}  ⏳ Esperando que $service_name esté saludable...${NC}"

    while [ $elapsed -lt $max_wait ]; do
        # Verificar si el contenedor existe
        if ! docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
            echo -e "${YELLOW}    Contenedor $container_name no encontrado, esperando...${NC}"
            sleep $HEALTH_INTERVAL
            elapsed=$((elapsed + HEALTH_INTERVAL))
            continue
        fi

        # Obtener estado de salud
        local health=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}no-healthcheck{{end}}' "$container_name" 2>/dev/null)

        case $health in
            healthy)
                echo -e "${GREEN}  ✅ $service_name está saludable${NC}"
                return 0
                ;;
            starting)
                echo -e "    $service_name iniciando... (${elapsed}s/${max_wait}s)"
                ;;
            unhealthy)
                echo -e "${YELLOW}    $service_name no saludable, reintentando...${NC}"
                ;;
            no-healthcheck)
                # Si no tiene healthcheck, verificar que está corriendo
                local status=$(docker inspect --format='{{.State.Status}}' "$container_name" 2>/dev/null)
                if [ "$status" = "running" ]; then
                    echo -e "${GREEN}  ✅ $service_name está corriendo (sin healthcheck)${NC}"
                    return 0
                fi
                ;;
        esac

        sleep $HEALTH_INTERVAL
        elapsed=$((elapsed + HEALTH_INTERVAL))
    done

    echo -e "${RED}  ❌ Timeout esperando $service_name después de ${max_wait}s${NC}"
    return 1
}

# Función para verificar conectividad de PostgreSQL
wait_for_postgres() {
    local container_name="edugo-postgres"
    local max_wait=$HEALTH_TIMEOUT
    local elapsed=0

    echo -e "${BLUE}  ⏳ Verificando conectividad de PostgreSQL...${NC}"

    while [ $elapsed -lt $max_wait ]; do
        if docker exec "$container_name" pg_isready -U edugo -d edugo &>/dev/null; then
            echo -e "${GREEN}  ✅ PostgreSQL acepta conexiones${NC}"
            return 0
        fi

        echo -e "    PostgreSQL no listo aún... (${elapsed}s/${max_wait}s)"
        sleep $HEALTH_INTERVAL
        elapsed=$((elapsed + HEALTH_INTERVAL))
    done

    echo -e "${RED}  ❌ Timeout esperando PostgreSQL${NC}"
    return 1
}

# Función para verificar conectividad de MongoDB
wait_for_mongodb() {
    local container_name="edugo-mongodb"
    local max_wait=$HEALTH_TIMEOUT
    local elapsed=0

    echo -e "${BLUE}  ⏳ Verificando conectividad de MongoDB...${NC}"

    while [ $elapsed -lt $max_wait ]; do
        if docker exec "$container_name" mongosh --eval "db.adminCommand('ping')" -u edugo -p edugo123 --authSource admin &>/dev/null; then
            echo -e "${GREEN}  ✅ MongoDB acepta conexiones${NC}"
            return 0
        fi

        echo -e "    MongoDB no listo aún... (${elapsed}s/${max_wait}s)"
        sleep $HEALTH_INTERVAL
        elapsed=$((elapsed + HEALTH_INTERVAL))
    done

    echo -e "${RED}  ❌ Timeout esperando MongoDB${NC}"
    return 1
}

# Función para verificar conectividad de RabbitMQ
wait_for_rabbitmq() {
    local container_name="edugo-rabbitmq"
    local max_wait=$HEALTH_TIMEOUT
    local elapsed=0

    echo -e "${BLUE}  ⏳ Verificando conectividad de RabbitMQ...${NC}"

    while [ $elapsed -lt $max_wait ]; do
        if docker exec "$container_name" rabbitmq-diagnostics ping &>/dev/null; then
            echo -e "${GREEN}  ✅ RabbitMQ acepta conexiones${NC}"
            return 0
        fi

        echo -e "    RabbitMQ no listo aún... (${elapsed}s/${max_wait}s)"
        sleep $HEALTH_INTERVAL
        elapsed=$((elapsed + HEALTH_INTERVAL))
    done

    echo -e "${RED}  ❌ Timeout esperando RabbitMQ${NC}"
    return 1
}

# Función para esperar que todos los servicios de infraestructura estén listos
wait_for_infrastructure() {
    echo ""
    echo -e "${BLUE}🔍 Verificando servicios de infraestructura...${NC}"

    local failed=0

    # Verificar PostgreSQL
    if ! wait_for_postgres; then
        failed=1
    fi

    # Verificar MongoDB
    if ! wait_for_mongodb; then
        failed=1
    fi

    # Verificar RabbitMQ
    if ! wait_for_rabbitmq; then
        failed=1
    fi

    if [ $failed -eq 1 ]; then
        echo ""
        echo -e "${RED}❌ Algunos servicios no están listos. Revisa los logs:${NC}"
        echo "   docker-compose logs postgres mongodb rabbitmq"
        return 1
    fi

    echo ""
    echo -e "${GREEN}✅ Todos los servicios de infraestructura están listos${NC}"
    return 0
}

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--profile)
            PROFILE="$2"
            shift 2
            ;;
        -s|--seed)
            SEED_DATA=true
            shift
            ;;
        -t|--timeout)
            HEALTH_TIMEOUT="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Opción desconocida: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🚀 EduGo - Setup de Ambiente de Desarrollo${NC}"
echo "=========================================="
echo ""
echo -e "${BLUE}📋 Configuración:${NC}"
echo "  - Perfil: $PROFILE"
echo "  - Seed data: $SEED_DATA"
echo "  - Timeout health checks: ${HEALTH_TIMEOUT}s"
echo ""

# Verificar que Docker está instalado
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker no está instalado.${NC}"
    echo "   Por favor instala Docker Desktop: https://docs.docker.com/desktop/install/"
    exit 1
fi

echo -e "${GREEN}✅ Docker está instalado${NC}"

# Verificar que Docker está corriendo
if ! docker info &> /dev/null; then
    echo -e "${RED}❌ Docker no está corriendo.${NC}"
    echo "   Por favor inicia Docker Desktop."
    exit 1
fi

echo -e "${GREEN}✅ Docker está corriendo${NC}"

# Crear archivo .env si no existe
if [ ! -f docker/.env ]; then
    echo -e "${BLUE}📝 Creando archivo .env desde .env.example...${NC}"
    cp docker/.env.example docker/.env
    echo -e "${GREEN}✅ Archivo .env creado${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANTE: Edita docker/.env si necesitas cambiar configuraciones${NC}"
    echo "   Especialmente OPENAI_API_KEY para que el worker funcione"
else
    echo -e "${GREEN}✅ Archivo .env ya existe${NC}"
fi

# Login a GitHub Container Registry
echo ""
echo -e "${BLUE}🔐 Configurando acceso a GitHub Container Registry...${NC}"
echo "Por favor ingresa tu GitHub Personal Access Token (con scope read:packages):"
echo "(El token debe tener formato: ghp_...)"
read -s GITHUB_TOKEN

if [ -z "$GITHUB_TOKEN" ]; then
    echo -e "${RED}❌ Token no puede estar vacío${NC}"
    exit 1
fi

echo ""
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin 2>&1 | grep -q "Login Succeeded"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Login exitoso a GitHub Container Registry${NC}"
else
    echo -e "${RED}❌ Error en login. Verifica tu token.${NC}"
    exit 1
fi

# Iniciar servicios con el perfil especificado
echo ""
echo -e "${BLUE}🐳 Iniciando servicios con perfil: ${PROFILE}${NC}"
cd docker
docker-compose --profile $PROFILE up -d

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Contenedores iniciados${NC}"
else
    echo -e "${RED}❌ Error al iniciar servicios${NC}"
    exit 1
fi

# Esperar a que los servicios de infraestructura estén listos (reemplaza sleep 10)
if ! wait_for_infrastructure; then
    echo -e "${RED}❌ Error: La infraestructura no está lista${NC}"
    exit 1
fi

# Mostrar estado de los servicios
echo ""
echo -e "${BLUE}📊 Estado de los servicios:${NC}"
docker-compose --profile $PROFILE ps

# Cargar datos de prueba si se especificó
if [ "$SEED_DATA" = true ]; then
    echo ""
    echo -e "${BLUE}🌱 Cargando datos via migrator (edugo-infrastructure)...${NC}"
    cd ..
    if [ -d migrator ]; then
        if (cd migrator && FORCE_MIGRATION=true go run cmd/main.go); then
            echo -e "${GREEN}✅ Base de datos recreada con todos los seeds${NC}"
        else
            echo -e "${YELLOW}⚠️  Hubo problemas ejecutando el migrator${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Directorio migrator/ no encontrado${NC}"
    fi
    cd docker
fi

# Mostrar URLs de los servicios
echo ""
echo -e "${GREEN}🎉 ¡Ambiente listo!${NC}"
echo ""
echo -e "${BLUE}📍 URLs de los servicios:${NC}"

case $PROFILE in
    full|api-only|mobile-only)
        echo "  - API Mobile: http://localhost:8081"
        echo "    - Swagger: http://localhost:8081/swagger/index.html"
        ;;
esac

case $PROFILE in
    full|api-only|admin-only)
        echo "  - API Admin: http://localhost:8082"
        echo "    - Swagger: http://localhost:8082/swagger/index.html"
        ;;
esac

case $PROFILE in
    full|db-only|api-only|mobile-only|admin-only|worker-only)
        echo "  - PostgreSQL: localhost:5432"
        echo "  - MongoDB: localhost:27017"
        ;;
esac

case $PROFILE in
    full|db-only|api-only|mobile-only|worker-only)
        echo "  - RabbitMQ Management: http://localhost:15672"
        echo "    - User: edugo / Pass: edugo123"
        ;;
esac

echo ""
echo -e "${BLUE}📝 Comandos útiles:${NC}"
echo "  - Ver logs: docker-compose --profile $PROFILE logs -f"
echo "  - Detener: docker-compose --profile $PROFILE down"
echo "  - Reiniciar: docker-compose --profile $PROFILE restart"
echo "  - Recrear DB: cd migrator && FORCE_MIGRATION=true go run cmd/main.go"
echo ""
