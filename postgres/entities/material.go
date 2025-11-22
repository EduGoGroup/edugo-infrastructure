package entities

import (
	"time"

	"github.com/google/uuid"
)

// Material representa la tabla 'materials' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 005_create_materials.up.sql
// Usada por: api-mobile, worker
type Material struct {
	ID                     uuid.UUID  `db:"id"`
	SchoolID               uuid.UUID  `db:"school_id"`
	UploadedByTeacherID    uuid.UUID  `db:"uploaded_by_teacher_id"`
	AcademicUnitID         *uuid.UUID `db:"academic_unit_id"` // NULL permitido
	Title                  string     `db:"title"`
	Description            *string    `db:"description"`
	Subject                *string    `db:"subject"`
	Grade                  *string    `db:"grade"`
	FileURL                string     `db:"file_url"`
	FileType               string     `db:"file_type"`
	FileSizeBytes          int64      `db:"file_size_bytes"`
	Status                 string     `db:"status"` // uploaded, processing, ready, failed
	ProcessingStartedAt    *time.Time `db:"processing_started_at"`
	ProcessingCompletedAt  *time.Time `db:"processing_completed_at"`
	IsPublic               bool       `db:"is_public"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
	DeletedAt              *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Material) TableName() string {
	return "materials"
}
