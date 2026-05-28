# Reglas para modificar migraciones

## OBLIGATORIO al cambiar cualquier archivo aqui

1. **Incrementar `SchemaVersion`** en `version.go` (mismo directorio). Sin excepcion.
2. **NUNCA usar ALTER TABLE.** Modifica la struct entity correspondiente en `../entities/*.go` (GORM AutoMigrate aplica el delta).
3. **NUNCA crear archivos de migracion incremental** (ej: `090_alter_users_add_column.sql`). Edita el entity o el archivo SQL maestro.
4. **Pre/post SQL**: lo que GORM no puede expresar (CHECK constraints con regex, indices parciales con WHERE, funciones plpgsql, vistas, triggers, FKs cross-schema) va en `sql/post_gorm.sql`. Schemas, extensiones, ENUMs y funciones compartidas van en `sql/pre_gorm.sql`.

## Arquitectura de migraciones

El pipeline `migrate.go` ejecuta en este orden:

1. `sql/pre_gorm.sql` — schemas, extensiones, ENUMs, funciones compartidas. Idempotente (`IF NOT EXISTS`, `DO $$ ... EXCEPTION WHEN duplicate_object`).
2. `gorm.AutoMigrate(entidades...)` — entidades Go en `../entities/*.go` crean/actualizan tablas, columnas, indices y constraints anotados con tags GORM.
3. `sql/post_gorm.sql` — todo lo que GORM no expresa. Idempotente (`CREATE OR REPLACE`, `IF NOT EXISTS`, `DO $$ ... duplicate_object`).

## Proceso de cambio

### Cambio pequeno (agregar columna, cambiar tipo, modificar default)
1. Editar la struct correspondiente en `../entities/*.go` (agregar campo + tags GORM).
2. Si requiere CHECK/regex/indice parcial → agregar a `sql/post_gorm.sql` (idempotente).
3. Incrementar `SchemaVersion` en `version.go`.
4. Recrear BD con `cd ../../edugo-dev-environment/migrator && make docker-recreate`.

### Cambio fuerte (nueva tabla, reestructuracion, cambio de relaciones)
1. Crear nuevo entity en `../entities/*.go` con `TableName()` apuntando al schema correcto (ej `iam.role_grants`).
2. Registrar el entity en `migrate.go` (`AllModels`).
3. Agregar constraints/indices/funciones SQL que GORM no exprese en `sql/post_gorm.sql`.
4. Incrementar `SchemaVersion` en `version.go`.
5. Recrear BD con `make docker-recreate`.

## Numeracion de bloques SQL (en pre_gorm.sql y post_gorm.sql)

- `000-009`: Schemas, extensiones, funciones compartidas, version
- `010-019`: auth (users, refresh_tokens, login_attempts)
- `020-029`: iam (resources, roles, permissions, role_permissions, user_roles, role_grants, user_grants)
- `030-039`: academic (schools, academic_units, memberships, subjects, guardian_relations, concept_types)
- `040-049`: content (materials, material_versions, progress)
- `050-059`: assessment (assessments, attempts, answers, materials)
- `060-069`: ui_config (screen_templates, screen_instances, resource_screens, preferences)
- `070-079`: cross-schema (foreign keys, functions, views)
- `080-089`: audit (audit_events)

## Validacion

El migrador (`../../edugo-dev-environment/migrator`) valida automaticamente:
- Version antes y despues de migrar
- Hash de `sql/pre_gorm.sql` + `sql/post_gorm.sql` via `ComputeFilesHash()` en `version.go`
- Si la version no se incremento pero hay cambios en SQL, el hash lo delata (los entities no entran en el hash; un cambio solo en entities igual requiere bump por la regla 1)
