package demo

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	SeedVersion         = "development-gorm-v3"
	defaultPasswordHash = "$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau"
)

func ApplyDemo(gdb *gorm.DB) error {
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := truncateDevelopmentData(tx); err != nil {
			return err
		}
		if err := seedSchools(tx); err != nil {
			return err
		}
		if err := seedAcademicUnits(tx); err != nil {
			return err
		}
		if err := seedUsers(tx); err != nil {
			return err
		}
		if err := seedMemberships(tx); err != nil {
			return err
		}
		if err := seedUserRoles(tx); err != nil {
			return err
		}
		if err := seedUserGrants(tx); err != nil {
			return err
		}
		if err := seedSubjects(tx); err != nil {
			return err
		}
		if err := seedMaterials(tx); err != nil {
			return err
		}
		if err := seedAssessments(tx); err != nil {
			return err
		}
		if err := seedAssessmentMaterials(tx); err != nil {
			return err
		}
		if err := seedQuestions(tx); err != nil {
			return err
		}
		if err := seedQuestionOptions(tx); err != nil {
			return err
		}
		if err := seedAssessmentAssignments(tx); err != nil {
			return err
		}
		if err := seedAssessmentAttempts(tx); err != nil {
			return err
		}
		if err := seedAssessmentAttemptAnswers(tx); err != nil {
			return err
		}
		if err := seedGuardianRelations(tx); err != nil {
			return err
		}
		if err := seedScreenUserPreferences(tx); err != nil {
			return err
		}
		if err := seedSchoolConcepts(tx); err != nil {
			return err
		}
		if err := seedProgress(tx); err != nil {
			return err
		}
		if err := seedAcademicPeriods(tx); err != nil {
			return err
		}
		// subject_offerings + enrollments dependen de academic_periods (FK
		// period_id), por eso van despues de seedAcademicPeriods.
		if err := seedSubjectOfferings(tx); err != nil {
			return err
		}
		if err := seedGrades(tx); err != nil {
			return err
		}
		if err := seedAttendance(tx); err != nil {
			return err
		}
		if err := seedSchedules(tx); err != nil {
			return err
		}
		if err := seedAnnouncements(tx); err != nil {
			return err
		}
		if err := seedCalendarEvents(tx); err != nil {
			return err
		}
		if err := seedCourses(tx); err != nil {
			return err
		}
		return nil
	})
}

func truncateDevelopmentData(tx *gorm.DB) error {
	guarded := []string{
		"assessment.attempt_analytics",
		"assessment.assessment_stats",
		// Sesiones de materia (plan 010 N1.7). Se truncan ANTES que
		// academic.memberships (mas abajo, con CASCADE): enrollments primero
		// (FK a offerings) y offerings despues (FK a memberships/periods).
		"academic.subject_offering_enrollments",
		"academic.subject_offerings",
	}
	for _, table := range guarded {
		if err := truncateIfExists(tx, table); err != nil {
			return err
		}
	}

	required := []string{
		"assessment.assessment_attempt_answer",
		"assessment.assessment_attempt",
		"assessment.assessment_materials",
		"assessment.assessment",
		"content.materials",
		"academic.memberships",
		"iam.user_grants",
		"iam.user_roles",
		"academic.academic_units",
	}
	// content.courses puede no existir en entornos sin la migración de learning
	if err := truncateIfExists(tx, "content.courses"); err != nil {
		return err
	}
	for _, table := range required {
		if err := truncateTable(tx, table); err != nil {
			return err
		}
	}

	if err := tx.Table("ui_config.screen_templates").
		Where("created_by IS NOT NULL").
		Update("created_by", nil).Error; err != nil {
		return fmt.Errorf("error desacoplando ui_config.screen_templates: %w", err)
	}
	if err := tx.Table("ui_config.screen_instances").
		Where("created_by IS NOT NULL").
		Update("created_by", nil).Error; err != nil {
		return fmt.Errorf("error desacoplando ui_config.screen_instances: %w", err)
	}

	guarded = []string{
		"academic.calendar_events",
		"academic.announcements",
		"academic.schedules",
		"academic.attendance",
		"academic.grades",
		"academic.academic_periods",
	}
	for _, table := range guarded {
		if err := truncateIfExists(tx, table); err != nil {
			return err
		}
	}

	required = []string{
		"academic.subjects",
		"academic.guardian_relations",
		"academic.school_concepts",
		"ui_config.screen_user_preferences",
		"auth.refresh_tokens",
		"auth.login_attempts",
		"academic.schools",
	}
	for _, table := range required {
		if err := truncateTable(tx, table); err != nil {
			return err
		}
	}

	if err := tx.Exec("DELETE FROM auth.users").Error; err != nil {
		return fmt.Errorf("error limpiando auth.users: %w", err)
	}

	return nil
}

func seedSchools(tx *gorm.DB) error {
	rows := []map[string]any{
		{
			"id":                mustUUID("b1000000-0000-0000-0000-000000000001"),
			"name":              "Colegio San Ignacio",
			"code":              "SCH_SI_001",
			"address":           "Av. Libertador 1500",
			"city":              "Santiago",
			"country":           "Chile",
			"phone":             "+56 2 2345 6789",
			"email":             "contacto@sanignacio.edugo.test",
			"concept_type_id":   mustUUID("c1000000-0000-0000-0000-000000000002"),
			"metadata":          mustJSON(`{"level":"secondary","demo":true,"founded_year":2018}`),
			"is_active":         true,
			"subscription_tier": "premium",
			"max_teachers":      20,
			"max_students":      300,
		},
		{
			"id":                mustUUID("b2000000-0000-0000-0000-000000000002"),
			"name":              "Taller CreArte",
			"code":              "SCH_CA_001",
			"address":           "Calle Artistas 234",
			"city":              "Valparaiso",
			"country":           "Chile",
			"phone":             "+56 32 2345 678",
			"email":             "contacto@crearte.edugo.test",
			"concept_type_id":   mustUUID("c1000000-0000-0000-0000-000000000005"),
			"metadata":          mustJSON(`{"level":"workshop","demo":true,"founded_year":2021}`),
			"is_active":         true,
			"subscription_tier": "basic",
			"max_teachers":      10,
			"max_students":      100,
		},
		{
			"id":                mustUUID("b3000000-0000-0000-0000-000000000003"),
			"name":              "Academia Global English",
			"code":              "SCH_GE_001",
			"address":           "Paseo Internacional 89",
			"city":              "Santiago",
			"country":           "Chile",
			"phone":             "+56 2 9876 5432",
			"email":             "contacto@globalenglish.edugo.test",
			"concept_type_id":   mustUUID("c1000000-0000-0000-0000-000000000003"),
			"metadata":          mustJSON(`{"level":"language_academy","demo":true,"founded_year":2020}`),
			"is_active":         true,
			"subscription_tier": "basic",
			"max_teachers":      10,
			"max_students":      150,
		},
	}
	return upsertMaps(
		tx,
		"academic.schools",
		rows,
		[]string{"id"},
		[]string{
			"name", "code", "address", "city", "country", "phone", "email",
			"concept_type_id", "metadata", "is_active", "subscription_tier", "max_teachers", "max_students",
		},
		true,
	)
}

