package layers

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL1UserRole liga al usuario viewer con su rol
// `announcement_viewer` en el contexto de la escuela L1-DEMO
// (scope=school, school_id NOT NULL — F3-REQ-6.2).
//
// Id fijo (L1_USER_ROLE_VIEWER_ID): como school_id NO es NULL en L1,
// el UNIQUE compuesto (user_id, role_id, school_id, academic_unit_id)
// dispara correctamente, pero usamos id fijo + OnConflict por id
// porque es más simple y suficiente para idempotencia.
func applyL1UserRole(tx *gorm.DB) error {
	id, err := uuid.Parse(L1_USER_ROLE_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1UserRole: parse id: %w", err)
	}
	userID, err := uuid.Parse(L1_USER_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1UserRole: parse user_id: %w", err)
	}
	roleID, err := uuid.Parse(L1_ROLE_ANNOUNCEMENT_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1UserRole: parse role_id: %w", err)
	}
	schoolID, err := uuid.Parse(L1_SCHOOL_DEMO_ID)
	if err != nil {
		return fmt.Errorf("applyL1UserRole: parse school_id: %w", err)
	}
	ur := entities.UserRole{
		ID:        id,
		UserID:    userID,
		RoleID:    roleID,
		SchoolID:  &schoolID,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&ur).Error
}
