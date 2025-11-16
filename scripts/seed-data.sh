#!/bin/bash
set -e

echo "🌱 Cargando seeds de datos de prueba..."

# Variables de entorno (leer de .env si existe)
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-edugo_dev}
DB_USER=${DB_USER:-edugo}
DB_PASSWORD=${DB_PASSWORD:-changeme}

export PGPASSWORD=$DB_PASSWORD

# PostgreSQL seeds
echo "📊 Cargando seeds de PostgreSQL..."

if [ -f "seeds/postgres/users.sql" ]; then
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seeds/postgres/users.sql
    echo "  ✅ users"
else
    echo "  ⚠️  seeds/postgres/users.sql no encontrado"
fi

if [ -f "seeds/postgres/schools.sql" ]; then
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seeds/postgres/schools.sql
    echo "  ✅ schools"
else
    echo "  ⚠️  seeds/postgres/schools.sql no encontrado"
fi

if [ -f "seeds/postgres/materials.sql" ]; then
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seeds/postgres/materials.sql
    echo "  ✅ materials"
else
    echo "  ⚠️  seeds/postgres/materials.sql no encontrado"
fi

# MongoDB seeds
echo ""
echo "🍃 Cargando seeds de MongoDB..."

if [ -f "seeds/mongodb/assessments.js" ]; then
    mongosh --host $DB_HOST --eval "load('seeds/mongodb/assessments.js')"
    echo "  ✅ assessments"
else
    echo "  ⚠️  seeds/mongodb/assessments.js no encontrado"
fi

echo ""
echo "✅ Seeds cargados correctamente"