func seedAcademicUnits(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("ac000000-0000-0000-0000-000000000001"), "parent_unit_id": nil, "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "Colegio San Ignacio", "code": "CSI-ROOT", "type": "school", "description": "Unidad raiz del Colegio San Ignacio", "level": "secondary", "academic_year": 0, "metadata": mustJSON(`{"is_root":true}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000007"), "parent_unit_id": nil, "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Taller CreArte", "code": "TCA-ROOT", "type": "school", "description": "Unidad raiz del Taller CreArte", "level": "workshop", "academic_year": 0, "metadata": mustJSON(`{"is_root":true}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000012"), "parent_unit_id": nil, "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Academia Global English", "code": "AGE-ROOT", "type": "school", "description": "Unidad raiz de la Academia Global English", "level": "language", "academic_year": 0, "metadata": mustJSON(`{"is_root":true}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000002"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to Basico", "code": "GRADE-05", "type": "grade", "description": "Quinto ano de educacion basica, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000005"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "6to Basico", "code": "GRADE-06", "type": "grade", "description": "Sexto ano de educacion basica, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"grade_number":6}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000008"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Modulo Pintura", "code": "MOD-PINT", "type": "grade", "description": "Modulo de tecnicas de pintura", "level": "workshop", "academic_year": 2026, "metadata": mustJSON(`{"module_type":"pintura"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000010"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Modulo Escultura", "code": "MOD-ESCL", "type": "grade", "description": "Modulo de fundamentos de escultura", "level": "workshop", "academic_year": 2026, "metadata": mustJSON(`{"module_type":"escultura"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000013"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000012"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Level A2", "code": "LVL-A2", "type": "grade", "description": "Elementary level A2", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"cefr_level":"A2"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000015"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000012"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Level B1", "code": "LVL-B1", "type": "grade", "description": "Intermediate level B1", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"cefr_level":"B1"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000003"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to A", "code": "5A", "type": "class", "description": "Seccion A del 5to Basico, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"section":"A","grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000004"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to B", "code": "5B", "type": "class", "description": "Seccion B del 5to Basico, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"section":"B","grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000006"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "6to A", "code": "6A", "type": "class", "description": "Seccion A del 6to Basico, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"section":"A","grade_number":6}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000009"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Grupo Manana", "code": "GRP-MAN", "type": "class", "description": "Grupo de la manana - Modulo Pintura", "level": "workshop", "academic_year": 2026, "metadata": mustJSON(`{"schedule":"morning"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000011"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000010"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Grupo Tarde", "code": "GRP-TAR", "type": "class", "description": "Grupo de la tarde - Modulo Escultura", "level": "workshop", "academic_year": 2026, "metadata": mustJSON(`{"schedule":"afternoon"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000014"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000013"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Class Monday", "code": "CLS-MON", "type": "class", "description": "Monday class - Level A2", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"day":"monday"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000016"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000015"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Class Tuesday", "code": "CLS-TUE", "type": "class", "description": "Tuesday class - Level B1", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"day":"tuesday"}`), "is_active": true},
	}

	return upsertMaps(
		tx,
		"academic.academic_units",
		rows,
		[]string{"school_id", "code", "academic_year"},
		[]string{"name", "parent_unit_id", "description", "level", "metadata", "is_active"},
		true,
	)
}

func seedUsers(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("00000000-0000-0000-0000-000000000001"), "email": "super@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Santiago", "last_name": "Ramirez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000002"), "email": "admin.sanignacio@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Carmen", "last_name": "Valdes", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000003"), "email": "admin.crearte@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Roberto", "last_name": "Silva", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000004"), "email": "coord.academico@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Lucia", "last_name": "Fernandez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000005"), "email": "prof.martinez@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Maria", "last_name": "Martinez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000006"), "email": "prof.gonzalez@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Pedro", "last_name": "Gonzalez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000007"), "email": "facilitador.ruiz@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Ana", "last_name": "Ruiz", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000008"), "email": "est.carlos@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Carlos", "last_name": "Mendoza", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000009"), "email": "est.sofia@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Sofia", "last_name": "Herrera", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000010"), "email": "est.diego@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Diego", "last_name": "Vargas", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000011"), "email": "est.valentina@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Valentina", "last_name": "Rojas", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000012"), "email": "est.mateo@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Mateo", "last_name": "Fuentes", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000013"), "email": "tutor.mendoza@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Ricardo", "last_name": "Mendoza", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000014"), "email": "tutora.herrera@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Patricia", "last_name": "Herrera", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000015"), "email": "admin.plataforma@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Elena", "last_name": "Torres", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000016"), "email": "director.sanignacio@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Miguel", "last_name": "Castillo", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000017"), "email": "asist.admin@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Laura", "last_name": "Pena", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000018"), "email": "asist.prof@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Andres", "last_name": "Gomez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000019"), "email": "observador@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Diana", "last_name": "Lopez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000020"), "email": "guardian.pendiente@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Fernando", "last_name": "Ruiz", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000021"), "email": "readonly@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Test", "last_name": "ReadOnly", "is_active": true},
	}

	return upsertMaps(
		tx,
		"auth.users",
		rows,
		[]string{"id"},
		[]string{"email", "password_hash", "first_name", "last_name", "is_active"},
		true,
	)
}

func seedMemberships(tx *gorm.DB) error {
	parse := func(v string) time.Time { return mustTimestamp(v) }
	rows := []map[string]any{
		{"id": mustUUID("bb000000-0000-0000-0000-000000000001"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000002"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000003"), "user_id": mustUUID("00000000-0000-0000-0000-000000000009"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000004"), "user_id": mustUUID("00000000-0000-0000-0000-000000000010"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000005"), "user_id": mustUUID("00000000-0000-0000-0000-000000000011"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000006"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000006"), "user_id": mustUUID("00000000-0000-0000-0000-000000000011"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000007"), "user_id": mustUUID("00000000-0000-0000-0000-000000000012"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "role": "student", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000008"), "user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000009"), "user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000010"), "user_id": mustUUID("00000000-0000-0000-0000-000000000006"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000011"), "user_id": mustUUID("00000000-0000-0000-0000-000000000006"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000006"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000012"), "user_id": mustUUID("00000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000013"), "user_id": mustUUID("00000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000011"), "role": "teacher", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000014"), "user_id": mustUUID("00000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "role": "admin", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-01-15 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000015"), "user_id": mustUUID("00000000-0000-0000-0000-000000000003"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": nil, "role": "admin", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-01-15 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000016"), "user_id": mustUUID("00000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "role": "coordinator", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000017"), "user_id": mustUUID("00000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": nil, "role": "coordinator", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000018"), "user_id": mustUUID("00000000-0000-0000-0000-000000000013"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "guardian", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000019"), "user_id": mustUUID("00000000-0000-0000-0000-000000000013"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "role": "guardian", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000020"), "user_id": mustUUID("00000000-0000-0000-0000-000000000014"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "guardian", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000021"), "user_id": mustUUID("00000000-0000-0000-0000-000000000014"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "role": "guardian", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000022"), "user_id": mustUUID("00000000-0000-0000-0000-000000000016"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "role": "admin", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-01-20 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000023"), "user_id": mustUUID("00000000-0000-0000-0000-000000000017"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "role": "assistant", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-01 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000024"), "user_id": mustUUID("00000000-0000-0000-0000-000000000018"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "assistant", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-15 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000025"), "user_id": mustUUID("00000000-0000-0000-0000-000000000019"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "assistant", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-20 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000026"), "user_id": mustUUID("00000000-0000-0000-0000-000000000019"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "role": "assistant", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-20 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000027"), "user_id": mustUUID("00000000-0000-0000-0000-000000000020"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "role": "guardian", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-10 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000028"), "user_id": mustUUID("00000000-0000-0000-0000-000000000021"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "role": "admin", "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-25 08:00:00+00")},
	}

	return upsertMaps(
		tx,
		"academic.memberships",
		rows,
		[]string{"id"},
		[]string{"user_id", "school_id", "academic_unit_id", "role", "metadata", "is_active", "enrolled_at", "withdrawn_at"},
		true,
	)
}

