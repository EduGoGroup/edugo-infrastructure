-- Mock Data: Escuelas de demostración para testing

INSERT INTO schools (id, name, address, city, state, country, postal_code, phone, email, created_at, updated_at) VALUES
('sch_demo_primary', 'Escuela Primaria Demo', 'Calle Principal 123', 'Buenos Aires', 'CABA', 'Argentina', '1000', '+54-11-1234-5678', 'contacto@primaria.test', NOW(), NOW()),
('sch_demo_secondary', 'Colegio Secundario Demo', 'Avenida Libertador 456', 'Buenos Aires', 'CABA', 'Argentina', '1001', '+54-11-8765-4321', 'info@secundario.test', NOW(), NOW()),
('sch_demo_tech', 'Instituto Técnico Demo', 'Boulevard Tecnológico 789', 'Córdoba', 'Córdoba', 'Argentina', '5000', '+54-351-999-8888', 'admin@tecnico.test', NOW(), NOW());
