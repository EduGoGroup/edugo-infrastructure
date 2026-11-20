package demo

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Los textos y tipos son los que escribe la API (learning,
// usecase/assignment/create.go, usecase/material_assignment/create.go y
// usecase/attempt/submit.go). Se copian tal cual, sin tildes donde la API no las lleva.
const (
	notifAssessmentAssigned = "assessment_assigned"
	notifMaterialAssigned   = "material_assigned"
	notifAttemptRecorded    = "assessment_attempt_recorded"

	titleAssessmentAssigned = "Nueva evaluacion asignada"
	titleMaterialAssigned   = "Nuevo material asignado"
	titleAttemptToReview    = "Evaluación por revisar"
)

// buildNotifications arma la bandeja de cada cuenta:
//   - Matías y Antonia: la de E-X (`assessment`), la de E-P
//     (`assessment_practice`, abre la práctica) y la de T1 (`material`).
//   - La docente: una «Evaluación por revisar» por intento de E-R.
//   - La tutora: ninguna (nada le notifica).
//
// Nada de lo que se asigna en escena (E-B), ninguna `assessment_reviewed` y
// ningún token de push. Todas sin leer, con escuela y unidad, y con el id de la
// fórmula de la API.
func buildNotifications(w *world) {
	assignedAt := map[uuid.UUID]time.Time{}
	for _, a := range w.assignments {
		assignedAt[a.AssessmentID] = a.AssignedAt
	}
	school, unit := SchoolID, UnitID
	add := func(user uuid.UUID, typ, title, body, resourceType string, resourceID uuid.UUID, at time.Time) {
		w.notifications = append(w.notifications, entities.Notification{
			ID: NotificationID(user, typ, resourceType, resourceID), UserID: user, Type: typ, Title: title,
			Body: ptr(body), ResourceType: ptr(resourceType), ResourceID: ptr(resourceID),
			SchoolID: &school, UnitID: &unit, IsRead: false, CreatedAt: at,
		})
	}

	exam := w.assessment(AssessmentExamID)
	practice := w.assessment(AssessmentPracticeID)
	material := w.materials[0]
	for _, user := range []string{Student1UserID, Student2UserID} {
		u := uuid.MustParse(user)
		add(u, notifAssessmentAssigned, titleAssessmentAssigned, "Te han asignado: "+exam.Title,
			"assessment", exam.ID, assignedAt[exam.ID])
		add(u, notifAssessmentAssigned, titleAssessmentAssigned, "Te han asignado: "+practice.Title,
			"assessment_practice", practice.ID, assignedAt[practice.ID])
		add(u, notifMaterialAssigned, titleMaterialAssigned, "Te han asignado: "+material.Title,
			"material", material.ID, w.materialAsgs[0].AssignedAt)
	}

	review := w.assessment(AssessmentReviewID)
	teacher := uuid.MustParse(TeacherUserID)
	for _, at := range w.attempts {
		add(teacher, notifAttemptRecorded, titleAttemptToReview,
			fmt.Sprintf("Un estudiante entregó y requiere tu revisión: %s", review.Title),
			"assessment_attempt", at.ID, *at.CompletedAt)
	}
}
