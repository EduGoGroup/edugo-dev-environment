#!/bin/bash

set -e

echo "🧹 Limpiando ambiente de desarrollo EduGo..."
echo "=========================================="
echo ""

# Verificar que estamos en el directorio correcto
if [ ! -f "docker/docker-compose.yml" ]; then
    echo "❌ Error: Ejecuta este script desde la raíz de edugo-dev-environment/"
    echo "   Ejemplo: ./scripts/cleanup.sh"
    exit 1
fi

cd docker

# Detener y eliminar contenedores
echo "🛑 Deteniendo contenedores..."
docker-compose down

echo ""
echo "✅ Contenedores detenidos"
echo ""

# Preguntar si eliminar volúmenes
read -p "¿Deseas eliminar los volúmenes (DATOS DE BD)? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Eliminando volúmenes..."
    docker-compose down -v
    echo "✅ Volúmenes eliminados (se perdieron datos de PostgreSQL, MongoDB y RabbitMQ)"
else
    echo "ℹ️  Volúmenes preservados (los datos persisten)"
fi

echo ""

# Preguntar si limpiar imágenes
read -p "¿Deseas limpiar imágenes Docker no usadas? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Limpiando imágenes no usadas..."
    docker image prune -f
    echo "✅ Imágenes limpiadas"
else
    echo "ℹ️  Imágenes preservadas"
fi

echo ""

# Preguntar si eliminar imágenes de EduGo
read -p "¿Deseas eliminar imágenes de EduGo (api-mobile, api-admin, worker)? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Eliminando imágenes de EduGo..."
    docker rmi ghcr.io/medinatello/api-mobile:latest 2>/dev/null || echo "   api-mobile ya eliminada"
    docker rmi ghcr.io/medinatello/api-administracion:latest 2>/dev/null || echo "   api-administracion ya eliminada"
    docker rmi ghcr.io/medinatello/worker:latest 2>/dev/null || echo "   worker ya eliminada"
    echo "✅ Imágenes de EduGo eliminadas"
    echo ""
    echo "   Para volver a usarlas, ejecuta: ./scripts/update-images.sh"
else
    echo "ℹ️  Imágenes de EduGo preservadas"
fi

echo ""
echo "✅ Limpieza completada!"
echo ""
echo "Para volver a iniciar el ambiente:"
echo "  cd docker"
echo "  docker-compose up -d"
echo ""
