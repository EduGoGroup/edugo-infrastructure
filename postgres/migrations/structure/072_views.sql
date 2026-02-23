-- ============================================================
-- 072: Views
-- Schema: academic
-- Vistas con referencias schema-qualified
-- ============================================================

-- Vista recursiva del árbol de unidades académicas
CREATE OR REPLACE VIEW academic.v_academic_unit_tree AS
WITH RECURSIVE unit_hierarchy AS (
    SELECT id, parent_unit_id, school_id, name, code, type, level, academic_year,
        1 AS depth, ARRAY[id] AS path, name::text AS full_path
    FROM academic.academic_units
    WHERE parent_unit_id IS NULL AND deleted_at IS NULL
    UNION ALL
    SELECT au.id, au.parent_unit_id, au.school_id, au.name, au.code, au.type, au.level, au.academic_year,
        uh.depth + 1, uh.path || au.id, uh.full_path || ' > ' || au.name::text
    FROM academic.academic_units au
    JOIN unit_hierarchy uh ON au.parent_unit_id = uh.id
    WHERE au.deleted_at IS NULL
)
SELECT uh.*, s.name AS school_name, s.code AS school_code
FROM unit_hierarchy uh
LEFT JOIN academic.schools s ON uh.school_id = s.id
ORDER BY uh.school_id, uh.path;
