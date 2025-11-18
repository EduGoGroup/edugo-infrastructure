-- Tabla: users (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin, api-mobile, worker

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


COMMENT ON TABLE users IS 'Usuarios del sistema (admin, docentes, estudiantes, apoderados)';
COMMENT ON COLUMN users.role IS 'Roles: admin, teacher, student, guardian';
