#!/usr/bin/env bash
set -euo pipefail

MODULES=("postgres" "mongodb" "schemas" "tools/mock-generator")
GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$GO_FILES" ]; then
  echo "No hay archivos Go staged; omitiendo validaciones de módulos"
  exit 0
fi

echo "Ejecutando validaciones pre-commit..."

declare -i ERRORS=0

for module in "${MODULES[@]}"; do
  if ! echo "$GO_FILES" | grep -q "^${module}/"; then
    continue
  fi

  echo ""
  echo "==> $module"

  if ! make -C "$module" fmt-check; then
    ERRORS+=1
  fi
  if ! make -C "$module" vet; then
    ERRORS+=1
  fi
  if ! make -C "$module" test; then
    ERRORS+=1
  fi
 done

if [ "$ERRORS" -ne 0 ]; then
  echo ""
  echo "Pre-commit falló con $ERRORS error(es)"
  exit 1
fi

echo ""
echo "Pre-commit completado correctamente"
