-- Migration: 014_create_units (DOWN)
-- Description: Revierte la creación de la tabla units
-- Created: 2025-11-22

DROP INDEX IF EXISTS idx_units_hierarchy;
DROP INDEX IF EXISTS idx_units_created_at;
DROP INDEX IF EXISTS idx_units_is_active;
DROP INDEX IF EXISTS idx_units_name;
DROP INDEX IF EXISTS idx_units_parent_unit_id;
DROP INDEX IF EXISTS idx_units_school_id;

DROP TABLE IF EXISTS units;
