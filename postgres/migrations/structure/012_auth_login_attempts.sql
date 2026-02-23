-- ============================================================
-- 012: auth.login_attempts
-- Schema: auth
-- Registro de intentos de login para rate limiting y auditoría
-- ============================================================

CREATE SEQUENCE auth.login_attempts_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

CREATE TABLE auth.login_attempts (
    id integer NOT NULL DEFAULT nextval('auth.login_attempts_id_seq'::regclass),
    identifier character varying(255) NOT NULL,
    attempt_type character varying(50) NOT NULL,
    successful boolean DEFAULT false NOT NULL,
    user_agent text,
    ip_address character varying(45),
    attempted_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT login_attempts_pkey PRIMARY KEY (id),
    CONSTRAINT chk_attempt_type CHECK (attempt_type IN ('email', 'ip'))
);

ALTER SEQUENCE auth.login_attempts_id_seq OWNED BY auth.login_attempts.id;

-- Indexes
CREATE INDEX idx_login_attempts_identifier ON auth.login_attempts USING btree (identifier);
CREATE INDEX idx_login_attempts_attempted_at ON auth.login_attempts USING btree (attempted_at);
CREATE INDEX idx_login_attempts_identifier_attempted_at ON auth.login_attempts USING btree (identifier, attempted_at);
CREATE INDEX idx_login_attempts_rate_limit ON auth.login_attempts USING btree (identifier, successful, attempted_at) WHERE (successful = false);
CREATE INDEX idx_login_attempts_successful ON auth.login_attempts USING btree (successful);

COMMENT ON TABLE auth.login_attempts IS 'Registro de intentos de login para rate limiting y auditoría';
