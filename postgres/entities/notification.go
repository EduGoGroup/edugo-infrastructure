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
	ID           uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID  `db:"user_id" gorm:"type:uuid;index;not null"`
	Type         string     `db:"type" gorm:"not null;type:varchar(50)"`
	Title        string     `db:"title" gorm:"not null;type:varchar(255)"`
	Body         *string    `db:"body" gorm:"default:null"`
	ResourceType *string    `db:"resource_type" gorm:"type:varchar(50)"`
	ResourceID   *uuid.UUID `db:"resource_id" gorm:"type:uuid"`
	IsRead       bool       `db:"is_read" gorm:"not null;default:false"`
	CreatedAt    time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	ReadAt       *time.Time `db:"read_at" gorm:"default:null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (Notification) TableName() string {
	return "notifications.notifications"
}
