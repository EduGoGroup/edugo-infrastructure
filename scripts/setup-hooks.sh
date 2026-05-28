#!/bin/bash
# Setup script para instalar pre-commit hooks en edugo-infrastructure

set -e

echo "🔧 Instalando pre-commit hooks en edugo-infrastructure..."
echo ""

# Verificar que estamos en el repo correcto
if [ ! -f ".github/workflows/ci.yml" ]; then
  echo "❌ Error: Este script debe ejecutarse desde la raíz de edugo-infrastructure"
  exit 1
fi

# Copiar hook
HOOK_SOURCE="scripts/pre-commit-hook.sh"
HOOK_DEST=".git/hooks/pre-commit"

if [ ! -f "$HOOK_SOURCE" ]; then
  echo "❌ Error: $HOOK_SOURCE no encontrado"
  echo "   Verifica que el archivo existe en el repositorio"
  exit 1
fi

# Backup del hook existente si hay
if [ -f "$HOOK_DEST" ]; then
  echo "⚠️  Hook pre-commit ya existe, creando backup..."
  mv "$HOOK_DEST" "$HOOK_DEST.backup.$(date +%Y%m%d-%H%M%S)"
  echo "   Backup guardado en $HOOK_DEST.backup.*"
fi

# Instalar hook
cp "$HOOK_SOURCE" "$HOOK_DEST"
chmod +x "$HOOK_DEST"

echo "✅ Pre-commit hook instalado exitosamente"
echo ""
echo "El hook ejecutará automáticamente antes de cada commit:"
echo "  1. go fmt (formato)"
echo "  2. go vet (análisis estático)"
echo "  3. go mod tidy check"
echo "  4. go test -short (tests unitarios)"
echo ""
echo "Para bypass temporal: git commit --no-verify"
echo ""
