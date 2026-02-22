package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RefreshToken representa la tabla 'refresh_tokens' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Usada por: api-mobile, api-administracion
//
// Nota: token_hash almacena el hash del refresh token (nunca el token en claro).
// client_info contiene metadatos del cliente (user-agent, IP, device, etc.).
// revoked_at es NULL si el token sigue vigente.
// replaced_by apunta al nuevo token cuando éste es rotado (refresh token rotation).
type RefreshToken struct {
	ID          uuid.UUID       `db:"id"`
	TokenHash   string          `db:"token_hash"`   // Hash del refresh token (nunca texto claro)
	UserID      uuid.UUID       `db:"user_id"`      // Usuario propietario del token
	ClientInfo  json.RawMessage `db:"client_info"`  // Metadata del cliente (JSONB, nullable)
	ExpiresAt   time.Time       `db:"expires_at"`   // Fecha de expiración del token
	CreatedAt   time.Time       `db:"created_at"`   // Fecha de creación
	RevokedAt   *time.Time      `db:"revoked_at"`   // NULL si el token sigue vigente
	ReplacedBy  *uuid.UUID      `db:"replaced_by"`  // UUID del nuevo token si fue rotado (nullable)
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsValid retorna true si el token no ha expirado ni ha sido revocado
func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}
