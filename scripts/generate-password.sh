#!/bin/bash

# Script para generar hash bcrypt de contraseñas
# Uso: ./scripts/generate-password.sh [password]

set -e

if [ "$#" -eq 0 ]; then
    echo "📝 Generador de Hash Bcrypt para Contraseñas"
    echo ""
    echo "Uso:"
    echo "  ./scripts/generate-password.sh <password>"
    echo ""
    echo "Ejemplo:"
    echo "  ./scripts/generate-password.sh mipassword123"
    echo ""
    exit 1
fi

PASSWORD="$1"

echo "🔐 Generando hash bcrypt..."
echo ""

cd "$(dirname "$0")/../migrator"

cat > /tmp/gen-hash-temp.go << 'GOEOF'
package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(hash))
}
GOEOF

HASH=$(go run /tmp/gen-hash-temp.go "$PASSWORD")
rm /tmp/gen-hash-temp.go

echo "Password: $PASSWORD"
echo "Hash:     $HASH"
echo ""
echo "✅ Hash generado exitosamente"
echo ""
echo "💡 Puedes usar este hash en:"
echo "   - Migraciones SQL (columna password_hash)"
echo "   - Tests de integración"
echo "   - Datos de prueba"
