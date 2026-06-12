package entities

import (
	"time"

	"github.com/google/uuid"
)

// System representa la tabla 'systems' en el schema iam.
//
// Catalogo de sistemas/apps del ecosistema (p.ej. 'kmp', 'admin-tool'). Modela
// en datos el acceso por sistema (MP-08): que roles entran a cada app se
// resuelve por la tabla puente iam.system_roles, no por nombres hardcodeados.
// Los valores se siembran en F1 (este F0 solo define el esquema).
//
// El indice UNIQUE de `key` lo materializa GORM desde el tag uniqueIndex; el
// trigger set_updated_at vive en sql/post_gorm.sql (seccion iam).
type System struct {
	ID          uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	Key         string    `db:"key" gorm:"uniqueIndex:systems_key_key;not null;type:varchar(50)" validate:"required,min=1,max=50"`
	Name        string    `db:"name" gorm:"not null;type:varchar(100)" validate:"required,min=1,max=100"`
	Description *string   `db:"description" gorm:"type:varchar(255);default:null" validate:"omitempty,max=255"`
	CreatedAt   time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (System) TableName() string {
	return "iam.systems"
}
