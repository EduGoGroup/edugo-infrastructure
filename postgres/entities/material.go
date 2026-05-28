package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Material representa la tabla 'materials' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 005_create_materials.up.sql
// Usada por: api-mobile, worker
type Material struct {
	ID                    uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID              uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null;constraint:materials_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	UploadedByTeacherID   uuid.UUID      `db:"uploaded_by_teacher_id" gorm:"type:uuid;index;not null;constraint:materials_teacher_fkey,OnDelete:RESTRICT" validate:"required,uuid"`
	AcademicUnitID        *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid;index;constraint:materials_unit_fkey,OnDelete:SET NULL" validate:"omitempty,uuid"`
	Title                 string         `db:"title" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Description           *string        `db:"description" gorm:"default:null" validate:"omitempty"`
	Subject               *string        `db:"subject" gorm:"default:null;size:100" validate:"omitempty"`
	Grade                 *string        `db:"grade" gorm:"default:null;size:50" validate:"omitempty"`
	FileURL               string         `db:"file_url" gorm:"not null;default:''" validate:"required,url"`
	FileType              string         `db:"file_type" gorm:"not null;size:100" validate:"required,max=100"`
	FileSizeBytes         int64          `db:"file_size_bytes" gorm:"not null;default:0" validate:"required"`
	Status                string         `db:"status" gorm:"not null;type:varchar(50);check:materials_status_check,status IN ('draft','uploaded','processing','ready','failed');index:idx_materials_status" validate:"required,oneof=draft uploaded processing ready failed"`
	ProcessingStartedAt   *time.Time     `db:"processing_started_at" gorm:"default:null"`
	ProcessingCompletedAt *time.Time     `db:"processing_completed_at" gorm:"default:null"`
	IsPublic              bool           `db:"is_public" gorm:"not null;default:false"`
	CreatedAt             time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt             time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	// NOTE: partial index idx_materials_status_active (WHERE deleted_at IS NULL) must be created in post_gorm.sql
	DeletedAt             gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Material) TableName() string {
	return "content.materials"
}
