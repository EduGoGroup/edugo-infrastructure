# postgres

Modulo responsable del activo relacional de `edugo-infrastructure`.

## Rol actual

- Define el schema SQL principal del dominio.
- Embebe scripts de estructura para inicializar una base desde cero.
- Embebe seeds de produccion y desarrollo.
- Publica entities Go alineadas con tablas y schemas actuales.
- Expone CLIs para migracion legacy y ejecucion por capas.

## Estado observado

- Tests cortos pasan.
- El paquete `migrations` y el paquete `seeds` son la superficie programatica mas confiable del modulo hoy.
- La versión actual del módulo es `postgres/v0.900.5` (plan 019: `entity-picker` en `assessments-form.subject_id`; `view_when` read-only fuera de borrador ADR 0022; baja de forms SDUI legacy `grades-form`/`user-roles`; fixture de conformidad del contrato entity-picker; `L4_SEED_VERSION` 1.54.0, `SchemaVersion` 3.58.0). Detalle en [CHANGELOG.md](CHANGELOG.md).

## Documentacion del modulo

- [docs/README.md](docs/README.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [../docs/releasing.md](../docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

## Entrada rapida

Comandos utiles observados:

- `make -C postgres build test fmt-check`
- `make -C postgres release-check`
- `make -C postgres migrate-status`
- `make -C postgres runner-up`
- `make -C postgres seed-all`

## Nota operativa

La integracion con otros repositorios ya esta documentada en `docs/ecosystem-integration.md`. El flujo de release se comparte a nivel repo en `../docs/releasing.md` para no duplicar reglas.