func seedUserRoles(tx *gorm.DB) error {
	parse := func(v string) time.Time { return mustTimestamp(v) }
	// uuidPtr: helper local — convierte un UUID literal en *uuid.UUID
	// para los campos opcionales (SchoolID, AcademicUnitID, GrantedBy)
	// del struct entities.UserRole.
	uuidPtr := func(v string) *uuid.UUID {
		u := mustUUID(v)
		return &u
	}
	// Mapeo legacy → L0/L4 (post-Fase-6 + PRE-4):
	//   super_admin (L0), school_admin, school_director, school_coordinator,
	//   school_assistant, teacher, assistant_teacher, student, guardian,
	//   observer, readonly_auditor.
	// (PRE-4) El rol `platform_admin` fue eliminado del catálogo; el
	// usuario que antes lo recibía ahora se mapea a `super_admin` (L0),
	// que aporta acceso global equivalente.
	//
	// P1-2: se usa el struct entities.UserRole (en lugar de map[string]any
	// vía upsertMaps) para que el hook BeforeSave registrado en el modelo
	// dispare el cómputo automático de scope_pattern desde
	// school_id/academic_unit_id. GORM no ejecuta hooks de modelo cuando
	// el input es un map.
	rows := []entities.UserRole{
		{ID: mustUUID("cc000000-0000-0000-0000-000000000001"), UserID: mustUUID("00000000-0000-0000-0000-000000000001"), RoleID: mustUUID(layers.L0_ROLE_SUPER_ADMIN_ID), SchoolID: nil, AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-01-01 00:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000002"), UserID: mustUUID("00000000-0000-0000-0000-000000000002"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-01-15 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000003"), UserID: mustUUID("00000000-0000-0000-0000-000000000003"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-01-15 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000004"), UserID: mustUUID("00000000-0000-0000-0000-000000000004"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_COORDINATOR_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-01 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000005"), UserID: mustUUID("00000000-0000-0000-0000-000000000004"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_COORDINATOR_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-02-01 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000006"), UserID: mustUUID("00000000-0000-0000-0000-000000000005"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000007"), UserID: mustUUID("00000000-0000-0000-0000-000000000005"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b3000000-0000-0000-0000-000000000003"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000008"), UserID: mustUUID("00000000-0000-0000-0000-000000000006"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000009"), UserID: mustUUID("00000000-0000-0000-0000-000000000007"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000010"), UserID: mustUUID("00000000-0000-0000-0000-000000000008"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000011"), UserID: mustUUID("00000000-0000-0000-0000-000000000008"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000012"), UserID: mustUUID("00000000-0000-0000-0000-000000000009"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000013"), UserID: mustUUID("00000000-0000-0000-0000-000000000010"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000014"), UserID: mustUUID("00000000-0000-0000-0000-000000000011"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000015"), UserID: mustUUID("00000000-0000-0000-0000-000000000011"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b3000000-0000-0000-0000-000000000003"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000016"), UserID: mustUUID("00000000-0000-0000-0000-000000000012"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000017"), UserID: mustUUID("00000000-0000-0000-0000-000000000013"), RoleID: mustUUID(l4.L4_ROLE_GUARDIAN_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000018"), UserID: mustUUID("00000000-0000-0000-0000-000000000013"), RoleID: mustUUID(l4.L4_ROLE_GUARDIAN_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000019"), UserID: mustUUID("00000000-0000-0000-0000-000000000014"), RoleID: mustUUID(l4.L4_ROLE_GUARDIAN_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		// PRE-4: usuario admin.plataforma@edugo.test re-mapeado de
		// platform_admin (eliminado) a super_admin (L0).
		{ID: mustUUID("cc000000-0000-0000-0000-000000000020"), UserID: mustUUID("00000000-0000-0000-0000-000000000015"), RoleID: mustUUID(layers.L0_ROLE_SUPER_ADMIN_ID), SchoolID: nil, AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-01-20 00:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000021"), UserID: mustUUID("00000000-0000-0000-0000-000000000016"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_DIRECTOR_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-01-20 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000022"), UserID: mustUUID("00000000-0000-0000-0000-000000000017"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_ASSISTANT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-01 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000023"), UserID: mustUUID("00000000-0000-0000-0000-000000000018"), RoleID: mustUUID(l4.L4_ROLE_ASSISTANT_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-15 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000024"), UserID: mustUUID("00000000-0000-0000-0000-000000000019"), RoleID: mustUUID(l4.L4_ROLE_OBSERVER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-20 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000025"), UserID: mustUUID("00000000-0000-0000-0000-000000000019"), RoleID: mustUUID(l4.L4_ROLE_OBSERVER_ID), SchoolID: uuidPtr("b2000000-0000-0000-0000-000000000002"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000003"), GrantedAt: parse("2026-02-20 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000026"), UserID: mustUUID("00000000-0000-0000-0000-000000000020"), RoleID: mustUUID(l4.L4_ROLE_GUARDIAN_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-10 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000027"), UserID: mustUUID("00000000-0000-0000-0000-000000000021"), RoleID: mustUUID(l4.L4_ROLE_READONLY_AUDITOR_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-03-25 08:00:00")},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&rows).Error; err != nil {
		return fmt.Errorf("seedUserRoles: %w", err)
	}
	return nil
}

// seedUserGrants — P4-2: overrides puntuales por usuario en iam.user_grants.
// Demuestra deny > allow (override prohibitivo sobre lectura de notas a un
// student) y allow temporal con expires_at (concede admin.users.create extra a
// un teacher por 30 días). Idempotente vía OnConflict.DoNothing sobre id.
func seedUserGrants(tx *gorm.DB) error {
	grantedBy := mustUUID("00000000-0000-0000-0000-000000000001")
	expiresIn30Days := mustTimestamp("2026-06-11 00:00:00")
	rows := []entities.UserGrant{
		{
			ID:                mustUUID("ee000000-0000-0000-0000-000000000001"),
			UserID:            mustUUID("00000000-0000-0000-0000-000000000008"),
			ScopePattern:      "*",
			PermissionPattern: "academic.grades.read",
			Effect:            "deny",
			GrantedBy:         &grantedBy,
		},
		{
			ID:                mustUUID("ee000000-0000-0000-0000-000000000002"),
			UserID:            mustUUID("00000000-0000-0000-0000-000000000005"),
			ScopePattern:      "*",
			PermissionPattern: "admin.users.create",
			Effect:            "allow",
			ExpiresAt:         &expiresIn30Days,
			GrantedBy:         &grantedBy,
		},
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&rows).Error; err != nil {
		return fmt.Errorf("seedUserGrants: %w", err)
	}
	return nil
}

func seedSubjects(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("dd000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "name": "Matematicas", "code": "MAT-5A", "description": "Matematicas para 5to A", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "name": "Ciencias Naturales", "code": "SCI-5A", "description": "Ciencias Naturales para 5to A", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000003"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "name": "Matematicas", "code": "MAT-5B", "description": "Matematicas para 5to B", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000006"), "name": "Historia", "code": "HIS-6A", "description": "Historia de Chile para 6to A", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "name": "Tecnicas de Pintura", "code": "PINT-GM", "description": "Taller de tecnicas de pintura", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000006"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000011"), "name": "Fundamentos de Escultura", "code": "ESCL-GT", "description": "Taller de fundamentos de escultura", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "name": "English Basics A2", "code": "ENG-A2", "description": "English course for level A2", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "name": "Ciencias Naturales", "code": "SCI-5B", "description": "Ciencias Naturales para 5to B", "is_active": true},
	}
	return upsertMaps(tx, "academic.subjects", rows, []string{"id"}, nil, false)
}

