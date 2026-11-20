package demo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// sampleNows son los momentos sobre los que se arma el mundo: siete días
// seguidos (cubre cualquier día de la semana) y el cambio de año.
func sampleNows() []time.Time {
	start := time.Date(2026, 11, 20, 14, 30, 0, 0, time.UTC)
	nows := make([]time.Time, 0, 9)
	for d := range 7 {
		nows = append(nows, start.AddDate(0, 0, d))
	}
	return append(nows,
		time.Date(2026, 12, 31, 23, 59, 0, 0, time.UTC),
		time.Date(2027, 1, 1, 0, 1, 0, 0, time.UTC))
}

func TestQuestionsAreCoherent(t *testing.T) {
	w := buildWorld(sampleNows()[0])
	opts := map[uuid.UUID][]string{}
	for _, o := range w.options {
		opts[o.QuestionID] = append(opts[o.QuestionID], o.OptionText)
	}
	for _, a := range w.assessments {
		qs := w.questionsOf(a.ID)
		if a.QuestionsCount != len(qs) {
			t.Errorf("%s: questions_count=%d, preguntas=%d", a.Title, a.QuestionsCount, len(qs))
		}
		if a.PassThreshold != int(a.PassingScore) {
			t.Errorf("%s: umbrales distintos %d/%d", a.Title, a.PassThreshold, a.PassingScore)
		}
		if a.ShuffleQuestions || a.IsTimed || a.TimeLimitMinutes != nil {
			t.Errorf("%s: con mezcla o cronómetro", a.Title)
		}
		for i, q := range qs {
			if q.SortOrder != i+1 {
				t.Errorf("%s: sort_order %d en la posición %d", a.Title, q.SortOrder, i)
			}
			checkCorrectAnswer(t, a.Title, q.QuestionType, q.CorrectAnswer, opts[q.ID])
		}
	}
}

func checkCorrectAnswer(t *testing.T, title, typ string, correct *string, options []string) {
	t.Helper()
	has := func(s string) bool {
		for _, o := range options {
			if o == s {
				return true
			}
		}
		return false
	}
	switch typ {
	case qMultipleChoice:
		if len(options) < 2 || correct == nil || !has(*correct) {
			t.Errorf("%s: opción múltiple sin ≥2 opciones o con respuesta que no es una opción", title)
		}
	case qMultipleSelect:
		var texts []string
		if correct == nil || json.Unmarshal([]byte(*correct), &texts) != nil || len(texts) < 2 || len(options) < 3 {
			t.Fatalf("%s: selección múltiple mal formada (%v, %d opciones)", title, correct, len(options))
		}
		for _, c := range texts {
			if !has(c) {
				t.Errorf("%s: %q no es una opción", title, c)
			}
		}
	case qTrueFalse:
		if correct == nil || (*correct != "true" && *correct != "false") || len(options) != 0 {
			t.Errorf("%s: V/F debe ser true|false y sin opciones", title)
		}
	case qShortAnswer:
		if correct == nil || strings.TrimSpace(*correct) == "" {
			t.Errorf("%s: respuesta corta sin respuesta", title)
		}
	case qOpenEnded:
		if correct != nil || len(options) != 0 {
			t.Errorf("%s: la abierta no lleva respuesta ni opciones", title)
		}
	default:
		t.Errorf("%s: tipo desconocido %s", title, typ)
	}
}

func TestSceneQuestionTypes(t *testing.T) {
	w := buildWorld(sampleNows()[0])
	types := func(id uuid.UUID) map[string]int {
		out := map[string]int{}
		for _, q := range w.questionsOf(id) {
			out[q.QuestionType]++
		}
		return out
	}
	if got := types(AssessmentDraftID); len(got) != 5 {
		t.Errorf("E-B debe traer los cinco tipos: %v", got)
	}
	for typ := range types(AssessmentExamID) {
		if typ != qMultipleChoice && typ != qMultipleSelect && typ != qTrueFalse {
			t.Errorf("E-X con una pregunta no autocalificable: %s", typ)
		}
	}
	if got := types(AssessmentReviewID); got[qShortAnswer] != 0 || got[qOpenEnded] != 1 || got[qMultipleChoice] != 2 {
		t.Errorf("E-R debe ser 2 de opción y 1 abierta, sin corta: %v", got)
	}
}

