package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// announcementFormSlotData es el slot_data JSON canónico de la
// pantalla L2 "announcement-form".
//
// F3.1 (2ª pasada, plan 004): migrada al patrón delta. Hereda los 3
// default_actions de form-basic-v1 (save_new create-only +
// permission=create, save edit-only + permission=update, delete
// edit-only + permission=delete) con $resource$ resuelto a
// "academic.announcements" desde required_permission. El desdoblado
// save_new/save (uno create-only con permission=create y otro edit-only
// con permission=update, para no ocultar el botón a usuarios con solo
// `create`) ya forma parte del template, así que se conserva. El scope
// inline "form" del legacy era render-equivalente al "form-submit" del
// template (normalizeScope en el FE). El conjunto semántico
// {event_id, permission} queda idéntico (verificado por el harness).
const announcementFormSlotData = `{
  "title": "Anuncio",
  "page_title": "Anuncio",
  "edit_title": "Editar anuncio",
  "fields": [
    { "key": "title", "label": "Título", "type": "text", "required": true },
    { "key": "body", "label": "Cuerpo", "type": "textarea", "required": true },
    { "key": "scope", "label": "Alcance", "type": "select", "required": true, "default": "school", "options": [
      { "label": "Toda la escuela", "value": "school" },
      { "label": "Solo la unidad", "value": "unit" }
    ] },
    { "key": "is_pinned", "label": "Fijar arriba", "type": "toggle", "default": "false" },
    { "key": "published_at", "label": "Fecha de publicación", "type": "datetime", "required": false }
  ],
  "api_prefix": "platform"
}`

// applyL2Screens siembra la ScreenInstance "announcement-form" de L2
// en ui_config.screen_instances. Idempotente: UPSERT DoNothing sobre id.
func applyL2Screens(tx *gorm.DB) error {
	description := "Formulario de creación/edición de anuncios"
	requiredPermission := "academic.announcements.read"

	instance := entities.ScreenInstance{
		ID:                 mustParseUUID(L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID, "L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID"),
		ScreenKey:          L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
		TemplateID:         mustParseUUID(L0_SCREEN_TPL_FORM_ID, "L0_SCREEN_TPL_FORM_ID"),
		Name:               "Formulario de anuncio",
		Description:        &description,
		SlotData:           json.RawMessage([]byte(announcementFormSlotData)),
		Scope:              "school",
		RequiredPermission: &requiredPermission,
		HandlerKey:         nil,
		IsActive:           true,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&instance).Error; err != nil {
		return fmt.Errorf("applyL2Screens: upsert screen_instances: %w", err)
	}
	return nil
}
