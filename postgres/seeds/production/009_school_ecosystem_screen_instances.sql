-- ====================================================================
-- SEEDS: Screen instances for School Ecosystem features
-- Periodos, Calificaciones, Asistencia, Horarios, Anuncios,
-- Calendario, Directorio y Boleta de Notas
-- Idempotente: usa ON CONFLICT (DO NOTHING o DO UPDATE SET)
-- ====================================================================

BEGIN;

-- ---------------------------------------------------------------
-- Academic Management - Periodos
-- ---------------------------------------------------------------

-- Instancia: Lista de Periodos Academicos
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000110', 'periods-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Periodos Academicos', 'Lista de periodos academicos de la escuela',
 '{"page_title": "Periodos Academicos", "search_placeholder": "Buscar periodo...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "calendar", "empty_state_title": "No hay periodos academicos", "empty_state_description": "No se encontraron periodos academicos configurados", "empty_action_label": "Crear Periodo"}'::jsonb,
 'school', 'periods:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Periodo Academico
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000111', 'periods-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Periodo Academico', 'Crear o editar un periodo academico',
 '{"page_title": "Nuevo Periodo", "edit_title": "Editar Periodo", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "name", "type": "text", "label": "Nombre del periodo", "placeholder": "Ej: Primer Semestre 2026", "required": true}, {"key": "code", "type": "text", "label": "Codigo", "placeholder": "Ej: 2026-S1"}, {"key": "type", "type": "select", "label": "Tipo", "required": true, "options": [{"value": "semester", "label": "Semestre"}, {"value": "trimester", "label": "Trimestre"}, {"value": "bimester", "label": "Bimestre"}, {"value": "quarter", "label": "Cuatrimestre"}]}, {"key": "start_date", "type": "date", "label": "Fecha inicio", "required": true}, {"key": "end_date", "type": "date", "label": "Fecha fin", "required": true}, {"key": "academic_year", "type": "number", "label": "Ano academico", "placeholder": "2026", "required": true}, {"key": "sort_order", "type": "number", "label": "Orden", "placeholder": "1"}]}'::jsonb,
 'school', 'periods:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Academic Management - Calificaciones
-- ---------------------------------------------------------------

-- Instancia: Lista de Calificaciones
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000112', 'grades-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Calificaciones', 'Lista de calificaciones de los estudiantes',
 '{"page_title": "Calificaciones", "search_placeholder": "Buscar estudiante...", "filter_all_label": "Todas", "filter_ready_label": "Aprobadas", "filter_processing_label": "Reprobadas", "empty_icon": "grade", "empty_state_title": "No hay calificaciones", "empty_state_description": "No se encontraron calificaciones registradas", "empty_action_label": "Registrar Calificacion"}'::jsonb,
 'unit', 'grades:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Calificacion
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000113', 'grades-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Calificacion', 'Registrar o editar una calificacion',
 '{"page_title": "Nueva Calificacion", "edit_title": "Editar Calificacion", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "membership_id", "type": "remote_select", "label": "Estudiante", "required": true, "remote_endpoint": "admin:/api/v1/memberships?unit_id={unitId}", "display_field": "user_name", "value_field": "id"}, {"key": "subject_id", "type": "remote_select", "label": "Materia", "required": true, "remote_endpoint": "admin:/api/v1/subjects", "display_field": "name", "value_field": "id"}, {"key": "period_id", "type": "remote_select", "label": "Periodo", "required": true, "remote_endpoint": "admin:/api/v1/periods", "display_field": "name", "value_field": "id"}, {"key": "grade_value", "type": "number", "label": "Nota", "placeholder": "0-100", "required": true}, {"key": "grade_letter", "type": "text", "label": "Letra", "placeholder": "Ej: A, B, C"}, {"key": "notes", "type": "textarea", "label": "Observaciones", "placeholder": "Observaciones sobre el desempeno"}]}'::jsonb,
 'unit', 'grades:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Academic Management - Asistencia
