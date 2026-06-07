package entities

import (
	"time"

	"github.com/google/uuid"
)

// GradeHistory representa la tabla 'academic.grade_history' en PostgreSQL
// (N4 / ADR 0020). Auditoria de override: registra cada cambio de valor sobre
// una nota (academic.grades) O un componente (academic.grade_item), con quien lo
// cambio y por que.
//
// El registro apunta a EXACTAMENTE UNO de grade_id / grade_item_id (los dos
// nullable). El invariante XOR se materializa via CHECK en post_gorm.sql:
//
//	CHECK (((grade_id IS NOT NULL)::int + (grade_item_id IS NOT NULL)::int) = 1)
//
// FKs: grade_id (→academic.grades CASCADE), grade_item_id (→academic.grade_item
// CASCADE) y changed_by_membership_id (→academic.memberships RESTRICT) se
// materializan en migrations/sql/post_gorm.sql (GORM no crea FKs desde el tag
// `constraint:` sin campo de relacion). El default now() de changed_at se aplica
// via tag GORM.
type GradeHistory struct {
	ID                    uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	GradeID               *uuid.UUID `db:"grade_id" gorm:"type:uuid;index:idx_grade_history_grade;constraint:grade_history_grade_fkey,OnDelete:CASCADE" validate:"omitempty,uuid"`
	GradeItemID           *uuid.UUID `db:"grade_item_id" gorm:"type:uuid;index:idx_grade_history_item;constraint:grade_history_item_fkey,OnDelete:CASCADE" validate:"omitempty,uuid"`
	OldValue              *float64   `db:"old_value" gorm:"type:decimal(5,2)" validate:"omitempty"`
	NewValue              *float64   `db:"new_value" gorm:"type:decimal(5,2)" validate:"omitempty"`
	ChangedByMembershipID uuid.UUID  `db:"changed_by_membership_id" gorm:"type:uuid;not null;constraint:grade_history_changed_by_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	ChangedAt             time.Time  `db:"changed_at" gorm:"not null;default:now()" validate:"-"`
	Reason                *string    `db:"reason" gorm:"type:text" validate:"omitempty"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (GradeHistory) TableName() string {
	return "academic.grade_history"
}
