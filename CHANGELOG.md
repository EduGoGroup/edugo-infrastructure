# Changelog

Este archivo registra cambios relevantes del repositorio a nivel documental y operativo.

El historial narrativo anterior a 2026-03-08 no se reconstruyo en esta fase. Los tags historicos siguen existiendo en Git, pero esta nueva serie de changelog comienza con la reorganizacion documental actual.

## [Unreleased]

### Added
- Nueva base documental de fase 1 para el repositorio.
- Estructura general `docs/` centrada en alcance, mapa, procesos, arquitectura, integracion ecosistemica, automatizacion y roadmap.
- Documentacion homogenea por modulo con `README.md`, `docs/`, integracion ecosistemica y `CHANGELOG.md`.
- Estandarizacion de `Makefile` por modulo para `build`, `test`, `lint`, `fmt`, `fmt-check`, `vet` y `release-check`.
- `docs/releasing.md` como fuente de verdad del corte de versiones y GitHub Releases por modulo.
- Helpers compartidos en `make/go-module.mk`, `make/module-release.mk` y `scripts/module-release.sh`.

### Changed
- La documentacion heredada fue reemplazada por una fuente de verdad basada en el estado real del repo.
- `.github/workflows/ci.yml` ahora valida calidad, build y tests por modulo.
- `.github/workflows/release.yml` ahora resuelve releases por tag de modulo.

### Removed
- Documentos heredados que describian estructuras o flujos ya no alineados con el codigo actual.
