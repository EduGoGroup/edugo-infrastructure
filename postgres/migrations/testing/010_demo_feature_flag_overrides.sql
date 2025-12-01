-- Mock Data: Overrides de feature flags para testing
-- Ejemplos de sobrescrituras específicas por usuario

INSERT INTO feature_flag_overrides (feature_flag_id, user_id, enabled, reason, expires_at, created_by) VALUES

-- Admin tiene acceso a dashboard nuevo (experimental) habilitado
(
    (SELECT id FROM feature_flags WHERE key = 'new_dashboard'),
    (SELECT id FROM users WHERE email = 'admin@edugo.test'),
    true,
    'Testing de dashboard nuevo - Admin beta tester',
    NULL,  -- Sin expiración
    (SELECT id FROM users WHERE email = 'admin@edugo.test')
),

-- Teacher Math tiene debug logs habilitado temporalmente para diagnóstico
(
    (SELECT id FROM feature_flags WHERE key = 'debug_logs'),
    (SELECT id FROM users WHERE email = 'teacher.math@edugo.test'),
    true,
    'Diagnóstico de problema de sincronización reportado',
    NOW() + INTERVAL '7 days',  -- Expira en 7 días
    (SELECT id FROM users WHERE email = 'admin@edugo.test')
),

-- Student1 tiene push notifications habilitado (beta testing)
(
    (SELECT id FROM feature_flags WHERE key = 'push_notifications'),
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    true,
    'Beta tester de notificaciones push',
    NOW() + INTERVAL '30 days',  -- Expira en 30 días
    (SELECT id FROM users WHERE email = 'admin@edugo.test')
),

-- Student2 tiene background sync deshabilitado (problemas de batería)
(
    (SELECT id FROM feature_flags WHERE key = 'background_sync'),
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    false,
    'Usuario reportó consumo excesivo de batería',
    NOW() + INTERVAL '14 days',  -- Expira en 14 días
    (SELECT id FROM users WHERE email = 'admin@edugo.test')
);
