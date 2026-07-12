# EduGo Infrastructure

Base de infraestructura y contratos compartidos de EduGo.

La base documental vigente ya cubre tres capas:

- Fase 1: estado real del repo visto de forma autocontenida.
- Fase 2: integracion entre modulos y con el ecosistema EduGo.
- Fase 3: operacion uniforme por modulo para validacion y release.

## Alcance actual

- Documentar el estado actual del repo desde su propio codigo, scripts y automatizacion.
- Documentar integraciones reales con el ecosistema a partir de codigo y repos relacionados.
- Estandarizar build, test, lint, fmt, vet y release-check por modulo.
- Centralizar el acceso a la documentacion por modulo.
- Reemplazar documentacion heredada o desalineada.

## Modulos documentados

| Modulo | Rol actual | Documentacion |
| --- | --- | --- |
| `postgres` | Schema relacional, seeds y entities | [postgres/README.md](postgres/README.md) |
| `mongodb` | Collections documentales y fixtures del worker | [mongodb/README.md](mongodb/README.md) |
| `schemas` | Validacion de contratos JSON Schema para eventos | [schemas/README.md](schemas/README.md) |
| `tools/mock-generator` | Generacion de datasets Go desde SQL | [tools/mock-generator/README.md](tools/mock-generator/README.md) |
| `docker` | Topologia local con Docker Compose | [docker/README.md](docker/README.md) |

## Documentacion general

- [docs/README.md](docs/README.md)
- [docs/repository-map.md](docs/repository-map.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/automation.md](docs/automation.md)
- [docs/releasing.md](docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

> Histórico: `docs/archive/` contiene `phase-1-scope.md`, `ecosystem-integration.md` y `roadmap.md`, de
> marzo (nombres de API viejos y carpeta `Common/`, ya inexistente). No son fuente de verdad.

## Seeds (MP-09)

`postgres/seeds/` tiene dos planos: `system/` = contrato puro y `playground_v2/` = único mundo de datos,
con `base` como fixture por defecto (2 escuelas, 9 usuarios `@edugo.test`, login `12345678`). Detalle en
[`AGENTS.md`](AGENTS.md).

Los comandos que **aplican** migraciones y seeds no viven aquí: están en el `Makefile` de
`../edugo-dev-environment/migrator/` — `make docker-recreate` (system + base) y
`make docker-playground-v2 P=<fixture>` (fixtures focalizados). Este repo aporta los schemas, entities,
migraciones y datos; el migrator los ejecuta.

## Scripts de automatizacion

- **[scripts/auto-release.sh](scripts/auto-release.sh)** - Release automatizado de modulos ([Documentacion](scripts/AUTO-RELEASE-README.md))
- [scripts/module-release.sh](scripts/module-release.sh) - Release manual de modulos
- [scripts/dev-setup.sh](scripts/dev-setup.sh) - Setup del ambiente local
- [scripts/seed-data.sh](scripts/seed-data.sh) - Aplicar seeds a las bases de datos

## Estado operativo observado

Validaciones ejecutadas localmente el 2026-03-08 sobre este repo:

- `make -C postgres release-check`
- `make -C mongodb release-check`
- `make -C schemas release-check`
- `make -C tools/mock-generator release-check`
- `make -C docker release-check`

Esas superficies pasan en local. El flujo de release por modulo tambien queda documentado y soportado por `scripts/module-release.sh`, `make/module-release.mk` y `.github/workflows/release.yml`.

## Principios de esta nueva documentacion

- El codigo y la estructura actual mandan.
- Los procesos tienen prioridad sobre la descripcion promocional.
- La arquitectura se explica despues del proceso.
- Las integraciones externas viven en documentos especificos, no mezcladas con la descripcion base del modulo.
- Las repeticiones se evitan enlazando a la documentacion por modulo.
