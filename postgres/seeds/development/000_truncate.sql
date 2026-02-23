-- =============================================================================
-- EduGo Development Seeds — 000_truncate.sql
-- =============================================================================
-- Limpia todas las tablas de datos de desarrollo en orden correcto
-- respetando las foreign keys (de hojas hacia raíces).
--
-- ADVERTENCIA: No truncar roles, permissions, resources ni role_permissions
-- ya que esos son seeds de producción que deben persistir.
--
-- Ejecutar antes de re-aplicar los seeds de desarrollo.
--
-- MAPA DE IDs FIJOS (solo hex válido para tipo UUID):
--   Escuelas        b1000000-...-001/002/003
--   Academic units  ac000000-...-001 a 009
--   Users           00000000-...-001 a 013
--   Memberships     bb000000-...-001 a 012
--   User roles      cc000000-...-001 a 013
--   Materials       aa100000-...-001 a 004
--   Assessments     aa200000-...-001 a 002
--   Attempts        aa300000-...-001 a 004
-- =============================================================================

BEGIN;

-- Nivel 5 — hojas
TRUNCATE TABLE assessment.assessment_attempt_answer CASCADE;

-- Nivel 4
TRUNCATE TABLE assessment.assessment_attempt CASCADE;

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

-- Nivel 0 — raíces (usar DELETE para auth.users para no cascadear a ui_config)
TRUNCATE TABLE auth.refresh_tokens CASCADE;
TRUNCATE TABLE auth.login_attempts CASCADE;
TRUNCATE TABLE academic.schools CASCADE;
DELETE FROM auth.users;

-- Tablas adicionales de desarrollo
TRUNCATE TABLE academic.subjects CASCADE;
TRUNCATE TABLE academic.guardian_relations CASCADE;
TRUNCATE TABLE ui_config.screen_user_preferences CASCADE;

COMMIT;
