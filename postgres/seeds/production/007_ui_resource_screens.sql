-- ====================================================================
-- SEEDS: Mapeos recurso-pantalla (vincula RBAC con screen instances)
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ====================================================================

BEGIN;

-- Materials -> materials-list (list, default)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'materials-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Materials -> material-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'material-detail', 'detail', false, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-teacher (dashboard, default)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-teacher', 'dashboard', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-superadmin
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-superadmin', 'dashboard-superadmin', false, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-schooladmin
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-schooladmin', 'dashboard-schooladmin', false, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-student
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-student', 'dashboard-student', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-guardian
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-guardian', 'dashboard-guardian', false, 5)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Recursos admin/academic/content/reports → pantallas
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default) VALUES
('c0000000-0000-0000-0000-000000000020', '20000000-0000-0000-0000-000000000010', 'users', 'users-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000021', '20000000-0000-0000-0000-000000000011', 'schools', 'schools-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000022', '20000000-0000-0000-0000-000000000012', 'roles', 'roles-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000023', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'permissions-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000024', '20000000-0000-0000-0000-000000000020', 'units', 'units-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000025', '20000000-0000-0000-0000-000000000021', 'memberships', 'memberships-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000026', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000027', '20000000-0000-0000-0000-000000000040', 'progress', 'progress-dashboard', 'dashboard', TRUE),
('c0000000-0000-0000-0000-000000000028', '20000000-0000-0000-0000-000000000041', 'stats', 'stats-dashboard', 'dashboard', TRUE),
('c0000000-0000-0000-0000-000000000029', '20000000-0000-0000-0000-000000000032', 'subjects', 'subjects-list', 'list', TRUE),
('c0000000-0000-0000-0000-000000000030', '20000000-0000-0000-0000-000000000032', 'subjects', 'subjects-form', 'form', FALSE)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Users -> users-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000031', '20000000-0000-0000-0000-000000000010', 'users', 'users-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Roles -> roles-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000032', '20000000-0000-0000-0000-000000000012', 'roles', 'roles-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Permissions -> permissions-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000033', '20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'permissions-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

COMMIT;
