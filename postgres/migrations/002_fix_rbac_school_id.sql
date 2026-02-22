-- ============================================================
-- MIGRACIÓN 002: Poblar user_roles.school_id desde memberships
-- Fecha: 2026-02-22
-- R4 del análisis arquitectónico Opus
-- Problema: user_roles.school_id está en NULL para todos los
--           registros con scope 'school'. Debe poblarse desde
--           la primera membresía activa del usuario.
-- ============================================================

-- Asigna school_id desde la primera membresía activa del usuario.
-- Solo actualiza registros donde school_id es NULL (no afecta super_admin
-- ni roles con scope 'system' que legítimamente no tienen school_id).
UPDATE user_roles ur
SET school_id = (
    SELECT m.school_id
    FROM memberships m
    WHERE m.user_id = ur.user_id
      AND m.is_active = true
    ORDER BY m.enrolled_at ASC
    LIMIT 1
)
WHERE ur.school_id IS NULL
  AND EXISTS (
    SELECT 1 FROM memberships m
    WHERE m.user_id = ur.user_id AND m.is_active = true
  );

-- Verificación: reporta cuántos quedaron poblados vs sin school_id
DO $$
DECLARE
    updated_count    INTEGER;
    still_null_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO updated_count    FROM user_roles WHERE school_id IS NOT NULL;
    SELECT COUNT(*) INTO still_null_count FROM user_roles WHERE school_id IS NULL;
    RAISE NOTICE 'user_roles con school_id poblado: %, sin school_id (super_admin / system): %',
        updated_count, still_null_count;
END $$;
