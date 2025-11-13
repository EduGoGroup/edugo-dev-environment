#!/bin/bash

set -e

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Variables por defecto
PROFILE="full"
REMOVE_VOLUMES=false

# Función para mostrar ayuda
show_help() {
    echo "🛑 EduGo - Detener Ambiente de Desarrollo"
    echo ""
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -p, --profile <profile>   Perfil de Docker Compose a detener (default: full)"
    echo "  -v, --volumes             Eliminar volúmenes de datos"
    echo "  -h, --help                Mostrar esta ayuda"
    echo ""
    echo "Perfiles disponibles:"
    echo "  full, db-only, api-only, mobile-only, admin-only, worker-only"
    echo ""
    echo "Ejemplos:"
    echo "  $0                      # Detiene todo (profile: full)"
    echo "  $0 --profile db-only    # Detiene solo bases de datos"
    echo "  $0 -v                   # Detiene todo y elimina volúmenes"
}

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--profile)
            PROFILE="$2"
            shift 2
            ;;
        -v|--volumes)
            REMOVE_VOLUMES=true
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

echo -e "${BLUE}🛑 EduGo - Detener Ambiente${NC}"
echo "============================"
echo ""
echo -e "${BLUE}📋 Configuración:${NC}"
echo "  - Perfil: $PROFILE"
echo "  - Eliminar volúmenes: $REMOVE_VOLUMES"
echo ""

cd docker

if [ "$REMOVE_VOLUMES" = true ]; then
    echo -e "${BLUE}🗑️  Deteniendo servicios y eliminando volúmenes...${NC}"
    docker-compose --profile $PROFILE down -v
else
    echo -e "${BLUE}⏸️  Deteniendo servicios (conservando datos)...${NC}"
    docker-compose --profile $PROFILE down
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Servicios detenidos correctamente${NC}"
    
    if [ "$REMOVE_VOLUMES" = true ]; then
        echo -e "${GREEN}✅ Volúmenes de datos eliminados${NC}"
    else
        echo -e "${BLUE}ℹ️  Los datos se conservaron en volúmenes de Docker${NC}"
        echo "   Para eliminarlos, usa: $0 --volumes"
    fi
else
    echo -e "${RED}❌ Error al detener servicios${NC}"
    exit 1
fi

echo ""
