# edugo-infrastructure

**Infraestructura compartida del ecosistema EduGo**

---

## 🎯 Propósito

Centraliza toda la infraestructura compartida entre proyectos:

- 🗄️ **Migraciones de BD** (PostgreSQL + MongoDB)
- 🐳 **Docker Compose** con perfiles
- 📋 **JSON Schemas** de eventos RabbitMQ
- 🛠️ **Scripts** automatizados
- 🌱 **Seeds** de datos de prueba

**Problema que resuelve:**
- ❌ Migraciones duplicadas entre proyectos
- ❌ Setup manual lento (1-2 horas)
- ❌ Eventos sin validación
- ❌ Configuración inconsistente

**Solución:**
- ✅ 1 fuente de verdad para infraestructura
- ✅ Setup en 5 minutos: `make dev-setup`
- ✅ Validación automática de eventos
- ✅ Ownership claro de tablas

---

## 🚀 Quick Start

```bash
# 1. Clonar
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Setup completo
make dev-setup

# ✅ Listo! Infraestructura corriendo
```

**Servicios disponibles:**
- PostgreSQL: `localhost:5432`
- MongoDB: `localhost:27017`
- RabbitMQ: `localhost:5672` (UI: http://localhost:15672)

---

## 📦 Estructura Modular

```
edugo-infrastructure/
├── database/              # Módulo: Migraciones
│   ├── migrations/
│   │   └── postgres/     # 8 migraciones SQL
│   ├── go.mod
│   └── TABLE_OWNERSHIP.md
│
├── docker/                # Módulo: Docker Compose
│   ├── docker-compose.yml  # Con profiles
│   └── README.md
│
├── schemas/               # Módulo: JSON Schemas
│   ├── events/            # 4 schemas de validación
│   ├── go.mod
│   └── README.md
│
├── scripts/               # Scripts automatizados
│   ├── dev-setup.sh
│   ├── seed-data.sh
│   └── validate-env.sh
│
├── seeds/                 # Datos de prueba
│   ├── postgres/          # users, schools, materials
│   └── mongodb/           # assessments
│
├── Makefile               # Comandos principales
├── .env.example
├── EVENT_CONTRACTS.md     # Contratos de eventos
└── README.md
```

---

## 🛠️ Comandos Principales

```bash
make help                 # Ver todos los comandos

# Desarrollo
make dev-setup            # Setup completo (primera vez)
make dev-up-core          # Solo PostgreSQL + MongoDB
make dev-up-messaging     # Core + RabbitMQ
make dev-up-full          # Todo + herramientas
make dev-teardown         # Limpiar todo

# Migraciones
make migrate-up           # Ejecutar migraciones
make migrate-status       # Ver estado
make seed                 # Cargar datos de prueba
```

---

## 🗄️ Módulo: database

**Propósito:** Migraciones centralizadas de PostgreSQL.

### Tablas Creadas

| Migración | Tabla | Owner | Usada por |
|-----------|-------|-------|-----------|
| 001 | users | infrastructure | api-admin, api-mobile, worker |
| 002 | schools | infrastructure | api-admin, api-mobile |
| 003 | academic_units | infrastructure | api-admin, api-mobile |
| 004 | memberships | infrastructure | api-admin, api-mobile |
| 005 | materials | infrastructure | api-mobile, worker |
| 006 | assessment | infrastructure | api-mobile, worker |
| 007 | assessment_attempt | infrastructure | api-mobile |
| 008 | assessment_attempt_answer | infrastructure | api-mobile |

**Ver:** `database/TABLE_OWNERSHIP.md`

### Crear Nueva Migración

```bash
cd database
go run migrate.go create "add_avatar_to_users"

# Editar archivos generados:
# - migrations/postgres/009_add_avatar_to_users.up.sql
# - migrations/postgres/009_add_avatar_to_users.down.sql

# Ejecutar
go run migrate.go up
```

---

## 🐳 Módulo: docker

**Propósito:** Docker Compose con perfiles para diferentes necesidades.

### Perfiles

| Perfil | Servicios | Cuándo usar |
|--------|-----------|-------------|
| **(default)** | PostgreSQL, MongoDB | api-admin |
| `messaging` | + RabbitMQ | api-mobile, worker |
| `cache` | + Redis | Si necesitas caché |
| `tools` | + PgAdmin, Mongo Express | Debugging |

### Ejemplos

```bash
# Solo core
docker-compose -f docker/docker-compose.yml up -d

# Core + RabbitMQ (para api-mobile, worker)
docker-compose -f docker/docker-compose.yml --profile messaging up -d

# Todo + herramientas de debugging
docker-compose -f docker/docker-compose.yml --profile messaging --profile tools up -d
```

---

## 📋 Módulo: schemas

**Propósito:** Validación automática de eventos RabbitMQ.

### Eventos Soportados

- `material.uploaded` v1.0 (api-mobile → worker)
- `assessment.generated` v1.0 (worker → api-mobile)
- `material.deleted` v1.0 (api-mobile → worker)
- `student.enrolled` v1.0 (api-admin → api-mobile)

### Uso

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

validator := schemas.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return err  // Evento inválido
}
publisher.Publish(event)  // ✅ Validado
```

**Ver:** `EVENT_CONTRACTS.md` para detalles completos

---

## 🔄 Workflow por Proyecto

### api-admin

```bash
cd edugo-infrastructure
make dev-up-core          # Solo PostgreSQL + MongoDB

