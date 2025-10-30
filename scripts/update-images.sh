#!/bin/bash

set -e

echo "🔄 Actualizando imágenes Docker de EduGo..."
echo "=========================================="
echo ""

# Verificar que Docker está corriendo
if ! docker info &> /dev/null; then
    echo "❌ Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi

echo "✅ Docker está corriendo"
echo ""

# Verificar login a ghcr.io
echo "🔐 Verificando acceso a GitHub Container Registry..."
if ! docker info 2>/dev/null | grep -q "ghcr.io"; then
    echo "⚠️  No has hecho login a ghcr.io"
    echo "   Ejecuta primero: ./setup.sh"
    echo ""
    echo "O haz login manualmente:"
    echo "   echo \$GITHUB_TOKEN | docker login ghcr.io -u medinatello --password-stdin"
    exit 1
fi

echo "✅ Acceso a ghcr.io configurado"
echo ""

# Pull de las últimas imágenes
echo "📦 Descargando imágenes más recientes..."
echo ""

echo "⬇️  Actualizando api-mobile..."
docker pull ghcr.io/medinatello/api-mobile:latest

echo ""
echo "⬇️  Actualizando api-administracion..."
docker pull ghcr.io/medinatello/api-administracion:latest

echo ""
echo "⬇️  Actualizando worker..."
docker pull ghcr.io/medinatello/worker:latest

echo ""
echo "✅ Imágenes actualizadas!"
echo ""
echo "📋 Para aplicar los cambios:"
echo ""
echo "   cd docker"
echo "   docker-compose down"
echo "   docker-compose up -d"
echo ""
echo "💡 Tip: Puedes ver qué versión tienes con:"
echo "   docker images | grep medinatello"
echo ""
