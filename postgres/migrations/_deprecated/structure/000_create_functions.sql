-- ========================================
-- FUNCIONES BASE PARA POSTGRESQL
-- ========================================
-- Este archivo contiene funciones auxiliares utilizadas
-- por múltiples tablas en el sistema EduGo
-- Debe ejecutarse ANTES que cualquier otra migración

-- ========================================
-- FUNCIÓN: update_updated_at_column
-- ========================================
-- Actualiza automáticamente el campo updated_at
-- con la fecha/hora actual cuando se modifica un registro
--
-- Uso: Se asocia a un trigger BEFORE UPDATE en tablas
-- que tienen campo updated_at TIMESTAMP WITH TIME ZONE
--
-- Ejemplo:
--   CREATE TRIGGER set_updated_at_tablename
--     BEFORE UPDATE ON tablename
--     FOR EACH ROW
--     EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Comentario de la función
COMMENT ON FUNCTION update_updated_at_column() IS
'Trigger function para actualizar automáticamente el campo updated_at con la fecha/hora actual';

-- ====================================================================
-- FUNCIÓN: get_user_permissions
-- DESCRIPCIÓN: Obtiene permisos de un usuario en un contexto específico
-- VERSIÓN: postgres/v0.17.0
-- ====================================================================

CREATE OR REPLACE FUNCTION get_user_permissions(
    p_user_id UUID,
    p_school_id UUID DEFAULT NULL,
    p_unit_id UUID DEFAULT NULL
) RETURNS TABLE(permission_name VARCHAR, permission_scope permission_scope) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT p.name::VARCHAR, p.scope
    FROM user_roles ur
    JOIN roles ro ON ur.role_id = ro.id
    JOIN role_permissions rp ON ro.id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    JOIN resources r ON p.resource_id = r.id
    WHERE ur.user_id = p_user_id
      AND ur.is_active = true
      AND ro.is_active = true
      AND p.is_active = true
      AND r.is_active = true
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
      AND (
          -- Permisos a nivel sistema (sin contexto)
          (ur.school_id IS NULL AND p_school_id IS NULL)
          OR
          -- Permisos a nivel escuela (coincide school_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
          OR
          -- Permisos a nivel unidad (coincide school_id y unit_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          OR
          -- Permisos globales siempre aplican (super_admin)
          (ur.school_id IS NULL)
      );
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_permissions IS 'Obtiene lista de permisos de un usuario en un contexto específico';

-- ====================================================================
-- FUNCIÓN: user_has_permission
-- DESCRIPCIÓN: Verifica si un usuario tiene un permiso específico
-- VERSIÓN: postgres/v0.17.0
-- ====================================================================

CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_permission_name VARCHAR,
    p_school_id UUID DEFAULT NULL,
    p_unit_id UUID DEFAULT NULL
) RETURNS BOOLEAN AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM user_roles ur
        JOIN roles ro ON ur.role_id = ro.id
        JOIN role_permissions rp ON ro.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        JOIN resources r ON p.resource_id = r.id
        WHERE ur.user_id = p_user_id
          AND p.name = p_permission_name
          AND ur.is_active = true
          AND ro.is_active = true
          AND p.is_active = true
          AND r.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND (
              (ur.school_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          )
    ) INTO has_perm;

    RETURN has_perm;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION user_has_permission IS 'Verifica si un usuario tiene un permiso específico en un contexto dado';

-- ====================================================================
-- FUNCIÓN: get_user_resources
-- DESCRIPCIÓN: Obtiene los resources visibles en menu para un usuario
-- VERSIÓN: postgres/v0.17.0
-- ====================================================================

CREATE OR REPLACE FUNCTION get_user_resources(
    p_user_id UUID,
    p_school_id UUID DEFAULT NULL,
    p_unit_id UUID DEFAULT NULL
)
RETURNS TABLE(resource_key VARCHAR, resource_display_name VARCHAR, resource_icon VARCHAR, resource_scope permission_scope, parent_id UUID, sort_order INT)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE
    -- 1. Leaf resources the user has permission to access
    leaf_resources AS (
        SELECT DISTINCT r.id
        FROM resources r
        JOIN permissions p ON p.resource_id = r.id
        JOIN role_permissions rp ON rp.permission_id = p.id
        JOIN user_roles ur ON ur.role_id = rp.role_id
        JOIN roles ro ON ur.role_id = ro.id AND ro.is_active = true
        WHERE ur.user_id = p_user_id
          AND ur.is_active = true
          AND r.is_active = true
          AND r.is_menu_visible = true
          AND p.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND (
              (r.scope = 'system')
              OR (r.scope = 'school' AND p_school_id IS NOT NULL AND ur.school_id = p_school_id)
              OR (r.scope = 'unit' AND p_unit_id IS NOT NULL AND ur.academic_unit_id = p_unit_id)
          )
    ),
    -- 2. Recursively find all ancestors to build the full tree
    resource_tree AS (
        -- Base: leaf resources
        SELECT r2.id, r2.parent_id
        FROM resources r2
        WHERE r2.id IN (SELECT lr.id FROM leaf_resources lr)

        UNION

        -- Recursive: parent nodes
        SELECT r3.id, r3.parent_id
        FROM resources r3
        INNER JOIN resource_tree rt ON rt.parent_id = r3.id
        WHERE r3.is_active = true
          AND r3.is_menu_visible = true
    )
    SELECT DISTINCT r4.key::VARCHAR, r4.display_name::VARCHAR, r4.icon::VARCHAR, r4.scope, r4.parent_id, r4.sort_order::INT
    FROM resources r4
    INNER JOIN resource_tree rt2 ON rt2.id = r4.id
    ORDER BY r4.sort_order;
END;
$$;

COMMENT ON FUNCTION get_user_resources IS 'Obtiene los resources visibles en menu para un usuario según sus permisos';
