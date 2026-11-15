# Changelog

Este archivo registra cambios relevantes del repositorio a nivel documental y operativo.

El historial narrativo anterior a 2026-03-08 no se reconstruyo en esta fase. Los tags historicos siguen existiendo en Git, pero esta nueva serie de changelog comienza con la reorganizacion documental actual.

## [Unreleased]

### Added
- `make test-integration` (raiz y `postgres/`): integracion de `postgres/` con Docker, `GOWORK=off`, `-tags integration` y `ENABLE_INTEGRATION_TESTS=true`. Gate local obligatorio antes de un PR que toque `postgres/`; la CI no lo corre (plan 054, S.6).
- `.golangci.yml`: `run.build-tags: [integration]`, para que el lint vea los ficheros `//go:build integration`; `issues.max-same-issues: 0` y `max-issues-per-linter: 0`, para que no recorte hallazgos, y `run.allow-parallel-runners: true`.
- `govulncheck` v1.9.0 en la CI (`ci.yml`, job por modulo) y en `make vuln` / `make ci-local`; una vulnerabilidad alcanzable pone la CI en rojo (plan 054, H.2).
- Nueva base documental de fase 1 para el repositorio.
- Estructura general `docs/` centrada en alcance, mapa, procesos, arquitectura, integracion ecosistemica, automatizacion y roadmap.
- Documentacion homogenea por modulo con `README.md`, `docs/`, integracion ecosistemica y `CHANGELOG.md`.
- Estandarizacion de `Makefile` por modulo para `build`, `test`, `lint`, `fmt`, `fmt-check`, `vet` y `release-check`.
- `docs/releasing.md` como fuente de verdad del corte de versiones y GitHub Releases por modulo.
- Helpers compartidos en `make/go-module.mk`, `make/module-release.mk` y `scripts/module-release.sh`.

### Changed
- `.golangci.yml`: fuera `linters-settings` e `issues.exclude-rules`, claves de la v1 que la v2 ignora y que hacian fallar `golangci-lint config verify`. El lint no cambia: esas exclusiones nunca se aplicaron (plan 054, ola S).
- `release.yml`: los `echo >> "$GITHUB_OUTPUT"` del paso de metadatos van agrupados (SC2129); `actionlint` queda en 0.
- La documentacion heredada fue reemplazada por una fuente de verdad basada en el estado real del repo.
- `.github/workflows/ci.yml` ahora valida calidad, build y tests por modulo.
- `.github/workflows/release.yml` ahora resuelve releases por tag de modulo.

### Removed
- Documentos heredados que describian estructuras o flujos ya no alineados con el codigo actual.
