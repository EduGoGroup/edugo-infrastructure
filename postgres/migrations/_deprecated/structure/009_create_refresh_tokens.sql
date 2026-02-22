-- Tabla para almacenar refresh tokens JWT
-- Permite invalidar tokens y rastrear sesiones de usuarios

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    client_info JSONB,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE,
    replaced_by UUID,
    
    CONSTRAINT fk_refresh_tokens_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_refresh_tokens_replaced_by 
        FOREIGN KEY (replaced_by) 
        REFERENCES refresh_tokens(id) 
        ON DELETE SET NULL
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at) WHERE revoked_at IS NOT NULL;

-- Comentarios
COMMENT ON TABLE refresh_tokens IS 'Almacena refresh tokens JWT para gestión de sesiones';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'Hash del refresh token (no se guarda el token en texto plano)';
COMMENT ON COLUMN refresh_tokens.client_info IS 'Información del cliente (navegador, IP, etc.)';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp cuando el token fue revocado';
COMMENT ON COLUMN refresh_tokens.replaced_by IS 'ID del nuevo token que reemplazó a este (rotation)';
