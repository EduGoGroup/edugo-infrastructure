# Sistema de Migraciones PostgreSQL — EduGo

## Estructura

```
postgres/
  MIGRATIONS.md             ← Este archivo
  migrations/               ← Cambios de schema versionados
    001_baseline.sql        ← Schema completo capturado desde producción (2026-02-22)
    002_fix_rbac_school_id.sql ← R4: poblar user_roles.school_id desde memberships
    003_fix_timestamps.sql  ← R5: convertir TIMESTAMP a TIMESTAMPTZ en tablas RBAC
  seeds/
    production/             ← Datos base del sistema. Van a TODOS los entornos.
      001_resources.sql     ← 18 recursos (menú, permisos, mobile)
      002_roles.sql         ← 11 roles del sistema
      003_permissions.sql   ← 52 permisos
      004_role_permissions.sql ← 109 asignaciones rol→permiso
    development/            ← Datos de prueba coherentes. Solo para dev/staging.
```

## Convenciones

| Regla | Descripción |
|-------|-------------|
| Numeración | `NNN_descripcion.sql` — tres dígitos, incrementales, sin saltos |
| Idempotencia | Todos los seeds usan `ON CONFLICT DO NOTHING` |
| Un cambio por archivo | Cada migración resuelve un problema concreto y documentado |
| Header obligatorio | Cada archivo inicia con bloque de comentario que indica qué hace y por qué |
| No bajar schema | Las migraciones solo avanzan; nunca revierten automáticamente |

## Cómo ejecutar

### Requisitos
- `psql` >= 14
- Variable de entorno `DATABASE_URL` o credenciales explícitas

### Ejecución manual (orden estricto)

```bash
# 1. Baseline (solo en DB nueva / vacía)
psql "$DATABASE_URL" -f postgres/migrations/001_baseline.sql

# 2. Correcciones de datos y schema
psql "$DATABASE_URL" -f postgres/migrations/002_fix_rbac_school_id.sql
psql "$DATABASE_URL" -f postgres/migrations/003_fix_timestamps.sql

# 3. Seeds de producción (idempotentes, seguros de re-ejecutar)
psql "$DATABASE_URL" -f postgres/seeds/production/001_resources.sql
psql "$DATABASE_URL" -f postgres/seeds/production/002_roles.sql
psql "$DATABASE_URL" -f postgres/seeds/production/003_permissions.sql
psql "$DATABASE_URL" -f postgres/seeds/production/004_role_permissions.sql

# 4. Seeds de desarrollo (solo en entornos no-productivos)
# psql "$DATABASE_URL" -f postgres/seeds/development/001_users.sql
# psql "$DATABASE_URL" -f postgres/seeds/development/002_schools.sql
```

### Ejecución con script (todos los seeds de producción)

```bash
for f in postgres/seeds/production/*.sql; do
  echo "Ejecutando $f..."
  psql "$DATABASE_URL" -f "$f"
done
```

## Descripción de cada migración

### 001_baseline.sql
Schema completo capturado con `pg_dump --schema-only` desde el entorno de
producción Neon el 2026-02-22. Incluye:
- Schemas: `public`, `ui_config`
- Tipos ENUM: `permission_scope`, `role_scope`
- Funciones: `get_user_permissions`, `update_updated_at_column`
- Todas las tablas, índices, constraints y triggers

**IMPORTANTE:** Solo ejecutar en una base de datos vacía. Si la DB ya tiene
el schema (porque se creó por otro medio), omitir este archivo y comenzar
desde 002.

### 002_fix_rbac_school_id.sql
**Problema (R4):** `user_roles.school_id` estaba NULL para todos los usuarios
con roles de scope `school`. El sistema de permisos no podía resolver
correctamente el contexto de escuela.

**Solución:** Pobla `school_id` desde la primera membresía activa del usuario
(`memberships.enrolled_at ASC`). Los usuarios con roles de sistema
(`super_admin`, `platform_admin`) quedan con `school_id = NULL`, que es
el comportamiento correcto.

### 003_fix_timestamps.sql
**Problema (R5):** Las tablas `roles`, `permissions`, `resources` y
`user_roles` usaban `TIMESTAMP WITHOUT TIME ZONE`. El resto del schema ya
usaba `TIMESTAMP WITH TIME ZONE`. Esta inconsistencia podía causar errores
silenciosos al comparar timestamps entre tablas en entornos con timezone
distinto a UTC.

**Solución:** Convierte todas las columnas de timestamp a `TIMESTAMPTZ`
asumiendo que los valores almacenados estaban en UTC (convención del sistema).

## Seeds de producción

Los seeds son **idempotentes** (`ON CONFLICT DO NOTHING`) y representan el
estado canónico de la configuración del sistema RBAC. Deben aplicarse en
**todos los entornos** (producción, staging, desarrollo).

| Archivo | Registros | Descripción |
|---------|-----------|-------------|
| 001_resources.sql | 18 | Recursos del sistema (menú + mobile) |
| 002_roles.sql | 11 | Roles: 2 system, 4 school, 5 unit |
| 003_permissions.sql | 52 | Permisos agrupados por recurso |
| 004_role_permissions.sql | 109 | Asignaciones: super_admin (52), platform_admin (10), school_admin (18), teacher (16), student (8), guardian (5) |

## Seeds de desarrollo

El directorio `seeds/development/` está reservado para datos de prueba
coherentes (usuarios de demo, escuelas de ejemplo, membresías, etc.).
Estos seeds **no deben aplicarse en producción**.

## Agregar una nueva migración

1. Crear el archivo con el siguiente número: `NNN_descripcion.sql`
2. Incluir el bloque de comentario estándar al inicio
3. Verificar en un entorno de staging antes de aplicar a producción
4. Si el cambio es reversible, documentar el rollback en el comentario

```sql
-- ============================================================
-- MIGRACIÓN NNN: Descripción breve
-- Fecha: YYYY-MM-DD
-- Autor: nombre
-- Motivo: por qué es necesario este cambio
-- Rollback: cómo revertirlo si es necesario
-- ============================================================

-- SQL aquí
```

## Herramientas de migración

Este repositorio usa migraciones SQL planas por simplicidad y portabilidad.
Si en el futuro se requiere tracking de versiones aplicadas, considerar:
- [golang-migrate](https://github.com/golang-migrate/migrate) — compatible con Go
- [goose](https://github.com/pressly/goose) — alternativa liviana en Go
- [Flyway](https://flywaydb.org/) — opción robusta con soporte de checksum

La estructura de archivos actual es compatible con `golang-migrate` usando
el driver `postgres` y formato de archivo `NNN_name.up.sql`.
