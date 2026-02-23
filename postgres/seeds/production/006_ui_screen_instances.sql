-- ====================================================================
-- SEEDS: Instancias de pantalla para las pantallas configuradas
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ====================================================================

BEGIN;

-- Instancia 1: Login
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000001', 'app-login',
 'a0000000-0000-0000-0000-000000000001', 'Login', 'Pantalla de inicio de sesion',
 '{"app_logo": "edugo_logo", "app_name": "EduGo", "app_tagline": "Learning made easy", "email_label": "Email", "password_label": "Password", "remember_label": "Remember me", "login_btn_label": "Sign In", "forgot_password_label": "Forgot password?", "divider_text": "or continue with", "google_btn_label": "Google"}'::jsonb,
 '[{"id": "submit-login", "trigger": "button_click", "triggerSlotId": "login_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "/v1/auth/login", "method": "POST", "fieldMapping": {"email": "email", "password": "password"}, "onSuccess": {"type": "NAVIGATE", "config": {"target": "dashboard-home"}}}}]'::jsonb,
 NULL, '{}'::jsonb, 'system', NULL, 'login')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 2: Dashboard Profesor
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000002', 'dashboard-teacher',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Profesor', 'Panel principal del profesor',
 '{"page_title": "Dashboard", "greeting_text": "Good morning, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Students", "kpi_materials_label": "Materials", "kpi_avg_score_label": "Avg Score", "kpi_completion_label": "Completion", "activity_title": "Recent Activity", "upload_label": "Upload Material", "progress_label": "View Progress"}'::jsonb,
 '[{"id": "navigate-materials", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "materials-list"}}, {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 '/v1/stats/global', '{"method": "GET"}'::jsonb, 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 3: Dashboard Estudiante
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000003', 'dashboard-student',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Estudiante', 'Panel principal del estudiante',
 '{"page_title": "Home", "greeting_text": "Hello, {user.firstName}!", "date_text": "{today_date}", "kpi_students_label": "Courses", "kpi_materials_label": "Materials", "kpi_avg_score_label": "My Score", "kpi_completion_label": "Progress", "activity_title": "Recent Activity", "upload_label": "My Materials", "progress_label": "My Progress"}'::jsonb,
 '[{"id": "navigate-materials", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "materials-list"}}, {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 '/v1/stats/student', '{"method": "GET"}'::jsonb, 'unit', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 4: Lista de Materiales
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000004', 'materials-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Materiales', 'Lista de materiales educativos',
 '{"page_title": "Materials", "search_placeholder": "Search materials...", "filter_all_label": "All", "filter_ready_label": "Ready", "filter_processing_label": "Processing", "empty_icon": "folder_open", "empty_state_title": "No materials yet", "empty_state_description": "Upload your first educational material", "empty_action_label": "Upload Material"}'::jsonb,
 '[{"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "material-detail", "params": {"id": "{item.id}"}}}, {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 '/v1/materials', '{"method": "GET", "pagination": {"type": "offset", "pageSize": 20}}'::jsonb, 'unit', 'materials:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 5: Detalle de Material
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000005', 'material-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Material', 'Detalle de un material educativo',
 '{"page_title": "Material Detail", "file_size_label": "File Size", "uploaded_label": "Uploaded", "status_label": "Status", "description_title": "Description", "summary_title": "AI Summary", "download_label": "Download", "quiz_label": "Take Quiz"}'::jsonb,
 '[{"id": "download-material", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "API_CALL", "config": {"endpoint": "/v1/materials/{id}/download-url", "method": "GET"}}, {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}]'::jsonb,
 '/v1/materials/{id}', '{"method": "GET"}'::jsonb, 'unit', 'materials:read', 'material-detail')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 6: Configuracion
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000006', 'app-settings',
 'a0000000-0000-0000-0000-000000000005', 'Configuracion', 'Pantalla de configuracion',
 '{"page_title": "Settings", "appearance_title": "Appearance", "dark_mode_label": "Dark Mode", "theme_label": "Theme Color", "notifications_title": "Notifications", "push_label": "Push Notifications", "email_label": "Email Notifications", "logout_label": "Sign Out"}'::jsonb,
 '[{"id": "logout", "trigger": "button_click", "triggerSlotId": "logout_btn", "type": "CONFIRM", "config": {"title": "Sign Out", "message": "Are you sure?", "onConfirm": {"type": "LOGOUT"}}}]'::jsonb,
 NULL, '{}'::jsonb, 'system', NULL, 'settings')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 7: Dashboard Superadmin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000010', 'dashboard-superadmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Superadmin', 'Panel principal del superadmin',
 '{"page_title": "System Dashboard", "greeting_text": "Welcome, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Total Schools", "kpi_materials_label": "Total Users", "kpi_avg_score_label": "Total Materials", "kpi_completion_label": "System Health", "activity_title": "System Activity", "upload_label": "Manage Schools", "progress_label": "View Stats"}'::jsonb,
 '[{"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 'admin:/v1/stats/global', '{"method": "GET"}'::jsonb, 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 8: Dashboard School Admin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000011', 'dashboard-schooladmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard School Admin', 'Panel principal del admin de escuela',
 '{"page_title": "School Dashboard", "greeting_text": "Welcome, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Teachers", "kpi_materials_label": "Students", "kpi_avg_score_label": "Units", "kpi_completion_label": "School Score", "activity_title": "School Activity", "upload_label": "Manage Users", "progress_label": "View Reports"}'::jsonb,
 '[{"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 '/v1/stats/global', '{"method": "GET"}'::jsonb, 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 9: Dashboard Guardian
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000080', 'dashboard-guardian',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Guardian', 'Panel principal del apoderado',
 '{"page_title": "Family Dashboard", "greeting_text": "Hello, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Children", "kpi_materials_label": "Activities", "kpi_avg_score_label": "Avg Score", "kpi_completion_label": "Overall Progress", "activity_title": "Recent Activity", "upload_label": "My Children", "progress_label": "View Progress"}'::jsonb,
 '[{"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}]'::jsonb,
 '/v1/guardians/me/stats', '{"method": "GET"}'::jsonb, 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

COMMIT;
