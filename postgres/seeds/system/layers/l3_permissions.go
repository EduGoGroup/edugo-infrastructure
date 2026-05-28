package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL3Permissions inserta los 3 permisos CRUD parcial sobre
// materials (read, create, update). NO incluye delete por design
// (F5-REQ-2.1): L3 valida que el sistema soporta recursos con un
// subconjunto de acciones.
//
// Todos comparten ResourceID (L3_RESOURCE_MATERIALS_ID) y
// Scope=unit (alineado con el resource).
// Idempotente vía ON CONFLICT DO NOTHING sobre id.
func applyL3Permissions(tx *gorm.DB) error {
	resourceID, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return fmt.Errorf("applyL3Permissions: parse resource_id: %w", err)
	}

	specs := []struct {
		idStr       string
		name        string
		displayName string
		action      string
	}{
		{L3_PERM_MATERIALS_READ_ID, "content.materials.read", "Ver Materiales", "read"},
		{L3_PERM_MATERIALS_CREATE_ID, "content.materials.create", "Crear Material", "create"},
		{L3_PERM_MATERIALS_UPDATE_ID, "content.materials.update", "Editar Material", "update"},
	}

	perms := make([]entities.Permission, 0, len(specs))
	for _, s := range specs {
		id, err := uuid.Parse(s.idStr)
		if err != nil {
			return fmt.Errorf("applyL3Permissions: parse id %s: %w", s.idStr, err)
		}
		perms = append(perms, entities.Permission{
			ID:          id,
			Name:        s.name,
			DisplayName: s.displayName,
			ResourceID:  resourceID,
			Action:      s.action,
			Scope:       "unit",
			IsActive:    true,
		})
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(&perms, 10).Error
}
