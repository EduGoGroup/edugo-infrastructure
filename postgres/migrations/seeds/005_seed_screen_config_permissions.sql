-- ====================================================================
-- SEEDS: Recursos y permisos para configuracion de pantallas (Dynamic UI)
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- NOTA: Se crean 3 recursos separados (screen_templates, screen_instances, screens)
-- porque el constraint uq_permissions_resource_action requiere UNIQUE(resource_id, action).
-- Cada sub-entidad necesita su propio resource_id para tener permisos CRUD independientes.

-- Recursos para Dynamic UI (no visibles en menu, solo para RBAC)
INSERT INTO public.resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000050', 'screen_templates', 'Templates de Pantalla', 'Templates base para configuracion de pantallas', 'settings_applications', '20000000-0000-0000-0000-000000000002', 5, false, 'system'),
('20000000-0000-0000-0000-000000000051', 'screen_instances', 'Instancias de Pantalla', 'Instancias configuradas de pantalla por escuela', 'devices', '20000000-0000-0000-0000-000000000002', 6, false, 'system'),
('20000000-0000-0000-0000-000000000052', 'screens', 'Pantallas (Mobile)', 'Lectura de pantallas desde aplicacion mobile', 'smartphone', NULL, 0, false, 'system')
ON CONFLICT (key) DO NOTHING;

-- Permisos de screen_templates (CRUD admin)
INSERT INTO public.permissions (name, display_name, description, resource_id, action, scope) VALUES
('screen_templates:read',   'Ver Templates de Pantalla',        'Ver templates de pantalla del sistema',      '20000000-0000-0000-0000-000000000050', 'read',   'system'),
('screen_templates:create', 'Crear Templates de Pantalla',      'Crear nuevos templates de pantalla',         '20000000-0000-0000-0000-000000000050', 'create', 'system'),
('screen_templates:update', 'Actualizar Templates de Pantalla', 'Modificar templates de pantalla existentes', '20000000-0000-0000-0000-000000000050', 'update', 'system'),
('screen_templates:delete', 'Eliminar Templates de Pantalla',   'Eliminar templates de pantalla del sistema', '20000000-0000-0000-0000-000000000050', 'delete', 'system')
ON CONFLICT (name) DO NOTHING;

-- Permisos de screen_instances (CRUD admin)
INSERT INTO public.permissions (name, display_name, description, resource_id, action, scope) VALUES
('screen_instances:read',   'Ver Instancias de Pantalla',       'Ver instancias de pantalla configuradas',      '20000000-0000-0000-0000-000000000051', 'read',   'system'),
('screen_instances:create', 'Crear Instancias de Pantalla',     'Crear nuevas instancias de pantalla',          '20000000-0000-0000-0000-000000000051', 'create', 'system'),
('screen_instances:update', 'Actualizar Instancias de Pantalla','Modificar instancias de pantalla existentes',  '20000000-0000-0000-0000-000000000051', 'update', 'system'),
('screen_instances:delete', 'Eliminar Instancias de Pantalla',  'Eliminar instancias de pantalla configuradas', '20000000-0000-0000-0000-000000000051', 'delete', 'system')
ON CONFLICT (name) DO NOTHING;

-- Permisos de screens (lectura mobile)
INSERT INTO public.permissions (name, display_name, description, resource_id, action, scope) VALUES
('screens:read', 'Leer Pantallas (Mobile)', 'Leer configuracion de pantallas desde mobile', '20000000-0000-0000-0000-000000000052', 'read', 'system')
ON CONFLICT (name) DO NOTHING;

-- Asignar todos los permisos de screen config al rol super_admin
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.roles r, public.permissions p
WHERE r.name = 'super_admin'
  AND p.resource_id IN (
    '20000000-0000-0000-0000-000000000050',
    '20000000-0000-0000-0000-000000000051',
    '20000000-0000-0000-0000-000000000052'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;
