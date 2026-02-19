-- ====================================================================
-- SEEDS: Instancias de pantalla para las pantallas de la Fase 1-3
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- UUIDs fijos para referenciar desde resource_screens (008):
--   app-login:          b0000000-0000-0000-0000-000000000001
--   dashboard-teacher:  b0000000-0000-0000-0000-000000000002
--   dashboard-student:  b0000000-0000-0000-0000-000000000003
--   materials-list:     b0000000-0000-0000-0000-000000000004
--   material-detail:    b0000000-0000-0000-0000-000000000005
--   app-settings:       b0000000-0000-0000-0000-000000000006
--   dashboard-superadmin: b0000000-0000-0000-0000-000000000010
--   dashboard-schooladmin: b0000000-0000-0000-0000-000000000011
--   material-create:    b0000000-0000-0000-0000-000000000020
--   material-edit:      b0000000-0000-0000-0000-000000000021
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
-- Template IDs (de 006_seed_screen_templates.sql):
--   login-basic-v1:     a0000000-0000-0000-0000-000000000001
--   dashboard-basic-v1: a0000000-0000-0000-0000-000000000002
--   list-basic-v1:      a0000000-0000-0000-0000-000000000003
--   detail-basic-v1:    a0000000-0000-0000-0000-000000000004
--   settings-basic-v1:  a0000000-0000-0000-0000-000000000005
--   form-basic-v1:      a0000000-0000-0000-0000-000000000006

