# Changelog

Este changelog comienza la nueva serie documental del modulo `postgres`.

Los tags historicos del modulo siguen existiendo en Git. El ultimo tag observado en esta fase es `postgres/v0.61.0`, pero el detalle narrativo de versiones anteriores no fue reconstruido aqui.

## [Unreleased]

## [0.74.0] - 2026-03-28
### Changed
- fix seeds

## [0.73.0] - 2026-03-27
### Added
- fix seeds

## [0.72.0] - 2026-03-27
### Added
- Agregar assessments:assign y assessments:review a permisos del sistema
- Asignar ambos a super_admin, school_admin, school_director, school_coordinator, teacher
- Agregar 2 assessments manuales de ejemplo con preguntas PG, opciones y asignaciones
- Bump SchemaVersion 1.2.0 → 1.3.0

### Changed
- 003_permissions.sql — 2 permisos nuevos (68 total)
- 004_role_permissions.sql — 10 asignaciones nuevas (5 roles x 2 permisos)
- 008_assessments.sql — 2 assessments manuales + 6 preguntas PG + opciones + 3 asignaciones
- version.go — SchemaVersion 1.3.0

## [0.71.0] - 2026-03-27

### Added
- Tablas para el sistema de evaluaciones: `assessment.questions`, `assessment.question_options`, `assessment.assessment_assignments`.
- Tablas para el sistema de revisiones: `assessment.attempt_reviews`.
- Tablas para el sistema de notificaciones: `notifications.notifications`.
- Entidades GORM para todas las nuevas tablas: `Question`, `QuestionOption`, `AssessmentAssignment`, `AttemptReview`, `Notification`.

### Changed
- Modificada la tabla `assessment.assessment` para soportar `source_type` (manual/ai_generated) y hacer `mongo_document_id` opcional.
- Actualizada la tabla `assessment.assessment_attempt_answer` con referencia a `question_id` y `review_status`.

## [0.70.0] - 2026-03-25

### Added
- Nuevo rol `readonly_auditor` (scope: `school`) en seeds de produccion.
- 27 permisos de solo lectura asignados a `readonly_auditor` (todos los recursos del ecosistema).
- Permisos `context:browse_schools` y `context:browse_units` para `school_coordinator` (total: 43).
- Permisos `context:browse_schools` y `context:browse_units` para `school_assistant` (total: 15).
- Usuario de prueba `readonly@edugo.test` (U-21) en seeds de desarrollo.
- Membership `m028` para usuario ReadOnly en San Ignacio (rol: `readonly_auditor`).
- Asignacion de rol `ur27` para usuario ReadOnly en seeds de desarrollo.

### Changed
- `SchemaVersion` bumpeado de `1.1.0` a `1.1.4`.

## [0.66.0] - 2026-03-23

### Changed
- Renombrada la clave de permiso `system:settings` a `system_settings:settings` en seeds de produccion.
- Actualizada la instancia de pantalla `system-settings` para usar la nueva clave de permiso y la plantilla `settings-system-v1`.

### Added
- Nueva plantilla de pantalla `settings-system-v1` en seeds de produccion.

## [0.65.0] - 2026-03-20

### Added
- Nueva documentacion fase 1 del modulo.
- Documentacion fase 2 de integracion ecosistemica del modulo.
- Indice local en `docs/` con procesos, arquitectura e integracion.
- `Makefile` uniforme con `release-check` y wrappers de release.
- CLI `cmd/seed` para ejecutar seeds embebidos sin scripts externos.
- `internal/dbutil`: paquete interno con `BuildDBURL` y `EnvFirst`, compartido por `cmd/runner` y `cmd/seed`.
- `internal/sqlutil`: paquete interno con `IsEmptyOrComment`, compartido por `migrations` y `seeds`.

### Changed
- `README.md` reescrito para representar el estado actual del modulo y no la documentacion heredada.
- `cmd/runner` ahora usa migraciones y seeds embebidos en lugar de rutas obsoletas del filesystem.
- `migrations/embed.go` y `seeds/embed.go` simplifican su API publica: se eliminaron `GetScript`, `ListScripts` y `GetScriptsByLayer` (sin callers confirmados en todo el ecosistema).

### Removed
- `cmd/migrate/`: CLI legacy de migraciones incrementales (`up/down/status/create/force`). Era incompatible con el modelo actual de recreacion completa de schema y no tenia callers en ningun proyecto del ecosistema.
- Targets de Makefile: `migrate-up`, `migrate-down`, `migrate-status`, `migrate-create`.
