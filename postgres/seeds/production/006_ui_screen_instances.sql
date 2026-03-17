-- ====================================================================
-- SEEDS: Instancias de pantalla para las pantallas configuradas
-- Idempotente: usa ON CONFLICT (DO NOTHING o DO UPDATE SET)
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
 '{"page_title": "Materials", "search_placeholder": "Search materials...", "filter_all_label": "All", "filter_ready_label": "Ready", "filter_processing_label": "Processing", "empty_icon": "folder_open", "empty_state_title": "No materials yet", "empty_state_description": "Upload your first educational material", "empty_action_label": "Upload Material", "data_endpoint": "mobile:/api/v1/materials"}'::jsonb,
 'unit', 'materials:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

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
 '{"page_title": "Permisos", "search_placeholder": "Buscar permiso...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "key", "empty_state_title": "No hay permisos", "empty_state_description": "No se encontraron permisos"}'::jsonb,
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
 '{"page_title": "Evaluaciones", "search_placeholder": "Buscar evaluación...", "filter_all_label": "Todas", "filter_ready_label": "Publicadas", "filter_processing_label": "Borradores", "empty_icon": "clipboard", "empty_state_title": "No hay evaluaciones", "empty_state_description": "No se encontraron evaluaciones"}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 17: Dashboard Progreso
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000027', 'progress-dashboard',
 'a0000000-0000-0000-0000-000000000002', 'Progreso Academico', 'Dashboard de progreso academico',
 '{"page_title": "Progreso", "greeting_text": "Progreso Académico", "date_text": "{today_date}", "kpi_students_label": "Estudiantes", "kpi_materials_label": "Completados", "kpi_avg_score_label": "Promedio", "kpi_completion_label": "Avance", "activity_title": "Actividad Reciente", "upload_label": "Ver Detalle", "progress_label": "Exportar", "data_endpoint": "admin:/api/v1/stats/global"}'::jsonb,
 'unit', 'progress:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

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
 '{"page_title": "Nueva Materia", "edit_title": "Editar Materia", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la materia", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la materia", "required": false}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true, "required": false}]}'::jsonb,
 'school', 'subjects:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 21: Formulario de Escuela
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000031', 'schools-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Escuela', 'Crear o editar una escuela',
 '{"page_title": "Nueva Escuela", "edit_title": "Editar Escuela", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "concept_type_id", "type": "remote_select", "label": "Tipo de institución", "placeholder": "Seleccione un tipo", "required": true, "options_endpoint": "admin:/api/v1/concept-types", "option_label": "name", "option_value": "id"}, {"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la escuela", "required": true}, {"key": "code", "type": "text", "label": "Código", "placeholder": "Código único (ej: ESC001)", "required": true}, {"key": "address", "type": "text", "label": "Dirección", "placeholder": "Dirección de la escuela", "required": true}, {"key": "city", "type": "text", "label": "Ciudad", "placeholder": "Ciudad", "required": true}, {"key": "country", "type": "text", "label": "País", "placeholder": "País", "required": true}, {"key": "contact_email", "type": "email", "label": "Email de contacto", "placeholder": "email@escuela.com", "required": true}, {"key": "contact_phone", "type": "text", "label": "Teléfono", "placeholder": "Teléfono de contacto", "required": false}, {"key": "subscription_tier", "type": "select", "label": "Plan", "placeholder": "Seleccione un plan", "required": true, "options": [{"label": "Gratuito", "value": "free"}, {"label": "Básico", "value": "basic"}, {"label": "Premium", "value": "premium"}, {"label": "Empresarial", "value": "enterprise"}]}, {"key": "max_teachers", "type": "number", "label": "Máx. Profesores", "placeholder": "10", "required": true}, {"key": "max_students", "type": "number", "label": "Máx. Estudiantes", "placeholder": "100", "required": true}]}'::jsonb,
 'system', 'schools:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 22: Formulario de Usuario
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000032', 'users-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Usuario', 'Crear o editar un usuario',
 '{"page_title": "Nuevo Usuario", "edit_title": "Editar Usuario", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "first_name", "type": "text", "label": "Nombre", "placeholder": "Nombre del usuario", "required": true}, {"key": "last_name", "type": "text", "label": "Apellido", "placeholder": "Apellido del usuario", "required": true}, {"key": "email", "type": "email", "label": "Email", "placeholder": "email@ejemplo.com", "required": true}, {"key": "password", "type": "password", "label": "Contraseña", "placeholder": "Mínimo 8 caracteres", "required": true, "condition": "create-only"}, {"key": "confirm_password", "type": "password", "label": "Confirmar Contraseña", "placeholder": "Repita la contraseña", "required": true, "condition": "create-only"}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true, "required": false}]}'::jsonb,
 'school', 'users:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    template_id         = EXCLUDED.template_id,
    name                = EXCLUDED.name,
    description         = EXCLUDED.description,
    slot_data           = EXCLUDED.slot_data,
    scope               = EXCLUDED.scope,
    required_permission = EXCLUDED.required_permission,
    handler_key         = EXCLUDED.handler_key;

