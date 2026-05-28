# AGENTS.md — edugo-infrastructure

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Norte actual del proyecto en `docs/ACTIVE.md`. Mapa de módulos y arquitectura en `README.md`
> y `docs/` (architecture, repository-map, processes, releasing, automation, roadmap).

## Propósito

**Base de infraestructura y contratos compartidos** de EduGo: el esquema relacional Postgres (schemas,
migraciones, entities, seeds), las collections de MongoDB, los **JSON Schema** que validan los eventos
del ecosistema, la topología local de Docker Compose, y herramientas de generación de datos. No es un
servicio de runtime: es el sustrato sobre el que corren las APIs y el worker.

## Módulos

| Módulo | Rol |
| --- | --- |
| `postgres` | Schema relacional (`auth`, `iam`, `academic`, `assessment`, `content`, `ui_config`, `notifications`, `audit`, ...), `migrations/`, `entities/`, `seeds/` (system/demo/e2e/playground), runner. Tiene su propio `go.mod`. |
| `mongodb` | Collections documentales y fixtures que consume el worker. |
| `schemas` | Validación de contratos vía JSON Schema para los eventos (`events/`); empareja con `edugo-shared/messaging/events`. |
| `tools/mock-generator` | Generación de datasets Go desde SQL. |
| `docker` | Topología local con Docker Compose. |

## Cómo usar

`Makefile` raíz + un `make` por módulo. Patrón observado: `make -C <modulo> release-check` valida cada
superficie (postgres, mongodb, schemas, tools/mock-generator, docker).
Scripts en `scripts/`: `dev-setup.sh` (ambiente local), `seed-data.sh` (aplicar seeds),
`module-release.sh` / `auto-release.sh` (release por módulo).

## Convenciones y gotchas locales

- **PROHIBIDO tocar seeds y migraciones sin confirmar** (regla global): `postgres/seeds/` y
  `postgres/migrations/` son la fuente de verdad de datos del ecosistema; un cambio aquí afecta a todas
  las APIs y al worker.
- **Schemas separados, una sola BD**: Postgres está particionado por schemas de dominio; cada API declara
  su `search_path` (p.ej. identity `auth,iam,academic,audit,public`; learning `content,assessment,auth,public`).
  Las entities y migraciones aquí definen esos schemas.
- **Permisos SDUI** se siembran vía `slot.permission` en los seeds (única fuente; ver `../../docs/adr/0003`).
- **Contratos de eventos**: si cambias un schema en `schemas/events/`, sincroniza con
  `edugo-shared/messaging/events` y con los procesadores del worker.
- Release por módulo: tags modulares y `release.yml`; cada módulo lleva su `CHANGELOG.md`.
- Reglas globales: código en inglés, logs/docs en español, fechas UTC en BD y transporte.
