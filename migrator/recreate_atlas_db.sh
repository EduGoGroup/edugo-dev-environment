#!/bin/bash

# ============================================================
# Script para recrear la base de datos MongoDB en Atlas
# Equivalente a: docker-compose down -v && docker-compose up -d (para MongoDB)
# ============================================================

set -e  # Salir si hay error

echo "🔄 Recreando base de datos MongoDB EduGo en Atlas..."
echo ""

# Verificar si se quiere aplicar mock data
APPLY_MOCK=${APPLY_MOCK_DATA:-false}

if [ "$APPLY_MOCK" = "true" ]; then
    echo "⚠️  ADVERTENCIA: Los datos de prueba de MongoDB tienen problemas de validación"
    echo "⚠️  Se recomienda usar APPLY_MOCK_DATA=false"
fi

if [ "$APPLY_MOCK" = "false" ]; then
    echo "⚠️  MODO: Solo estructura (sin datos de prueba)"
else
    echo "⚠️  MODO: Estructura + datos de prueba"
fi

echo ""
read -p "⚠️  Esto ELIMINARÁ todos los datos de MongoDB. ¿Continuar? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Operación cancelada"
    exit 1
fi

echo ""
echo "🔥 Eliminando y recreando base de datos MongoDB..."

# Ejecutar migración con FORCE_MIGRATION=true
FORCE_MIGRATION=true APPLY_MOCK_DATA=$APPLY_MOCK go run migrate_to_atlas.go

echo ""
echo "✅ Base de datos MongoDB recreada exitosamente"
echo ""
echo "📋 Para recrear sin datos de prueba (recomendado), usa:"
echo "   APPLY_MOCK_DATA=false ./recreate_atlas_db.sh"
