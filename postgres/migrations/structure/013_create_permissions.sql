-- ====================================================================
-- TABLA: permissions
-- DESCRIPCION: Catalogo maestro de permisos del sistema RBAC
-- VERSION: postgres/v0.17.0
-- ====================================================================
-- NOTA: El tipo permission_scope se define en 000_base_types.sql
-- NOTA: Requiere tabla resources (011_create_resources.sql)

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    resource_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    scope permission_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE permissions IS 'Catalogo maestro de permisos del sistema RBAC';
COMMENT ON COLUMN permissions.name IS 'Nombre unico del permiso en formato resource:action (ej: users:create)';
COMMENT ON COLUMN permissions.resource_id IS 'FK al recurso sobre el que aplica el permiso';
COMMENT ON COLUMN permissions.action IS 'Accion que se puede realizar (create, read, update, delete, etc.)';
