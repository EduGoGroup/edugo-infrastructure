# mongodb ecosystem integration

## Rol ecosistemico

`mongodb` soporta el lado documental del ecosistema, especialmente para artefactos derivados del procesamiento asincrono y de evaluaciones.

## Consumidores observados

### `edugo-api-mobile-new`

Consumo observado:

- dependencia en `go.mod` a `github.com/EduGoGroup/edugo-infrastructure/mongodb v0.53.0`
- uso de `mongodb/entities` en servicios de assessment, summary, scoring y repositorios Mongo

Superficies usadas:

- `MaterialAssessment`
- `MaterialSummary`

### `edugo-worker`

Consumo observado:

- dependencia en `go.mod` a `mongodb v0.53.0`
- uso de `mongodb/entities` en processors, servicios de dominio y repositorios
- uso de `mongodb/migrations` en tests de integracion

Superficies usadas:

- `MaterialAssessment`
- `MaterialSummary`
- `MaterialEvent`
- setup de estructura para tests

### `edugo-dev-environment/migrator`

Consumo observado:

- importa `mongodb/migrations`
- ejecuta `ApplyAll(ctx, db)` como parte de la recreacion de entorno

## Integracion con reglas del ecosistema

En el ecosistema, los cambios documentales siguen el mismo principio que Postgres: el cambio se hace aqui y luego se propaga via `go.work` local o release publicada.

En la practica actual, el cambio ecosistemico pasa por:

1. modificar `mongodb/migrations` o `mongodb/entities`
2. validar consumidores `mobile` y `worker`
3. recrear Mongo con `edugo-dev-environment/migrator`
4. publicar release si otros repos no trabajan en local con `go.work`

## Integracion interna con otros modulos del repo

### Con `postgres`

La integracion interna se apoya en IDs y referencias cruzadas:

- `material_id` compartido
- `mongo_document_id` expuesto en el lado relacional de assessments
- seeds de Mongo alineados con materiales de Postgres

### Con `schemas`

La relacion es de proceso:

- eventos del ecosistema disparan o representan trabajo que termina reflejado en documentos Mongo
- `assessment.generated` referencia explicitamente un `mongo_document_id`

## Hallazgo importante

`ecosistema.md` lista varias colecciones historicas, pero la superficie activa de este modulo hoy esta concentrada en tres collections creadas por `mongodb/migrations/embed.go`:

- `material_summary`
- `material_assessment_worker`
- `material_event`

Para integracion de fase 2, eso importa mas que el inventario heredado.
