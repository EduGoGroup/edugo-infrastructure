#!/bin/bash
# reproduce-failures.sh
# Intenta reproducir los fallos identificados localmente

set -e

echo "🔬 Reproduciendo fallos de CI localmente..."
echo "Versión de Go: $(go version)"
echo ""

# Módulos de infrastructure
MODULES=(
  "postgres"
  "mongodb"
  "messaging"
  "schemas"
)

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SUCCESS=0
FAILED=0

for module in "${MODULES[@]}"; do
  if [ ! -d "$module" ]; then
    echo -e "${YELLOW}⚠️  Módulo $module no encontrado${NC}"
    continue
  fi

  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🧪 Testeando módulo: $module"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  cd "$module"

  # Paso 1: Verificar go.mod
  echo "1️⃣  Verificando go.mod..."
  if go mod verify; then
    echo -e "${GREEN}✅ go.mod válido${NC}"
  else
    echo -e "${RED}❌ go.mod inválido${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi

  # Paso 2: Descargar dependencias
  echo ""
  echo "2️⃣  Descargando dependencias..."
  if go mod download; then
    echo -e "${GREEN}✅ Dependencias descargadas${NC}"
  else
    echo -e "${RED}❌ Error descargando dependencias${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi

  # Paso 3: Compilar
  echo ""
  echo "3️⃣  Compilando módulo..."
  if go build ./...; then
    echo -e "${GREEN}✅ Compilación exitosa${NC}"
  else
    echo -e "${RED}❌ Error de compilación${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi

  # Paso 4: Tests unitarios (sin integración)
  echo ""
  echo "4️⃣  Ejecutando tests unitarios (con -short)..."
  mkdir -p ../logs
  if go test -short -v ./... 2>&1 | tee "../logs/test-$module.log"; then
    echo -e "${GREEN}✅ Tests unitarios pasaron${NC}"
    SUCCESS=$((SUCCESS + 1))
  else
    echo -e "${RED}❌ Tests unitarios fallaron${NC}"
    echo "    Ver logs/test-$module.log para detalles"
    FAILED=$((FAILED + 1))
  fi

  cd ..
  echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESUMEN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Exitosos: $SUCCESS${NC}"
echo -e "${RED}❌ Fallidos: $FAILED${NC}"
echo "📦 Total: ${#MODULES[@]}"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}🎉 Todos los módulos pasaron localmente${NC}"
  echo ""
  echo "⚠️  NOTA: Los fallos de CI pueden ser por:"
  echo "   - Tests de integración (requieren servicios externos)"
  echo "   - Diferencias de ambiente (GitHub Actions vs local)"
  echo "   - Race conditions en CI"
  exit 0
else
  echo -e "${RED}⚠️  Algunos módulos fallaron${NC}"
  echo ""
  echo "📋 Próximos pasos:"
  echo "   1. Revisar logs en logs/test-*.log"
  echo "   2. Identificar diferencias con CI"
  echo "   3. Corregir en Tarea 2.1"
  exit 1
fi
