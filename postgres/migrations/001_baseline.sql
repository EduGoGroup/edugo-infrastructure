-- ============================================================
-- MIGRACION 001: Baseline Schema (generado desde produccion)
-- Fecha: 2026-02-22
-- Fuente: Neon (ep-green-frost-ado4abbi-pooler)
-- NOTA: Este es el estado real de produccion como punto de partida.
--       Nuevas migraciones van en archivos 002_*.sql en adelante.
-- ============================================================

--
-- PostgreSQL database dump
--

\restrict 1VVQyfTmOcunn5xtJsSNMATFzJSOXGGHvpq3ytbuCfqVyvb88cKuC8SYajU97FN

-- Dumped from database version 17.8 (6108b59)
-- Dumped by pg_dump version 18.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', 'public', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- Name: ui_config; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA ui_config;


--
-- Name: permission_scope; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.permission_scope AS ENUM (
    'system',
    'school',
    'unit'
);


--
-- Name: TYPE permission_scope; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.permission_scope IS 'Define los alcances posibles de un permiso en el sistema RBAC';


--
-- Name: role_scope; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.role_scope AS ENUM (
    'system',
    'school',
    'unit'
);


--
-- Name: get_user_permissions(uuid, uuid, uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_user_permissions(p_user_id uuid, p_school_id uuid DEFAULT NULL::uuid, p_unit_id uuid DEFAULT NULL::uuid) RETURNS TABLE(permission_name character varying, permission_scope permission_scope)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT p.name::VARCHAR, p.scope
    FROM user_roles ur
    JOIN roles ro ON ur.role_id = ro.id
    JOIN role_permissions rp ON ro.id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    JOIN resources r ON p.resource_id = r.id
    WHERE ur.user_id = p_user_id
      AND ur.is_active = true
      AND ro.is_active = true
      AND p.is_active = true
      AND r.is_active = true
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
      AND (
          -- Permisos a nivel sistema (sin contexto)
          (ur.school_id IS NULL AND p_school_id IS NULL)
          OR
          -- Permisos a nivel escuela (coincide school_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
          OR
          -- Permisos a nivel unidad (coincide school_id y unit_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          OR
          -- Permisos globales siempre aplican (super_admin)
          (ur.school_id IS NULL)
      );
END;
$$;


--
-- Name: FUNCTION get_user_permissions(p_user_id uuid, p_school_id uuid, p_unit_id uuid); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.get_user_permissions(p_user_id uuid, p_school_id uuid, p_unit_id uuid) IS 'Obtiene lista de permisos de un usuario en un contexto específico';


--
-- Name: get_user_resources(uuid, uuid, uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_user_resources(p_user_id uuid, p_school_id uuid DEFAULT NULL::uuid, p_unit_id uuid DEFAULT NULL::uuid) RETURNS TABLE(resource_key character varying, resource_display_name character varying, resource_icon character varying, resource_scope permission_scope, parent_id uuid, sort_order integer)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE
    -- 1. Leaf resources the user has permission to access
    leaf_resources AS (
        SELECT DISTINCT r.id
        FROM resources r
        JOIN permissions p ON p.resource_id = r.id
        JOIN role_permissions rp ON rp.permission_id = p.id
        JOIN user_roles ur ON ur.role_id = rp.role_id
        JOIN roles ro ON ur.role_id = ro.id AND ro.is_active = true
        WHERE ur.user_id = p_user_id
          AND ur.is_active = true
          AND r.is_active = true
          AND r.is_menu_visible = true
          AND p.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND (
              (r.scope = 'system')
              OR (r.scope = 'school' AND p_school_id IS NOT NULL AND ur.school_id = p_school_id)
              OR (r.scope = 'unit' AND p_unit_id IS NOT NULL AND ur.academic_unit_id = p_unit_id)
          )
    ),
    -- 2. Recursively find all ancestors to build the full tree
    resource_tree AS (
        -- Base: leaf resources
        SELECT r2.id, r2.parent_id
        FROM resources r2
        WHERE r2.id IN (SELECT lr.id FROM leaf_resources lr)

        UNION

        -- Recursive: parent nodes
        SELECT r3.id, r3.parent_id
        FROM resources r3
        INNER JOIN resource_tree rt ON rt.parent_id = r3.id
        WHERE r3.is_active = true
          AND r3.is_menu_visible = true
    )
    SELECT DISTINCT r4.key::VARCHAR, r4.display_name::VARCHAR, r4.icon::VARCHAR, r4.scope, r4.parent_id, r4.sort_order::INT
    FROM resources r4
    INNER JOIN resource_tree rt2 ON rt2.id = r4.id
    ORDER BY r4.sort_order;
END;
$$;


--
-- Name: FUNCTION get_user_resources(p_user_id uuid, p_school_id uuid, p_unit_id uuid); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.get_user_resources(p_user_id uuid, p_school_id uuid, p_unit_id uuid) IS 'Obtiene los resources visibles en menu para un usuario según sus permisos';


--
-- Name: prevent_academic_unit_cycles(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.prevent_academic_unit_cycles() RETURNS trigger
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
        FROM academic_units
        WHERE id = current_parent_id;

        depth := depth + 1;
    END LOOP;

    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;

    RETURN NEW;
END;
$$;


--
-- Name: sync_questions_count(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.sync_questions_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.total_questions IS NOT NULL THEN
        NEW.questions_count := NEW.total_questions;
    ELSIF NEW.questions_count IS NOT NULL THEN
        NEW.total_questions := NEW.questions_count;
    ELSE
        NEW.total_questions := 0;
        NEW.questions_count := 0;
    END IF;
    RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


--
-- Name: FUNCTION update_updated_at_column(); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.update_updated_at_column() IS 'Trigger function para actualizar automáticamente el campo updated_at con la fecha/hora actual';


--
-- Name: user_has_permission(uuid, character varying, uuid, uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.user_has_permission(p_user_id uuid, p_permission_name character varying, p_school_id uuid DEFAULT NULL::uuid, p_unit_id uuid DEFAULT NULL::uuid) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM user_roles ur
        JOIN roles ro ON ur.role_id = ro.id
        JOIN role_permissions rp ON ro.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        JOIN resources r ON p.resource_id = r.id
        WHERE ur.user_id = p_user_id
          AND p.name = p_permission_name
          AND ur.is_active = true
          AND ro.is_active = true
          AND p.is_active = true
          AND r.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND (
              (ur.school_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          )
    ) INTO has_perm;

    RETURN has_perm;
END;
$$;


--
-- Name: FUNCTION user_has_permission(p_user_id uuid, p_permission_name character varying, p_school_id uuid, p_unit_id uuid); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.user_has_permission(p_user_id uuid, p_permission_name character varying, p_school_id uuid, p_unit_id uuid) IS 'Verifica si un usuario tiene un permiso específico en un contexto dado';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: academic_units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.academic_units (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    parent_unit_id uuid,
    school_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    type character varying(50) NOT NULL,
    description text,
    level character varying(50),
    academic_year integer DEFAULT 0,
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT academic_units_no_self_reference CHECK ((id <> parent_unit_id)),
    CONSTRAINT academic_units_type_check CHECK (((type)::text = ANY ((ARRAY['school'::character varying, 'grade'::character varying, 'class'::character varying, 'section'::character varying, 'club'::character varying, 'department'::character varying])::text[])))
);


--
-- Name: TABLE academic_units; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.academic_units IS 'Unidades académicas con soporte de jerarquía opcional';


--
-- Name: COLUMN academic_units.parent_unit_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.academic_units.parent_unit_id IS 'Unidad padre (jerarquía: Facultad → Departamento). NULL = raíz';


--
-- Name: COLUMN academic_units.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.academic_units.type IS 'Tipo: school, grade, class, section, club, department';


--
-- Name: COLUMN academic_units.description; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.academic_units.description IS 'Descripción de la unidad académica';


--
-- Name: COLUMN academic_units.academic_year; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.academic_units.academic_year IS 'Año académico. 0 = sin año específico (usado por api-admin)';


--
-- Name: COLUMN academic_units.metadata; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.academic_units.metadata IS 'Metadata extensible';


--
-- Name: assessment; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.assessment (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    material_id uuid NOT NULL,
    mongo_document_id character varying(24) NOT NULL,
    questions_count integer DEFAULT 0 NOT NULL,
    status character varying(50) DEFAULT 'generated'::character varying NOT NULL,
    title character varying(255),
    pass_threshold integer DEFAULT 70,
    max_attempts integer,
    time_limit_minutes integer,
    total_questions integer,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT assessment_pass_threshold_check CHECK (((pass_threshold >= 0) AND (pass_threshold <= 100))),
    CONSTRAINT assessment_status_check CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'generated'::character varying, 'published'::character varying, 'archived'::character varying, 'closed'::character varying])::text[])))
);


--
-- Name: TABLE assessment; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.assessment IS 'Assessments/Quizzes generados por IA (contenido en MongoDB)';


--
-- Name: COLUMN assessment.mongo_document_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.mongo_document_id IS 'ObjectId del documento en MongoDB material_assessment';


--
-- Name: COLUMN assessment.title; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.title IS 'Título del assessment (opcional si se usa metadata de MongoDB)';


--
-- Name: COLUMN assessment.pass_threshold; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.pass_threshold IS 'Porcentaje mínimo para aprobar (0-100)';


--
-- Name: COLUMN assessment.max_attempts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.max_attempts IS 'Máximo de intentos permitidos (NULL = ilimitado)';


--
-- Name: COLUMN assessment.time_limit_minutes; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.time_limit_minutes IS 'Límite de tiempo en minutos (NULL = sin límite)';


--
-- Name: COLUMN assessment.total_questions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment.total_questions IS 'Total de preguntas (sincronizado con questions_count)';


--
-- Name: assessment_attempt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.assessment_attempt (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    assessment_id uuid NOT NULL,
    student_id uuid NOT NULL,
    started_at timestamp with time zone DEFAULT now() NOT NULL,
    completed_at timestamp with time zone,
    score numeric(5,2),
    max_score numeric(5,2),
    percentage numeric(5,2),
    status character varying(50) DEFAULT 'in_progress'::character varying NOT NULL,
    time_spent_seconds integer,
    idempotency_key character varying(64),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT assessment_attempt_status_check CHECK (((status)::text = ANY ((ARRAY['in_progress'::character varying, 'completed'::character varying, 'abandoned'::character varying])::text[]))),
    CONSTRAINT assessment_attempt_time_spent_seconds_check CHECK (((time_spent_seconds IS NULL) OR ((time_spent_seconds > 0) AND (time_spent_seconds <= 7200)))),
    CONSTRAINT check_attempt_time_logical CHECK (((completed_at IS NULL) OR (completed_at > started_at)))
);


--
-- Name: TABLE assessment_attempt; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.assessment_attempt IS 'Intentos de estudiantes en assessments';


--
-- Name: COLUMN assessment_attempt.time_spent_seconds; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment_attempt.time_spent_seconds IS 'Tiempo total del intento en segundos (max 2 horas)';


--
-- Name: COLUMN assessment_attempt.idempotency_key; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment_attempt.idempotency_key IS 'Clave para prevenir intentos duplicados';


--
-- Name: CONSTRAINT check_attempt_time_logical ON assessment_attempt; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON CONSTRAINT check_attempt_time_logical ON public.assessment_attempt IS 'Validar que completed_at > started_at';


--
-- Name: assessment_attempt_answer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.assessment_attempt_answer (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    attempt_id uuid NOT NULL,
    question_index integer NOT NULL,
    student_answer text,
    is_correct boolean,
    points_earned numeric(5,2),
    max_points numeric(5,2),
    time_spent_seconds integer,
    answered_at timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT assessment_attempt_answer_time_spent_seconds_check CHECK ((time_spent_seconds >= 0))
);


--
-- Name: TABLE assessment_attempt_answer; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.assessment_attempt_answer IS 'Respuestas individuales de estudiantes por pregunta';


--
-- Name: COLUMN assessment_attempt_answer.question_index; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment_attempt_answer.question_index IS 'Índice de la pregunta (0-based). APIs mapean a question_id según necesidad.';


--
-- Name: COLUMN assessment_attempt_answer.student_answer; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment_attempt_answer.student_answer IS 'Respuesta del estudiante (TEXT flexible: JSON, string, etc). APIs mapean a selected_answer_id según necesidad.';


--
-- Name: COLUMN assessment_attempt_answer.time_spent_seconds; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.assessment_attempt_answer.time_spent_seconds IS 'Tiempo que tomó responder esta pregunta en segundos';


--
-- Name: login_attempts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.login_attempts (
    id integer NOT NULL,
    identifier character varying(255) NOT NULL,
    attempt_type character varying(50) NOT NULL,
    successful boolean DEFAULT false NOT NULL,
    user_agent text,
    ip_address character varying(45),
    attempted_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_attempt_type CHECK (((attempt_type)::text = ANY ((ARRAY['email'::character varying, 'ip'::character varying])::text[])))
);


--
-- Name: TABLE login_attempts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.login_attempts IS 'Registro de intentos de login para rate limiting y auditoría';


--
-- Name: COLUMN login_attempts.identifier; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.login_attempts.identifier IS 'Email o IP address dependiendo de attempt_type';


--
-- Name: COLUMN login_attempts.attempt_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.login_attempts.attempt_type IS 'Tipo de intento: email (por usuario) o ip (por dirección IP)';


--
-- Name: COLUMN login_attempts.successful; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.login_attempts.successful IS 'Indica si el intento de login fue exitoso';


--
-- Name: COLUMN login_attempts.user_agent; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.login_attempts.user_agent IS 'User agent del navegador/cliente';


--
-- Name: COLUMN login_attempts.ip_address; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.login_attempts.ip_address IS 'Dirección IP del cliente';


--
-- Name: login_attempts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.login_attempts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: login_attempts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.login_attempts_id_seq OWNED BY public.login_attempts.id;


--
-- Name: materials; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materials (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    uploaded_by_teacher_id uuid NOT NULL,
    academic_unit_id uuid,
    title character varying(255) NOT NULL,
    description text,
    subject character varying(100),
    grade character varying(50),
    file_url text NOT NULL,
    file_type character varying(100) NOT NULL,
    file_size_bytes bigint NOT NULL,
    status character varying(50) DEFAULT 'uploaded'::character varying NOT NULL,
    processing_started_at timestamp with time zone,
    processing_completed_at timestamp with time zone,
    is_public boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT materials_status_check CHECK (((status)::text = ANY ((ARRAY['uploaded'::character varying, 'processing'::character varying, 'ready'::character varying, 'failed'::character varying])::text[])))
);


--
-- Name: TABLE materials; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.materials IS 'Materiales educativos subidos por docentes';


--
-- Name: COLUMN materials.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.materials.status IS 'Estado: uploaded, processing, ready, failed';


--
-- Name: memberships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.memberships (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    school_id uuid NOT NULL,
    academic_unit_id uuid,
    role character varying(50) NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    enrolled_at timestamp with time zone DEFAULT now() NOT NULL,
    withdrawn_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT memberships_role_check CHECK (((role)::text = ANY ((ARRAY['teacher'::character varying, 'student'::character varying, 'guardian'::character varying, 'coordinator'::character varying, 'admin'::character varying, 'assistant'::character varying])::text[])))
);


--
-- Name: TABLE memberships; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.memberships IS 'Relación usuario-escuela-unidad académica';


--
-- Name: COLUMN memberships.role; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.memberships.role IS 'Rol: teacher, student, guardian, coordinator, admin, assistant';


--
-- Name: COLUMN memberships.metadata; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.memberships.metadata IS 'Metadata extensible: permisos específicos, configuración, historial';


--
-- Name: COLUMN memberships.enrolled_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.memberships.enrolled_at IS 'Fecha de inicio de membresía';


--
-- Name: COLUMN memberships.withdrawn_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.memberships.withdrawn_at IS 'Fecha de fin de membresía (NULL = activo)';


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permissions (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(150) NOT NULL,
    description text,
    resource_id uuid NOT NULL,
    action character varying(50) NOT NULL,
    scope permission_scope DEFAULT 'school'::permission_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_permission_name_format CHECK (((name)::text ~* '^[a-z_]+:[a-z_]+(:[a-z_]+)?$'::text))
);


--
-- Name: TABLE permissions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.permissions IS 'Catalogo maestro de permisos del sistema RBAC';


--
-- Name: COLUMN permissions.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.permissions.name IS 'Nombre unico del permiso en formato resource:action (ej: users:create)';


--
-- Name: COLUMN permissions.resource_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.permissions.resource_id IS 'FK al recurso sobre el que aplica el permiso';


--
-- Name: COLUMN permissions.action; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.permissions.action IS 'Accion que se puede realizar (create, read, update, delete, etc.)';


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    user_id uuid NOT NULL,
    client_info jsonb,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    revoked_at timestamp with time zone,
    replaced_by uuid
);


--
-- Name: TABLE refresh_tokens; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.refresh_tokens IS 'Almacena refresh tokens JWT para gestión de sesiones';


--
-- Name: COLUMN refresh_tokens.token_hash; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.token_hash IS 'Hash del refresh token (no se guarda el token en texto plano)';


--
-- Name: COLUMN refresh_tokens.client_info; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.client_info IS 'Información del cliente (navegador, IP, etc.)';


--
-- Name: COLUMN refresh_tokens.revoked_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.revoked_at IS 'Timestamp cuando el token fue revocado';


--
-- Name: COLUMN refresh_tokens.replaced_by; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.refresh_tokens.replaced_by IS 'ID del nuevo token que reemplazó a este (rotation)';


--
-- Name: resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resources (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    key character varying(50) NOT NULL,
    display_name character varying(150) NOT NULL,
    description text,
    icon character varying(100),
    parent_id uuid,
    sort_order integer DEFAULT 0 NOT NULL,
    is_menu_visible boolean DEFAULT true NOT NULL,
    scope permission_scope DEFAULT 'school'::permission_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: TABLE resources; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.resources IS 'Catalogo de recursos/modulos del sistema para RBAC y generacion de menu';


--
-- Name: COLUMN resources.key; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resources.key IS 'Identificador unico del recurso (ej: users, schools, materials)';


--
-- Name: COLUMN resources.icon; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resources.icon IS 'Nombre del icono para UI (ej: users, school, book)';


--
-- Name: COLUMN resources.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resources.parent_id IS 'FK a resources.id para jerarquia de menu';


--
-- Name: COLUMN resources.sort_order; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resources.sort_order IS 'Orden de aparicion dentro de su nivel de menu';


--
-- Name: COLUMN resources.is_menu_visible; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.resources.is_menu_visible IS 'Si el recurso aparece como item de menu';


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_permissions (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: TABLE role_permissions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.role_permissions IS 'Relación N:N entre roles y permisos (RBAC)';


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(50) NOT NULL,
    display_name character varying(100) NOT NULL,
    description text,
    scope role_scope DEFAULT 'school'::role_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: TABLE roles; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.roles IS 'Catálogo maestro de roles del sistema RBAC';


--
-- Name: COLUMN roles.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.roles.name IS 'Nombre único del rol (snake_case)';


--
-- Name: COLUMN roles.display_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.roles.display_name IS 'Nombre para mostrar en UI';


--
-- Name: COLUMN roles.scope; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.roles.scope IS 'Alcance del rol: system (global), school (institución), unit (clase/sección)';


--
-- Name: schools; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schools (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    address text,
    city character varying(100),
    country character varying(100) DEFAULT 'Chile'::character varying NOT NULL,
    phone character varying(50),
    email character varying(255),
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    subscription_tier character varying(50) DEFAULT 'free'::character varying NOT NULL,
    max_teachers integer DEFAULT 10 NOT NULL,
    max_students integer DEFAULT 100 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT schools_subscription_tier_check CHECK (((subscription_tier)::text = ANY ((ARRAY['free'::character varying, 'basic'::character varying, 'premium'::character varying, 'enterprise'::character varying])::text[])))
);


--
-- Name: TABLE schools; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.schools IS 'Escuelas/Instituciones educativas';


--
-- Name: COLUMN schools.metadata; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.schools.metadata IS 'Metadata extensible: logo, configuración institucional, etc.';


--
-- Name: COLUMN schools.subscription_tier; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.schools.subscription_tier IS 'Nivel de subscripción: free, basic, premium, enterprise';


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    school_id uuid,
    academic_unit_id uuid,
    is_active boolean DEFAULT true NOT NULL,
    granted_by uuid,
    granted_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_user_roles_unit_requires_school CHECK (((academic_unit_id IS NULL) OR (school_id IS NOT NULL)))
);


--
-- Name: TABLE user_roles; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.user_roles IS 'Asignación de roles a usuarios en contextos específicos (RBAC)';


--
-- Name: COLUMN user_roles.school_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_roles.school_id IS 'Escuela en la que aplica el rol. NULL = rol a nivel sistema';


--
-- Name: COLUMN user_roles.academic_unit_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_roles.academic_unit_id IS 'Unidad académica en la que aplica el rol. NULL = rol a nivel escuela';


--
-- Name: COLUMN user_roles.granted_by; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_roles.granted_by IS 'Usuario que otorgó el rol (auditoría)';


--
-- Name: COLUMN user_roles.expires_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_roles.expires_at IS 'Fecha de expiración del rol. NULL = no expira';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.users IS 'Usuarios del sistema (admin, docentes, estudiantes, apoderados)';


--
-- Name: v_academic_unit_tree; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_academic_unit_tree AS
 WITH RECURSIVE unit_hierarchy AS (
         SELECT academic_units.id,
            academic_units.parent_unit_id,
            academic_units.school_id,
            academic_units.name,
            academic_units.code,
            academic_units.type,
            academic_units.level,
            academic_units.academic_year,
            1 AS depth,
            ARRAY[academic_units.id] AS path,
            (academic_units.name)::text AS full_path
           FROM academic_units
          WHERE ((academic_units.parent_unit_id IS NULL) AND (academic_units.deleted_at IS NULL))
        UNION ALL
         SELECT au.id,
            au.parent_unit_id,
            au.school_id,
            au.name,
            au.code,
            au.type,
            au.level,
            au.academic_year,
            (uh_1.depth + 1),
            (uh_1.path || au.id),
            ((uh_1.full_path || ' > '::text) || (au.name)::text) AS text
           FROM (academic_units au
             JOIN unit_hierarchy uh_1 ON ((au.parent_unit_id = uh_1.id)))
          WHERE (au.deleted_at IS NULL)
        )
 SELECT uh.id,
    uh.parent_unit_id,
    uh.school_id,
    uh.name,
    uh.code,
    uh.type,
    uh.level,
    uh.academic_year,
    uh.depth,
    uh.path,
    uh.full_path,
    s.name AS school_name,
    s.code AS school_code
   FROM (unit_hierarchy uh
     LEFT JOIN schools s ON ((uh.school_id = s.id)))
  ORDER BY uh.school_id, uh.path;


--
-- Name: VIEW v_academic_unit_tree; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON VIEW public.v_academic_unit_tree IS 'Vista con árbol jerárquico completo de unidades académicas';


--
-- Name: resource_screens; Type: TABLE; Schema: ui_config; Owner: -
--

CREATE TABLE ui_config.resource_screens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    resource_id uuid NOT NULL,
    resource_key character varying(100) NOT NULL,
    screen_key character varying(100) NOT NULL,
    screen_type character varying(50) NOT NULL,
    is_default boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: screen_instances; Type: TABLE; Schema: ui_config; Owner: -
--

CREATE TABLE ui_config.screen_instances (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    screen_key character varying(100) NOT NULL,
    template_id uuid NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    slot_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    actions jsonb DEFAULT '[]'::jsonb NOT NULL,
    data_endpoint character varying(500),
    data_config jsonb DEFAULT '{}'::jsonb,
    scope character varying(20) DEFAULT 'school'::character varying,
    required_permission character varying(100),
    handler_key character varying(100) DEFAULT NULL::character varying,
    is_active boolean DEFAULT true,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: screen_templates; Type: TABLE; Schema: ui_config; Owner: -
--

CREATE TABLE ui_config.screen_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pattern character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    version integer DEFAULT 1 NOT NULL,
    definition jsonb NOT NULL,
    is_active boolean DEFAULT true,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: screen_user_preferences; Type: TABLE; Schema: ui_config; Owner: -
--

CREATE TABLE ui_config.screen_user_preferences (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    screen_instance_id uuid NOT NULL,
    user_id uuid NOT NULL,
    preferences jsonb DEFAULT '{}'::jsonb NOT NULL,
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: login_attempts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login_attempts ALTER COLUMN id SET DEFAULT nextval('login_attempts_id_seq'::regclass);


--
-- Name: academic_units academic_units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.academic_units
    ADD CONSTRAINT academic_units_pkey PRIMARY KEY (id);


--
-- Name: academic_units academic_units_unique_code; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.academic_units
    ADD CONSTRAINT academic_units_unique_code UNIQUE (school_id, code, academic_year);


--
-- Name: assessment_attempt_answer assessment_attempt_answer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt_answer
    ADD CONSTRAINT assessment_attempt_answer_pkey PRIMARY KEY (id);


--
-- Name: assessment_attempt_answer assessment_attempt_answer_unique_question; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt_answer
    ADD CONSTRAINT assessment_attempt_answer_unique_question UNIQUE (attempt_id, question_index);


--
-- Name: assessment_attempt assessment_attempt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt
    ADD CONSTRAINT assessment_attempt_pkey PRIMARY KEY (id);


--
-- Name: assessment assessment_mongo_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment
    ADD CONSTRAINT assessment_mongo_unique UNIQUE (mongo_document_id);


--
-- Name: assessment assessment_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment
    ADD CONSTRAINT assessment_pkey PRIMARY KEY (id);


--
-- Name: login_attempts login_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.login_attempts
    ADD CONSTRAINT login_attempts_pkey PRIMARY KEY (id);


--
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- Name: memberships memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_pkey PRIMARY KEY (id);


--
-- Name: memberships memberships_unique_membership; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_unique_membership UNIQUE (user_id, school_id, academic_unit_id, role);


--
-- Name: permissions permissions_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_key UNIQUE (name);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash);


--
-- Name: resources resources_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_key_key UNIQUE (key);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schools schools_code_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schools
    ADD CONSTRAINT schools_code_unique UNIQUE (code);


--
-- Name: schools schools_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schools
    ADD CONSTRAINT schools_pkey PRIMARY KEY (id);


--
-- Name: assessment_attempt unique_idempotency_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt
    ADD CONSTRAINT unique_idempotency_key UNIQUE (idempotency_key);


--
-- Name: permissions uq_permissions_resource_action; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT uq_permissions_resource_action UNIQUE (resource_id, action);


--
-- Name: role_permissions uq_role_permission; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id);


--
-- Name: user_roles uq_user_role_context; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT uq_user_role_context UNIQUE (user_id, role_id, school_id, academic_unit_id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_email_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_unique UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: resource_screens resource_screens_pkey; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.resource_screens
    ADD CONSTRAINT resource_screens_pkey PRIMARY KEY (id);


--
-- Name: resource_screens resource_screens_resource_id_screen_type_key; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.resource_screens
    ADD CONSTRAINT resource_screens_resource_id_screen_type_key UNIQUE (resource_id, screen_type);


--
-- Name: screen_instances screen_instances_pkey; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_instances
    ADD CONSTRAINT screen_instances_pkey PRIMARY KEY (id);


--
-- Name: screen_instances screen_instances_screen_key_key; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_instances
    ADD CONSTRAINT screen_instances_screen_key_key UNIQUE (screen_key);


--
-- Name: screen_templates screen_templates_name_version_key; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_templates
    ADD CONSTRAINT screen_templates_name_version_key UNIQUE (name, version);


--
-- Name: screen_templates screen_templates_pkey; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_templates
    ADD CONSTRAINT screen_templates_pkey PRIMARY KEY (id);


--
-- Name: screen_user_preferences screen_user_preferences_pkey; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_user_preferences
    ADD CONSTRAINT screen_user_preferences_pkey PRIMARY KEY (id);


--
-- Name: screen_user_preferences screen_user_preferences_screen_instance_id_user_id_key; Type: CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_user_preferences
    ADD CONSTRAINT screen_user_preferences_screen_instance_id_user_id_key UNIQUE (screen_instance_id, user_id);


--
-- Name: idx_login_attempts_attempted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_login_attempts_attempted_at ON public.login_attempts USING btree (attempted_at);


--
-- Name: idx_login_attempts_identifier; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_login_attempts_identifier ON public.login_attempts USING btree (identifier);


--
-- Name: idx_login_attempts_identifier_attempted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_login_attempts_identifier_attempted_at ON public.login_attempts USING btree (identifier, attempted_at);


--
-- Name: idx_login_attempts_rate_limit; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_login_attempts_rate_limit ON public.login_attempts USING btree (identifier, successful, attempted_at) WHERE (successful = false);


--
-- Name: idx_login_attempts_successful; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_login_attempts_successful ON public.login_attempts USING btree (successful);


--
-- Name: idx_permissions_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_active ON public.permissions USING btree (is_active);


--
-- Name: idx_permissions_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_name ON public.permissions USING btree (name);


--
-- Name: idx_permissions_resource; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_resource ON public.permissions USING btree (resource_id);


--
-- Name: idx_permissions_scope; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permissions_scope ON public.permissions USING btree (scope);


--
-- Name: idx_refresh_tokens_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--
-- Name: idx_refresh_tokens_revoked_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_revoked_at ON public.refresh_tokens USING btree (revoked_at) WHERE (revoked_at IS NOT NULL);


--
-- Name: idx_refresh_tokens_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_resources_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_active ON public.resources USING btree (is_active);


--
-- Name: idx_resources_key; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_key ON public.resources USING btree (key);


--
-- Name: idx_resources_menu_visible; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_menu_visible ON public.resources USING btree (is_menu_visible);


--
-- Name: idx_resources_parent; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_parent ON public.resources USING btree (parent_id);


--
-- Name: idx_resources_sort; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_resources_sort ON public.resources USING btree (sort_order);


--
-- Name: idx_role_permissions_permission; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_permission ON public.role_permissions USING btree (permission_id);


--
-- Name: idx_role_permissions_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_role ON public.role_permissions USING btree (role_id);


--
-- Name: idx_roles_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_active ON public.roles USING btree (is_active);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_roles_scope; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_scope ON public.roles USING btree (scope);


--
-- Name: idx_user_roles_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_active ON public.user_roles USING btree (is_active);


--
-- Name: idx_user_roles_context; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_context ON public.user_roles USING btree (user_id, school_id, academic_unit_id);


--
-- Name: idx_user_roles_expires; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_expires ON public.user_roles USING btree (expires_at) WHERE (expires_at IS NOT NULL);


--
-- Name: idx_user_roles_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_role ON public.user_roles USING btree (role_id);


--
-- Name: idx_user_roles_school; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_school ON public.user_roles USING btree (school_id);


--
-- Name: idx_user_roles_unit; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_unit ON public.user_roles USING btree (academic_unit_id);


--
-- Name: idx_user_roles_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_user ON public.user_roles USING btree (user_id);


--
-- Name: idx_user_roles_user_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_user_active ON public.user_roles USING btree (user_id, is_active);


--
-- Name: idx_resource_screens_resource; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_resource_screens_resource ON ui_config.resource_screens USING btree (resource_id);


--
-- Name: idx_resource_screens_resource_key; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_resource_screens_resource_key ON ui_config.resource_screens USING btree (resource_key);


--
-- Name: idx_resource_screens_screen_key; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_resource_screens_screen_key ON ui_config.resource_screens USING btree (screen_key);


--
-- Name: idx_screen_instances_active; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_instances_active ON ui_config.screen_instances USING btree (is_active) WHERE (is_active = true);


--
-- Name: idx_screen_instances_handler_key; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_instances_handler_key ON ui_config.screen_instances USING btree (handler_key) WHERE (handler_key IS NOT NULL);


--
-- Name: idx_screen_instances_scope; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_instances_scope ON ui_config.screen_instances USING btree (scope);


--
-- Name: idx_screen_instances_slot_data; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_instances_slot_data ON ui_config.screen_instances USING gin (slot_data);


--
-- Name: idx_screen_instances_template; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_instances_template ON ui_config.screen_instances USING btree (template_id);


--
-- Name: idx_screen_templates_active; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_templates_active ON ui_config.screen_templates USING btree (is_active) WHERE (is_active = true);


--
-- Name: idx_screen_templates_definition; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_templates_definition ON ui_config.screen_templates USING gin (definition);


--
-- Name: idx_screen_templates_pattern; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_templates_pattern ON ui_config.screen_templates USING btree (pattern);


--
-- Name: idx_screen_user_prefs_screen; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_user_prefs_screen ON ui_config.screen_user_preferences USING btree (screen_instance_id);


--
-- Name: idx_screen_user_prefs_user; Type: INDEX; Schema: ui_config; Owner: -
--

CREATE INDEX idx_screen_user_prefs_user ON ui_config.screen_user_preferences USING btree (user_id);


--
-- Name: permissions set_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at BEFORE UPDATE ON public.permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: resources set_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at BEFORE UPDATE ON public.resources FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: roles set_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at BEFORE UPDATE ON public.roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: user_roles set_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at BEFORE UPDATE ON public.user_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: assessment trg_sync_questions_count; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_sync_questions_count BEFORE INSERT OR UPDATE ON public.assessment FOR EACH ROW EXECUTE FUNCTION sync_questions_count();


--
-- Name: TRIGGER trg_sync_questions_count ON assessment; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TRIGGER trg_sync_questions_count ON public.assessment IS 'Mantiene sincronizado questions_count y total_questions durante transición';


--
-- Name: academic_units trigger_prevent_academic_unit_cycles; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_prevent_academic_unit_cycles BEFORE INSERT OR UPDATE OF parent_unit_id ON public.academic_units FOR EACH ROW EXECUTE FUNCTION prevent_academic_unit_cycles();


--
-- Name: resource_screens update_resource_screens_updated_at; Type: TRIGGER; Schema: ui_config; Owner: -
--

CREATE TRIGGER update_resource_screens_updated_at BEFORE UPDATE ON ui_config.resource_screens FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: screen_instances update_screen_instances_updated_at; Type: TRIGGER; Schema: ui_config; Owner: -
--

CREATE TRIGGER update_screen_instances_updated_at BEFORE UPDATE ON ui_config.screen_instances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: screen_templates update_screen_templates_updated_at; Type: TRIGGER; Schema: ui_config; Owner: -
--

CREATE TRIGGER update_screen_templates_updated_at BEFORE UPDATE ON ui_config.screen_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: screen_user_preferences update_screen_user_prefs_updated_at; Type: TRIGGER; Schema: ui_config; Owner: -
--

CREATE TRIGGER update_screen_user_prefs_updated_at BEFORE UPDATE ON ui_config.screen_user_preferences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


--
-- Name: academic_units academic_units_parent_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.academic_units
    ADD CONSTRAINT academic_units_parent_fkey FOREIGN KEY (parent_unit_id) REFERENCES academic_units(id) ON DELETE SET NULL;


--
-- Name: academic_units academic_units_school_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.academic_units
    ADD CONSTRAINT academic_units_school_fkey FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;


--
-- Name: assessment_attempt_answer assessment_attempt_answer_attempt_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt_answer
    ADD CONSTRAINT assessment_attempt_answer_attempt_fkey FOREIGN KEY (attempt_id) REFERENCES assessment_attempt(id) ON DELETE CASCADE;


--
-- Name: assessment_attempt assessment_attempt_assessment_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt
    ADD CONSTRAINT assessment_attempt_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE;


--
-- Name: assessment_attempt assessment_attempt_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment_attempt
    ADD CONSTRAINT assessment_attempt_student_fkey FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE;


--
-- Name: assessment assessment_material_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assessment
    ADD CONSTRAINT assessment_material_fkey FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;


--
-- Name: permissions fk_permissions_resource; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT fk_permissions_resource FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE RESTRICT;


--
-- Name: refresh_tokens fk_refresh_tokens_replaced_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_replaced_by FOREIGN KEY (replaced_by) REFERENCES refresh_tokens(id) ON DELETE SET NULL;


--
-- Name: refresh_tokens fk_refresh_tokens_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


--
-- Name: resources fk_resources_parent; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT fk_resources_parent FOREIGN KEY (parent_id) REFERENCES resources(id) ON DELETE SET NULL;


--
-- Name: role_permissions fk_role_permissions_permission; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;


--
-- Name: role_permissions fk_role_permissions_role; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;


--
-- Name: user_roles fk_user_roles_granted_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;


--
-- Name: user_roles fk_user_roles_role; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;


--
-- Name: user_roles fk_user_roles_school; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_school FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;


--
-- Name: user_roles fk_user_roles_unit; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_unit FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE CASCADE;


--
-- Name: user_roles fk_user_roles_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


--
-- Name: materials materials_school_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_school_fkey FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;


--
-- Name: materials materials_teacher_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_teacher_fkey FOREIGN KEY (uploaded_by_teacher_id) REFERENCES users(id) ON DELETE RESTRICT;


--
-- Name: materials materials_unit_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE SET NULL;


--
-- Name: memberships memberships_school_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_school_fkey FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;


--
-- Name: memberships memberships_unit_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE CASCADE;


--
-- Name: memberships memberships_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_user_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


--
-- Name: resource_screens fk_resource_screens_resource; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.resource_screens
    ADD CONSTRAINT fk_resource_screens_resource FOREIGN KEY (resource_id) REFERENCES resources(id);


--
-- Name: resource_screens fk_resource_screens_screen_key; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.resource_screens
    ADD CONSTRAINT fk_resource_screens_screen_key FOREIGN KEY (screen_key) REFERENCES ui_config.screen_instances(screen_key);


--
-- Name: screen_instances fk_screen_instances_created_by; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_instances
    ADD CONSTRAINT fk_screen_instances_created_by FOREIGN KEY (created_by) REFERENCES users(id);


--
-- Name: screen_instances fk_screen_instances_template; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_instances
    ADD CONSTRAINT fk_screen_instances_template FOREIGN KEY (template_id) REFERENCES ui_config.screen_templates(id);


--
-- Name: screen_templates fk_screen_templates_created_by; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_templates
    ADD CONSTRAINT fk_screen_templates_created_by FOREIGN KEY (created_by) REFERENCES users(id);


--
-- Name: screen_user_preferences fk_screen_user_prefs_instance; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_user_preferences
    ADD CONSTRAINT fk_screen_user_prefs_instance FOREIGN KEY (screen_instance_id) REFERENCES ui_config.screen_instances(id);


--
-- Name: screen_user_preferences fk_screen_user_prefs_user; Type: FK CONSTRAINT; Schema: ui_config; Owner: -
--

ALTER TABLE ONLY ui_config.screen_user_preferences
    ADD CONSTRAINT fk_screen_user_prefs_user FOREIGN KEY (user_id) REFERENCES users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict 1VVQyfTmOcunn5xtJsSNMATFzJSOXGGHvpq3ytbuCfqVyvb88cKuC8SYajU97FN

