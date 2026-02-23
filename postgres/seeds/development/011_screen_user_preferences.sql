-- =============================================================================
-- EduGo Development Seeds — 011_screen_user_preferences.sql
-- =============================================================================
-- Preferencias de pantalla para usuarios de prueba.
-- Usa subquery con SELECT para tolerancia si screen_instances no tiene el registro.
-- =============================================================================

DO $$
BEGIN
  INSERT INTO ui_config.screen_user_preferences (id, screen_instance_id, user_id, preferences)
  SELECT 'ff000000-0000-0000-0000-000000000001', si.id, '00000000-0000-0000-0000-000000000001'::uuid, '{"dark_mode": true, "language": "es"}'::jsonb
  FROM ui_config.screen_instances si WHERE si.screen_key = 'app-settings'
  ON CONFLICT (screen_instance_id, user_id) DO NOTHING;

  INSERT INTO ui_config.screen_user_preferences (id, screen_instance_id, user_id, preferences)
  SELECT 'ff000000-0000-0000-0000-000000000002', si.id, '00000000-0000-0000-0000-000000000005'::uuid, '{"dark_mode": false, "language": "es", "push_enabled": true}'::jsonb
  FROM ui_config.screen_instances si WHERE si.screen_key = 'app-settings'
  ON CONFLICT (screen_instance_id, user_id) DO NOTHING;
EXCEPTION WHEN OTHERS THEN
  RAISE NOTICE 'screen_user_preferences: skipped (%) ', SQLERRM;
END;
$$;
