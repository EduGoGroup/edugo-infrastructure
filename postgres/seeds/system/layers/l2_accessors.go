package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Accessors públicos de L2 — espejan el patrón de
// `system/l4/accessors.go`. L2 sólo siembra screen_instances y
// resource_screens; esos son los únicos accessors expuestos.

// L2ScreenInstances retorna la instancia `announcement-form`
// sembrada por L2 (1 fila). Mirror determinístico de applyL2Screens
// en l2_screens.go.
func L2ScreenInstances() ([]entities.ScreenInstance, error) {
	id, err := uuid.Parse(L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L2ScreenInstances: parse id: %w", err)
	}
	templateID, err := uuid.Parse(L0_SCREEN_TPL_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L2ScreenInstances: parse template id: %w", err)
	}
	description := "Formulario de creación/edición de anuncios"
	requiredPermission := "academic.announcements.read"
	return []entities.ScreenInstance{
		{
			ID:                 id,
			ScreenKey:          L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
			TemplateID:         templateID,
			Name:               "Formulario de anuncio",
			Description:        &description,
			SlotData:           json.RawMessage([]byte(announcementFormSlotData)),
			Scope:              "school",
			RequiredPermission: &requiredPermission,
			HandlerKey:         nil,
			IsActive:           true,
		},
	}, nil
}

// L2ResourceScreens retorna el mapping announcements ↔ announcement-form
// sembrado por L2 (1 fila). Mirror determinístico de
// applyL2ResourceScreens en l2_resource_screens.go.
func L2ResourceScreens() ([]entities.ResourceScreen, error) {
	id, err := uuid.Parse(L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L2ResourceScreens: parse id: %w", err)
	}
	resourceID, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return nil, fmt.Errorf("L2ResourceScreens: parse resource_id: %w", err)
	}
	return []entities.ResourceScreen{
		{
			ID:          id,
			ResourceID:  resourceID,
			ResourceKey: L0_RESOURCE_ANNOUNCEMENTS_KEY,
			ScreenKey:   L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
			ScreenType:  "form",
			IsDefault:   false,
			SortOrder:   1,
			IsActive:    true,
		},
	}, nil
}
