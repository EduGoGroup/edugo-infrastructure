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
// El screen_key `materials-list` NO tiene ScreenInstance: la pantalla de
// material en la app es NATIVA (Compose) y no consume slot_data SDUI.
// El resolver solo necesita que el menú exponga el screen_key — mismo
// patrón que `material-detail` / `join-requests-inbox` / `batch-enroll`.
// El mapping `form` (material-form) fue podado junto con su ScreenInstance
// (poda SDUI material 2026-06-07; ver l3_screens.go).
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
