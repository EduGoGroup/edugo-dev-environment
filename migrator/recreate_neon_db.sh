#!/bin/bash

# ============================================================
# Script para recrear la base de datos en Neon
# Equivalente a: docker-compose down -v && docker-compose up -d
# ============================================================

set -e  # Salir si hay error

echo "🔄 Recreando base de datos EduGo en Neon..."
echo ""

# Verificar si se quiere aplicar mock data
APPLY_MOCK=${APPLY_MOCK_DATA:-true}

if [ "$APPLY_MOCK" = "false" ]; then
    echo "⚠️  MODO: Solo estructura y seeds (sin datos de prueba)"
else
    echo "⚠️  MODO: Estructura + seeds + datos de prueba"
fi

echo ""
read -p "⚠️  Esto ELIMINARÁ todos los datos. ¿Continuar? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Operación cancelada"
    exit 1
fi

echo ""
echo "🔥 Eliminando y recreando base de datos..."

# Ejecutar migración con FORCE_MIGRATION=true
FORCE_MIGRATION=true APPLY_MOCK_DATA=$APPLY_MOCK go run migrate_to_neon.go

echo ""
echo "✅ Base de datos recreada exitosamente"
echo ""
echo "📋 Para recrear SIN datos de prueba, usa:"
echo "   APPLY_MOCK_DATA=false ./recreate_neon_db.sh"
