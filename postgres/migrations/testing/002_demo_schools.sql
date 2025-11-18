-- Mock Data: Escuelas de demostración para testing

INSERT INTO schools (id, name, code, address, city, country, phone, email, created_at, updated_at) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Escuela Primaria Demo', 'SCH_PRI_001', 'Calle Principal 123', 'Buenos Aires', 'Argentina', '+54-11-1234-5678', 'contacto@primaria.test', NOW(), NOW()),
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Colegio Secundario Demo', 'SCH_SEC_001', 'Avenida Libertador 456', 'Buenos Aires', 'Argentina', '+54-11-8765-4321', 'info@secundario.test', NOW(), NOW()),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Instituto Técnico Demo', 'SCH_TEC_001', 'Boulevard Tecnológico 789', 'Córdoba', 'Argentina', '+54-351-999-8888', 'admin@tecnico.test', NOW(), NOW());
