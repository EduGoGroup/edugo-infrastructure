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
type LoginAttempt struct {
	ID          int        `db:"id" gorm:"primaryKey;autoIncrement"`
	Identifier  string     `db:"identifier" gorm:"not null"`
	AttemptType string     `db:"attempt_type" gorm:"not null;type:varchar(50)"`
	Successful  bool       `db:"successful" gorm:"not null;default:false"`
	UserAgent   *string    `db:"user_agent" gorm:"default:null"`
	IPAddress   *string    `db:"ip_address" gorm:"default:null"`
	AttemptedAt time.Time  `db:"attempted_at" gorm:"not null;autoCreateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (LoginAttempt) TableName() string {
	return "auth.login_attempts"
}