-- ---------------------------------------------------------------

-- Instancia: Lista de Asistencia
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000114', 'attendance-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Asistencia', 'Registro de asistencia de la unidad',
 '{"page_title": "Asistencia", "search_placeholder": "Buscar estudiante...", "filter_all_label": "Todos", "filter_ready_label": "Presentes", "filter_processing_label": "Ausentes", "empty_icon": "check_circle", "empty_state_title": "No hay registros de asistencia", "empty_state_description": "No se encontraron registros de asistencia", "empty_action_label": "Registrar Asistencia"}'::jsonb,
 'unit', 'attendance:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Registro de Asistencia por Lote
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000115', 'attendance-batch',
 'a0000000-0000-0000-0000-000000000006', 'Registrar Asistencia', 'Formulario para registro masivo de asistencia',
 '{"page_title": "Registrar Asistencia", "submit_label": "Guardar Asistencia", "delete_label": "Cancelar", "fields": [{"key": "date", "type": "date", "label": "Fecha", "required": true}, {"key": "records", "type": "attendance-grid", "label": "Estudiantes", "required": true, "status_options": [{"value": "present", "label": "Presente"}, {"value": "absent", "label": "Ausente"}, {"value": "late", "label": "Tarde"}, {"value": "excused", "label": "Justificado"}]}]}'::jsonb,
 'unit', 'attendance:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Asistencia (alias para attendance-batch, usado por dynamic-ui)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000125', 'attendance-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Asistencia', 'Formulario para registro masivo de asistencia',
 '{"page_title": "Registrar Asistencia", "submit_label": "Guardar Asistencia", "delete_label": "Cancelar", "fields": [{"key": "date", "type": "date", "label": "Fecha", "required": true}, {"key": "records", "type": "attendance-grid", "label": "Estudiantes", "required": true, "status_options": [{"value": "present", "label": "Presente"}, {"value": "absent", "label": "Ausente"}, {"value": "late", "label": "Tarde"}, {"value": "excused", "label": "Justificado"}]}]}'::jsonb,
 'unit', 'attendance:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Resumen de Asistencia
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000116', 'attendance-summary',
 'a0000000-0000-0000-0000-000000000003', 'Resumen de Asistencia', 'Resumen estadistico de asistencia por unidad',
 '{"page_title": "Resumen de Asistencia", "search_placeholder": "Buscar estudiante...", "filter_all_label": "Todos", "filter_ready_label": "Regulares", "filter_processing_label": "Irregulares", "empty_icon": "bar_chart", "empty_state_title": "Sin datos de asistencia", "empty_state_description": "No hay datos suficientes para generar el resumen", "show_add_button": false, "readonly": true}'::jsonb,
 'unit', 'attendance:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Academic Management - Horarios
-- ---------------------------------------------------------------

-- Instancia: Lista de Horarios
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000117', 'schedules-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Horarios', 'Horarios de clases de la unidad',
 '{"page_title": "Horarios", "search_placeholder": "Buscar materia o profesor...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "schedule", "empty_state_title": "No hay horarios configurados", "empty_state_description": "No se encontraron bloques horarios para esta unidad", "empty_action_label": "Crear Bloque Horario"}'::jsonb,
 'unit', 'schedules:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Bloque Horario
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000118', 'schedules-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Bloque Horario', 'Crear o editar un bloque horario',
 '{"page_title": "Nuevo Bloque Horario", "edit_title": "Editar Bloque Horario", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "subject_id", "type": "remote_select", "label": "Materia", "required": true, "remote_endpoint": "admin:/api/v1/subjects", "display_field": "name", "value_field": "id"}, {"key": "teacher_membership_id", "type": "remote_select", "label": "Profesor", "required": true, "remote_endpoint": "admin:/api/v1/memberships/by-role?role=teacher&unit_id={unitId}", "display_field": "user_name", "value_field": "id"}, {"key": "day_of_week", "type": "select", "label": "Dia", "required": true, "options": [{"value": "1", "label": "Lunes"}, {"value": "2", "label": "Martes"}, {"value": "3", "label": "Miercoles"}, {"value": "4", "label": "Jueves"}, {"value": "5", "label": "Viernes"}, {"value": "6", "label": "Sabado"}, {"value": "0", "label": "Domingo"}]}, {"key": "start_time", "type": "text", "label": "Hora inicio", "placeholder": "08:00", "required": true}, {"key": "end_time", "type": "text", "label": "Hora fin", "placeholder": "09:30", "required": true}, {"key": "room", "type": "text", "label": "Salon", "placeholder": "Ej: Aula 101"}]}'::jsonb,
 'unit', 'schedules:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Communication - Anuncios
