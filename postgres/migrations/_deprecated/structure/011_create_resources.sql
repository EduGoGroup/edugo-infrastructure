-- ====================================================================
-- TABLA: resources
-- DESCRIPCION: Catalogo de recursos/modulos del sistema para RBAC y menu
-- VERSION: postgres/v0.17.0
-- ====================================================================

CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    parent_id UUID,
    sort_order INT DEFAULT 0 NOT NULL,
    is_menu_visible BOOLEAN DEFAULT true NOT NULL,
    scope permission_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE resources IS 'Catalogo de recursos/modulos del sistema para RBAC y generacion de menu';
COMMENT ON COLUMN resources.key IS 'Identificador unico del recurso (ej: users, schools, materials)';
COMMENT ON COLUMN resources.icon IS 'Nombre del icono para UI (ej: users, school, book)';
COMMENT ON COLUMN resources.parent_id IS 'FK a resources.id para jerarquia de menu';
COMMENT ON COLUMN resources.sort_order IS 'Orden de aparicion dentro de su nivel de menu';
COMMENT ON COLUMN resources.is_menu_visible IS 'Si el recurso aparece como item de menu';
