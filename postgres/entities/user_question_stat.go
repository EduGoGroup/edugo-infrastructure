package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserQuestionStat representa la tabla 'assessment.user_question_stat' en
// PostgreSQL (plan 035 D-035.4). Es el ACUMULADOR ACOTADO del alumno: una fila
// por (alumno × pregunta) que se ACTUALIZA (no crece con la practica). Tamaño
// O(alumnos × preguntas), no O(prácticas). La logica adaptativa (F2) lee SOLO
// esta tabla; es barata para siempre.
//
// Lo alimentan AMBOS planos de forma SINCRONA (D-035.3/D-035.4): cada respuesta
// de practica (en su tx) y cada respuesta autocalificada del submit de examen
// (post-commit del intento; error loggeado sin revertir — el examen jamas falla
// por la traza). Es historia del alumno, no de la practica: sobrevive al
// assessment.
//
// question_id es NULLABLE con ON DELETE SET NULL (el historial NO cuelga del
// assessment; question_snapshot conserva el enunciado legible si borran la
// pregunta). next_review_at / interval_days son campos SRS SEMBRADOS: F1 no los
// lee, F2 si (repaso espaciado tipo Leitner).
//
// UNIQUE (membership_id, question_id); indices (school_id, membership_id,
// subject_id) y (membership_id, next_review_at) para F2. FKs cross-schema
// (school_id, membership_id, subject_id, question_id SET NULL) en post_gorm.sql
// (GORM no crea FK desde el tag `constraint:` sin campo de relacion).
type UserQuestionStat struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID         uuid.UUID  `db:"school_id" gorm:"type:uuid;not null;index:idx_user_question_stat_scope,priority:1" validate:"required,uuid"`
	MembershipID     uuid.UUID  `db:"membership_id" gorm:"type:uuid;not null;index:idx_user_question_stat_scope,priority:2;index:idx_user_question_stat_review,priority:1;uniqueIndex:uq_user_question_stat,priority:1" validate:"required,uuid"`
	QuestionID       *uuid.UUID `db:"question_id" gorm:"type:uuid;default:null;uniqueIndex:uq_user_question_stat,priority:2" validate:"omitempty"`
	SubjectID        uuid.UUID  `db:"subject_id" gorm:"type:uuid;not null;index:idx_user_question_stat_scope,priority:3" validate:"required,uuid"`
	QuestionSnapshot string     `db:"question_snapshot" gorm:"type:text;not null"`
	TimesSeen        int        `db:"times_seen" gorm:"not null;default:0"`
	TimesCorrect     int        `db:"times_correct" gorm:"not null;default:0"`
	StreakCorrect    int        `db:"streak_correct" gorm:"not null;default:0"`
	LastResult       *bool      `db:"last_result" gorm:"default:null"`
	LastSeenAt       time.Time  `db:"last_seen_at" gorm:"not null;default:now()"`
	NextReviewAt     *time.Time `db:"next_review_at" gorm:"default:null;index:idx_user_question_stat_review,priority:2"`
	IntervalDays     int        `db:"interval_days" gorm:"not null;default:0"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (UserQuestionStat) TableName() string {
	return "assessment.user_question_stat"
}
