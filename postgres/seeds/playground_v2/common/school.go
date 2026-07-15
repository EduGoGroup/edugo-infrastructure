package common

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
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
//   - DefaultLandingScreenKey: "dashboard-home" por default (ADR 0024 F0);
//     pantalla de inicio predeterminada de la escuela si el rol no tiene la suya.
type SchoolSpec struct {
	ID                      uuid.UUID
	Name                    string
	Code                    string
	Country                 string     // default "Chile" si vacío
	SubscriptionTier        string     // default "basic" si vacío
	GradeProfile            string     // default "basic" si vacío ("basic"|"detailed")
	DefaultLandingScreenKey string     // default "dashboard-home" si vacío (ADR 0024 F0)
	ConceptTypeID           *uuid.UUID // nil = sin concept type
	MaxTeachers             int
	MaxStudents             int
	Metadata                json.RawMessage // default `{}` si nil
	// Settings es la configuración clave/valor opcional de la escuela (plan 039):
	// política LLM por carril, límites de import, etc. nil/vacío = la escuela cae
	// a los defaults de plataforma (resolución por env → default duro). Cada
	// clave/valor se valida contra el catálogo (entities) antes de insertarse.
	Settings map[string]string
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
	landing := spec.DefaultLandingScreenKey
	if landing == "" {
		landing = "dashboard-home"
	}
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.School{
		ID:                      spec.ID,
		Name:                    spec.Name,
		Code:                    spec.Code,
		Country:                 country,
		SubscriptionTier:        tier,
		GradeProfile:            profile,
		DefaultLandingScreenKey: &landing,
		ConceptTypeID:           spec.ConceptTypeID,
		MaxTeachers:             spec.MaxTeachers,
		MaxStudents:             spec.MaxStudents,
		IsActive:                true,
		Metadata:                metadata,
	}
}

// SeedSchool inserta la escuela aplicando los defaults del playground y, acto
// seguido, sus equivalencias tipo-de-invitación→rol por defecto (MP-08).
//
// Es el PUNTO ÚNICO donde nacen las escuelas de playground, así que aquí se
// engancha SeedDefaultSchoolInvitationRoles para que TODA escuela sembrada
// reciba sus 6 equivalencias sin duplicar el mapeo (shared over inline). Los
// invitation_types los siembra la capa L4 del system seed, que SIEMPRE corre
// antes que cualquier playground → la FK invitation_type_id existe. Idempotente:
// si el id de la escuela ya existe, no actualiza ni falla.
func SeedSchool(tx *gorm.DB, spec SchoolSpec) error {
	school := buildSchool(spec)
	if err := onConflictIgnore(tx, &school); err != nil {
		return err
	}
	if err := seedSchoolSettings(tx, spec.ID, spec.Settings); err != nil {
		return err
	}
	return l4.SeedDefaultSchoolInvitationRoles(tx, spec.ID)
}

// seedSchoolSettings valida cada setting contra el catálogo (entities) e inserta
// las filas en academic.school_settings (onConflict ignore por la PK compuesta
// (school_id, key)). Una clave/valor fuera de catálogo es un bug del fixture, no
// un error recuperable: se devuelve error para no sembrar configuración inválida.
func seedSchoolSettings(tx *gorm.DB, schoolID uuid.UUID, settings map[string]string) error {
	for key, value := range settings {
		if err := entities.ValidateSetting(key, value); err != nil {
			return fmt.Errorf("seedSchoolSettings escuela %s: %w", schoolID, err)
		}
		row := entities.SchoolSetting{SchoolID: schoolID, Key: key, Value: value}
		if err := onConflictIgnore(tx, &row); err != nil {
			return err
		}
	}
	return nil
}
