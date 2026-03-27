# Changelog

Este changelog inicia la nueva serie documental del modulo `mongodb`.

Los tags historicos del modulo siguen existiendo en Git. El ultimo tag observado en esta fase es `mongodb/v0.53.0`, sin backfill narrativo de versiones anteriores.

## [Unreleased]

## [0.56.0] - 2026-03-27

### Changed
- Actualización de dependencias internas.

### Added
- Documentacion fase 1 del modulo.
- Documentacion fase 2 de integracion ecosistemica del modulo.
- Estructura local de `docs/` con procesos, arquitectura e integracion.
- `Makefile` uniforme con release-check, runner y seed runner.
- CLIs `cmd/runner` y `cmd/seed` para estructura, constraints, seeds canonicos y mock data.
- `internal/mongodbutil`: paquete interno con `BuildMongoURI()` y `EnvFirst()`, compartido por `cmd/runner` y `cmd/seed`.

### Changed
- `README.md` reescrito sobre el estado actual del paquete y no sobre la arquitectura heredada.
- El flujo operativo del modulo ya no depende de `mongodb/migrations/cmd/runner.go`.
- `migrations/embed.go` simplifica su API publica: eliminada `ListFunctions()` (sin callers en el ecosistema).
- `docs/architecture.md` y `docs/processes.md` actualizados para reflejar el estado actual del modulo.

### Removed
- `cmd/migrate/`: CLI de estado y force sobre coleccion `schema_migrations`. La coleccion nunca es poblada por el flujo real de migraciones (embed.go), por lo que `status` siempre devuelvia "sin migraciones registradas". Cero callers en el ecosistema.
- Targets de Makefile: `migrate-status`, `migrate-force`.
- `migrations/_deprecated/`: 6 archivos Go de estructura y constraints ya inlinados en `embed.go`. Eran ignorados por el compilador (prefijo `_`) y no tenian callers.
- `migrations/001_setup_collections.js` y `seeds/development/001_assessments.js`: scripts JS heredados, nunca embebidos ni ejecutados por ningun mecanismo Go.
- Directorio `testing/` vacio.
