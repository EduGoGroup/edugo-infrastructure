-- ====================================================================
-- TESTING: Datos de prueba para screen config (Dynamic UI)
-- VERSION: postgres/v0.18.0
-- ====================================================================
-- Requiere que seeds 006, 007, 008 esten ejecutados.
-- Usuarios de referencia (de testing/001_demo_users.sql):
--   admin:     a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
--   teacher 1: a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22
--   teacher 2: a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33
--   student 1: a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44

-- Template de prueba adicional: lista con variante grid
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-0000000000f1', 'list', 'list-grid-v1', 'Lista con variante de grilla para desktop', 1, '{
  "navigation": {
    "topBar": {"title": "slot:page_title", "showBack": false}
  },
  "zones": [
    {
      "id": "search_zone",
      "type": "container",
      "slots": [
        {"id": "search_bar", "controlType": "search-bar", "bind": "slot:search_placeholder"}
      ]
    },
    {
      "id": "list_content",
      "type": "simple-list",
      "distribution": "grid",
      "itemLayout": {
        "slots": [
          {"id": "item_image", "controlType": "image", "field": "thumbnail"},
          {"id": "item_title", "controlType": "label", "style": "headline-small", "field": "title"},
          {"id": "item_subtitle", "controlType": "label", "style": "body-small", "field": "subtitle"}
        ]
      }
    }
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Instancia de prueba: lista de evaluaciones (reutiliza template list-basic-v1)
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-0000000000f1', 'assessments-list',
 'a0000000-0000-0000-0000-000000000003', 'Lista de Evaluaciones', 'Pantalla de prueba para evaluaciones',
 '{
   "page_title": "Assessments",
   "search_placeholder": "Search assessments...",
   "filter_all_label": "All",
   "filter_ready_label": "Published",
   "filter_processing_label": "Draft",
   "empty_icon": "quiz",
   "empty_state_title": "No assessments yet",
   "empty_state_description": "Create your first assessment",
   "empty_action_label": "Create Assessment"
 }'::jsonb,
 '[
   {"id": "item-click", "trigger": "item_click", "type": "NAVIGATE", "config": {"target": "assessment-detail", "params": {"id": "{item.id}"}}},
   {"id": "pull-refresh", "trigger": "pull_refresh", "type": "REFRESH"}
 ]'::jsonb,
 '/v1/assessments',
 '{"method": "GET", "pagination": {"type": "offset", "pageSize": 20}}'::jsonb,
 'unit', 'assessments:read')
ON CONFLICT (screen_key) DO NOTHING;

-- Instancia de prueba: login en espanol
INSERT INTO ui_config.screen_instances (id, screen_key, template_id, name, description, slot_data, actions, data_endpoint, data_config, scope, required_permission) VALUES
('b0000000-0000-0000-0000-0000000000f2', 'app-login-es',
 'a0000000-0000-0000-0000-000000000001', 'Login (Espanol)', 'Variante de login con textos en espanol',
 '{
   "app_logo": "edugo_logo",
   "app_name": "EduGo",
   "app_tagline": "Aprender es facil",
   "email_label": "Correo electronico",
   "password_label": "Contrasena",
   "remember_label": "Recordarme",
   "login_btn_label": "Iniciar Sesion",
   "forgot_password_label": "Olvidaste tu contrasena?",
   "divider_text": "o continua con",
   "google_btn_label": "Google"
 }'::jsonb,
 '[
   {
     "id": "submit-login",
     "trigger": "button_click",
     "triggerSlotId": "login_btn",
     "type": "SUBMIT_FORM",
     "config": {
       "endpoint": "/v1/auth/login",
       "method": "POST",
       "fieldMapping": {"email": "email", "password": "password"},
       "onSuccess": {"type": "NAVIGATE", "config": {"target": "dashboard-home"}}
     }
   }
 ]'::jsonb,
 NULL, '{}'::jsonb, 'system', NULL)
ON CONFLICT (screen_key) DO NOTHING;

-- Preferencias de usuario de ejemplo para teacher (a2eebc99...)
INSERT INTO ui_config.screen_user_preferences (screen_instance_id, user_id, preferences) VALUES
('b0000000-0000-0000-0000-000000000002', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
 '{"compact_view": true, "hide_kpis": false}'::jsonb)
ON CONFLICT (screen_instance_id, user_id) DO NOTHING;

INSERT INTO ui_config.screen_user_preferences (screen_instance_id, user_id, preferences) VALUES
('b0000000-0000-0000-0000-000000000006', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
 '{"dark_mode": true, "theme": "indigo", "language": "es", "push_enabled": true, "email_enabled": false}'::jsonb)
ON CONFLICT (screen_instance_id, user_id) DO NOTHING;

-- Preferencias de usuario de ejemplo para admin (a1eebc99...)
INSERT INTO ui_config.screen_user_preferences (screen_instance_id, user_id, preferences) VALUES
('b0000000-0000-0000-0000-000000000006', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
 '{"dark_mode": false, "language": "en", "push_enabled": true, "email_enabled": true}'::jsonb)
ON CONFLICT (screen_instance_id, user_id) DO NOTHING;
