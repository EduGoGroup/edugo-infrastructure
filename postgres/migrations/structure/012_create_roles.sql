-- ====================================================================
-- TABLA: roles
-- DESCRIPCIÓN: Catálogo maestro de roles del sistema RBAC
-- VERSIÓN: postgres/v0.15.0
-- FECHA: 2026-02-13
-- ====================================================================

-- Tipo ENUM para scope de roles
CREATE TYPE role_scope AS ENUM ('system', 'school', 'unit');

-- Tabla de roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    scope role_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Comentarios
COMMENT ON TABLE roles IS 'Catálogo maestro de roles del sistema RBAC';
COMMENT ON COLUMN roles.name IS 'Nombre único del rol (snake_case)';
COMMENT ON COLUMN roles.display_name IS 'Nombre para mostrar en UI';
COMMENT ON COLUMN roles.scope IS 'Alcance del rol: system (global), school (institución), unit (clase/sección)';
