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
	ID                    uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	SchoolID              uuid.UUID      `db:"school_id" gorm:"type:uuid;index;not null"`
	UploadedByTeacherID   uuid.UUID      `db:"uploaded_by_teacher_id" gorm:"type:uuid;index;not null"`
	AcademicUnitID        *uuid.UUID     `db:"academic_unit_id" gorm:"type:uuid;index"`
	Title                 string         `db:"title" gorm:"not null"`
	Description           *string        `db:"description" gorm:"default:null"`
	Subject               *string        `db:"subject" gorm:"default:null"`
	Grade                 *string        `db:"grade" gorm:"default:null"`
	FileURL               string         `db:"file_url" gorm:"not null"`
	FileType              string         `db:"file_type" gorm:"not null"`
	FileSizeBytes         int64          `db:"file_size_bytes" gorm:"not null;default:0"`
	Status                string         `db:"status" gorm:"not null;type:varchar(50)"`
	ProcessingStartedAt   *time.Time     `db:"processing_started_at" gorm:"default:null"`
	ProcessingCompletedAt *time.Time     `db:"processing_completed_at" gorm:"default:null"`
	IsPublic              bool           `db:"is_public" gorm:"not null;default:false"`
	CreatedAt             time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt             time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt             gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Material) TableName() string {
	return "content.materials"
}
