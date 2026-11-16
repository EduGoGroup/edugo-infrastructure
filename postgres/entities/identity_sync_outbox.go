package entities

import (
	"time"

	"github.com/google/uuid"
)

// IdentitySyncOutbox representa la tabla 'auth.identity_sync_outbox' en
// PostgreSQL (plan 053 T3.7, ADR 0040). Es el outbox de la sincronización de
// accesos con identity-core: la transacción de negocio solo apunta aquí que los
// accesos de un usuario cambiaron, y un proceso aparte hace la llamada externa
// fuera de la transacción.
//
// Hay una fila por usuario (user_id es la PK): un apunte nuevo sobre un usuario
// ya pendiente no crea otra fila, sube `generation`. Quien procesa la fila solo
// la borra si `generation` sigue siendo la que leyó; si subió entre medias, la
// fila queda para la siguiente vuelta.
//
// user_id NO lleva FK a auth.users a propósito (contrato de T3.7): la fila
// solo dice «propagar los accesos de este usuario» y no debe atar ni frenar
// el ciclo de vida de auth.users. Los DEFAULT viven en la BD (no
// autoCreateTime) porque los consumidores escriben con SQL crudo. El trigger
// set_updated_at y el COMMENT ON TABLE están en post_gorm.sql.
type IdentitySyncOutbox struct {
	UserID        uuid.UUID  `db:"user_id" gorm:"type:uuid;primaryKey" validate:"required,uuid"`
	Generation    int64      `db:"generation" gorm:"type:bigint;not null;default:1"`
	Attempts      int32      `db:"attempts" gorm:"type:integer;not null;default:0"`
	NextAttemptAt time.Time  `db:"next_attempt_at" gorm:"type:timestamptz;not null;default:now();index:idx_identity_sync_outbox_next_attempt"`
	LockedUntil   *time.Time `db:"locked_until" gorm:"type:timestamptz;default:null"`
	LastError     *string    `db:"last_error" gorm:"type:text;default:null" validate:"omitempty"`
	CreatedAt     time.Time  `db:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time  `db:"updated_at" gorm:"type:timestamptz;not null;default:now()"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (IdentitySyncOutbox) TableName() string {
	return "auth.identity_sync_outbox"
}
