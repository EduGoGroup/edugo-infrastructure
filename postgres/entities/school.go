package entities

import (
	"time"

	"github.com/google/uuid"
)

// School representa la tabla 'schools' en PostgreSQL.
// Esta entity es el reflejo exacto del schema de BD.
//
// Migración: 002_create_schools.up.sql
// Usada por: api-mobile, api-administracion, worker
type School struct {
	ID               uuid.UUID  `db:"id"`
	Name             string     `db:"name"`
	Code             string     `db:"code"`
	Address          *string    `db:"address"`
	City             *string    `db:"city"`
	Country          string     `db:"country"`
	Phone            *string    `db:"phone"`
	Email            *string    `db:"email"`
	Metadata         []byte     `db:"metadata"` // JSONB stored as []byte
	IsActive         bool       `db:"is_active"`
	SubscriptionTier string     `db:"subscription_tier"` // free, basic, premium, enterprise
	MaxTeachers      int        `db:"max_teachers"`
	MaxStudents      int        `db:"max_students"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
}

// TableName retorna el nombre de la tabla en PostgreSQL
func (School) TableName() string {
	return "schools"
}
