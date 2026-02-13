-- ====================================================================
-- SEEDS: Permisos predefinidos del sistema
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Permisos sobre USUARIOS
INSERT INTO permissions (name, display_name, description, resource, action, scope) VALUES
('users:create', 'Crear Usuarios', 'Crear nuevos usuarios en el sistema', 'users', 'create', 'system'),
('users:read', 'Ver Usuarios', 'Ver información de usuarios', 'users', 'read', 'school'),
('users:update', 'Editar Usuarios', 'Modificar datos de usuarios', 'users', 'update', 'school'),
('users:delete', 'Eliminar Usuarios', 'Eliminar usuarios del sistema', 'users', 'delete', 'system'),
('users:read:own', 'Ver Perfil Propio', 'Ver propio perfil de usuario', 'users', 'read:own', 'system'),
('users:update:own', 'Editar Perfil Propio', 'Modificar propio perfil', 'users', 'update:own', 'system'),

-- Permisos sobre ESCUELAS
('schools:create', 'Crear Escuelas', 'Crear nuevas instituciones educativas', 'schools', 'create', 'system'),
('schools:read', 'Ver Escuelas', 'Ver información de escuelas', 'schools', 'read', 'system'),
('schools:update', 'Editar Escuelas', 'Modificar datos de escuelas', 'schools', 'update', 'school'),
('schools:delete', 'Eliminar Escuelas', 'Eliminar escuelas del sistema', 'schools', 'delete', 'system'),
('schools:manage', 'Gestionar Escuela', 'Control total de la escuela', 'schools', 'manage', 'school'),

-- Permisos sobre UNIDADES ACADÉMICAS
('units:create', 'Crear Unidades', 'Crear unidades académicas (clases, grados)', 'units', 'create', 'school'),
('units:read', 'Ver Unidades', 'Ver unidades académicas', 'units', 'read', 'school'),
('units:update', 'Editar Unidades', 'Modificar unidades académicas', 'units', 'update', 'school'),
('units:delete', 'Eliminar Unidades', 'Eliminar unidades académicas', 'units', 'delete', 'school'),

-- Permisos sobre MATERIALES
('materials:create', 'Crear Materiales', 'Crear materiales educativos', 'materials', 'create', 'unit'),
('materials:read', 'Ver Materiales', 'Ver materiales educativos', 'materials', 'read', 'unit'),
('materials:update', 'Editar Materiales', 'Modificar materiales', 'materials', 'update', 'unit'),
('materials:delete', 'Eliminar Materiales', 'Eliminar materiales', 'materials', 'delete', 'unit'),
('materials:publish', 'Publicar Materiales', 'Publicar materiales para estudiantes', 'materials', 'publish', 'unit'),
('materials:download', 'Descargar Materiales', 'Descargar materiales educativos', 'materials', 'download', 'unit'),

-- Permisos sobre EVALUACIONES
('assessments:create', 'Crear Evaluaciones', 'Crear evaluaciones y exámenes', 'assessments', 'create', 'unit'),
('assessments:read', 'Ver Evaluaciones', 'Ver evaluaciones', 'assessments', 'read', 'unit'),
('assessments:update', 'Editar Evaluaciones', 'Modificar evaluaciones', 'assessments', 'update', 'unit'),
('assessments:delete', 'Eliminar Evaluaciones', 'Eliminar evaluaciones', 'assessments', 'delete', 'unit'),
('assessments:publish', 'Publicar Evaluaciones', 'Publicar evaluaciones para estudiantes', 'assessments', 'publish', 'unit'),
('assessments:grade', 'Calificar Evaluaciones', 'Calificar respuestas de estudiantes', 'assessments', 'grade', 'unit'),
('assessments:attempt', 'Rendir Evaluaciones', 'Intentar evaluaciones como estudiante', 'assessments', 'attempt', 'unit'),
('assessments:view_results', 'Ver Resultados', 'Ver resultados propios', 'assessments', 'view_results', 'unit'),

-- Permisos sobre PROGRESO
('progress:read', 'Ver Progreso', 'Ver progreso académico', 'progress', 'read', 'unit'),
('progress:update', 'Actualizar Progreso', 'Actualizar progreso de estudiantes', 'progress', 'update', 'unit'),
('progress:read:own', 'Ver Progreso Propio', 'Ver propio progreso', 'progress', 'read:own', 'unit'),

-- Permisos sobre ESTADÍSTICAS
('stats:global', 'Estadísticas Globales', 'Ver estadísticas de toda la plataforma', 'stats', 'global', 'system'),
('stats:school', 'Estadísticas de Escuela', 'Ver estadísticas de la institución', 'stats', 'school', 'school'),
('stats:unit', 'Estadísticas de Unidad', 'Ver estadísticas de la clase', 'stats', 'unit', 'unit');
