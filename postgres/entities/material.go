package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Material representa la tabla 'content.materials' en PostgreSQL (N4 / ADR 0019).
//
// Cambio vs viejo:
//   - se elimina subject/grade (varchar libre) y se gana subject_id (→academic.subjects)
//     nullable: un material puede ser general de la escuela; la guia adjunta a una
//     evaluacion se valida con subject en la capa de aplicacion (decision F1).
//   - uploaded_by_teacher_id (→auth.users) → uploaded_by_membership_id (→academic.memberships RESTRICT).
//
// Se conservan school_id (NOT NULL, CASCADE), academic_unit_id (nullable, SET NULL),
// title/description/file_*/status/is_public/processing_*, deleted_at. FKs
// cross-schema (school_id, subject_id, academic_unit_id, uploaded_by_membership_id)
// y el indice parcial idx_materials_active (WHERE deleted_at IS NULL) en post_gorm.sql.
type Material struct {
	ID                    uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID              uuid.UUID  `db:"school_id" gorm:"type:uuid;index;not null;constraint:materials_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	UploadedByMembershipID uuid.UUID `db:"uploaded_by_membership_id" gorm:"type:uuid;index;not null;constraint:materials_membership_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	SubjectID             *uuid.UUID `db:"subject_id" gorm:"type:uuid;index;constraint:materials_subject_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	AcademicUnitID        *uuid.UUID `db:"academic_unit_id" gorm:"type:uuid;index;constraint:materials_unit_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Title                 string     `db:"title" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Description           *string    `db:"description" gorm:"default:null" validate:"omitempty"`
	FileURL               string     `db:"file_url" gorm:"not null;default:''" validate:"required,url"`
	FileType              string     `db:"file_type" gorm:"not null;size:100" validate:"required,max=100"`
	FileSizeBytes         int64      `db:"file_size_bytes" gorm:"not null;default:0"`
	Status                string     `db:"status" gorm:"not null;type:varchar(50);index;default:'draft';check:materials_status_check,status IN ('draft','uploaded','processing','ready','failed')" validate:"required,oneof=draft uploaded processing ready failed"`
	ProcessingStartedAt   *time.Time `db:"processing_started_at" gorm:"default:null"`
	ProcessingCompletedAt *time.Time `db:"processing_completed_at" gorm:"default:null"`
	IsPublic              bool       `db:"is_public" gorm:"not null;default:false"`
	CreatedAt             time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt             time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt             gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Material) TableName() string {
	return "content.materials"
}
