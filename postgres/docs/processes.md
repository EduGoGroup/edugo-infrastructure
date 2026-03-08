# postgres processes

## Procesos propios del modulo

### 1. Definicion del schema canonico

El modulo define un schema relacional segmentado por domains SQL:

- `auth`
- `iam`
- `academic`
- `content`
- `assessment`
- `ui_config`
- `audit`

La estructura vive en `postgres/migrations/structure/*.sql` y hoy suma 33 scripts ordenados por prefijo numerico.

### 2. Inicializacion de base desde SQL embebido

El paquete `postgres/migrations` embebe `structure/*.sql` y expone:

- `ApplyAll(db)`
- `GetScript(name)`
- `ListScripts()`
- `GetScriptsByLayer(layer)`

En el estado actual, la capa programatica publica y estable es `structure`.

### 3. Seeds de produccion

El paquete `postgres/seeds` embebe `production/*.sql` y `development/*.sql`.

Seeds de produccion observados:

- `001_resources.sql`
- `002_roles.sql`
- `003_permissions.sql`
- `004_role_permissions.sql`
- `005_ui_screen_templates.sql`
- `006_ui_screen_instances.sql`
- `007_ui_resource_screens.sql`
- `008_concept_types.sql`

Estos scripts conforman la configuracion canonica del sistema y no deben tratarse como datos de demo.

### 4. Seeds de desarrollo

El modulo incluye 13 scripts en `postgres/seeds/development/`.

Procesos cubiertos por esos seeds:

- truncado controlado del dataset de desarrollo
- creacion de escuelas
- creacion de unidades academicas jerarquicas
- usuarios y memberships
- asignacion RBAC
- materias
- materiales
- assessments
- intentos y respuestas
- relaciones guardian-estudiante
- preferencias UI
- conceptos por escuela

### 5. Exposicion de entities Go

El directorio `postgres/entities/` contiene 27 structs, entre ellos:

- autenticacion: `User`, `RefreshToken`, `LoginAttempt`
- IAM: `Resource`, `Role`, `Permission`, `RolePermission`, `UserRole`
- academico: `School`, `AcademicUnit`, `Membership`, `Subject`, `GuardianRelation`, `ConceptType`, `ConceptDefinition`, `SchoolConcept`
- contenido: `Material`, `MaterialVersion`, `Progress`
- evaluacion: `Assessment`, `AssessmentAttempt`, `AssessmentAttemptAnswer`, `AssessmentMaterial`
- UI config: `ScreenTemplate`, `ScreenInstance`, `ResourceScreen`

### 6. Operacion por CLI

El modulo tiene dos entrypoints operativos:

- `cmd/migrate/migrate.go`: CLI legacy con `up`, `down`, `status`, `create`, `force`
- `cmd/runner/runner.go`: runner por capas que intenta ejecutar `structure`, `constraints`, `seeds`, `testing`

### 7. Tests del modulo

El modulo tiene tests de integracion en `postgres/migrations/migrations_integration_test.go` y los tests cortos pasan en el estado observado de esta fase.

## Realidades que importan documentar

- El paquete embebido y probado es `postgres/migrations` sobre `structure/*.sql`.
- La narrativa vieja de una arquitectura de 4 capas ya no representa exactamente el layout actual del modulo.
- `postgres/Makefile` conserva un target `seed` con una referencia a `../scripts/get-db-url.go`, archivo que hoy no existe.

## Superficie recomendada en fase 1

Si alguien necesita usar el modulo desde el propio repo, la superficie mas confiable hoy es:

- `postgres/migrations` para estructura
- `postgres/seeds` para datasets embebidos
- `postgres/entities` para tipos de datos
