# Releasing

## Convencion

- Cada modulo se versiona con tags `modulo/vX.Y.Z`.
- Las notas salen de `<modulo>/CHANGELOG.md`.
- `## [Unreleased]` es obligatoria en cada changelog.

## Prerrequisitos

- `git` con acceso de push al repositorio.
- `gh` autenticado si se va a crear GitHub Release desde la CLI.
- `golangci-lint` disponible para modulos Go.
- Go instalado para `postgres`, `mongodb`, `schemas` y `tools/mock-generator`.
- Docker Compose disponible si se valida el modulo `docker`.

## Flujo recomendado

1. Registrar cambios del modulo en `CHANGELOG.md` bajo `## [Unreleased]`.
2. Validar el modulo:
   - `make -C <modulo> release-check`
   - o `make release-check MODULE=<modulo>`
3. Congelar el changelog para la version:
   - `make -C <modulo> release-prepare VERSION=vX.Y.Z`
4. Revisar el diff del changelog y confirmar la seccion nueva.
5. Crear el tag local:
   - `make -C <modulo> release-tag VERSION=vX.Y.Z`
6. Publicar el tag:
   - `make -C <modulo> release-push-tag VERSION=vX.Y.Z`
7. Crear el GitHub Release:
   - `make -C <modulo> release-github VERSION=vX.Y.Z`

## Ejemplos

### `postgres`

```bash
make -C postgres release-check
make -C postgres release-prepare VERSION=v0.62.0
make -C postgres release-tag VERSION=v0.62.0
make -C postgres release-push-tag VERSION=v0.62.0
make -C postgres release-github VERSION=v0.62.0
```

### `docker`

```bash
make -C docker release-check
make -C docker release-prepare VERSION=v0.2.0
make -C docker release-tag VERSION=v0.2.0
make -C docker release-push-tag VERSION=v0.2.0
make -C docker release-github VERSION=v0.2.0
```

## Workflow remoto

Cuando llega un tag con formato `modulo/vX.Y.Z`, `.github/workflows/release.yml`:

1. resuelve el modulo y la version;
2. ejecuta `make -C <modulo> release-check`;
3. extrae notas desde el changelog del modulo;
4. crea el GitHub Release con el mismo tag.

En modulos Go agrega ademas un bloque de instalacion con `go get github.com/EduGoGroup/edugo-infrastructure/<modulo>@<version>`.

## Alcance por modulo

- `postgres`: valida build, tests, lint, vet, fmt-check y changelog.
- `mongodb`: valida build, tests, lint, vet, fmt-check y changelog.
- `schemas`: valida build, tests, lint, vet, fmt-check y changelog.
- `tools/mock-generator`: valida build, tests, lint, vet, fmt-check y changelog.
- `docker`: valida `docker compose config -q` y changelog.