func seedMaterials(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("aa100000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "uploaded_by_teacher_id": mustUUID("00000000-0000-0000-0000-000000000005"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "title": "Introduccion a las Fracciones", "description": "Material introductorio sobre fracciones simples, equivalentes y operaciones basicas.", "subject": "Matematicas", "grade": "5to Basico", "file_url": "s3://edugo-dev/materials/mat001.pdf", "file_type": "application/pdf", "file_size_bytes": 2048000, "status": "ready", "processing_started_at": mustTimestamp("2026-02-10 10:00:00+00"), "processing_completed_at": mustTimestamp("2026-02-10 10:05:32+00"), "is_public": true},
		{"id": mustUUID("aa100000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "uploaded_by_teacher_id": mustUUID("00000000-0000-0000-0000-000000000006"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "title": "El Sistema Solar", "description": "Descripcion de los planetas, el Sol y sus caracteristicas principales.", "subject": "Ciencias Naturales", "grade": "5to Basico", "file_url": "s3://edugo-dev/materials/mat002.pdf", "file_type": "application/pdf", "file_size_bytes": 3145728, "status": "ready", "processing_started_at": mustTimestamp("2026-02-12 11:00:00+00"), "processing_completed_at": mustTimestamp("2026-02-12 11:04:18+00"), "is_public": true},
		{"id": mustUUID("aa100000-0000-0000-0000-000000000003"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "uploaded_by_teacher_id": mustUUID("00000000-0000-0000-0000-000000000006"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000006"), "title": "Historia de Chile: Independencia", "description": "Resumen de los principales procesos de la independencia de Chile.", "subject": "Historia", "grade": "6to Basico", "file_url": "s3://edugo-dev/materials/mat003.pdf", "file_type": "application/pdf", "file_size_bytes": 5242880, "status": "ready", "processing_started_at": mustTimestamp("2026-02-15 14:00:00+00"), "processing_completed_at": mustTimestamp("2026-02-15 14:06:45+00"), "is_public": false},
		{"id": mustUUID("aa100000-0000-0000-0000-000000000004"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "uploaded_by_teacher_id": mustUUID("00000000-0000-0000-0000-000000000007"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "title": "Teoria del Color", "description": "Fundamentos de la teoria del color: colores primarios, secundarios, complementarios.", "subject": "Pintura", "grade": "Modulo Pintura", "file_url": "s3://edugo-dev/materials/mat004.pdf", "file_type": "application/pdf", "file_size_bytes": 1800000, "status": "ready", "processing_started_at": mustTimestamp("2026-02-14 09:00:00+00"), "processing_completed_at": mustTimestamp("2026-02-14 09:03:22+00"), "is_public": true},
		{"id": mustUUID("aa100000-0000-0000-0000-000000000005"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "uploaded_by_teacher_id": mustUUID("00000000-0000-0000-0000-000000000005"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "title": "English Grammar Basics", "description": "Introduction to basic English grammar: articles, pronouns, simple tenses.", "subject": "English", "grade": "Level A2", "file_url": "s3://edugo-dev/materials/mat005.pdf", "file_type": "application/pdf", "file_size_bytes": 1500000, "status": "ready", "processing_started_at": mustTimestamp("2026-02-16 10:00:00+00"), "processing_completed_at": mustTimestamp("2026-02-16 10:04:10+00"), "is_public": true},
	}
	return upsertMaps(tx, "content.materials", rows, []string{"id"}, []string{
		"title", "description", "subject", "grade", "file_url", "file_type", "file_size_bytes",
		"status", "processing_started_at", "processing_completed_at", "is_public",
	}, true)
}

func seedAssessments(tx *gorm.DB) error {
	now := time.Now().UTC()
	rows := []map[string]any{
		{"id": mustUUID("aa200000-0000-0000-0000-000000000001"), "mongo_document_id": "aaaaaa000000000000000001", "source_type": "ai_generated", "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "title": "Examen Fracciones", "description": "Evaluacion sobre operaciones basicas con fracciones: suma, resta y equivalencias.", "questions_count": 5, "pass_threshold": 60.0, "max_attempts": 3, "time_limit_minutes": 30.0, "is_timed": true, "shuffle_questions": true, "show_correct_answers": true, "available_from": now.AddDate(0, 0, -7), "available_until": now.AddDate(0, 0, 30), "status": "published"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000002"), "mongo_document_id": "aaaaaa000000000000000002", "source_type": "ai_generated", "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000006"), "title": "Quiz Ciencias: Sistema Solar", "description": "Quiz rapido sobre los planetas del sistema solar y sus caracteristicas.", "questions_count": 4, "pass_threshold": 50.0, "max_attempts": 2, "time_limit_minutes": nil, "is_timed": false, "shuffle_questions": false, "show_correct_answers": false, "available_from": nil, "available_until": nil, "status": "published"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000003"), "mongo_document_id": "aaaaaa000000000000000003", "source_type": "ai_generated", "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000007"), "title": "Ejercicio Color y Forma", "description": "Ejercicio practico sobre teoria del color y composicion visual.", "questions_count": 3, "pass_threshold": 70.0, "max_attempts": 2, "time_limit_minutes": nil, "is_timed": false, "shuffle_questions": false, "show_correct_answers": true, "available_from": now.AddDate(0, 0, -5), "available_until": now.AddDate(0, 0, 60), "status": "published"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000004"), "mongo_document_id": "aaaaaa000000000000000004", "source_type": "ai_generated", "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "title": "English Grammar Test", "description": "Test on basic English grammar: articles, pronouns, and simple tenses.", "questions_count": 4, "pass_threshold": 60.0, "max_attempts": 2, "time_limit_minutes": 20.0, "is_timed": true, "shuffle_questions": true, "show_correct_answers": true, "available_from": now.AddDate(0, 0, -3), "available_until": now.AddDate(0, 0, 30), "status": "published"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000005"), "mongo_document_id": "aaaaaa000000000000000005", "source_type": "ai_generated", "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000006"), "title": "Evaluacion Historia Chile", "description": "Evaluacion sobre los principales procesos de la independencia de Chile.", "questions_count": 3, "pass_threshold": 70.0, "max_attempts": nil, "time_limit_minutes": 45.0, "is_timed": true, "shuffle_questions": false, "show_correct_answers": true, "available_from": now.AddDate(0, 0, 7), "available_until": nil, "status": "draft"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000006"), "mongo_document_id": "aaaaaa000000000000000006", "source_type": "ai_generated", "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000007"), "title": "Proyecto Final Escultura", "description": "Proyecto final del modulo de escultura: crear una pieza original.", "questions_count": 0, "pass_threshold": 60.0, "max_attempts": 1, "time_limit_minutes": nil, "is_timed": false, "shuffle_questions": false, "show_correct_answers": false, "available_from": nil, "available_until": nil, "status": "draft"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000007"), "mongo_document_id": nil, "source_type": "manual", "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "title": "Evaluacion Manual: Operaciones Basicas", "description": "Evaluacion manual creada por la profesora Maria. Incluye 4 tipos de pregunta: opcion multiple, verdadero/falso, respuesta corta y abierta.", "questions_count": 4, "pass_threshold": 60.0, "max_attempts": 2, "time_limit_minutes": 25.0, "is_timed": true, "shuffle_questions": false, "show_correct_answers": true, "available_from": now.AddDate(0, 0, -2), "available_until": now.AddDate(0, 0, 14), "status": "published"},
		{"id": mustUUID("aa200000-0000-0000-0000-000000000008"), "mongo_document_id": nil, "source_type": "manual", "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "created_by_user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "title": "Quiz Manual: Ciencias Naturales", "description": "Quiz manual en borrador sobre fotosintesis y ecosistemas.", "questions_count": 2, "pass_threshold": 50.0, "max_attempts": 1, "time_limit_minutes": nil, "is_timed": false, "shuffle_questions": false, "show_correct_answers": true, "available_from": nil, "available_until": nil, "status": "draft"},
	}

	return upsertMaps(tx, "assessment.assessment", rows, []string{"id"}, []string{
		"title", "description", "source_type", "questions_count", "pass_threshold", "max_attempts",
		"time_limit_minutes", "is_timed", "shuffle_questions", "show_correct_answers",
		"available_from", "available_until", "status",
	}, true)
}

func seedAssessmentMaterials(tx *gorm.DB) error {
	rows := []map[string]any{
		{"assessment_id": mustUUID("aa200000-0000-0000-0000-000000000001"), "material_id": mustUUID("aa100000-0000-0000-0000-000000000001"), "sort_order": 0},
		{"assessment_id": mustUUID("aa200000-0000-0000-0000-000000000002"), "material_id": mustUUID("aa100000-0000-0000-0000-000000000002"), "sort_order": 0},
		{"assessment_id": mustUUID("aa200000-0000-0000-0000-000000000003"), "material_id": mustUUID("aa100000-0000-0000-0000-000000000004"), "sort_order": 0},
		{"assessment_id": mustUUID("aa200000-0000-0000-0000-000000000004"), "material_id": mustUUID("aa100000-0000-0000-0000-000000000005"), "sort_order": 0},
		{"assessment_id": mustUUID("aa200000-0000-0000-0000-000000000005"), "material_id": mustUUID("aa100000-0000-0000-0000-000000000003"), "sort_order": 0},
	}
	return upsertMaps(tx, "assessment.assessment_materials", rows, []string{"assessment_id", "material_id"}, []string{"sort_order"}, false)
}

func seedQuestions(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("ba000000-0000-0000-0000-000000000001"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "question_text": "Cuanto es 15 + 27?", "question_type": "multiple_choice", "correct_answer": "42", "explanation": "Se suman las unidades (5+7=12, llevamos 1) y las decenas (1+2+1=4). Resultado: 42.", "points": 2.0, "difficulty": "easy", "sort_order": 0},
		{"id": mustUUID("ba000000-0000-0000-0000-000000000002"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "question_text": "El resultado de multiplicar cualquier numero por cero es cero.", "question_type": "true_false", "correct_answer": "Verdadero", "explanation": "La propiedad absorbente de la multiplicacion establece que a x 0 = 0 para todo numero a.", "points": 1.0, "difficulty": "easy", "sort_order": 1},
		{"id": mustUUID("ba000000-0000-0000-0000-000000000003"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "question_text": "Como se llama el resultado de una resta?", "question_type": "short_answer", "correct_answer": "diferencia", "explanation": "El resultado de una resta se denomina diferencia.", "points": 1.5, "difficulty": "medium", "sort_order": 2},
		{"id": mustUUID("ba000000-0000-0000-0000-000000000004"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "question_text": "Explica con un ejemplo de la vida cotidiana donde usarias la multiplicacion.", "question_type": "open_ended", "correct_answer": nil, "explanation": "Respuesta abierta. Se evalua la capacidad de relacionar operaciones matematicas con situaciones reales.", "points": 5.0, "difficulty": "hard", "sort_order": 3},
		{"id": mustUUID("ba000000-0000-0000-0000-000000000005"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000008"), "question_text": "Que gas absorben las plantas durante la fotosintesis?", "question_type": "multiple_choice", "correct_answer": "Dioxido de carbono", "explanation": "Las plantas absorben CO2 y liberan O2 durante la fotosintesis.", "points": 2.0, "difficulty": "easy", "sort_order": 0},
		{"id": mustUUID("ba000000-0000-0000-0000-000000000006"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000008"), "question_text": "Como se llama la capa de la atmosfera donde vivimos?", "question_type": "short_answer", "correct_answer": "troposfera", "explanation": "La troposfera es la capa mas baja de la atmosfera terrestre.", "points": 2.0, "difficulty": "medium", "sort_order": 1},
	}
	return upsertMaps(tx, "assessment.questions", rows, []string{"id"}, nil, false)
}

func seedQuestionOptions(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("bf000000-0000-0000-0000-000000000001"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000001"), "option_text": "32", "sort_order": 0},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000002"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000001"), "option_text": "42", "sort_order": 1},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000003"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000001"), "option_text": "52", "sort_order": 2},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000004"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000001"), "option_text": "41", "sort_order": 3},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000005"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000005"), "option_text": "Oxigeno", "sort_order": 0},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000006"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000005"), "option_text": "Dioxido de carbono", "sort_order": 1},
		{"id": mustUUID("bf000000-0000-0000-0000-000000000007"), "question_id": mustUUID("ba000000-0000-0000-0000-000000000005"), "option_text": "Nitrogeno", "sort_order": 2},
	}
	return upsertMaps(tx, "assessment.question_options", rows, []string{"id"}, nil, false)
}

func seedAssessmentAssignments(tx *gorm.DB) error {
	now := time.Now().UTC()
	rows := []map[string]any{
		{"id": mustUUID("ab000000-0000-0000-0000-000000000001"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "assigned_by": mustUUID("00000000-0000-0000-0000-000000000005"), "assigned_at": now.Add(-24 * time.Hour)},
		{"id": mustUUID("ab000000-0000-0000-0000-000000000002"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "student_id": mustUUID("00000000-0000-0000-0000-000000000009"), "assigned_by": mustUUID("00000000-0000-0000-0000-000000000005"), "assigned_at": now.Add(-24 * time.Hour)},
		{"id": mustUUID("ab000000-0000-0000-0000-000000000003"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000007"), "student_id": mustUUID("00000000-0000-0000-0000-000000000010"), "assigned_by": mustUUID("00000000-0000-0000-0000-000000000005"), "assigned_at": now.Add(-24 * time.Hour)},
	}
	return upsertMaps(tx, "assessment.assessment_assignments", rows, []string{"id"}, nil, false)
}

func seedAssessmentAttempts(tx *gorm.DB) error {
	now := time.Now().UTC()

	start1 := now.Add(-72 * time.Hour)
	start2 := now.Add(-48 * time.Hour)
	start3 := now.Add(-52 * time.Hour)
	start4 := now.Add(-30 * time.Hour)
	start5 := now.Add(-26 * time.Hour)
	start6 := now.Add(-24 * time.Hour)
	start7 := now.Add(-12 * time.Hour)
	start10 := now.Add(-5 * time.Minute)

	rows := []map[string]any{
		{"id": mustUUID("aa300000-0000-0000-0000-000000000001"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000001"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "started_at": start1, "completed_at": start1.Add(25 * time.Minute), "score": 80.0, "max_score": 100.0, "percentage": 80.0, "status": "completed", "time_spent_seconds": 1520, "idempotency_key": "idem_att001_carlos_ass001_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000002"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000001"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "started_at": start2, "completed_at": start2.Add(20 * time.Minute), "score": 92.0, "max_score": 100.0, "percentage": 92.0, "status": "completed", "time_spent_seconds": 1200, "idempotency_key": "idem_att002_carlos_ass001_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000003"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000001"), "student_id": mustUUID("00000000-0000-0000-0000-000000000009"), "started_at": start3, "completed_at": start3.Add(30 * time.Minute), "score": 68.0, "max_score": 100.0, "percentage": 68.0, "status": "completed", "time_spent_seconds": 1800, "idempotency_key": "idem_att003_sofia_ass001_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000004"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000002"), "student_id": mustUUID("00000000-0000-0000-0000-000000000010"), "started_at": start4, "completed_at": start4.Add(15 * time.Minute), "score": 75.0, "max_score": 80.0, "percentage": 93.75, "status": "completed", "time_spent_seconds": 900, "idempotency_key": "idem_att004_diego_ass002_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000005"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000003"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "started_at": start5, "completed_at": start5.Add(20 * time.Minute), "score": 60.0, "max_score": 100.0, "percentage": 60.0, "status": "completed", "time_spent_seconds": 1200, "idempotency_key": "idem_att005_carlos_ass003_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000006"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000003"), "student_id": mustUUID("00000000-0000-0000-0000-000000000012"), "started_at": start6, "completed_at": start6.Add(10 * time.Minute), "score": 90.0, "max_score": 100.0, "percentage": 90.0, "status": "completed", "time_spent_seconds": 600, "idempotency_key": "idem_att006_mateo_ass003_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000007"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000004"), "student_id": mustUUID("00000000-0000-0000-0000-000000000011"), "started_at": start7, "completed_at": start7.Add(15 * time.Minute), "score": 85.0, "max_score": 100.0, "percentage": 85.0, "status": "completed", "time_spent_seconds": 1100, "idempotency_key": "idem_att007_valentina_ass004_v2"},
		{"id": mustUUID("aa300000-0000-0000-0000-000000000010"), "assessment_id": mustUUID("aa200000-0000-0000-0000-000000000001"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "started_at": start10, "completed_at": nil, "score": nil, "max_score": nil, "percentage": nil, "status": "in_progress", "time_spent_seconds": nil, "idempotency_key": "idem_att010_carlos_ass001_inprogress"},
	}

	return upsertMaps(tx, "assessment.assessment_attempt", rows, []string{"idempotency_key"}, []string{
		"score", "max_score", "percentage", "status", "completed_at", "time_spent_seconds",
	}, true)
}

func seedAssessmentAttemptAnswers(tx *gorm.DB) error {
	now := time.Now().UTC()
	attemptStart := now.Add(-72 * time.Hour)

	rows := []map[string]any{
		{"attempt_id": mustUUID("aa300000-0000-0000-0000-000000000001"), "question_index": 0, "student_answer": "1/2", "is_correct": true, "points_earned": 20.0, "max_points": 20.0, "time_spent_seconds": 280, "answered_at": attemptStart.Add(5 * time.Minute)},
		{"attempt_id": mustUUID("aa300000-0000-0000-0000-000000000001"), "question_index": 1, "student_answer": "3/4", "is_correct": true, "points_earned": 20.0, "max_points": 20.0, "time_spent_seconds": 310, "answered_at": attemptStart.Add(10 * time.Minute)},
		{"attempt_id": mustUUID("aa300000-0000-0000-0000-000000000001"), "question_index": 2, "student_answer": "2/6", "is_correct": false, "points_earned": 0.0, "max_points": 20.0, "time_spent_seconds": 420, "answered_at": attemptStart.Add(17 * time.Minute)},
		{"attempt_id": mustUUID("aa300000-0000-0000-0000-000000000001"), "question_index": 3, "student_answer": "2/5", "is_correct": true, "points_earned": 20.0, "max_points": 20.0, "time_spent_seconds": 265, "answered_at": attemptStart.Add(21 * time.Minute)},
		{"attempt_id": mustUUID("aa300000-0000-0000-0000-000000000001"), "question_index": 4, "student_answer": "3/8", "is_correct": true, "points_earned": 20.0, "max_points": 20.0, "time_spent_seconds": 245, "answered_at": attemptStart.Add(25 * time.Minute)},
		{"id": mustUUID("aa000000-0000-0000-0000-000000000020"), "attempt_id": mustUUID("aa300000-0000-0000-0000-000000000010"), "question_index": 0, "student_answer": "A", "is_correct": nil, "points_earned": nil, "max_points": nil, "time_spent_seconds": 5, "answered_at": now.Add(-4 * time.Minute)},
		{"id": mustUUID("aa000000-0000-0000-0000-000000000021"), "attempt_id": mustUUID("aa300000-0000-0000-0000-000000000010"), "question_index": 1, "student_answer": "B", "is_correct": nil, "points_earned": nil, "max_points": nil, "time_spent_seconds": 8, "answered_at": now.Add(-3 * time.Minute)},
	}

	return upsertMaps(tx, "assessment.assessment_attempt_answer", rows, []string{"attempt_id", "question_index"}, []string{
		"student_answer", "is_correct", "points_earned", "max_points", "time_spent_seconds",
	}, true)
}

func seedGuardianRelations(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("ee000000-0000-0000-0000-000000000001"), "guardian_id": mustUUID("00000000-0000-0000-0000-000000000013"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "relationship_type": "parent", "is_primary": true, "is_active": true, "status": "active"},
		{"id": mustUUID("ee000000-0000-0000-0000-000000000002"), "guardian_id": mustUUID("00000000-0000-0000-0000-000000000014"), "student_id": mustUUID("00000000-0000-0000-0000-000000000009"), "relationship_type": "parent", "is_primary": true, "is_active": true, "status": "active"},
		{"id": mustUUID("ee000000-0000-0000-0000-000000000003"), "guardian_id": mustUUID("00000000-0000-0000-0000-000000000014"), "student_id": mustUUID("00000000-0000-0000-0000-000000000010"), "relationship_type": "guardian", "is_primary": false, "is_active": true, "status": "active"},
		{"id": mustUUID("ee000000-0000-0000-0000-000000000004"), "guardian_id": mustUUID("00000000-0000-0000-0000-000000000020"), "student_id": mustUUID("00000000-0000-0000-0000-000000000008"), "relationship_type": "guardian", "is_primary": false, "is_active": true, "status": "pending"},
		{"id": mustUUID("ee000000-0000-0000-0000-000000000005"), "guardian_id": mustUUID("00000000-0000-0000-0000-000000000014"), "student_id": mustUUID("00000000-0000-0000-0000-000000000011"), "relationship_type": "guardian", "is_primary": false, "is_active": true, "status": "pending"},
	}
	return upsertMaps(tx, "academic.guardian_relations", rows, []string{"id"}, nil, false)
}

func seedScreenUserPreferences(tx *gorm.DB) error {
	exists, err := tableExists(tx, "ui_config.screen_user_preferences")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	exists, err = tableExists(tx, "ui_config.screen_instances")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("ff000000-0000-0000-0000-000000000001"), "screen_key": "app-settings", "user_id": mustUUID("00000000-0000-0000-0000-000000000001"), "preferences": mustJSON(`{"dark_mode":true,"language":"es"}`)},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000002"), "screen_key": "app-settings", "user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "preferences": mustJSON(`{"dark_mode":false,"language":"es","push_enabled":true}`)},
	}

	return upsertMaps(tx, "ui_config.screen_user_preferences", rows, []string{"screen_key", "user_id"}, nil, false)
}

func seedSchoolConcepts(tx *gorm.DB) error {
	mappings := []struct {
		SchoolID      uuid.UUID
		ConceptTypeID uuid.UUID
	}{
		{SchoolID: mustUUID("b1000000-0000-0000-0000-000000000001"), ConceptTypeID: mustUUID("c1000000-0000-0000-0000-000000000002")},
		{SchoolID: mustUUID("b2000000-0000-0000-0000-000000000002"), ConceptTypeID: mustUUID("c1000000-0000-0000-0000-000000000005")},
		{SchoolID: mustUUID("b3000000-0000-0000-0000-000000000003"), ConceptTypeID: mustUUID("c1000000-0000-0000-0000-000000000003")},
	}

	rows := make([]map[string]any, 0, 128)
	for _, mapping := range mappings {
		var definitions []struct {
			TermKey   string `gorm:"column:term_key"`
			TermValue string `gorm:"column:term_value"`
			Category  string `gorm:"column:category"`
		}
		if err := tx.Table("academic.concept_definitions").
			Select("term_key", "term_value", "category").
			Where("concept_type_id = ?", mapping.ConceptTypeID).
			Find(&definitions).Error; err != nil {
			return fmt.Errorf("error leyendo concept_definitions: %w", err)
		}
		for _, definition := range definitions {
			rows = append(rows, map[string]any{
				"school_id":  mapping.SchoolID,
				"term_key":   definition.TermKey,
				"term_value": definition.TermValue,
				"category":   definition.Category,
			})
		}
	}

	return upsertMaps(tx, "academic.school_concepts", rows, []string{"school_id", "term_key"}, nil, false)
}

func seedProgress(tx *gorm.DB) error {
	rows := []map[string]any{
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000001"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "progress_percentage": 100.0, "last_position": mustJSON(`{"page":24,"section":"ejercicios-finales"}`), "completed_at": mustTimestamp("2026-03-10 15:30:00+00")},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000002"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "progress_percentage": 65.0, "last_position": mustJSON(`{"page":12,"section":"planetas-exteriores"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000004"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "progress_percentage": 30.0, "last_position": mustJSON(`{"page":5,"section":"colores-primarios"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000001"), "user_id": mustUUID("00000000-0000-0000-0000-000000000009"), "progress_percentage": 80.0, "last_position": mustJSON(`{"page":19,"section":"fracciones-equivalentes"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000002"), "user_id": mustUUID("00000000-0000-0000-0000-000000000009"), "progress_percentage": 45.0, "last_position": mustJSON(`{"page":8,"section":"planetas-interiores"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000002"), "user_id": mustUUID("00000000-0000-0000-0000-000000000010"), "progress_percentage": 90.0, "last_position": mustJSON(`{"page":20,"section":"resumen"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000005"), "user_id": mustUUID("00000000-0000-0000-0000-000000000011"), "progress_percentage": 70.0, "last_position": mustJSON(`{"page":14,"section":"simple-tenses"}`), "completed_at": nil},
		{"material_id": mustUUID("aa100000-0000-0000-0000-000000000004"), "user_id": mustUUID("00000000-0000-0000-0000-000000000012"), "progress_percentage": 55.0, "last_position": mustJSON(`{"page":9,"section":"colores-secundarios"}`), "completed_at": nil},
	}
	return upsertMaps(tx, "content.progress", rows, []string{"material_id", "user_id"}, []string{
		"progress_percentage", "last_position", "completed_at",
	}, true)
}

// seedSubjectOfferings siembra las "sesiones de materia" (subject_offerings)
// y las inscripciones de alumnos (subject_offering_enrollments). Reemplaza al
// antiguo seedMembershipSubjects: el sentido "docente-dicta-materia" pasa al
// teacher_membership_id de la oferta, y "alumno-cursa-materia" pasa a las
// inscripciones. Las ofertas usan el periodo ACTIVO de cada colegio:
//
//	colegio San Ignacio (b1…01) → ff…01 ; CreArte (b2…02) → ff…03 ;
//	Global English (b3…03) → ff…05.
//
// La antigua fila docente=asistente (bb…24) se descarta: el asistente no es un
// docente valido para teacher_membership_id.
func seedSubjectOfferings(tx *gorm.DB) error {
	offeringsExist, err := tableExists(tx, "academic.subject_offerings")
	if err != nil {
		return err
	}
	if !offeringsExist {
		return nil
	}

	// IDs de oferta en rango propio c5000000-… (sesion = materia+unidad+periodo+docente).
	const (
		offMat5A = "c5000000-0000-0000-0000-000000000001" // Matematicas 5to A (San Ignacio)
		offEngA2 = "c5000000-0000-0000-0000-000000000002" // English A2 (Global English)
		offMat5B = "c5000000-0000-0000-0000-000000000003" // Matematicas 5to B (San Ignacio)
		offSci5B = "c5000000-0000-0000-0000-000000000004" // Ciencias 5to B (San Ignacio)
		offHis6A = "c5000000-0000-0000-0000-000000000005" // Historia 6to A (San Ignacio)
		offPint  = "c5000000-0000-0000-0000-000000000006" // Tecnicas de Pintura (CreArte)
		offEscl  = "c5000000-0000-0000-0000-000000000007" // Fundamentos de Escultura (CreArte)
	)

	offerings := []map[string]any{
		{"id": mustUUID(offMat5A), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offEngA2), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000007"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000005"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000009"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offMat5B), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offSci5B), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000008"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offHis6A), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000004"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000006"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000011"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offPint), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000005"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000009"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000003"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000012"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offEscl), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000006"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000011"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000003"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000013"), "is_active": true, "metadata": mustJSON(`{}`)},
	}

	if err := upsertMaps(
		tx,
		"academic.subject_offerings",
		offerings,
		[]string{"id"},
		[]string{"teacher_membership_id", "is_active", "metadata"},
		true,
	); err != nil {
		return err
	}

	enrollmentsExist, err := tableExists(tx, "academic.subject_offering_enrollments")
	if err != nil {
		return err
	}
	if !enrollmentsExist {
		return nil
	}

	// Inscripciones de los alumnos demo a las sesiones de su unidad.
	// enrolled_at se setea explicito: la columna es NOT NULL sin default de BD y
	// el autoCreateTime de la entity solo aplica en inserts via struct (no via
	// map, como hace upsertMaps). Mismo valor que el enrolled_at de las
	// membresias de alumnos: inicio del periodo academico activo (2026-03-01 UTC).
	enrollments := []map[string]any{
		{"offering_id": mustUUID(offMat5A), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Carlos (5to A)
		{"offering_id": mustUUID(offMat5A), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Sofia (5to A)
		{"offering_id": mustUUID(offMat5B), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B)
		{"offering_id": mustUUID(offSci5B), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B)
		{"offering_id": mustUUID(offHis6A), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000005"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Valentina (6to A)
		{"offering_id": mustUUID(offEngA2), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000006"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Valentina (Global English)
		{"offering_id": mustUUID(offPint), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000007"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")},  // Mateo (CreArte pintura)
	}

	return upsertMaps(tx, "academic.subject_offering_enrollments", enrollments, []string{"offering_id", "student_membership_id"}, nil, false)
}

func seedAcademicPeriods(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.academic_periods")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("ff000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "Primer Semestre 2026", "code": "S1-2026", "type": "semester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-07-15"), "is_active": true, "academic_year": 2026, "sort_order": 1},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "Segundo Semestre 2026", "code": "S2-2026", "type": "semester", "start_date": mustDate("2026-08-01"), "end_date": mustDate("2026-12-15"), "is_active": false, "academic_year": 2026, "sort_order": 2},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000003"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Primer Trimestre 2026", "code": "T1-2026", "type": "trimester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-05-31"), "is_active": true, "academic_year": 2026, "sort_order": 1},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "name": "Segundo Trimestre 2026", "code": "T2-2026", "type": "trimester", "start_date": mustDate("2026-06-01"), "end_date": mustDate("2026-08-31"), "is_active": false, "academic_year": 2026, "sort_order": 2},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Bimestre 1", "code": "B1-2026", "type": "bimester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-04-30"), "is_active": true, "academic_year": 2026, "sort_order": 1},
	}

	return upsertMaps(tx, "academic.academic_periods", rows, []string{"id"}, nil, false)
}

func seedGrades(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.grades")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("a0000000-0000-0000-0000-000000000001"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 6.5, "grade_letter": "B+", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "notes": "Buen rendimiento", "finalized_at": mustTimestamp("2026-03-20 10:00:00+00")},
		{"id": mustUUID("a0000000-0000-0000-0000-000000000002"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000002"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 5.8, "grade_letter": "B", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "notes": nil, "finalized_at": nil},
		{"id": mustUUID("a0000000-0000-0000-0000-000000000003"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 7.0, "grade_letter": "A-", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "notes": "Excelente", "finalized_at": mustTimestamp("2026-03-20 10:00:00+00")},
		{"id": mustUUID("a0000000-0000-0000-0000-000000000004"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000003"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 5.0, "grade_letter": "C+", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "notes": "Debe mejorar", "finalized_at": nil},
	}

	return upsertMaps(tx, "academic.grades", rows, []string{"id"}, nil, false)
}

func seedAttendance(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.attendance")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("a1000000-0000-0000-0000-000000000001"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000002"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000003"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "late", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000004"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000005"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-19"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000006"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-19"), "status": "absent", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000007"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000003"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000008"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000003"), "date": mustDate("2026-03-18"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
	}

	return upsertMaps(tx, "academic.attendance", rows, []string{"id"}, nil, false)
}

func seedSchedules(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.schedules")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("a2000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "day_of_week": 1, "start_time": mustTimestamp("2026-03-03 08:00:00+00"), "end_time": mustTimestamp("2026-03-03 09:30:00+00"), "room": "Sala 101", "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "is_active": true},
		{"id": mustUUID("a2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000002"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "day_of_week": 3, "start_time": mustTimestamp("2026-03-05 10:00:00+00"), "end_time": mustTimestamp("2026-03-05 11:30:00+00"), "room": "Sala 102", "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "is_active": true},
		{"id": mustUUID("a2000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000003"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "day_of_week": 2, "start_time": mustTimestamp("2026-03-04 08:00:00+00"), "end_time": mustTimestamp("2026-03-04 09:30:00+00"), "room": "Sala 103", "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "is_active": true},
	}

	return upsertMaps(tx, "academic.schedules", rows, []string{"id"}, nil, false)
}

func seedAnnouncements(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.announcements")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{
			"id":               mustUUID("a3000000-0000-0000-0000-000000000001"),
			"school_id":        mustUUID("b1000000-0000-0000-0000-000000000001"),
			"academic_unit_id": nil,
			"author_id":        mustUUID("00000000-0000-0000-0000-000000000002"),
			"title":            "Reunion de Apoderados Marzo 2026",
			"body":             "Se convoca a reunion general de apoderados el viernes 28 de marzo a las 18:00 hrs en el auditorio principal. Favor confirmar asistencia.",
			"scope":            "school",
			"target_roles":     pq.StringArray{"guardian", "teacher"},
			"is_pinned":        true,
			"published_at":     mustTimestamp("2026-03-20 10:00:00+00"),
			"expires_at":       mustTimestamp("2026-03-28 23:59:59+00"),
		},
		{
			"id":               mustUUID("a3000000-0000-0000-0000-000000000002"),
			"school_id":        mustUUID("b1000000-0000-0000-0000-000000000001"),
			"academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"),
			"author_id":        mustUUID("00000000-0000-0000-0000-000000000005"),
			"title":            "Prueba de Matematicas - 5to A",
			"body":             "La prueba de Matematicas del primer semestre se realizara el lunes 31 de marzo. Estudiar capitulos 1 al 4.",
			"scope":            "unit",
			"target_roles":     pq.StringArray{"student", "guardian"},
			"is_pinned":        false,
			"published_at":     mustTimestamp("2026-03-21 14:00:00+00"),
			"expires_at":       mustTimestamp("2026-03-31 23:59:59+00"),
		},
	}

	return upsertMaps(tx, "academic.announcements", rows, []string{"id"}, nil, false)
}

func seedCalendarEvents(tx *gorm.DB) error {
	exists, err := tableExists(tx, "academic.calendar_events")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{"id": mustUUID("a4000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "title": "Semana Santa - Sin Clases", "description": "Receso por Semana Santa. Se retoman clases el lunes 6 de abril.", "event_type": "holiday", "start_date": mustDate("2026-04-02"), "end_date": mustDate("2026-04-05"), "is_all_day": true, "created_by": mustUUID("00000000-0000-0000-0000-000000000002")},
		{"id": mustUUID("a4000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "title": "Examenes Primer Semestre - Matematicas", "description": "Examenes de matematicas para todos los cursos de 5to y 6to.", "event_type": "exam", "start_date": mustDate("2026-03-31"), "end_date": mustDate("2026-03-31"), "is_all_day": true, "created_by": mustUUID("00000000-0000-0000-0000-000000000002")},
	}

	return upsertMaps(tx, "academic.calendar_events", rows, []string{"id"}, nil, false)
}

func seedCourses(tx *gorm.DB) error {
	// content.courses puede no existir si la migración de la API de learning no se ha aplicado.
	// En ese caso, se omite silenciosamente.
	exists, err := tableExists(tx, "content.courses")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rows := []map[string]any{
		{
			"id":          mustUUID("c0000000-0000-0000-0000-000000000001"),
			"unit_id":     mustUUID("ac000000-0000-0000-0000-000000000003"),
			"name":        "Curso de Matematicas 5to A",
			"description": "Curso de matematicas para 5to A - seed de desarrollo para tests de integración",
			"status":      "active",
		},
		{
			"id":          mustUUID("c0000000-0000-0000-0000-000000000002"),
			"unit_id":     mustUUID("ac000000-0000-0000-0000-000000000004"),
			"name":        "Curso de Ciencias 5to B",
			"description": "Curso de ciencias para 5to B - seed de desarrollo para tests de integración",
			"status":      "active",
		},
	}
	return upsertMaps(tx, "content.courses", rows, []string{"id"}, []string{"name", "description", "status"}, true)
}

func upsertMaps(tx *gorm.DB, table string, rows []map[string]any, conflictColumns, updateColumns []string, touchUpdatedAt bool) error {
	if len(rows) == 0 {
		return nil
	}

	columns, err := tableColumnSet(tx, table)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, row := range rows {
		if columns["id"] {
			if _, exists := row["id"]; !exists {
				row["id"] = uuid.New()
			}
		}
		if columns["created_at"] {
			if _, exists := row["created_at"]; !exists {
				row["created_at"] = now
			}
		}
		if columns["updated_at"] {
			if _, exists := row["updated_at"]; !exists {
				row["updated_at"] = now
			}
		}
	}

	conflict := clause.OnConflict{
		Columns: toColumns(conflictColumns...),
	}
	if len(updateColumns) == 0 && !touchUpdatedAt {
		conflict.DoNothing = true
	} else {
		assignments := make(map[string]any, len(updateColumns)+1)
		for _, column := range updateColumns {
			assignments[column] = gorm.Expr("EXCLUDED." + column)
		}
		if touchUpdatedAt {
			assignments["updated_at"] = gorm.Expr("NOW()")
		}
		conflict.DoUpdates = clause.Assignments(assignments)
	}

	return tx.Table(table).Clauses(conflict).CreateInBatches(rows, 200).Error
}

func tableColumnSet(tx *gorm.DB, table string) (map[string]bool, error) {
	schema, name, err := splitQualifiedTable(table)
	if err != nil {
		return nil, err
	}

	type col struct {
		ColumnName string `gorm:"column:column_name"`
	}
	var cols []col
	if err := tx.Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
	`, schema, name).Scan(&cols).Error; err != nil {
		return nil, fmt.Errorf("error leyendo columnas de %s: %w", table, err)
	}

	result := make(map[string]bool, len(cols))
	for _, c := range cols {
		result[c.ColumnName] = true
	}
	return result, nil
}

func splitQualifiedTable(table string) (schema, name string, err error) {
	parts := strings.Split(table, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("tabla inválida (esperada schema.table): %s", table)
	}
	return parts[0], parts[1], nil
}

func toColumns(columns ...string) []clause.Column {
	cols := make([]clause.Column, 0, len(columns))
	for _, column := range columns {
		cols = append(cols, clause.Column{Name: column})
	}
	return cols
}

func truncateTable(tx *gorm.DB, table string) error {
	return tx.Exec("TRUNCATE TABLE " + table + " CASCADE").Error
}

func truncateIfExists(tx *gorm.DB, table string) error {
	exists, err := tableExists(tx, table)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return truncateTable(tx, table)
}

func tableExists(tx *gorm.DB, table string) (bool, error) {
	var exists bool
	if err := tx.Raw("SELECT to_regclass(?) IS NOT NULL", table).Scan(&exists).Error; err != nil {
		return false, fmt.Errorf("error verificando tabla %s: %w", table, err)
	}
	return exists, nil
}

func mustUUID(v string) uuid.UUID {
	parsed, err := uuid.Parse(v)
	if err != nil {
		panic(err)
	}
	return parsed
}

func mustJSON(v string) json.RawMessage {
	return json.RawMessage(v)
}

func mustDate(v string) time.Time {
	parsed, err := time.Parse("2006-01-02", v)
	if err != nil {
		panic(err)
	}
	return parsed.UTC()
}

func mustTimestamp(v string) time.Time {
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05-07",
		"2006-01-02 15:04:05-0700",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, v); err == nil {
			return parsed.UTC()
		}
	}
	panic(fmt.Sprintf("timestamp inválido: %s", v))
}
