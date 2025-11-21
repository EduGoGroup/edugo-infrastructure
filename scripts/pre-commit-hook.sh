#!/bin/bash
# Pre-commit hook para edugo-infrastructure
# Ejecuta checks básicos antes de permitir commit

set -e
set -o pipefail

echo "🔍 Ejecutando pre-commit hooks..."
echo ""

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Módulos de infrastructure
MODULES=("postgres" "mongodb" "messaging" "schemas")

# Contador de errores
ERRORS=0

# 1. go fmt
echo "1️⃣  Verificando formato (go fmt)..."
for module in "${MODULES[@]}"; do
  if [ -d "$module" ]; then
    cd "$module"

    UNFORMATTED=$(gofmt -l . 2>&1 | grep -v "vendor/" || true)
    if [ -n "$UNFORMATTED" ]; then
      echo -e "${RED}❌ Archivos sin formatear en $module:${NC}"
      echo "$UNFORMATTED"
      ERRORS=$((ERRORS + 1))
    fi

    cd ..
  fi
done

if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}✅ Formato correcto${NC}"
fi
echo ""

# 2. go vet
echo "2️⃣  Ejecutando go vet..."
for module in "${MODULES[@]}"; do
  if [ -d "$module" ]; then
    cd "$module"

    if ! go vet ./... 2>&1; then
      echo -e "${RED}❌ go vet falló en $module${NC}"
      ERRORS=$((ERRORS + 1))
    fi

    cd ..
  fi
done

if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}✅ go vet pasó${NC}"
fi
echo ""

# 3. go mod tidy check
echo "3️⃣  Verificando go.mod actualizado..."
for module in "${MODULES[@]}"; do
  if [ -d "$module" ]; then
    cd "$module"

    # Guardar estado actual
    cp go.mod go.mod.bak
    cp go.sum go.sum.bak 2>/dev/null || true

    # Ejecutar go mod tidy
    go mod tidy 2>/dev/null || true

    # Comparar
    if ! diff -q go.mod go.mod.bak >/dev/null 2>&1; then
      echo -e "${YELLOW}⚠️  $module/go.mod necesita go mod tidy${NC}"
      # Restaurar
      mv go.mod.bak go.mod
      mv go.sum.bak go.sum 2>/dev/null || true
      ERRORS=$((ERRORS + 1))
    else
      rm go.mod.bak
      rm go.sum.bak 2>/dev/null || true
    fi

    cd ..
  fi
done

if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}✅ go.mod actualizados${NC}"
fi
echo ""

# 4. Tests unitarios (solo si hay archivos .go staged)
GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -n "$GO_FILES" ]; then
  echo "4️⃣  Ejecutando tests unitarios (short)..."

  for module in "${MODULES[@]}"; do
    # Verificar si hay archivos Go staged en este módulo
    MODULE_FILES=$(echo "$GO_FILES" | grep "^$module/" || true)

    if [ -n "$MODULE_FILES" ] && [ -d "$module" ]; then
      cd "$module"

      echo "   Testing $module..."
      if ! go test -short ./... 2>&1 | grep -E "(PASS|FAIL|ok|FAIL)"; then
        echo -e "${RED}❌ Tests fallaron en $module${NC}"
        ERRORS=$((ERRORS + 1))
      fi

      cd ..
    fi
  done

  if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Tests pasaron${NC}"
  fi
else
  echo "4️⃣  No hay archivos Go modificados, skipeando tests"
fi
echo ""

# Resumen
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}✅ Pre-commit hooks pasaron${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  exit 0
else
  echo -e "${RED}❌ Pre-commit hooks fallaron ($ERRORS errores)${NC}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""
  echo "🔧 Para arreglar:"
  echo "   1. Ejecuta: go fmt ./..."
  echo "   2. Ejecuta: go vet ./..."
  echo "   3. Ejecuta: go mod tidy"
  echo "   4. Ejecuta: go test -short ./..."
  echo ""
  echo "O bypass con: git commit --no-verify (NO RECOMENDADO)"
  exit 1
fi
