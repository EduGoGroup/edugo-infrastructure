# Changelog - edugo-infrastructure

## [0.7.0] - 2025-11-17 - 🏗️ SCHEMA EXTENSION RELEASE

### 🚨 BREAKING CHANGES

Este release extiende las migraciones existentes 002, 003, 004 con campos adicionales y validaciones extendidas.

#### Migración Requerida

Los proyectos que usen infrastructure v0.5.0 deben:
1. Recrear base de datos (estamos en desarrollo)
2. Actualizar a v0.7.0: `go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.7.0`

### Added (postgres)

#### Soporte completo de jerarquía académica

**Migration 003 - academic_units:**
- Campo `parent_unit_id UUID` para estructura jerárquica (auto-referencia)
- Campo `description TEXT` para descripciones detalladas
- Campo `metadata JSONB` para datos flexibles adicionales
- CHECK constraint extendido con tipos: 'school', 'grade', 'class', 'section', 'club', 'department'
- Función `prevent_academic_unit_cycles()` para prevenir ciclos en jerarquía
- Trigger `prevent_cycles` que valida antes de INSERT/UPDATE
- Vista `v_academic_unit_tree` con CTE recursivo para consultar árbol completo

**Migration 002 - schools:**
- Campo `metadata JSONB` para configuraciones específicas por escuela

**Migration 004 - memberships:**
- Campo `metadata JSONB` para datos adicionales de membresía
- CHECK constraint extendido con roles: 'teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant'

#### Seeds actualizados

**postgres/seeds/academic_units.sql:**
- Datos de ejemplo con jerarquía completa
- Escuela → Grado → Sección
- Escuela → Departamento → Clase
- Ejemplos de metadata JSONB

**postgres/seeds/memberships.sql:**
- Datos de ejemplo con todos los roles
- Ejemplos de metadata JSONB

### Changed (postgres)

#### Migration 003 - academic_units
- `academic_year` ahora es NULLABLE con DEFAULT 0 (antes NOT NULL)
  - `0` = sin año académico específico (para departamentos, clubes)
  - `>0` = año académico específico (para grados, clases)

### Migration Guide

#### Si tienes datos existentes

**OPCIÓN 1: Desarrollo (Recomendado)**
```bash
# Recrear base de datos con nuevo schema
cd postgres
make migrate-down
make migrate-up
make seed
```

**OPCIÓN 2: Producción (Cuando aplique)**
Estos son cambios en migraciones base. En producción futura se requerirá:
```sql
-- Agregar columnas nuevas (ejecutar manualmente)
ALTER TABLE schools ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE academic_units ADD COLUMN parent_unit_id UUID REFERENCES academic_units(id);
ALTER TABLE academic_units ADD COLUMN description TEXT;
ALTER TABLE academic_units ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE academic_units ALTER COLUMN academic_year DROP NOT NULL;
ALTER TABLE memberships ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;

-- Actualizar CHECK constraints...
```

#### Actualizar código

```go
// Ahora puedes usar jerarquía
type AcademicUnit struct {
    ID           uuid.UUID
    ParentUnitID *uuid.UUID  // NUEVO: nullable
    SchoolID     uuid.UUID
    Type         string      // Tipos extendidos: school, grade, class, section, club, department
    Description  *string     // NUEVO: nullable
    Metadata     json.RawMessage  // NUEVO
    AcademicYear int         // CAMBIADO: ahora puede ser 0
}

// Consultar árbol completo
rows, err := db.Query("SELECT * FROM v_academic_unit_tree WHERE root_unit_id = $1", rootID)
```

---


## [0.5.0] - 2025-11-16 - 🔄 MODULAR ARCHITECTURE RELEASE

### 🚨 BREAKING CHANGES

Este release reorganiza completamente la estructura del proyecto en módulos independientes por tecnología.

#### Migración Requerida

**Antes (v0.3.0):**
```go
import "github.com/EduGoGroup/edugo-infrastructure/database"
import "github.com/EduGoGroup/edugo-infrastructure/schemas"
```

