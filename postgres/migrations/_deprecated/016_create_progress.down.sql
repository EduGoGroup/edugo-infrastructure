-- Migration: 016_create_progress (DOWN)
-- Description: Revierte la creación de la tabla progress
-- Created: 2025-11-22

DROP INDEX IF EXISTS idx_progress_material_status;
DROP INDEX IF EXISTS idx_progress_user_status;
DROP INDEX IF EXISTS idx_progress_percentage;
DROP INDEX IF EXISTS idx_progress_last_accessed_at;
DROP INDEX IF EXISTS idx_progress_status;
DROP INDEX IF EXISTS idx_progress_material_id;
DROP INDEX IF EXISTS idx_progress_user_id;

DROP TABLE IF EXISTS progress;
