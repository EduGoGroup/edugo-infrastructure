-- Tabla: academic_units (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin, api-mobile

CREATE TABLE IF NOT EXISTS academic_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('grade', 'class', 'section')),
    level VARCHAR(50),
    academic_year INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(school_id, code, academic_year)
);

CREATE INDEX idx_academic_units_school ON academic_units(school_id);
CREATE INDEX idx_academic_units_type ON academic_units(type);
CREATE INDEX idx_academic_units_year ON academic_units(academic_year);
CREATE INDEX idx_academic_units_active ON academic_units(is_active);

COMMENT ON TABLE academic_units IS 'Unidades académicas (cursos, clases, secciones)';
COMMENT ON COLUMN academic_units.type IS 'Tipo: grade (grado), class (clase), section (sección)';
