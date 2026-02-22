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
	ID          int        `db:"id"`           // Serial autoincremental
	Identifier  string     `db:"identifier"`   // Email u otro identificador usado en el intento
	AttemptType string     `db:"attempt_type"` // Tipo de autenticación (password, oauth, etc.)
	Successful  bool       `db:"successful"`   // true si el intento fue exitoso
	UserAgent   *string    `db:"user_agent"`   // User-Agent del cliente (nullable)
	IPAddress   *string    `db:"ip_address"`   // Dirección IP del cliente IPv4/IPv6 (nullable)
	AttemptedAt time.Time  `db:"attempted_at"` // Timestamp del intento
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (LoginAttempt) TableName() string {
	return "login_attempts"
}
