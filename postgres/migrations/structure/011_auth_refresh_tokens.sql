-- ============================================================
-- 011: auth.refresh_tokens
-- Schema: auth
-- Almacena refresh tokens JWT para gestión de sesiones
-- ============================================================

CREATE TABLE auth.refresh_tokens (
    id uuid NOT NULL,
    token_hash character varying(255) NOT NULL,
    user_id uuid NOT NULL,
    client_info jsonb,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    revoked_at timestamp with time zone,
    replaced_by uuid,
    CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash)
);

-- Intra-schema FK
ALTER TABLE auth.refresh_tokens ADD CONSTRAINT fk_refresh_tokens_replaced_by
    FOREIGN KEY (replaced_by) REFERENCES auth.refresh_tokens(id) ON DELETE SET NULL;

-- Indexes
CREATE INDEX idx_refresh_tokens_token_hash ON auth.refresh_tokens USING btree (token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON auth.refresh_tokens USING btree (user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON auth.refresh_tokens USING btree (expires_at);
CREATE INDEX idx_refresh_tokens_revoked_at ON auth.refresh_tokens USING btree (revoked_at) WHERE (revoked_at IS NOT NULL);

COMMENT ON TABLE auth.refresh_tokens IS 'Almacena refresh tokens JWT para gestión de sesiones';
