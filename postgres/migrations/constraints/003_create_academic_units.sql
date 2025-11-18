-- Constraints para tabla academic_units

ALTER TABLE academic_units ADD CONSTRAINT academic_units_parent_fkey 
    FOREIGN KEY (parent_unit_id) REFERENCES academic_units(id) ON DELETE SET NULL;

ALTER TABLE academic_units ADD CONSTRAINT academic_units_school_fkey 
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;

ALTER TABLE academic_units ADD CONSTRAINT academic_units_unique_code 
    UNIQUE(school_id, code, academic_year);

ALTER TABLE academic_units ADD CONSTRAINT academic_units_no_self_reference 
    CHECK (id != parent_unit_id);

ALTER TABLE academic_units ADD CONSTRAINT academic_units_type_check 
    CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department'));

-- Función para prevenir ciclos en jerarquía
CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
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
        FROM academic_units
        WHERE id = current_parent_id;

        depth := depth + 1;
    END LOOP;

    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prevent_academic_unit_cycles
    BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_units
    FOR EACH ROW
    EXECUTE FUNCTION prevent_academic_unit_cycles();

-- Vista para árbol jerárquico (CTE recursivo)
CREATE OR REPLACE VIEW v_academic_unit_tree AS
WITH RECURSIVE unit_hierarchy AS (
    SELECT
        id,
        parent_unit_id,
        school_id,
        name,
        code,
        type,
        level,
        academic_year,
        1 AS depth,
        ARRAY[id] AS path,
        name::TEXT AS full_path
    FROM academic_units
    WHERE parent_unit_id IS NULL
      AND deleted_at IS NULL

    UNION ALL

    SELECT
        au.id,
        au.parent_unit_id,
        au.school_id,
        au.name,
        au.code,
        au.type,
        au.level,
        au.academic_year,
        uh.depth + 1,
        uh.path || au.id,
        (uh.full_path || ' > ' || au.name)::TEXT
    FROM academic_units au
    INNER JOIN unit_hierarchy uh ON au.parent_unit_id = uh.id
    WHERE au.deleted_at IS NULL
)
SELECT
    uh.*,
    s.name AS school_name,
    s.code AS school_code
FROM unit_hierarchy uh
LEFT JOIN schools s ON uh.school_id = s.id
ORDER BY uh.school_id, uh.path;

COMMENT ON VIEW v_academic_unit_tree IS 'Vista con árbol jerárquico completo de unidades académicas';
