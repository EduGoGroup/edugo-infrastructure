package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Question representa la tabla 'questions' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 054_assessment_questions.sql
// Usada por: api-mobile, worker
type Question struct {
	ID            uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey"`
	AssessmentID  uuid.UUID      `db:"assessment_id" gorm:"type:uuid;index;not null"`
	SortOrder     int            `db:"sort_order" gorm:"not null;default:0"`
	QuestionText  string         `db:"question_text" gorm:"not null"`
	QuestionType  string         `db:"question_type" gorm:"not null;type:varchar(50)"`
	CorrectAnswer *string        `db:"correct_answer" gorm:"default:null"`
	Explanation   *string        `db:"explanation" gorm:"default:null"`
	Points        float64        `db:"points" gorm:"type:numeric(5,2);not null;default:1"`
	Difficulty    *string        `db:"difficulty" gorm:"type:varchar(20)"`
	Tags          pq.StringArray `db:"tags" gorm:"type:text[]"`
	CreatedAt     time.Time      `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime"`

	Options []QuestionOption `gorm:"foreignKey:QuestionID"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Question) TableName() string {
	return "assessment.questions"
}
