package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL1Role siembra el rol `announcement_viewer` de L1 con
// scope=school. Idempotente vía ON CONFLICT (id) DO NOTHING.
//
// Replica el patrón de upsertL0Role en l0_roles.go.
func applyL1Role(tx *gorm.DB) error {
	id, err := uuid.Parse(L1_ROLE_ANNOUNCEMENT_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1Role: parse id: %w", err)
	}
	desc := "Rol read-only: solo puede ver anuncios. Usado para validar gating de UI en Fase 3."
	// landing_screen_key (ADR 0024 sub-deuda "herencia del landing"): scope=school
	// → dashboard-schooladmin. Sin landing propio caía a school.default
	// (= "dashboard-home", el dashboard básico genérico); con landing propio
	// aterriza en el dashboard de su superficie en vez del home genérico.
	landing := "dashboard-schooladmin"
	role := entities.Role{
		ID:               id,
		Name:             L1_ROLE_ANNOUNCEMENT_VIEWER_NAME,
		DisplayName:      "Visualizador de Anuncios",
		Description:      &desc,
		Scope:            "school",
		LandingScreenKey: &landing,
		IsActive:         true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&role).Error
}
