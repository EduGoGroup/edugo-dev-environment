#!/bin/bash

# edugo-dev-environment - Script de Validación
# Valida sintaxis de docker-compose.yml

set -e  # Exit on error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Validando configuración de Docker Compose...${NC}"
echo ""

# Verificar que docker-compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ ERROR: docker-compose no está instalado${NC}"
    echo "Instalar desde: https://docs.docker.com/compose/install/"
    exit 1
fi

echo -e "${GREEN}✅ docker-compose instalado${NC}"
echo -e "   Versión: $(docker-compose --version)"

# Cambiar al directorio docker donde están los compose files
DOCKER_DIR="docker"
if [ ! -d "$DOCKER_DIR" ]; then
    echo -e "${RED}❌ ERROR: directorio $DOCKER_DIR no encontrado${NC}"
    echo "Ejecutar desde el directorio raíz del proyecto"
    exit 1
fi

cd "$DOCKER_DIR"

# Lista de archivos docker-compose a validar
COMPOSE_FILES=(
    "docker-compose.yml"
    "docker-compose.full.yml"
    "docker-compose.local.yml"
)

VALIDATED_COUNT=0
FAILED_COUNT=0

for compose_file in "${COMPOSE_FILES[@]}"; do
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📄 Validando: $compose_file${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    if [ ! -f "$compose_file" ]; then
        echo -e "${YELLOW}⚠️  Archivo no encontrado: $compose_file (saltando)${NC}"
        continue
    fi

    echo -e "${GREEN}✅ Archivo encontrado${NC}"

    # Validar sintaxis YAML
    echo ""
    echo "📝 Validando sintaxis YAML..."
    if docker-compose -f "$compose_file" config > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Sintaxis YAML válida${NC}"
    else
        echo -e "${RED}❌ ERROR: Sintaxis YAML inválida${NC}"
        echo ""
        docker-compose -f "$compose_file" config
        FAILED_COUNT=$((FAILED_COUNT + 1))
        continue
    fi

    # Verificar servicios definidos
    echo ""
    echo "🔍 Servicios encontrados:"
    docker-compose -f "$compose_file" config --services | while read service; do
        echo -e "  ${GREEN}✓${NC} $service"
    done

    # Verificar volúmenes definidos
    echo ""
    echo "💾 Volúmenes encontrados:"
    VOLUMES=$(docker-compose -f "$compose_file" config --volumes 2>/dev/null)
    if [ -z "$VOLUMES" ]; then
        echo -e "  ${YELLOW}(ninguno)${NC}"
    else
        echo "$VOLUMES" | while read volume; do
            echo -e "  ${GREEN}✓${NC} $volume"
        done
    fi

    # Verificar puertos expuestos
    echo ""
    echo "🌐 Puertos expuestos:"
    PORTS=$(docker-compose -f "$compose_file" config | grep -A 1 "ports:" | grep -o "[0-9]*:[0-9]*" | sort -u)
    if [ -z "$PORTS" ]; then
        echo -e "  ${YELLOW}(ninguno)${NC}"
    else
        echo "$PORTS" | while read port; do
            echo -e "  ${GREEN}✓${NC} $port"
        done
    fi

    # Verificar variables de entorno requeridas
    echo ""
    echo "🔐 Variables de entorno:"
    # Usar grep compatible con macOS (BSD grep) en lugar de GNU grep -P
    ENV_VARS=$(docker-compose -f "$compose_file" config | grep -o '\${[A-Z_]*[^}]*}' | sed 's/\${//g' | sed 's/}.*//g' | sort -u)
    if [ -z "$ENV_VARS" ]; then
        echo -e "  ${GREEN}✓${NC} No requiere variables de entorno"
    else
        echo "$ENV_VARS" | while read var; do
            if [ -z "${!var}" ]; then
                echo -e "  ${YELLOW}⚠${NC}  $var (no definida)"
            else
                echo -e "  ${GREEN}✓${NC} $var (definida)"
            fi
        done
    fi

    VALIDATED_COUNT=$((VALIDATED_COUNT + 1))
done

# Volver al directorio raíz
cd ..

# Verificar que .env existe (si es requerido)
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🔐 Verificando archivos de configuración${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ -f "$DOCKER_DIR/.env.example" ] && [ ! -f "$DOCKER_DIR/.env" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  ADVERTENCIA: .env no existe${NC}"
    echo "Crear desde .env.example:"
    echo -e "  ${BLUE}cp $DOCKER_DIR/.env.example $DOCKER_DIR/.env${NC}"
else
    echo -e "${GREEN}✅ Archivo .env existe${NC}"
fi

# Resumen final
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 Resumen de Validación${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "Archivos validados: ${GREEN}$VALIDATED_COUNT${NC}"
echo -e "Archivos con errores: ${RED}$FAILED_COUNT${NC}"

if [ $FAILED_COUNT -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Validación completada exitosamente${NC}"
    echo ""
    echo "Próximo paso:"
    echo -e "  ${BLUE}cd docker && docker-compose up -d${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}❌ Validación falló. Corregir errores antes de continuar.${NC}"
    exit 1
fi
