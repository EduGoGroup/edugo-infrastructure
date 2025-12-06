# 👨‍💻 Development Guide - EduGo

Guía para desarrolladores que trabajan con el proyecto edugo-infrastructure.

---

## 🚀 Primeros Pasos

### Requisitos Previos

| Herramienta | Versión Mínima | Instalación |
|-------------|----------------|-------------|
| **Go** | 1.22+ | [golang.org](https://golang.org/dl/) |
| **Docker** | 24+ | [docker.com](https://www.docker.com/get-started) |
| **Docker Compose** | 2.0+ | Incluido en Docker Desktop |
| **Make** | 3.81+ | Incluido en macOS/Linux |
| **Git** | 2.30+ | [git-scm.com](https://git-scm.com/) |

### Setup Inicial

```bash
# 1. Clonar repositorio
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con valores locales

# 3. Verificar requisitos
go version          # go1.22+
docker --version    # Docker 24+
make --version      # GNU Make 3.81+

# 4. Setup completo (primera vez)
make dev-setup
```

---

## 📁 Estructura del Proyecto

```
edugo-infrastructure/
├── .github/                   # GitHub Actions, templates
├── docker/                    # Docker Compose
│   └── docker-compose.yml
├── postgres/                  # Módulo PostgreSQL
│   ├── entities/              # Go structs
│   ├── migrations/            # SQL migrations
│   └── cmd/                   # CLI tools
├── mongodb/                   # Módulo MongoDB
│   ├── entities/              # Go structs
│   ├── migrations/            # Index scripts
│   └── seeds/                 # Test data
├── schemas/                   # JSON Schemas
│   └── events/                # Event schemas
├── messaging/                 # Event validation
│   └── events/                # Event types
├── seeds/                     # Seed data
│   ├── postgres/
│   └── mongodb/
├── scripts/                   # Utility scripts
├── tools/                     # Internal tools
├── documents/                 # Documentation
├── .env.example               # Environment template
├── Makefile                   # Build commands
└── README.md                  # Project readme
```

---

## 🔨 Comandos Make Principales

### Desarrollo Local

```bash
make help               # Ver todos los comandos disponibles

# Docker
make dev-up-core        # Levantar PostgreSQL + MongoDB
make dev-up-messaging   # + RabbitMQ
make dev-up-full        # Todos los servicios
make dev-down           # Detener servicios
make dev-teardown       # Detener y eliminar volúmenes
make dev-ps             # Estado de containers
make dev-logs           # Ver logs

# Migraciones
make migrate-up         # Ejecutar migraciones
make migrate-down       # Revertir última migración
make migrate-status     # Ver estado
make migrate-create NAME="nombre"  # Crear migración

# Seeds
make seed               # Cargar datos de prueba
make seed-minimal       # Solo datos mínimos

# Calidad
make lint               # Ejecutar linter
make fmt                # Formatear código
make vet                # Analizar código
```

---

## 🗃️ Trabajar con Migraciones

### Crear Nueva Migración

```bash
# Sintaxis
make migrate-create NAME="descripcion_de_la_migracion"

# Ejemplo
make migrate-create NAME="add_phone_to_users"
```

Esto crea dos archivos:
- `postgres/migrations/XXX_add_phone_to_users.up.sql`
- `postgres/migrations/XXX_add_phone_to_users.down.sql`

### Escribir Migraciones

**UP (aplicar cambios):**
```sql
-- postgres/migrations/XXX_add_phone_to_users.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_phone ON users(phone);
```

**DOWN (revertir cambios):**
```sql
-- postgres/migrations/XXX_add_phone_to_users.down.sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Ejecutar Migraciones

```bash
# Aplicar todas las pendientes
make migrate-up

# Revertir última
make migrate-down

# Ver estado
make migrate-status
```

---

## 📝 Crear Nuevas Entities

### PostgreSQL Entity

```go
// postgres/entities/new_entity.go
package entities

import (
    "time"
    "github.com/google/uuid"
)

// NewEntity representa la tabla 'new_entities' en PostgreSQL.
//
// Migración: XXX_create_new_entities.up.sql
// Usada por: api-mobile, api-administracion
type NewEntity struct {
    ID          uuid.UUID  `db:"id"`
    Name        string     `db:"name"`
    Description *string    `db:"description"` // Nullable
    Metadata    []byte     `db:"metadata"`    // JSONB
    IsActive    bool       `db:"is_active"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at"`
    DeletedAt   *time.Time `db:"deleted_at"`  // Soft delete
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (NewEntity) TableName() string {
    return "new_entities"
}
```

### MongoDB Entity

```go
// mongodb/entities/new_document.go
package entities

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// NewDocument representa la collection 'new_documents' en MongoDB.
//
// Seed: mongodb/seeds/new_documents.js
// Usada por: worker
type NewDocument struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    RefID     string             `bson:"ref_id"` // UUID de PostgreSQL
    Data      []DataItem         `bson:"data"`
    Status    string             `bson:"status"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// DataItem tipo embebido
type DataItem struct {
    Key   string      `bson:"key"`
    Value interface{} `bson:"value"`
}

// CollectionName retorna el nombre de la collection en MongoDB
func (NewDocument) CollectionName() string {
    return "new_documents"
}
```

---

## 📬 Crear Nuevos Eventos

### 1. Crear JSON Schema

```json
// schemas/events/new-event-v1.schema.json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://github.com/EduGoGroup/edugo-infrastructure/schemas/events/new-event-v1",
  "title": "New Event",
  "type": "object",
  "required": ["event_id", "event_type", "event_version", "timestamp", "payload"],
  "properties": {
    "event_id": { "type": "string", "format": "uuid" },
    "event_type": { "type": "string", "const": "new.event" },
    "event_version": { "type": "string", "const": "1.0" },
    "timestamp": { "type": "string", "format": "date-time" },
    "payload": {
      "type": "object",
      "required": ["entity_id"],
      "properties": {
        "entity_id": { "type": "string", "format": "uuid" },
        "action": { "type": "string" }
      }
    }
  }
}
```

### 2. Crear Tipo Go

```go
// messaging/events/new_event.go
package events

import "time"

type NewEvent struct {
    EventID      string           `json:"event_id"`
    EventType    string           `json:"event_type"`
    EventVersion string           `json:"event_version"`
    Timestamp    time.Time        `json:"timestamp"`
    Payload      NewEventPayload  `json:"payload"`
}

type NewEventPayload struct {
    EntityID string `json:"entity_id"`
    Action   string `json:"action,omitempty"`
}

// NewNewEvent crea un evento con valores por defecto
func NewNewEvent(entityID, action string) *NewEvent {
    return &NewEvent{
        EventID:      generateUUID(),
        EventType:    "new.event",
        EventVersion: "1.0",
        Timestamp:    time.Now().UTC(),
        Payload: NewEventPayload{
            EntityID: entityID,
            Action:   action,
        },
    }
}
```

### 3. Registrar en Validador

```go
// messaging/validator.go
func init() {
    RegisterSchema("new.event", "1.0", "new-event-v1.schema.json")
}
```

---

## 🧪 Testing

### Ejecutar Tests

```bash
# Tests de un módulo
cd postgres && go test ./...
cd mongodb && go test ./...
cd messaging && go test ./...

# Tests con coverage
cd postgres && go test -cover ./...

# Tests verbosos
cd postgres && go test -v ./...
```

### Escribir Tests

```go
// postgres/entities/user_test.go
package entities_test

import (
    "testing"
    "github.com/google/uuid"
    "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestUser_TableName(t *testing.T) {
    user := &entities.User{}
    if got := user.TableName(); got != "users" {
        t.Errorf("TableName() = %v, want %v", got, "users")
    }
}

func TestUser_Fields(t *testing.T) {
    user := &entities.User{
        ID:        uuid.New(),
        Email:     "test@example.com",
        FirstName: "Test",
        LastName:  "User",
        Role:      "student",
        IsActive:  true,
    }
    
    if user.Email != "test@example.com" {
        t.Errorf("Email = %v, want %v", user.Email, "test@example.com")
    }
}
```

---

## 🔍 Linting y Formateo

### Configuración de golangci-lint

El proyecto usa `.golangci.yml`:

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused

linters-settings:
  errcheck:
    check-blank: true
```

### Ejecutar Linter

```bash
# Todo el proyecto
make lint

# Módulo específico
cd postgres && golangci-lint run

# Con fix automático
golangci-lint run --fix
```

### Formatear Código

```bash
make fmt
# o
go fmt ./...
```

---

## 📊 Convenciones de Código

### Naming

| Tipo | Convención | Ejemplo |
|------|------------|---------|
| Packages | lowercase | `entities` |
| Types | PascalCase | `User`, `MaterialAssessment` |
| Functions | PascalCase (exported) | `GetUser()` |
| Variables | camelCase | `userID`, `createdAt` |
| Constants | PascalCase or SCREAMING_SNAKE | `MaxRetries`, `DEFAULT_TIMEOUT` |

### Tags de Struct

```go
// PostgreSQL con sqlx
Field string `db:"field_name"`

// MongoDB con bson
Field string `bson:"field_name"`

// JSON para APIs
Field string `json:"field_name"`
```

### Documentación

```go
// User representa un usuario del sistema.
// Esta entity refleja la tabla 'users' en PostgreSQL.
//
// Migración: 001_create_users.up.sql
// Usada por: api-mobile, api-administracion, worker
type User struct {
    // ...
}
```

---

## 🔄 Flujo de Trabajo Git

### Branches

| Branch | Propósito |
|--------|-----------|
| `main` | Producción estable |
| `develop` | Integración de features |
| `feature/*` | Nuevas características |
| `fix/*` | Corrección de bugs |
| `hotfix/*` | Fixes urgentes para prod |

### Commits

Usar Conventional Commits:

```
feat: add phone field to users entity
fix: correct assessment query
docs: update architecture diagram
chore: update dependencies
refactor: simplify migration runner
```

### Pull Requests

1. Crear branch desde `develop`
2. Hacer cambios y commits
3. Push y crear PR
4. Pasar CI checks
5. Code review
6. Merge a `develop`

---

## 🐛 Debugging

### Logs de Docker

```bash
# Todos los logs
make dev-logs

# Logs de servicio específico
docker-compose -f docker/docker-compose.yml logs -f postgres
docker-compose -f docker/docker-compose.yml logs -f mongodb
```

### Conectar a Bases de Datos

```bash
# PostgreSQL
docker exec -it edugo-postgres psql -U edugo -d edugo_dev

# MongoDB
docker exec -it edugo-mongodb mongosh edugo

# Redis
docker exec -it edugo-redis redis-cli
```

### Verificar Estado

```bash
make status
make dev-ps
```

---

## 🚀 CI/CD

### GitHub Actions

El proyecto incluye:
- **Lint:** Ejecuta golangci-lint
- **Test:** Ejecuta tests de todos los módulos
- **Build:** Verifica compilación

### Pre-commit Hooks (Recomendado)

```bash
# Instalar pre-commit
pip install pre-commit

# Configurar hooks
pre-commit install
```

Archivo `.pre-commit-config.yaml`:
```yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: go fmt ./...
        language: system
        types: [go]
      - id: go-vet
        name: go vet
        entry: go vet ./...
        language: system
        types: [go]
```

---

## 📚 Recursos Adicionales

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [MongoDB Manual](https://www.mongodb.com/docs/manual/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/tutorials)
- [JSON Schema](https://json-schema.org/)

---

**Última actualización:** Diciembre 2024
