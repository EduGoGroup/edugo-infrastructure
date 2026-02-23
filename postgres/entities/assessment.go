package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Assessment representa la tabla 'assessment' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 006_create_assessments.up.sql, 009_extend_assessment_schema.up.sql
// Usada por: api-mobile, worker
//
// Nota: El contenido completo de las preguntas se almacena en MongoDB (material_assessment).
// Esta tabla solo mantiene metadata y referencia al documento MongoDB.
type Assessment struct {
	ID               uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	MaterialID       uuid.UUID      `db:"material_id" gorm:"type:uuid;index;not null"`
	MongoDocumentID  string         `db:"mongo_document_id" gorm:"not null"`
	QuestionsCount   int            `db:"questions_count" gorm:"not null;default:0"`
	Title            *string        `db:"title" gorm:"default:null"`
	PassThreshold    *int           `db:"pass_threshold" gorm:"default:null"`
	MaxAttempts      *int           `db:"max_attempts" gorm:"default:null"`
	TimeLimitMinutes *int           `db:"time_limit_minutes" gorm:"default:null"`
	Status           string         `db:"status" gorm:"not null;type:varchar(50)"`
	CreatedAt        time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `db:"deleted_at" gorm:"index"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment.assessment"
}
