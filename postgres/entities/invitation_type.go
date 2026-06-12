package entities

import (
	"time"

	"github.com/google/uuid"
)

// InvitationType representa la tabla 'invitation_types' en el schema academic.
//
// Catalogo GLOBAL de tipos de invitacion (MP-08): teacher/student/guardian/
// coordinator/admin/assistant. Reemplaza el CHECK inline hardcodeado que
// vivia en el tag de school_invitations/school_join_requests/memberships; esas
// entities ahora referencian este tipo por id (invitation_type_id). Los valores
// (key/label/requires_unit) se siembran en F1 (este F0 solo define el esquema).
//
// El indice UNIQUE de `key` lo materializa GORM desde el tag uniqueIndex; el
// trigger set_updated_at vive en sql/post_gorm.sql (seccion academic).
type InvitationType struct {
	ID           uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Key          string    `db:"key" gorm:"uniqueIndex:invitation_types_key_key;not null;type:varchar(50)" validate:"required,min=1,max=50"`
	Label        string    `db:"label" gorm:"not null;type:varchar(100)" validate:"required,min=1,max=100"`
	RequiresUnit bool      `db:"requires_unit" gorm:"not null;default:false"`
	CreatedAt    time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt    time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (InvitationType) TableName() string {
	return "academic.invitation_types"
}
