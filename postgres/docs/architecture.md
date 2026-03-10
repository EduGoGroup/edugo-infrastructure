# postgres architecture

## Mapa interno

```text
postgres/
|-- cmd/
|   |-- migrate/
|   `-- runner/
|-- docs/
|-- entities/
|-- migrations/
|   `-- structure/
|-- seeds/
|   |-- development/
|   `-- production/
|-- Makefile
|-- README.md
`-- CHANGELOG.md
```

## Activos principales

| Activo | Funcion |
| --- | --- |
| `migrations/structure/*.sql` | define schemas, tablas, funciones, vistas y FK |
| `migrations/embed.go` | embebe y ejecuta SQL de estructura |
| `seeds/embed.go` | embebe y ejecuta seeds de produccion y desarrollo |
| `entities/*.go` | representa tablas como structs Go |
| `cmd/migrate` | administra migraciones versionadas legacy |
| `cmd/runner` | intenta ejecutar capas desde el filesystem |

## Diagrama local

```mermaid
flowchart TB
    PG[postgres]
    PG --> SQL[structure SQL]
    PG --> SD[seed packages]
    PG --> ENT[entities]
    PG --> CLI1[cmd/migrate]
    PG --> CLI2[cmd/runner]

    SQL --> AUTH[auth]
    SQL --> IAM[iam]
    SQL --> ACD[academic]
    SQL --> CNT[content]
    SQL --> ASM[assessment]
    SQL --> UIC[ui_config]
    SQL --> AUD[audit]

    SD --> PROD[production seeds]
    SD --> DEV[development seeds]
```

## Decisiones estructurales visibles

- El schema se ordena por prefijos numericos y por domain.
- Los seeds se separan entre canonicos y de desarrollo.
- Los entities viven como reflejo tipado del schema.
- Las capacidades programaticas embebidas son mas confiables que algunos scripts heredados de shell o Makefile.
