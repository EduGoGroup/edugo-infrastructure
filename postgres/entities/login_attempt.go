package entities

import (
	"time"
)

// LoginAttempt representa la tabla 'login_attempts' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Usada por: api-mobile, api-administracion
//
// Nota: identifier puede ser email, username u otro identificador usado en el intento.
// attempt_type distingue el tipo de autenticación (password, oauth, magic_link, etc.).
// ip_address soporta IPv4 (hasta 15 chars) e IPv6 (hasta 45 chars).
// NOTE: partial index idx_login_attempts_rate_limit (WHERE successful = false) must be created in post_gorm.sql
type LoginAttempt struct {
	ID          int       `db:"id" gorm:"primaryKey;autoIncrement" validate:"required"`
	Identifier  string    `db:"identifier" gorm:"not null;size:255;index:idx_login_attempts_identifier;index:idx_login_attempts_identifier_attempted_at" validate:"required"`
	AttemptType string    `db:"attempt_type" gorm:"not null;type:varchar(50);check:chk_attempt_type,attempt_type IN ('email','ip')" validate:"required,oneof=email ip"`
	Successful  bool      `db:"successful" gorm:"not null;default:false;index:idx_login_attempts_successful"`
	UserAgent   *string   `db:"user_agent" gorm:"default:null" validate:"omitempty"`
	IPAddress   *string   `db:"ip_address" gorm:"default:null;size:45" validate:"omitempty"`
	AttemptedAt time.Time `db:"attempted_at" gorm:"not null;autoCreateTime;index:idx_login_attempts_attempted_at;index:idx_login_attempts_identifier_attempted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (LoginAttempt) TableName() string {
	return "auth.login_attempts"
}