**Ahora (v0.5.0):**
```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres"
import "github.com/EduGoGroup/edugo-infrastructure/mongodb"
import "github.com/EduGoGroup/edugo-infrastructure/messaging"
```

### Added

#### Nuevos Módulos Go Independientes

- **postgres/** - Módulo de migraciones PostgreSQL
  - `go.mod`: github.com/EduGoGroup/edugo-infrastructure/postgres
  - `migrate.go`: CLI sin build tags (simplificado)
  - `migrations/`: 8 migraciones SQL
  - `seeds/`: Datos de prueba PostgreSQL
  - `Makefile`: Comandos específicos del módulo
  - `README.md`: Documentación completa

- **mongodb/** - Módulo de migraciones MongoDB
  - `go.mod`: github.com/EduGoGroup/edugo-infrastructure/mongodb
  - `migrate.go`: CLI sin build tags (simplificado)
  - `migrations/`: 6 migraciones JavaScript
  - `seeds/`: Datos de prueba MongoDB
  - `Makefile`: Comandos específicos del módulo
  - `README.md`: Documentación completa

- **messaging/** - Módulo de validación de eventos
  - `go.mod`: github.com/EduGoGroup/edugo-infrastructure/messaging
  - `validator.go`: Validador de eventos RabbitMQ
  - `events/`: 4 JSON Schemas
  - `Makefile`: Tests y benchmarks
  - `README.md`: Documentación completa

#### Makefiles Específicos por Módulo

- `postgres/Makefile`: migrate-up, migrate-down, migrate-status, seed, test
- `mongodb/Makefile`: migrate-up, migrate-down, migrate-status, seed, test
- `messaging/Makefile`: test, coverage, benchmark

#### Documentación Reorganizada

- `README.md` principal actualizado con arquitectura modular
- Sección "Uso por Proyecto" explicando importaciones selectivas
- Guías de migración desde v0.3.0
- Ejemplos de uso para api-admin, api-mobile, worker

### Changed

#### Estructura de Directorios

**Antes:**
```
edugo-infrastructure/
├── database/
│   ├── migrate.go (build tag: !mongodb)
│   ├── mongodb_migrate.go (build tag: mongodb)
│   └── migrations/
│       ├── postgres/
│       └── mongodb/
└── schemas/
```

**Ahora:**
```
edugo-infrastructure/
├── postgres/        # Módulo independiente
├── mongodb/         # Módulo independiente
└── messaging/       # Módulo independiente
```

#### Dependencias Optimizadas

- Proyectos pueden importar solo módulos necesarios
- `api-admin`: Solo postgres (sin mongo-driver, ~5MB menos)
- `api-mobile`: postgres + mongodb + messaging
- `worker`: postgres + mongodb + messaging

#### CLI Simplificado

- Removidos build tags (`!mongodb`, `mongodb`)
- Cada módulo tiene su propio `migrate.go` standalone
- Paths de migraciones simplificados (`migrations/` en lugar de `migrations/postgres/`)

### Removed

- ❌ Módulo `database/` (separado en `postgres/` y `mongodb/`)
- ❌ Módulo `schemas/` (renombrado a `messaging/`)
- ❌ Build tags complejos para compilación
- ❌ Directorio `seeds/` global (movido a cada módulo)

### Fixed

- Conflictos de compilación entre migrate.go y mongodb_migrate.go
- Dependencias innecesarias en proyectos que no usan todas las tecnologías
- Complejidad en la estructura de directorios

### Migration Guide

#### Actualizar Imports

```bash
# En tus proyectos (api-admin, api-mobile, worker)
find . -name "*.go" -type f -exec sed -i '' 's|edugo-infrastructure/database|edugo-infrastructure/postgres|g' {} +
find . -name "*.go" -type f -exec sed -i '' 's|edugo-infrastructure/schemas|edugo-infrastructure/messaging|g' {} +
```

#### Actualizar go.mod

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.5.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.5.0
go get github.com/EduGoGroup/edugo-infrastructure/messaging@v0.5.0
go mod tidy
```

#### Actualizar Scripts

Si usabas:
```bash
cd database && go run migrate.go up
cd database && go run -tags mongodb mongodb_migrate.go up
```

Ahora usa:
```bash
cd postgres && make migrate-up
cd mongodb && make migrate-up
```

---

## [0.1.0] - 2025-11-15 - 🎉 INITIAL RELEASE

### Added

#### Módulo database
- **8 migraciones PostgreSQL** con UP y DOWN
  - 001: users (roles: admin, teacher, student, guardian)
  - 002: schools (instituciones educativas)
  - 003: academic_units (cursos, clases, secciones)
  - 004: memberships (relación usuario-escuela-curso)
  - 005: materials (materiales educativos)
  - 006: assessment (quizzes con referencia a MongoDB)
  - 007: assessment_attempt (intentos de estudiantes)
  - 008: assessment_attempt_answer (respuestas individuales)
- **migrate.go** CLI para ejecutar migraciones
  - Comandos: up, down, status, create, force
  - Soporte para variables de entorno
- **TABLE_OWNERSHIP.md** documentando ownership claro

#### Módulo docker
- **docker-compose.yml** con 4 perfiles
  - default: PostgreSQL 15 + MongoDB 7.0
  - messaging: + RabbitMQ 3.12
  - cache: + Redis 7
  - tools: + PgAdmin + Mongo Express
- Healthchecks en todos los servicios
- Network compartida `edugo-network`

#### Módulo schemas
- **4 JSON Schemas** de validación de eventos
  - material.uploaded v1.0
  - assessment.generated v1.0
  - material.deleted v1.0
  - student.enrolled v1.0
- **validator.go** validador automático
  - Soporte para validar objetos Go
  - Soporte para validar JSON bytes
  - Schemas embebidos en binario
- **example_test.go** ejemplos de uso

#### Scripts
- **dev-setup.sh** setup completo automatizado
- **seed-data.sh** carga datos de prueba
- **validate-env.sh** validación de variables

#### Seeds
- **PostgreSQL seeds**
  - 3 usuarios (admin, teacher, student)
  - 2 escuelas
  - 3 materiales de prueba
- **MongoDB seeds**
  - 2 assessments de ejemplo

#### Documentación
- **README.md** documentación principal
- **EVENT_CONTRACTS.md** contratos de eventos completos
- **Makefile** con 20+ comandos útiles
- **.env.example** con todas las variables necesarias

---

## Formato de Versiones

- **MAJOR** (1.x.x): Breaking changes en schemas o migraciones
- **MINOR** (x.1.x): Nuevas features (nuevas migraciones, schemas)
- **PATCH** (x.x.1): Bug fixes

---

## [0.3.0] - 2025-11-16 - 🗄️ MONGODB MIGRATIONS RELEASE

### Added - database

#### Migraciones MongoDB (6 colecciones)
- **material_assessment** (001) - Contenido de assessments/quizzes generados por IA
  - Preguntas, opciones, respuestas correctas
  - Validación JSON Schema y índices
  - Relacionada con tabla PostgreSQL: assessment
- **material_content** (002) - Contenido procesado de materiales educativos
  - Texto extraído, estructura parseada, resumen IA
  - Full-text search en español
  - Relacionada con tabla PostgreSQL: materials
- **assessment_attempt_result** (003) - Resultados detallados de intentos
  - Respuestas del estudiante, tiempo por pregunta, score
  - Relacionada con tabla PostgreSQL: assessment_attempt
- **audit_logs** (004) - Logs de auditoría del sistema
  - Eventos de usuarios, recursos, sistema
  - TTL: 90 días de retención
- **notifications** (005) - Notificaciones para usuarios
  - In-app, push, email con prioridades y categorías
  - TTL: 30 días para archivadas
- **analytics_events** (006) - Eventos de analítica y comportamiento
  - Navegación, sesiones, interacciones
  - TTL: 365 días de retención

#### CLI MongoDB
- **mongodb_migrate.go** - CLI completo para migraciones MongoDB
  - Comandos: up, down, status, create, force
  - Ejecuta scripts JavaScript via mongosh
  - Tracking en colección schema_migrations
  - Patrón idéntico a migrate.go de PostgreSQL

#### Seeds MongoDB
- 6 archivos de seeds con datos de prueba
  - material_content.js (2 documentos)
  - assessment_attempt_result.js (2 documentos)
  - audit_logs.js (5 documentos)
  - notifications.js (4 documentos)
  - analytics_events.js (6 documentos)

#### Documentación
- **MONGODB_SCHEMA.md** - Schema completo de las 6 colecciones
  - Estructura, índices, validaciones
  - Relación con PostgreSQL
  - Queries de ejemplo y guía de uso
- **README.md** actualizado con sección MongoDB
  - Comandos de CLI con build tags
  - Variables de entorno
  - Referencias a documentación

### Changed
- **Build tags agregados** para resolución de conflictos de compilación
  - migrate.go con tag `!mongodb` (PostgreSQL, por defecto)
  - mongodb_migrate.go con tag `mongodb` (requiere `-tags mongodb`)
- **Dependencias actualizadas**
  - Agregado go.mongodb.org/mongo-driver v1.17.3

### Fixed
- Sintaxis de seeds MongoDB para compatibilidad con JavaScript
  - Reemplazado `use edugo;` por `db = db.getSiblingDB('edugo');`
  - Corregidos comentarios de ejecución en todos los seeds

---

## [0.2.0] - 2025-11-16 - 🧪 TESTS & VALIDATION RELEASE

#### Added - database
- **Tests de integración** con Testcontainers para PostgreSQL
  - 9 tests de integración: migrateUp, migrateDown, showStatus, rollback
  - Tests de transacciones, migraciones parciales, idempotencia
  - Tests de edge cases: SQL inválido, errores de conexión
- **Cobertura de tests:** 55.7% total
  - Funciones críticas >68% (migrateUp: 72.4%, showStatus: 81.2%)
- **Dependencias agregadas:** testcontainers-go v0.40.0

#### Added - schemas
- **Tests exhaustivos** para validator.go
  - 11 funciones de test con 40+ subtests
  - Tests para los 4 schemas: material.uploaded, assessment.generated, material.deleted, student.enrolled
  - Edge cases: event_type faltante, UUIDs inválidos, timestamps incorrectos
  - Tests de ValidateJSON y ValidateWithType
- **Benchmarks de performance**
  - BenchmarkValidation: ~10µs por validación
  - BenchmarkValidation10000: 10,000 eventos en ~102ms (<1s objetivo)
- **Cobertura de tests:** 92.5% (>90% objetivo superado)
  - Validate: 100%, ValidateJSON: 100%, ValidateWithType: 92.9%

#### Changed
- README.md actualizado con sección completa de Testing
  - Instrucciones para ejecutar tests
  - Métricas de cobertura documentadas
  - Ejemplos de benchmarks

#### Documentation
- Documentación de tests en README.md
- Métricas de cobertura y performance

---

**Mantenedor:** Equipo EduGo  
**Repositorio:** https://github.com/EduGoGroup/edugo-infrastructure

## [0.1.1] - 2025-11-16

### Added
- **CONTRIBUTING.md** con guía completa de desarrollo
- **sync-main-to-dev.yml** workflow de sincronización automática
- CI/CD mejorado siguiendo patrón de edugo-shared

### Changed
- **release.yml** ahora valida todos los módulos antes de publicar
- **ci.yml** con matrix strategy para Go 1.24 y 1.25
- Release workflow extrae changelog automáticamente

### Documentation
- Workflow completo documentado: feature → dev → main → tags
- Convenciones de commits estandarizadas
