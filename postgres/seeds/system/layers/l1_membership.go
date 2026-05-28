package layers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyL1Membership siembra la membresía del viewer en la escuela
// L1-DEMO dentro de academic.memberships.
//
// Justificación: identity API (POST /auth/switch-context →
// usecase/context/switch_context.go) valida que exista una fila
// activa en academic.memberships antes de emitir un JWT con el
// contexto del rol. Sin esta fila, switch-context retorna 403
// NO_MEMBERSHIP y el viewer nunca carga sus permisos efectivos.
//
// El CHECK constraint memberships_role_check solo admite
// (teacher|student|guardian|coordinator|admin|assistant). El rol
// semántico del viewer vive en iam.user_roles vía
// announcement_viewer; aquí usamos "assistant" como valor mínimo
// que satisface el CHECK sin implicar permisos extra — la
// resolución real de permisos sigue ocurriendo en
// FindUserContextForSchool sobre iam.user_roles.
//
// Idempotente vía OnConflict por id.
func applyL1Membership(tx *gorm.DB) error {
	id, err := uuid.Parse(L1_MEMBERSHIP_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1Membership: parse id: %w", err)
	}
	userID, err := uuid.Parse(L1_USER_VIEWER_ID)
	if err != nil {
		return fmt.Errorf("applyL1Membership: parse user_id: %w", err)
	}
	schoolID, err := uuid.Parse(L1_SCHOOL_DEMO_ID)
	if err != nil {
		return fmt.Errorf("applyL1Membership: parse school_id: %w", err)
	}
	m := entities.Membership{
		ID:         id,
		UserID:     userID,
		SchoolID:   schoolID,
		Role:       L1_MEMBERSHIP_ROLE,
		Metadata:   json.RawMessage(`{}`),
		IsActive:   true,
		EnrolledAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}
