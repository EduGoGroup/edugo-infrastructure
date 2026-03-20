# PostgreSQL Migrations

Scripts SQL que definen la estructura completa de la base de datos.

## Estructura

```
structure/
  000_schemas_and_extensions.sql   — Schemas de dominio + extensiones
  001_shared_functions.sql         — Funciones trigger compartidas
  002_schema_version.sql           — Tabla de tracking de versiones
  010-012_auth_*.sql               — Autenticacion (users, tokens, login)
  020-024_iam_*.sql                — IAM (resources, roles, permissions)
  030-037_academic_*.sql           — Academico (schools, units, members, subjects)
  040-042_content_*.sql            — Contenido (materials, versions, progress)
  050-053_assessment_*.sql         — Evaluaciones (assessments, attempts, answers)
  060-063_ui_config_*.sql          — UI dinamica (templates, instances, preferences)
  070-072_cross_schema_*.sql       — FK entre schemas, funciones IAM, vistas
  080_audit_events.sql             — Auditoria
```

## Version

La version actual esta en `version.go` → constante `SchemaVersion`.

El migrador valida version + hash SHA256 de todos los archivos antes y despues de cada migracion.

## Como ejecutar

```bash
cd /Users/jhoanmedina/source/EduGo/EduBack/edugo-dev-environment

# Ver version actual de Neon
make neon-status

# Recrear Neon desde cero (borra todo)
make neon-recreate

# Recrear BD local (Docker)
make db-recreate
```

## Regla de oro

**NUNCA ALTER TABLE.** Siempre editar el CREATE TABLE original y recrear.
