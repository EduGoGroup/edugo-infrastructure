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
- [docs/phase-1-scope.md](docs/phase-1-scope.md)
- [docs/repository-map.md](docs/repository-map.md)
- [docs/processes.md](docs/processes.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/ecosystem-integration.md](docs/ecosystem-integration.md)
- [docs/automation.md](docs/automation.md)
- [docs/releasing.md](docs/releasing.md)
- [docs/roadmap.md](docs/roadmap.md)
- [CHANGELOG.md](CHANGELOG.md)

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
