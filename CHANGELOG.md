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

## Próximas Versiones Planeadas

### [0.2.0] - TBD
- Agregar módulo de MongoDB migrations
- Agregar más eventos (material.updated, assessment.failed)
- Mejorar CLI de migraciones (rollback múltiple)

---

**Mantenedor:** Equipo EduGo  
**Repositorio:** https://github.com/EduGoGroup/edugo-infrastructure
