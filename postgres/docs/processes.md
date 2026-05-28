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

### 3. Seeds del sistema

El paquete `postgres/seeds/system/` provee la capa programatica del seed system, organizada en capas que cumplen la interfaz `Layer` (`Name()`, `SeedVersion()`, `Apply()`).

- `system/layer.go` — interfaz `Layer`.
- `system/system.go` — `Layers()` y `ApplySystem(db, upTo)`.
- `system/layers/` — capas L0..L4 del rebuild (Fase 6 cerrada).
- `system/l4/` — sub-paquete con los datos de L4 por dominio + accessors publicos consumidos por el cross-checker.

Cobertura del seed system (L0..L4):

- recursos
- roles
- permisos
- role-permissions
- screen templates
- screen instances
- resource-screens
- concept types
- concept definitions

Estos datos conforman la configuracion canonica del sistema y se aplican siempre.

### 4. Seeds de demo

El modulo incluye el paquete `postgres/seeds/demo/` (ex `seeds/development/`), que expone `ApplyDemo(gdb)`.

Procesos cubiertos por esos seeds:

- truncado controlado del dataset de demo
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
- progreso de estudiantes

### 5. Exposicion de entities Go

El directorio `postgres/entities/` contiene 27 structs, entre ellos:

- autenticacion: `User`, `RefreshToken`, `LoginAttempt`
- IAM: `Resource`, `Role`, `Permission`, `RolePermission`, `UserRole`
- academico: `School`, `AcademicUnit`, `Membership`, `Subject`, `GuardianRelation`, `ConceptType`, `ConceptDefinition`, `SchoolConcept`, `Unit`
- contenido: `Material`, `MaterialVersion`, `Progress`
- evaluacion: `Assessment`, `AssessmentAttempt`, `AssessmentAttemptAnswer`, `AssessmentMaterial`
- UI config: `ScreenTemplate`, `ScreenInstance`, `ResourceScreen`

### 6. Operacion por CLI

El modulo tiene dos entrypoints operativos:

- `cmd/runner/runner.go`: ejecuta estructura + seeds; aplica migraciones DDL, luego `system.ApplySystem(db, "")` y opcionalmente `demo.ApplyDemo(gdb)`.
- `cmd/seed/main.go`: ejecuta solo seeds; expone `system.ApplySystem` y `demo.ApplyDemo` segun flag.

Consultar `--help` de cada binario para los flags vigentes.

### 7. Tests del modulo

El modulo tiene tests de integracion en `postgres/migrations/migrations_integration_test.go` y los tests cortos pasan en el estado observado de esta fase.

## Superficie recomendada

Si alguien necesita usar el modulo desde el propio repo, la superficie mas confiable hoy es:

- `postgres/migrations` para estructura
- `postgres/seeds` para datasets embebidos
- `postgres/entities` para tipos de datos
