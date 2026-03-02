-- ============================================================
-- SEED: permissions (60 registros)
-- Fecha: 2026-02-24
-- Fuente: Datos reales de producción (Neon)
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ============================================================

INSERT INTO iam.permissions (id, resource_id, name, display_name, description, action, scope, is_active)
VALUES
  -- assessments (resource: 20000000-0000-0000-0000-000000000031)
  ('8c8d7a5b-2688-4646-9888-bc53600dbbc0', '20000000-0000-0000-0000-000000000031', 'assessments:attempt',      'Rendir Evaluaciones',               'Intentar evaluaciones como estudiante',       'attempt',      'unit',   TRUE),
  ('0e72f5de-e2df-4bc6-b6ec-68266855d1e8', '20000000-0000-0000-0000-000000000031', 'assessments:create',       'Crear Evaluaciones',                 'Crear evaluaciones y exámenes',               'create',       'unit',   TRUE),
  ('8b39f083-c630-4d1c-92ec-58afabb3376b', '20000000-0000-0000-0000-000000000031', 'assessments:delete',       'Eliminar Evaluaciones',              'Eliminar evaluaciones',                       'delete',       'unit',   TRUE),
  ('d477a2fd-f996-41b5-9dd2-d66cda3460d6', '20000000-0000-0000-0000-000000000031', 'assessments:grade',        'Calificar Evaluaciones',             'Calificar respuestas de estudiantes',          'grade',        'unit',   TRUE),
  ('1b69ea9a-1e15-4c38-a5e1-68c408bb1b97', '20000000-0000-0000-0000-000000000031', 'assessments:publish',      'Publicar Evaluaciones',              'Publicar evaluaciones para estudiantes',       'publish',      'unit',   TRUE),
  ('2606efba-8615-48ef-bc2a-6bf576a9158c', '20000000-0000-0000-0000-000000000031', 'assessments:read',         'Ver Evaluaciones',                   'Ver evaluaciones',                            'read',         'unit',   TRUE),
  ('91c8ae21-a955-4d59-b52e-6cf74bd6532b', '20000000-0000-0000-0000-000000000031', 'assessments:update',       'Editar Evaluaciones',                'Modificar evaluaciones',                      'update',       'unit',   TRUE),
  ('b457a385-29bc-4e5e-a79d-06c74f81d23e', '20000000-0000-0000-0000-000000000031', 'assessments:view_results', 'Ver Resultados',                     'Ver resultados propios',                      'view_results', 'unit',   TRUE),

  -- materials (resource: 20000000-0000-0000-0000-000000000030)
  ('a0ffd5c0-8781-48db-bc14-29933923c1b1', '20000000-0000-0000-0000-000000000030', 'materials:create',         'Crear Materiales',                   'Crear materiales educativos',                  'create',       'unit',   TRUE),
  ('6358de3d-11ef-49c0-be42-da51ffcdfbc1', '20000000-0000-0000-0000-000000000030', 'materials:delete',         'Eliminar Materiales',                'Eliminar materiales',                         'delete',       'unit',   TRUE),
  ('9b0c10e0-0a7b-4e73-af9a-2d5adf99790f', '20000000-0000-0000-0000-000000000030', 'materials:download',       'Descargar Materiales',               'Descargar materiales educativos',              'download',     'unit',   TRUE),
  ('9aba2d20-ca23-403e-b127-6bc967eec751', '20000000-0000-0000-0000-000000000030', 'materials:publish',        'Publicar Materiales',                'Publicar materiales para estudiantes',         'publish',      'unit',   TRUE),
  ('bd681ea0-a974-4d16-86dd-30e9e07e9970', '20000000-0000-0000-0000-000000000030', 'materials:read',           'Ver Materiales',                     'Ver materiales educativos',                    'read',         'unit',   TRUE),
  ('3436645f-af5c-4a3e-b09d-ff580c119427', '20000000-0000-0000-0000-000000000030', 'materials:update',         'Editar Materiales',                  'Modificar materiales',                        'update',       'unit',   TRUE),

  -- memberships (resource: 20000000-0000-0000-0000-000000000021)
  ('989bbeb5-9884-4728-8c98-d87d9d27f088', '20000000-0000-0000-0000-000000000021', 'memberships:create',       'Crear Membresías',                   'Asignar usuarios a unidades académicas',       'create',       'school', TRUE),
  ('c78ca0cf-0b6d-4e70-94d8-5ee0d021233f', '20000000-0000-0000-0000-000000000021', 'memberships:delete',       'Eliminar Membresías',                'Eliminar membresías de unidades',              'delete',       'school', TRUE),
  ('28dfb6b5-c680-4442-8530-73b67199fbcb', '20000000-0000-0000-0000-000000000021', 'memberships:read',         'Ver Membresías',                     'Ver membresías de unidades académicas',        'read',         'school', TRUE),
  ('0f53cce3-0133-4f93-9b8c-7c62b3a8eb3c', '20000000-0000-0000-0000-000000000021', 'memberships:update',       'Editar Membresías',                  'Modificar membresías',                        'update',       'school', TRUE),

  -- permissions_mgmt (resource: 20000000-0000-0000-0000-000000000013)
  ('31000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt:create',  'Crear Permisos',                      'Crear nuevos permisos en el sistema',           'create',       'system', TRUE),
  ('31000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt:delete',  'Eliminar Permisos',                   'Eliminar permisos del sistema',                 'delete',       'system', TRUE),
  ('6cfd0e01-6834-4fde-94e9-536b663b5be4', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt:read',    'Ver Configuración de Permisos',       'Ver recursos y permisos del sistema',          'read',         'system', TRUE),
  ('59da389c-a246-4d43-a6ce-f316584b2be7', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt:update',  'Editar Configuración de Permisos',    'Modificar recursos y permisos del sistema',    'update',       'system', TRUE),

  -- progress (resource: 20000000-0000-0000-0000-000000000040)
  ('89336e44-3636-4744-a056-aea878f57b18', '20000000-0000-0000-0000-000000000040', 'progress:read',            'Ver Progreso',                       'Ver progreso académico',                       'read',         'unit',   TRUE),
  ('d033dfd0-47a5-476b-b51a-f52f5fc66d7a', '20000000-0000-0000-0000-000000000040', 'progress:read:own',        'Ver Progreso Propio',                'Ver propio progreso',                         'read:own',     'unit',   TRUE),
  ('19d017d1-ee5b-4fc3-828a-2d12056631b4', '20000000-0000-0000-0000-000000000040', 'progress:update',          'Actualizar Progreso',                'Actualizar progreso de estudiantes',           'update',       'unit',   TRUE),

  -- roles (resource: 20000000-0000-0000-0000-000000000012)
  ('a28c2133-9d49-46ae-8cb7-fe59ff0246df', '20000000-0000-0000-0000-000000000012', 'roles:read',               'Ver Roles',                          'Ver roles del sistema',                        'read',         'system', TRUE),
  ('31000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000012', 'roles:create',             'Crear Roles',                        'Crear nuevos roles en el sistema',              'create',       'system', TRUE),
  ('31000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000012', 'roles:delete',             'Eliminar Roles',                     'Eliminar roles del sistema',                   'delete',       'system', TRUE),
  ('dfcc999e-fd49-4611-9d52-a2baf8b95851', '20000000-0000-0000-0000-000000000012', 'roles:update',             'Editar Roles',                       'Modificar configuración de roles',             'update',       'system', TRUE),

  -- schools (resource: 20000000-0000-0000-0000-000000000011)
  ('611df7ce-b4cd-474f-901d-9bfd8873a9c1', '20000000-0000-0000-0000-000000000011', 'schools:create',           'Crear Escuelas',                     'Crear nuevas instituciones educativas',        'create',       'system', TRUE),
  ('5bd8088b-1506-4b22-aa7e-9e4eb50de24e', '20000000-0000-0000-0000-000000000011', 'schools:delete',           'Eliminar Escuelas',                  'Eliminar escuelas del sistema',               'delete',       'system', TRUE),
  ('8545c3be-3117-40a1-b1fb-da78d6233ae1', '20000000-0000-0000-0000-000000000011', 'schools:manage',           'Gestionar Escuela',                  'Control total de la escuela',                  'manage',       'school', TRUE),
  ('bc15c7a1-f203-46e0-80be-2850fad94b0e', '20000000-0000-0000-0000-000000000011', 'schools:read',             'Ver Escuelas',                       'Ver información de escuelas',                  'read',         'system', TRUE),
  ('2b823ad1-d875-4951-9c85-3baafa3f1f65', '20000000-0000-0000-0000-000000000011', 'schools:update',           'Editar Escuelas',                    'Modificar datos de escuelas',                  'update',       'school', TRUE),

  -- screen_instances (resource: 20000000-0000-0000-0000-000000000051)
  ('fa35b956-665f-48f4-a51e-ad1393e72652', '20000000-0000-0000-0000-000000000051', 'screen_instances:create',  'Crear Instancias de Pantalla',        'Crear nuevas instancias de pantalla',          'create',       'system', TRUE),
  ('4096f489-b3f8-49bd-8ecb-6e3588a85f84', '20000000-0000-0000-0000-000000000051', 'screen_instances:delete',  'Eliminar Instancias de Pantalla',     'Eliminar instancias de pantalla configuradas', 'delete',       'system', TRUE),
  ('ebfd0911-43bf-42ef-9523-8dd93079db47', '20000000-0000-0000-0000-000000000051', 'screen_instances:read',    'Ver Instancias de Pantalla',          'Ver instancias de pantalla configuradas',      'read',         'system', TRUE),
  ('1ad07392-4c86-4ec9-b249-b66be3f97ce8', '20000000-0000-0000-0000-000000000051', 'screen_instances:update',  'Actualizar Instancias de Pantalla',   'Modificar instancias de pantalla existentes',  'update',       'system', TRUE),

  -- screen_templates (resource: 20000000-0000-0000-0000-000000000050)
  ('52011396-5981-4c59-a772-1f353d10a3e9', '20000000-0000-0000-0000-000000000050', 'screen_templates:create',  'Crear Templates de Pantalla',         'Crear nuevos templates de pantalla',           'create',       'system', TRUE),
  ('b6db1991-4a2c-429a-9c45-0ed177b6e3ed', '20000000-0000-0000-0000-000000000050', 'screen_templates:delete',  'Eliminar Templates de Pantalla',      'Eliminar templates de pantalla del sistema',   'delete',       'system', TRUE),
  ('3d89c941-cbe5-4c1b-8cf0-0b55b4aaa313', '20000000-0000-0000-0000-000000000050', 'screen_templates:read',    'Ver Templates de Pantalla',           'Ver templates de pantalla del sistema',        'read',         'system', TRUE),
  ('e5bf88e6-73ff-40d4-93a4-8c787d3930af', '20000000-0000-0000-0000-000000000050', 'screen_templates:update',  'Actualizar Templates de Pantalla',    'Modificar templates de pantalla existentes',   'update',       'system', TRUE),

  -- screens / mobile (resource: 20000000-0000-0000-0000-000000000052)
  ('2b31df13-4c54-43fc-8bcd-8a9265fba1a0', '20000000-0000-0000-0000-000000000052', 'screens:read',             'Leer Pantallas (Mobile)',             'Leer configuracion de pantallas desde mobile', 'read',         'system', TRUE),

  -- subjects (resource: 20000000-0000-0000-0000-000000000032)
  ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000032', 'subjects:create',          'Crear Materia',                      'Crear materias en el plan de estudios',         'create',       'school', TRUE),
  ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000032', 'subjects:read',            'Ver Materias',                       'Ver materias del plan de estudios',             'read',         'school', TRUE),
  ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000032', 'subjects:update',          'Editar Materia',                     'Modificar datos de materias',                  'update',       'school', TRUE),
  ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000032', 'subjects:delete',          'Eliminar Materia',                   'Eliminar materias del plan de estudios',        'delete',       'school', TRUE),

  -- stats (resource: 20000000-0000-0000-0000-000000000041)
  ('8a9fbae4-1b64-4870-ad14-41c436348bcc', '20000000-0000-0000-0000-000000000041', 'stats:global',             'Estadísticas Globales',               'Ver estadísticas de toda la plataforma',       'global',       'system', TRUE),
  ('f35d45c1-9539-422d-974f-5075d8f9b296', '20000000-0000-0000-0000-000000000041', 'stats:school',             'Estadísticas de Escuela',             'Ver estadísticas de la institución',           'school',       'school', TRUE),
  ('f47983b3-a721-461e-8de5-05fea4eda3fe', '20000000-0000-0000-0000-000000000041', 'stats:unit',               'Estadísticas de Unidad',              'Ver estadísticas de la clase',                 'unit',         'unit',   TRUE),

  -- units (resource: 20000000-0000-0000-0000-000000000020)
  ('619f6f66-6806-4894-a965-3c266a483be3', '20000000-0000-0000-0000-000000000020', 'units:create',             'Crear Unidades',                     'Crear unidades académicas (clases, grados)',    'create',       'school', TRUE),
  ('8d3a079d-5b6c-452f-ab49-b725547a052c', '20000000-0000-0000-0000-000000000020', 'units:delete',             'Eliminar Unidades',                  'Eliminar unidades académicas',                 'delete',       'school', TRUE),
  ('61633d6c-aa56-40c1-a048-8e21f2893058', '20000000-0000-0000-0000-000000000020', 'units:read',               'Ver Unidades',                       'Ver unidades académicas',                      'read',         'school', TRUE),
  ('4809f4d8-16dc-4222-9e12-1fca5f3c7ab7', '20000000-0000-0000-0000-000000000020', 'units:update',             'Editar Unidades',                    'Modificar unidades académicas',               'update',       'school', TRUE),

  -- users (resource: 20000000-0000-0000-0000-000000000010)
  ('eff25f87-711d-43a5-b8d3-1e3fb6be6a19', '20000000-0000-0000-0000-000000000010', 'users:create',             'Crear Usuarios',                     'Crear nuevos usuarios en el sistema',          'create',       'system', TRUE),
  ('4129c4b5-89a1-4908-8b21-b6289e1ad095', '20000000-0000-0000-0000-000000000010', 'users:delete',             'Eliminar Usuarios',                  'Eliminar usuarios del sistema',               'delete',       'system', TRUE),
  ('1ae1ad50-857c-4601-8378-5fd25128f11d', '20000000-0000-0000-0000-000000000010', 'users:read',               'Ver Usuarios',                       'Ver información de usuarios',                  'read',         'school', TRUE),
  ('8098577f-e5d8-4e07-aee7-4c1521cbe88b', '20000000-0000-0000-0000-000000000010', 'users:read:own',           'Ver Perfil Propio',                  'Ver propio perfil de usuario',                 'read:own',     'system', TRUE),
  ('813077d4-4624-4817-b3f2-69d60f6cb7a9', '20000000-0000-0000-0000-000000000010', 'users:update',             'Editar Usuarios',                    'Modificar datos de usuarios',                  'update',       'school', TRUE),
  ('668b0d86-f4cb-45cd-a8c4-afa8f6cfe9b6', '20000000-0000-0000-0000-000000000010', 'users:update:own',         'Editar Perfil Propio',               'Modificar propio perfil',                      'update:own',   'system', TRUE)
ON CONFLICT (id) DO NOTHING;
