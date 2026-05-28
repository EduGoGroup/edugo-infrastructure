package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SubjectOffering representa la tabla 'subject_offerings' en PostgreSQL.
//
// Una "sesión de materia" (oferta): materia + sección + período + docente.
// Es la unidad de enseñanza e inscripción (ADR 0009 / plan 010 N1.7). La
// inscripción del alumno apunta a una sesión, no a la materia pelada.
//
// Las FKs (school/subject/academic_unit/period/teacher) se materializan en
// migrations/sql/post_gorm.sql: GORM no crea FKs desde el tag `constraint:`
// cuando la entity no declara un campo de relación (mismo caso que
// academic.subjects / academic.schedules / academic.school_invitations).
type SubjectOffering struct {
	ID                  uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID            uuid.UUID       `db:"school_id" gorm:"type:uuid;not null;index;constraint:subject_offerings_school_fkey,OnDelete:CASCADE;uniqueIndex:uq_subject_offerings_natural" validate:"required,uuid"`
	SubjectID           uuid.UUID       `db:"subject_id" gorm:"type:uuid;not null;index;constraint:subject_offerings_subject_fkey,OnDelete:CASCADE;uniqueIndex:uq_subject_offerings_natural" validate:"required,uuid"`
	AcademicUnitID      uuid.UUID       `db:"academic_unit_id" gorm:"type:uuid;not null;index;constraint:subject_offerings_unit_fkey,OnDelete:CASCADE;uniqueIndex:uq_subject_offerings_natural" validate:"required,uuid"`
	SectionLabel        *string         `db:"section_label" gorm:"type:varchar(10);uniqueIndex:uq_subject_offerings_natural" validate:"omitempty,max=10"`
	PeriodID            uuid.UUID       `db:"period_id" gorm:"type:uuid;not null;index;constraint:subject_offerings_period_fkey,OnDelete:CASCADE;uniqueIndex:uq_subject_offerings_natural" validate:"required,uuid"`
	TeacherMembershipID *uuid.UUID      `db:"teacher_membership_id" gorm:"type:uuid;index;constraint:subject_offerings_teacher_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	// Capacity esta RESERVADO para el caso universidad (cupo); NO se usa en
	// N1.7 (la politica de inscripcion es no-op auto-aprobar). Ver ADR 0009.
	Capacity  *int            `db:"capacity" gorm:"default:null" validate:"omitempty"`
	IsActive  bool            `db:"is_active" gorm:"not null;default:true"`
	Metadata  json.RawMessage `db:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time       `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SubjectOffering) TableName() string {
	return "academic.subject_offerings"
}
