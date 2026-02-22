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

-- Nivel 5 — tablas sin dependientes (hojas)
TRUNCATE TABLE public.assessment_attempt_answer CASCADE;

-- Nivel 4
TRUNCATE TABLE public.assessment_attempt CASCADE;

-- Nivel 3
TRUNCATE TABLE public.assessment CASCADE;

-- Nivel 2
TRUNCATE TABLE public.materials CASCADE;
TRUNCATE TABLE public.memberships CASCADE;
TRUNCATE TABLE public.user_roles CASCADE;

-- Nivel 1
TRUNCATE TABLE public.academic_units CASCADE;

-- Nivel 0 — raíces (schools y users no tienen dependencias entre sí excepto a través de las tablas ya truncadas)
TRUNCATE TABLE public.refresh_tokens CASCADE;
TRUNCATE TABLE public.login_attempts CASCADE;
TRUNCATE TABLE public.schools CASCADE;
TRUNCATE TABLE public.users CASCADE;

COMMIT;
