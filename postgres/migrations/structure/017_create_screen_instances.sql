CREATE TABLE ui_config.screen_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    screen_key VARCHAR(100) NOT NULL UNIQUE,
    template_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    slot_data JSONB NOT NULL DEFAULT '{}',
    actions JSONB NOT NULL DEFAULT '[]',
    data_endpoint VARCHAR(500),
    data_config JSONB DEFAULT '{}',
    scope VARCHAR(20) DEFAULT 'school',
    required_permission VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
