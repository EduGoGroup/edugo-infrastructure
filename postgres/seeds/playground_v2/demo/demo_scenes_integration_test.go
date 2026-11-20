//go:build integration
// +build integration

package demo_test

import (
	"math"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/demo"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestDemoScenes afirma, escena por escena del libreto, lo que cada cuenta tiene
// que encontrar en la base recién sembrada con demo.
func TestDemoScenes(t *testing.T) {
	gdb := startDemoDB(t)
	eb, ex, ep, er := demo.AssessmentDraftID.String(), demo.AssessmentExamID.String(),
		demo.AssessmentPracticeID.String(), demo.AssessmentReviewID.String()

	t.Run("E-B en borrador", func(t *testing.T) {
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment WHERE id = ? AND status = 'draft' AND max_attempts IS NOT NULL`, eb)
		expectInt(t, gdb, 5, `SELECT count(DISTINCT question_type) FROM assessment.question WHERE assessment_id = ?`, eb)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.assessment_assignment WHERE assessment_id = ?`, eb)
		expectInt(t, gdb, 0, `SELECT count(*) FROM notifications.notifications WHERE resource_id = ?`, eb)
	})

	t.Run("E-X por rendir", func(t *testing.T) {
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment
			WHERE id = ? AND status = 'published' AND purpose = 'exam' AND max_attempts = 2 AND subject_id = ?`,
			ex, demo.SubjectMathID.String())
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment_assignment WHERE assessment_id = ? AND subject_offering_id = ?`,
			ex, demo.OfferingMathID.String())
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.question
			WHERE assessment_id = ? AND question_type NOT IN ('multiple_choice','multiple_select','true_false')`, ex)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.assessment_attempt WHERE assessment_id = ?`, ex)
	})

	t.Run("E-P de práctica", func(t *testing.T) {
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment
			WHERE id = ? AND status = 'published' AND purpose = 'practice' AND max_attempts IS NULL`, ep)
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment_assignment WHERE assessment_id = ?`, ep)
	})

	t.Run("E-R con intentos por corregir", func(t *testing.T) {
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.assessment_assignment WHERE assessment_id = ? AND subject_offering_id = ?`,
			er, demo.OfferingLanguageID.String())
		n := scalarInt(t, gdb, `SELECT count(*) FROM assessment.assessment_attempt WHERE assessment_id = ? AND status = 'pending_review'`, er)
		if n < 3 || n > 5 {
			t.Errorf("E-R con %d intentos pending_review, esperado 3 a 5", n)
		}
		// Uno por alumno, de relleno (sin login) e inscrito en Lenguaje.
		expectInt(t, gdb, n, `SELECT count(DISTINCT a.student_membership_id) FROM assessment.assessment_attempt a
			JOIN academic.memberships m ON m.id = a.student_membership_id
			JOIN academic.subject_offering_enrollments e ON e.student_membership_id = m.id AND e.offering_id = ?
			WHERE a.assessment_id = ? AND m.user_id NOT IN ?`, demo.OfferingLanguageID.String(), er, accountIDs)
		// Una respuesta por pregunta; la abierta pending con 0 y sin is_correct; las de elegir auto_graded.
		nq := scalarInt(t, gdb, `SELECT count(*) FROM assessment.question WHERE assessment_id = ?`, er)
		expectInt(t, gdb, n*nq, `SELECT count(*) FROM assessment.assessment_attempt_answer r
			JOIN assessment.assessment_attempt a ON a.id = r.attempt_id WHERE a.assessment_id = ?`, er)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.assessment_attempt_answer r
			JOIN assessment.question q ON q.id = r.question_id
			WHERE q.assessment_id = ? AND NOT (
				(q.question_type = 'open_ended' AND r.review_status = 'pending' AND r.points_earned = 0 AND r.is_correct IS NULL)
				OR (q.question_type <> 'open_ended' AND r.review_status = 'auto_graded' AND r.is_correct IS NOT NULL
					AND r.points_earned = CASE WHEN r.is_correct THEN q.points ELSE 0 END))`, er)
		checkAttemptTotals(t, gdb, er)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.assessment_attempt
			WHERE assessment_id = ? AND NOT (completed_at > started_at AND time_spent_seconds BETWEEN 1 AND 7200)`, er)
		expectInt(t, gdb, 0, `SELECT count(*) FROM academic.grade_item`)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.attempt_review`)
		expectInt(t, gdb, 0, `SELECT count(*) FROM assessment.assessment_attempt WHERE status = 'in_progress'`)
	})

	t.Run("avisos", func(t *testing.T) {
		checkNotificationIDs(t, gdb)
		for _, student := range []string{demo.Student1UserID, demo.Student2UserID} {
			for _, rt := range []struct{ typ, resType, id string }{
				{"assessment_assigned", "assessment", ex},
				{"assessment_assigned", "assessment_practice", ep},
				{"material_assigned", "material", demo.MaterialID.String()},
			} {
				expectInt(t, gdb, 1, `SELECT count(*) FROM notifications.notifications
					WHERE user_id = ? AND type = ? AND resource_type = ? AND resource_id = ?
					AND NOT is_read AND school_id = ? AND unit_id = ?`,
					student, rt.typ, rt.resType, rt.id, demo.SchoolID.String(), demo.UnitID.String())
			}
		}
		expectInt(t, gdb, scalarInt(t, gdb, `SELECT count(*) FROM assessment.assessment_attempt WHERE assessment_id = ?`, er),
			`SELECT count(*) FROM notifications.notifications n
			JOIN assessment.assessment_attempt a ON a.id = n.resource_id
			WHERE n.user_id = ? AND n.type = 'assessment_attempt_recorded' AND n.resource_type = 'assessment_attempt'
			AND n.title = 'Evaluación por revisar' AND NOT n.is_read`, demo.TeacherUserID)
		expectInt(t, gdb, 0, `SELECT count(*) FROM notifications.notifications WHERE user_id = ?`, demo.GuardianUserID)
		expectInt(t, gdb, 0, `SELECT count(*) FROM notifications.notifications WHERE type = 'assessment_reviewed'`)
		expectInt(t, gdb, 0, `SELECT count(*) FROM notifications.device_tokens`)
	})

	t.Run("T1 con su PDF", func(t *testing.T) {
		mat := demo.MaterialID.String()
		expectInt(t, gdb, 1, `SELECT count(*) FROM content.materials WHERE id = ? AND status = 'ready' AND is_public AND subject_id = ?`,
			mat, demo.SubjectMathID.String())
		expectInt(t, gdb, 1, `SELECT count(*) FROM content.material_file
			WHERE id = ? AND material_id = ? AND file_url = ? AND file_type = 'application/pdf' AND file_size_bytes > 0`,
			demo.MaterialFileID.String(), mat, demo.MaterialFileKey())
		expectInt(t, gdb, 1, `SELECT count(*) FROM content.material_assignment
			WHERE material_id = ? AND subject_offering_id = ? AND available_from IS NULL AND available_until IS NULL`,
			mat, demo.OfferingMathID.String())
		expectInt(t, gdb, 0, `SELECT count(*) FROM content.material_pipeline_job`)
	})

	t.Run("notas", func(t *testing.T) {
		students := int64(2 + demo.FillerCount)
		expectInt(t, gdb, students*3, `SELECT count(*) FROM academic.grades WHERE period_id = ?`, demo.PeriodID.String())
		expectInt(t, gdb, students, `SELECT count(*) FROM (SELECT membership_id FROM academic.grades
			GROUP BY membership_id HAVING count(DISTINCT subject_id) = 3) x`)
		expectInt(t, gdb, 0, `SELECT count(*) FROM academic.grades WHERE grade_value < 1 OR grade_value > 7 OR grade_value IS NULL`)
		if scalarInt(t, gdb, `SELECT count(*) FROM academic.grades WHERE finalized_at IS NOT NULL`) == 0 ||
			scalarInt(t, gdb, `SELECT count(*) FROM academic.grades WHERE finalized_at IS NULL`) == 0 {
			t.Errorf("las notas deben ser unas finalizadas y otras no")
		}
	})

	t.Run("asistencia", func(t *testing.T) {
		expectInt(t, gdb, 0, `SELECT count(*) FROM academic.attendance
			WHERE subject_id IS NULL OR date >= current_date OR date < current_date - 21`)
		for _, m := range []string{demo.Student1MembershipID.String(), demo.Student2MembershipID.String()} {
			for _, status := range []string{"absent", "late"} {
				if scalarInt(t, gdb, `SELECT count(*) FROM academic.attendance WHERE membership_id = ? AND status = ?`, m, status) == 0 {
					t.Errorf("la cuenta %s sin asistencia %s", m, status)
				}
			}
		}
	})

	t.Run("anuncios", func(t *testing.T) {
		expectInt(t, gdb, 3, `SELECT count(*) FROM academic.announcements
			WHERE school_id = ? AND author_id = ? AND published_at >= now() - interval '14 days' AND published_at <= now()`,
			demo.SchoolID.String(), demo.TeacherUserID)
		expectInt(t, gdb, 1, `SELECT count(*) FROM academic.announcements WHERE scope = 'unit' AND academic_unit_id = ?`,
			demo.UnitID.String())
		expectInt(t, gdb, 1, `SELECT count(*) FROM academic.announcements WHERE scope = 'school' AND is_pinned`)
	})

	t.Run("recorte V5", func(t *testing.T) {
		expectInt(t, gdb, 3, `SELECT count(*) FROM iam.user_grants`)
		expectInt(t, gdb, 0, `SELECT count(*) FROM iam.user_grants WHERE effect <> 'deny'`)
		expectInt(t, gdb, 1, `SELECT count(*) FROM iam.user_grants WHERE user_id = ? AND permission_pattern = 'messaging.*'`,
			demo.TeacherUserID)
		expectInt(t, gdb, 2, `SELECT count(*) FROM iam.user_grants WHERE user_id = ?
			AND permission_pattern IN ('content.materials.read','content.materials.download')`, demo.GuardianUserID)
	})
}

