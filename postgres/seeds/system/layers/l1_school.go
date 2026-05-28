package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL1School siembra la escuela mínima L1-DEMO en academic.schools.
//
// Justificación (ADR-7): el contrato `scope=school` exige que el
// UserRole del viewer tenga `school_id IS NOT NULL` (ver
// post_gorm.sql:~311 — chk_user_roles_unit_requires_school y la
// regla de scope=school). Sin esta escuela L1 no podría ligar al
// usuario viewer a su rol scope=school. Esto NO es scope creep —
// es requisito de integridad referencial.
//
// Idempotente vía ON CONFLICT (id) DO NOTHING. Replica el patrón de
// upsertL0User para minimizar superficie y mantener consistencia.
func applyL1School(tx *gorm.DB) error {
	id, err := uuid.Parse(L1_SCHOOL_DEMO_ID)
	if err != nil {
		return fmt.Errorf("applyL1School: parse id: %w", err)
	}
	school := entities.School{
		ID:               id,
		Name:             L1_SCHOOL_DEMO_NAME,
		Code:             L1_SCHOOL_DEMO_CODE,
		Country:          "Chile",
		SubscriptionTier: "basic",
		MaxTeachers:      0,
		MaxStudents:      0,
		IsActive:         true,
		// ConceptTypeID nil → NULL (concept_type_id es nullable).
		Metadata: json.RawMessage(`{}`),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&school).Error
}
