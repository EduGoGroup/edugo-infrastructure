-- ============================================================
-- 010: auth.users
-- Schema: auth
-- Usuarios del sistema (admin, docentes, estudiantes, apoderados)
-- ============================================================

CREATE TABLE auth.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_unique UNIQUE (email)
);

COMMENT ON TABLE auth.users IS 'Usuarios del sistema (admin, docentes, estudiantes, apoderados)';

CREATE TRIGGER set_updated_at BEFORE UPDATE ON auth.users
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
