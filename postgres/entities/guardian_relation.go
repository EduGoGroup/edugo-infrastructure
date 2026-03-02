package entities

import (
	"time"

	"github.com/google/uuid"
)

// GuardianRelation representa la tabla 'guardian_relations' en PostgreSQL
// Esta entity es el reflejo exacto del schema de BD definido en:
// - postgres/migrations/structure/034_academic_guardian_relations.sql
//
// Representa la relación entre un apoderado (guardian) y un estudiante.
// Define el tipo de relación familiar o legal entre ellos.
type GuardianRelation struct {
	ID               uuid.UUID  `db:"id" gorm:"type:uuid;primaryKey"`
	GuardianID       uuid.UUID  `db:"guardian_id" gorm:"type:uuid;index;not null"`
	StudentID        uuid.UUID  `db:"student_id" gorm:"type:uuid;index;not null"`
	RelationshipType string     `db:"relationship_type" gorm:"not null;type:varchar(50)"`
	IsPrimary        bool       `db:"is_primary" gorm:"not null;default:false"`
	IsActive         bool       `db:"is_active" gorm:"not null;default:true"`
	Status           string     `db:"status" gorm:"not null;type:varchar(20);default:'active'"`
	CreatedBy        *uuid.UUID `db:"created_by" gorm:"type:uuid"`
	CreatedAt        time.Time  `db:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt        time.Time  `db:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (GuardianRelation) TableName() string {
	return "academic.guardian_relations"
}
