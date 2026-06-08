package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Poda SDUI material (2026-06-07) + corrección F2 (2026-06-08):
//
// La poda de la sesión anterior eliminó AMBAS ScreenInstances L3
// (`materials-list` y `material-form`) por considerarlas código muerto.
// Pero `materials-list` SÍ tiene un mapping `resource_screens` (item de
// menú, is_default), y existe la FK
// `fk_resource_screens_screen_key`: resource_screens.screen_key →
// screen_instances.screen_key. Sin la screen_instance, un recreate limpio
// falla en L3 con violación 23503 (la FK no tiene destino).
//
// Patrón correcto del codebase para una pantalla NATIVA que es ítem de
// menú: mantener una screen_instance MÍNIMA solo para satisfacer la FK
// (el SDUI engine NUNCA la renderiza; el FE intercepta el screen_key y
// pinta el composable nativo). Lo hacen `batch-enroll` y
// `join-requests-inbox` en L4 (ver system/l4/screen_instances_*.go).
//
// Por eso `materials-list` se RESTAURA aquí como screen_instance mínima.
//
// `material-form` SE QUEDA PODADO: no tiene mapping resource_screen → no
// hay FK que satisfacer → su screen_instance sigue eliminada. NO se
// resiembra.

// materialsListSlotData es el slot_data MÍNIMO de la screen_instance
// `materials-list`. NO se renderiza (la pantalla es NATIVA Compose); se
// conserva mínimo y válido (list-basic-v1) por higiene, igual que
// batch-enroll / join-requests-inbox en L4. El `api_prefix` correcto
// (learning) vive en el código nativo, no aquí.
const materialsListSlotData = `{
  "title": "Materiales",
  "columns": [
    {"key": "title", "label": "Tema"},
    {"key": "status", "label": "Estado"}
  ]
}`

// applyL3Screens siembra la screen_instance MÍNIMA `materials-list` en
// ui_config.screen_instances. Existe SOLO para satisfacer la FK
// `fk_resource_screens_screen_key` del mapping de menú homónimo
// (applyL3ResourceScreens); el SDUI engine no la renderiza porque la
// pantalla de material en la app es NATIVA (Compose). Mismo patrón que
// `batch-enroll` / `join-requests-inbox` en L4.
//
// Idempotente: UPSERT DoNothing sobre id.
func applyL3Screens(tx *gorm.DB) error {
	description := "Listado de materiales (pantalla nativa; instancia mínima para satisfacer la FK del menú)"
	requiredPermission := "content.materials.read"

	instance := entities.ScreenInstance{
		ID:                 mustParseUUID(L3_SCREEN_INSTANCE_MATERIALS_LIST_ID, "L3_SCREEN_INSTANCE_MATERIALS_LIST_ID"),
		ScreenKey:          L3_SCREEN_KEY_MATERIALS_LIST,
		TemplateID:         mustParseUUID(L0_SCREEN_TPL_LIST_ID, "L0_SCREEN_TPL_LIST_ID"),
		Name:               "Materiales",
		Description:        &description,
		SlotData:           json.RawMessage([]byte(materialsListSlotData)),
		Scope:              "unit",
		RequiredPermission: &requiredPermission,
		HandlerKey:         nil,
		IsActive:           true,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&instance).Error; err != nil {
		return fmt.Errorf("applyL3Screens: upsert screen_instances: %w", err)
	}
	return nil
}
