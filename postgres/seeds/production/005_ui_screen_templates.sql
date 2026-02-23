-- ====================================================================
-- SEEDS: Templates base para los 6 patterns de la Fase 1
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ====================================================================

BEGIN;

-- Template 1: Login
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000001', 'login', 'login-basic-v1', 'Login con marca, formulario, autenticacion social', 1, '{
  "navigation": {"topBar": null},
  "zones": [
    {"id": "brand", "type": "container", "slots": [
      {"id": "app_logo", "controlType": "icon", "bind": "slot:app_logo"},
      {"id": "app_name", "controlType": "label", "style": "headline-large", "bind": "slot:app_name"},
      {"id": "app_tagline", "controlType": "label", "style": "body", "bind": "slot:app_tagline"}
    ]},
    {"id": "form", "type": "form-section", "slots": [
      {"id": "email", "controlType": "email-input", "bind": "slot:email_label", "required": true},
      {"id": "password", "controlType": "password-input", "bind": "slot:password_label", "required": true, "secureTextEntry": true},
      {"id": "remember_me", "controlType": "checkbox", "bind": "slot:remember_label"},
      {"id": "login_btn", "controlType": "filled-button", "bind": "slot:login_btn_label"},
      {"id": "forgot_password", "controlType": "text-button", "bind": "slot:forgot_password_label"}
    ]},
    {"id": "social", "type": "container", "slots": [
      {"id": "divider_text", "controlType": "label", "bind": "slot:divider_text"},
      {"id": "google_btn", "controlType": "outlined-button", "bind": "slot:google_btn_label", "icon": "google"}
    ]}
  ],
  "platformOverrides": {
    "desktop": {"distribution": "side-by-side", "weights": [0.4, 0.6], "zones": {"brand": {"panel": "left"}, "form": {"panel": "right"}, "social": {"panel": "right"}}},
    "web": {"distribution": "centered-card", "maxWidth": 480},
    "ios": {"distribution": "stacked", "zones": {"brand": {"alignment": "center"}, "social": {"visible": false}}}
  }
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Template 2: Dashboard
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000002', 'dashboard', 'dashboard-basic-v1', 'Dashboard con saludo, KPIs, actividad, acciones rapidas', 1, '{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": false}},
  "zones": [
    {"id": "greeting", "type": "container", "slots": [
      {"id": "greeting_text", "controlType": "label", "style": "headline-large", "bind": "slot:greeting_text"},
      {"id": "date_text", "controlType": "label", "style": "body", "bind": "slot:date_text"}
    ]},
    {"id": "kpis", "type": "metric-grid", "distribution": "grid", "slots": [
      {"id": "total_students", "controlType": "metric-card", "bind": "slot:kpi_students_label", "field": "total_students", "icon": "people"},
      {"id": "total_materials", "controlType": "metric-card", "bind": "slot:kpi_materials_label", "field": "total_materials", "icon": "folder"},
      {"id": "avg_score", "controlType": "metric-card", "bind": "slot:kpi_avg_score_label", "field": "avg_score", "icon": "trending_up"},
      {"id": "completion_rate", "controlType": "metric-card", "bind": "slot:kpi_completion_label", "field": "completion_rate", "icon": "check_circle"}
    ]},
    {"id": "recent_activity", "type": "simple-list", "slots": [
      {"id": "section_title", "controlType": "label", "style": "title-medium", "bind": "slot:activity_title"}
    ], "itemLayout": {"slots": [
      {"id": "activity_icon", "controlType": "icon", "field": "type_icon"},
      {"id": "activity_text", "controlType": "label", "style": "body", "field": "description"},
      {"id": "activity_time", "controlType": "label", "style": "caption", "field": "time_ago"}
    ]}},
    {"id": "quick_actions", "type": "action-group", "distribution": "flow-row", "slots": [
      {"id": "upload_material", "controlType": "outlined-button", "bind": "slot:upload_label", "icon": "upload"},
      {"id": "view_progress", "controlType": "outlined-button", "bind": "slot:progress_label", "icon": "bar_chart"}
    ]}
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Template 3: List
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000003', 'list', 'list-basic-v1', 'Lista con busqueda, filtros, estado vacio, elementos', 1, '{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": false}},
  "zones": [
    {"id": "search_zone", "type": "container", "slots": [
      {"id": "search_bar", "controlType": "search-bar", "bind": "slot:search_placeholder"}
    ]},
    {"id": "filters", "type": "container", "distribution": "flow-row", "slots": [
      {"id": "filter_all", "controlType": "chip", "bind": "slot:filter_all_label", "selected": true},
      {"id": "filter_ready", "controlType": "chip", "bind": "slot:filter_ready_label"},
      {"id": "filter_processing", "controlType": "chip", "bind": "slot:filter_processing_label"}
    ]},
    {"id": "empty_state", "type": "container", "condition": "data.isEmpty", "slots": [
      {"id": "empty_icon", "controlType": "icon", "bind": "slot:empty_icon"},
      {"id": "empty_title", "controlType": "label", "style": "headline", "bind": "slot:empty_state_title"},
      {"id": "empty_desc", "controlType": "label", "style": "body", "bind": "slot:empty_state_description"},
      {"id": "empty_action", "controlType": "filled-button", "bind": "slot:empty_action_label"}
    ]},
    {"id": "list_content", "type": "simple-list", "condition": "!data.isEmpty", "itemLayout": {"slots": [
      {"id": "item_icon", "controlType": "icon", "field": "file_type_icon"},
      {"id": "item_title", "controlType": "label", "style": "headline-small", "field": "title"},
      {"id": "item_subtitle", "controlType": "label", "style": "body-small", "field": "subtitle"},
      {"id": "item_status", "controlType": "chip", "field": "status"},
      {"id": "item_date", "controlType": "label", "style": "caption", "field": "created_at"}
    ]}}
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Template 4: Detail
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000004', 'detail', 'detail-basic-v1', 'Detalle con hero, secciones de contenido, acciones', 1, '{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": true, "actions": []}},
  "zones": [
    {"id": "hero", "type": "container", "slots": [
      {"id": "file_type_icon", "controlType": "icon", "style": "large", "field": "file_type"},
      {"id": "status_badge", "controlType": "chip", "field": "status"}
    ]},
    {"id": "header", "type": "container", "slots": [
      {"id": "title", "controlType": "label", "style": "headline-large", "field": "title"},
      {"id": "subject", "controlType": "label", "style": "body", "field": "subject"},
      {"id": "grade", "controlType": "label", "style": "body-small", "field": "grade"}
    ]},
    {"id": "details", "type": "container", "slots": [
      {"id": "file_size", "controlType": "label", "bind": "slot:file_size_label", "field": "file_size_display"},
      {"id": "uploaded_date", "controlType": "label", "bind": "slot:uploaded_label", "field": "created_at"},
      {"id": "status", "controlType": "label", "bind": "slot:status_label", "field": "status"}
    ]},
    {"id": "description", "type": "container", "slots": [
      {"id": "section_title", "controlType": "label", "style": "title-medium", "bind": "slot:description_title"},
      {"id": "description_text", "controlType": "label", "style": "body", "field": "description"}
    ]},
    {"id": "summary", "type": "container", "condition": "data.summary != null", "slots": [
      {"id": "summary_title", "controlType": "label", "style": "title-medium", "bind": "slot:summary_title"},
      {"id": "summary_content", "controlType": "label", "style": "body", "field": "summary.main_ideas"}
    ]},
    {"id": "actions", "type": "action-group", "slots": [
      {"id": "download_btn", "controlType": "filled-button", "bind": "slot:download_label", "icon": "download"},
      {"id": "take_quiz_btn", "controlType": "outlined-button", "bind": "slot:quiz_label", "icon": "quiz"}
    ]}
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Template 5: Settings
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000005', 'settings', 'settings-basic-v1', 'Configuracion con secciones agrupadas', 1, '{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": false}},
  "zones": [
    {"id": "user_card", "type": "container", "slots": [
      {"id": "avatar", "controlType": "avatar", "field": "user.avatar_url"},
      {"id": "user_name", "controlType": "label", "style": "headline-small", "field": "user.full_name"},
      {"id": "user_email", "controlType": "label", "style": "body-small", "field": "user.email"},
      {"id": "user_role", "controlType": "chip", "field": "user.role"}
    ]},
    {"id": "section_appearance", "type": "form-section", "slots": [
      {"id": "appearance_title", "controlType": "label", "style": "title-medium", "bind": "slot:appearance_title"},
      {"id": "dark_mode", "controlType": "switch", "bind": "slot:dark_mode_label", "field": "preferences.dark_mode"},
      {"id": "theme_color", "controlType": "list-item-navigation", "bind": "slot:theme_label", "field": "preferences.theme"}
    ]},
    {"id": "section_notifications", "type": "form-section", "slots": [
      {"id": "notifications_title", "controlType": "label", "style": "title-medium", "bind": "slot:notifications_title"},
      {"id": "push_notifications", "controlType": "switch", "bind": "slot:push_label", "field": "preferences.push_enabled"},
      {"id": "email_notifications", "controlType": "switch", "bind": "slot:email_label", "field": "preferences.email_enabled"}
    ]},
    {"id": "logout", "type": "container", "slots": [
      {"id": "logout_btn", "controlType": "filled-button", "bind": "slot:logout_label", "style": "error"}
    ]}
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

-- Template 6: Form
INSERT INTO ui_config.screen_templates (id, pattern, name, description, version, definition) VALUES
('a0000000-0000-0000-0000-000000000006', 'form', 'form-basic-v1', 'Formulario generico con campos dinamicos, validacion, submit/cancel', 1, '{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": true}},
  "zones": [
    {"id": "form_header", "type": "container", "slots": [
      {"id": "form_title", "controlType": "label", "style": "headline-medium", "bind": "slot:form_title"},
      {"id": "form_description", "controlType": "label", "style": "body", "bind": "slot:form_description"}
    ]},
    {"id": "form_fields", "type": "form-section", "slots": []},
    {"id": "form_actions", "type": "action-group", "distribution": "flow-row", "slots": [
      {"id": "cancel_btn", "controlType": "outlined-button", "bind": "slot:cancel_label"},
      {"id": "submit_btn", "controlType": "filled-button", "bind": "slot:submit_label"}
    ]}
  ]
}'::jsonb)
ON CONFLICT (name, version) DO NOTHING;

COMMIT;
