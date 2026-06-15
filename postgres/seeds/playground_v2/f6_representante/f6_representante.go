// Package f6_representante es un fixture de la línea v2 que siembra EVALUACIONES
// asignadas a la sesión del hijo del representante del mundo `base`, para validar
// E2E la pantalla "Evaluaciones del hijo" (plan 024 micro-plan / GET
// /me/wards/assessments).
//
// Contexto: en `base` el esquema de evaluación (assessment.assessment /
// assessment.assessment_assignment) queda VACÍO, así que el representante no tiene
// nada que ver en la pestaña de evaluaciones del hijo. Este fixture compone ENCIMA
// de `base`: invoca primero base.Apply (idempotente: trunca y resiembra el mundo
// de dev) y reutiliza sus escuelas, unidades, materias, membresías, sesiones e
// inscripciones tal cual; luego agrega filas a assessment.assessment y
// assessment.assessment_assignment. Por eso
// `make docker-playground-v2 P=f6_representante` deja un mundo completo y
// consistente — el migrator aplica `system` + ESTE fixture, que a su vez siembra
// `base` adentro (el runner NO encadena `base` automáticamente).
//
// Tejido representante↔hijo↔sesión en `base` que este fixture explota
// (ver seedGuardianTejido + seedSubjectOfferings de base):
//   - Laura Mendoza (tutor.mendoza@edugo.test, user …0011) ↔ hija Sofia (user
//     …0009) en S1 (Colegio San Ignacio, b1…01), relación guardian activa.
//   - Sofia (membership bb…03, unidad ac…03 = 5to A) está inscrita en la sesión
//     de Matemáticas 5to A (offering offMat5A = c5…01, subject dd…01), dictada por
//     la docente Martínez (membership bb…08).
//
// Lo que siembra (2 evaluaciones + 2 asignaciones a la sesión de Sofia):
//   - "Examen Final — Fracciones" (kind='final', status='published'): genera
//     componente de nota en el expediente (el worker ramifica por kind).
//   - "Práctica guiada — Geometría" (kind='practice', status='published'): NO va
//     al expediente; su resultado se guarda en academic.practice_result.
//   Ambas creadas por la docente de Mate5A (created_by_membership_id = bb…08),
//   subject_id = Matemáticas (dd…01), school_id = San Ignacio (b1…01).
//
// Contrato del lector (GET /me/wards/assessments): el guardián ve las
// evaluaciones de la sesión en la que su acudido está inscrito vía
// assessment.assessment_assignment (target = subject_offering_id). La entrega NO
// crea filas por alumno: los destinatarios se resuelven de
// academic.subject_offering_enrollments — por eso basta asignar a la OFERTA de
// Sofia (offMat5A) para que ella (student) y por ende Laura (guardian) las vean.
// Las dos quedan published + asignadas con due_date a futuro para que estén en
// estado consumible.
//
// Idempotente: OnConflict DoNothing por id. Rango UUID propio
// f6000000-0000-0000-0000-0000000000NN (LIBRE; no colisiona con los demás
// fixtures).
package f6_representante

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/base"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// IDs reales del mundo `base` (Colegio San Ignacio) que este fixture explota.
	schoolID    = "b1000000-0000-0000-0000-000000000001" // Colegio San Ignacio
	subjectMath = "dd000000-0000-0000-0000-000000000001" // Matematicas (escuela)
	// Sesión de Matemáticas 5to A donde Sofia está inscrita (subject_offering).
	offeringMat5A = "c5000000-0000-0000-0000-000000000001" // offMat5A
	// Membresía de la docente Martínez en Mate5A (creadora/asignadora).
	teacherMembID = "bb000000-0000-0000-0000-000000000008"

	// Evaluaciones del fixture (rango f6… propio).
	assessmentFinalID    = "f6000000-0000-0000-0000-000000000001"
	assessmentPracticeID = "f6000000-0000-0000-0000-000000000002"

	// Asignaciones de cada evaluación a la sesión de Sofia.
	assignmentFinalID    = "f6000000-0000-0000-0000-000000000011"
	assignmentPracticeID = "f6000000-0000-0000-0000-000000000012"
)

