#!/bin/bash

echo "🔍 Validando variables de entorno..."

# Leer .env si existe
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  .env no encontrado, usando valores por defecto"
fi

ERRORS=0

# Validar PostgreSQL
if [ -z "$DB_HOST" ]; then
    echo "❌ DB_HOST no definido"
    ERRORS=$((ERRORS+1))
else
    echo "✅ DB_HOST=$DB_HOST"
fi

if [ -z "$DB_NAME" ]; then
    echo "❌ DB_NAME no definido"
    ERRORS=$((ERRORS+1))
else
    echo "✅ DB_NAME=$DB_NAME"
fi

if [ -z "$DB_USER" ]; then
    echo "❌ DB_USER no definido"
    ERRORS=$((ERRORS+1))
else
    echo "✅ DB_USER=$DB_USER"
fi

# Validar MongoDB
if [ -z "$MONGO_HOST" ]; then
    echo "⚠️  MONGO_HOST no definido (usando localhost)"
else
    echo "✅ MONGO_HOST=$MONGO_HOST"
fi

# Validar RabbitMQ (opcional)
if [ -z "$RABBITMQ_HOST" ]; then
    echo "⚠️  RABBITMQ_HOST no definido (opcional para desarrollo core)"
else
    echo "✅ RABBITMQ_HOST=$RABBITMQ_HOST"
fi

echo ""
if [ $ERRORS -eq 0 ]; then
    echo "✅ Todas las variables críticas están configuradas"
    exit 0
else
    echo "❌ $ERRORS variables críticas faltantes"
    echo "   Copiar .env.example a .env y configurar"
    exit 1
fi
