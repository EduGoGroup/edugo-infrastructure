package common

import (
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnitSpec describe la unidad académica a sembrar. Type por default "class"
// (n0n1_escuelas usa "grade"). AcademicYear se aplica tal cual (no hay default
// implícito: el playground siempre lo conoce).
type UnitSpec struct {
	ID           uuid.UUID
	SchoolID     uuid.UUID
	Name         string
	Code         string
	Type         string // default "class" si vacío ("school"|"grade"|"class"|"section"|"club"|"department")
	AcademicYear int
	Metadata     json.RawMessage // default `{}` si nil
}

func buildUnit(spec UnitSpec) entities.AcademicUnit {
	typ := spec.Type
	if typ == "" {
		typ = "class"
	}
	metadata := spec.Metadata
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	return entities.AcademicUnit{
		ID:           spec.ID,
		SchoolID:     spec.SchoolID,
		Name:         spec.Name,
		Code:         spec.Code,
		Type:         typ,
		AcademicYear: spec.AcademicYear,
		Metadata:     metadata,
		IsActive:     true,
	}
}

// SeedAcademicUnit inserta la unidad académica aplicando defaults. Idempotente
// por id.
func SeedAcademicUnit(tx *gorm.DB, spec UnitSpec) error {
	u := buildUnit(spec)
	return onConflictIgnore(tx, &u)
}

// PeriodSpec describe un período académico. Por default queda ACTIVO
// (IsActive=true) y de tipo "semester"; hay un índice único parcial por
// (school_id, academic_unit_id) WHERE is_active, así que sólo puede haber uno
// activo por (escuela, unidad).
//
//   - AcademicUnitID: el zero value uuid.Nil es válido (la columna es not-null
//     pero acepta el UUID todo-ceros). n0n1_escuelas lo deja en cero (período
//     por escuela sin anclar a unidad); los demás lo setean a la unidad real.
//   - Start/End: si ambos son zero time, se aplican defaults 2026-03-01 →
//     2026-07-31 (Semestre 1 2026), igual que los playgrounds existentes.
type PeriodSpec struct {
	ID             uuid.UUID
	SchoolID       uuid.UUID
	AcademicUnitID uuid.UUID // uuid.Nil válido (período sin anclar a unidad)
	Name           string    // default "Semestre 1 2026" si vacío
	Code           string    // si vacío, Code queda nil (puntero)
	Type           string    // default "semester" si vacío
	StartDate      time.Time // default 2026-03-01 si zero y EndDate también zero
	EndDate        time.Time // default 2026-07-31 si zero y StartDate también zero
	AcademicYear   int
	SortOrder      int
}

func buildPeriod(spec PeriodSpec) entities.AcademicPeriod {
	name := spec.Name
	if name == "" {
		name = "Semestre 1 2026"
	}
	typ := spec.Type
	if typ == "" {
		typ = "semester"
	}
	start := spec.StartDate
	end := spec.EndDate
	if start.IsZero() && end.IsZero() {
		start = time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	}
	var code *string
	if spec.Code != "" {
		c := spec.Code
		code = &c
	}
	return entities.AcademicPeriod{
		ID:             spec.ID,
		SchoolID:       spec.SchoolID,
		AcademicUnitID: spec.AcademicUnitID,
		Name:           name,
		Code:           code,
		Type:           typ,
		StartDate:      start,
		EndDate:        end,
		IsActive:       true,
		AcademicYear:   spec.AcademicYear,
		SortOrder:      spec.SortOrder,
	}
}

// SeedActivePeriod inserta un período académico ACTIVO aplicando defaults.
// Idempotente por id.
func SeedActivePeriod(tx *gorm.DB, spec PeriodSpec) error {
	p := buildPeriod(spec)
	return onConflictIgnore(tx, &p)
}

// SubjectSpec describe una materia. Por convención v2 (ADR 0016) la materia es
// catálogo de ESCUELA: AcademicUnitID queda nil. Code: si vacío, queda nil.
type SubjectSpec struct {
	ID       uuid.UUID
	SchoolID uuid.UUID
	Name     string
	Code     string // si vacío, Code queda nil (puntero)
}

func buildSubject(spec SubjectSpec) entities.Subject {
	var code *string
	if spec.Code != "" {
		c := spec.Code
		code = &c
	}
	return entities.Subject{
		ID:             spec.ID,
		SchoolID:       spec.SchoolID,
		AcademicUnitID: nil,
		Name:           spec.Name,
		Code:           code,
		IsActive:       true,
	}
}

// SeedSubject inserta la materia (scope escuela, AcademicUnitID nil).
// Idempotente por id.
func SeedSubject(tx *gorm.DB, spec SubjectSpec) error {
	subj := buildSubject(spec)
	return onConflictIgnore(tx, &subj)
}
