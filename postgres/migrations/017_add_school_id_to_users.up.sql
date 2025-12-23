-- Migración: 017_add_school_id_to_users
-- Descripción: Agrega school_id a la tabla users para asociar usuarios a su escuela principal
-- Nota: Es nullable para super_admin que no pertenecen a una escuela específica
-- Usado por: api-admin (JWT con school_id), api-mobile (multi-tenant)

ALTER TABLE users ADD COLUMN school_id UUID REFERENCES schools(id) ON DELETE SET NULL;

-- Índice para búsquedas por escuela
CREATE INDEX idx_users_school_id ON users(school_id);

COMMENT ON COLUMN users.school_id IS 'Escuela principal del usuario. NULL para super_admin sin escuela específica';
