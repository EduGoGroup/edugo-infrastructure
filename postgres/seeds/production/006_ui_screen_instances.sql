-- ====================================================================
-- SEEDS: Instancias de pantalla para las pantallas configuradas
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ====================================================================

BEGIN;

-- Instancia 1: Login
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000001', 'app-login',
 'a0000000-0000-0000-0000-000000000001', 'Login', 'Pantalla de inicio de sesion',
 '{"app_logo": "edugo_logo", "app_name": "EduGo", "app_tagline": "Learning made easy", "email_label": "Email", "password_label": "Password", "remember_label": "Remember me", "login_btn_label": "Sign In", "forgot_password_label": "Forgot password?", "divider_text": "or continue with", "google_btn_label": "Google"}'::jsonb,
 'system', NULL, 'login')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 2: Dashboard Profesor
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000002', 'dashboard-teacher',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Profesor', 'Panel principal del profesor',
 '{"page_title": "Dashboard", "greeting_text": "Good morning, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Students", "kpi_materials_label": "Materials", "kpi_avg_score_label": "Avg Score", "kpi_completion_label": "Completion", "activity_title": "Recent Activity", "upload_label": "Upload Material", "progress_label": "View Progress"}'::jsonb,
 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 3: Dashboard Estudiante
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000003', 'dashboard-student',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Estudiante', 'Panel principal del estudiante',
 '{"page_title": "Home", "greeting_text": "Hello, {user.firstName}!", "date_text": "{today_date}", "kpi_students_label": "Courses", "kpi_materials_label": "Materials", "kpi_avg_score_label": "My Score", "kpi_completion_label": "Progress", "activity_title": "Recent Activity", "upload_label": "My Materials", "progress_label": "My Progress"}'::jsonb,
 'unit', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 4: Lista de Materiales
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000004', 'materials-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Materiales', 'Lista de materiales educativos',
 '{"page_title": "Materials", "search_placeholder": "Search materials...", "filter_all_label": "All", "filter_ready_label": "Ready", "filter_processing_label": "Processing", "empty_icon": "folder_open", "empty_state_title": "No materials yet", "empty_state_description": "Upload your first educational material", "empty_action_label": "Upload Material"}'::jsonb,
 'unit', 'materials:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 5: Detalle de Material
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000005', 'material-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Material', 'Detalle de un material educativo',
 '{"page_title": "Material Detail", "file_size_label": "File Size", "uploaded_label": "Uploaded", "status_label": "Status", "description_title": "Description", "summary_title": "AI Summary", "download_label": "Download", "quiz_label": "Take Quiz"}'::jsonb,
 'unit', 'materials:read', 'material-detail')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 6: Configuracion
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000006', 'app-settings',
 'a0000000-0000-0000-0000-000000000005', 'Configuracion', 'Pantalla de configuracion',
 '{"page_title": "Settings", "appearance_title": "Appearance", "dark_mode_label": "Dark Mode", "theme_label": "Theme Color", "notifications_title": "Notifications", "push_label": "Push Notifications", "email_label": "Email Notifications", "logout_label": "Sign Out"}'::jsonb,
 'system', NULL, 'settings')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 7: Dashboard Superadmin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000010', 'dashboard-superadmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Superadmin', 'Panel principal del superadmin',
 '{"page_title": "System Dashboard", "greeting_text": "Welcome, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Total Schools", "kpi_materials_label": "Total Users", "kpi_avg_score_label": "Total Materials", "kpi_completion_label": "System Health", "activity_title": "System Activity", "upload_label": "Manage Schools", "progress_label": "View Stats"}'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 8: Dashboard School Admin
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000011', 'dashboard-schooladmin',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard School Admin', 'Panel principal del admin de escuela',
 '{"page_title": "School Dashboard", "greeting_text": "Welcome, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Teachers", "kpi_materials_label": "Students", "kpi_avg_score_label": "Units", "kpi_completion_label": "School Score", "activity_title": "School Activity", "upload_label": "Manage Users", "progress_label": "View Reports"}'::jsonb,
 'school', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 9: Dashboard Guardian
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000080', 'dashboard-guardian',
 'a0000000-0000-0000-0000-000000000002', 'Dashboard Guardian', 'Panel principal del apoderado',
 '{"page_title": "Family Dashboard", "greeting_text": "Hello, {user.firstName}", "date_text": "{today_date}", "kpi_students_label": "Children", "kpi_materials_label": "Activities", "kpi_avg_score_label": "Avg Score", "kpi_completion_label": "Overall Progress", "activity_title": "Recent Activity", "upload_label": "My Children", "progress_label": "View Progress"}'::jsonb,
 'system', NULL, NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 10: Lista de Usuarios
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000020', 'users-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Usuarios', 'Lista de usuarios del sistema',
 '{"page_title": "Usuarios", "search_placeholder": "Buscar usuario...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "users", "empty_state_title": "No hay usuarios", "empty_state_description": "No se encontraron usuarios", "empty_action_label": "Crear Usuario"}'::jsonb,
 'school', 'users:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 11: Lista de Escuelas
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000021', 'schools-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Escuelas', 'Lista de escuelas del sistema',
 '{"page_title": "Escuelas", "search_placeholder": "Buscar escuela...", "filter_all_label": "Todas", "filter_ready_label": "Activas", "filter_processing_label": "Inactivas", "empty_icon": "school", "empty_state_title": "No hay escuelas", "empty_state_description": "No se encontraron escuelas", "empty_action_label": "Crear Escuela"}'::jsonb,
 'system', 'schools:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 12: Lista de Roles
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000022', 'roles-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Roles', 'Lista de roles del sistema',
 '{"page_title": "Roles", "search_placeholder": "Buscar rol...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "shield", "empty_state_title": "No hay roles", "empty_state_description": "No se encontraron roles", "empty_action_label": "Crear Rol"}'::jsonb,
 'system', 'roles:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 13: Lista de Permisos
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000023', 'permissions-list',
 'a0000000-0000-0000-0000-000000000003', 'Gestion de Permisos', 'Lista de permisos del sistema',
 '{"page_title": "Permisos", "search_placeholder": "Buscar permiso...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "key", "empty_state_title": "No hay permisos", "empty_state_description": "No se encontraron permisos", "empty_action_label": ""}'::jsonb,
 'system', 'permissions_mgmt:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 14: Unidades Academicas
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000024', 'units-list',
 'a0000000-0000-0000-0000-000000000003', 'Unidades Academicas', 'Lista de unidades academicas',
 '{"page_title": "Unidades Académicas", "search_placeholder": "Buscar unidad...", "filter_all_label": "Todas", "filter_ready_label": "Activas", "filter_processing_label": "Inactivas", "empty_icon": "layers", "empty_state_title": "No hay unidades", "empty_state_description": "No se encontraron unidades académicas", "empty_action_label": "Crear Unidad"}'::jsonb,
 'school', 'units:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 15: Lista de Miembros
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000025', 'memberships-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Miembros', 'Lista de miembros de unidades',
 '{"page_title": "Miembros", "search_placeholder": "Buscar miembro...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "user-plus", "empty_state_title": "No hay miembros", "empty_state_description": "No se encontraron miembros asignados", "empty_action_label": "Asignar Miembro"}'::jsonb,
 'school', 'memberships:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 16: Lista de Evaluaciones
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000026', 'assessments-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Evaluaciones', 'Lista de evaluaciones',
 '{"page_title": "Evaluaciones", "search_placeholder": "Buscar evaluación...", "filter_all_label": "Todas", "filter_ready_label": "Publicadas", "filter_processing_label": "Borradores", "empty_icon": "clipboard", "empty_state_title": "No hay evaluaciones", "empty_state_description": "No se encontraron evaluaciones", "empty_action_label": "Crear Evaluación"}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 17: Dashboard Progreso
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000027', 'progress-dashboard',
 'a0000000-0000-0000-0000-000000000002', 'Progreso Academico', 'Dashboard de progreso academico',
 '{"page_title": "Progreso", "greeting_text": "Progreso Académico", "date_text": "{today_date}", "kpi_students_label": "Estudiantes", "kpi_materials_label": "Completados", "kpi_avg_score_label": "Promedio", "kpi_completion_label": "Avance", "activity_title": "Actividad Reciente", "upload_label": "Ver Detalle", "progress_label": "Exportar"}'::jsonb,
 'unit', 'progress:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 18: Dashboard Estadisticas
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000028', 'stats-dashboard',
 'a0000000-0000-0000-0000-000000000002', 'Estadisticas', 'Dashboard de estadisticas del sistema',
 '{"page_title": "Estadísticas", "greeting_text": "Estadísticas del Sistema", "date_text": "{today_date}", "kpi_students_label": "Usuarios", "kpi_materials_label": "Materiales", "kpi_avg_score_label": "Evaluaciones", "kpi_completion_label": "Escuelas", "activity_title": "Resumen", "upload_label": "Ver Detalle", "progress_label": "Exportar"}'::jsonb,
 'school', 'stats:school', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 19: Lista de Materias (CRUD ejemplo)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000029', 'subjects-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Materias', 'Lista de materias del plan de estudios',
 '{"page_title": "Materias", "search_placeholder": "Buscar materia...", "filter_all_label": "Todas", "filter_ready_label": "Activas", "filter_processing_label": "Inactivas", "empty_icon": "book", "empty_state_title": "No hay materias registradas", "empty_state_description": "Crea la primera materia del plan de estudios", "empty_action_label": "Crear Materia", "columns": ["Nombre", "Descripción", "Estado"]}'::jsonb,
 'school', 'subjects:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 20: Formulario de Materia
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000030', 'subjects-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Materia', 'Formulario para crear/editar materias',
 '{"page_title": "Nueva Materia", "edit_title": "Editar Materia", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la materia", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la materia", "required": false}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true, "required": false}]}'::jsonb,
 'school', 'subjects:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 21: Formulario de Escuela
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000031', 'schools-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Escuela', 'Crear o editar una escuela',
 '{"page_title": "Nueva Escuela", "edit_title": "Editar Escuela", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la escuela", "required": true}, {"key": "code", "type": "text", "label": "Código", "placeholder": "Código único (ej: ESC001)", "required": true}, {"key": "address", "type": "text", "label": "Dirección", "placeholder": "Dirección de la escuela", "required": true}, {"key": "city", "type": "text", "label": "Ciudad", "placeholder": "Ciudad", "required": true}, {"key": "country", "type": "text", "label": "País", "placeholder": "País", "required": true}, {"key": "contact_email", "type": "email", "label": "Email de contacto", "placeholder": "email@escuela.com", "required": true}, {"key": "contact_phone", "type": "text", "label": "Teléfono", "placeholder": "Teléfono de contacto", "required": false}, {"key": "subscription_tier", "type": "text", "label": "Plan", "placeholder": "basic / premium / enterprise", "required": true}, {"key": "max_teachers", "type": "number", "label": "Máx. Profesores", "placeholder": "10", "required": true}, {"key": "max_students", "type": "number", "label": "Máx. Estudiantes", "placeholder": "100", "required": true}]}'::jsonb,
 'system', 'schools:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 22: Formulario de Usuario
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000032', 'users-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Usuario', 'Crear o editar un usuario',
 '{"page_title": "Nuevo Usuario", "edit_title": "Editar Usuario", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "first_name", "type": "text", "label": "Nombre", "placeholder": "Nombre del usuario", "required": true}, {"key": "last_name", "type": "text", "label": "Apellido", "placeholder": "Apellido del usuario", "required": true}, {"key": "email", "type": "email", "label": "Email", "placeholder": "email@ejemplo.com", "required": true}, {"key": "password", "type": "password", "label": "Contraseña", "placeholder": "Mínimo 8 caracteres", "required": true}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true, "required": false}]}'::jsonb,
 'school', 'users:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 23: Formulario de Rol
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000033', 'roles-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Rol', 'Crear o editar un rol del sistema',
 '{"page_title": "Nuevo Rol", "edit_title": "Editar Rol", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "name", "type": "text", "label": "Nombre clave", "placeholder": "ej: school_coordinator", "required": true}, {"key": "display_name", "type": "text", "label": "Nombre visible", "placeholder": "Nombre para mostrar", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del rol"}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true}]}'::jsonb,
 'system', 'roles:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 24: Formulario de Permiso
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000034', 'permissions-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Permiso', 'Crear o editar un permiso del sistema',
 '{"page_title": "Nuevo Permiso", "edit_title": "Editar Permiso", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "name", "type": "text", "label": "Nombre (resource:action)", "placeholder": "ej: users:create", "required": true}, {"key": "display_name", "type": "text", "label": "Nombre visible", "placeholder": "Nombre para mostrar", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción"}, {"key": "resource_id", "type": "remote_select", "label": "Recurso", "required": true, "remote_endpoint": "/api/v1/resources", "display_field": "display_name", "value_field": "id"}, {"key": "action", "type": "text", "label": "Acción", "placeholder": "ej: create, read, update, delete", "required": true}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true}]}'::jsonb,
 'system', 'permissions_mgmt:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 25: Formulario de Unidad Académica (Fase 2.1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000040', 'units-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Unidad Academica', 'Crear o editar unidad academica',
 '{"page_title": "Nueva Unidad", "edit_title": "Editar Unidad", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la unidad", "required": true}, {"key": "code", "type": "text", "label": "Código", "placeholder": "Código único", "required": false}, {"key": "type", "type": "select", "label": "Tipo", "required": true, "options": [{"value": "grade", "label": "Grado"}, {"value": "class", "label": "Clase"}, {"value": "section", "label": "Sección"}, {"value": "department", "label": "Departamento"}, {"value": "club", "label": "Club"}]}, {"key": "academic_year", "type": "number", "label": "Año académico", "placeholder": "2026", "required": false}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la unidad", "required": false}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true}]}'::jsonb,
 'school', 'units:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 26: Formulario de Miembros (Fase 2.2)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000041', 'memberships-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Miembro', 'Asignar miembro a unidad academica',
 '{"page_title": "Asignar Miembro", "edit_title": "Editar Miembro", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "user_id", "type": "remote_select", "label": "Usuario", "required": true, "remote_endpoint": "/api/v1/users", "display_field": "email", "value_field": "id"}, {"key": "academic_unit_id", "type": "remote_select", "label": "Unidad Académica", "required": true, "remote_endpoint": "/api/v1/schools/{schoolId}/units", "display_field": "name", "value_field": "id"}, {"key": "role", "type": "select", "label": "Rol", "required": true, "options": [{"value": "student", "label": "Estudiante"}, {"value": "teacher", "label": "Profesor"}, {"value": "assistant", "label": "Asistente"}]}, {"key": "enrolled_at", "type": "text", "label": "Fecha de inscripción", "placeholder": "YYYY-MM-DD", "required": false}]}'::jsonb,
 'school', 'memberships:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 27: Formulario de Crear Material (Fase 2.3)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000042', 'material-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Material', 'Formulario para crear material educativo',
 '{"page_title": "Nuevo Material", "submit_label": "Crear", "cancel_label": "Cancelar", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Título del material", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del material", "required": false}, {"key": "subject", "type": "text", "label": "Materia", "placeholder": "Materia relacionada", "required": false}, {"key": "grade", "type": "text", "label": "Grado", "placeholder": "Grado o nivel", "required": false}, {"key": "is_public", "type": "toggle", "label": "Público", "default": false}]}'::jsonb,
 'unit', 'materials:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 28: Formulario de Editar Material (Fase 2.3)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000043', 'material-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Material', 'Formulario para editar material educativo',
 '{"page_title": "Editar Material", "edit_title": "Editar Material", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Título del material", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del material", "required": false}, {"key": "subject", "type": "text", "label": "Materia", "placeholder": "Materia relacionada", "required": false}, {"key": "grade", "type": "text", "label": "Grado", "placeholder": "Grado o nivel", "required": false}, {"key": "is_public", "type": "toggle", "label": "Público", "default": false}]}'::jsonb,
 'unit', 'materials:update', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 29: Lista de Templates de Pantalla (Fase 5.5)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000050', 'screen-templates-list',
 'a0000000-0000-0000-0000-000000000003', 'Templates de Pantalla', 'Lista de templates base del sistema SDUI',
 '{"page_title": "Templates de Pantalla", "search_placeholder": "Buscar template...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "settings_applications", "empty_state_title": "No hay templates", "empty_state_description": "No se encontraron templates de pantalla", "empty_action_label": ""}'::jsonb,
 'system', 'screen_templates:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 30: Lista de Instancias de Pantalla (Fase 5.5)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000051', 'screen-instances-list',
 'a0000000-0000-0000-0000-000000000003', 'Instancias de Pantalla', 'Lista de instancias configuradas de pantalla',
 '{"page_title": "Instancias de Pantalla", "search_placeholder": "Buscar instancia...", "filter_all_label": "Todas", "filter_ready_label": "Activas", "filter_processing_label": "Inactivas", "empty_icon": "devices", "empty_state_title": "No hay instancias", "empty_state_description": "No se encontraron instancias de pantalla", "empty_action_label": "Crear Instancia"}'::jsonb,
 'system', 'screen_instances:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 31: Formulario de Instancia de Pantalla (Fase 5.5)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000052', 'screen-instances-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Instancia', 'Editar configuracion de instancia de pantalla',
 '{"page_title": "Nueva Instancia", "edit_title": "Editar Instancia", "submit_label": "Guardar", "cancel_label": "Cancelar", "fields": [{"key": "screen_key", "type": "text", "label": "Screen Key", "placeholder": "ej: my-screen-list", "required": true}, {"key": "template_id", "type": "remote_select", "label": "Template", "required": true, "remote_endpoint": "/api/v1/screen-config/templates", "display_field": "name", "value_field": "id"}, {"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la instancia", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción"}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "required_permission", "type": "text", "label": "Permiso requerido", "placeholder": "ej: screens:read"}, {"key": "handler_key", "type": "text", "label": "Handler Key", "placeholder": "Clave del handler"}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true}]}'::jsonb,
 'system', 'screen_instances:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

COMMIT;
