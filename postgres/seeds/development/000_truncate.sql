-- =============================================================================
-- EduGo Development Seeds v2 — 000_truncate.sql
-- =============================================================================
-- Limpia todas las tablas de datos de desarrollo en orden correcto
-- respetando las foreign keys (de hojas hacia raices).
--
-- ADVERTENCIA: No truncar roles, permissions, resources ni role_permissions
-- ya que esos son seeds de produccion que deben persistir.
--
-- MAPA DE IDs FIJOS (v2):
--   Escuelas        b1000000-...-001 (San Ignacio) / b2000000-...-002 (CreArte) / b3000000-...-003 (Academia)
--   Academic units  ac000000-...-001 a 016
--   Users           00000000-...-001 a 014
--   Memberships     bb000000-...-001 a 021
--   User roles      cc000000-...-001 a 019
--   Subjects        dd000000-...-001 a 007
--   Materials       aa100000-...-001 a 005
--   Assessments     aa200000-...-001 a 006
--   Ass. Materials  ab100000-...-001 a 005
--   Attempts        aa300000-...-001 a 007
--   Guardian rels   ee000000-...-001 a 003
-- =============================================================================

BEGIN;

-- Nivel 5 — hojas
TRUNCATE TABLE assessment.assessment_attempt_answer CASCADE;

-- Nivel 4
TRUNCATE TABLE assessment.assessment_attempt CASCADE;

-- Nivel 3.5
TRUNCATE TABLE assessment.assessment_materials CASCADE;

-- Nivel 3
TRUNCATE TABLE assessment.assessment CASCADE;

-- Nivel 2
TRUNCATE TABLE content.materials CASCADE;
TRUNCATE TABLE academic.memberships CASCADE;
TRUNCATE TABLE iam.user_roles CASCADE;

-- Nivel 1
TRUNCATE TABLE academic.academic_units CASCADE;

-- Desacoplar ui_config de auth.users para evitar que CASCADE destruya datos de produccion
UPDATE ui_config.screen_templates SET created_by = NULL WHERE created_by IS NOT NULL;
UPDATE ui_config.screen_instances SET created_by = NULL WHERE created_by IS NOT NULL;

-- Tablas adicionales de desarrollo (hojas)
TRUNCATE TABLE academic.subjects CASCADE;
TRUNCATE TABLE academic.guardian_relations CASCADE;
TRUNCATE TABLE academic.school_concepts CASCADE;
TRUNCATE TABLE ui_config.screen_user_preferences CASCADE;

-- Nivel 0 — raices (usar DELETE para auth.users para no cascadear a ui_config)
TRUNCATE TABLE auth.refresh_tokens CASCADE;
TRUNCATE TABLE auth.login_attempts CASCADE;
TRUNCATE TABLE academic.schools CASCADE;
DELETE FROM auth.users;

COMMIT;
