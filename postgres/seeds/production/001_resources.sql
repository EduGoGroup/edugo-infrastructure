-- ============================================================
-- SEED: resources (28 registros)
-- Fecha: 2026-02-24
-- Fuente: Datos reales de producción (Neon)
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ============================================================

INSERT INTO iam.resources (id, key, display_name, description, icon, scope, parent_id, sort_order, is_menu_visible, is_active)
VALUES
  -- Recursos raíz (sin parent)
  ('20000000-0000-0000-0000-000000000001', 'dashboard',        'Dashboard',              'Panel principal',                               'dashboard',             'system', NULL,                                   1, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000002', 'admin',            'Administración',         'Módulo de administración',                      'settings',              'system', NULL,                                   2, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000003', 'academic',         'Académico',              'Módulo académico',                              'graduation-cap',        'school', NULL,                                   3, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000004', 'content',          'Contenido',              'Contenido educativo',                           'book-open',             'unit',   NULL,                                   4, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000005', 'reports',          'Reportes',               'Reportes y estadísticas',                       'bar-chart',             'school', NULL,                                   5, TRUE,  TRUE),

  -- Hijos de admin (20000000-0000-0000-0000-000000000002)
  ('20000000-0000-0000-0000-000000000010', 'users',            'Usuarios',               'Gestión de usuarios',                           'users',                 'school', '20000000-0000-0000-0000-000000000002',  1, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000011', 'schools',          'Escuelas',               'Gestión de escuelas',                           'school',                'system', '20000000-0000-0000-0000-000000000002',  2, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000012', 'roles',            'Roles',                  'Gestión de roles',                              'shield',                'system', '20000000-0000-0000-0000-000000000002',  3, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'Permisos',               'Gestión de permisos',                           'key',                   'system', '20000000-0000-0000-0000-000000000002',  4, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000050', 'screen_templates', 'Templates de Pantalla',  'Templates base para configuracion de pantallas', 'settings_applications', 'system', '20000000-0000-0000-0000-000000000002',  5, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000051', 'screen_instances', 'Instancias de Pantalla', 'Instancias configuradas de pantalla por escuela','devices',               'system', '20000000-0000-0000-0000-000000000002',  6, TRUE,  TRUE),

  -- Hijos de academic (20000000-0000-0000-0000-000000000003)
  ('20000000-0000-0000-0000-000000000020', 'units',                'Unidades Académicas',    'Gestión de clases',                             'layers',                'school', '20000000-0000-0000-0000-000000000003',  1, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000021', 'memberships',          'Miembros',               'Asignación de miembros',                        'user-plus',             'school', '20000000-0000-0000-0000-000000000003',  2, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000032', 'subjects',             'Materias',               'Gestión de materias',                           'book',                  'school', '20000000-0000-0000-0000-000000000003',  3, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000060', 'guardian_relations',   'Vínculos Guardian',      'Gestión de vínculos guardian-estudiante',        'user-check',            'school', '20000000-0000-0000-0000-000000000003',  4, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000034', 'periods',              'Periodos Académicos',    'Gestión de periodos académicos',                'calendar-range',        'school', '20000000-0000-0000-0000-000000000003',  5, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000035', 'grades',               'Calificaciones',         'Gestión de calificaciones',                     'award',                 'unit',   '20000000-0000-0000-0000-000000000003',  6, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000036', 'attendance',           'Asistencia',             'Registro de asistencia',                        'check-square',          'unit',   '20000000-0000-0000-0000-000000000003',  7, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000037', 'schedules',            'Horarios',               'Horarios semanales por unidad',                 'clock',                 'unit',   '20000000-0000-0000-0000-000000000003',  8, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000038', 'announcements',        'Anuncios',               'Comunicaciones escuela-familia',                'megaphone',             'school', '20000000-0000-0000-0000-000000000003',  9, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000039', 'calendar',             'Calendario',             'Calendario escolar',                            'calendar',              'school', '20000000-0000-0000-0000-000000000003', 10, TRUE,  TRUE),

  -- Hijos de content (20000000-0000-0000-0000-000000000004)
  ('20000000-0000-0000-0000-000000000030', 'materials',           'Materiales',             'Materiales educativos',                                        'file-text',   'unit', '20000000-0000-0000-0000-000000000004', 1, TRUE, TRUE),
  ('20000000-0000-0000-0000-000000000031', 'assessments',         'Evaluaciones',           'Evaluaciones y exámenes',                                      'clipboard',   'unit', '20000000-0000-0000-0000-000000000004', 2, TRUE, TRUE),
  ('20000000-0000-0000-0000-000000000033', 'assessments_student', 'Tomar Evaluación',       'Vista de evaluaciones desde perspectiva del estudiante',        'play-circle', 'unit', '20000000-0000-0000-0000-000000000004', 3, TRUE, TRUE),

  -- Hijos de reports (20000000-0000-0000-0000-000000000005)
  ('20000000-0000-0000-0000-000000000040', 'progress',         'Progreso',               'Seguimiento de progreso',                       'trending-up',           'unit',   '20000000-0000-0000-0000-000000000005',  1, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000041', 'stats',            'Estadísticas',           'Estadísticas del sistema',                      'pie-chart',             'school', '20000000-0000-0000-0000-000000000005',  2, TRUE,  TRUE),

  -- Recurso mobile (sin parent)
  ('20000000-0000-0000-0000-000000000052', 'screens',          'Pantallas (Mobile)',      'Lectura de pantallas desde aplicacion mobile',  'smartphone',            'system', NULL,                                   0, FALSE, TRUE),

  -- Hijo de admin (20000000-0000-0000-0000-000000000002)
  ('20000000-0000-0000-0000-000000000070', 'audit',            'Auditoría',               'Registro de auditoría del sistema',             'file-search',           'system', '20000000-0000-0000-0000-000000000002',  7, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000080', 'concept_types',    'Tipos de Concepto',       'Tipos de institución y terminología',           'tag',                   'system', '20000000-0000-0000-0000-000000000002',  8, TRUE,  TRUE),
  ('20000000-0000-0000-0000-000000000090', 'system_settings',  'Configuración',           'Configuración y mantenimiento del sistema',     'settings',              'system', '20000000-0000-0000-0000-000000000002',  9, TRUE,  TRUE),

  -- Recurso de contexto (no visible en menu, solo para permisos de browsing)
  ('20000000-0000-0000-0000-0000000000A0', 'context',          'Contexto',                'Exploración de escuelas y unidades para selección de contexto', 'swap_horiz', 'system', NULL, 99, FALSE, TRUE)
ON CONFLICT (id) DO NOTHING;
