#!/bin/bash

# EduGo Dev Environment - Script de Diagnóstico
# Verifica el estado del ambiente y detecta problemas comunes

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo ""
echo -e "${BLUE}🔍 EduGo - Diagnóstico del Ambiente${NC}"
echo "======================================"
echo ""

ISSUES_FOUND=0

# ==========================================
# 1. VERIFICAR DOCKER
# ==========================================
echo -e "${BLUE}1. Docker${NC}"
echo "----------"

if ! command -v docker &> /dev/null; then
    echo -e "  ${RED}❌ Docker no instalado${NC}"
    echo "     Instalar: https://docs.docker.com/desktop/"
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
else
    echo -e "  ${GREEN}✅ Docker instalado${NC}"
    
    if ! docker info &> /dev/null 2>&1; then
        echo -e "  ${RED}❌ Docker no está corriendo${NC}"
        echo "     Solución: open -a Docker"
        ISSUES_FOUND=$((ISSUES_FOUND + 1))
    else
        echo -e "  ${GREEN}✅ Docker está corriendo${NC}"
        
        # Espacio en disco de Docker
        DOCKER_DISK=$(docker system df --format "{{.Size}}" 2>/dev/null | head -1)
        echo -e "  ${BLUE}ℹ️  Uso de disco Docker: $DOCKER_DISK${NC}"
    fi
fi
echo ""

# ==========================================
# 2. VERIFICAR CONTENEDORES
# ==========================================
echo -e "${BLUE}2. Contenedores${NC}"
echo "----------------"

check_container() {
    local name=$1
    local port=$2
    
    if docker ps --format '{{.Names}}' | grep -q "^${name}$"; then
        STATUS=$(docker inspect --format='{{.State.Health.Status}}' $name 2>/dev/null || echo "unknown")
        if [ "$STATUS" = "healthy" ]; then
            echo -e "  ${GREEN}✅ $name (healthy)${NC}"
        elif [ "$STATUS" = "unknown" ]; then
            echo -e "  ${GREEN}✅ $name (running)${NC}"
        else
            echo -e "  ${YELLOW}⚠️  $name ($STATUS)${NC}"
            ISSUES_FOUND=$((ISSUES_FOUND + 1))
        fi
    else
        if docker ps -a --format '{{.Names}}' | grep -q "^${name}$"; then
            echo -e "  ${RED}❌ $name (detenido)${NC}"
            echo "     Solución: cd docker && docker-compose up -d"
        else
            echo -e "  ${YELLOW}⚠️  $name (no existe)${NC}"
        fi
        ISSUES_FOUND=$((ISSUES_FOUND + 1))
    fi
}

check_container "edugo-postgres"
check_container "edugo-mongodb"
check_container "edugo-rabbitmq"
check_container "edugo-api-mobile"
check_container "edugo-api-administracion"
check_container "edugo-worker"
echo ""

# ==========================================
# 3. VERIFICAR PUERTOS
# ==========================================
echo -e "${BLUE}3. Puertos${NC}"
echo "-----------"

check_port() {
    local port=$1
    local service=$2
    
    if lsof -i :$port -sTCP:LISTEN &> /dev/null; then
        PROCESS=$(lsof -i :$port -sTCP:LISTEN | tail -1 | awk '{print $1}')
        if [[ "$PROCESS" == *"docker"* ]] || [[ "$PROCESS" == "com.docke"* ]]; then
            echo -e "  ${GREEN}✅ Puerto $port ($service) - Docker${NC}"
        else
            echo -e "  ${YELLOW}⚠️  Puerto $port ocupado por: $PROCESS${NC}"
            echo "     Solución: lsof -ti:$port | xargs kill -9"
            ISSUES_FOUND=$((ISSUES_FOUND + 1))
        fi
    else
        echo -e "  ${YELLOW}⚠️  Puerto $port ($service) - no en uso${NC}"
    fi
}

check_port 5432 "PostgreSQL"
check_port 27017 "MongoDB"
check_port 5672 "RabbitMQ"
check_port 15672 "RabbitMQ UI"
check_port 8081 "API Mobile"
check_port 8082 "API Admin"
echo ""

