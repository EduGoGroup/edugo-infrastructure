package entities

import (
	"time"

	"github.com/google/uuid"
)

// Notification representa la tabla 'notifications' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migraciones: 095_notifications.sql
// Usada por: api-mobile
type Notification struct {
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID       uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null;constraint:notifications_user_fkey,OnDelete:CASCADE;index:idx_notif_user_all" validate:"required,uuid"`
	Type         string     `db:"type" gorm:"not null;type:varchar(50)" validate:"required"`
	Title        string     `db:"title" gorm:"not null;type:varchar(255)" validate:"required,min=2,max=255"`
	Body         *string    `db:"body" gorm:"default:null" validate:"omitempty"`
	ResourceType *string    `db:"resource_type" gorm:"type:varchar(50)" validate:"omitempty"`
	ResourceID   *uuid.UUID `db:"resource_id" gorm:"type:uuid" validate:"omitempty,uuid"`
	// NOTE: partial index idx_notif_user_unread ON (user_id, created_at DESC) WHERE is_read = FALSE must be created in post_gorm.sql
	IsRead       bool       `db:"is_read" gorm:"not null;default:false"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime;index:idx_notif_user_all" validate:"-"`
	ReadAt       *time.Time `db:"read_at" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Notification) TableName() string {
	return "notifications.notifications"
}
