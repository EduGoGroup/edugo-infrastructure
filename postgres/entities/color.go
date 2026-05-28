package entities

import (
	"time"

	"github.com/google/uuid"
)

// Color representa la tabla 'colors' bajo el schema academic.
// Recurso CRUD plano mínimo creado para validar que agregar un CRUD nuevo
// en EduGo NO requiere código Kotlin nuevo: la pantalla colors-list /
// colors-form se resuelve vía GenericListContract / GenericFormContract
// (Fase 3 — sdui-refactor-spec, F3-REQ-4).
//
// El CHECK constraint del campo `hex` (`^#[0-9A-Fa-f]{6}$`) y el UNIQUE
// (school_id, name) viven en migrations/sql/post_gorm.sql porque GORM no
// expresa CHECK con regex ni constraints compuestos con nombre.
type Color struct {
	ID        uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	SchoolID  uuid.UUID `db:"school_id" gorm:"type:uuid;not null;index:idx_colors_school;constraint:colors_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Name      string    `db:"name" gorm:"not null;type:varchar(120)" validate:"required,min=1,max=120"`
	Hex       string    `db:"hex" gorm:"not null;type:varchar(7)" validate:"required,len=7"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

func (Color) TableName() string {
	return "academic.colors"
}