cd ../edugo-api-admin
make run                  # Correr API
```

### api-mobile

```bash
cd edugo-infrastructure
make dev-up-messaging     # PostgreSQL + MongoDB + RabbitMQ

cd ../edugo-api-mobile
make run
```

### worker

```bash
cd edugo-infrastructure
make dev-up-messaging     # PostgreSQL + MongoDB + RabbitMQ

cd ../edugo-worker
make run
```

---

## 📊 Variables de Entorno

```bash
cp .env.example .env
# Editar .env si necesitas cambiar valores

# Validar configuración
make validate-env
```

**Principales variables:**
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `MONGO_URI`
- `RABBITMQ_URL`

Ver `.env.example` para lista completa.

---

## 🧪 Testing

### Tests de Integración

Este proyecto incluye tests exhaustivos con alta cobertura:

**database/migrate.go:**
- 9 tests de integración con Testcontainers
- Cobertura: 55.7% total (funciones críticas >68%)
- Tests: migrateUp, migrateDown, showStatus, rollback, idempotencia

**schemas/validator.go:**
- 11 tests exhaustivos + 40+ subtests
- Cobertura: 92.5% (>90% objetivo superado)
- Benchmarks: 10,000 validaciones en ~102ms
- Tests para los 4 schemas (material.uploaded, assessment.generated, material.deleted, student.enrolled)

### Ejecutar Tests

```bash
# Tests de database (requiere Docker)
cd database
go test -v ./...
go test -coverprofile=coverage.out

# Tests de schemas (no requiere servicios)
cd schemas
go test -v ./...
go test -bench=. -benchmem

# Benchmarks específicos
go test -bench=BenchmarkValidation10000 -benchtime=1x
```

### Tests en Otros Proyectos

Los tests de integración en api-admin, api-mobile y worker usan **Testcontainers** (no necesitan este docker-compose).

Este docker-compose es para:
- ✅ Desarrollo local manual
- ✅ Debugging con herramientas visuales
- ✅ Demos y pruebas exploratorias

---

## 📚 Documentación

- **Ownership de tablas:** `database/TABLE_OWNERSHIP.md`
- **Contratos de eventos:** `EVENT_CONTRACTS.md`
- **Docker Compose:** `docker/README.md`
- **JSON Schemas:** `schemas/README.md`

---

## 🤝 Contribuir

### Agregar Nueva Tabla

```bash
cd database
go run migrate.go create "create_nueva_tabla"

# Editar SQL generado
# Actualizar database/TABLE_OWNERSHIP.md
```

### Agregar Nuevo Evento

```bash
cd schemas/events
cp material-uploaded-v1.schema.json nuevo-evento-v1.schema.json

# Editar schema
# Actualizar EVENT_CONTRACTS.md
```

---

## 📞 Soporte

**Issues:** https://github.com/EduGoGroup/edugo-infrastructure/issues  
**Documentación completa:** Ver archivos en cada módulo

---

**Versión:** 0.1.0  
**Última actualización:** 15 de Noviembre, 2025  
**Mantenedores:** Equipo EduGo
