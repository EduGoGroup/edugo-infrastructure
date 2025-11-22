-- Migration: 013_create_subjects (DOWN)
-- Description: Revierte la creación de la tabla subjects
-- Created: 2025-11-22

DROP INDEX IF EXISTS idx_subjects_metadata;
DROP INDEX IF EXISTS idx_subjects_created_at;
DROP INDEX IF EXISTS idx_subjects_is_active;
DROP INDEX IF EXISTS idx_subjects_name;

DROP TABLE IF EXISTS subjects;
