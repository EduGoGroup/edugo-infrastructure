-- ====================================================================
-- SEEDS: Instancias de pantalla para las 6 pantallas de la Fase 1
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- UUIDs fijos para referenciar desde resource_screens (008):
--   app-login:          b0000000-0000-0000-0000-000000000001
--   dashboard-teacher:  b0000000-0000-0000-0000-000000000002
--   dashboard-student:  b0000000-0000-0000-0000-000000000003
--   materials-list:     b0000000-0000-0000-0000-000000000004
--   material-detail:    b0000000-0000-0000-0000-000000000005
--   app-settings:       b0000000-0000-0000-0000-000000000006
--
-- Template IDs (de 006_seed_screen_templates.sql):
--   login-basic-v1:     a0000000-0000-0000-0000-000000000001
--   dashboard-basic-v1: a0000000-0000-0000-0000-000000000002
--   list-basic-v1:      a0000000-0000-0000-0000-000000000003
--   detail-basic-v1:    a0000000-0000-0000-0000-000000000004
--   settings-basic-v1:  a0000000-0000-0000-0000-000000000005

-- Instancia 1: Login
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000001', 'app-login',
 'a0000000-0000-0000-0000-000000000001', 'Login', 'Pantalla de inicio de sesion',
 '{
   "app_logo": "edugo_logo",
   "app_name": "EduGo",
   "app_tagline": "Learning made easy",
   "email_label": "Email",
   "password_label": "Password",
   "remember_label": "Remember me",
   "login_btn_label": "Sign In",
   "forgot_password_label": "Forgot password?",
   "divider_text": "or continue with",
   "google_btn_label": "Google"
 }'::jsonb,
 '[
   {
     "id": "submit-login",
     "trigger": "button_click",
     "triggerSlotId": "login_btn",
     "type": "SUBMIT_FORM",
     "config": {
       "endpoint": "/v1/auth/login",
       "method": "POST",
       "fieldMapping": {"email": "email", "password": "password"},
       "onSuccess": {"type": "NAVIGATE", "config": {"target": "dashboard-home"}}
     }
   }
 ]'::jsonb,
 NULL, '{}'::jsonb, 'system', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 2: Dashboard Profesor
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000002', 'dashboard-teacher',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Profesor', 'Panel principal del profesor',
 '{
   "page_title": "Dashboard",
   "greeting_text": "Good morning, {user.firstName}",
   "date_text": "{today_date}",
   "kpi_students_label": "Students",
   "kpi_materials_label": "Materials",
   "kpi_avg_score_label": "Avg Score",
   "kpi_completion_label": "Completion",
   "activity_title": "Recent Activity",
   "upload_label": "Upload Material",
   "progress_label": "View Progress"
 }'::jsonb,
 '[
   {"id": "navigate-materials", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/stats/global',
 '{"method": "GET", "fieldMapping": {"total_students": "total_students", "total_materials": "total_materials", "avg_score": "avg_score", "completion_rate": "completion_rate"}}'::jsonb,
 'school', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 3: Dashboard Estudiante
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000003', 'dashboard-student',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Estudiante', 'Panel principal del estudiante',
 '{
   "page_title": "Home",
   "greeting_text": "Hello, {user.firstName}!",
   "date_text": "{today_date}",
   "kpi_students_label": "Courses",
   "kpi_materials_label": "Materials",
   "kpi_avg_score_label": "My Score",
   "kpi_completion_label": "Progress",
   "activity_title": "Recent Activity",
   "upload_label": "My Materials",
   "progress_label": "My Progress"
 }'::jsonb,
 '[
   {"id": "navigate-materials", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/stats/student',
 '{"method": "GET", "fieldMapping": {"total_students": "enrolled_courses", "total_materials": "available_materials", "avg_score": "my_avg_score", "completion_rate": "my_completion_rate"}}'::jsonb,
 'unit', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 4: Lista de Materiales
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000004', 'materials-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Materiales', 'Lista de materiales educativos',
 '{
   "page_title": "Materials",
   "search_placeholder": "Search materials...",
   "filter_all_label": "All",
   "filter_ready_label": "Ready",
   "filter_processing_label": "Processing",
   "empty_icon": "folder_open",
   "empty_state_title": "No materials yet",
   "empty_state_description": "Upload your first educational material",
   "empty_action_label": "Upload Material"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "material-detail", "params": {"id": "{item.id}"}}},
   {"id": "search", "trigger": "button_click", "triggerSlotId": "search_bar", "type": "REFRESH", "config": {"addParams": {"search": "{search_bar.value}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/materials',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "defaultParams": {"sort": "created_at", "order": "desc"},
   "fieldMapping": {"title": "title", "subtitle": "subject", "status": "status", "file_type_icon": "file_type", "created_at": "created_at", "id": "id"}
 }'::jsonb,
 'unit', 'materials:read')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 5: Detalle de Material
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000005', 'material-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Material', 'Detalle de un material educativo',
 '{
   "page_title": "Material Detail",
   "file_size_label": "File Size",
   "uploaded_label": "Uploaded",
   "status_label": "Status",
   "description_title": "Description",
   "summary_title": "AI Summary",
   "download_label": "Download",
   "quiz_label": "Take Quiz"
 }'::jsonb,
 '[
   {"id": "download-material", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "API_CALL", "config": {"endpoint": "/v1/materials/{id}/download-url", "method": "GET", "onSuccess": {"type": "OPEN_URL", "config": {"url": "{response.url}"}}}},
   {"id": "take-quiz", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "assessment-view", "params": {"materialId": "{item.id}"}}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/materials/{id}',
 '{"method": "GET", "fieldMapping": {"title": "title", "subject": "subject", "grade": "grade", "status": "status", "file_type": "file_type", "file_size_display": "file_size_display", "description": "description", "created_at": "created_at", "summary": "summary"}}'::jsonb,
 'unit', 'materials:read')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 6: Configuracion
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-000000000006', 'app-settings',
 'a0000000-0000-0000-0000-000000000005', 'Configuracion', 'Pantalla de configuracion de la aplicacion',
 '{
   "page_title": "Settings",
   "appearance_title": "Appearance",
   "dark_mode_label": "Dark Mode",
   "theme_label": "Theme Color",
   "notifications_title": "Notifications",
   "push_label": "Push Notifications",
   "email_label": "Email Notifications",
   "account_title": "Account",
   "change_password_label": "Change Password",
   "language_label": "Language",
   "about_title": "About",
   "version_label": "App Version 1.0.0",
   "privacy_label": "Privacy Policy",
   "terms_label": "Terms of Service",
   "logout_label": "Sign Out"
 }'::jsonb,
 '[
   {"id": "toggle-dark-mode", "trigger": "button_click", "triggerSlotId": "dark_mode", "type": "API_CALL", "config": {"handler": "theme_toggle", "local": true}},
   {"id": "logout", "trigger": "button_click", "triggerSlotId": "logout_btn", "type": "CONFIRM", "config": {"title": "Sign Out", "message": "Are you sure you want to sign out?", "confirmLabel": "Sign Out", "onConfirm": {"type": "LOGOUT"}}}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'system', NULL)
ON CONFLICT (screen_key) DO NOTHING;
