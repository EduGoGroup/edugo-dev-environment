#!/bin/bash

set -e

echo "🚀 EduGo - Setup de Ambiente de Desarrollo"
echo "=========================================="
echo ""

# Verificar que Docker está instalado
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado. Por favor instala Docker Desktop."
    echo "   https://docs.docker.com/desktop/install/mac-install/"
    exit 1
fi

echo "✅ Docker está instalado"

# Verificar que Docker está corriendo
if ! docker info &> /dev/null; then
    echo "❌ Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi

echo "✅ Docker está corriendo"

# Crear archivo .env si no existe
if [ ! -f docker/.env ]; then
    echo "📝 Creando archivo .env desde .env.example..."
    cp docker/.env.example docker/.env
    echo "✅ Archivo .env creado"
    echo ""
    echo "⚠️  IMPORTANTE: Edita docker/.env si necesitas cambiar configuraciones"
    echo "   Especialmente OPENAI_API_KEY para que el worker funcione"
else
    echo "✅ Archivo .env ya existe"
fi

# Login a GitHub Container Registry
echo ""
echo "🔐 Configurando acceso a GitHub Container Registry..."
echo "Por favor ingresa tu GitHub Personal Access Token (con scope read:packages):"
echo "(El token debe tener formato: ghp_...)"
read -s GITHUB_TOKEN

if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ Token no puede estar vacío"
    exit 1
fi

echo ""
echo "Intentando login a ghcr.io..."
echo "$GITHUB_TOKEN" | docker login ghcr.io -u medinatello --password-stdin

if [ $? -eq 0 ]; then
    echo "✅ Login exitoso a ghcr.io"
else
    echo "❌ Error al hacer login. Verifica tu token."
    echo "   El token debe tener permisos: read:packages"
    exit 1
fi

# Pull de las imágenes más recientes
echo ""
echo "📦 Descargando imágenes Docker más recientes..."
echo "   Esto puede tomar varios minutos la primera vez..."
echo ""

docker pull ghcr.io/medinatello/api-mobile:latest || echo "⚠️  Imagen api-mobile no disponible aún (normal si no has hecho deploy)"
docker pull ghcr.io/medinatello/api-administracion:latest || echo "⚠️  Imagen api-administracion no disponible aún"
docker pull ghcr.io/medinatello/worker:latest || echo "⚠️  Imagen worker no disponible aún"

echo ""
echo "✅ Setup completado!"
echo ""
echo "📋 Próximos pasos:"
echo ""
echo "1. Editar configuración (opcional):"
echo "   nano docker/.env"
echo ""
echo "2. Iniciar el ambiente:"
echo "   cd docker"
echo "   docker-compose up -d"
echo ""
echo "3. Ver los logs:"
echo "   docker-compose logs -f"
echo ""
echo "4. Verificar servicios:"
echo "   docker-compose ps"
echo ""
echo "5. Probar endpoints:"
echo "   curl http://localhost:8081/health  # API Mobile"
echo "   curl http://localhost:8082/health  # API Admin"
echo "   open http://localhost:15672        # RabbitMQ UI (user: edugo, pass: edugo123)"
echo ""
echo "Para detener:"
echo "   docker-compose down"
echo ""
