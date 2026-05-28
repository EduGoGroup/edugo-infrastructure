package common

import (
	"encoding/json"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchoolSpec describe la escuela a sembrar. Los campos opcionales con cero
// valor adoptan los defaults que usan los playgrounds existentes (Country
// "Chile", SubscriptionTier "basic", Metadata "{}").
type SchoolSpec struct {
	ID               uuid.UUID
	Name             string
	Code             string
	Country          string // default "Chile" si vacío
	SubscriptionTier string // default "basic" si vacío
	MaxTeachers      int
	MaxStudents      int
	Metadata         json.RawMessage // default `{}` si nil
}

// buildSchool mapea SchoolSpec a entities.School aplicando defaults.
// Extraído como función no exportada para poder testear los defaults
// sin necesidad de schema real (los enum/jsonb de Postgres no migran a
// SQLite).
func buildSchool(spec SchoolSpec) entities.School {
	country := spec.Country
	if country == "" {
		country = "Chile"
	}
	tier := spec.SubscriptionTier
	if tier == "" {
		tier = "basic"
	}
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.School{
		ID:               spec.ID,
		Name:             spec.Name,
		Code:             spec.Code,
		Country:          country,
		SubscriptionTier: tier,
		MaxTeachers:      spec.MaxTeachers,
		MaxStudents:      spec.MaxStudents,
		IsActive:         true,
		Metadata:         metadata,
	}
}

// SeedSchool inserta la escuela aplicando los defaults del playground.
// Idempotente: si el ID ya existe, no actualiza ni falla.
func SeedSchool(tx *gorm.DB, spec SchoolSpec) error {
	school := buildSchool(spec)
	return OnConflictIgnore(tx, &school)
}
