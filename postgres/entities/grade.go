package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Grade representa la tabla 'grades' en PostgreSQL.
//
// Migracion: 090_academic_grades.sql
type Grade struct {
	ID               uuid.UUID       `db:"id" gorm:"type:uuid;primaryKey"`
	MembershipID     uuid.UUID       `db:"membership_id" gorm:"type:uuid;index;not null"`
	SubjectID        uuid.UUID       `db:"subject_id" gorm:"type:uuid;index;not null"`
	PeriodID         uuid.UUID       `db:"period_id" gorm:"type:uuid;index;not null"`
	GradeValue       *float64        `db:"grade_value" gorm:"type:decimal(5,2)"`
	GradeLetter      *string         `db:"grade_letter" gorm:"type:varchar(5)"`
	AssessmentScores json.RawMessage `db:"assessment_scores" gorm:"type:jsonb;default:'[]'"`
	TeacherID        *uuid.UUID      `db:"teacher_id" gorm:"type:uuid"`
	Notes            *string         `db:"notes"`
	FinalizedAt      *time.Time      `db:"finalized_at"`
	CreatedAt        time.Time       `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time       `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

func (Grade) TableName() string {
	return "academic.grades"
}
