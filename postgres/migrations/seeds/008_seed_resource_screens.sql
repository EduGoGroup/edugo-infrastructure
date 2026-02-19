-- ====================================================================
-- SEEDS: Mapeos recurso-pantalla (vincula RBAC con screen instances)
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- Recursos existentes (de 001_seed_resources.sql):
--   materials:    20000000-0000-0000-0000-000000000030 (key: 'materials')
--   dashboard:    20000000-0000-0000-0000-000000000001 (key: 'dashboard')
--   assessments:  20000000-0000-0000-0000-000000000031 (key: 'assessments')
--   progress:     20000000-0000-0000-0000-000000000040 (key: 'progress')
--   users:        20000000-0000-0000-0000-000000000010 (key: 'users')
--   schools:      20000000-0000-0000-0000-000000000011 (key: 'schools')
--   units:        20000000-0000-0000-0000-000000000020 (key: 'units')
--   memberships:  20000000-0000-0000-0000-000000000021 (key: 'memberships')
--   roles:        20000000-0000-0000-0000-000000000012 (key: 'roles')
--   permissions_mgmt: 20000000-0000-0000-0000-000000000013 (key: 'permissions_mgmt')
--
-- Instancias (de 007_seed_screen_instances.sql):
--   materials-list:     b0000000-0000-0000-0000-000000000004
--   material-detail:    b0000000-0000-0000-0000-000000000005
--   material-create:    b0000000-0000-0000-0000-000000000020
--   material-edit:      b0000000-0000-0000-0000-000000000021
--   dashboard-teacher:  b0000000-0000-0000-0000-000000000002
--   dashboard-superadmin: b0000000-0000-0000-0000-000000000010
--   dashboard-schooladmin: b0000000-0000-0000-0000-000000000011
--   dashboard-student:  b0000000-0000-0000-0000-000000000003
--   app-settings:       b0000000-0000-0000-0000-000000000006
--   assessments-list:   b0000000-0000-0000-0000-000000000030
--   assessment-take:    b0000000-0000-0000-0000-000000000031
--   assessment-result:  b0000000-0000-0000-0000-000000000032
--   attempts-history:   b0000000-0000-0000-0000-000000000033
--   progress-my:        b0000000-0000-0000-0000-000000000040
--   progress-unit-list: b0000000-0000-0000-0000-000000000041
--   progress-student-detail: b0000000-0000-0000-0000-000000000042
--   users-list:         b0000000-0000-0000-0000-000000000050
--   user-detail:        b0000000-0000-0000-0000-000000000051
--   user-create:        b0000000-0000-0000-0000-000000000052
--   user-edit:          b0000000-0000-0000-0000-000000000053
--   schools-list:       b0000000-0000-0000-0000-000000000060
--   school-detail:      b0000000-0000-0000-0000-000000000061
--   school-create:      b0000000-0000-0000-0000-000000000062
--   school-edit:        b0000000-0000-0000-0000-000000000063
--   units-list:         b0000000-0000-0000-0000-000000000070
--   unit-detail:        b0000000-0000-0000-0000-000000000071
--   unit-create:        b0000000-0000-0000-0000-000000000072
--   unit-edit:          b0000000-0000-0000-0000-000000000073
--   memberships-list:   b0000000-0000-0000-0000-000000000074
--   membership-add:     b0000000-0000-0000-0000-000000000075
--   dashboard-guardian:  b0000000-0000-0000-0000-000000000080
--   children-list:      b0000000-0000-0000-0000-000000000081
--   child-progress:     b0000000-0000-0000-0000-000000000082
--   roles-list:         b0000000-0000-0000-0000-000000000090
--   role-detail:        b0000000-0000-0000-0000-000000000091
--   resources-list:     b0000000-0000-0000-0000-000000000092
--   permissions-list:   b0000000-0000-0000-0000-000000000093
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

-- Materials -> material-create (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'material-create', 'form', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Materials -> material-edit (form-edit)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'material-edit', 'form-edit', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessments-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-take (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-take', 'form', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-result (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-result', 'detail', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> attempts-history (history)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000031', 'assessments', 'attempts-history', 'history', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Progress -> progress-my (dashboard)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000040', 'progress', 'progress-my', 'dashboard', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Progress -> progress-unit-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000040', 'progress', 'progress-unit-list', 'list', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Progress -> progress-student-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000040', 'progress', 'progress-student-detail', 'detail', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Users -> users-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'users-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Users -> user-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'user-detail', 'detail', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Users -> user-create (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'user-create', 'form', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Users -> user-edit (form-edit)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'user-edit', 'form-edit', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schools -> schools-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000011', 'schools', 'schools-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schools -> school-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000011', 'schools', 'school-detail', 'detail', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schools -> school-create (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000011', 'schools', 'school-create', 'form', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schools -> school-edit (form-edit)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000011', 'schools', 'school-edit', 'form-edit', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Units -> units-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'units-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Units -> unit-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'unit-detail', 'detail', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Units -> unit-create (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'unit-create', 'form', true, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Units -> unit-edit (form-edit)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'unit-edit', 'form-edit', false, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Memberships -> memberships-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000021', 'memberships', 'memberships-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Memberships -> membership-add (form)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000021', 'memberships', 'membership-add', 'form', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Dashboard -> dashboard-guardian
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'dashboard-guardian', 'dashboard-guardian', false, 5)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Roles -> roles-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000012', 'roles', 'roles-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Roles -> role-detail (detail)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000012', 'roles', 'role-detail', 'detail', true, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Permissions -> resources-list (list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'resources-list', 'list', true, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Permissions -> permissions-list (detail-list)
INSERT INTO ui_config.resource_screens (resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'permissions-list', 'detail-list', false, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
