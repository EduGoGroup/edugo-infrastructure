-- Migration: 012_create_material_versions (DOWN)
-- Description: Revierte la creación de la tabla material_versions
-- Created: 2025-11-22

DROP INDEX IF EXISTS idx_material_versions_created_at;
DROP INDEX IF EXISTS idx_material_versions_changed_by;
DROP INDEX IF EXISTS idx_material_versions_version_number;
DROP INDEX IF EXISTS idx_material_versions_material_id;

DROP TABLE IF EXISTS material_versions;
