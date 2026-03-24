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

-- Units -> units-form (form) [Fase 2.1]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000040', '20000000-0000-0000-0000-000000000020', 'units', 'units-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Memberships -> memberships-form (form) [Fase 2.2]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000041', '20000000-0000-0000-0000-000000000021', 'memberships', 'memberships-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Materials -> material-create (form) [Fase 2.3]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000042', '20000000-0000-0000-0000-000000000030', 'materials', 'material-create', 'form', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Screen Templates -> screen-templates-list (list) [Fase 5.5]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000050', '20000000-0000-0000-0000-000000000050', 'screen_templates', 'screen-templates-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Screen Instances -> screen-instances-list (list) [Fase 5.5]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000051', '20000000-0000-0000-0000-000000000051', 'screen_instances', 'screen-instances-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Screen Instances -> screen-instances-form (form) [Fase 5.5]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000052', '20000000-0000-0000-0000-000000000051', 'screen_instances', 'screen-instances-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Guardian Relations -> guardian-requests-list (list, default) [Fase 4.1]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000060', '20000000-0000-0000-0000-000000000060', 'guardian_relations', 'guardian-requests-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Guardian Relations -> children-list (detail-children) [Fase 4.1]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000061', '20000000-0000-0000-0000-000000000060', 'guardian_relations', 'children-list', 'detail-children', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Guardian Relations -> child-progress (detail) [Fase 4.1]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000062', '20000000-0000-0000-0000-000000000060', 'guardian_relations', 'child-progress', 'detail', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessments-form (form) [Fase 3]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000070', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-take (detail) [Fase 3]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000071', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-take', 'detail', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments Student View -> assessments-list (superadmin puede ver como estudiante)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000094', '20000000-0000-0000-0000-000000000033', 'assessments_student', 'assessments-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessments-management-list (superadmin view) [Fase 3 - Gestión]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000091', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-management-list', 'superadmin', FALSE, 7)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessments-management-list (teacher view) [Fase 3 - Gestión]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000092', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-management-list', 'teacher', FALSE, 8)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessments-management-list (school_coordinator view) [Fase 3 - Gestión]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000093', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessments-management-list', 'schoolcoordinator', FALSE, 9)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Audit -> audit-events-list (list, default) [Fase 6]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000072', '20000000-0000-0000-0000-000000000070', 'audit', 'audit-events-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-questions-list (questions) [Assessment CRUD]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000073', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-questions-list', 'questions', FALSE, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-question-form (question-form) [Assessment CRUD]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000074', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-question-form', 'question-form', FALSE, 5)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Assessments -> assessment-result (result) [Fase 3]
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000090', '20000000-0000-0000-0000-000000000031', 'assessments', 'assessment-result', 'result', FALSE, 6)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Concept Types -> concept-types-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000095', '20000000-0000-0000-0000-000000000080', 'concept_types', 'concept-types-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Concept Types -> concept-types-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000096', '20000000-0000-0000-0000-000000000080', 'concept_types', 'concept-types-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Audit -> audit-detail (detail)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000097', '20000000-0000-0000-0000-000000000070', 'audit', 'audit-detail', 'detail', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Materials -> materials-form (create-form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000098', '20000000-0000-0000-0000-000000000030', 'materials', 'materials-form', 'create-form', FALSE, 4)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Screen Instances -> screens-form (create-form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000099', '20000000-0000-0000-0000-000000000051', 'screen_instances', 'screens-form', 'create-form', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Progress -> progress-detail (detail)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000100', '20000000-0000-0000-0000-000000000040', 'progress', 'progress-detail', 'detail', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Stats -> stats-detail (detail)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000101', '20000000-0000-0000-0000-000000000041', 'stats', 'stats-detail', 'detail', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- System Settings -> system-settings (settings, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000102', '20000000-0000-0000-0000-000000000090', 'system_settings', 'system-settings', 'settings', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- ====================================================================
-- School Ecosystem — Fase 2 resources (periods, grades, attendance, schedules, announcements)
-- ====================================================================

-- Periods -> periods-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000110', '20000000-0000-0000-0000-000000000034', 'periods', 'periods-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Periods -> periods-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000111', '20000000-0000-0000-0000-000000000034', 'periods', 'periods-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Grades -> grades-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000112', '20000000-0000-0000-0000-000000000035', 'grades', 'grades-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Grades -> grades-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000113', '20000000-0000-0000-0000-000000000035', 'grades', 'grades-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Attendance -> attendance-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000114', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Attendance -> attendance-batch (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000115', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-batch', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Attendance -> attendance-summary (summary)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000116', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-summary', 'summary', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schedules -> schedules-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000117', '20000000-0000-0000-0000-000000000037', 'schedules', 'schedules-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schedules -> schedules-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000118', '20000000-0000-0000-0000-000000000037', 'schedules', 'schedules-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Announcements -> announcements-list (list, default)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000119', '20000000-0000-0000-0000-000000000038', 'announcements', 'announcements-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Announcements -> announcements-form (form)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000120', '20000000-0000-0000-0000-000000000038', 'announcements', 'announcements-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

COMMIT;
