# schemas

Modulo responsable de validar contratos JSON Schema para eventos.

## Rol actual

- Embebe archivos JSON Schema desde `schemas/events`.
- Construye un indice interno `event_type:event_version`.
- Valida eventos desde struct, bytes JSON o tipo/version explicitos.
- Mantiene tests de validacion y benchmarks.

## Estado observado

- Tests cortos pasan.
- El ultimo tag observado en Git es `schemas/v0.51.0`.
- La API publica esta concentrada en `validator.go`.

## Documentacion del modulo

- [docs/README.md](docs/README.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [../docs/releasing.md](../docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

## Entrada rapida

- `make -C schemas build test fmt-check`
- `make -C schemas release-check`

## Nota operativa

La integracion con servicios externos ya esta en `docs/ecosystem-integration.md`. El flujo de release se documenta una sola vez en `../docs/releasing.md`.
