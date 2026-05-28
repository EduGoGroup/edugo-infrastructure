package entities

import (
	"time"

	"github.com/google/uuid"
)

// NOTE: FK fk_resource_screens_screen_key (screen_key→ui_config.screen_instances.screen_key) must be created in post_gorm.sql
type ResourceScreen struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ResourceID  uuid.UUID `db:"resource_id" gorm:"type:uuid;index;not null;constraint:fk_resource_screens_resource;uniqueIndex:resource_screens_resource_id_screen_type_key;index:idx_resource_screens_resource" validate:"required,uuid"`
	ResourceKey string    `db:"resource_key" gorm:"not null;size:100;index:idx_resource_screens_resource_key" validate:"required,max=100"`
	ScreenKey   string    `db:"screen_key" gorm:"not null;size:100;index:idx_resource_screens_screen_key" validate:"required,max=100"`
	ScreenType  string    `db:"screen_type" gorm:"not null;type:varchar(50);uniqueIndex:resource_screens_resource_id_screen_type_key" validate:"required"`
	IsDefault   bool      `db:"is_default" gorm:"not null;default:false"`
	SortOrder   int       `db:"sort_order" gorm:"not null;default:0" validate:"required"`
	IsActive    bool      `db:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (ResourceScreen) TableName() string {
	return "ui_config.resource_screens"
}
