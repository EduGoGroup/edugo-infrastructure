package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserMaterialTag representa la tabla 'content.user_material_tags' en
// PostgreSQL. Es una ETIQUETA PERSONAL que un usuario aplica a un material para
// organizar su propia biblioteca; jamas se comparte entre usuarios.
//
// Sirve a las tres vistas (plan 033, D-B2.6): la biblioteca del profesor (lo que
// crea/importa), el catalogo (se etiqueta lo propio) y la biblioteca del alumno
// (lo que recibe). Habilita el filtro por etiqueta en los listados.
//
// El par (user_id, material_id) no lleva FK dura aqui: user_id vive en el schema
// auth y material_id en content; la unicidad de la etiqueta se garantiza con el
// UNIQUE (user_id, material_id, tag) y el filtrado por usuario se apoya en el
// indice idx_user_material_tags_user.
type UserMaterialTag struct {
	ID         uuid.UUID `db:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	UserID     uuid.UUID `db:"user_id" gorm:"type:uuid;not null;index:idx_user_material_tags_user;uniqueIndex:uq_user_material_tag,priority:1" validate:"required,uuid"`
	MaterialID uuid.UUID `db:"material_id" gorm:"type:uuid;not null;uniqueIndex:uq_user_material_tag,priority:2" validate:"required,uuid"`
	Tag        string    `db:"tag" gorm:"not null;type:varchar(50);uniqueIndex:uq_user_material_tag,priority:3" validate:"required,min=1,max=50"`
	CreatedAt  time.Time `db:"created_at" gorm:"not null;autoCreateTime" validate:"-"`
	UpdatedAt  time.Time `db:"updated_at" gorm:"not null;autoUpdateTime" validate:"-"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (UserMaterialTag) TableName() string {
	return "content.user_material_tags"
}
