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

COMMIT;
