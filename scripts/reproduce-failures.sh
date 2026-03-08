#!/usr/bin/env bash
set -euo pipefail

MODULES=("postgres" "mongodb" "schemas" "tools/mock-generator")
SUCCESS=0
FAILED=0

echo "Reproduciendo validaciones de módulos..."

for module in "${MODULES[@]}"; do
  echo ""
  echo "========================================"
  echo "Módulo: $module"
  echo "========================================"

  if make -C "$module" release-check; then
    SUCCESS=$((SUCCESS + 1))
  else
    FAILED=$((FAILED + 1))
  fi
 done

echo ""
echo "Resumen"
echo "  OK:   $SUCCESS"
echo "  FAIL: $FAILED"

if [ "$FAILED" -ne 0 ]; then
  exit 1
fi
