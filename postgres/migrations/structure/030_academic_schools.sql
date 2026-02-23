-- ============================================================
-- 030: academic.schools
-- Schema: academic
-- Establecimientos educacionales
-- ============================================================

CREATE TABLE academic.schools (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    address text,
    city character varying(100),
    country character varying(100) DEFAULT 'Chile' NOT NULL,
    phone character varying(50),
    email character varying(255),
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    subscription_tier character varying(50) DEFAULT 'free' NOT NULL,
    max_teachers integer DEFAULT 10 NOT NULL,
    max_students integer DEFAULT 100 NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT schools_pkey PRIMARY KEY (id),
    CONSTRAINT schools_code_unique UNIQUE (code),
    CONSTRAINT schools_subscription_tier_check CHECK (subscription_tier IN ('free', 'basic', 'premium', 'enterprise'))
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.schools
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
