package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ServiceClient representa la tabla 'service_clients' en PostgreSQL.
// Registra clientes M2M autorizados a llamar rutas /api/v1/internal/*.
//
// Usada por: edugo-api-identity (client credentials), edugo-api-platform (validación)
// secret_hash almacena el hash del client_secret (nunca el secret en claro).
type ServiceClient struct {
	ID          uuid.UUID      `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ClientID    string         `db:"client_id" gorm:"uniqueIndex:service_clients_client_id_key;not null;size:64" validate:"required,max=64"`
	SecretHash  string         `db:"secret_hash" gorm:"not null;type:text" validate:"required"`
	Scopes      pq.StringArray `db:"scopes" gorm:"type:text[];not null;default:'{}'" validate:"required"`
	IsActive    bool           `db:"is_active" gorm:"not null;default:true"`
	Description *string        `db:"description" gorm:"type:varchar(255)" validate:"omitempty,max=255"`
	CreatedAt   time.Time      `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time      `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	LastUsedAt  *time.Time     `db:"last_used_at" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ServiceClient) TableName() string {
	return "auth.service_clients"
}
