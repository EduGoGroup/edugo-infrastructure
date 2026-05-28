# Architecture

## Vista de repositorio

El repositorio se organiza como una plataforma de infraestructura compartida. Cada modulo encapsula un tipo de activo distinto: SQL, documentos Mongo, contratos JSON, generacion de dataset o runtime local.

```mermaid
flowchart TB
    R[edugo-infrastructure]

    R --> PG[postgres]
    R --> MG[mongodb]
    R --> SC[schemas]
    R --> MK[tools/mock-generator]
    R --> DK[docker]
    R --> AU[automation]

    PG --> PGA[SQL structure]
    PG --> PGB[production seeds]
    PG --> PGC[development seeds]
    PG --> PGD[entities and CLIs]

    MG --> MGA[Go migrations]
    MG --> MGB[seeds and mock data]
    MG --> MGC[entities]

    SC --> SCA[JSON Schemas]
    SC --> SCB[validator package]

    MK --> MKA[SQL parser]
    MK --> MKB[dataset generator]

    DK --> DKA[docker compose profiles]
    DK --> DKB[local services]

    AU --> AUA[root Makefile]
    AU --> AUB[scripts]
    AU --> AUC[GitHub workflows]
    AU --> AUD[composite actions]
```

## Capas observadas

| Capa | Superficies |
| --- | --- |
| Datos relacionales | `postgres` |
| Datos documentales | `mongodb` |
| Contratos | `schemas` |
| Runtime local | `docker` |
| Generacion auxiliar | `tools/mock-generator` |
| Calidad y entrega | `Makefile`, `scripts`, `.github` |

## Regla de modelado usada en esta documentacion

- Primero se documenta el proceso interno del modulo.
- Luego se describe su arquitectura local.
- La integracion entre modulos queda diferida para fase 2.
