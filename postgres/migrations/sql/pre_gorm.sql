-- ============================================================
-- pre_gorm.sql: Schemas, Extensions, ENUM types, Shared Functions
-- Runs BEFORE gorm.AutoMigrate()
-- All statements are idempotent (IF NOT EXISTS / CREATE OR REPLACE / DO blocks)
-- ============================================================

-- Domain schemas
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS iam;
CREATE SCHEMA IF NOT EXISTS academic;
CREATE SCHEMA IF NOT EXISTS content;
CREATE SCHEMA IF NOT EXISTS assessment;
CREATE SCHEMA IF NOT EXISTS ui_config;
CREATE SCHEMA IF NOT EXISTS audit;
CREATE SCHEMA IF NOT EXISTS notifications;

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schema version table (lives in public to survive partial schema drops)
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

-- ENUM types (idempotent via DO blocks — PostgreSQL raises duplicate_object on re-creation)
DO $$ BEGIN
    CREATE TYPE iam.permission_scope AS ENUM ('system', 'school', 'unit');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE iam.role_scope AS ENUM ('system', 'school', 'unit');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Shared trigger function: auto-update updated_at on every row update
CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- Shared trigger function: prevent cycles in academic_units hierarchy
CREATE OR REPLACE FUNCTION public.prevent_academic_unit_cycles()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    current_parent_id UUID;
    visited_ids UUID[];
    depth INTEGER := 0;
    max_depth INTEGER := 50;
BEGIN
    IF NEW.parent_unit_id IS NULL THEN
        RETURN NEW;
    END IF;

    current_parent_id := NEW.parent_unit_id;
    visited_ids := ARRAY[]::UUID[];

    IF NEW.id IS NOT NULL THEN
        visited_ids := array_append(visited_ids, NEW.id);
    END IF;

    WHILE current_parent_id IS NOT NULL AND depth < max_depth LOOP
        IF current_parent_id = ANY(visited_ids) THEN
            RAISE EXCEPTION 'Ciclo detectado en jerarquía: no se puede asignar % como padre de %',
                NEW.parent_unit_id, NEW.id;
        END IF;

        visited_ids := array_append(visited_ids, current_parent_id);

        SELECT parent_unit_id INTO current_parent_id
        FROM academic.academic_units
        WHERE id = current_parent_id;

        depth := depth + 1;
    END LOOP;

    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;

    RETURN NEW;
END;
$$;
