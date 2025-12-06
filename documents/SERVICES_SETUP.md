# 🔧 Services Setup Guide - EduGo

Guía completa para configurar y ejecutar los servicios necesarios para el desarrollo local.

---

## 📋 Servicios Requeridos

| Servicio | Puerto | Versión | Propósito |
|----------|--------|---------|-----------|
| **PostgreSQL** | 5432 | 15 | Base de datos relacional |
| **MongoDB** | 27017 | 7.0 | Base de datos de documentos |
| **RabbitMQ** | 5672 / 15672 | 3.12 | Mensajería (opcional para dev básico) |
| **Redis** | 6379 | 7 | Cache (opcional) |
| **PgAdmin** | 5050 | latest | UI para PostgreSQL (opcional) |
| **Mongo Express** | 8082 | latest | UI para MongoDB (opcional) |

---

## 🚀 Quick Start

### Opción 1: Servicios Core (Recomendado para empezar)

Solo PostgreSQL + MongoDB:

```bash
# Levantar servicios core
make dev-up-core

# Verificar que estén corriendo
make dev-ps
```

**Output esperado:**
```
PostgreSQL: localhost:5432
MongoDB: localhost:27017
```

### Opción 2: Core + Mensajería

Incluye RabbitMQ para testing de eventos:

```bash
make dev-up-messaging
```

**Output esperado:**
```
PostgreSQL: localhost:5432
MongoDB: localhost:27017
RabbitMQ: localhost:5672
RabbitMQ UI: http://localhost:15672
```

### Opción 3: Stack Completo

Todos los servicios + herramientas de desarrollo:

```bash
make dev-up-full
```

**Output esperado:**
```
PostgreSQL: localhost:5432
MongoDB: localhost:27017
RabbitMQ: localhost:5672 (UI: http://localhost:15672)
Redis: localhost:6379
PgAdmin: http://localhost:5050
Mongo Express: http://localhost:8082
```

---

## ⚙️ Configuración de Variables de Entorno

### 1. Crear archivo `.env`

```bash
cp .env.example .env
```

### 2. Variables Principales

```bash
# ===================
# PostgreSQL
# ===================
DB_HOST=localhost
DB_PORT=5432
DB_NAME=edugo_dev
DB_USER=edugo
DB_PASSWORD=changeme
DB_SSL_MODE=disable

# ===================
# MongoDB
# ===================
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=edugo

# ===================
# RabbitMQ
# ===================
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# ===================
# Redis (Opcional)
# ===================
REDIS_HOST=localhost
REDIS_PORT=6379

# ===================
# JWT
# ===================
JWT_SECRET=your-secret-key-at-least-32-characters-long-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# ===================
# OpenAI (para worker)
# ===================
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4-turbo-preview
OPENAI_MAX_TOKENS=2000

# ===================
# AWS S3
# ===================
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET=edugo-materials-dev

# ===================
# Environment
# ===================
ENVIRONMENT=local
LOG_LEVEL=debug
LOG_FORMAT=json
```

---

## 🐘 PostgreSQL

### Conexión

```bash
# Via psql
psql -h localhost -U edugo -d edugo_dev

# Connection string
postgres://edugo:changeme@localhost:5432/edugo_dev?sslmode=disable
```

### Ejecutar Migraciones

```bash
# Ejecutar todas las migraciones pendientes
make migrate-up

# Revertir última migración
make migrate-down

# Ver estado de migraciones
make migrate-status

# Crear nueva migración
make migrate-create NAME="add_new_column"
```

### Cargar Seeds

```bash
# Seeds completos
make seed

# Solo datos mínimos
make seed-minimal
```

### PgAdmin (UI)

1. Abrir http://localhost:5050
2. Login:
   - Email: `admin@edugo.com`
   - Password: `changeme`
3. Agregar servidor:
   - Host: `postgres` (nombre del container)
   - Port: `5432`
   - Database: `edugo_dev`
   - User: `edugo`
   - Password: `changeme`

---

## 🍃 MongoDB

### Conexión

```bash
# Via mongosh
mongosh mongodb://localhost:27017/edugo

# Connection string
mongodb://localhost:27017/edugo
```

### Verificar Collections

```javascript
// En mongosh
use edugo
show collections

// Esperado:
// material_assessment_worker
// material_summary
// material_event
```

### Mongo Express (UI)

1. Abrir http://localhost:8082
2. Login:
   - User: `admin`
   - Password: `changeme`
3. Seleccionar database `edugo`

---

## 🐰 RabbitMQ

### Conexión

```bash
# AMQP URL
amqp://guest:guest@localhost:5672/
```

### Management UI

1. Abrir http://localhost:15672
2. Login:
   - User: `guest`
   - Password: `guest`

