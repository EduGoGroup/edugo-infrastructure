# Reglas para modificar migraciones

## OBLIGATORIO al cambiar cualquier archivo aqui

1. **Incrementar `SchemaVersion`** en `version.go` (mismo directorio). Sin excepcion.
2. **NUNCA usar ALTER TABLE.** Modifica el `CREATE TABLE` original en `structure/`.
3. **NUNCA crear archivos de migracion incremental** (ej: `090_alter_users_add_column.sql`). Edita el archivo original.

## Proceso de cambio

### Cambio pequeno (agregar columna, cambiar tipo, modificar default)
1. Editar el `CREATE TABLE` en `structure/XXX_nombre.sql`
2. Incrementar `SchemaVersion` en `version.go`
3. Si es urgente: aplicar el cambio directo en Neon con SQL manual, validando que el script quedo actualizado para futuras recreaciones
4. Si no es urgente: recrear BD con `make neon-recreate` desde `edugo-dev-environment`

### Cambio fuerte (nueva tabla, reestructuracion, cambio de relaciones)
1. Editar/crear archivos en `structure/`
2. Incrementar `SchemaVersion` en `version.go`
3. Recrear BD con `make neon-recreate` desde `edugo-dev-environment`

## Numeracion de archivos

- `000-009`: Schemas, extensiones, funciones compartidas, version
- `010-019`: auth (users, refresh_tokens, login_attempts)
- `020-029`: iam (resources, roles, permissions, role_permissions, user_roles)
- `030-039`: academic (schools, academic_units, memberships, subjects, guardian_relations, concept_types)
- `040-049`: content (materials, material_versions, progress)
- `050-059`: assessment (assessments, attempts, answers, materials)
- `060-069`: ui_config (screen_templates, screen_instances, resource_screens, preferences)
- `070-079`: cross-schema (foreign keys, functions, views)
- `080-089`: audit (audit_events)

## Validacion

El migrador (`edugo-dev-environment/migrator`) valida automaticamente:
- Version antes y despues de migrar
- Hash de todos los archivos SQL (detecta cambios sin bump de version)
- Si la version no se incremento pero hay cambios, el hash lo delata
