package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Assessment representa la tabla 'assessment' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 050_assessment_assessment.sql
// Usada por: api-mobile, worker
//
// Nota: El contenido completo de las preguntas se almacena en MongoDB (material_assessment).
// Esta tabla solo mantiene metadata y referencia al documento MongoDB.
// La relacion con materiales es N:N via assessment_materials (053).
type Assessment struct {
	ID                 uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	MongoDocumentID    string         `db:"mongo_document_id" gorm:"not null"`
	SchoolID           *uuid.UUID     `db:"school_id" gorm:"type:uuid;index"`
	CreatedByUserID    *uuid.UUID     `db:"created_by_user_id" gorm:"type:uuid"`
	QuestionsCount     int            `db:"questions_count" gorm:"not null;default:0"`
	Title              *string        `db:"title" gorm:"default:null"`
	Description        *string        `db:"description" gorm:"default:null"`
	PassThreshold      *float64       `db:"pass_threshold" gorm:"type:numeric(5,2);default:null"`
	MaxAttempts        *int           `db:"max_attempts" gorm:"default:null"`
	TimeLimitMinutes   *float64       `db:"time_limit_minutes" gorm:"type:numeric(7,2);default:null"`
	IsTimed            bool           `db:"is_timed" gorm:"not null;default:false"`
	ShuffleQuestions   bool           `db:"shuffle_questions" gorm:"not null;default:false"`
	ShowCorrectAnswers bool           `db:"show_correct_answers" gorm:"not null;default:true"`
	AvailableFrom      *time.Time     `db:"available_from" gorm:"default:null"`
	AvailableUntil     *time.Time     `db:"available_until" gorm:"default:null"`
	Status             string         `db:"status" gorm:"not null;type:varchar(50)"`
	CreatedAt          time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `db:"deleted_at" gorm:"index"`

	Materials []AssessmentMaterial `gorm:"foreignKey:AssessmentID"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment.assessment"
}
