# edugo-infrastructure

**Infraestructura compartida modular del ecosistema EduGo**

---

## 🎯 Propósito

Centraliza toda la infraestructura compartida entre proyectos con módulos independientes:

- 🐘 **postgres/** - Migraciones PostgreSQL
- 🍃 **mongodb/** - Migraciones MongoDB
- 📨 **messaging/** - Validación de eventos RabbitMQ
- 🐳 **docker/** - Docker Compose con perfiles
- 🛠️ **scripts/** - Scripts automatizados

**Problema que resuelve:**
- ❌ Migraciones duplicadas entre proyectos
- ❌ Dependencias innecesarias (cada proyecto solo usa lo que necesita)
- ❌ Setup manual lento
- ❌ Eventos sin validación

**Solución:**
- ✅ Módulos independientes por tecnología
- ✅ Importar solo lo necesario
- ✅ Setup en 5 minutos
- ✅ Validación automática de eventos

---

## 🚀 Quick Start

```bash
# 1. Clonar
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Levantar servicios (docker)
make dev-up-core          # PostgreSQL + MongoDB
make dev-up-messaging     # + RabbitMQ

# 3. Ejecutar migraciones
cd postgres && make migrate-up
cd ../mongodb && make migrate-up

# ✅ Listo!
```

---

## 📦 Estructura Modular

```
edugo-infrastructure/
├── postgres/              # Módulo Go: Migraciones PostgreSQL
│   ├── go.mod            # github.com/EduGoGroup/edugo-infrastructure/postgres
│   ├── migrate.go        # CLI de migraciones
│   ├── migrations/       # 8 migraciones SQL
│   ├── seeds/            # Datos de prueba
│   ├── Makefile
│   └── README.md
│
├── mongodb/               # Módulo Go: Migraciones MongoDB
│   ├── go.mod            # github.com/EduGoGroup/edugo-infrastructure/mongodb
│   ├── migrate.go        # CLI de migraciones
│   ├── migrations/       # 6 migraciones JavaScript
│   ├── seeds/            # Datos de prueba
│   ├── Makefile
│   └── README.md
│
├── messaging/             # Módulo Go: Validación de eventos
│   ├── go.mod            # github.com/EduGoGroup/edugo-infrastructure/messaging
│   ├── validator.go      # Validador de eventos
│   ├── events/           # 4 JSON Schemas
│   ├── Makefile
│   └── README.md
│
├── docker/                # Docker Compose con perfiles
│   ├── docker-compose.yml
│   └── README.md
│
├── scripts/               # Scripts automatizados
│   ├── dev-setup.sh
│   └── validate-env.sh
│
├── docs/                  # Documentación
│   ├── TABLE_OWNERSHIP.md
│   ├── MONGODB_SCHEMA.md
│   └── ...
│
├── Makefile               # Comandos globales
├── .env.example
├── EVENT_CONTRACTS.md
└── README.md
```

---

## 🛠️ Uso por Proyecto

### api-admin (solo PostgreSQL)

```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres"

// Solo importa postgres, sin dependencias de MongoDB
```

```bash
cd edugo-infrastructure
make dev-up-core          # Solo PostgreSQL + MongoDB (básico)

cd postgres
make migrate-up
```

### api-mobile (PostgreSQL + MongoDB + RabbitMQ)

```go
import (
    "github.com/EduGoGroup/edugo-infrastructure/postgres"
    "github.com/EduGoGroup/edugo-infrastructure/mongodb"
    "github.com/EduGoGroup/edugo-infrastructure/messaging"
)
```

```bash
cd edugo-infrastructure
make dev-up-messaging     # PostgreSQL + MongoDB + RabbitMQ

cd postgres && make migrate-up
cd ../mongodb && make migrate-up
```

### worker (PostgreSQL + MongoDB + RabbitMQ)

```go
import (
    "github.com/EduGoGroup/edugo-infrastructure/postgres"
    "github.com/EduGoGroup/edugo-infrastructure/mongodb"
    "github.com/EduGoGroup/edugo-infrastructure/messaging"
)
```

---

## 📋 Módulos Disponibles

### 1. postgres/

**Propósito:** Migraciones de PostgreSQL

**Tablas:** users, schools, academic_units, memberships, materials, assessment, assessment_attempt, assessment_attempt_answer

**Uso:**
```bash
cd postgres
make migrate-up          # Ejecutar migraciones
make migrate-status      # Ver estado
make migrate-create name="nueva_tabla"
```

**Importar:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres"
```

**Ver:** [postgres/README.md](postgres/README.md)

---

### 2. mongodb/

**Propósito:** Migraciones de MongoDB

**Colecciones:** material_assessment, material_content, assessment_attempt_result, audit_logs, notifications, analytics_events

**Uso:**
```bash
cd mongodb
make migrate-up          # Ejecutar migraciones
make migrate-status      # Ver estado
make migrate-create name="nueva_coleccion"
```

**Importar:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/mongodb"
```

**Ver:** [mongodb/README.md](mongodb/README.md)

---

### 3. messaging/

**Propósito:** Validación de eventos RabbitMQ

**Eventos:** material.uploaded, assessment.generated, material.deleted, student.enrolled

**Uso:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

validator := messaging.NewEventValidator()
if err := validator.Validate(event); err != nil {
    return err
}
```

**Ver:** [messaging/README.md](messaging/README.md)

---

## 🐳 Docker

Perfiles disponibles:

| Perfil | Servicios | Cuándo usar |
|--------|-----------|-------------|
| **core** | PostgreSQL, MongoDB | api-admin |
| **messaging** | + RabbitMQ | api-mobile, worker |
| **cache** | + Redis | Si necesitas caché |
| **tools** | + PgAdmin, Mongo Express | Debugging |

```bash
make dev-up-core          # PostgreSQL + MongoDB
make dev-up-messaging     # + RabbitMQ
make dev-up-cache         # + Redis
make dev-up-full          # Todo
make dev-teardown         # Limpiar
```

---

## 🧪 Testing

Cada módulo tiene sus propios tests:

```bash
# PostgreSQL
cd postgres && make test

# MongoDB
cd mongodb && make test

# Messaging
cd messaging && make test
cd messaging && make benchmark
```

---

## 📚 Documentación

- **PostgreSQL Tables:** [docs/TABLE_OWNERSHIP.md](docs/TABLE_OWNERSHIP.md)
- **MongoDB Schemas:** [docs/MONGODB_SCHEMA.md](docs/MONGODB_SCHEMA.md)
- **Event Contracts:** [EVENT_CONTRACTS.md](EVENT_CONTRACTS.md)
- **Integration Guide:** [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)

---

## 🤝 Contribuir

### Agregar migración PostgreSQL

```bash
cd postgres
make migrate-create name="add_column_to_users"
# Editar archivos SQL generados
make migrate-up
```

### Agregar migración MongoDB

```bash
cd mongodb
make migrate-create name="add_new_collection"
# Editar archivos JavaScript generados
make migrate-up
```

### Agregar evento

```bash
cd messaging/events
cp material-uploaded-v1.schema.json nuevo-evento-v1.schema.json
# Editar schema
# Actualizar EVENT_CONTRACTS.md
```

---

## 🔄 Versionamiento

**Versión actual:** 0.5.0

Este proyecto usa **versionamiento único** para el repositorio completo, aunque está organizado en módulos Go independientes.

**Semantic Versioning:**
- **MAJOR (1.x.x):** Breaking changes en estructura modular o APIs
- **MINOR (x.1.x):** Nuevas features (nuevas migraciones, schemas, módulos)
- **PATCH (x.x.1):** Bug fixes

---

## 📞 Soporte

**Issues:** https://github.com/EduGoGroup/edugo-infrastructure/issues  
**Versión:** 0.5.0  
**Última actualización:** 16 de Noviembre, 2025  
**Mantenedores:** Equipo EduGo
