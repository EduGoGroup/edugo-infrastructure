-- ========================================
-- FUNCIONES BASE PARA POSTGRESQL
-- ========================================
-- Este archivo contiene funciones auxiliares utilizadas
-- por múltiples tablas en el sistema EduGo
-- Debe ejecutarse ANTES que cualquier otra migración

-- ========================================
-- FUNCIÓN: update_updated_at_column
-- ========================================
-- Actualiza automáticamente el campo updated_at
-- con la fecha/hora actual cuando se modifica un registro
--
-- Uso: Se asocia a un trigger BEFORE UPDATE en tablas
-- que tienen campo updated_at TIMESTAMP WITH TIME ZONE
--
-- Ejemplo:
--   CREATE TRIGGER set_updated_at_tablename
--     BEFORE UPDATE ON tablename
--     FOR EACH ROW
--     EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Comentario de la función
COMMENT ON FUNCTION update_updated_at_column() IS 
'Trigger function para actualizar automáticamente el campo updated_at con la fecha/hora actual';