// checkAttemptTotals: puntaje, máximo y porcentaje de cada intento son la suma
// de sus respuestas, como los calcula la API (domain.ComputeScore).
func checkAttemptTotals(t *testing.T, gdb *gorm.DB, assessmentID string) {
	t.Helper()
	var rows []struct {
		ID                          string
		Score, MaxScore, Percentage float64
		Earned, Max                 float64
	}
	if err := gdb.Raw(`SELECT a.id::text AS id, a.score, a.max_score, a.percentage,
			sum(r.points_earned) AS earned, sum(r.max_points) AS max
		FROM assessment.assessment_attempt a
		JOIN assessment.assessment_attempt_answer r ON r.attempt_id = a.id
		WHERE a.assessment_id = ? GROUP BY a.id`, assessmentID).Scan(&rows).Error; err != nil {
		t.Fatalf("leyendo intentos: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("sin intentos que revisar")
	}
	for _, r := range rows {
		want := math.Round(r.Earned/r.Max*100*100) / 100
		if r.Score != r.Earned || r.MaxScore != r.Max || r.Percentage != want {
			t.Errorf("intento %s: %v/%v/%v, esperado %v/%v/%v", r.ID, r.Score, r.MaxScore, r.Percentage, r.Earned, r.Max, want)
		}
	}
}

// checkNotificationIDs recalcula, desde cada fila, el id que daría la fórmula de
// la API.
func checkNotificationIDs(t *testing.T, gdb *gorm.DB) {
	t.Helper()
	var rows []struct{ ID, UserID, Type, ResourceType, ResourceID string }
	if err := gdb.Raw(`SELECT id::text AS id, user_id::text AS user_id, type, resource_type,
		resource_id::text AS resource_id FROM notifications.notifications`).Scan(&rows).Error; err != nil {
		t.Fatalf("leyendo avisos: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("sin avisos")
	}
	for _, r := range rows {
		want := uuid.NewSHA1(uuid.NameSpaceOID, []byte(r.UserID+r.Type+r.ResourceType+r.ResourceID))
		if r.ID != want.String() {
			t.Errorf("aviso %s: la fórmula da %s", r.ID, want)
		}
	}
}
