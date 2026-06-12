package common

import (
	"encoding/json"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchoolSpec describe la escuela a sembrar. Los campos opcionales con cero valor
// adoptan los defaults que usan los playgrounds v2 (Country "Chile",
// SubscriptionTier "basic", GradeProfile "basic", Metadata "{}").
//
//   - GradeProfile: "basic" por default; n4_evaluacion lo pone "detailed" para
//     habilitar el desglose por componentes (grade_item) del cierre N4 (ADR 0020).
//   - ConceptTypeID: nil por default; n0n1_escuelas lo setea (concept type
//     "university") para escuelas con concepto de universidad.
//   - Country: "Chile" por default; n0n1_escuelas usa "Venezuela".
type SchoolSpec struct {
	ID               uuid.UUID
	Name             string
	Code             string
	Country          string     // default "Chile" si vacío
	SubscriptionTier string     // default "basic" si vacío
	GradeProfile     string     // default "basic" si vacío ("basic"|"detailed")
	ConceptTypeID    *uuid.UUID // nil = sin concept type
	MaxTeachers      int
	MaxStudents      int
	Metadata         json.RawMessage // default `{}` si nil
}

// buildSchool mapea SchoolSpec a entities.School aplicando defaults.
func buildSchool(spec SchoolSpec) entities.School {
	country := spec.Country
	if country == "" {
		country = "Chile"
	}
	tier := spec.SubscriptionTier
	if tier == "" {
		tier = "basic"
	}
	profile := spec.GradeProfile
	if profile == "" {
		profile = "basic"
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
		GradeProfile:     profile,
		ConceptTypeID:    spec.ConceptTypeID,
		MaxTeachers:      spec.MaxTeachers,
		MaxStudents:      spec.MaxStudents,
		IsActive:         true,
		Metadata:         metadata,
	}
}

// SeedSchool inserta la escuela aplicando los defaults del playground.
// Idempotente: si el id ya existe, no actualiza ni falla.
func SeedSchool(tx *gorm.DB, spec SchoolSpec) error {
	school := buildSchool(spec)
	return onConflictIgnore(tx, &school)
}
