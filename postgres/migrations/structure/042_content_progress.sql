-- ============================================================
-- 042: content.progress
-- Schema: content
-- Progreso de usuarios en materiales educativos
-- Cross-schema FK (user_id -> auth.users) va en 070
-- ============================================================

CREATE TABLE content.progress (
    material_id uuid NOT NULL,
    user_id uuid NOT NULL,
    progress_percentage numeric(5,2) DEFAULT 0 NOT NULL,
    last_position jsonb DEFAULT '{}'::jsonb,
    completed_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT progress_pkey PRIMARY KEY (material_id, user_id),
    -- Intra-schema FK
    CONSTRAINT progress_material_fkey FOREIGN KEY (material_id) REFERENCES content.materials(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON content.progress
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
