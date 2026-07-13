# Repository Map

## Vista rapida

Este repositorio contiene cuatro superficies documentadas y un conjunto transversal de automatizacion.

| Superficie | Tipo | Estado observado | Fuente principal |
| --- | --- | --- | --- |
| `postgres/` | Modulo Go + SQL | Operativo en tests cortos | SQL embebido, seeds, entities, CLIs |
| `schemas/` | Modulo Go | Operativo en tests cortos | JSON Schemas embebidos + validator |
| `tools/mock-generator/` | Modulo Go CLI | Compila via tests | parser SQL + generador de dataset |
| `docker/` | Soporte de runtime local | Configurado | `docker-compose.yml` y `.env.example` |
| `.github/`, `scripts/`, `Makefile` | Automatizacion transversal | Parcialmente alineada | workflows, actions, scripts, comandos raiz |

## Conteos observados

- `postgres/migrations/structure`: 33 scripts SQL
- `postgres/seeds/system/layers` + `postgres/seeds/system/l4`: capas L0..L4 del rebuild (post-Fase-6)
- `postgres/seeds/demo`: seed de desarrollo (ex `development/`)
- `postgres/entities`: 27 structs Go
- `schemas/events`: 4 contratos JSON Schema
- `.github/workflows`: 9 workflows
- `.github/actions`: 3 composite actions

## Mapa de carpetas relevante

```text
.
|-- Makefile
|-- .env.example
|-- .golangci.yml
|-- docs/
|-- docker/
|-- postgres/
|-- schemas/
|-- scripts/
|-- tools/mock-generator/
`-- .github/
```

## Modulos y lectura sugerida

- `postgres`: ver [../postgres/docs/processes.md](../postgres/docs/processes.md)
- `schemas`: ver [../schemas/docs/processes.md](../schemas/docs/processes.md)
- `tools/mock-generator`: ver [../tools/mock-generator/docs/processes.md](../tools/mock-generator/docs/processes.md)
- `docker`: ver [../docker/docs/processes.md](../docker/docs/processes.md)
