-- ============================================================
-- post_gorm.sql: Triggers, Views, IAM Functions, Partial Indexes,
--                Analytics tables (no entity files), extra FKs
-- Runs AFTER gorm.AutoMigrate()
-- All statements are idempotent (CREATE OR REPLACE / IF NOT EXISTS)
-- Requires PostgreSQL 14+ for CREATE OR REPLACE TRIGGER
-- ============================================================

-- ============================================================
-- Analitica de evaluacion: DIFERIDA en N4 (ADR 0019).
-- Las tablas viejas assessment.attempt_analytics y assessment.assessment_stats
-- (llaveadas a auth.users, modelo global muerto) se ELIMINAN. El resumen del
-- docente se computa on-the-fly; si luego se requiere analitica materializada,
-- se rehace por membership (deuda anotada).
-- ============================================================

-- ============================================================
-- Extra FK that GORM cannot express (cross-schema, non-entity-annotated)
-- ============================================================

-- ui_config.resource_screens.screen_key → ui_config.screen_instances.screen_key
-- (noted in resource_screen.go entity comment)
DO $$ BEGIN
    ALTER TABLE ui_config.resource_screens
        ADD CONSTRAINT fk_resource_screens_screen_key
            FOREIGN KEY (screen_key) REFERENCES ui_config.screen_instances(screen_key);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- academic.school_invitations / academic.school_join_requests (onboarding, plan 005)
-- GORM no materializa FKs desde el tag `constraint:` cuando la entity no
-- declara un campo de relación (mismo caso que academic.guardian_relations,
-- que termina sin FKs). Por eso se declaran aquí explícitamente, espejando
-- los nombres de constraint anotados en las entities. Idempotente.