// Apply siembra las evaluaciones del hijo del representante sobre el mundo `base`.
// Primero compone `base` (trunca y resiembra; idempotente) para garantizar que
// existan la escuela/materia/sesión/inscripción/membresías a las que apuntan las
// FKs; luego inserta las filas en assessment.assessment y
// assessment.assessment_assignment. Idempotente por id.
func Apply(tx *gorm.DB) error {
	// Compone el mundo de datos `base` (trunca y resiembra; idempotente). El
	// runner aplica `system` + este fixture, pero NO encadena `base` solo.
	if err := base.Apply(tx); err != nil {
		return fmt.Errorf("playground_v2/f6_representante: base: %w", err)
	}

	schoolUUID := common.MustParseUUID(schoolID)
	subjectUUID := common.MustParseUUID(subjectMath)
	offeringUUID := common.MustParseUUID(offeringMat5A)
	teacherMembUUID := common.MustParseUUID(teacherMembID)

	// due_date a futuro (fijo, UTC ISO-Z) para que ambas queden consumibles por el
	// alumno; el read-path del guardián las muestra como pendientes de tomar.
	dueDate := time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)

	// 1) Evaluaciones (assessment.assessment). Ambas published, creadas por la
	// docente de Mate5A (bb…08), subject Matemáticas (dd…01), school San Ignacio.
	// kind explícito: 'final' (va al expediente) vs 'practice' (no va).
	assessments := []entities.Assessment{
		{
			ID:                    common.MustParseUUID(assessmentFinalID),
			SchoolID:              schoolUUID,
			CreatedByMembershipID: teacherMembUUID,
			SubjectID:             subjectUUID,
			Title:                 "Examen Final — Fracciones",
			SourceType:            "manual",
			Status:                "published",
			Kind:                  "final",
			QuestionsCount:        0,
			PassThreshold:         70,
			ShowCorrectAnswers:    true,
		},
		{
			ID:                    common.MustParseUUID(assessmentPracticeID),
			SchoolID:              schoolUUID,
			CreatedByMembershipID: teacherMembUUID,
			SubjectID:             subjectUUID,
			Title:                 "Práctica guiada — Geometría",
			SourceType:            "manual",
			Status:                "published",
			Kind:                  "practice",
			QuestionsCount:        0,
			PassThreshold:         70,
			ShowCorrectAnswers:    true,
		},
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&assessments).Error; err != nil {
		return fmt.Errorf("playground_v2/f6_representante: assessments: %w", err)
	}

	// 2) Asignaciones (assessment.assessment_assignment). Target = OFERTA de Sofia
	// (offMat5A): los destinatarios se resuelven de subject_offering_enrollments,
	// así que asignar a la sesión basta para que Sofia (y Laura) las vean.
	// assigned_by = la misma docente (bb…08). due_date a futuro.
	assignments := []entities.AssessmentAssignment{
		{
			ID:                     common.MustParseUUID(assignmentFinalID),
			AssessmentID:           common.MustParseUUID(assessmentFinalID),
			SubjectOfferingID:      offeringUUID,
			AssignedByMembershipID: teacherMembUUID,
			DueDate:                &dueDate,
		},
		{
			ID:                     common.MustParseUUID(assignmentPracticeID),
			AssessmentID:           common.MustParseUUID(assessmentPracticeID),
			SubjectOfferingID:      offeringUUID,
			AssignedByMembershipID: teacherMembUUID,
			DueDate:                &dueDate,
		},
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&assignments).Error; err != nil {
		return fmt.Errorf("playground_v2/f6_representante: assignments: %w", err)
	}

	return nil
}