-- ---------------------------------------------------------------

-- Instancia: Lista de Anuncios
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000119', 'announcements-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Anuncios', 'Anuncios de la escuela',
 '{"page_title": "Anuncios", "search_placeholder": "Buscar anuncio...", "filter_all_label": "Todos", "filter_ready_label": "Fijados", "filter_processing_label": "Normales", "empty_icon": "campaign", "empty_state_title": "No hay anuncios", "empty_state_description": "No se encontraron anuncios publicados", "empty_action_label": "Crear Anuncio"}'::jsonb,
 'school', 'announcements:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Anuncio
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000120', 'announcements-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Anuncio', 'Crear o editar un anuncio',
 '{"page_title": "Nuevo Anuncio", "edit_title": "Editar Anuncio", "submit_label": "Publicar", "delete_label": "Eliminar", "fields": [{"key": "title", "type": "text", "label": "Titulo", "placeholder": "Titulo del anuncio", "required": true}, {"key": "body", "type": "textarea", "label": "Contenido", "placeholder": "Escriba el contenido del anuncio...", "required": true}, {"key": "scope", "type": "select", "label": "Alcance", "required": true, "options": [{"value": "school", "label": "Toda la Escuela"}, {"value": "unit", "label": "Unidad Especifica"}]}, {"key": "is_pinned", "type": "toggle", "label": "Fijar anuncio", "default": false}]}'::jsonb,
 'school', 'announcements:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Communication - Calendario Escolar
-- ---------------------------------------------------------------

-- Instancia: Lista de Eventos del Calendario
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000121', 'calendar-list',
 'a0000000-0000-0000-0000-000000000003', 'Calendario Escolar', 'Eventos del calendario escolar',
 '{"page_title": "Calendario Escolar", "search_placeholder": "Buscar evento...", "filter_all_label": "Todos", "filter_ready_label": "Proximos", "filter_processing_label": "Pasados", "empty_icon": "event", "empty_state_title": "No hay eventos en el calendario", "empty_state_description": "No se encontraron eventos programados", "empty_action_label": "Crear Evento"}'::jsonb,
 'school', 'calendar:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Formulario de Evento de Calendario
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000122', 'calendar-form',
 'a0000000-0000-0000-0000-000000000006', 'Formulario de Evento', 'Crear o editar un evento del calendario',
 '{"page_title": "Nuevo Evento", "edit_title": "Editar Evento", "submit_label": "Guardar", "delete_label": "Eliminar", "fields": [{"key": "title", "type": "text", "label": "Titulo", "placeholder": "Titulo del evento", "required": true}, {"key": "description", "type": "textarea", "label": "Descripcion", "placeholder": "Descripcion del evento"}, {"key": "event_type", "type": "select", "label": "Tipo", "required": true, "options": [{"value": "holiday", "label": "Feriado"}, {"value": "exam", "label": "Examen"}, {"value": "meeting", "label": "Reunion"}, {"value": "activity", "label": "Actividad"}, {"value": "deadline", "label": "Fecha Limite"}]}, {"key": "start_date", "type": "date", "label": "Fecha inicio", "required": true}, {"key": "end_date", "type": "date", "label": "Fecha fin"}, {"key": "is_all_day", "type": "toggle", "label": "Todo el dia", "default": false}]}'::jsonb,
 'school', 'calendar:create', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ---------------------------------------------------------------
