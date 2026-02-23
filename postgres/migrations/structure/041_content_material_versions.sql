-- ============================================================
-- 041: content.material_versions
-- Schema: content
-- Versiones de materiales educativos
-- Cross-schema FK (created_by -> auth.users) va en 070
-- ============================================================

CREATE TABLE content.material_versions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    material_id uuid NOT NULL,
    version_number integer NOT NULL,
    file_url text NOT NULL,
    file_type character varying(100) NOT NULL,
    file_size_bytes bigint NOT NULL,
    created_by uuid,
    created_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT material_versions_pkey PRIMARY KEY (id),
    CONSTRAINT material_versions_unique UNIQUE (material_id, version_number),
    -- Intra-schema FK
    CONSTRAINT material_versions_material_fkey FOREIGN KEY (material_id) REFERENCES content.materials(id) ON DELETE CASCADE
);
