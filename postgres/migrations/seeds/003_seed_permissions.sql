-- ====================================================================
-- SEEDS: Permisos predefinidos del sistema
-- VERSION: postgres/v0.17.0
-- ====================================================================

-- Permisos sobre USUARIOS
INSERT INTO permissions (name, display_name, description, resource_id, action, scope) VALUES
('users:create', 'Crear Usuarios', 'Crear nuevos usuarios en el sistema', (SELECT id FROM resources WHERE key = 'users'), 'create', 'system'),
('users:read', 'Ver Usuarios', 'Ver información de usuarios', (SELECT id FROM resources WHERE key = 'users'), 'read', 'school'),
('users:update', 'Editar Usuarios', 'Modificar datos de usuarios', (SELECT id FROM resources WHERE key = 'users'), 'update', 'school'),
('users:delete', 'Eliminar Usuarios', 'Eliminar usuarios del sistema', (SELECT id FROM resources WHERE key = 'users'), 'delete', 'system'),
('users:read:own', 'Ver Perfil Propio', 'Ver propio perfil de usuario', (SELECT id FROM resources WHERE key = 'users'), 'read:own', 'system'),
('users:update:own', 'Editar Perfil Propio', 'Modificar propio perfil', (SELECT id FROM resources WHERE key = 'users'), 'update:own', 'system'),

-- Permisos sobre ESCUELAS
('schools:create', 'Crear Escuelas', 'Crear nuevas instituciones educativas', (SELECT id FROM resources WHERE key = 'schools'), 'create', 'system'),
('schools:read', 'Ver Escuelas', 'Ver información de escuelas', (SELECT id FROM resources WHERE key = 'schools'), 'read', 'system'),
('schools:update', 'Editar Escuelas', 'Modificar datos de escuelas', (SELECT id FROM resources WHERE key = 'schools'), 'update', 'school'),
('schools:delete', 'Eliminar Escuelas', 'Eliminar escuelas del sistema', (SELECT id FROM resources WHERE key = 'schools'), 'delete', 'system'),
('schools:manage', 'Gestionar Escuela', 'Control total de la escuela', (SELECT id FROM resources WHERE key = 'schools'), 'manage', 'school'),

-- Permisos sobre UNIDADES ACADÉMICAS
('units:create', 'Crear Unidades', 'Crear unidades académicas (clases, grados)', (SELECT id FROM resources WHERE key = 'units'), 'create', 'school'),
('units:read', 'Ver Unidades', 'Ver unidades académicas', (SELECT id FROM resources WHERE key = 'units'), 'read', 'school'),
('units:update', 'Editar Unidades', 'Modificar unidades académicas', (SELECT id FROM resources WHERE key = 'units'), 'update', 'school'),
('units:delete', 'Eliminar Unidades', 'Eliminar unidades académicas', (SELECT id FROM resources WHERE key = 'units'), 'delete', 'school'),

-- Permisos sobre MATERIALES
('materials:create', 'Crear Materiales', 'Crear materiales educativos', (SELECT id FROM resources WHERE key = 'materials'), 'create', 'unit'),
('materials:read', 'Ver Materiales', 'Ver materiales educativos', (SELECT id FROM resources WHERE key = 'materials'), 'read', 'unit'),
('materials:update', 'Editar Materiales', 'Modificar materiales', (SELECT id FROM resources WHERE key = 'materials'), 'update', 'unit'),
('materials:delete', 'Eliminar Materiales', 'Eliminar materiales', (SELECT id FROM resources WHERE key = 'materials'), 'delete', 'unit'),
('materials:publish', 'Publicar Materiales', 'Publicar materiales para estudiantes', (SELECT id FROM resources WHERE key = 'materials'), 'publish', 'unit'),
('materials:download', 'Descargar Materiales', 'Descargar materiales educativos', (SELECT id FROM resources WHERE key = 'materials'), 'download', 'unit'),

-- Permisos sobre EVALUACIONES
('assessments:create', 'Crear Evaluaciones', 'Crear evaluaciones y exámenes', (SELECT id FROM resources WHERE key = 'assessments'), 'create', 'unit'),
('assessments:read', 'Ver Evaluaciones', 'Ver evaluaciones', (SELECT id FROM resources WHERE key = 'assessments'), 'read', 'unit'),
('assessments:update', 'Editar Evaluaciones', 'Modificar evaluaciones', (SELECT id FROM resources WHERE key = 'assessments'), 'update', 'unit'),
('assessments:delete', 'Eliminar Evaluaciones', 'Eliminar evaluaciones', (SELECT id FROM resources WHERE key = 'assessments'), 'delete', 'unit'),
('assessments:publish', 'Publicar Evaluaciones', 'Publicar evaluaciones para estudiantes', (SELECT id FROM resources WHERE key = 'assessments'), 'publish', 'unit'),
('assessments:grade', 'Calificar Evaluaciones', 'Calificar respuestas de estudiantes', (SELECT id FROM resources WHERE key = 'assessments'), 'grade', 'unit'),
('assessments:attempt', 'Rendir Evaluaciones', 'Intentar evaluaciones como estudiante', (SELECT id FROM resources WHERE key = 'assessments'), 'attempt', 'unit'),
('assessments:view_results', 'Ver Resultados', 'Ver resultados propios', (SELECT id FROM resources WHERE key = 'assessments'), 'view_results', 'unit'),

-- Permisos sobre PROGRESO
('progress:read', 'Ver Progreso', 'Ver progreso académico', (SELECT id FROM resources WHERE key = 'progress'), 'read', 'unit'),
('progress:update', 'Actualizar Progreso', 'Actualizar progreso de estudiantes', (SELECT id FROM resources WHERE key = 'progress'), 'update', 'unit'),
('progress:read:own', 'Ver Progreso Propio', 'Ver propio progreso', (SELECT id FROM resources WHERE key = 'progress'), 'read:own', 'unit'),

-- Permisos sobre ESTADÍSTICAS
('stats:global', 'Estadísticas Globales', 'Ver estadísticas de toda la plataforma', (SELECT id FROM resources WHERE key = 'stats'), 'global', 'system'),
('stats:school', 'Estadísticas de Escuela', 'Ver estadísticas de la institución', (SELECT id FROM resources WHERE key = 'stats'), 'school', 'school'),
('stats:unit', 'Estadísticas de Unidad', 'Ver estadísticas de la clase', (SELECT id FROM resources WHERE key = 'stats'), 'unit', 'unit'),

-- Permisos sobre GESTIÓN DE PERMISOS/RESOURCES
('permissions_mgmt:read', 'Ver Configuración de Permisos', 'Ver recursos y permisos del sistema', (SELECT id FROM resources WHERE key = 'permissions_mgmt'), 'read', 'system'),
('permissions_mgmt:update', 'Editar Configuración de Permisos', 'Modificar recursos y permisos del sistema', (SELECT id FROM resources WHERE key = 'permissions_mgmt'), 'update', 'system');
