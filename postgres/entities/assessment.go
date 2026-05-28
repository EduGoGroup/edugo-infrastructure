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
	ID                 uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	MongoDocumentID    *string        `db:"mongo_document_id" gorm:"default:null;size:24;uniqueIndex:assessment_mongo_unique" validate:"omitempty"`
	SourceType         string         `db:"source_type" gorm:"not null;type:varchar(20);default:'manual';check:assessment_source_type_check,source_type IN ('manual','ai_generated')" validate:"required,oneof=manual ai_generated"`
	SchoolID           *uuid.UUID     `db:"school_id" gorm:"type:uuid;index" validate:"omitempty,uuid"`
	CreatedByUserID    *uuid.UUID     `db:"created_by_user_id" gorm:"type:uuid" validate:"omitempty,uuid"`
	QuestionsCount     int            `db:"questions_count" gorm:"not null;default:0" validate:"required"`
	Title              *string        `db:"title" gorm:"default:null;size:255" validate:"omitempty"`
	Description        *string        `db:"description" gorm:"default:null" validate:"omitempty"`
	PassThreshold      *float64       `db:"pass_threshold" gorm:"type:numeric(5,2);default:70;check:assessment_pass_threshold_check,pass_threshold >= 0 AND pass_threshold <= 100" validate:"omitempty"`
	MaxAttempts        *int           `db:"max_attempts" gorm:"default:null" validate:"omitempty"`
	TimeLimitMinutes   *float64       `db:"time_limit_minutes" gorm:"type:numeric(7,2);default:null" validate:"omitempty"`
	IsTimed            bool           `db:"is_timed" gorm:"not null;default:false"`
	ShuffleQuestions   bool           `db:"shuffle_questions" gorm:"not null;default:false"`
	ShowCorrectAnswers bool           `db:"show_correct_answers" gorm:"not null;default:true"`
	AvailableFrom      *time.Time     `db:"available_from" gorm:"default:null"`
	AvailableUntil     *time.Time     `db:"available_until" gorm:"default:null;check:assessment_available_dates_check,available_until IS NULL OR available_from IS NULL OR available_until > available_from"`
	Status             string         `db:"status" gorm:"not null;type:varchar(50);check:assessment_status_check,status IN ('draft','generated','published','archived','closed');default:'generated'" validate:"required,oneof=draft generated published archived closed"`
	CreatedAt          time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt          time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	DeletedAt          gorm.DeletedAt `db:"deleted_at" gorm:"index" validate:"-"`

	Materials []AssessmentMaterial `gorm:"foreignKey:AssessmentID" validate:"-"`
	Questions []Question           `gorm:"foreignKey:AssessmentID" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment.assessment"
}
