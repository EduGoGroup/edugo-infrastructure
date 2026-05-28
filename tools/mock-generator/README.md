# tools/mock-generator

CLI para generar codigo Go de dataset a partir de scripts SQL con `INSERT`.

## Rol actual

- Parsea directorios de SQL.
- Extrae filas de `INSERT`.
- Genera archivos Go para un paquete `dataset`.
- Intenta formatear el codigo generado.

## Estado observado

- El modulo compila via tests.
- No tiene tests propios mas alla de compilacion de paquetes.
- El ultimo tag observado en Git es `tools/mock-generator/v0.51.0`.

## Documentacion del modulo

- [docs/README.md](docs/README.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [../../docs/releasing.md](../../docs/releasing.md)
- [CHANGELOG.md](CHANGELOG.md)

## Entrada rapida

- `make -C tools/mock-generator build test fmt-check`
- `make -C tools/mock-generator release-check`
- `make -C tools/mock-generator help`

## Nota operativa

La integracion con el ecosistema ya esta documentada en `docs/ecosystem-integration.md`. El release sigue las mismas reglas compartidas del repo en `../../docs/releasing.md`.
