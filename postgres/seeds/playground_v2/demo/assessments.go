package demo

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Umbral de aprobación: los dos campos iguales (deuda 024: si difieren, el
// alumno puede ver «reprobado» y «Ya aprobaste» a la vez).
const passThreshold = 60

// Tipos de pregunta (CHECK question_type_check).
const (
	qMultipleChoice = "multiple_choice"
	qMultipleSelect = "multiple_select"
	qTrueFalse      = "true_false"
	qShortAnswer    = "short_answer"
	qOpenEnded      = "open_ended"
)

// questionDef es una pregunta antes de tener id. Correct sigue el formato que
// valida learning (usecase/question/validate.go): el texto exacto de una opción,
// un array JSON de textos, "true"/"false", un texto, o nil en la abierta.
type questionDef struct {
	Type    string
	Text    string
	Correct *string
	Options []string
	Points  int
}

// assessmentDef es una evaluación de demo con sus preguntas y, si se asigna, la
// sesión y los días atrás de la asignación.
type assessmentDef struct {
	ID          uuid.UUID
	SubjectID   uuid.UUID
	Title       string
	Description string
	Status      string
	Purpose     string
	MaxAttempts *int
	CreatedDays int
	Offering    *uuid.UUID // nil = sin asignación
	AssignDays  int
	Questions   []questionDef
}

func assessmentDefs() []assessmentDef {
	one, two := 1, 2
	return []assessmentDef{
		{
			// E-B: se edita, se publica y se asigna en vivo. Una pregunta de cada tipo.
			ID: AssessmentDraftID, SubjectID: SubjectMathID, Title: "Repaso de la unidad: números enteros",
			Description: "Borrador para revisar y publicar en clase.", Status: "draft", Purpose: "exam",
			MaxAttempts: &one, CreatedDays: -3,
			Questions: []questionDef{
				{Type: qMultipleChoice, Text: "¿Cuánto es -7 + 12?", Correct: ptr("5"), Options: []string{"5", "-5", "19"}, Points: 2},
				{Type: qMultipleSelect, Text: "¿Cuáles de estos números son negativos?", Correct: ptr(`["-3","-10"]`), Options: []string{"-3", "4", "-10", "0"}, Points: 2},
				{Type: qTrueFalse, Text: "El cero es un número entero.", Correct: ptr("true"), Points: 1},
				{Type: qShortAnswer, Text: "Escribe el opuesto de 8.", Correct: ptr("-8"), Points: 1},
				{Type: qOpenEnded, Text: "Explica con un ejemplo de la vida diaria para qué sirven los números negativos.", Points: 2},
			},
		},
		{
			// E-X: la rinden Matías (escritorio) y Antonia (web). Solo de elegir.
			ID: AssessmentExamID, SubjectID: SubjectMathID, Title: "Prueba de fracciones",
			Description: "Fracciones equivalentes, simplificación y comparación.", Status: "published", Purpose: "exam",
			MaxAttempts: &two, CreatedDays: -6, Offering: ptr(OfferingMathID), AssignDays: -2,
			Questions: []questionDef{
				{Type: qMultipleChoice, Text: "¿Qué fracción es equivalente a 2/4?", Correct: ptr("1/2"), Options: []string{"1/2", "2/8", "4/2", "3/4"}, Points: 2},
				{Type: qMultipleSelect, Text: "¿Cuáles de estas fracciones son mayores que 1/2?", Correct: ptr(`["3/4","5/6"]`), Options: []string{"3/4", "1/3", "5/6", "2/5"}, Points: 2},
				{Type: qTrueFalse, Text: "La fracción 6/9 simplificada es 2/3.", Correct: ptr("true"), Points: 1},
				{Type: qMultipleChoice, Text: "¿Cuánto es 1/4 + 2/4?", Correct: ptr("3/4"), Options: []string{"3/8", "3/4", "1/2", "2/4"}, Points: 2},
			},
		},
		{
			// E-P: práctica ilimitada de Ciencias Naturales.
			ID: AssessmentPracticeID, SubjectID: SubjectScienceID, Title: "Práctica: el sistema digestivo",
			Description: "Repaso libre antes de la prueba.", Status: "published", Purpose: "practice",
			MaxAttempts: nil, CreatedDays: -8, Offering: ptr(OfferingScienceID), AssignDays: -3,
			Questions: []questionDef{
				{Type: qMultipleChoice, Text: "¿En qué órgano comienza la digestión?", Correct: ptr("La boca"), Options: []string{"La boca", "El estómago", "El intestino delgado"}, Points: 1},
				{Type: qTrueFalse, Text: "El intestino delgado absorbe la mayor parte de los nutrientes.", Correct: ptr("true"), Points: 1},
				{Type: qMultipleChoice, Text: "¿Qué órgano produce la bilis?", Correct: ptr("El hígado"), Options: []string{"El páncreas", "El hígado", "El estómago", "El esófago"}, Points: 1},
				{Type: qTrueFalse, Text: "El esófago digiere las proteínas.", Correct: ptr("false"), Points: 1},
				{Type: qMultipleSelect, Text: "¿Cuáles forman parte del tubo digestivo?", Correct: ptr(`["El esófago","El estómago"]`), Options: []string{"El esófago", "El pulmón", "El estómago", "El riñón"}, Points: 2},
			},
		},
		{
			// E-R: intentos entregados en revisión. Sin respuesta corta (el docente no
			// puede puntuarla).
			ID: AssessmentReviewID, SubjectID: SubjectLanguageID, Title: "Comprensión lectora: la fábula",
			Description: "Lectura de «La cigarra y la hormiga».", Status: "published", Purpose: "exam",
			MaxAttempts: &one, CreatedDays: -9, Offering: ptr(OfferingLanguageID), AssignDays: -5,
			Questions: []questionDef{
				{Type: qMultipleChoice, Text: "¿Qué hacía la hormiga durante el verano?", Correct: ptr("Guardaba comida"), Options: []string{"Cantaba", "Guardaba comida", "Dormía"}, Points: 2},
				{Type: qMultipleChoice, Text: "¿Qué tipo de texto es una fábula?", Correct: ptr("Narrativo"), Options: []string{"Informativo", "Narrativo", "Instructivo"}, Points: 2},
				{Type: qOpenEnded, Text: "¿Cuál es la moraleja de la fábula? Explícala con tus palabras.", Points: 4},
			},
		},
	}
}

// buildAssessments arma las cuatro evaluaciones, sus preguntas y opciones, y las
// asignaciones a las sesiones de la docente. Ids de pregunta y opción por
// contador dentro de su dominio.
func buildAssessments(w *world) {
	qn, on, an := 0, 0, 0
	for _, d := range assessmentDefs() {
		created := w.day(d.CreatedDays)
		w.assessments = append(w.assessments, entities.Assessment{
			ID: d.ID, SchoolID: SchoolID, CreatedByMembershipID: TeacherMembershipID, SubjectID: d.SubjectID,
			Title: d.Title, Description: ptr(d.Description), SourceType: "manual", Status: d.Status,
			Purpose: d.Purpose, QuestionsCount: len(d.Questions), PassThreshold: passThreshold,
			PassingScore: passThreshold, MaxAttempts: d.MaxAttempts, IsTimed: false, IsPublic: false,
			ShuffleQuestions: false, ShowCorrectAnswers: true, CreatedAt: created, UpdatedAt: created,
		})
		for i, q := range d.Questions {
			qn++
			qID := newID(domQuestion, qn)
			w.questions = append(w.questions, entities.Question{
				ID: qID, AssessmentID: d.ID, SortOrder: i + 1, QuestionText: q.Text, QuestionType: q.Type,
				CorrectAnswer: q.Correct, Points: q.Points, LLMPrepStatus: "none",
				CreatedAt: created, UpdatedAt: created,
			})
			for j, text := range q.Options {
				on++
				w.options = append(w.options, entities.QuestionOption{
					ID: newID(domOption, on), QuestionID: qID, OptionText: text, SortOrder: j + 1,
					CreatedAt: created, UpdatedAt: created,
				})
			}
		}
		if d.Offering != nil {
			an++
			at := w.day(d.AssignDays).Add(9 * time.Hour) // 09:00 UTC
			w.assignments = append(w.assignments, entities.AssessmentAssignment{
				ID: newID(domAssignment, an), AssessmentID: d.ID, SubjectOfferingID: *d.Offering,
				AssignedByMembershipID: TeacherMembershipID, AssignedAt: at, CreatedAt: at, UpdatedAt: at,
			})
		}
	}
}

// questionsOf devuelve las preguntas de una evaluación en su orden natural.
func (w *world) questionsOf(assessmentID uuid.UUID) []entities.Question {
	var out []entities.Question
	for _, q := range w.questions {
		if q.AssessmentID == assessmentID {
			out = append(out, q)
		}
	}
	return out
}

// assessment devuelve la evaluación de demo con ese id.
func (w *world) assessment(id uuid.UUID) entities.Assessment {
	for _, a := range w.assessments {
		if a.ID == id {
			return a
		}
	}
	panic(fmt.Sprintf("demo: evaluación desconocida %s", id))
}

func writeAssessments(tx *gorm.DB, w *world) error {
	if err := insert(tx, &w.assessments); err != nil {
		return fmt.Errorf("assessment: %w", err)
	}
	if err := insert(tx, &w.questions); err != nil {
		return fmt.Errorf("question: %w", err)
	}
	if err := insert(tx, &w.options); err != nil {
		return fmt.Errorf("question_option: %w", err)
	}
	if err := insert(tx, &w.assignments); err != nil {
		return fmt.Errorf("assessment_assignment: %w", err)
	}
	return nil
}