func TestOwnIDsAreUniqueAndPrefixed(t *testing.T) {
	w := buildWorld(sampleNows()[0])
	var ids []uuid.UUID
	ids = append(ids, w.school.ID, w.period.ID, w.guardian.ID)
	for _, x := range w.units {
		ids = append(ids, x.ID)
	}
	for _, x := range w.subjects {
		ids = append(ids, x.ID)
	}
	for _, x := range w.offerings {
		ids = append(ids, x.ID)
	}
	for _, x := range w.memberships {
		ids = append(ids, x.ID)
	}
	for _, x := range w.userRoles {
		ids = append(ids, x.ID)
	}
	for _, x := range w.grants {
		ids = append(ids, x.ID)
	}
	for _, x := range w.assessments {
		ids = append(ids, x.ID)
	}
	for _, x := range w.questions {
		ids = append(ids, x.ID)
	}
	for _, x := range w.options {
		ids = append(ids, x.ID)
	}
	for _, x := range w.assignments {
		ids = append(ids, x.ID)
	}
	for _, x := range w.attempts {
		ids = append(ids, x.ID)
	}
	for _, x := range w.answers {
		ids = append(ids, x.ID)
	}
	for _, x := range w.materials {
		ids = append(ids, x.ID)
	}
	for _, x := range w.materialFiles {
		ids = append(ids, x.ID)
	}
	for _, x := range w.materialAsgs {
		ids = append(ids, x.ID)
	}
	for _, x := range w.grades {
		ids = append(ids, x.ID)
	}
	for _, x := range w.attendance {
		ids = append(ids, x.ID)
	}
	for _, x := range w.announcements {
		ids = append(ids, x.ID)
	}
	accounts := map[string]bool{TeacherUserID: true, Student1UserID: true, Student2UserID: true, GuardianUserID: true}
	for _, u := range w.users {
		if !accounts[u.ID.String()] {
			ids = append(ids, u.ID)
		}
	}
	seen := map[uuid.UUID]bool{}
	for _, id := range ids {
		if !strings.HasPrefix(id.String(), IDPrefix) {
			t.Errorf("id fuera del rango propio: %s", id)
		}
		if seen[id] {
			t.Errorf("id repetido: %s", id)
		}
		seen[id] = true
	}
	if len(w.users) != 4+FillerCount || FillerCount != 18 {
		t.Errorf("usuarios = %d (relleno %d), esperado 4 + 18", len(w.users), FillerCount)
	}
}

