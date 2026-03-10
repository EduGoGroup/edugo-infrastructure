# mongodb

Modulo responsable de las collections MongoDB y de los documentos derivados usados por el worker.

## Rol actual

- Define collections y validators via codigo Go.
- Define indices via codigo Go.
- Embebe seeds canonicos y mock data.
- Publica entities Go para las collections activas.
- Expone un CLI de estado de migraciones, no un runner completo de ejecucion.

## Estado observado

- Tests cortos pasan.
- El ultimo tag observado en Git es `mongodb/v0.53.0`.
- La superficie publica estable del modulo hoy es el paquete `mongodb/migrations`.

## Documentacion del modulo

- [docs/README.md](docs/README.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [../docs/releasing.md](../docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

## Entrada rapida

Comandos utiles observados:

- `make -C mongodb build test fmt-check`
- `make -C mongodb release-check`
- `make -C mongodb migrate-status`
- `make -C mongodb runner-up`
- `make -C mongodb seed-all`

## Nota operativa

La integracion con otros modulos y repositorios ya esta documentada en `docs/ecosystem-integration.md`. El flujo de release se comparte en `../docs/releasing.md`.
