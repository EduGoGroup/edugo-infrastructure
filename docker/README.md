# docker

Superficie de runtime local para `edugo-infrastructure`.

## Rol actual

- Define servicios core para desarrollo local.
- Expone perfiles opcionales para mensajeria, cache y herramientas visuales.
- Usa `.env.example` como contrato base de configuracion local.

## Estado observado

- El ultimo tag observado en Git es `docker/v0.1.0`.
- El modulo no es un modulo Go, pero si una pieza operativa central del repo.

## Documentacion del modulo

- [docs/README.md](docs/README.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [../docs/releasing.md](../docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

## Entrada rapida

- `make -C docker validate`
- `make -C docker up-core`
- `make -C docker up-messaging`
- `make -C docker up-full`
- `make -C docker release-check`

## Nota operativa

La integracion con `edugo-dev-environment` y el ecosistema ya esta documentada en `docs/ecosystem-integration.md`. El flujo de versionado se comparte en `../docs/releasing.md`.
