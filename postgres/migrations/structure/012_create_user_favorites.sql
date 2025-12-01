-- Tabla para almacenar materiales marcados como favoritos por usuarios
-- Permite acceso rápido a contenido frecuentemente usado
-- Parte de FASE 1 UI Roadmap - Funcionalidad de favoritos en apps

CREATE TABLE IF NOT EXISTS user_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    material_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT uq_user_favorites_user_material UNIQUE(user_id, material_id),
    CONSTRAINT fk_user_favorites_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_favorites_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE CASCADE
);

-- Índices para queries frecuentes
CREATE INDEX idx_user_favorites_user ON user_favorites(user_id);
CREATE INDEX idx_user_favorites_material ON user_favorites(material_id);
CREATE INDEX idx_user_favorites_created ON user_favorites(created_at DESC);

-- Comentarios
COMMENT ON TABLE user_favorites IS 'Materiales marcados como favoritos por usuarios';
COMMENT ON COLUMN user_favorites.user_id IS 'Usuario que marcó el favorito';
COMMENT ON COLUMN user_favorites.material_id IS 'Material marcado como favorito';
COMMENT ON COLUMN user_favorites.created_at IS 'Fecha cuando fue agregado a favoritos (para ordenar por reciente)';
