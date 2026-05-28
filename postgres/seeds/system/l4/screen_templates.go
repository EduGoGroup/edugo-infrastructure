package l4

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UUIDs propios de L4 para los templates adicionales (§6 design.md).
// Se elige el rango a4xxxxxx para distinguirlos visualmente de los
// UUIDs L0 (3xxxxxxx) y de los hipotéticos L1..L3.
// Poda F2 (plan 004-permisologia-mvp): el template L4
// master-detail-basic-v1 (UUID a4000000-...004, 0 instancias) se
// retiró del MVP. assessments-form usa el template master-detail de L0
// (L0_SCREEN_TPL_MASTER_DETAIL_ID_REF = 30000000-...004), no este. El
// UUID a4000000-...004 queda libre para reuso futuro.
const (
	l4TplLoginV1ID          = "a4000000-0000-0000-0000-000000000001"
	l4TplDashboardV1ID      = "a4000000-0000-0000-0000-000000000002"
	l4TplSettingsV1ID       = "a4000000-0000-0000-0000-000000000003"
	l4TplSettingsSystemV1ID = "a4000000-0000-0000-0000-000000000005"
)

// loginBasicV1Definition — template de login con marca + formulario +
// autenticación social. Derivado del inventario legacy
// [archivado pre-Fase-6] data.go:605 (no copy-paste literal — se
// reformatea y se preserva la estructura `zones` que valida el SDUI
// engine del KMP en modules/sdui-engine/.../model/Zone.kt).
const loginBasicV1Definition = `{
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
}`

// dashboardBasicV1Definition — template de dashboard con greeting,
// KPIs, actividad reciente y acciones rápidas. Derivado del legacy
// [archivado pre-Fase-6] data.go:606.
const dashboardBasicV1Definition = `{
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
}`

// settingsBasicV1Definition — template de configuración del usuario
// (apariencia, notificaciones, logout). Derivado del legacy
// [archivado pre-Fase-6] data.go:609.
const settingsBasicV1Definition = `{
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
}`

// Poda F2 (plan 004-permisologia-mvp): la definición
// masterDetailBasicV1Definition se retiró junto con el template
// master-detail-basic-v1 (0 instancias L4).

// settingsSystemV1Definition — template de configuración del sistema
// (cache, info de versiones). Derivado del legacy
// [archivado pre-Fase-6] data.go:612.
const settingsSystemV1Definition = `{
  "navigation": {"topBar": {"title": "slot:page_title", "showBack": false}},
  "zones": [
    {"id": "section_cache", "type": "form-section", "slots": [
      {"id": "cache_title", "controlType": "label", "style": "title-medium", "bind": "slot:cache_title"},
      {"id": "cache_description", "controlType": "label", "style": "body-small", "bind": "slot:cache_description"},
      {"id": "clear_cache_btn", "controlType": "filled-button", "bind": "slot:clear_cache_label", "event_id": "clear_all_cache", "icon": "trash-2", "style": "error"}
    ]},
    {"id": "section_info", "type": "form-section", "slots": [
      {"id": "info_title", "controlType": "label", "style": "title-medium", "bind": "slot:info_title"},
      {"id": "app_version", "controlType": "list-item", "bind": "slot:app_version_label", "value": "slot:app_version_value"},
      {"id": "schema_version", "controlType": "list-item", "bind": "slot:schema_version_label", "value": "slot:schema_version_value"}
    ]}
  ]
}`

// (helper strPtr definido en concept_types.go).

// ApplyScreenTemplates siembra los 4 templates adicionales de L4
// (login, dashboard, settings-user, settings-system). Los 3 templates
// base (list/detail/form-basic-v1) están en L0 y su `definition`
// canónica vive en seeds/system/layers/l0_screens.go (refactor de
// excepción aplicado en B3 — ver phase-6-layer-l4/design.md §4).
//
// Inventario fuente: [archivado pre-Fase-6] data.go:604-613
// (screenTemplateSeedRows). Decisiones de redefinición:
//   - login-basic-v1: preservado tal cual (v1).
//   - dashboard-basic-v1: preservado (v1). Único pattern dashboard.
//   - settings-basic-v1: preservado (v1) como settings del usuario.
//   - settings-system-v1: preservado (v1) — pantalla de mantenimiento
//     diferenciada de settings-basic-v1. Mismo pattern "settings" con
//     name único permite que ambos coexistan vía el unique
//     (name, version).
//
// Poda F2 (plan 004-permisologia-mvp): el template master-detail-basic-v1
// (version=2, 0 instancias L4) se retiró. El master-detail que usa
// assessments-form vive en L0.
//
// Idempotente: UPSERT con conflict target (name, version) DoNothing,
// consistente con upsertL0ScreenTemplates.
func ApplyScreenTemplates(tx *gorm.DB) error {
	templates := buildL4ScreenTemplates()
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "version"}},
		DoNothing: true,
	}).Create(&templates).Error; err != nil {
		return fmt.Errorf("ApplyScreenTemplates: upsert screen_templates: %w", err)
	}
	return nil
}

// buildL4ScreenTemplates materializa los 4 templates adicionales de L4
// como entities.ScreenTemplate. Helper compartido por
// ApplyScreenTemplates y por el accessor público l4.ScreenTemplates().
func buildL4ScreenTemplates() []entities.ScreenTemplate {
	return []entities.ScreenTemplate{
		{
			ID:          mustParseL4UUID(l4TplLoginV1ID, "l4TplLoginV1ID"),
			Pattern:     "login",
			Name:        "login-basic-v1",
			Description: strPtr("Login con marca, formulario, autenticación social"),
			Version:     1,
			Definition:  json.RawMessage([]byte(loginBasicV1Definition)),
			IsActive:    true,
		},
		{
			ID:          mustParseL4UUID(l4TplDashboardV1ID, "l4TplDashboardV1ID"),
			Pattern:     "dashboard",
			Name:        "dashboard-basic-v1",
			Description: strPtr("Dashboard con saludo, KPIs, actividad reciente, acciones rápidas"),
			Version:     1,
			Definition:  json.RawMessage([]byte(dashboardBasicV1Definition)),
			IsActive:    true,
		},
		{
			ID:          mustParseL4UUID(l4TplSettingsV1ID, "l4TplSettingsV1ID"),
			Pattern:     "settings",
			Name:        "settings-basic-v1",
			Description: strPtr("Configuración del usuario con secciones agrupadas"),
			Version:     1,
			Definition:  json.RawMessage([]byte(settingsBasicV1Definition)),
			IsActive:    true,
		},
		{
			ID:          mustParseL4UUID(l4TplSettingsSystemV1ID, "l4TplSettingsSystemV1ID"),
			Pattern:     "settings",
			Name:        "settings-system-v1",
			Description: strPtr("Configuración del sistema con secciones, acciones y datos informativos"),
			Version:     1,
			Definition:  json.RawMessage([]byte(settingsSystemV1Definition)),
			IsActive:    true,
		},
	}
}

// mustParseL4UUID convierte un literal UUID hardcodeado en este
// archivo a uuid.UUID, paniqueando si está corrupto. Coincide con la
// política de mustParseUUID en seeds/system/layers/l0_screens.go: los
// literales son validados visualmente y un fallo aquí indica
// corrupción del archivo, no error runtime.
func mustParseL4UUID(s, name string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("ApplyScreenTemplates: parse %s=%q: %v", name, s, err))
	}
	return u
}
