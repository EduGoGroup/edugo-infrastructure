-- Tabla para rastrear intentos de login (exitosos y fallidos)
-- Implementa rate limiting y seguridad contra ataques de fuerza bruta

CREATE TABLE IF NOT EXISTS login_attempts (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL,
    attempt_type VARCHAR(50) NOT NULL,
    successful BOOLEAN NOT NULL DEFAULT false,
    user_agent TEXT,
    ip_address VARCHAR(45),
    attempted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_attempt_type 
        CHECK (attempt_type IN ('email', 'ip'))
);

-- Índices para rate limiting eficiente
CREATE INDEX idx_login_attempts_identifier ON login_attempts(identifier);
CREATE INDEX idx_login_attempts_attempted_at ON login_attempts(attempted_at);
CREATE INDEX idx_login_attempts_identifier_attempted_at ON login_attempts(identifier, attempted_at);
CREATE INDEX idx_login_attempts_successful ON login_attempts(successful);

-- Índice para consultas de rate limiting
CREATE INDEX idx_login_attempts_rate_limit 
    ON login_attempts(identifier, successful, attempted_at) 
    WHERE successful = false;

-- Comentarios
COMMENT ON TABLE login_attempts IS 'Registro de intentos de login para rate limiting y auditoría';
COMMENT ON COLUMN login_attempts.identifier IS 'Email o IP address dependiendo de attempt_type';
COMMENT ON COLUMN login_attempts.attempt_type IS 'Tipo de intento: email (por usuario) o ip (por dirección IP)';
COMMENT ON COLUMN login_attempts.successful IS 'Indica si el intento de login fue exitoso';
COMMENT ON COLUMN login_attempts.user_agent IS 'User agent del navegador/cliente';
COMMENT ON COLUMN login_attempts.ip_address IS 'Dirección IP del cliente';
