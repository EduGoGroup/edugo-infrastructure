# postgres ecosystem integration

## Rol ecosistemico

`postgres` es el modulo con mayor radio de impacto del repo. Define el modelo relacional que consumen los servicios backend y el migrador del ambiente.

## Consumidores observados

### `edugo-api-iam-platform`

Consumo observado:

- dependencia en `go.mod` a `github.com/EduGoGroup/edugo-infrastructure/postgres v0.61.0`
- uso intensivo de `postgres/entities`

Superficies usadas:

- roles
- permisos
- recursos
- user roles
- school concepts
- screen config

### `edugo-api-admin-new`

Consumo observado:

- dependencia en `go.mod` a `postgres v0.61.0`
- uso de `postgres/entities` en DTOs, servicios, repositorios y tests

Superficies usadas:

- schools
- academic units
- memberships
- subjects
- guardian relations
- concept types
- materiales relacionales auxiliares

### `edugo-api-mobile-new`

Consumo observado:

- dependencia en `go.mod` a `postgres v0.61.0`
- uso de `postgres/entities` en servicios, repositorios, handlers y tests

Superficies usadas:

- materials
- assessment
- assessment attempts
- progress
- guardians
- screen config

### `edugo-dev-environment/migrator`

Consumo observado:

- importa `postgres/migrations`
- importa `postgres/seeds`

Ese repo es hoy el ejecutor canonico del reseteo y recreacion de base de datos del ecosistema.

## Integracion con reglas del ecosistema

Segun `ecosistema.md`, un cambio de BD debe hacerse aqui editando el `CREATE` completo. En la practica actual, eso implica tocar:

- `postgres/migrations/structure/*.sql`
- `postgres/entities/*.go` cuando cambian los tipos compartidos
- eventualmente `postgres/seeds/*` si el cambio altera datos canonicos o de desarrollo

Luego el ecosistema espera:

1. validacion local por `go.work`
2. recreacion de BD mediante `edugo-dev-environment/migrator`
3. release del modulo `postgres` cuando el cambio debe ser consumido por otros repos

## Integracion interna con otros modulos del repo

### Con `mongodb`

- `assessment.assessment` contiene `mongo_document_id`
- seeds de desarrollo de assessments y materiales se alinean con documentos Mongo

### Con `schemas`

Los contratos de eventos usan IDs que nacen del modelo relacional: usuarios, escuelas, memberships, materiales.

### Con `tools/mock-generator`

El generador importa `postgres/entities`, por lo que cualquier cambio estructural en `entities` tiene impacto potencial en la herramienta.

## Riesgos de integracion

- cambios en `entities` rompen compilacion de IAM, Admin y Mobile casi de inmediato
- cambios en seeds de produccion impactan bootstrap de permisos, recursos y UI dinamica
- cambios en `structure/*.sql` deben seguir siendo compatibles con el migrador y con la regla ecosistemica de recreacion completa
