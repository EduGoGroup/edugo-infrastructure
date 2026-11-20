package demo

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConceptTypePrimarySchoolID es el tipo de concepto «primary_school» que siembra
// system/ (l4/concept_types.go). Los c1000000-… de `base` no existen en system/.
var ConceptTypePrimarySchoolID = uuid.MustParse("c4000000-0000-0000-0000-000000000001")

// Ventana del período activo, en días respecto a hoy (propuesta 8 de la ficha).
const (
	periodStartDays = -60
	periodEndDays   = 120
)

// buildSchool arma la escuela, sus dos unidades, el período activo (atado a la
// raíz, como en `base`), las tres materias y las tres sesiones de la docente.
func buildSchool(w *world) {
	w.school = common.SchoolSpec{
		ID:            SchoolID,
		Name:          "Colegio Los Aromos",
		Code:          "DEMO-AROMOS",
		GradeProfile:  "basic",
		ConceptTypeID: &ConceptTypePrimarySchoolID,
		MaxTeachers:   10,
		MaxStudents:   100,
		Metadata:      json.RawMessage(`{"demo":true}`),
	}

	year := w.today.Year()
	w.units = []entities.AcademicUnit{
		{
			ID: RootUnitID, SchoolID: SchoolID, Name: "Colegio Los Aromos", Code: "AROMOS-ROOT",
			Type: "school", Description: ptr("Unidad raíz del Colegio Los Aromos"), Level: ptr("primary"),
			AcademicYear: 0, Metadata: json.RawMessage(`{"is_root":true}`), IsActive: true,
		},
		{
			ID: UnitID, ParentUnitID: ptr(RootUnitID), SchoolID: SchoolID, Name: "8.º Básico A", Code: "8A",
			Type: "class", Description: ptr("Octavo básico, sección A"), Level: ptr("primary"),
			AcademicYear: year, Metadata: json.RawMessage(`{"section":"A","grade_number":8}`), IsActive: true,
		},
	}

	w.period = common.PeriodSpec{
		ID:             PeriodID,
		SchoolID:       SchoolID,
		AcademicUnitID: RootUnitID,
		Name:           "Semestre en curso",
		Code:           "DEMO-SEM",
		Type:           "semester",
		StartDate:      w.day(periodStartDays),
		EndDate:        w.day(periodEndDays),
		AcademicYear:   year,
		SortOrder:      1,
	}

	w.subjects = []common.SubjectSpec{
		{ID: SubjectMathID, SchoolID: SchoolID, Name: "Matemáticas", Code: "MAT-8"},
		{ID: SubjectLanguageID, SchoolID: SchoolID, Name: "Lenguaje", Code: "LEN-8"},
		{ID: SubjectScienceID, SchoolID: SchoolID, Name: "Ciencias Naturales", Code: "CN-8"},
	}

	teacher := TeacherMembershipID
	w.offerings = []common.OfferingSpec{
		{ID: OfferingMathID, SchoolID: SchoolID, SubjectID: SubjectMathID, AcademicUnitID: UnitID, PeriodID: PeriodID, TeacherMembershipID: &teacher},
		{ID: OfferingLanguageID, SchoolID: SchoolID, SubjectID: SubjectLanguageID, AcademicUnitID: UnitID, PeriodID: PeriodID, TeacherMembershipID: &teacher},
		{ID: OfferingScienceID, SchoolID: SchoolID, SubjectID: SubjectScienceID, AcademicUnitID: UnitID, PeriodID: PeriodID, TeacherMembershipID: &teacher},
	}
}

// writeSchool escribe la escuela (con sus equivalencias de invitación, que
// common.SeedSchool siembra en la misma llamada), sus conceptos, las unidades,
// el período y las materias.
func writeSchool(tx *gorm.DB, w *world) error {
	if err := common.SeedSchool(tx, w.school); err != nil {
		return fmt.Errorf("escuela: %w", err)
	}
	if err := writeSchoolConcepts(tx); err != nil {
		return err
	}
	if err := insert(tx, &w.units); err != nil {
		return fmt.Errorf("unidades: %w", err)
	}
	if err := common.SeedActivePeriod(tx, w.period); err != nil {
		return fmt.Errorf("período: %w", err)
	}
	for _, s := range w.subjects {
		if err := common.SeedSubject(tx, s); err != nil {
			return fmt.Errorf("materia %s: %w", s.Name, err)
		}
	}
	return nil
}

// writeSchoolConcepts copia a la escuela las definiciones del tipo de concepto,
// como hace `base` (el bundle de sincronización las sirve).
func writeSchoolConcepts(tx *gorm.DB) error {
	var defs []entities.ConceptDefinition
	if err := tx.Where("concept_type_id = ?", ConceptTypePrimarySchoolID).
		Order("sort_order, term_key").Find(&defs).Error; err != nil {
		return fmt.Errorf("leyendo concept_definitions: %w", err)
	}
	rows := make([]entities.SchoolConcept, 0, len(defs))
	for i, d := range defs {
		rows = append(rows, entities.SchoolConcept{
			ID:        newID(domSchoolConcept, i+1),
			SchoolID:  SchoolID,
			TermKey:   d.TermKey,
			TermValue: d.TermValue,
			Category:  d.Category,
			SortOrder: d.SortOrder,
		})
	}
	if err := insert(tx, &rows); err != nil {
		return fmt.Errorf("school_concepts: %w", err)
	}
	return nil
}

// writeOfferings escribe las sesiones; va después de las membresías (FK a la
// membresía docente).
func writeOfferings(tx *gorm.DB, w *world) error {
	for _, o := range w.offerings {
		if err := common.SeedOffering(tx, o); err != nil {
			return fmt.Errorf("sesión %s: %w", o.ID, err)
		}
	}
	return nil
}
