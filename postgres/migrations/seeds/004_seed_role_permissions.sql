-- ====================================================================
-- SEEDS: Asignacion de permisos a roles
-- VERSION: postgres/v0.17.0
-- ====================================================================

-- SUPER_ADMIN: Todos los permisos
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'super_admin'),
    id
FROM permissions;

-- PLATFORM_ADMIN: Gestion de escuelas y usuarios + permisos mgmt
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'platform_admin'),
    id
FROM permissions
WHERE name IN (
    'users:create', 'users:read', 'users:update',
    'schools:create', 'schools:read', 'schools:update', 'schools:delete',
    'stats:global',
    'permissions_mgmt:read', 'permissions_mgmt:update'
);

-- SCHOOL_ADMIN: Control total de escuela
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'school_admin'),
    id
FROM permissions
WHERE name IN (
    'users:read', 'users:update',
    'schools:read', 'schools:update', 'schools:manage',
    'units:create', 'units:read', 'units:update', 'units:delete',
    'materials:read', 'materials:update', 'materials:delete',
    'assessments:read', 'assessments:update', 'assessments:delete',
    'progress:read', 'progress:update',
    'stats:school'
);

-- TEACHER: Gestion de clase
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'teacher'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'units:read',
    'materials:create', 'materials:read', 'materials:update', 'materials:publish', 'materials:download',
    'assessments:create', 'assessments:read', 'assessments:update', 'assessments:publish', 'assessments:grade',
    'progress:read', 'progress:update',
    'stats:unit'
);

-- STUDENT: Consumo de contenido
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'student'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'materials:read', 'materials:download',
    'assessments:read', 'assessments:attempt', 'assessments:view_results',
    'progress:read:own'
);

-- GUARDIAN: Ver progreso de estudiantes
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'guardian'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'materials:read',
    'assessments:view_results',
    'progress:read'
);
