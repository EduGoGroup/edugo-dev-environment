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

# Función para mostrar ayuda
show_help() {
    echo "🚀 EduGo - Setup de Ambiente de Desarrollo"
    echo ""
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -p, --profile <profile>   Perfil de Docker Compose a usar (default: full)"
    echo "  -s, --seed                Cargar datos de prueba después de iniciar"
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
    echo -e "${GREEN}✅ Servicios iniciados correctamente${NC}"
else
    echo -e "${RED}❌ Error al iniciar servicios${NC}"
    exit 1
fi

# Esperar a que los servicios estén listos
echo ""
echo -e "${BLUE}⏳ Esperando a que los servicios estén listos...${NC}"
sleep 10

# Mostrar estado de los servicios
echo ""
echo -e "${BLUE}📊 Estado de los servicios:${NC}"
docker-compose --profile $PROFILE ps

# Cargar datos de prueba si se especificó
if [ "$SEED_DATA" = true ]; then
    echo ""
    echo -e "${BLUE}🌱 Cargando datos de prueba...${NC}"
    cd ..
    if [ -f scripts/seed-data.sh ]; then
        bash scripts/seed-data.sh
    else
        echo -e "${YELLOW}⚠️  Script de seed no encontrado: scripts/seed-data.sh${NC}"
    fi
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
echo ""
