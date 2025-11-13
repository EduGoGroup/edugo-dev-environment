#!/bin/bash

set -e

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🌱 EduGo - Carga de Datos de Prueba${NC}"
echo "========================================"
echo ""

# Verificar que los containers estén corriendo
if ! docker ps | grep -q edugo-postgres; then
    echo -e "${YELLOW}⚠️  PostgreSQL no está corriendo${NC}"
    echo "   Inicia los servicios primero: ./scripts/setup.sh"
    exit 1
fi

echo -e "${BLUE}📊 Cargando datos en PostgreSQL...${NC}"
echo ""

# Verificar si existen seeds de PostgreSQL
if [ -d "seeds/postgresql" ]; then
    for sql_file in seeds/postgresql/*.sql; do
        if [ -f "$sql_file" ]; then
            filename=$(basename "$sql_file")
            echo -e "${BLUE}  ↳ Ejecutando: $filename${NC}"
            docker exec -i edugo-postgres psql -U edugo -d edugo < "$sql_file"
            echo -e "${GREEN}  ✅ $filename aplicado${NC}"
        fi
    done
else
    echo -e "${YELLOW}⚠️  No se encontró carpeta seeds/postgresql/${NC}"
    echo "   Los seeds se agregarán en la siguiente fase"
fi

# Cargar datos en MongoDB si está corriendo
if docker ps | grep -q edugo-mongodb; then
    echo ""
    echo -e "${BLUE}📊 Cargando datos en MongoDB...${NC}"
    echo ""
    
    if [ -d "seeds/mongodb" ]; then
        for js_file in seeds/mongodb/*.js; do
            if [ -f "$js_file" ]; then
                filename=$(basename "$js_file")
                echo -e "${BLUE}  ↳ Ejecutando: $filename${NC}"
                docker exec -i edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo < "$js_file"
                echo -e "${GREEN}  ✅ $filename aplicado${NC}"
            fi
        done
    else
        echo -e "${YELLOW}⚠️  No se encontró carpeta seeds/mongodb/${NC}"
        echo "   Los seeds se agregarán en la siguiente fase"
    fi
fi

echo ""
echo -e "${GREEN}🎉 Datos de prueba cargados exitosamente${NC}"
echo ""
