-- ============================================================
-- 002: Schema Version Tracking
-- Tabla para rastrear versiones de migracion y validar integridad.
-- Vive en public schema para sobrevivir a recreaciones parciales.
-- ============================================================

CREATE TABLE IF NOT EXISTS public.schema_version (
    id              SERIAL PRIMARY KEY,
    version         VARCHAR(20)  NOT NULL,
    content_hash    VARCHAR(64)  NOT NULL,
    execution_id    UUID         NOT NULL DEFAULT gen_random_uuid(),
    forced          BOOLEAN      DEFAULT false,
    applied_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    applied_by      VARCHAR(100) DEFAULT 'migrator',
    description     TEXT
);

COMMENT ON TABLE public.schema_version IS 'Tracking de versiones de migracion - el migrador valida version y execution_id antes/despues de cada ejecucion';
COMMENT ON COLUMN public.schema_version.version IS 'Version semantica definida en migrations/version.go — se incrementa manualmente al cambiar scripts';
COMMENT ON COLUMN public.schema_version.content_hash IS 'SHA256 de todos los archivos SQL — detecta cambios sin bump de version';
COMMENT ON COLUMN public.schema_version.execution_id IS 'UUID aleatorio generado al insertar — NUNCA sera igual, prueba que la migracion realmente se ejecuto';
COMMENT ON COLUMN public.schema_version.forced IS 'true si se ejecuto con FORCE_MIGRATION=true';
