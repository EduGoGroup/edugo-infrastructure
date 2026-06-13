package layers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
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
// invitation_type_id referencia academic.invitation_types (MP-08). El tipo
// semántico del viewer vive en iam.user_roles vía announcement_viewer; aquí
// usamos el tipo "assistant" (L1_MEMBERSHIP_ROLE) como valor mínimo sin permisos
// extra — la resolución real de permisos sigue ocurriendo en
// FindUserContextForSchool sobre iam.user_roles. El catálogo invitation_types se
// siembra en el paso 0 de l1Layer.Apply (l4.ApplyInvitationTypes).
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
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, L1_MEMBERSHIP_ROLE)
	if err != nil {
		return fmt.Errorf("applyL1Membership: %w", err)
	}
	m := entities.Membership{
		ID:               id,
		UserID:           userID,
		SchoolID:         schoolID,
		InvitationTypeID: invitationTypeID,
		Metadata:         json.RawMessage(`{}`),
		IsActive:         true,
		EnrolledAt:       time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}
