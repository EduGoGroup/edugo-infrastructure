#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

echo "EduGo Infrastructure - Setup local"
echo ""

command -v docker >/dev/null 2>&1 || { echo "Docker no está instalado"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go no está instalado"; exit 1; }

docker info >/dev/null 2>&1 || { echo "Docker daemon no está corriendo"; exit 1; }

if [ ! -f ".env" ]; then
  cp .env.example .env
  echo "Se creó .env desde .env.example"
else
  echo ".env ya existe"
fi

echo ""
echo "Levantando PostgreSQL + MongoDB..."
make dev-up-core

echo ""
echo "Esperando que los servicios estén listos..."
sleep 10

echo ""
echo "Bootstrap de estructura y datos locales..."
make db-bootstrap

echo ""
echo "Setup completado"
echo ""
echo "Servicios disponibles:"
echo "  PostgreSQL: localhost:5432"
echo "  MongoDB:    localhost:27017"
echo ""
echo "Notas:"
echo "  - Este setup sirve para trabajar edugo-infrastructure de forma aislada."
echo "  - Para el ecosistema completo, el flujo canónico sigue siendo edugo-dev-environment."
