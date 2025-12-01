-- Mock Data: Feature Flags iniciales para testing
-- Basado en BACKEND-SPEC-FEATURE-FLAGS.md
-- 11 feature flags organizados por categoría

INSERT INTO feature_flags (
    key,
    name,
    description,
    enabled,
    enabled_globally,
    category,
    priority,
    requires_restart,
    is_debug_only,
    affects_security,
    is_experimental
) VALUES

-- Security Features (Prioridad 98-100)
(
    'biometric_login',
    'Login Biométrico',
    'Habilita Face ID/Touch ID para autenticación',
    true,
    true,
    'security',
    100,
    false,
    false,
    true,
    false
),
(
    'certificate_pinning',
    'Certificate Pinning',
    'Habilita certificate pinning SSL para mayor seguridad',
    true,
    true,
    'security',
    99,
    true,
    false,
    true,
    false
),
(
    'login_rate_limiting',
    'Rate Limiting Login',
    'Límite de intentos de login para prevenir ataques de fuerza bruta',
    true,
    true,
    'security',
    98,
    false,
    false,
    true,
    false
),

-- Core Features (Prioridad 48-50)
(
    'offline_mode',
    'Modo Offline',
    'Habilita funcionalidad offline con sincronización posterior',
    true,
    true,
    'features',
    50,
    false,
    false,
    false,
    false
),
(
    'background_sync',
    'Sync Background',
    'Sincronización de datos en segundo plano',
    false,
    false,
    'features',
    49,
    false,
    false,
    false,
    false
),
(
    'push_notifications',
    'Notificaciones Push',
    'Habilita notificaciones push para eventos importantes',
    false,
    false,
    'features',
    48,
    false,
    false,
    false,
    false
),

-- UI Features (Prioridad 10-30)
(
    'auto_dark_mode',
    'Tema Oscuro Automático',
    'Tema oscuro según configuración del sistema',
    true,
    true,
    'ui',
    30,
    false,
    false,
    false,
    false
),
(
    'new_dashboard',
    'Dashboard Nuevo',
    'Dashboard rediseñado con nuevas métricas (experimental)',
    false,
    false,
    'ui',
    10,
    false,
    false,
    false,
    true
),
(
    'transition_animations',
    'Animaciones de Transición',
    'Animaciones suaves entre pantallas',
    true,
    true,
    'ui',
    20,
    false,
    false,
    false,
    false
),

-- Debug Features (Prioridad 1-5) - Solo desarrollo
(
    'debug_logs',
    'Logs de Debug',
    'Logs de debug detallados en producción (solo para diagnóstico)',
    false,
    false,
    'debug',
    5,
    true,
    true,
    false,
    false
),
(
    'mock_api',
    'API Mock',
    'Usar API mock en lugar de backend real (solo desarrollo)',
    false,
    false,
    'debug',
    1,
    true,
    true,
    false,
    false
);
