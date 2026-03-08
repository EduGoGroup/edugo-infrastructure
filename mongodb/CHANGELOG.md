# Changelog

Este changelog inicia la nueva serie documental del modulo `mongodb`.

Los tags historicos del modulo siguen existiendo en Git. El ultimo tag observado en esta fase es `mongodb/v0.53.0`, sin backfill narrativo de versiones anteriores.

## [Unreleased]

### Added
- Documentacion fase 1 del modulo.
- Documentacion fase 2 de integracion ecosistemica del modulo.
- Estructura local de `docs/` con procesos, arquitectura e integracion.
- `Makefile` uniforme con release-check, runner y seed runner.
- CLIs `cmd/runner` y `cmd/seed` para estructura, constraints, seeds canonicos y mock data.

### Changed
- `README.md` reescrito sobre el estado actual del paquete y no sobre la arquitectura heredada.
- El flujo operativo del modulo ya no depende de `mongodb/migrations/cmd/runner.go`.