func TestWorldHoldsOnAnyDay(t *testing.T) {
	for _, now := range sampleNows() {
		w := buildWorld(now)
		today := now.Truncate(24 * time.Hour)
		if w.period.StartDate.After(today) || w.period.EndDate.Before(today.AddDate(0, 0, 90)) {
			t.Errorf("%s: el período %s–%s no cubre hoy y 90 días más", now, w.period.StartDate, w.period.EndDate)
		}
		if w.period.AcademicUnitID != RootUnitID {
			t.Errorf("%s: el período no está en la raíz", now)
		}
		counts := map[uuid.UUID]map[string]int{}
		for _, a := range w.attendance {
			if counts[a.MembershipID] == nil {
				counts[a.MembershipID] = map[string]int{}
			}
			counts[a.MembershipID][a.Status]++
			if a.SubjectID == nil || !a.Date.Before(today) || a.Date.Before(today.AddDate(0, 0, -attendanceWindowDays)) {
				t.Fatalf("%s: asistencia fuera de la ventana o sin materia: %+v", now, a)
			}
			if wd := a.Date.Weekday(); wd == time.Saturday || wd == time.Sunday {
				t.Fatalf("%s: asistencia en fin de semana: %s", now, a.Date)
			}
		}
		for _, m := range []uuid.UUID{Student1MembershipID, Student2MembershipID} {
			if counts[m]["absent"] < 1 || counts[m]["late"] < 1 {
				t.Errorf("%s: la cuenta %s sin ausencia o sin atraso: %v", now, m, counts[m])
			}
		}
		for _, an := range w.announcements {
			if an.PublishedAt.Before(today.AddDate(0, 0, -14)) || an.PublishedAt.After(now) {
				t.Errorf("%s: anuncio fuera de las dos últimas semanas: %s", now, an.PublishedAt)
			}
		}
		for _, at := range w.attempts {
			if !at.CompletedAt.After(at.StartedAt) || at.CompletedAt.After(now) {
				t.Errorf("%s: intento con fechas incoherentes: %s → %s", now, at.StartedAt, *at.CompletedAt)
			}
		}
	}
}

func TestAttemptScoresAddUp(t *testing.T) {
	w := buildWorld(sampleNows()[0])
	if n := len(w.attempts); n < 3 || n > 5 {
		t.Fatalf("E-R con %d intentos, esperado 3 a 5", n)
	}
	nq := len(w.questionsOf(AssessmentReviewID))
	for _, at := range w.attempts {
		var score, maxScore float64
		n := 0
		for _, ans := range w.answers {
			if ans.AttemptID != at.ID {
				continue
			}
			n++
			score += *ans.PointsEarned
			maxScore += *ans.MaxPoints
		}
		if n != nq {
			t.Errorf("intento %s con %d respuestas, esperado %d", at.ID, n, nq)
		}
		if *at.Score != score || *at.MaxScore != maxScore || *at.Percentage != percentage(score, maxScore) {
			t.Errorf("intento %s: %v/%v/%v, esperado %v/%v/%v", at.ID,
				*at.Score, *at.MaxScore, *at.Percentage, score, maxScore, percentage(score, maxScore))
		}
		if s := *at.TimeSpentSeconds; s < 1 || s > 7200 {
			t.Errorf("intento %s: time_spent_seconds=%d", at.ID, s)
		}
	}
}

// TestNotificationIDFormula compara la fórmula con un caso calculado aparte con
// uuid.NewSHA1(uuid.NameSpaceOID, …) sobre la cadena concatenada a mano.
func TestNotificationIDFormula(t *testing.T) {
	user := uuid.MustParse(Student1UserID)
	resource := uuid.MustParse("de00000a-0000-0000-0000-000000000002")
	raw := "d4544906-e877-4f8d-8d0d-0e41d252fdbf" + "assessment_assigned" + "assessment" + "de00000a-0000-0000-0000-000000000002"
	want := uuid.NewSHA1(uuid.NameSpaceOID, []byte(raw))
	if got := NotificationID(user, "assessment_assigned", "assessment", resource); got != want {
		t.Fatalf("NotificationID = %s, esperado %s", got, want)
	}
	if resource != AssessmentExamID {
		t.Fatalf("el caso de prueba supone que E-X es %s, y es %s", resource, AssessmentExamID)
	}
	w := buildWorld(sampleNows()[0])
	for _, n := range w.notifications {
		if n.ID != NotificationID(n.UserID, n.Type, *n.ResourceType, *n.ResourceID) {
			t.Errorf("aviso %s no sigue la fórmula", n.ID)
		}
	}
}

func TestMaterialFileKey(t *testing.T) {
	want := "materials/" + MaterialID.String() + "/files/" + MaterialFileID.String()
	if MaterialFileKey() != want {
		t.Fatalf("MaterialFileKey() = %s, esperado %s", MaterialFileKey(), want)
	}
}
