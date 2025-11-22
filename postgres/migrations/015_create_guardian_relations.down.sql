-- Migration: 015_create_guardian_relations (DOWN)
-- Description: Revierte la creación de la tabla guardian_relations
-- Created: 2025-11-22

DROP INDEX IF EXISTS idx_guardian_relations_active_student;
DROP INDEX IF EXISTS idx_guardian_relations_active_guardian;
DROP INDEX IF EXISTS idx_guardian_relations_created_at;
DROP INDEX IF EXISTS idx_guardian_relations_is_active;
DROP INDEX IF EXISTS idx_guardian_relations_relationship_type;
DROP INDEX IF EXISTS idx_guardian_relations_student_id;
DROP INDEX IF EXISTS idx_guardian_relations_guardian_id;

DROP TABLE IF EXISTS guardian_relations;
