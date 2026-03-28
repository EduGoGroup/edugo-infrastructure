-- ============================================================
-- 095: notifications.notifications
-- Schema: notifications
-- Notificaciones para usuarios del sistema
-- Cross-schema FK (user_id -> auth.users) definida al final de este archivo
-- ============================================================

CREATE SCHEMA IF NOT EXISTS notifications;

CREATE TABLE notifications.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    body text,
    resource_type character varying(50),
    resource_id uuid,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    read_at timestamptz,
    CONSTRAINT notifications_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_notif_user_unread ON notifications.notifications(user_id, created_at DESC) WHERE is_read = FALSE;
CREATE INDEX idx_notif_user_all ON notifications.notifications(user_id, created_at DESC);

-- FK cross-schema: notifications.notifications -> auth.users
ALTER TABLE notifications.notifications ADD CONSTRAINT notifications_user_fkey
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;
