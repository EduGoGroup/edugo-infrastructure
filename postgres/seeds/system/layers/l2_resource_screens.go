package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL2ResourceScreens vincula el resource "announcements" con la
// instance "announcement-form" como pantalla secundaria de tipo "form".
// Idempotente: UPSERT DoNothing sobre (resource_id, screen_type).
func applyL2ResourceScreens(tx *gorm.DB) error {
	mapping := entities.ResourceScreen{
		ID:          mustParseUUID(L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID, "L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID"),
		ResourceID:  mustParseUUID(L0_RESOURCE_ANNOUNCEMENTS_ID, "L0_RESOURCE_ANNOUNCEMENTS_ID"),
		ResourceKey: L0_RESOURCE_ANNOUNCEMENTS_KEY,
		ScreenKey:   L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
		ScreenType:  "form",
		IsDefault:   false,
		SortOrder:   1,
		IsActive:    true,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "screen_type"}},
		DoNothing: true,
	}).Create(&mapping).Error; err != nil {
		return fmt.Errorf("applyL2ResourceScreens: upsert resource_screens: %w", err)
	}
	return nil
}