-- Directory & Reports
-- ---------------------------------------------------------------

-- Instancia: Directorio de Unidad
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000123', 'unit-directory',
 'a0000000-0000-0000-0000-000000000003', 'Directorio', 'Directorio de miembros de la unidad agrupado por rol',
 '{"page_title": "Directorio", "search_placeholder": "Buscar por nombre o email...", "filter_all_label": "Todos", "filter_ready_label": "Activos", "filter_processing_label": "Inactivos", "empty_icon": "contacts", "empty_state_title": "Directorio vacio", "empty_state_description": "No se encontraron miembros en esta unidad", "show_add_button": false, "readonly": true}'::jsonb,
 'unit', 'memberships:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- Instancia: Boleta de Notas
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, scope, required_permission, handler_key) VALUES
('b0000000-0000-0000-0000-000000000124', 'report-card',
 'a0000000-0000-0000-0000-000000000004', 'Boleta de Notas', 'Boleta de notas del estudiante por periodo',
 '{"page_title": "Boleta de Notas", "file_size_label": "Promedio General", "uploaded_label": "Periodo", "status_label": "Estado", "description_title": "Calificaciones por Materia", "summary_title": "Observaciones", "download_label": "Descargar PDF", "quiz_label": "Ver Historial", "sections": [{"key": "student_name", "label": "Estudiante", "field": "student_name"}, {"key": "period_name", "label": "Periodo", "field": "period_name"}, {"key": "overall_average", "label": "Promedio General", "field": "overall_average"}, {"key": "grades", "label": "Calificaciones", "field": "grades", "type": "table"}]}'::jsonb,
 'school', 'reports:read', NULL)
ON CONFLICT (screen_key) DO UPDATE SET
    slot_data = EXCLUDED.slot_data;

-- ================================================================
-- Resource-screen mappings for screens defined in this file
-- (must be after screen_instances due to FK on screen_key)
-- ================================================================

-- Periods
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000110', '20000000-0000-0000-0000-000000000034', 'periods', 'periods-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000111', '20000000-0000-0000-0000-000000000034', 'periods', 'periods-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Grades
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000112', '20000000-0000-0000-0000-000000000035', 'grades', 'grades-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000113', '20000000-0000-0000-0000-000000000035', 'grades', 'grades-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Attendance
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000114', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000115', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-batch', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000116', '20000000-0000-0000-0000-000000000036', 'attendance', 'attendance-summary', 'summary', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Schedules
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000117', '20000000-0000-0000-0000-000000000037', 'schedules', 'schedules-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000118', '20000000-0000-0000-0000-000000000037', 'schedules', 'schedules-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Announcements
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000119', '20000000-0000-0000-0000-000000000038', 'announcements', 'announcements-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000120', '20000000-0000-0000-0000-000000000038', 'announcements', 'announcements-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Calendar
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000121', '20000000-0000-0000-0000-000000000039', 'calendar', 'calendar-list', 'list', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000122', '20000000-0000-0000-0000-000000000039', 'calendar', 'calendar-form', 'form', FALSE, 2)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Memberships -> unit-directory (directory, readonly)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000123', '20000000-0000-0000-0000-000000000021', 'memberships', 'unit-directory', 'directory', FALSE, 3)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

-- Reports -> report-card (detail)
INSERT INTO ui_config.resource_screens (id, resource_id, resource_key, screen_key, screen_type, is_default, sort_order) VALUES
('c0000000-0000-0000-0000-000000000124', '20000000-0000-0000-0000-000000000005', 'reports', 'report-card', 'detail', TRUE, 1)
ON CONFLICT (resource_id, screen_type) DO NOTHING;

COMMIT;
