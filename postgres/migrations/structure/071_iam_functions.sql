-- ============================================================
-- 071: IAM Functions
-- Schema: iam
-- Funciones PL/pgSQL para consultas RBAC con referencias schema-qualified
-- ============================================================

-- Obtiene lista de permisos de un usuario en un contexto específico
CREATE OR REPLACE FUNCTION iam.get_user_permissions(
    p_user_id uuid,
    p_school_id uuid DEFAULT NULL,
    p_unit_id uuid DEFAULT NULL
)
RETURNS TABLE(permission_name character varying, permission_scope iam.permission_scope)
LANGUAGE plpgsql
SET search_path = iam, auth, academic
AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT p.name::VARCHAR, p.scope
    FROM iam.user_roles ur
    JOIN iam.roles ro ON ur.role_id = ro.id
    JOIN iam.role_permissions rp ON ro.id = rp.role_id
    JOIN iam.permissions p ON rp.permission_id = p.id
    JOIN iam.resources r ON p.resource_id = r.id
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
$$;

COMMENT ON FUNCTION iam.get_user_permissions(uuid, uuid, uuid)
    IS 'Obtiene lista de permisos de un usuario en un contexto específico';

-- Obtiene los resources visibles en menú para un usuario según sus permisos
CREATE OR REPLACE FUNCTION iam.get_user_resources(
    p_user_id uuid,
    p_school_id uuid DEFAULT NULL,
    p_unit_id uuid DEFAULT NULL
)
RETURNS TABLE(
    resource_key character varying,
    resource_display_name character varying,
    resource_icon character varying,
    resource_scope iam.permission_scope,
    parent_id uuid,
    sort_order integer
)
LANGUAGE plpgsql
SET search_path = iam, auth, academic
AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE
    -- 1. Leaf resources the user has permission to access
    leaf_resources AS (
        SELECT DISTINCT r.id
        FROM iam.resources r
        JOIN iam.permissions p ON p.resource_id = r.id
        JOIN iam.role_permissions rp ON rp.permission_id = p.id
        JOIN iam.user_roles ur ON ur.role_id = rp.role_id
        JOIN iam.roles ro ON ur.role_id = ro.id AND ro.is_active = true
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
        FROM iam.resources r2
        WHERE r2.id IN (SELECT lr.id FROM leaf_resources lr)

        UNION

        -- Recursive: parent nodes
        SELECT r3.id, r3.parent_id
        FROM iam.resources r3
        INNER JOIN resource_tree rt ON rt.parent_id = r3.id
        WHERE r3.is_active = true
          AND r3.is_menu_visible = true
    )
    SELECT DISTINCT r4.key::VARCHAR, r4.display_name::VARCHAR, r4.icon::VARCHAR, r4.scope, r4.parent_id, r4.sort_order::INT
    FROM iam.resources r4
    INNER JOIN resource_tree rt2 ON rt2.id = r4.id
    ORDER BY r4.sort_order;
END;
$$;

COMMENT ON FUNCTION iam.get_user_resources(uuid, uuid, uuid)
    IS 'Obtiene los resources visibles en menu para un usuario según sus permisos';

-- Verifica si un usuario tiene un permiso específico en un contexto dado
CREATE OR REPLACE FUNCTION iam.user_has_permission(
    p_user_id uuid,
    p_permission_name character varying,
    p_school_id uuid DEFAULT NULL,
    p_unit_id uuid DEFAULT NULL
)
RETURNS boolean
LANGUAGE plpgsql
SET search_path = iam, auth, academic
AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM iam.user_roles ur
        JOIN iam.roles ro ON ur.role_id = ro.id
        JOIN iam.role_permissions rp ON ro.id = rp.role_id
        JOIN iam.permissions p ON rp.permission_id = p.id
        JOIN iam.resources r ON p.resource_id = r.id
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
$$;

COMMENT ON FUNCTION iam.user_has_permission(uuid, character varying, uuid, uuid)
    IS 'Verifica si un usuario tiene un permiso específico en un contexto dado';
