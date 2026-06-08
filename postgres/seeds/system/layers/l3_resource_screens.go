package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL3ResourceScreens vincula el resource "materials" con su pantalla
// default (F5-REQ-3.3):
//   - materials-list (screen_type=list, is_default=true)
//
// El screen_key `materials-list` apunta a una pantalla NATIVA (Compose)
// que NO consume slot_data SDUI, pero SÍ requiere una screen_instance
// MÍNIMA homónima para satisfacer la FK fk_resource_screens_screen_key
// (resource_screens.screen_key → screen_instances.screen_key). Esa
// screen_instance la siembra applyL3Screens (debe correr ANTES); mismo
// patrón que `join-requests-inbox` / `batch-enroll` en L4. El SDUI engine
// no la renderiza: el FE intercepta el screen_key.
// El mapping `form` (material-form) fue podado junto con su ScreenInstance
// (poda SDUI material 2026-06-07; ver l3_screens.go) y NO se resiembra:
// no tiene mapping → no hay FK que satisfacer.
//
// IDs derivados determinísticamente vía SHA1 sobre (resource_id,
// screen_type) — replica el patrón de upsertL0ResourceScreens. Esto
// hace que los IDs sean estables entre re-seeds y consistentes en el
// hash de exports JSON.
//
// Conflict target (resource_id, screen_type) DoNothing.
func applyL3ResourceScreens(tx *gorm.DB) error {
	materialsID, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return fmt.Errorf("applyL3ResourceScreens: parse resource_id: %w", err)
	}

	idList := uuid.NewSHA1(uuid.NameSpaceOID, []byte(materialsID.String()+":list"))

	mappings := []entities.ResourceScreen{
		{
			ID:          idList,
			ResourceID:  materialsID,
			ResourceKey: L3_RESOURCE_MATERIALS_KEY,
			ScreenKey:   L3_SCREEN_KEY_MATERIALS_LIST,
			ScreenType:  "list",
			IsDefault:   true,
			SortOrder:   0,
			IsActive:    true,
		},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "screen_type"}},
		DoNothing: true,
	}).Create(&mappings).Error; err != nil {
		return fmt.Errorf("applyL3ResourceScreens: upsert resource_screens: %w", err)
	}
	return nil
}
