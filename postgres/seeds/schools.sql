-- Seeds de escuelas de prueba
-- Ejecutar después de migración 002_create_schools

INSERT INTO schools (id, name, code, city, country, subscription_tier, max_teachers, max_students) VALUES
-- Escuela 1
('44444444-4444-4444-4444-444444444444', 'Liceo Técnico Santiago', 'LTS-001', 'Santiago', 'Chile', 'premium', 50, 500),

-- Escuela 2
('55555555-5555-5555-5555-555555555555', 'Colegio Valparaíso', 'CV-002', 'Valparaíso', 'Chile', 'basic', 20, 200)

ON CONFLICT (code) DO NOTHING;
