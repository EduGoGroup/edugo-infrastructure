# AGENTS.md — edugo-infrastructure

> Detalle local. Reglas globales del ecosistema en `../../AGENTS.md` (no las repitas).
> Norte actual del proyecto en `docs/ACTIVE.md`. Mapa de módulos y arquitectura en `README.md`
> y `docs/` (architecture, repository-map, processes, releasing, automation, roadmap).

## Propósito

**Base de infraestructura y contratos compartidos** de EduGo: el esquema relacional Postgres (schemas,
migraciones, entities, seeds), los **JSON Schema** que validan los eventos del ecosistema, la topología
local de Docker Compose, y herramientas de generación de datos. No es un servicio de runtime: es el
sustrato sobre el que corren las APIs y el worker.

## Módulos

| Módulo | Rol |
| --- | --- |
| `postgres` | Schema relacional (`auth`, `iam`, `academic`, `assessment`, `content`, `ui_config`, `notifications`, `audit`, ...), `migrations/`, `entities/`, `seeds/` (`system/` = contrato + `playground_v2/` = datos, con `base` por defecto; ver MP-09), runner. Tiene su propio `go.mod`. |
| `schemas` | Validación de contratos vía JSON Schema para los eventos (`events/`); empareja con `edugo-shared/messaging/events`. |
| `tools/mock-generator` | Generación de datasets Go desde SQL. |
| `docker` | Topología local con Docker Compose. |

## Cómo usar

`Makefile` raíz + un `make` por módulo. Patrón observado: `make -C <modulo> release-check` valida cada
superficie (postgres, schemas, tools/mock-generator, docker).
Scripts en `scripts/`: `dev-setup.sh` (ambiente local), `seed-data.sh` (aplicar seeds),
`module-release.sh` / `auto-release.sh` (release por módulo).

## Convenciones y gotchas locales

- **PROHIBIDO tocar seeds y migraciones sin confirmar** (regla global): `postgres/seeds/` y
  `postgres/migrations/` son la fuente de verdad de datos del ecosistema; un cambio aquí afecta a todas
  las APIs y al worker.
- **Un solo mundo de seeds (MP-09, 2026-06-14):** `system/` (L0–L4) = **contrato puro**; `playground_v2/`
  = **único mundo de datos**, con `base` como fixture **por defecto**. `make docker-recreate` sin flags →
  `system` + `base` (2 escuelas, 9 usuarios `@edugo.test`, login `12345678`). Fixtures focalizados:
  `make docker-playground-v2 P=<fixture>`. `seeds/demo/` y `seeds/playground/` (v1) **fueron eliminados**.
- **Schemas separados, una sola BD**: Postgres está particionado por schemas de dominio; cada API declara
  su `search_path` (p.ej. identity `auth,iam,academic,audit,public`; learning `content,assessment,auth,public`).
  Las entities y migraciones aquí definen esos schemas.
- **Permisos SDUI** se siembran vía `slot.permission` en los seeds (única fuente; ver `../../docs/adr/0003`).
- **Contratos de eventos**: si cambias un schema en `schemas/events/`, sincroniza con
  `edugo-shared/messaging/events` y con los procesadores del worker.
- Release por módulo: tags modulares y `release.yml`; cada módulo lleva su `CHANGELOG.md`.
- Reglas globales: código en inglés, logs/docs en español, fechas UTC en BD y transporte.

## Versionado de dependencias `edugo-shared`: migrar hacia adelante, nunca re-taguear

> **Por qué:** los módulos Go de aquí (`postgres`, `migrator`, `schemas`, `tools/mock-generator`)
> consumen `edugo-shared`. Si un módulo de shared se **re-taguea** (mismo número, contenido distinto —
> p. ej. tras un "clean reset"), el build **local pasa** (el `go.work` lo resuelve del filesystem sin
> checksum) pero **CI/cloud falla**: Go descarga el tag de GitHub, lo hashea y lo compara contra el
> `go.sum` → `SECURITY ERROR: checksum mismatch`. El "historial" vive en el `go.sum` commiteado (repos
> privados: no hay sumdb global).

**Regla:** el estándar de versión de los módulos shared es **`0.900.X`** (→ `1.0.0` en producción).
- **Nunca re-pushees un tag existente.** Si cambia el contenido, **sube el número** (tag nuevo e inmutable).
- **Si vas a modificar (o tu cambio arrastra) un módulo shared aún en `v0.1.0`, migrálo a `0.900.X`**: bump
  del tag + bump del `require` en **todos** los consumidores + `go mod tidy`. No lo dejes "bajo estándar".
- Tras migrar, **retira el `replace` transitorio del `go.work`** (es muleta local, no fix). El objetivo es
  que CI/cloud resuelva contra el tag real.

> Detalle y el caso `audit`/`audit/postgres` que originó la regla: `../edugo-shared/AGENTS.md`
> (sección "Versionado de módulos").