# ==========================================
# 4. VERIFICAR CONECTIVIDAD APIs
# ==========================================
echo -e "${BLUE}4. Health Checks${NC}"
echo "-----------------"

check_api() {
    local url=$1
    local name=$2
    
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 3 $url 2>/dev/null || echo "000")
    
    if [ "$RESPONSE" = "200" ]; then
        echo -e "  ${GREEN}✅ $name - OK${NC}"
    elif [ "$RESPONSE" = "000" ]; then
        echo -e "  ${RED}❌ $name - No responde${NC}"
        ISSUES_FOUND=$((ISSUES_FOUND + 1))
    else
        echo -e "  ${YELLOW}⚠️  $name - HTTP $RESPONSE${NC}"
        ISSUES_FOUND=$((ISSUES_FOUND + 1))
    fi
}

check_api "http://localhost:8081/health" "API Mobile"
check_api "http://localhost:8082/health" "API Admin"
echo ""

# ==========================================
# 5. VERIFICAR CONFIGURACIÓN
# ==========================================
echo -e "${BLUE}5. Configuración${NC}"
echo "-----------------"

if [ -f "docker/.env" ]; then
    echo -e "  ${GREEN}✅ docker/.env existe${NC}"
    
    # Verificar OPENAI_API_KEY
    if grep -q "OPENAI_API_KEY=sk-" docker/.env; then
        echo -e "  ${GREEN}✅ OPENAI_API_KEY configurada${NC}"
    else
        echo -e "  ${YELLOW}⚠️  OPENAI_API_KEY no configurada (Worker no procesará PDFs)${NC}"
    fi
else
    echo -e "  ${RED}❌ docker/.env no existe${NC}"
    echo "     Solución: cp docker/.env.example docker/.env"
    ISSUES_FOUND=$((ISSUES_FOUND + 1))
fi
echo ""

# ==========================================
# 6. VERIFICAR AUTENTICACIÓN GHCR
# ==========================================
echo -e "${BLUE}6. GitHub Container Registry${NC}"
echo "-----------------------------"

if docker login ghcr.io --get-login &> /dev/null 2>&1; then
    echo -e "  ${GREEN}✅ Autenticado en ghcr.io${NC}"
else
    echo -e "  ${YELLOW}⚠️  No autenticado en ghcr.io${NC}"
    echo "     Solución: docker login ghcr.io"
fi
echo ""

# ==========================================
# 7. LOGS RECIENTES DE ERRORES
# ==========================================
echo -e "${BLUE}7. Errores Recientes${NC}"
echo "---------------------"

if docker ps --format '{{.Names}}' | grep -q "edugo-api-mobile"; then
    ERRORS=$(docker logs edugo-api-mobile 2>&1 | grep -i "error\|panic\|fatal" | tail -3)
    if [ -n "$ERRORS" ]; then
        echo -e "  ${YELLOW}⚠️  Errores en API Mobile:${NC}"
        echo "$ERRORS" | while read line; do
            echo "     $line"
        done
        ISSUES_FOUND=$((ISSUES_FOUND + 1))
    else
        echo -e "  ${GREEN}✅ Sin errores recientes en API Mobile${NC}"
    fi
else
    echo -e "  ${YELLOW}⚠️  API Mobile no está corriendo${NC}"
fi
echo ""

# ==========================================
# RESUMEN
# ==========================================
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}📊 Resumen${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""

if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ Todo funciona correctamente${NC}"
    echo ""
    echo "URLs disponibles:"
    echo "  - API Mobile:  http://localhost:8081"
    echo "  - API Admin:   http://localhost:8082"
    echo "  - Swagger:     http://localhost:8081/swagger/index.html"
    echo "  - RabbitMQ:    http://localhost:15672"
else
    echo -e "${YELLOW}⚠️  Se encontraron $ISSUES_FOUND problema(s)${NC}"
    echo ""
    echo "Comandos útiles:"
    echo "  - Ver logs:    make logs"
    echo "  - Reiniciar:   make restart"
    echo "  - Reset:       make reset"
fi
echo ""