-- school_invitations → schools / academic_units / users (created_by)
DO $$ BEGIN
    ALTER TABLE academic.school_invitations
        ADD CONSTRAINT school_invitations_school_fkey
            FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_invitations
        ADD CONSTRAINT school_invitations_unit_fkey
            FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_invitations
        ADD CONSTRAINT school_invitations_created_by_fkey
            FOREIGN KEY (created_by) REFERENCES auth.users(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- school_join_requests → users / schools / academic_units / invitation + aprobadores
DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_user_fkey
            FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_school_fkey
            FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_unit_fkey
            FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_invitation_fkey
            FOREIGN KEY (invitation_id) REFERENCES academic.school_invitations(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_school_approver_fkey
            FOREIGN KEY (school_approved_by) REFERENCES auth.users(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_unit_approver_fkey
            FOREIGN KEY (unit_approved_by) REFERENCES auth.users(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.school_join_requests
        ADD CONSTRAINT school_join_requests_rejected_by_fkey
            FOREIGN KEY (rejected_by) REFERENCES auth.users(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- academic.subject_offerings / academic.subject_offering_enrollments
-- (sesiones de materia + inscripcion, ADR 0009 / plan 010 N1.7).
-- GORM no materializa FKs desde el tag `constraint:` sin campo de relacion
-- (mismo caso que academic.subjects). Se declaran aqui espejando los nombres
-- de constraint anotados en las entities.
-- Idempotente.

-- subject_offerings → schools / subjects / academic_units / periods / teacher membership
DO $$ BEGIN
    ALTER TABLE academic.subject_offerings
        ADD CONSTRAINT subject_offerings_school_fkey
            FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.subject_offerings
        ADD CONSTRAINT subject_offerings_subject_fkey
            FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.subject_offerings
        ADD CONSTRAINT subject_offerings_unit_fkey
            FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.subject_offerings
        ADD CONSTRAINT subject_offerings_period_fkey
            FOREIGN KEY (period_id) REFERENCES academic.academic_periods(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Docente = SET NULL: borrar/expirar la membresia del docente no debe borrar
-- la sesion ni desinscribir alumnos; la sesion queda sin docente asignado.
DO $$ BEGIN
    ALTER TABLE academic.subject_offerings
        ADD CONSTRAINT subject_offerings_teacher_fkey
            FOREIGN KEY (teacher_membership_id) REFERENCES academic.memberships(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- subject_offering_enrollments → subject_offerings / membership del alumno
DO $$ BEGIN
    ALTER TABLE academic.subject_offering_enrollments
        ADD CONSTRAINT subject_offering_enrollments_offering_fkey
            FOREIGN KEY (offering_id) REFERENCES academic.subject_offerings(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.subject_offering_enrollments
        ADD CONSTRAINT subject_offering_enrollments_student_fkey
            FOREIGN KEY (student_membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- subject_id denormalizado (copia inmutable de la oferta) → materia. Respalda el
-- invariante "una oferta por materia por alumno" (uniqueIndex
-- uq_enrollment_student_subject, materializado por GORM); bug 0036.
DO $$ BEGIN
    ALTER TABLE academic.subject_offering_enrollments
        ADD CONSTRAINT subject_offering_enrollments_subject_fkey
            FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- academic.grades.teacher_id → academic.memberships (docente que califica).
-- GORM no materializa esta FK desde el tag `constraint:grades_teacher_fkey` de la
-- entity Grade porque no declara campo de relacion (mismo caso que
-- subject_offerings.teacher_membership_id). teacher_id es nullable, asi que el
-- docente que califica se desvincula con SET NULL al borrar/expirar su membresia:
-- la nota persiste sin docente asignado (paridad con subject_offerings_teacher_fkey).
DO $$ BEGIN
    ALTER TABLE academic.grades
        ADD CONSTRAINT grades_teacher_fkey
            FOREIGN KEY (teacher_id) REFERENCES academic.memberships(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================
-- academic.grade_item / academic.grade_history (N4 / ADR 0020) — componentes de
-- nota, procedencia y auditoria de override. GORM no materializa FKs desde el tag
-- `constraint:` sin campo de relacion (mismo patron que academic.grades.teacher_id),
-- por eso TODAS las FKs (academic y cross-schema a assessment.*), el CHECK XOR y el
-- UNIQUE parcial viven aqui. Idempotente.
-- ============================================================

-- academic.grade_item → memberships (CASCADE) / subjects (CASCADE) / periods
-- (CASCADE) / membership autor (RESTRICT)
DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_membership_fkey
            FOREIGN KEY (membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_subject_fkey
            FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_period_fkey
            FOREIGN KEY (period_id) REFERENCES academic.academic_periods(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_created_by_fkey
            FOREIGN KEY (created_by_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Cross-schema: trazabilidad al origen auto_scored/auto_llm (SET NULL: la nota
-- persiste si se borra el intento/evaluacion de origen).
DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_source_attempt_fkey
            FOREIGN KEY (source_attempt_id) REFERENCES assessment.assessment_attempt(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_item
        ADD CONSTRAINT grade_item_source_assessment_fkey
            FOREIGN KEY (source_assessment_id) REFERENCES assessment.assessment(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- academic.grade_history → grades (CASCADE) / grade_item (CASCADE) / membership
-- que cambia (RESTRICT)
DO $$ BEGIN
    ALTER TABLE academic.grade_history
        ADD CONSTRAINT grade_history_grade_fkey
            FOREIGN KEY (grade_id) REFERENCES academic.grades(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_history
        ADD CONSTRAINT grade_history_item_fkey
            FOREIGN KEY (grade_item_id) REFERENCES academic.grade_item(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE academic.grade_history
        ADD CONSTRAINT grade_history_changed_by_fkey
            FOREIGN KEY (changed_by_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- XOR: cada fila de auditoria apunta a EXACTAMENTE UNO de grade_id / grade_item_id.
DO $$ BEGIN
    ALTER TABLE academic.grade_history
        ADD CONSTRAINT grade_history_target_xor_check
            CHECK (((grade_id IS NOT NULL)::int + (grade_item_id IS NOT NULL)::int) = 1);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================
-- assessment.* y content.* (N4 / ADR 0019) — esquema de evaluacion/contenido
-- anclado al modelo de sesion. GORM no materializa FKs desde el tag
-- `constraint:` sin campo de relacion (mismo caso que subject_offerings),
-- por eso TODAS las FKs cross-schema y el UNIQUE de assignment se declaran
-- aqui, espejando los nombres de constraint anotados en las entities.
-- Idempotente.
-- ============================================================

-- content.materials → schools (CASCADE) / subjects (SET NULL, nullable) /
-- academic_units (SET NULL, nullable) / membership autor (RESTRICT)
DO $$ BEGIN
    ALTER TABLE content.materials
        ADD CONSTRAINT materials_school_fkey
            FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE content.materials
        ADD CONSTRAINT materials_subject_fkey
            FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE content.materials
        ADD CONSTRAINT materials_unit_fkey
            FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE content.materials
        ADD CONSTRAINT materials_membership_fkey
            FOREIGN KEY (uploaded_by_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- content.material_file (rediseño F2 plan 018): su unica FK material_id →
-- content.materials(id) ON DELETE CASCADE es same-schema y la materializa GORM
-- desde el tag `constraint:material_file_material_fkey`, por eso NO se declara
-- aqui (a diferencia de las FKs cross-schema). La vieja content.material_version
-- (y sus FKs material_version_material_fkey / material_version_membership_fkey)
-- fue ELIMINADA en F2.

-- assessment.assessment → schools (CASCADE) / membership autor (RESTRICT) /
-- subjects (RESTRICT, catalogo de escuela)
DO $$ BEGIN
    ALTER TABLE assessment.assessment
        ADD CONSTRAINT assessment_school_fkey
            FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment
        ADD CONSTRAINT assessment_created_by_fkey
            FOREIGN KEY (created_by_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment
        ADD CONSTRAINT assessment_subject_fkey
            FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.question → assessment (CASCADE)
DO $$ BEGIN
    ALTER TABLE assessment.question
        ADD CONSTRAINT question_assessment_fkey
            FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.question_option → question (CASCADE)
DO $$ BEGIN
    ALTER TABLE assessment.question_option
        ADD CONSTRAINT question_option_question_fkey
            FOREIGN KEY (question_id) REFERENCES assessment.question(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.assessment_material (N:N) → assessment (CASCADE) / materials (CASCADE)
DO $$ BEGIN
    ALTER TABLE assessment.assessment_material
        ADD CONSTRAINT assessment_material_assessment_fkey
            FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment_material
        ADD CONSTRAINT assessment_material_material_fkey
            FOREIGN KEY (material_id) REFERENCES content.materials(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.assessment_assignment → assessment (CASCADE) / subject_offering
-- (CASCADE, la sesion/grupo) / membership que asigna (RESTRICT)
DO $$ BEGIN
    ALTER TABLE assessment.assessment_assignment
        ADD CONSTRAINT assessment_assignment_assessment_fkey
            FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment_assignment
        ADD CONSTRAINT assessment_assignment_offering_fkey
            FOREIGN KEY (subject_offering_id) REFERENCES academic.subject_offerings(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment_assignment
        ADD CONSTRAINT assessment_assignment_assigned_by_fkey
            FOREIGN KEY (assigned_by_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- No duplicar la misma evaluacion a la misma sesion.
DO $$ BEGIN
    ALTER TABLE assessment.assessment_assignment
        ADD CONSTRAINT uq_assignment_assessment_offering
            UNIQUE (assessment_id, subject_offering_id);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.assessment_attempt → assessment (CASCADE) / membership del alumno (CASCADE)
DO $$ BEGIN
    ALTER TABLE assessment.assessment_attempt
        ADD CONSTRAINT assessment_attempt_assessment_fkey
            FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment_attempt
        ADD CONSTRAINT assessment_attempt_student_fkey
            FOREIGN KEY (student_membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.assessment_attempt_answer → attempt (CASCADE) / question (SET NULL)
DO $$ BEGIN
    ALTER TABLE assessment.assessment_attempt_answer
        ADD CONSTRAINT assessment_attempt_answer_attempt_fkey
            FOREIGN KEY (attempt_id) REFERENCES assessment.assessment_attempt(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.assessment_attempt_answer
        ADD CONSTRAINT assessment_attempt_answer_question_fkey
            FOREIGN KEY (question_id) REFERENCES assessment.question(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- assessment.attempt_review → attempt_answer (CASCADE) / membership revisor (RESTRICT)
DO $$ BEGIN
    ALTER TABLE assessment.attempt_review
        ADD CONSTRAINT attempt_review_answer_fkey
            FOREIGN KEY (attempt_answer_id) REFERENCES assessment.assessment_attempt_answer(id) ON DELETE CASCADE;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE assessment.attempt_review
        ADD CONSTRAINT attempt_review_reviewer_fkey
            FOREIGN KEY (reviewer_membership_id) REFERENCES academic.memberships(id) ON DELETE RESTRICT;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================
-- Triggers (CREATE OR REPLACE TRIGGER — PostgreSQL 14+)
-- ============================================================

-- auth.users
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON auth.users
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- iam.resources
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON iam.resources
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- iam.roles
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON iam.roles
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- iam.permissions
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON iam.permissions
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- iam.user_roles
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON iam.user_roles
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.schools
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.schools
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.academic_units
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.academic_units
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE OR REPLACE TRIGGER trigger_prevent_academic_unit_cycles
    BEFORE INSERT OR UPDATE OF parent_unit_id ON academic.academic_units
    FOR EACH ROW EXECUTE FUNCTION public.prevent_academic_unit_cycles();

-- academic.memberships
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.memberships
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.subjects
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.subjects
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.subject_offerings (sesiones de materia; enrollments no tiene
-- updated_at, solo enrolled_at, por eso no lleva trigger).
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.subject_offerings
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.guardian_relations
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.guardian_relations
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.school_invitations
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.school_invitations
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.school_join_requests
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.school_join_requests
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.concept_types
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.concept_types
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.concept_definitions
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.concept_definitions
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.school_concepts
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.school_concepts
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.academic_periods
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.academic_periods
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.grades
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.grades
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.grade_item (N4 / ADR 0020). grade_history no tiene updated_at (es
-- append-only: changed_at lo fija el insert), por eso no lleva trigger.
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.grade_item
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- academic.announcements
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON academic.announcements
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- content.materials
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON content.materials
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.assessment
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.assessment
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.question
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.question
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.question_option
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.question_option
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.assessment_assignment
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.assessment_assignment
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.assessment_attempt
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.assessment_attempt
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.assessment_attempt_answer
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.assessment_attempt_answer
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- assessment.attempt_review
CREATE OR REPLACE TRIGGER set_updated_at
    BEFORE UPDATE ON assessment.attempt_review
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- ui_config.screen_templates
CREATE OR REPLACE TRIGGER update_screen_templates_updated_at
    BEFORE UPDATE ON ui_config.screen_templates
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- ui_config.screen_instances
CREATE OR REPLACE TRIGGER update_screen_instances_updated_at
    BEFORE UPDATE ON ui_config.screen_instances
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- ui_config.resource_screens
CREATE OR REPLACE TRIGGER update_resource_screens_updated_at
    BEFORE UPDATE ON ui_config.resource_screens
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- ui_config.screen_user_preferences
CREATE OR REPLACE TRIGGER update_screen_user_prefs_updated_at
    BEFORE UPDATE ON ui_config.screen_user_preferences
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

-- ============================================================
-- Views
-- ============================================================

CREATE OR REPLACE VIEW academic.v_academic_unit_tree AS
WITH RECURSIVE unit_hierarchy AS (
    SELECT id, parent_unit_id, school_id, name, code, type, level, academic_year,
           1 AS depth, ARRAY[id] AS path, name::text AS full_path
    FROM academic.academic_units
    WHERE parent_unit_id IS NULL AND deleted_at IS NULL
    UNION ALL
    SELECT au.id, au.parent_unit_id, au.school_id, au.name, au.code, au.type, au.level, au.academic_year,
           uh.depth + 1, uh.path || au.id, uh.full_path || ' > ' || au.name::text
    FROM academic.academic_units au
    JOIN unit_hierarchy uh ON au.parent_unit_id = uh.id
    WHERE au.deleted_at IS NULL
)
SELECT uh.*, s.name AS school_name, s.code AS school_code
FROM unit_hierarchy uh
LEFT JOIN academic.schools s ON uh.school_id = s.id
ORDER BY uh.school_id, uh.path;

-- ============================================================
-- IAM permissions redesign (P1-1)
--
-- P4-1 (plan B): se eliminaron las funciones iam.get_user_permissions(),
-- iam.get_user_resources() y iam.user_has_permission(). Estaban basadas
-- en la tabla legacy iam.role_permissions (enumeración 1:1 rol×permiso)
-- que ya no existe. El modelo nuevo de permisos vive en iam.role_grants
-- (patterns wildcard) e iam.user_grants. La lógica de resolución del
-- permiso efectivo se computa en backend (edugo-shared/auth.PermissionMatches)
-- y no en SQL. Ningún consumer en Go llamaba estas funciones.
-- Functions iam.permission_matches() e iam.scope_covers(),
-- CHECK constraints regex sobre iam.role_grants/iam.user_grants.
-- ============================================================

CREATE OR REPLACE FUNCTION iam.permission_matches(pattern TEXT, request TEXT)
RETURNS BOOLEAN
LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE AS $$
DECLARE
    pos     INT;
    head    TEXT;
    tail    TEXT;
    middle  TEXT;
    suffix  TEXT;
BEGIN
    -- Wildcard global
    IF pattern = '*' THEN RETURN TRUE; END IF;

    -- Match exacto
    IF pattern = request THEN RETURN TRUE; END IF;

    -- Glob prefix: 'academic.*' matches 'academic' OR 'academic.<anything>'
    IF right(pattern, 2) = '.*' THEN
        RETURN request = left(pattern, length(pattern)-2)
            OR request LIKE left(pattern, length(pattern)-1) || '%';
    END IF;

    -- Wildcard leading: '*.suffix' matches '<algo>.suffix' (sin importar profundidad).
    IF left(pattern, 2) = '*.' THEN
        suffix := substr(pattern, 2); -- ej ".create"
        RETURN length(request) > length(suffix)
           AND right(request, length(suffix)) = suffix;
    END IF;

    -- Wildcard medio: 'prefix.*.suffix' matches 'prefix.<algo>.suffix'
    -- con uno o más segmentos intermedios.
    pos := position('.*.' IN pattern);
    IF pos > 1 THEN
        head := left(pattern, pos);           -- 'prefix.'
        tail := substr(pattern, pos + 2);     -- '.suffix'
        IF position('*' IN head) > 0 OR position('*' IN tail) > 0 THEN
            RETURN FALSE;
        END IF;
        IF left(request, length(head)) <> head
           OR right(request, length(tail)) <> tail THEN
            RETURN FALSE;
        END IF;
        middle := substr(request, length(head) + 1,
                         length(request) - length(head) - length(tail));
        RETURN length(middle) > 0
           AND left(middle, 1) <> '.'
           AND right(middle, 1) <> '.';
    END IF;

    RETURN FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION iam.scope_covers(pattern TEXT, request TEXT)
RETURNS BOOLEAN
LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE AS $$
DECLARE
    p_segs TEXT[];
    r_segs TEXT[];
    i INT;
BEGIN
    IF pattern = '*' THEN RETURN TRUE; END IF;
    IF request = '*' THEN RETURN FALSE; END IF;

    p_segs := string_to_array(pattern, '/');
    r_segs := string_to_array(request, '/');

    IF array_length(p_segs, 1) > array_length(r_segs, 1) THEN
        RETURN FALSE;
    END IF;

    FOR i IN 1..array_length(p_segs, 1) LOOP
        IF p_segs[i] = r_segs[i] THEN CONTINUE; END IF;
        IF split_part(p_segs[i], ':', 2) = '*'
           AND split_part(p_segs[i], ':', 1) = split_part(r_segs[i], ':', 1)
        THEN CONTINUE; END IF;
        RETURN FALSE;
    END LOOP;

    RETURN TRUE;
END;
$$;

-- CHECK constraints on iam.role_grants
-- P1-2: regex acepta ambos formatos: legacy (recurso:accion[:own]) y
-- path-based (recurso.accion[.*][:own]). Mirror desde role_permissions
-- usa nombres legacy con `:`.
DO $$ BEGIN
    ALTER TABLE iam.role_grants
        ADD CONSTRAINT role_grants_pattern_format
        CHECK (pattern ~ '^(\*|[a-z_]+(\.[a-z_]+){0,2}(\.\*)?|\*\.[a-z_]+|[a-z_]+\.\*\.[a-z_]+|[a-z_]+(:[a-z_]+){0,1})(:own)?$');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE iam.role_grants
        ADD CONSTRAINT role_grants_effect_format
        CHECK (effect IN ('allow','deny'));
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- CHECK constraints on iam.user_grants
DO $$ BEGIN
    ALTER TABLE iam.user_grants
        ADD CONSTRAINT user_grants_scope_format
        CHECK (scope_pattern ~ '^(\*|(school|unit|section|subject):[^/]+(/(school|unit|section|subject):[^/]+)*)$');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE iam.user_grants
        ADD CONSTRAINT user_grants_permission_format
        CHECK (permission_pattern ~ '^(\*|[a-z_]+(\.[a-z_]+){0,2}(\.\*)?|\*\.[a-z_]+|[a-z_]+\.\*\.[a-z_]+|[a-z_]+(:[a-z_]+){0,1})(:own)?$');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE iam.user_grants
        ADD CONSTRAINT user_grants_effect_format
        CHECK (effect IN ('allow','deny'));
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE INDEX IF NOT EXISTS idx_user_grants_user_active
    ON iam.user_grants (user_id, effect, expires_at);

-- ADR-6 (herencia de roles): FK self-referencial nullable de
-- iam.roles.parent_role_id → iam.roles(id). GORM no la emite para una
-- columna UUID sin campo de asociación (mismo caso que resources.parent_id),
-- así que se declara aquí. ON DELETE SET NULL: borrar un canónico no
-- rompe sus alias (quedan sin parent). Idempotente.
DO $$ BEGIN
    ALTER TABLE iam.roles
        ADD CONSTRAINT fk_roles_parent
        FOREIGN KEY (parent_role_id) REFERENCES iam.roles(id) ON DELETE SET NULL;
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================
-- Partial Indexes (WHERE clauses that GORM tags cannot express)
-- ============================================================

-- auth
CREATE INDEX IF NOT EXISTS idx_users_active
    ON auth.users (is_active) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at
    ON auth.refresh_tokens (revoked_at) WHERE revoked_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_service_clients_active
    ON auth.service_clients (client_id) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_login_attempts_rate_limit
    ON auth.login_attempts (identifier, successful, attempted_at) WHERE successful = false;

-- iam
CREATE INDEX IF NOT EXISTS idx_user_roles_active
    ON iam.user_roles (user_id, school_id) WHERE is_active = true AND expires_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_roles_expires
    ON iam.user_roles (expires_at) WHERE expires_at IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_roles_user_active_covering
    ON iam.user_roles (user_id, school_id, academic_unit_id, role_id) WHERE is_active = true;

-- academic
CREATE INDEX IF NOT EXISTS idx_memberships_unit_role_active
    ON academic.memberships (academic_unit_id, role) WHERE is_active = true;

CREATE UNIQUE INDEX IF NOT EXISTS idx_academic_periods_active
    ON academic.academic_periods (school_id, academic_unit_id) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_concept_types_active
    ON academic.concept_types (is_active) WHERE is_active = true;

-- Una sola solicitud de ingreso pendiente por (usuario, colegio, unidad).
CREATE UNIQUE INDEX IF NOT EXISTS idx_join_requests_pending_unique
    ON academic.school_join_requests (user_id, school_id, academic_unit_id) WHERE status = 'pending';

-- Un solo componente auto_scored por (alumno, materia, periodo, intento de origen):
-- previene duplicar el grade_item derivado del mismo intento (N4 / ADR 0020). Los
-- componentes manuales (source_attempt_id NULL) quedan fuera del indice parcial.
CREATE UNIQUE INDEX IF NOT EXISTS uq_grade_item_attempt
    ON academic.grade_item (membership_id, subject_id, period_id, source_attempt_id) WHERE source_attempt_id IS NOT NULL;

-- content
CREATE INDEX IF NOT EXISTS idx_materials_status_active
    ON content.materials (status) WHERE deleted_at IS NULL;

-- assessment (N4 / ADR 0019: llaveado al modelo de sesion)
CREATE INDEX IF NOT EXISTS idx_attempt_completed
    ON assessment.assessment_attempt (assessment_id, percentage) WHERE status = 'completed';

CREATE INDEX IF NOT EXISTS idx_attempt_pending_review
    ON assessment.assessment_attempt (assessment_id) WHERE status = 'pending_review';

-- Un solo intento ACTIVO por (evaluacion, alumno): soporta "reusar intento
-- in_progress". Reemplaza los indices viejos por student_id/academic_unit_id
-- (modelo global muerto) que ya no existen.
CREATE UNIQUE INDEX IF NOT EXISTS idx_attempt_active_unique
    ON assessment.assessment_attempt (assessment_id, student_membership_id) WHERE status = 'in_progress';

-- evaluacion activa (soft delete)
CREATE INDEX IF NOT EXISTS idx_assessment_active
    ON assessment.assessment (id) WHERE deleted_at IS NULL;

-- ui_config
CREATE INDEX IF NOT EXISTS idx_screen_templates_active
    ON ui_config.screen_templates (is_active) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_screen_instances_active
    ON ui_config.screen_instances (is_active) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_screen_instances_handler_key
    ON ui_config.screen_instances (handler_key) WHERE handler_key IS NOT NULL;

-- audit
CREATE INDEX IF NOT EXISTS idx_audit_events_severity
    ON audit.events (severity, created_at DESC) WHERE severity != 'info';

-- notifications
CREATE INDEX IF NOT EXISTS idx_notif_user_unread
    ON notifications.notifications (user_id, created_at DESC) WHERE is_read = FALSE;

CREATE INDEX IF NOT EXISTS idx_device_tokens_user_active
    ON notifications.device_tokens (user_id) WHERE revoked_at IS NULL;
