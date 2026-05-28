# Automation

## Superficies actuales

| Superficie | Funcion | Estado observado |
| --- | --- | --- |
| `Makefile` raiz | Orquestacion local, multi-modulo y release wrapper | Alineado |
| `make/go-module.mk` | Contrato comun para modulos Go | Alineado |
| `make/module-release.mk` | Contrato comun de changelog, tags y GitHub Release | Alineado |
| `postgres/Makefile` | build, test, lint, fmt, vet, migraciones, seeds y runner | Alineado |
| `mongodb/Makefile` | build, test, lint, fmt, vet, migraciones, seeds y runner | Alineado |
| `schemas/Makefile` | validacion de modulo | Alineado |
| `tools/mock-generator/Makefile` | validacion de modulo | Alineado |
| `docker/Makefile` | validacion y operacion de compose local | Alineado |
| `scripts/*.sh` | setup local, hooks y release helper | Alineado |
| `.github/workflows/ci.yml` | calidad, build y test por modulo | Alineado |
| `.github/workflows/release.yml` | release por tag de modulo | Alineado |
| `.github/actions/*` | setup Go, build Docker, coverage | Alineadas a su rol |

## Contrato comun por modulo Go

- `make build`
- `make test`
- `make lint`
- `make fmt`
- `make fmt-check`
- `make vet`
- `make tidy`
- `make check`
- `make release-check`

Los modulos `postgres` y `mongodb` agregan sus targets operativos propios para migraciones, runners y seeds. `docker/` adopta solo el contrato de release y validacion que aplica a su superficie.

## Release por modulo

- Convencion de tag: `modulo/vX.Y.Z`.
- Fuente de notas: `<modulo>/CHANGELOG.md`.
- Helper comun: `scripts/module-release.sh`.
- Wrapper comun: `make/module-release.mk`.
- Workflow remoto: `.github/workflows/release.yml`.

Targets disponibles por modulo:

- `make release-check`
- `make release-prepare VERSION=vX.Y.Z`
- `make release-notes VERSION=vX.Y.Z`
- `make release-tag VERSION=vX.Y.Z`
- `make release-push-tag VERSION=vX.Y.Z`
- `make release-github VERSION=vX.Y.Z`

El `Makefile` raiz expone wrappers equivalentes con `MODULE=<ruta-del-modulo>`.

## Validaciones locales ejecutadas

Ejecutadas el 2026-03-08:

- `make -C postgres release-check`
- `make -C mongodb release-check`
- `make -C schemas release-check`
- `make -C tools/mock-generator release-check`
- `make -C docker release-check`

Las validaciones locales pasaron despues de corregir la comprobacion de formato para que no dependiera de `mapfile`, ya que el `bash` de macOS no la soporta.

## Hooks y scripts

- `scripts/dev-setup.sh`: levanta `docker/`, espera disponibilidad basica y ejecuta `make db-bootstrap`.
- `scripts/seed-data.sh`: aplica seeds embebidos de PostgreSQL y MongoDB.
- `scripts/pre-commit-hook.sh`: valida `fmt-check`, `vet` y `test` solo para modulos Go con cambios staged.
- `scripts/reproduce-failures.sh`: reejecuta `release-check` para todos los modulos Go.
- `scripts/module-release.sh`: prepara changelog, imprime notas y crea GitHub Releases.

## Limites abiertos

- `release-prepare` actualiza el changelog, pero no crea commits ni hace push. Ese corte sigue siendo una decision humana.
- `release-github` depende de `gh` autenticado en la maquina que lo ejecute.
- `docker/Makefile` valida el compose local, pero no publica imagenes ni artefactos propios.
- El primer release real por modulo en GitHub sigue siendo la validacion final del flujo remoto.

## Reusable actions

- `setup-edugo-go`: estandariza Go y acceso a repos privados.
- `coverage-check`: calcula y valida cobertura.
- `docker-build-edugo`: build y push Docker multi-arch.
