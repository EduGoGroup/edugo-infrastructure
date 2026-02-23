-- ============================================================
-- 001: Shared Functions (public schema)
-- Funciones compartidas entre múltiples schemas
-- ============================================================

-- Trigger function para actualizar automáticamente updated_at
CREATE OR REPLACE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

COMMENT ON FUNCTION public.update_updated_at_column() IS 'Trigger function para actualizar automáticamente el campo updated_at con la fecha/hora actual';

-- Trigger function para prevenir ciclos en jerarquía de academic_units
CREATE OR REPLACE FUNCTION public.prevent_academic_unit_cycles() RETURNS trigger
    LANGUAGE plpgsql
AS $$
DECLARE
    current_parent_id UUID;
    visited_ids UUID[];
    depth INTEGER := 0;
    max_depth INTEGER := 50;
BEGIN
    IF NEW.parent_unit_id IS NULL THEN
        RETURN NEW;
    END IF;

    current_parent_id := NEW.parent_unit_id;
    visited_ids := ARRAY[]::UUID[];

    IF NEW.id IS NOT NULL THEN
        visited_ids := array_append(visited_ids, NEW.id);
    END IF;

    WHILE current_parent_id IS NOT NULL AND depth < max_depth LOOP
        IF current_parent_id = ANY(visited_ids) THEN
            RAISE EXCEPTION 'Ciclo detectado en jerarquía: no se puede asignar % como padre de %',
                NEW.parent_unit_id, NEW.id;
        END IF;

        visited_ids := array_append(visited_ids, current_parent_id);

        SELECT parent_unit_id INTO current_parent_id
        FROM academic.academic_units
        WHERE id = current_parent_id;

        depth := depth + 1;
    END LOOP;

    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;

    RETURN NEW;
END;
$$;
