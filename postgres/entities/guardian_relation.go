package entities

import (
	"time"

	"github.com/google/uuid"
)

// GuardianRelation representa la tabla 'guardian_relations' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/034_academic_guardian_relations.sql
//
// Representa la relación entre un apoderado (guardian) y un estudiante.
// Define el tipo de relación familiar o legal entre ellos.
type GuardianRelation struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	GuardianID       uuid.UUID  `db:"guardian_id" gorm:"type:uuid;index;not null;constraint:guardian_relations_guardian_fkey,OnDelete:CASCADE;uniqueIndex:guardian_relations_unique" validate:"required,uuid"`
	StudentID        uuid.UUID  `db:"student_id" gorm:"type:uuid;index;not null;constraint:guardian_relations_student_fkey,OnDelete:CASCADE;uniqueIndex:guardian_relations_unique" validate:"required,uuid"`
	SchoolID         uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null;constraint:guardian_relations_school_fkey,OnDelete:CASCADE;uniqueIndex:guardian_relations_unique" validate:"required,uuid"`
	AcademicUnitID   *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index;constraint:guardian_relations_unit_fkey,OnDelete:CASCADE" validate:"omitempty,uuid"`
	RelationshipType string     `db:"relationship_type" gorm:"not null;type:varchar(50);default:'parent';check:guardian_relations_type_check,relationship_type IN ('parent','guardian','tutor','other')" validate:"required,oneof=parent guardian tutor other"`
	IsPrimary        bool       `db:"is_primary" gorm:"not null;default:false"`
	IsActive         bool       `db:"is_active" gorm:"not null;default:true"`
	Status           string     `db:"status" gorm:"not null;type:varchar(20);default:'active';index:idx_guardian_relations_status;check:guardian_relations_status_check,status IN ('pending','active','rejected','revoked')" validate:"required,oneof=pending active rejected revoked"`
	CreatedBy        *uuid.UUID `db:"created_by" gorm:"type:uuid;constraint:guardian_relations_created_by_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (GuardianRelation) TableName() string {
	return "academic.guardian_relations"
}
