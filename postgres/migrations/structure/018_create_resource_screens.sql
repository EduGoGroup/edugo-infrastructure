CREATE TABLE ui_config.resource_screens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID NOT NULL,
    resource_key VARCHAR(100) NOT NULL,
    screen_key VARCHAR(100) NOT NULL,
    screen_type VARCHAR(50) NOT NULL,
    is_default BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(resource_id, screen_type)
);
