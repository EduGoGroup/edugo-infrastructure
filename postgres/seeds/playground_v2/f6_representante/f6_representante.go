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
//   - "Práctica guiada — Geometría" (purpose='practice', status='published'): NO va
//     al expediente (la trazabilidad de práctica vive en el plano
//     assessment.practice_session, plan 035 F1a; academic.practice_result se
//     eliminó en el plan 037 F1g).
//     Ambas creadas por la docente de Mate5A (created_by_membership_id = bb…08),
//     subject_id = Matemáticas (dd…01), school_id = San Ignacio (b1…01).
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
	"github.com/google/uuid"
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

	// Preguntas de la evaluación final (Fracciones), sort_order 1..3.
	qFinal1ID = "f6000000-0000-0000-0000-000000000101" // multiple_choice
	qFinal2ID = "f6000000-0000-0000-0000-000000000102" // true_false
	qFinal3ID = "f6000000-0000-0000-0000-000000000103" // multiple_choice

	// Preguntas de la práctica (Geometría), sort_order 1..3.
	qPractice1ID = "f6000000-0000-0000-0000-000000000201" // multiple_choice
	qPractice2ID = "f6000000-0000-0000-0000-000000000202" // true_false
	qPractice3ID = "f6000000-0000-0000-0000-000000000203" // multiple_choice
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
	// purpose explícito (plan 035 D-035.1): 'exam' (va al expediente) vs
	// 'practice' (no va).
	assessments := []entities.Assessment{
		{
			ID:                    common.MustParseUUID(assessmentFinalID),
			SchoolID:              schoolUUID,
			CreatedByMembershipID: teacherMembUUID,
			SubjectID:             subjectUUID,
			Title:                 "Examen Final — Fracciones",
			SourceType:            "manual",
			Status:                "published",
			Purpose:               "exam",
			QuestionsCount:        3,
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
			Purpose:               "practice",
			QuestionsCount:        3,
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

	// 1.b) Preguntas (assessment.question + assessment.question_option). 3 por
	// evaluación, tipos sencillos (multiple_choice / true_false) para que la
	// pantalla de toma del alumno tenga contenido. Temática coherente con cada
	// evaluación (Fracciones / Geometría). Idempotente por id.
	if err := seedFinalQuestions(tx); err != nil {
		return fmt.Errorf("playground_v2/f6_representante: questions_final: %w", err)
	}
	if err := seedPracticeQuestions(tx); err != nil {
		return fmt.Errorf("playground_v2/f6_representante: questions_practice: %w", err)
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

// strPtr es un helper para el campo *string correct_answer del entity Question.
func strPtr(s string) *string { return &s }

// optionUUID deriva el id de una opción a partir del question_id + sort_order,
// de forma determinística e idempotente (no enumera a mano los ids de opciones,
// pero se mantienen dentro del rango lógico del fixture sin colisionar con los
// f6… explícitos de preguntas/evaluaciones/asignaciones).
func optionUUID(questionID uuid.UUID, sortOrder int) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, fmt.Appendf(nil, "f6_representante:opt:%s:%d", questionID, sortOrder))
}

// seedFinalQuestions siembra las 3 preguntas de "Examen Final — Fracciones".
// Tipos: multiple_choice (con opciones), true_false (con opciones Verdadero/
// Falso). correct_answer replica el formato del form de autoría / la pantalla de
// toma (texto de la opción correcta para choice; "true"/"false" para true_false).
func seedFinalQuestions(tx *gorm.DB) error {
	assessmentID := common.MustParseUUID(assessmentFinalID)

	// Q1 — multiple_choice
	q1 := common.MustParseUUID(qFinal1ID)
	if err := upsertQuestion(tx, q1, assessmentID, 1, "multiple_choice",
		"¿Cuánto es 1/2 + 1/4?", strPtr("3/4"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"3/4", "2/6", "1/6", "1/2"} {
		if err := upsertQuestionOption(tx, q1, i, opt); err != nil {
			return err
		}
	}

	// Q2 — true_false
	q2 := common.MustParseUUID(qFinal2ID)
	if err := upsertQuestion(tx, q2, assessmentID, 2, "true_false",
		"La fracción 2/4 es equivalente a 1/2.", strPtr("true"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"Verdadero", "Falso"} {
		if err := upsertQuestionOption(tx, q2, i, opt); err != nil {
			return err
		}
	}

	// Q3 — multiple_choice
	q3 := common.MustParseUUID(qFinal3ID)
	if err := upsertQuestion(tx, q3, assessmentID, 3, "multiple_choice",
		"¿Cuál de estas fracciones es la mayor?", strPtr("3/4"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"1/4", "1/2", "3/4", "2/8"} {
		if err := upsertQuestionOption(tx, q3, i, opt); err != nil {
			return err
		}
	}

	return nil
}

// seedPracticeQuestions siembra las 3 preguntas de "Práctica guiada — Geometría".
// Mismos tipos sencillos que la final.
func seedPracticeQuestions(tx *gorm.DB) error {
	assessmentID := common.MustParseUUID(assessmentPracticeID)

	// Q1 — multiple_choice
	q1 := common.MustParseUUID(qPractice1ID)
	if err := upsertQuestion(tx, q1, assessmentID, 1, "multiple_choice",
		"¿Cuántos lados tiene un triángulo?", strPtr("3"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"3", "4", "5", "6"} {
		if err := upsertQuestionOption(tx, q1, i, opt); err != nil {
			return err
		}
	}

	// Q2 — true_false
	q2 := common.MustParseUUID(qPractice2ID)
	if err := upsertQuestion(tx, q2, assessmentID, 2, "true_false",
		"Un cuadrado tiene los cuatro lados iguales.", strPtr("true"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"Verdadero", "Falso"} {
		if err := upsertQuestionOption(tx, q2, i, opt); err != nil {
			return err
		}
	}

	// Q3 — multiple_choice
	q3 := common.MustParseUUID(qPractice3ID)
	if err := upsertQuestion(tx, q3, assessmentID, 3, "multiple_choice",
		"¿Cómo se llama el perímetro de un círculo?", strPtr("Circunferencia"), 10); err != nil {
		return err
	}
	for i, opt := range []string{"Diámetro", "Radio", "Circunferencia", "Área"} {
		if err := upsertQuestionOption(tx, q3, i, opt); err != nil {
			return err
		}
	}

	return nil
}

// upsertQuestion siembra una pregunta (assessment.question). Idempotente por id.
func upsertQuestion(tx *gorm.DB, id, assessmentID uuid.UUID, sortOrder int, qType, text string, correctAnswer *string, points int) error {
	q := entities.Question{
		ID:            id,
		AssessmentID:  assessmentID,
		SortOrder:     sortOrder,
		QuestionText:  text,
		QuestionType:  qType,
		CorrectAnswer: correctAnswer,
		Points:        points,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&q).Error
}

// upsertQuestionOption siembra una opción de respuesta (assessment.question_option).
// Id derivado del question_id + sort_order → idempotente. La opción correcta NO se
// marca aquí: se referencia desde question.correct_answer (ver entity).
func upsertQuestionOption(tx *gorm.DB, questionID uuid.UUID, sortOrder int, text string) error {
	o := entities.QuestionOption{
		ID:         optionUUID(questionID, sortOrder),
		QuestionID: questionID,
		OptionText: text,
		SortOrder:  sortOrder,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&o).Error
}
