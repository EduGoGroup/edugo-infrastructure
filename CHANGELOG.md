# Changelog - edugo-infrastructure

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
