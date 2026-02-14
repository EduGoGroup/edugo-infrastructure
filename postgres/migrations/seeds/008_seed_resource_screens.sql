-- ====================================================================
-- SEEDS: Mapeos recurso-pantalla (vincula RBAC con screen instances)
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- Recursos existentes (de 001_seed_resources.sql):
--   materials:  20000000-0000-0000-0000-000000000030 (key: 'materials')
--   dashboard:  20000000-0000-0000-0000-000000000001 (key: 'dashboard')
--
-- Instancias (de 007_seed_screen_instances.sql):
--   materials-list:     b0000000-0000-0000-0000-000000000004
--   material-detail:    b0000000-0000-0000-0000-000000000005
--   dashboard-teacher:  b0000000-0000-0000-0000-000000000002
--   app-settings:       b0000000-0000-0000-0000-000000000006
--
-- NOTA: settings no tiene recurso RBAC propio; se accede directamente.
-- El recurso 'dashboard' mapea al dashboard del profesor como default.

-- Materials -> materials-list (list, default)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'materials-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Materials -> material-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'material-detail', 'detail', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-teacher (dashboard, default)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-teacher', 'dashboard', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
