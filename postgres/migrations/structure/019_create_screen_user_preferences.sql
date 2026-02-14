CREATE TABLE ui_config.screen_user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    screen_instance_id UUID NOT NULL,
    user_id UUID NOT NULL,
    preferences JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(screen_instance_id, user_id)
);
