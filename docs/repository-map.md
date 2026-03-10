# Repository Map

## Vista rapida

Este repositorio contiene cinco superficies documentadas y un conjunto transversal de automatizacion.

| Superficie | Tipo | Estado observado | Fuente principal |
| --- | --- | --- | --- |
| `postgres/` | Modulo Go + SQL | Operativo en tests cortos | SQL embebido, seeds, entities, CLIs |
| `mongodb/` | Modulo Go + Mongo | Operativo en tests cortos | migraciones Go, seeds, mock data, entities |
| `schemas/` | Modulo Go | Operativo en tests cortos | JSON Schemas embebidos + validator |
| `tools/mock-generator/` | Modulo Go CLI | Compila via tests | parser SQL + generador de dataset |
| `docker/` | Soporte de runtime local | Configurado | `docker-compose.yml` y `.env.example` |
| `.github/`, `scripts/`, `Makefile` | Automatizacion transversal | Parcialmente alineada | workflows, actions, scripts, comandos raiz |

## Conteos observados

- `postgres/migrations/structure`: 33 scripts SQL
- `postgres/seeds/production`: 8 scripts SQL
- `postgres/seeds/development`: 13 scripts SQL
- `postgres/entities`: 27 structs Go
- `mongodb/entities`: 3 structs Go
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
|-- mongodb/
|-- postgres/
|-- schemas/
|-- scripts/
|-- tools/mock-generator/
`-- .github/
```

## Modulos y lectura sugerida

- `postgres`: ver [../postgres/docs/processes.md](../postgres/docs/processes.md)
- `mongodb`: ver [../mongodb/docs/processes.md](../mongodb/docs/processes.md)
- `schemas`: ver [../schemas/docs/processes.md](../schemas/docs/processes.md)
- `tools/mock-generator`: ver [../tools/mock-generator/docs/processes.md](../tools/mock-generator/docs/processes.md)
- `docker`: ver [../docker/docs/processes.md](../docker/docs/processes.md)
