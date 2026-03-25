-- =============================================================================
-- EduGo Development Seeds v2 — 003_users.sql
-- =============================================================================
-- 21 usuarios de prueba con contrasena unificada: "12345678"
-- Hash bcrypt (cost=10): $2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau
--
-- U-01: super@edugo.test             — Santiago Ramirez   — Super Admin
-- U-02: admin.sanignacio@...         — Carmen Valdes      — School Admin (San Ignacio)
-- U-03: admin.crearte@...            — Roberto Silva      — School Admin (CreArte)
-- U-04: coord.academico@...          — Lucia Fernandez    — Coordinator (San Ignacio + CreArte)
-- U-05: prof.martinez@...            — Maria Martinez     — Teacher (San Ignacio + Academia)
-- U-06: prof.gonzalez@...            — Pedro Gonzalez     — Teacher (San Ignacio)
-- U-07: facilitador.ruiz@...         — Ana Ruiz           — Teacher (CreArte)
-- U-08: est.carlos@...               — Carlos Mendoza     — Student (San Ignacio + CreArte)
-- U-09: est.sofia@...                — Sofia Herrera      — Student (San Ignacio)
-- U-10: est.diego@...                — Diego Vargas       — Student (San Ignacio)
-- U-11: est.valentina@...            — Valentina Rojas    — Student (San Ignacio + Academia)
-- U-12: est.mateo@...                — Mateo Fuentes      — Student (CreArte)
-- U-13: tutor.mendoza@...            — Ricardo Mendoza    — Guardian (San Ignacio + CreArte)
-- U-14: tutora.herrera@...           — Patricia Herrera   — Guardian (San Ignacio)
-- U-15: admin.plataforma@...         — Elena Torres       — Platform Admin
-- U-16: director.sanignacio@...      — Miguel Castillo    — School Director (San Ignacio)
-- U-17: asist.admin@...              — Laura Pena         — School Assistant (San Ignacio)
-- U-18: asist.prof@...               — Andres Gomez       — Assistant Teacher (San Ignacio/5to A)
-- U-19: observador@...               — Diana Lopez        — Observer (San Ignacio + CreArte)
-- U-20: guardian.pendiente@...       — Fernando Ruiz      — Guardian with pending request
-- U-21: readonly@...                 — Test ReadOnly       — Readonly Tester (San Ignacio)
-- =============================================================================

BEGIN;

INSERT INTO auth.users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active
) VALUES

-- Sistema
(
    '00000000-0000-0000-0000-000000000001',
    'super@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Santiago',
    'Ramirez',
    true
),

-- Administradores
(
    '00000000-0000-0000-0000-000000000002',
    'admin.sanignacio@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Carmen',
    'Valdes',
    true
),
(
    '00000000-0000-0000-0000-000000000003',
    'admin.crearte@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Roberto',
    'Silva',
    true
),

-- Coordinadora multi-escuela
(
    '00000000-0000-0000-0000-000000000004',
    'coord.academico@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Lucia',
    'Fernandez',
    true
),

-- Docentes
(
    '00000000-0000-0000-0000-000000000005',
    'prof.martinez@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Maria',
    'Martinez',
    true
),
(
    '00000000-0000-0000-0000-000000000006',
    'prof.gonzalez@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Pedro',
    'Gonzalez',
    true
),
(
    '00000000-0000-0000-0000-000000000007',
    'facilitador.ruiz@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Ana',
    'Ruiz',
    true
),

-- Estudiantes
(
    '00000000-0000-0000-0000-000000000008',
    'est.carlos@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Carlos',
    'Mendoza',
    true
),
(
    '00000000-0000-0000-0000-000000000009',
    'est.sofia@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Sofia',
    'Herrera',
    true
),
(
    '00000000-0000-0000-0000-000000000010',
    'est.diego@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Diego',
    'Vargas',
    true
),
(
    '00000000-0000-0000-0000-000000000011',
    'est.valentina@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Valentina',
    'Rojas',
    true
),
(
    '00000000-0000-0000-0000-000000000012',
    'est.mateo@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Mateo',
    'Fuentes',
    true
),

-- Tutores
(
    '00000000-0000-0000-0000-000000000013',
    'tutor.mendoza@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Ricardo',
    'Mendoza',
    true
),
(
    '00000000-0000-0000-0000-000000000014',
    'tutora.herrera@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Patricia',
    'Herrera',
    true
),

-- Nuevos roles (U-15 a U-20)
(
    '00000000-0000-0000-0000-000000000015',
    'admin.plataforma@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Elena',
    'Torres',
    true
),
(
    '00000000-0000-0000-0000-000000000016',
    'director.sanignacio@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Miguel',
    'Castillo',
    true
),
(
    '00000000-0000-0000-0000-000000000017',
    'asist.admin@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Laura',
    'Pena',
    true
),
(
    '00000000-0000-0000-0000-000000000018',
    'asist.prof@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Andres',
    'Gomez',
    true
),
(
    '00000000-0000-0000-0000-000000000019',
    'observador@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Diana',
    'Lopez',
    true
),
(
    '00000000-0000-0000-0000-000000000020',
    'guardian.pendiente@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Fernando',
    'Ruiz',
    true
),

-- Readonly Tester
(
    '00000000-0000-0000-0000-000000000021',
    'readonly@edugo.test',
    '$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau',
    'Test',
    'ReadOnly',
    true
)

ON CONFLICT (id) DO UPDATE SET
    email         = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    first_name    = EXCLUDED.first_name,
    last_name     = EXCLUDED.last_name,
    is_active     = EXCLUDED.is_active,
    updated_at    = now();

COMMIT;
