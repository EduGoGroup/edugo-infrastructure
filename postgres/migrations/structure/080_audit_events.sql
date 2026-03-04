-- ============================================================
-- 080: Audit Events
-- Tabla principal de auditoría para registrar acciones del sistema
-- ============================================================

CREATE TABLE IF NOT EXISTS audit.events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- WHO
    actor_id        UUID NOT NULL,
    actor_email     VARCHAR(255) NOT NULL,
    actor_role      VARCHAR(100) NOT NULL,
    actor_ip        VARCHAR(45),
    actor_user_agent TEXT,
    -- WHERE (context)
    school_id       UUID,
    unit_id         UUID,
    service_name    VARCHAR(50) NOT NULL,
    -- WHAT
    action          VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(100) NOT NULL,
    resource_id     VARCHAR(255),
    permission_used VARCHAR(100),
    -- REQUEST DETAILS
    request_method  VARCHAR(10),
    request_path    VARCHAR(500),
    request_id      VARCHAR(100),
    status_code     INTEGER,
    changes         JSONB,
    metadata        JSONB,
    error_message   TEXT,
    -- WHEN
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- CLASSIFICATION
    severity        VARCHAR(20) NOT NULL DEFAULT 'info',
    category        VARCHAR(50) NOT NULL DEFAULT 'data',

    -- Constraints
    CONSTRAINT audit_events_severity_check CHECK (severity IN ('info', 'warning', 'critical')),
    CONSTRAINT audit_events_category_check CHECK (category IN ('auth', 'data', 'config', 'admin'))
);

-- Indexes for frequent queries
CREATE INDEX IF NOT EXISTS idx_audit_events_actor
    ON audit.events (actor_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_resource
    ON audit.events (resource_type, resource_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_action
    ON audit.events (action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_school
    ON audit.events (school_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_created
    ON audit.events (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_events_severity
    ON audit.events (severity, created_at DESC)
    WHERE severity != 'info';

CREATE INDEX IF NOT EXISTS idx_audit_events_category
    ON audit.events (category, created_at DESC);

COMMENT ON TABLE audit.events IS 'Audit trail for all system actions in EduGo';
COMMENT ON COLUMN audit.events.actor_id IS 'UUID of the user who performed the action (from JWT)';
COMMENT ON COLUMN audit.events.service_name IS 'Service that recorded the event: iam-platform, admin-api, mobile-api, worker';
COMMENT ON COLUMN audit.events.severity IS 'Event severity: info (normal ops), warning (deletions, failures), critical (RBAC changes)';
COMMENT ON COLUMN audit.events.category IS 'Event category: auth (login/logout), data (CRUD), config (settings), admin (RBAC)';
