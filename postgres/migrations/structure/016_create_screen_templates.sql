-- Schema para configuración de UI
CREATE SCHEMA IF NOT EXISTS ui_config;

-- Templates de pantalla
CREATE TABLE ui_config.screen_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pattern VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    version INT NOT NULL DEFAULT 1,
    definition JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(name, version)
);
