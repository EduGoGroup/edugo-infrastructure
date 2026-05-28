# schemas ecosystem integration

## Rol ecosistemico

`schemas` es la capa potencial de contratos versionados para eventos del ecosistema.

## Consumo observado

En los repos escaneados del ecosistema no aparecio un consumo directo de `github.com/EduGoGroup/edugo-infrastructure/schemas` en APIs, worker ni migrator.

Repos revisados:

- `edugo-api-iam-platform`
- `edugo-api-admin-new`
- `edugo-api-mobile-new`
- `edugo-worker`
- `edugo-dev-environment`
- `kmp_new`
- `apple_new`

## Interpretacion

Eso sugiere una de dos situaciones actuales:

1. el modulo existe como capa de contrato disponible, pero todavia no esta adoptado de forma directa por los servicios
2. parte de la semantica de eventos vive hoy en `edugo-shared/messaging/events` y no en este modulo

## Integracion semantica con el ecosistema

Aunque no haya consumo directo observado, los contratos presentes en `schemas/events` coinciden con procesos del ecosistema:

- `material.uploaded`
- `assessment.generated`
- `material.deleted`
- `student.enrolled`

Esos nombres encajan con el dominio descrito en `ecosistema.md` y con servicios que manipulan materiales, evaluaciones, miembros y procesamiento asincrono.

## Integracion interna con otros modulos del repo

### Con `postgres`

Los payloads usan IDs y conceptos que nacen del modelo relacional.

### Con `mongodb`

`assessment.generated` incluye `mongo_document_id`, puente hacia el modelo documental.

## Implicancia para fase 3

Si el ecosistema decide adoptar este modulo de forma directa, habra que formalizar:

- quien valida al publicar
- quien valida al consumir
- como se versionan cambios compatibles e incompatibles
- como se sincroniza con `edugo-shared/messaging/events`
