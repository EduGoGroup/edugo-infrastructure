-- ====================================================================
-- SEEDS: Resources del sistema (catalogo de modulos para RBAC y menu)
-- VERSION: postgres/v0.17.0
-- ====================================================================

-- Recursos raiz (nivel 1)
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000001', 'dashboard', 'Dashboard', 'Panel principal', 'dashboard', NULL, 1, true, 'system'),
('20000000-0000-0000-0000-000000000002', 'admin', 'Administracion', 'Modulo de administracion', 'settings', NULL, 2, true, 'system'),
('20000000-0000-0000-0000-000000000003', 'academic', 'Academico', 'Modulo academico', 'graduation-cap', NULL, 3, true, 'school'),
('20000000-0000-0000-0000-000000000004', 'content', 'Contenido', 'Contenido educativo', 'book-open', NULL, 4, true, 'unit'),
('20000000-0000-0000-0000-000000000005', 'reports', 'Reportes', 'Reportes y estadisticas', 'bar-chart', NULL, 5, true, 'school');

-- Hijos de Administracion
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000010', 'users', 'Usuarios', 'Gestion de usuarios', 'users', '20000000-0000-0000-0000-000000000002', 1, true, 'school'),
('20000000-0000-0000-0000-000000000011', 'schools', 'Escuelas', 'Gestion de escuelas', 'school', '20000000-0000-0000-0000-000000000002', 2, true, 'system'),
('20000000-0000-0000-0000-000000000012', 'roles', 'Roles', 'Gestion de roles', 'shield', '20000000-0000-0000-0000-000000000002', 3, true, 'system'),
('20000000-0000-0000-0000-000000000013', 'permissions_mgmt', 'Permisos', 'Gestion de permisos', 'key', '20000000-0000-0000-0000-000000000002', 4, true, 'system');

-- Hijos de Academico
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000020', 'units', 'Unidades Academicas', 'Gestion de clases', 'layers', '20000000-0000-0000-0000-000000000003', 1, true, 'school'),
('20000000-0000-0000-0000-000000000021', 'memberships', 'Miembros', 'Asignacion de miembros', 'user-plus', '20000000-0000-0000-0000-000000000003', 2, true, 'school');

-- Hijos de Contenido
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000030', 'materials', 'Materiales', 'Materiales educativos', 'file-text', '20000000-0000-0000-0000-000000000004', 1, true, 'unit'),
('20000000-0000-0000-0000-000000000031', 'assessments', 'Evaluaciones', 'Evaluaciones y examenes', 'clipboard', '20000000-0000-0000-0000-000000000004', 2, true, 'unit');

-- Hijos de Reportes
INSERT INTO resources (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope) VALUES
('20000000-0000-0000-0000-000000000040', 'progress', 'Progreso', 'Seguimiento de progreso', 'trending-up', '20000000-0000-0000-0000-000000000005', 1, true, 'unit'),
('20000000-0000-0000-0000-000000000041', 'stats', 'Estadisticas', 'Estadisticas del sistema', 'pie-chart', '20000000-0000-0000-0000-000000000005', 2, true, 'school');