### Exchanges y Queues

Crear manualmente o via código en primer run:

```
Exchanges:
- edugo.materials (topic)
- edugo.assessments (topic)
- edugo.students (topic)

Queues:
- worker.materials.process
- api-mobile.assessments.ready
- api-mobile.students.enrolled
```

---

## 🔴 Redis

### Conexión

```bash
# Via redis-cli
redis-cli -h localhost -p 6379

# Connection string
redis://localhost:6379
```

### Verificar

```bash
redis-cli ping
# Output: PONG
```

---

## 🐳 Docker Compose Commands

### Comandos Make

```bash
# Levantar servicios
make dev-up-core        # Solo PostgreSQL + MongoDB
make dev-up-messaging   # + RabbitMQ
make dev-up-full        # Todos los servicios

# Estado y logs
make dev-ps             # Ver estado de containers
make dev-logs           # Ver logs en tiempo real

# Detener
make dev-down           # Detener (mantener datos)
make dev-teardown       # Detener y eliminar volúmenes

# Reset completo
make dev-reset          # Teardown + setup desde cero
```

### Comandos Docker Compose Directos

```bash
cd docker

# Ver containers
docker-compose ps

# Logs de un servicio específico
docker-compose logs -f postgres
docker-compose logs -f mongodb

# Reiniciar un servicio
docker-compose restart postgres

# Ejecutar comando en container
docker-compose exec postgres psql -U edugo -d edugo_dev
docker-compose exec mongodb mongosh edugo
```

---

## 🔍 Health Checks

### Verificar PostgreSQL

```bash
# Desde host
pg_isready -h localhost -p 5432 -U edugo

# Desde container
docker exec edugo-postgres pg_isready -U edugo
```

### Verificar MongoDB

```bash
# Desde container
docker exec edugo-mongodb mongosh --eval "db.adminCommand('ping')"
```

### Verificar RabbitMQ

```bash
# Desde container
docker exec edugo-rabbitmq rabbitmq-diagnostics ping
```

### Verificar Todos

```bash
make status
```

---

## 🔧 Troubleshooting

### Puerto ya en uso

```bash
# Verificar qué proceso usa el puerto
lsof -i :5432
lsof -i :27017

# Matar proceso
kill -9 <PID>
```

### Containers no inician

```bash
# Ver logs detallados
docker-compose logs postgres
docker-compose logs mongodb

# Reiniciar Docker
# macOS: Reiniciar Docker Desktop

# Limpiar todo y reintentar
make dev-teardown
make dev-setup
```

### Migraciones fallan

```bash
# Verificar conexión a PostgreSQL
psql -h localhost -U edugo -d edugo_dev -c "SELECT 1"

# Ver estado de migraciones
make migrate-status

# Forzar versión (solo si necesario)
cd postgres && go run cmd/migrate/migrate.go force <version>
```

### MongoDB sin collections

```bash
# Verificar que MongoDB está corriendo
docker exec edugo-mongodb mongosh --eval "db.stats()"

# Crear collections manualmente
docker exec edugo-mongodb mongosh edugo --eval "db.createCollection('material_assessment_worker')"
```

---

## 📊 Configuración por Ambiente

### Local (Desarrollo)

```bash
ENVIRONMENT=local
DB_HOST=localhost
MONGO_HOST=localhost
LOG_LEVEL=debug
```

### Docker (CI/Testing)

```bash
ENVIRONMENT=docker
DB_HOST=postgres          # Nombre del container
MONGO_HOST=mongodb        # Nombre del container
RABBITMQ_HOST=rabbitmq    # Nombre del container
```

### Producción

```bash
ENVIRONMENT=prod
DB_HOST=<rds-endpoint>
MONGO_HOST=<documentdb-endpoint>
RABBITMQ_HOST=<amazonmq-endpoint>
LOG_LEVEL=info
```

---

## 🔐 Credenciales por Defecto (Solo Desarrollo)

| Servicio | Usuario | Password |
|----------|---------|----------|
| PostgreSQL | edugo | changeme |
| MongoDB | (sin auth) | - |
| RabbitMQ | guest | guest |
| PgAdmin | admin@edugo.com | changeme |
| Mongo Express | admin | changeme |

⚠️ **IMPORTANTE:** Cambiar todas las credenciales en producción.

---

## 🌐 URLs de Acceso Local

| Servicio | URL |
|----------|-----|
| PostgreSQL | `localhost:5432` |
| MongoDB | `localhost:27017` |
| RabbitMQ | `localhost:5672` |
| RabbitMQ Management | http://localhost:15672 |
| Redis | `localhost:6379` |
| PgAdmin | http://localhost:5050 |
| Mongo Express | http://localhost:8082 |

---

**Última actualización:** Diciembre 2024
