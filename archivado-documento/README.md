# edugo-infrastructure

**Infraestructura compartida modular del ecosistema EduGo**

![CI Status](https://github.com/EduGoGroup/edugo-infrastructure/workflows/CI/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

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

### CI/CD

**Workflows automáticos:**
- ✅ Tests unitarios en cada PR (`-short` flag)
- ✅ Race detection habilitado (`-race`)
- ✅ Go 1.25 estandarizado
- ✅ Pre-commit hooks para calidad de código

**Ver configuración completa:** [docs/WORKFLOWS.md](docs/WORKFLOWS.md)

### Tests Locales

```bash
# Tests unitarios (rápidos, sin servicios externos)
cd postgres && go test -short -v ./...
cd mongodb && go test -short -v ./...
cd messaging && go test -short -v ./...

# Tests de integración (requieren Docker)
cd postgres && ENABLE_INTEGRATION_TESTS=true go test -v ./...
cd mongodb && ENABLE_INTEGRATION_TESTS=true go test -v ./...

# Benchmarks
cd messaging && go test -bench=. -benchmem
```

### Pre-commit Hooks

Instala hooks locales para validar código antes de commit:

```bash
# Una sola vez por clon del repo
./scripts/setup-hooks.sh

# Los hooks ejecutarán automáticamente:
# 1. go fmt (formato)
# 2. go vet (análisis estático)
# 3. go mod tidy check
# 4. go test -short (tests unitarios)
```

---

## 📚 Documentación

### Infraestructura y Base de Datos
- **PostgreSQL Tables:** [docs/TABLE_OWNERSHIP.md](docs/TABLE_OWNERSHIP.md)
- **MongoDB Schemas:** [docs/MONGODB_SCHEMA.md](docs/MONGODB_SCHEMA.md)
- **Event Contracts:** [EVENT_CONTRACTS.md](EVENT_CONTRACTS.md)
- **Integration Guide:** [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)

### CI/CD y Desarrollo
- **Workflows y Testing:** [docs/WORKFLOWS.md](docs/WORKFLOWS.md) ⭐
- **Sprint Planning:** [docs/cicd/](docs/cicd/)
- **Pre-commit Hooks:** [scripts/pre-commit-hook.sh](scripts/pre-commit-hook.sh)

---

## 🤝 Contribuir

### Setup Inicial

```bash
# 1. Clonar repo
git clone git@github.com:EduGoGroup/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Instalar pre-commit hooks
./scripts/setup-hooks.sh

# 3. Verificar Go version
go version  # Debe ser 1.25+

# 4. Validar setup
for module in postgres mongodb messaging schemas; do
  cd $module && go mod download && cd ..
done
```

### Agregar migración PostgreSQL

```bash
cd postgres
make migrate-create name="add_column_to_users"
# Editar archivos SQL generados
make migrate-up

# Validar
go test -short -v ./...
```

### Agregar migración MongoDB

```bash
cd mongodb
make migrate-create name="add_new_collection"
# Editar archivos JavaScript generados
make migrate-up

# Validar
go test -short -v ./...
```

### Agregar evento

```bash
cd messaging/events
cp material-uploaded-v1.schema.json nuevo-evento-v1.schema.json
# Editar schema
# Actualizar EVENT_CONTRACTS.md

# Validar
cd ../
go test -short -v ./...
```

### Workflow de Contribución

1. **Crear branch:**
   ```bash
   git checkout -b feature/descripcion-breve
   ```

2. **Hacer cambios** (los pre-commit hooks validarán automáticamente)

3. **Ejecutar tests:**
   ```bash
   # Unit tests (obligatorio)
   go test -short -race -v ./...

   # Integration tests (recomendado antes de merge)
   ENABLE_INTEGRATION_TESTS=true go test -v ./...
   ```

4. **Commit con conventional commits:**
   ```bash
   git commit -m "feat(postgres): add new migration for users"
   git commit -m "fix(mongodb): correct schema validation"
   git commit -m "docs: update WORKFLOWS.md"
   ```

5. **Push y crear PR:**
   ```bash
   git push -u origin feature/descripcion-breve
   ```

6. **Esperar CI** antes de merge (debe estar ✅ verde)

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

---

## 📋 Último Plan de Trabajo

**FASE 1: UI Database Infrastructure** - [Ver plan completo](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/README.md)

Implementación de 3 nuevas tablas PostgreSQL para soportar UI Roadmap de EduGo:
- **`user_active_context`** - Contexto/escuela activa del usuario para filtrado en UI
- **`user_favorites`** - Materiales marcados como favoritos
- **`user_activity_log`** - Log de actividades del usuario para analytics e historial

**Estado**: 🔄 En planificación  
**Rama**: `feature/fase1-ui-database-infrastructure`  
**Fecha**: 1 de Diciembre, 2025  
**Bloquea**: FASE 2 (APIs), FASE 4 (UI Estudiantes)

**Documentación del plan**:
- [Resumen ejecutivo](./docs/specs/fase1-ui-database/README.md)
- [Análisis técnico detallado](./docs/specs/fase1-ui-database/ANALISIS-TECNICO.md)
- [Plan de fases y pasos](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner.md)
- [Estrategia de commits](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Planner-commit.md)
- [Archivos afectados](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Files-affected.md)
- [Tests unitarios](./docs/specs/fase1-ui-database/Plan-Fase1-UI-Database/Test-unit.md)
