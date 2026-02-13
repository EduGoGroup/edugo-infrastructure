-- ====================================================================
-- SEEDS: Resources del sistema (catalogo de modulos para RBAC y menu)
-- VERSION: postgres/v0.17.0
-- ====================================================================

-- Recursos raiz (nivel 1)
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'Dashboard', 'Panel principal', 'dashboard', NULL, 1, true, 'system'),
('20000000-0000-0000-0000-000000000002', 'admin', 'Administración', 'Módulo de administración', 'settings', NULL, 2, true, 'system'),
('20000000-0000-0000-0000-000000000003', 'academic', 'Académico', 'Módulo académico', 'graduation-cap', NULL, 3, true, 'school'),
('20000000-0000-0000-0000-000000000004', 'content', 'Contenido', 'Contenido educativo', 'book-open', NULL, 4, true, 'unit'),
('20000000-0000-0000-0000-000000000005', 'reports', 'Reportes', 'Reportes y estadísticas', 'bar-chart', NULL, 5, true, 'school');

-- Hijos de Administración
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'Usuarios', 'Gestión de usuarios', 'users', '20000000-0000-0000-0000-000000000002', 1, true, 'school'),
('20000000-0000-0000-0000-000000000011', 'schools', 'Escuelas', 'Gestión de escuelas', 'school', '20000000-0000-0000-0000-000000000002', 2, true, 'system'),
('20000000-0000-0000-0000-000000000012', 'roles', 'Roles', 'Gestión de roles', 'shield', '20000000-0000-0000-0000-000000000002', 3, true, 'system'),
('20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'Permisos', 'Gestión de permisos', 'key', '20000000-0000-0000-0000-000000000002', 4, true, 'system');

-- Hijos de Académico
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'Unidades Académicas', 'Gestión de clases', 'layers', '20000000-0000-0000-0000-000000000003', 1, true, 'school'),
('20000000-0000-0000-0000-000000000021', 'memberships', 'Miembros', 'Asignación de miembros', 'user-plus', '20000000-0000-0000-0000-000000000003', 2, true, 'school');

-- Hijos de Contenido
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'Materiales', 'Materiales educativos', 'file-text', '20000000-0000-0000-0000-000000000004', 1, true, 'unit'),
('20000000-0000-0000-0000-000000000031', 'assessments', 'Evaluaciones', 'Evaluaciones y exámenes', 'clipboard', '20000000-0000-0000-0000-000000000004', 2, true, 'unit');

-- Hijos de Reportes
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000040', 'progress', 'Progreso', 'Seguimiento de progreso', 'trending-up', '20000000-0000-0000-0000-000000000005', 1, true, 'unit'),
('20000000-0000-0000-0000-000000000041', 'stats', 'Estadísticas', 'Estadísticas del sistema', 'pie-chart', '20000000-0000-0000-0000-000000000005', 2, true, 'school');
