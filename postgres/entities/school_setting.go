package entities

import (
	"time"

	"github.com/google/uuid"
)

// SchoolSetting representa la tabla 'school_settings' en PostgreSQL (schema
// academic). Es la configuración clave/valor POR ESCUELA (plan 039, D-039.1): la
// política LLM por carril y los límites de import viven aquí como FILAS, no como
// columnas de academic.schools. Un carril o límite nuevo es un INSERT de una
// clave nueva, sin migración.
//
// PK compuesta (school_id, key): una fila por (escuela, clave). El catálogo de
// claves válidas y sus valores permitidos vive en código
// (school_setting_catalog.go), única puerta de escritura: la tabla clave/valor
// NO lleva CHECK por clave. Value es texto libre validado contra el catálogo por
// el admin y por el endpoint M2M de academic.
//
// La FK school_id → academic.schools (ON DELETE CASCADE) se materializa en
// migrations/sql/post_gorm.sql: GORM no la crea desde el tag `constraint:` sin un
// campo de relación (mismo patrón que academic.school_invitation_roles).
type SchoolSetting struct {
	SchoolID  uuid.UUID `db:"school_id" gorm:"type:uuid;primaryKey;constraint:school_settings_school_fkey,OnDelete:CASCADE" validate:"required,uuid"`
	Key       string    `db:"key" gorm:"primaryKey;not null;type:varchar(64)" validate:"required,min=1,max=64"`
	Value     string    `db:"value" gorm:"not null;type:varchar(255)" validate:"required"`
	CreatedAt time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (SchoolSetting) TableName() string {
	return "academic.school_settings"
}
