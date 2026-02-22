-- ============================================================
-- MIGRACIÓN 003: Unificar timestamps a TIMESTAMPTZ en tablas RBAC
-- Fecha: 2026-02-22
-- R5 del análisis arquitectónico Opus
-- Problema: las tablas roles, permissions, resources y user_roles
--           usan TIMESTAMP WITHOUT TIME ZONE. El resto del schema
--           ya usa TIMESTAMP WITH TIME ZONE. Esta migración unifica.
--
-- Tablas afectadas:
--   roles       → created_at, updated_at
--   permissions → created_at, updated_at
--   resources   → created_at, updated_at
--   user_roles  → granted_at, expires_at, created_at, updated_at
--
-- Nota: Las columnas se convierten asumiendo que los valores
--       almacenados están en UTC (convención del sistema).
-- ============================================================

-- roles: no tiene granted_at (solo created_at y updated_at)
ALTER TABLE roles
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE
        USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE
        USING updated_at AT TIME ZONE 'UTC';

-- permissions
ALTER TABLE permissions
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE
        USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE
        USING updated_at AT TIME ZONE 'UTC';

-- resources
ALTER TABLE resources
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE
        USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE
        USING updated_at AT TIME ZONE 'UTC';

-- user_roles: tiene granted_at y expires_at adicionales
ALTER TABLE user_roles
    ALTER COLUMN granted_at TYPE TIMESTAMP WITH TIME ZONE
        USING granted_at AT TIME ZONE 'UTC',
    ALTER COLUMN expires_at TYPE TIMESTAMP WITH TIME ZONE
        USING expires_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE
        USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at TYPE TIMESTAMP WITH TIME ZONE
        USING updated_at AT TIME ZONE 'UTC';

-- Verificación: muestra tipos actuales de las columnas afectadas
DO $$
DECLARE
    col_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO col_count
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND table_name IN ('roles', 'permissions', 'resources', 'user_roles')
      AND column_name IN ('created_at', 'updated_at', 'granted_at', 'expires_at')
      AND data_type = 'timestamp with time zone';

    RAISE NOTICE 'Columnas de timestamp convertidas a TIMESTAMPTZ: % (esperado: 8)', col_count;
END $$;