-- Instancia 1: Login
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 NULL, '{}'::jsonb, 'system', NULL, 'login')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 2: Dashboard Profesor
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 3: Dashboard Estudiante
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 'unit', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 4: Lista de Materiales
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 'unit', 'materials:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 5: Detalle de Material
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 'unit', 'materials:read', 'material-detail')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 6: Configuracion
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
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
 NULL, '{}'::jsonb, 'system', NULL, 'settings')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 7: Dashboard Superadmin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000010', 'dashboard-superadmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Superadmin', 'Panel principal del superadmin',
 '{
   "page_title": "System Dashboard",
   "greeting_text": "Welcome, {user.firstName}",
   "date_text": "{today_date}",
   "kpi_students_label": "Total Schools",
   "kpi_materials_label": "Total Users",
   "kpi_avg_score_label": "Total Materials",
   "kpi_completion_label": "System Health",
   "activity_title": "System Activity",
   "upload_label": "Manage Schools",
   "progress_label": "View Stats"
 }'::jsonb,
 '[
   {"id": "navigate-schools", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "schools-list"}},
   {"id": "navigate-stats", "trigger": "button_click", "triggerSlotId": "view_progress", "type": "NAVIGATE", "config": {"target": "stats-global"}},
   {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/stats/global',
 '{"method": "GET", "fieldMapping": {"total_students": "total_schools", "total_materials": "total_users", "avg_score": "total_materials", "completion_rate": "system_health"}}'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 8: Dashboard School Admin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000011', 'dashboard-schooladmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard School Admin', 'Panel principal del admin de escuela',
 '{
   "page_title": "School Dashboard",
   "greeting_text": "Welcome, {user.firstName}",
   "date_text": "{today_date}",
   "kpi_students_label": "Teachers",
   "kpi_materials_label": "Students",
   "kpi_avg_score_label": "Units",
   "kpi_completion_label": "School Score",
   "activity_title": "School Activity",
   "upload_label": "Manage Users",
   "progress_label": "View Reports"
 }'::jsonb,
 '[
   {"id": "navigate-users", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "users-list"}},
   {"id": "navigate-reports", "trigger": "button_click", "triggerSlotId": "view_progress", "type": "NAVIGATE", "config": {"target": "progress-school"}},
   {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/stats/global',
 '{"method": "GET", "fieldMapping": {"total_students": "total_teachers", "total_materials": "total_students", "avg_score": "total_units", "completion_rate": "school_score"}}'::jsonb,
 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 9: Material Create
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000020', 'material-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Material', 'Formulario para crear un nuevo material',
 '{
   "page_title": "New Material",
   "form_title": "Create Material",
   "form_description": "Fill in the details to create a new educational material",
   "cancel_label": "Cancel",
   "submit_label": "Create Material",
   "form_fields": [
     {"id": "title", "controlType": "text-input", "label": "Title", "placeholder": "Enter material title", "required": true},
     {"id": "subject", "controlType": "text-input", "label": "Subject", "placeholder": "e.g. Mathematics, Science"},
     {"id": "grade", "controlType": "text-input", "label": "Grade", "placeholder": "e.g. 5th Grade"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the material content", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-material", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "/v1/materials", "method": "POST", "fieldMapping": {"title": "title", "subject": "subject", "grade": "grade", "description": "description"}}},
   {"id": "cancel-create", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'unit', 'materials:write', 'material-create')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 10: Material Edit
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000021', 'material-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Material', 'Formulario para editar un material existente',
 '{
   "page_title": "Edit Material",
   "form_title": "Edit Material",
   "form_description": "Update the material details",
   "cancel_label": "Cancel",
   "submit_label": "Save Changes",
   "form_fields": [
     {"id": "title", "controlType": "text-input", "label": "Title", "placeholder": "Enter material title", "required": true},
     {"id": "subject", "controlType": "text-input", "label": "Subject", "placeholder": "e.g. Mathematics, Science"},
     {"id": "grade", "controlType": "text-input", "label": "Grade", "placeholder": "e.g. 5th Grade"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the material content", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-edit", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "/v1/materials/{id}", "method": "PUT", "fieldMapping": {"title": "title", "subject": "subject", "grade": "grade", "description": "description"}}},
   {"id": "cancel-edit", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/materials/{id}',
 '{"method": "GET", "fieldMapping": {"title": "title", "subject": "subject", "grade": "grade", "description": "description"}}'::jsonb,
 'unit', 'materials:write', 'material-edit')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 11: Assessments List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000030', 'assessments-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Evaluaciones', 'Lista de evaluaciones disponibles para un material',
 '{
   "page_title": "Assessments",
   "search_placeholder": "Search assessments...",
   "filter_all_label": "All",
   "filter_ready_label": "Available",
   "filter_processing_label": "Completed",
   "empty_icon": "clipboard",
   "empty_state_title": "No assessments yet",
   "empty_state_description": "Assessments will appear when materials are processed",
   "empty_action_label": "Back to Materials"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "assessment-take", "params": {"materialId": "{item.material_id}", "assessmentId": "{item.id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/materials/{materialId}/assessment',
 '{"method": "GET", "fieldMapping": {"title": "title", "subtitle": "description", "status": "status", "file_type_icon": "clipboard", "created_at": "created_at", "id": "id", "material_id": "material_id"}}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 12: Assessment Take (Quiz)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000031', 'assessment-take',
 'a0000000-0000-0000-0000-000000000006', 'Rendir Evaluacion', 'Formulario para rendir una evaluacion',
 '{
   "page_title": "Assessment",
   "form_title": "Take Assessment",
   "form_description": "Answer all questions and submit when ready",
   "cancel_label": "Cancel",
   "submit_label": "Submit Answers"
 }'::jsonb,
 '[
   {"id": "submit-assessment", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "/v1/materials/{materialId}/assessment/attempts", "method": "POST"}},
   {"id": "cancel-assessment", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/materials/{materialId}/assessment',
 '{"method": "GET"}'::jsonb,
 'unit', 'assessments:read', 'assessment-take')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 13: Assessment Result
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000032', 'assessment-result',
 'a0000000-0000-0000-0000-000000000004', 'Resultado de Evaluacion', 'Resultado de un intento de evaluacion',
 '{
   "page_title": "Assessment Result",
   "file_size_label": "Score",
   "uploaded_label": "Completed at",
   "status_label": "Status",
   "description_title": "Feedback",
   "summary_title": "Details",
   "download_label": "Try Again",
   "quiz_label": "Back to Materials"
 }'::jsonb,
 '[
   {"id": "retry", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "assessment-take", "params": {"materialId": "{item.material_id}"}}},
   {"id": "back-materials", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/attempts/{id}/results',
 '{"method": "GET", "fieldMapping": {"title": "assessment_title", "subject": "material_subject", "grade": "score_display", "status": "result_status", "file_type": "clipboard", "file_size_display": "score", "description": "feedback", "created_at": "completed_at", "summary": "details"}}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 14: Attempts History
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000033', 'attempts-history',
 'a0000000-0000-0000-0000-000000000003', 'Historial de Intentos', 'Historial de intentos de evaluacion del usuario',
 '{
   "page_title": "My Attempts",
   "search_placeholder": "Search attempts...",
   "filter_all_label": "All",
   "filter_ready_label": "Passed",
   "filter_processing_label": "Failed",
   "empty_icon": "history",
   "empty_state_title": "No attempts yet",
   "empty_state_description": "Your assessment attempts will appear here",
   "empty_action_label": "Take an Assessment"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "assessment-result", "params": {"id": "{item.id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/users/me/attempts',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "assessment_title", "subtitle": "score_display", "status": "result_status", "file_type_icon": "clipboard", "created_at": "completed_at", "id": "id"}
 }'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 15: My Progress
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000040', 'progress-my',
 'a0000000-0000-0000-0000-000000000002', 'Mi Progreso', 'Dashboard de progreso personal del estudiante',
 '{
   "page_title": "My Progress",
   "greeting_text": "Your Progress, {user.firstName}",
   "date_text": "{today_date}",
   "kpi_students_label": "Completed",
   "kpi_materials_label": "In Progress",
   "kpi_avg_score_label": "Avg Score",
   "kpi_completion_label": "Overall",
   "activity_title": "Recent Progress",
   "upload_label": "View Materials",
   "progress_label": "View Attempts"
 }'::jsonb,
 '[
   {"id": "navigate-materials", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "navigate-attempts", "trigger": "button_click", "triggerSlotId": "view_progress", "type": "NAVIGATE", "config": {"target": "attempts-history"}},
   {"id": "refresh-progress", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/stats/global',
 '{"method": "GET", "fieldMapping": {"total_students": "completed_count", "total_materials": "in_progress_count", "avg_score": "avg_score", "completion_rate": "overall_completion"}}'::jsonb,
 'unit', 'progress:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 16: Progress Unit List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000041', 'progress-unit-list',
 'a0000000-0000-0000-0000-000000000003', 'Progreso por Unidad', 'Lista de progreso agrupado por unidad academica',
 '{
   "page_title": "Unit Progress",
   "search_placeholder": "Search units...",
   "filter_all_label": "All",
   "filter_ready_label": "Active",
   "filter_processing_label": "Completed",
   "empty_icon": "trending_up",
   "empty_state_title": "No progress data",
   "empty_state_description": "Progress will appear as students complete materials",
   "empty_action_label": "Back to Dashboard"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "progress-student-detail", "params": {"unitId": "{item.unit_id}", "studentId": "{item.student_id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/progress',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "unit_name", "subtitle": "completion_summary", "status": "status", "file_type_icon": "trending_up", "created_at": "last_activity", "id": "id", "unit_id": "unit_id", "student_id": "student_id"}
 }'::jsonb,
 'school', 'progress:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 17: Student Progress Detail
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000042', 'progress-student-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Progreso', 'Detalle del progreso de un estudiante',
 '{
   "page_title": "Student Progress",
   "file_size_label": "Completion",
   "uploaded_label": "Last Activity",
   "status_label": "Status",
   "description_title": "Progress Details",
   "summary_title": "Assessment Results",
   "download_label": "View Materials",
   "quiz_label": "View Attempts"
 }'::jsonb,
 '[
   {"id": "view-materials", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "view-attempts", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "attempts-history", "params": {"studentId": "{item.student_id}"}}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/progress/{studentId}',
 '{"method": "GET", "fieldMapping": {"title": "student_name", "subject": "unit_name", "grade": "completion_percentage", "status": "status", "file_type": "trending_up", "file_size_display": "completion_display", "description": "progress_details", "created_at": "last_activity", "summary": "assessment_summary"}}'::jsonb,
 'school', 'progress:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 18: Users List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000050', 'users-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Usuarios', 'Lista de usuarios del sistema',
 '{
   "page_title": "Users",
   "search_placeholder": "Search users...",
   "filter_all_label": "All",
   "filter_ready_label": "Active",
   "filter_processing_label": "Inactive",
   "empty_icon": "people",
   "empty_state_title": "No users found",
   "empty_state_description": "Add users to your school",
   "empty_action_label": "Add User"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "user-detail", "params": {"id": "{item.id}"}}},
   {"id": "create-user", "trigger": "fab_click", "type": "NAVIGATE", "config": {"target": "user-create"}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/users',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "defaultParams": {"sort": "created_at", "order": "desc"},
   "fieldMapping": {"title": "full_name", "subtitle": "email", "status": "status", "file_type_icon": "person", "created_at": "created_at", "id": "id"}
 }'::jsonb,
 'school', 'users:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 19: User Detail
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000051', 'user-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Usuario', 'Detalle de un usuario',
 '{
   "page_title": "User Detail",
   "file_size_label": "Role",
   "uploaded_label": "Created",
   "status_label": "Status",
   "description_title": "Contact Information",
   "summary_title": "Activity",
   "download_label": "Edit User",
   "quiz_label": "Back to Users"
 }'::jsonb,
 '[
   {"id": "edit-user", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "user-edit", "params": {"id": "{item.id}"}}},
   {"id": "back-users", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "users-list"}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/users/{id}',
 '{"method": "GET", "fieldMapping": {"title": "full_name", "subject": "email", "grade": "phone", "status": "status", "file_type": "person", "file_size_display": "role_name", "description": "email", "created_at": "created_at"}}'::jsonb,
 'school', 'users:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 20: User Create
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000052', 'user-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Usuario', 'Formulario para crear un nuevo usuario',
 '{
   "page_title": "New User",
   "form_title": "Create User",
   "form_description": "Fill in the details to create a new user",
   "cancel_label": "Cancel",
   "submit_label": "Create User",
   "form_fields": [
     {"id": "first_name", "controlType": "text-input", "label": "First Name", "placeholder": "Enter first name", "required": true},
     {"id": "last_name", "controlType": "text-input", "label": "Last Name", "placeholder": "Enter last name", "required": true},
     {"id": "email", "controlType": "email-input", "label": "Email", "placeholder": "Enter email address", "required": true},
     {"id": "password", "controlType": "password-input", "label": "Password", "placeholder": "Enter password", "required": true}
   ]
 }'::jsonb,
 '[
   {"id": "submit-user", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/users", "method": "POST", "fieldMapping": {"first_name": "first_name", "last_name": "last_name", "email": "email", "password": "password"}}},
   {"id": "cancel-create", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'school', 'users:write', 'user-create')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 21: User Edit
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000053', 'user-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Usuario', 'Formulario para editar un usuario',
 '{
   "page_title": "Edit User",
   "form_title": "Edit User",
   "form_description": "Update user details",
   "cancel_label": "Cancel",
   "submit_label": "Save Changes",
   "form_fields": [
     {"id": "first_name", "controlType": "text-input", "label": "First Name", "placeholder": "Enter first name", "required": true},
     {"id": "last_name", "controlType": "text-input", "label": "Last Name", "placeholder": "Enter last name", "required": true},
     {"id": "email", "controlType": "email-input", "label": "Email", "placeholder": "Enter email address", "required": true}
   ]
 }'::jsonb,
 '[
   {"id": "submit-edit", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/users/{id}", "method": "PUT", "fieldMapping": {"first_name": "first_name", "last_name": "last_name", "email": "email"}}},
   {"id": "cancel-edit", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/users/{id}',
 '{"method": "GET", "fieldMapping": {"first_name": "first_name", "last_name": "last_name", "email": "email"}}'::jsonb,
 'school', 'users:write', 'user-edit')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 22: Schools List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000060', 'schools-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Escuelas', 'Lista de escuelas del sistema',
 '{
   "page_title": "Schools",
   "search_placeholder": "Search schools...",
   "filter_all_label": "All",
   "filter_ready_label": "Active",
   "filter_processing_label": "Inactive",
   "empty_icon": "school",
   "empty_state_title": "No schools found",
   "empty_state_description": "Create a school to get started",
   "empty_action_label": "Add School"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "school-detail", "params": {"id": "{item.id}"}}},
   {"id": "create-school", "trigger": "fab_click", "type": "NAVIGATE", "config": {"target": "school-create"}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/schools',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "defaultParams": {"sort": "name", "order": "asc"},
   "fieldMapping": {"title": "name", "subtitle": "address", "status": "status", "file_type_icon": "school", "created_at": "created_at", "id": "id"}
 }'::jsonb,
 'system', 'schools:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 23: School Detail
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000061', 'school-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Escuela', 'Detalle de una escuela',
 '{
   "page_title": "School Detail",
   "file_size_label": "Students",
   "uploaded_label": "Created",
   "status_label": "Status",
   "description_title": "Information",
   "summary_title": "Units",
   "download_label": "Edit School",
   "quiz_label": "View Units"
 }'::jsonb,
 '[
   {"id": "edit-school", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "school-edit", "params": {"id": "{item.id}"}}},
   {"id": "view-units", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "units-list", "params": {"schoolId": "{item.id}"}}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/schools/{id}',
 '{"method": "GET", "fieldMapping": {"title": "name", "subject": "address", "grade": "phone", "status": "status", "file_type": "school", "file_size_display": "total_students", "description": "description", "created_at": "created_at"}}'::jsonb,
 'system', 'schools:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 24: School Create
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000062', 'school-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Escuela', 'Formulario para crear una nueva escuela',
 '{
   "page_title": "New School",
   "form_title": "Create School",
   "form_description": "Fill in the details to create a new school",
   "cancel_label": "Cancel",
   "submit_label": "Create School",
   "form_fields": [
     {"id": "name", "controlType": "text-input", "label": "School Name", "placeholder": "Enter school name", "required": true},
     {"id": "address", "controlType": "text-input", "label": "Address", "placeholder": "Enter address"},
     {"id": "phone", "controlType": "text-input", "label": "Phone", "placeholder": "Enter phone number"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the school", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-school", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/schools", "method": "POST", "fieldMapping": {"name": "name", "address": "address", "phone": "phone", "description": "description"}}},
   {"id": "cancel-create", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'system', 'schools:write', 'school-create')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 25: School Edit
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000063', 'school-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Escuela', 'Formulario para editar una escuela',
 '{
   "page_title": "Edit School",
   "form_title": "Edit School",
   "form_description": "Update school details",
   "cancel_label": "Cancel",
   "submit_label": "Save Changes",
   "form_fields": [
     {"id": "name", "controlType": "text-input", "label": "School Name", "placeholder": "Enter school name", "required": true},
     {"id": "address", "controlType": "text-input", "label": "Address", "placeholder": "Enter address"},
     {"id": "phone", "controlType": "text-input", "label": "Phone", "placeholder": "Enter phone number"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the school", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-edit", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/schools/{id}", "method": "PUT", "fieldMapping": {"name": "name", "address": "address", "phone": "phone", "description": "description"}}},
   {"id": "cancel-edit", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/schools/{id}',
 '{"method": "GET", "fieldMapping": {"name": "name", "address": "address", "phone": "phone", "description": "description"}}'::jsonb,
 'system', 'schools:write', 'school-edit')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 26: Units List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000070', 'units-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Unidades', 'Lista de unidades academicas',
 '{
   "page_title": "Academic Units",
   "search_placeholder": "Search units...",
   "filter_all_label": "All",
   "filter_ready_label": "Active",
   "filter_processing_label": "Archived",
   "empty_icon": "layers",
   "empty_state_title": "No units found",
   "empty_state_description": "Create an academic unit to organize classes",
   "empty_action_label": "Add Unit"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "unit-detail", "params": {"schoolId": "{item.school_id}", "id": "{item.id}"}}},
   {"id": "create-unit", "trigger": "fab_click", "type": "NAVIGATE", "config": {"target": "unit-create"}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/schools/{schoolId}/units',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "name", "subtitle": "description", "status": "status", "file_type_icon": "layers", "created_at": "created_at", "id": "id", "school_id": "school_id"}
 }'::jsonb,
 'school', 'units:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 27: Unit Detail
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000071', 'unit-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Unidad', 'Detalle de una unidad academica',
 '{
   "page_title": "Unit Detail",
   "file_size_label": "Members",
   "uploaded_label": "Created",
   "status_label": "Status",
   "description_title": "Description",
   "summary_title": "Members",
   "download_label": "Edit Unit",
   "quiz_label": "View Members"
 }'::jsonb,
 '[
   {"id": "edit-unit", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "unit-edit", "params": {"schoolId": "{item.school_id}", "id": "{item.id}"}}},
   {"id": "view-members", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "memberships-list", "params": {"unitId": "{item.id}"}}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/schools/{schoolId}/units/{id}',
 '{"method": "GET", "fieldMapping": {"title": "name", "subject": "school_name", "grade": "grade_level", "status": "status", "file_type": "layers", "file_size_display": "total_members", "description": "description", "created_at": "created_at"}}'::jsonb,
 'school', 'units:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 28: Unit Create
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000072', 'unit-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Unidad', 'Formulario para crear una unidad academica',
 '{
   "page_title": "New Unit",
   "form_title": "Create Academic Unit",
   "form_description": "Fill in the details to create a new academic unit",
   "cancel_label": "Cancel",
   "submit_label": "Create Unit",
   "form_fields": [
     {"id": "name", "controlType": "text-input", "label": "Unit Name", "placeholder": "e.g. 5th Grade - Section A", "required": true},
     {"id": "grade_level", "controlType": "text-input", "label": "Grade Level", "placeholder": "e.g. 5th Grade"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the unit", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-unit", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/schools/{schoolId}/units", "method": "POST", "fieldMapping": {"name": "name", "grade_level": "grade_level", "description": "description"}}},
   {"id": "cancel-create", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'school', 'units:write', 'unit-create')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 29: Unit Edit
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000073', 'unit-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Unidad', 'Formulario para editar una unidad academica',
 '{
   "page_title": "Edit Unit",
   "form_title": "Edit Academic Unit",
   "form_description": "Update unit details",
   "cancel_label": "Cancel",
   "submit_label": "Save Changes",
   "form_fields": [
     {"id": "name", "controlType": "text-input", "label": "Unit Name", "placeholder": "e.g. 5th Grade - Section A", "required": true},
     {"id": "grade_level", "controlType": "text-input", "label": "Grade Level", "placeholder": "e.g. 5th Grade"},
     {"id": "description", "controlType": "text-input", "label": "Description", "placeholder": "Describe the unit", "style": "multiline"}
   ]
 }'::jsonb,
 '[
   {"id": "submit-edit", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/schools/{schoolId}/units/{id}", "method": "PUT", "fieldMapping": {"name": "name", "grade_level": "grade_level", "description": "description"}}},
   {"id": "cancel-edit", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/schools/{schoolId}/units/{id}',
 '{"method": "GET", "fieldMapping": {"name": "name", "grade_level": "grade_level", "description": "description"}}'::jsonb,
 'school', 'units:write', 'unit-edit')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 30: Memberships List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000074', 'memberships-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Miembros', 'Lista de miembros de una unidad academica',
 '{
   "page_title": "Members",
   "search_placeholder": "Search members...",
   "filter_all_label": "All",
   "filter_ready_label": "Teachers",
   "filter_processing_label": "Students",
   "empty_icon": "user_plus",
   "empty_state_title": "No members yet",
   "empty_state_description": "Add members to this unit",
   "empty_action_label": "Add Member"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "user-detail", "params": {"id": "{item.user_id}"}}},
   {"id": "add-member", "trigger": "fab_click", "type": "NAVIGATE", "config": {"target": "membership-add", "params": {"unitId": "{unitId}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/memberships',
 '{
   "method": "GET",
   "defaultParams": {"unit_id": "{unitId}"},
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "user_full_name", "subtitle": "role_name", "status": "status", "file_type_icon": "person", "created_at": "joined_at", "id": "id", "user_id": "user_id"}
 }'::jsonb,
 'school', 'memberships:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 31: Membership Add
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000075', 'membership-add',
 'a0000000-0000-0000-0000-000000000006', 'Agregar Miembro', 'Formulario para agregar un miembro a la unidad',
 '{
   "page_title": "Add Member",
   "form_title": "Add Member to Unit",
   "form_description": "Select a user and role to add to this unit",
   "cancel_label": "Cancel",
   "submit_label": "Add Member",
   "form_fields": [
     {"id": "user_email", "controlType": "email-input", "label": "User Email", "placeholder": "Enter user email", "required": true},
     {"id": "role", "controlType": "text-input", "label": "Role", "placeholder": "e.g. teacher, student", "required": true}
   ]
 }'::jsonb,
 '[
   {"id": "submit-membership", "trigger": "button_click", "triggerSlotId": "submit_btn", "type": "SUBMIT_FORM", "config": {"endpoint": "admin:/v1/memberships", "method": "POST", "fieldMapping": {"user_email": "user_email", "role": "role", "unit_id": "{unitId}"}}},
   {"id": "cancel-add", "trigger": "button_click", "triggerSlotId": "cancel_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 NULL, '{}'::jsonb, 'school', 'memberships:write', 'membership-add')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 32: Dashboard Guardian
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000080', 'dashboard-guardian',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Guardian', 'Panel principal del apoderado/guardian',
 '{
   "page_title": "Family Dashboard",
   "greeting_text": "Hello, {user.firstName}",
   "date_text": "{today_date}",
   "kpi_students_label": "Children",
   "kpi_materials_label": "Activities",
   "kpi_avg_score_label": "Avg Score",
   "kpi_completion_label": "Overall Progress",
   "activity_title": "Recent Activity",
   "upload_label": "My Children",
   "progress_label": "View Progress"
 }'::jsonb,
 '[
   {"id": "navigate-children", "trigger": "button_click", "triggerSlotId": "upload_material", "type": "NAVIGATE", "config": {"target": "children-list"}},
   {"id": "navigate-progress", "trigger": "button_click", "triggerSlotId": "view_progress", "type": "NAVIGATE", "config": {"target": "child-progress"}},
   {"id": "refresh-dashboard", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/guardians/me/stats',
 '{"method": "GET", "fieldMapping": {"total_students": "total_children", "total_materials": "total_activities", "avg_score": "avg_score", "completion_rate": "overall_progress"}}'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 33: Children List (Guardian view)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000081', 'children-list',
 'a0000000-0000-0000-0000-000000000003', 'Mis Hijos', 'Lista de hijos/estudiantes vinculados al guardian',
 '{
   "page_title": "My Children",
   "search_placeholder": "Search...",
   "filter_all_label": "All",
   "filter_ready_label": "Active",
   "filter_processing_label": "Inactive",
   "empty_icon": "people",
   "empty_state_title": "No children linked",
   "empty_state_description": "Contact your school to link your children",
   "empty_action_label": "Back to Dashboard"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "child-progress", "params": {"childId": "{item.id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/guardians/me/relations',
 '{
   "method": "GET",
   "fieldMapping": {"title": "child_full_name", "subtitle": "school_name", "status": "status", "file_type_icon": "person", "created_at": "enrolled_at", "id": "child_id"}
 }'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 34: Child Progress (Guardian view)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000082', 'child-progress',
 'a0000000-0000-0000-0000-000000000004', 'Progreso del Hijo', 'Vista de progreso de un hijo para el guardian',
 '{
   "page_title": "Child Progress",
   "file_size_label": "Completion",
   "uploaded_label": "Last Activity",
   "status_label": "Status",
   "description_title": "Progress Details",
   "summary_title": "Recent Scores",
   "download_label": "View Materials",
   "quiz_label": "View Attempts"
 }'::jsonb,
 '[
   {"id": "view-materials", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "materials-list"}},
   {"id": "view-attempts", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "attempts-history", "params": {"studentId": "{item.child_id}"}}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 '/v1/guardians/me/children/{childId}/progress',
 '{"method": "GET", "fieldMapping": {"title": "child_name", "subject": "school_name", "grade": "completion_percentage", "status": "status", "file_type": "trending_up", "file_size_display": "completion_display", "description": "progress_summary", "created_at": "last_activity", "summary": "recent_scores"}}'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 35: Roles List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000090', 'roles-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Roles', 'Lista de roles del sistema',
 '{
   "page_title": "Roles",
   "search_placeholder": "Search roles...",
   "filter_all_label": "All",
   "filter_ready_label": "System",
   "filter_processing_label": "Custom",
   "empty_icon": "shield",
   "empty_state_title": "No roles found",
   "empty_state_description": "Roles will be configured by the system",
   "empty_action_label": "Back to Admin"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "role-detail", "params": {"id": "{item.id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/roles',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 20, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "display_name", "subtitle": "description", "status": "scope", "file_type_icon": "shield", "created_at": "created_at", "id": "id"}
 }'::jsonb,
 'system', 'roles:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 36: Role Detail
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000091', 'role-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Rol', 'Detalle de un rol con sus permisos',
 '{
   "page_title": "Role Detail",
   "file_size_label": "Permissions",
   "uploaded_label": "Created",
   "status_label": "Scope",
   "description_title": "Description",
   "summary_title": "Assigned Permissions",
   "download_label": "View Permissions",
   "quiz_label": "Back to Roles"
 }'::jsonb,
 '[
   {"id": "view-permissions", "trigger": "button_click", "triggerSlotId": "download_btn", "type": "NAVIGATE", "config": {"target": "permissions-list"}},
   {"id": "back-roles", "trigger": "button_click", "triggerSlotId": "take_quiz_btn", "type": "NAVIGATE", "config": {"target": "roles-list"}},
   {"id": "go-back", "trigger": "button_click", "triggerSlotId": "back_btn", "type": "NAVIGATE_BACK"}
 ]'::jsonb,
 'admin:/v1/roles/{id}',
 '{"method": "GET", "fieldMapping": {"title": "display_name", "subject": "key", "grade": "scope", "status": "scope", "file_type": "shield", "file_size_display": "permissions_count", "description": "description", "created_at": "created_at", "summary": "permissions_summary"}}'::jsonb,
 'system', 'roles:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 37: Resources List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000092', 'resources-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Recursos', 'Lista de recursos del sistema para RBAC',
 '{
   "page_title": "Resources",
   "search_placeholder": "Search resources...",
   "filter_all_label": "All",
   "filter_ready_label": "Visible",
   "filter_processing_label": "Hidden",
   "empty_icon": "key",
   "empty_state_title": "No resources found",
   "empty_state_description": "System resources are managed automatically",
   "empty_action_label": "Back to Admin"
 }'::jsonb,
 '[
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/resources',
 '{
   "method": "GET",
   "fieldMapping": {"title": "display_name", "subtitle": "description", "status": "scope", "file_type_icon": "key", "created_at": "key", "id": "id"}
 }'::jsonb,
 'system', 'permissions:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 38: Permissions List
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000093', 'permissions-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Permisos', 'Lista de permisos del sistema',
 '{
   "page_title": "Permissions",
   "search_placeholder": "Search permissions...",
   "filter_all_label": "All",
   "filter_ready_label": "Read",
   "filter_processing_label": "Write",
   "empty_icon": "key",
   "empty_state_title": "No permissions found",
   "empty_state_description": "Permissions are managed by the system",
   "empty_action_label": "Back to Admin"
 }'::jsonb,
 '[
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 'admin:/v1/permissions',
 '{
   "method": "GET",
   "pagination": {"type": "offset", "pageSize": 50, "pageParam": "offset", "limitParam": "limit"},
   "fieldMapping": {"title": "display_name", "subtitle": "key", "status": "resource_name", "file_type_icon": "key", "created_at": "key", "id": "id"}
 }'::jsonb,
 'system', 'permissions:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;