-- Instancia 23: Formulario de Rol
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000033', 'roles-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Rol', 'Crear o editar un rol del sistema',
 '{"page_title": "Nuevo Rol", "edit_title": "Editar Rol", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre clave", "placeholder": "ej: school_coordinator", "required": true}, {"key": "display_name", "type": "text", "label": "Nombre visible", "placeholder": "Nombre para mostrar", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del rol"}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true}]}'::jsonb,
 'system', 'roles:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 24: Formulario de Permiso
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000034', 'permissions-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Permiso', 'Crear o editar un permiso del sistema',
 '{"page_title": "Nuevo Permiso", "edit_title": "Editar Permiso", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre (resource:action)", "placeholder": "ej: users:create", "required": true}, {"key": "display_name", "type": "text", "label": "Nombre visible", "placeholder": "Nombre para mostrar", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción"}, {"key": "resource_id", "type": "remote_select", "label": "Recurso", "required": true, "remote_endpoint": "/api/v1/resources", "display_field": "display_name", "value_field": "id"}, {"key": "action", "type": "text", "label": "Acción", "placeholder": "ej: create, read, update, delete", "required": true}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true}]}'::jsonb,
 'system', 'permissions_mgmt:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 25: Formulario de Unidad Académica (Fase 2.1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000040', 'units-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Unidad Academica', 'Crear o editar unidad academica',
 '{"page_title": "Nueva Unidad", "edit_title": "Editar Unidad", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la unidad", "required": true}, {"key": "code", "type": "text", "label": "Código", "placeholder": "Código único", "required": false}, {"key": "type", "type": "select", "label": "Tipo", "required": true, "options": [{"value": "grade", "label": "Grado"}, {"value": "class", "label": "Clase"}, {"value": "section", "label": "Sección"}, {"value": "department", "label": "Departamento"}, {"value": "club", "label": "Club"}]}, {"key": "academic_year", "type": "number", "label": "Año académico", "placeholder": "2026", "required": false}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la unidad", "required": false}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true}]}'::jsonb,
 'school', 'units:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 26: Formulario de Miembros (Fase 2.2)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000041', 'memberships-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Miembro', 'Asignar miembro a unidad academica',
 '{"page_title": "Asignar Miembro", "edit_title": "Editar Miembro", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "user_id", "type": "remote_select", "label": "Usuario", "required": true, "remote_endpoint": "admin:/api/v1/users", "display_field": "email", "value_field": "id"}, {"key": "academic_unit_id", "type": "remote_select", "label": "Unidad Académica", "required": true, "remote_endpoint": "/api/v1/schools/{schoolId}/units", "display_field": "name", "value_field": "id"}, {"key": "role", "type": "select", "label": "Rol", "required": true, "options": [{"value": "student", "label": "Estudiante"}, {"value": "teacher", "label": "Profesor"}, {"value": "assistant", "label": "Asistente"}]}, {"key": "enrolled_at", "type": "text", "label": "Fecha de inscripción", "placeholder": "YYYY-MM-DD", "required": false}]}'::jsonb,
 'school', 'memberships:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia 27: Formulario de Crear Material (Fase 2.3)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000042', 'material-create',
 'a0000000-0000-0000-0000-000000000006', 'Crear Material', 'Formulario para crear material educativo',
 '{"page_title": "Nuevo Material", "submit_label": "Crear", "delete_label": "Eliminar", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Título del material", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del material", "required": false}, {"key": "subject", "type": "text", "label": "Materia", "placeholder": "Materia relacionada", "required": false}, {"key": "grade", "type": "text", "label": "Grado", "placeholder": "Grado o nivel", "required": false}, {"key": "is_public", "type": "toggle", "label": "Público", "default": false}]}'::jsonb,
 'unit', 'materials:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 28: Formulario de Editar Material (Fase 2.3)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000043', 'material-edit',
 'a0000000-0000-0000-0000-000000000006', 'Editar Material', 'Formulario para editar material educativo',
 '{"page_title": "Editar Material", "edit_title": "Editar Material", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Título del material", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del material", "required": false}, {"key": "subject", "type": "text", "label": "Materia", "placeholder": "Materia relacionada", "required": false}, {"key": "grade", "type": "text", "label": "Grado", "placeholder": "Grado o nivel", "required": false}, {"key": "is_public", "type": "toggle", "label": "Público", "default": false}]}'::jsonb,
 'unit', 'materials:update', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 29: Lista de Templates de Pantalla (Fase 5.5)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000050', 'screen-templates-list',
 'a0000000-0000-0000-0000-000000000003', 'Templates de Pantalla', 'Lista de templates base del sistema SDUI',
 '{"page_title": "Templates de Pantalla", "search_placeholder": "Buscar template...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "settings_applications", "empty_state_title": "No hay templates", "empty_state_description": "No se encontraron templates de pantalla"}'::jsonb,
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
 '{"page_title": "Nueva Instancia", "edit_title": "Editar Instancia", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "screen_key", "type": "text", "label": "Screen Key", "placeholder": "ej: my-screen-list", "required": true}, {"key": "template_id", "type": "remote_select", "label": "Template", "required": true, "remote_endpoint": "/api/v1/screen-config/templates", "display_field": "name", "value_field": "id"}, {"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la instancia", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción"}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "required_permission", "type": "text", "label": "Permiso requerido", "placeholder": "ej: screens:read"}, {"key": "handler_key", "type": "text", "label": "Handler Key", "placeholder": "Clave del handler"}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true}]}'::jsonb,
 'system', 'screen_instances:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 32: Lista de Hijos (Guardian - Fase 4.1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000060', 'children-list',
 'a0000000-0000-0000-0000-000000000003', 'Mis Hijos', 'Lista de hijos vinculados al guardian',
 '{"page_title": "Mis Hijos", "search_placeholder": "Buscar hijo...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Pendientes", "empty_icon": "users", "empty_state_title": "No hay hijos vinculados", "empty_state_description": "Aún no tienes estudiantes vinculados", "empty_action_label": "Solicitar Vínculo"}'::jsonb,
 'school', 'guardian_relations:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 33: Progreso del Hijo (Guardian - Fase 4.1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000061', 'child-progress',
 'a0000000-0000-0000-0000-000000000004', 'Progreso del Hijo', 'Detalle de progreso academico del hijo',
 '{"page_title": "Progreso", "file_size_label": "Promedio", "uploaded_label": "Última actividad", "status_label": "Estado", "description_title": "Resumen", "summary_title": "Progreso Detallado", "download_label": "Ver Materiales", "quiz_label": "Ver Evaluaciones"}'::jsonb,
 'school', 'guardian_relations:read', 'child-progress')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 34: Lista de Solicitudes de Vinculación Guardian (Fase 4.1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000062', 'guardian-requests-list',
 'a0000000-0000-0000-0000-000000000003', 'Solicitudes de Vinculación', 'Lista de solicitudes de vinculación guardian-estudiante',
 '{"page_title": "Solicitudes de Vinculación", "search_placeholder": "Buscar solicitud...", "filter_all_label": "Todas", "filter_ready_label": "Aprobadas", "filter_processing_label": "Pendientes", "empty_icon": "user-check", "empty_state_title": "No hay solicitudes", "empty_state_description": "No se encontraron solicitudes de vinculación", "item_actions": [{"action_id": "approve-request", "label": "Aprobar", "icon": "check", "style": "primary"}, {"action_id": "reject-request", "label": "Rechazar", "icon": "close", "style": "destructive"}]}'::jsonb,
 'school', 'guardian_relations:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 35: Formulario de Evaluación (Fase 3 - Master-Detail)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000070', 'assessments-form',
 'a0000000-0000-0000-0000-000000000007', 'Formulario de Evaluación', 'Crear o editar una evaluación con panel de preguntas',
 '{"page_title": "Nueva Evaluación", "edit_title": "Editar Evaluación", "submit_label": "Guardar", "delete_label": "Eliminar", "questions_label": "Preguntas", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Nombre de la evaluación", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la evaluación"}, {"key": "material_ids", "type": "multi-chip-select", "label": "Materiales", "options_endpoint": "mobile:/api/v1/materials", "option_label": "title", "option_value": "id"}, {"key": "pass_threshold", "type": "number", "label": "Umbral aprobación (%)", "default": 70, "min": 0, "max": 100}, {"key": "max_attempts", "type": "number", "label": "Intentos máximos", "placeholder": "Sin límite"}, {"key": "is_timed", "type": "toggle", "label": "Con cronómetro", "default": false}, {"key": "time_limit_minutes", "type": "number", "label": "Tiempo límite (min)", "placeholder": "Sin límite"}, {"key": "shuffle_questions", "type": "toggle", "label": "Aleatorizar preguntas", "default": false}, {"key": "show_correct_answers", "type": "toggle", "label": "Mostrar respuestas correctas", "default": true}, {"key": "available_from", "type": "date", "label": "Disponible desde"}, {"key": "available_until", "type": "date", "label": "Disponible hasta"}, {"key": "status", "type": "select", "label": "Estado", "required": true, "options": [{"value": "draft", "label": "Borrador"}, {"value": "published", "label": "Publicado"}]}], "detailConfig": {"screenKey": "assessment-questions-list", "modalScreenKey": "assessment-question-form", "parentIdParam": "assessmentId", "childIdField": "question_id", "title": "Preguntas", "masterWeight": 0.4, "detailWeight": 0.6}}'::jsonb,
 'unit', 'assessments:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 36: Tomar Evaluación (Fase 3 - Assessment Take)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000071', 'assessment-take',
 'a0000000-0000-0000-0000-000000000004', 'Tomar Evaluación', 'Pantalla para rendir una evaluación',
 '{"page_title": "Evaluación"}'::jsonb,
 'unit', 'assessments:read', 'assessment-take')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 41: Lista de Gestión de Evaluaciones (pantalla para profesores/admins)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000091', 'assessments-management-list',
 'a0000000-0000-0000-0000-000000000003', 'Gestión de Evaluaciones', 'Lista de evaluaciones para gestión (crear, editar, publicar)',
 '{"page_title": "Gestión de Evaluaciones", "search_placeholder": "Buscar evaluación...", "filter_all_label": "Todas", "filter_ready_label": "Publicadas", "filter_processing_label": "Borradores", "empty_icon": "clipboard", "empty_state_title": "No hay evaluaciones", "empty_state_description": "Crea la primera evaluación para tu unidad", "empty_action_label": "Crear Evaluación"}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 37: Lista de Eventos de Auditoría (Fase 6 - Audit)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000072', 'audit-events-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Eventos de Auditoría', 'Lista de eventos de auditoría del sistema',
 '{"page_title": "Auditoría", "search_placeholder": "Buscar evento...", "filter_all_label": "Todos", "filter_ready_label": "Normales", "filter_processing_label": "Críticos", "empty_icon": "file-search", "empty_state_title": "No hay eventos de auditoría", "empty_state_description": "No se encontraron eventos de auditoría"}'::jsonb,
 'system', 'audit:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 38: Lista de Preguntas de Evaluación (Assessment CRUD)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000075', 'assessment-questions-list',
 'a0000000-0000-0000-0000-000000000003', 'Preguntas de Evaluación', 'Lista de preguntas de una evaluación',
 '{"page_title": "Preguntas", "search_placeholder": "Buscar pregunta...", "filter_all_label": "Todas", "filter_ready_label": "Fácil", "filter_processing_label": "Difícil", "empty_icon": "help_outline", "empty_state_title": "Sin preguntas", "empty_state_description": "Agrega la primera pregunta a esta evaluación", "empty_action_label": "Agregar Pregunta"}'::jsonb,
 'unit', 'assessments:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 39: Formulario de Pregunta de Evaluación (Assessment CRUD)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000076', 'assessment-question-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Pregunta', 'Crear o editar una pregunta de evaluación',
 '{"page_title": "Nueva Pregunta", "edit_title": "Editar Pregunta", "submit_label": "Guardar Pregunta", "delete_label": "Eliminar", "fields": [{"key": "question_text", "type": "textarea", "label": "Texto de la pregunta", "placeholder": "Escribe la pregunta...", "required": true}, {"key": "question_type", "type": "select", "label": "Tipo de pregunta", "required": true, "options": [{"value": "multiple_choice", "label": "Opción múltiple"}, {"value": "true_false", "label": "Verdadero/Falso"}, {"value": "open", "label": "Respuesta abierta"}]}, {"key": "options", "type": "option-list", "label": "Opciones de respuesta", "required": true}, {"key": "difficulty", "type": "select", "label": "Dificultad", "required": true, "options": [{"value": "easy", "label": "Fácil"}, {"value": "medium", "label": "Media"}, {"value": "hard", "label": "Difícil"}]}, {"key": "points", "type": "number", "label": "Puntos", "placeholder": "10", "required": true}, {"key": "explanation", "type": "textarea", "label": "Explicación (feedback)", "placeholder": "Por qué esta es la respuesta correcta..."}]}'::jsonb,
 'unit', 'assessments:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 40: Resultado de Evaluación (Fase 3 - Assessment Result)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000090', 'assessment-result',
 'a0000000-0000-0000-0000-000000000004', 'Resultado de Evaluacion', 'Pantalla de resultados de evaluación',
 '{"page_title": "Resultado de Evaluación", "sections": [{"key": "score_hero", "title": "Puntuación"}, {"key": "answers_detail", "title": "Detalle de Respuestas"}]}'::jsonb,
 'unit', 'assessments:view_results', 'assessment-result')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 42: Lista de Tipos de Concepto
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000092', 'concept-types-list',
 'a0000000-0000-0000-0000-000000000003', 'Tipos de Concepto', 'Lista de tipos de institución y terminología',
 '{"page_title": "Tipos de Concepto", "search_placeholder": "Buscar tipo...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "tag", "empty_state_title": "No hay tipos de concepto", "empty_state_description": "No se encontraron tipos de concepto", "empty_action_label": "Crear Tipo"}'::jsonb,
 'system', 'concept_types:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 43: Formulario de Tipo de Concepto
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000093', 'concept-types-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Tipo de Concepto', 'Crear o editar un tipo de concepto',
 '{"page_title": "Nuevo Tipo de Concepto", "edit_title": "Editar Tipo de Concepto", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre del tipo", "required": true}, {"key": "code", "type": "text", "label": "Código", "placeholder": "Código único (ej: school)", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del tipo de concepto", "required": false}, {"key": "is_active", "type": "toggle", "label": "Activo", "default": true}]}'::jsonb,
 'system', 'concept_types:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 44: Formulario de Screen Instance (screens-form) — admin crear instancias via UI
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000094', 'screens-form',
 'a0000000-0000-0000-0000-000000000006', 'Nueva Screen Instance', 'Formulario para crear una nueva screen instance',
 '{"page_title": "Nueva Screen Instance", "edit_title": "Editar Screen Instance", "submit_label": "Guardar", "delete_label": "Eliminar", "submit_endpoint": "admin:/api/v1/screen-instances", "fields": [{"key": "screen_key", "type": "text", "label": "Screen Key", "placeholder": "ej: my-screen-list", "required": true}, {"key": "template_id", "type": "remote_select", "label": "Template", "required": true, "remote_endpoint": "admin:/api/v1/screen-config/templates", "display_field": "name", "value_field": "id"}, {"key": "name", "type": "text", "label": "Nombre", "placeholder": "Nombre de la instancia", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción de la pantalla"}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "system", "label": "Sistema"}, {"value": "school", "label": "Escuela"}, {"value": "unit", "label": "Unidad"}]}, {"key": "required_permission", "type": "text", "label": "Permiso requerido", "placeholder": "ej: screen_config:read"}, {"key": "handler_key", "type": "text", "label": "Handler Key", "placeholder": "Clave del handler (opcional)"}, {"key": "is_active", "type": "toggle", "label": "Activa", "default": true}]}'::jsonb,
 'system', 'screen_config:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 45: Detalle de Evento de Auditoría (audit-detail)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000095', 'audit-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Auditoría', 'Detalle de un evento de auditoría',
 '{"page_title": "Detalle de Auditoría", "data_endpoint": "admin:/api/v1/audit/{id}", "file_size_label": "Recurso", "uploaded_label": "Fecha", "status_label": "Acción", "description_title": "Información del Evento", "summary_title": "Metadata", "download_label": "Exportar", "quiz_label": "Ver contexto", "sections": [{"key": "action", "label": "Acción", "field": "action"}, {"key": "resource_type", "label": "Tipo de Recurso", "field": "resource_type"}, {"key": "actor_email", "label": "Usuario", "field": "actor_email"}, {"key": "created_at", "label": "Fecha y Hora", "field": "created_at", "type": "datetime"}, {"key": "metadata", "label": "Metadata", "field": "metadata", "type": "json"}]}'::jsonb,
 'system', 'audit:read', 'audit-detail')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 46: Formulario de Material (materials-form)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000096', 'materials-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Material', 'Crear un nuevo material educativo',
 '{"page_title": "Nuevo Material", "edit_title": "Editar Material", "submit_label": "Guardar", "delete_label": "Eliminar", "submit_endpoint": "mobile:/api/v1/materials", "fields": [{"key": "title", "type": "text", "label": "Título", "placeholder": "Título del material", "required": true}, {"key": "description", "type": "textarea", "label": "Descripción", "placeholder": "Descripción del material", "required": false}, {"key": "subject", "type": "text", "label": "Materia", "placeholder": "Materia relacionada", "required": false}, {"key": "status", "type": "select", "label": "Estado", "required": true, "options": [{"value": "uploaded", "label": "Subido"}, {"value": "processing", "label": "Procesando"}, {"value": "ready", "label": "Listo"}, {"value": "failed", "label": "Fallido"}]}, {"key": "file_url", "type": "text", "label": "URL del Archivo", "placeholder": "https://...", "required": true}]}'::jsonb,
 'school', 'materials:create', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 47: Detalle de Progreso
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000097', 'progress-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Progreso', 'Detalle de progreso académico por escuela',
 '{"page_title": "Detalle de Progreso", "data_endpoint": "mobile:/api/v1/stats/global?school_id={schoolId}", "sections": [{"key": "total_materials", "label": "Total Materiales", "field": "total_materials"}, {"key": "completed_progress", "label": "Progreso Completado", "field": "completed_progress"}, {"key": "average_attempt_score", "label": "Promedio de Evaluaciones", "field": "average_attempt_score"}]}'::jsonb,
 'school', 'progress:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia 48: Detalle de Estadísticas
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000098', 'stats-detail',
 'a0000000-0000-0000-0000-000000000004', 'Detalle de Estadísticas', 'Estadísticas detalladas del sistema',
 '{"page_title": "Estadísticas del Sistema", "data_endpoint": "admin:/api/v1/stats/global", "sections": [{"key": "total_users", "label": "Usuarios Totales", "field": "total_users"}, {"key": "total_active_users", "label": "Usuarios Activos", "field": "total_active_users"}, {"key": "total_schools", "label": "Escuelas", "field": "total_schools"}, {"key": "total_subjects", "label": "Materias", "field": "total_subjects"}, {"key": "total_guardian_relations", "label": "Relaciones de Tutor", "field": "total_guardian_relations"}]}'::jsonb,
 'system', 'stats:read', NULL)
ON CONFLICT (screen_key) DO NOTHING;

COMMIT;
