package entities

import (
	"time"

	"github.com/google/uuid"
)

// AIFailure representa la tabla 'assessment.ai_failure' en PostgreSQL (MVP-IA
// HU-IA.6, ADR 0050 del repo raíz: el rescate de lo fallado). Una fila por
// DESTINO del carril de IA —un intento cuya revisión falló o una pregunta cuya
// preparación falló— que se reabre si el destino vuelve a fallar: el UNIQUE
// (target_type, target_id) es la memoria que evita republicar sin fin.
//
// target_id NO lleva FK (P12): apunta a assessment.assessment_attempt o a
// assessment.question según target_type; si el destino se borra, la fila queda
// inerte. retry_count cuenta las publicaciones de la rutina de rescate, no las
// ejecuciones del consumidor. next_retry_at solo tiene valor en 'retrying';
// closed_at, al pasar a 'resolved' o 'cancelled'. last_error lo acota learning
// a 1024 bytes. Fechas en UTC.
//
// Los tres CHECK, el UNIQUE y el índice (status, next_retry_at) los materializa
// GORM desde los tags; el trigger set_updated_at va en post_gorm.sql.
type AIFailure struct {
	ID          uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	TargetType  string     `db:"target_type" gorm:"not null;type:varchar(20);uniqueIndex:uq_ai_failure_target,priority:1;check:ai_failure_target_type_check,target_type IN ('attempt','question')" validate:"required,oneof=attempt question"`
	TargetID    uuid.UUID  `db:"target_id" gorm:"type:uuid;not null;uniqueIndex:uq_ai_failure_target,priority:2" validate:"required,uuid"`
	Status      string     `db:"status" gorm:"not null;type:varchar(20);index:idx_ai_failure_due,priority:1;check:ai_failure_status_check,status IN ('retrying','failed','resolved','cancelled')" validate:"required,oneof=retrying failed resolved cancelled"`
	CauseClass  string     `db:"cause_class" gorm:"not null;type:varchar(20);check:ai_failure_cause_class_check,cause_class IN ('service','data')" validate:"required,oneof=service data"`
	LastError   *string    `db:"last_error" gorm:"type:text;default:null" validate:"omitempty"`
	RetryCount  int        `db:"retry_count" gorm:"type:integer;not null;default:0" validate:"gte=0"`
	NextRetryAt *time.Time `db:"next_retry_at" gorm:"index:idx_ai_failure_due,priority:2;default:null"`
	ClosedAt    *time.Time `db:"closed_at" gorm:"default:null"`
	CreatedAt   time.Time  `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (AIFailure) TableName() string {
	return "assessment.ai_failure"
}
