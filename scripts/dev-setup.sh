#!/bin/bash
set -e

echo "🚀 EduGo Infrastructure - Setup de Desarrollo"
echo ""

# 1. Validar Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker no está instalado"
    echo "   Instalar desde: https://www.docker.com/products/docker-desktop"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo "❌ Docker daemon no está corriendo"
    echo "   Iniciar Docker Desktop primero"
    exit 1
fi

echo "✅ Docker OK"

# 2. Validar que estamos en el directorio correcto
if [ ! -f "Makefile" ]; then
    echo "❌ Ejecutar desde la raíz de edugo-infrastructure"
    exit 1
fi

# 3. Crear .env si no existe
if [ ! -f ".env" ]; then
    echo "📝 Creando .env desde .env.example..."
    cp .env.example .env
    echo "⚠️  Revisar .env y ajustar valores si es necesario"
else
    echo "✅ .env ya existe"
fi

# 4. Levantar servicios core
echo ""
echo "🐳 Levantando servicios Docker (PostgreSQL + MongoDB)..."
cd docker
docker-compose up -d postgres mongodb
cd ..

# 5. Esperar que servicios estén listos
echo ""
echo "⏳ Esperando que servicios estén listos..."
sleep 10

# 6. Verificar conectividad
echo ""
echo "🔍 Verificando servicios..."
if docker-compose -f docker/docker-compose.yml ps postgres | grep -q "Up"; then
    echo "✅ PostgreSQL corriendo"
else
    echo "❌ PostgreSQL no está corriendo"
    exit 1
fi

if docker-compose -f docker/docker-compose.yml ps mongodb | grep -q "Up"; then
    echo "✅ MongoDB corriendo"
else
    echo "❌ MongoDB no está corriendo"
    exit 1
fi

# 7. Ejecutar migraciones (si existe migrate.go)
echo ""
if [ -f "database/migrate.go" ]; then
    echo "📊 Ejecutando migraciones de PostgreSQL..."
    cd database
    go run migrate.go up || echo "⚠️  Migraciones no ejecutadas (verificar migrate.go)"
    cd ..
else
    echo "⚠️  migrate.go no encontrado, omitiendo migraciones automáticas"
fi

# 8. Cargar seeds
echo ""
echo "🌱 Cargando datos de prueba..."
./scripts/seed-data.sh || echo "⚠️  Seeds no cargados (verificar script)"

# 9. Resumen final
echo ""
echo "✅ ¡Setup completo!"
echo ""
echo "📊 Servicios disponibles:"
echo "   PostgreSQL: localhost:5432"
echo "     - Database: edugo_dev"
echo "     - User: edugo"
echo "     - Password: changeme"
echo ""
echo "   MongoDB: localhost:27017"
echo "     - Database: edugo"
echo ""
echo "🚀 Próximos pasos:"
echo "   1. cd ../edugo-api-admin && make run"
echo "   2. cd ../edugo-api-mobile && make run"
echo "   3. cd ../edugo-worker && make run"
echo ""
echo "🛑 Para detener: make dev-teardown"
