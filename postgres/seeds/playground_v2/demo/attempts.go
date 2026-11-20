package demo

import (
	"fmt"
	"math"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
)

// reviewAttempt es un intento entregado de E-R: qué eligió el alumno en cada
// pregunta de opción y qué escribió en la abierta.
type reviewAttempt struct {
	Filler  int      // índice del alumno de relleno
	Choices []string // una por pregunta de opción, en orden natural
	Open    string
}

// ReviewAttemptCount es el número de intentos `pending_review` de E-R.
var ReviewAttemptCount = len(reviewAttempts())

func reviewAttempts() []reviewAttempt {
	return []reviewAttempt{
		{Filler: 0, Choices: []string{"Guardaba comida", "Narrativo"},
			Open: "Que hay que trabajar y prepararse antes de que lleguen los tiempos difíciles."},
		{Filler: 1, Choices: []string{"Guardaba comida", "Informativo"},
			Open: "La cigarra no guardó comida y después pasó hambre, entonces hay que ser previsor."},
		{Filler: 2, Choices: []string{"Cantaba", "Narrativo"},
			Open: "Que la hormiga es mala porque no comparte."},
		{Filler: 3, Choices: []string{"Guardaba comida", "Instructivo"},
			Open: "Hay que pensar en el futuro y no solo divertirse."},
	}
}

// buildAttempts arma los intentos de E-R tal como los habría escrito la API al
// entregar (learning, usecase/attempt/submit.go): una respuesta por pregunta; las
// de opción `auto_graded` con sus puntos; la abierta `pending`, 0 puntos y
// is_correct nulo; el intento `pending_review` con puntaje, máximo y porcentaje
// sumados de sus respuestas (domain.ComputeScore). Entregados ayer, sin
// grade_item ni attempt_review: eso lo escribe la revisión.
func buildAttempts(w *world) {
	questions := w.questionsOf(AssessmentReviewID)
	ansN := 0
	for i, ra := range reviewAttempts() {
		attemptID := newID(domAttempt, i+1)
		started := w.day(-1).Add(15*time.Hour + time.Duration(i)*25*time.Minute)
		spent := (9 + 3*i) * 60
		completed := started.Add(time.Duration(spent) * time.Second)

		var score, maxScore float64
		choice := 0
		for idx, q := range questions {
			ansN++
			qID := q.ID
			maxPts := float64(q.Points)
			ans := entities.AssessmentAttemptAnswer{
				ID: newID(domAnswer, ansN), AttemptID: attemptID, QuestionIndex: idx, QuestionID: &qID,
				MaxPoints: &maxPts, CreatedAt: completed, UpdatedAt: completed,
			}
			if q.QuestionType == qOpenEnded {
				ans.StudentAnswer = ptr(ra.Open)
				ans.PointsEarned = ptr(0.0)
				ans.ReviewStatus = "pending"
			} else {
				picked := ra.Choices[choice]
				choice++
				correct := q.CorrectAnswer != nil && picked == *q.CorrectAnswer
				earned := 0.0
				if correct {
					earned = maxPts
				}
				ans.StudentAnswer = ptr(picked)
				ans.IsCorrect = ptr(correct)
				ans.PointsEarned = ptr(earned)
				ans.ReviewStatus = "auto_graded"
			}
			score += *ans.PointsEarned
			maxScore += maxPts
			w.answers = append(w.answers, ans)
		}

		w.attempts = append(w.attempts, entities.AssessmentAttempt{
			ID: attemptID, AssessmentID: AssessmentReviewID, StudentMembershipID: fillerMembershipID(ra.Filler),
			StartedAt: started, CompletedAt: ptr(completed), Score: ptr(score), MaxScore: ptr(maxScore),
			Percentage: ptr(percentage(score, maxScore)), TimeSpentSeconds: ptr(spent),
			Status: "pending_review", CreatedAt: started, UpdatedAt: completed,
		})
	}
}

// percentage es score/max*100 redondeado a dos decimales, como lo guarda la
// columna decimal(5,2).
func percentage(score, maxScore float64) float64 {
	if maxScore == 0 {
		return 0
	}
	return math.Round(score/maxScore*100*100) / 100
}

func writeAttempts(tx *gorm.DB, w *world) error {
	if err := insert(tx, &w.attempts); err != nil {
		return fmt.Errorf("assessment_attempt: %w", err)
	}
	if err := insert(tx, &w.answers); err != nil {
		return fmt.Errorf("assessment_attempt_answer: %w", err)
	}
	return nil
}
