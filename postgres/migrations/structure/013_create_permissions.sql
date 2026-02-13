-- ====================================================================
-- TABLA: permissions
-- DESCRIPCIÓN: Catálogo maestro de permisos del sistema RBAC
-- VERSIÓN: postgres/v0.15.0
-- FECHA: 2026-02-13
-- ====================================================================

-- Tipo ENUM para scope de permisos
CREATE TYPE permission_scope AS ENUM ('system', 'school', 'unit');

-- Tabla de permisos
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    scope permission_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Comentarios
COMMENT ON TABLE permissions IS 'Catálogo maestro de permisos del sistema RBAC';
COMMENT ON COLUMN permissions.name IS 'Nombre único del permiso en formato resource:action (ej: users:create)';
COMMENT ON COLUMN permissions.resource IS 'Recurso sobre el que aplica el permiso (users, materials, schools, etc.)';
COMMENT ON COLUMN permissions.action IS 'Acción que se puede realizar sobre el recurso (create, read, update, delete, etc.)';
