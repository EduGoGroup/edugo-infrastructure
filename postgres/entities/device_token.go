package entities

import (
	"time"

	"github.com/google/uuid"
)

// DeviceToken representa la tabla 'device_tokens' en PostgreSQL.
// Registra tokens push (FCM/APNs) por usuario y plataforma.
//
// Usada por: edugo-api-platform (Notification Gateway)
// revoked_at NULL = token activo; no NULL = revocado (RF-02.5).
type DeviceToken struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID      uuid.UUID `db:"user_id" gorm:"type:uuid;index;not null;constraint:device_tokens_user_fkey,OnDelete:CASCADE;uniqueIndex:uq_device_tokens_user_token" validate:"required,uuid"`
	DeviceToken string    `db:"device_token" gorm:"not null;type:text;uniqueIndex:uq_device_tokens_user_token" validate:"required,max=512"`
	Platform    string    `db:"platform" gorm:"not null;type:varchar(16);check:device_tokens_platform_check,platform IN ('android','ios')" validate:"required,oneof=android ios"`
	AppVersion  *string   `db:"app_version" gorm:"type:varchar(32)" validate:"omitempty,max=32"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
	// NOTE: partial index idx_device_tokens_user_active ON (user_id) WHERE revoked_at IS NULL must be created in post_gorm.sql
	RevokedAt *time.Time `db:"revoked_at" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (DeviceToken) TableName() string {
	return "notifications.device_tokens"
}

// IsActive retorna true si el token no ha sido revocado
func (dt *DeviceToken) IsActive() bool {
	return dt.RevokedAt == nil
}
