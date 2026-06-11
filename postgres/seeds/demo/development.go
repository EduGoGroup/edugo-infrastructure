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
	// v4 (ADR 0016): materias migradas a scope de ESCUELA (academic_unit_id =
	// NULL) y deduplicadas por (school_id, name); referencias en offerings/
	// grades/attendance/schedules repuntadas al id sobreviviente.
	// v5: los academic_periods ganan academic_unit_id (período atado a la
	// unidad raíz de cada colegio); el activo es exclusivo por (school, unit).
	SeedVersion         = "development-gorm-v5"
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
		// N4 F1 (plan 015 / ADR 0019): los datos demo de evaluación/contenido
		// (materials, assessments, assessment_materials, questions,
		// question_options, assignments, attempts, attempt_answers, progress)
		// se ELIMINARON: sembraban el contrato viejo (created_by_user_id,
		// subject/grade texto-libre, student_id→auth.users), ahora demolido. Su
		// reconstrucción sobre el esquema nuevo (FKs a memberships/subjects/
		// subject_offerings) es DATA, no esquema → F2/F4.
		if err := seedGuardianRelations(tx); err != nil {
			return err
		}
		if err := seedScreenUserPreferences(tx); err != nil {
			return err
		}
		if err := seedSchoolConcepts(tx); err != nil {
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
		return nil
	})
}

func truncateDevelopmentData(tx *gorm.DB) error {
	guarded := []string{
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

	// N4 F1: las tablas analíticas viejas (assessment.attempt_analytics /
	// assessment_stats) se eliminaron del esquema; ya no se truncan. Las tablas
	// de evaluación/contenido se renombraron (assessment.question(_option),
	// assessment.assessment_material) y ya no llevan datos demo (reconstrucción
	// en F2). Se truncan igual por idempotencia ante datos previos.
	required := []string{
		"assessment.attempt_review",
		"assessment.assessment_attempt_answer",
		"assessment.assessment_attempt",
		"assessment.assessment_assignment",
		"assessment.assessment_material",
		"assessment.question_option",
		"assessment.question",
		"assessment.assessment",
		"content.progress",
		"content.material_file",
		"content.materials",
		"academic.memberships",
		"iam.user_grants",
		"iam.user_roles",
		"academic.academic_units",
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

// seedSubjects siembra el catálogo de materias de ESCUELA (ADR 0016):
// academic_unit_id = NULL en todas. Una materia es catálogo de la escuela
// (reutilizable en cualquier unidad vía sus sesiones), y cumple
// UNIQUE(school_id, name) — no puede repetirse el mismo nombre dentro de una
// escuela. Por eso se eliminaron las materias duplicadas del modelo viejo
// materia=unidad: en San Ignacio (b1…01) "Matematicas" y "Ciencias Naturales"
// existían dos veces (5A y 5B). Quedan una sola de cada (dd…01 Matematicas,
// dd…02 Ciencias Naturales); las antiguas dd…03 (Matematicas 5B) y dd…08
// (Ciencias 5B) se colapsan en ellas y sus referencias (offerings/grades/
// attendance) se repuntan al id sobreviviente. La diferencia de unidad que antes
// distinguía 5A de 5B la aporta ahora la sesión (subject_offering.academic_unit_id),
// no la materia.
func seedSubjects(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("dd000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "name": "Matematicas", "code": "MAT-5A", "description": "Matematicas (5to A/B)", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "name": "Ciencias Naturales", "code": "SCI-5A", "description": "Ciencias Naturales (5to A/B)", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "name": "Historia", "code": "HIS-6A", "description": "Historia de Chile para 6to A", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": nil, "name": "Tecnicas de Pintura", "code": "PINT-GM", "description": "Taller de tecnicas de pintura", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000006"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": nil, "name": "Fundamentos de Escultura", "code": "ESCL-GT", "description": "Taller de fundamentos de escultura", "is_active": true},
		{"id": mustUUID("dd000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": nil, "name": "English Basics A2", "code": "ENG-A2", "description": "English course for level A2", "is_active": true},
	}
	return upsertMaps(tx, "academic.subjects", rows, []string{"id"}, nil, false)
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
		// offMat5B/offSci5B repuntan a la materia ESCUELA sobreviviente (dd…01/
		// dd…02) tras colapsar las duplicadas 5B (ADR 0016). La unidad 5to B
		// (ac…04) la lleva la sesión, NO la materia; la natural key de la oferta
		// (school+subject+unit+section+period) no colisiona con la sesión 5A
		// (offMat5A/dd…01) porque la unidad difiere (ac…03 vs ac…04).
		{"id": mustUUID(offMat5B), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "is_active": true, "metadata": mustJSON(`{}`)},
		{"id": mustUUID(offSci5B), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "is_active": true, "metadata": mustJSON(`{}`)},
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
	// subject_id y period_id se copian de la oferta correspondiente (mismo
	// archivo, slice offerings de arriba): son denormalizados e inmutables y
	// respaldan el invariante una-oferta-por-materia-por-período (bug 0036).
	// Periodos: offMat5A/offMat5B/offSci5B/offHis6A → ff…01; offEngA2 → ff…05;
	// offPint → ff…03.
	enrollments := []map[string]any{
		{"offering_id": mustUUID(offMat5A), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Carlos (5to A, Matematicas)
		{"offering_id": mustUUID(offMat5A), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Sofia (5to A, Matematicas)
		{"offering_id": mustUUID(offMat5B), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B, Matematicas)
		{"offering_id": mustUUID(offSci5B), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000002"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B, Ciencias)
		{"offering_id": mustUUID(offHis6A), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000004"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000005"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Valentina (6to A, Historia)
		{"offering_id": mustUUID(offEngA2), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000007"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000005"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000006"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Valentina (Global English)
		{"offering_id": mustUUID(offPint), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000005"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000003"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000007"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")},  // Mateo (CreArte pintura)
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
		{"id": mustUUID("ff000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000001"), "name": "Primer Semestre 2026", "code": "S1-2026", "type": "semester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-07-15"), "is_active": true, "academic_year": 2026, "sort_order": 1},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000001"), "name": "Segundo Semestre 2026", "code": "S2-2026", "type": "semester", "start_date": mustDate("2026-08-01"), "end_date": mustDate("2026-12-15"), "is_active": false, "academic_year": 2026, "sort_order": 2},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000003"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000007"), "name": "Primer Trimestre 2026", "code": "T1-2026", "type": "trimester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-05-31"), "is_active": true, "academic_year": 2026, "sort_order": 1},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b2000000-0000-0000-0000-000000000002"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000007"), "name": "Segundo Trimestre 2026", "code": "T2-2026", "type": "trimester", "start_date": mustDate("2026-06-01"), "end_date": mustDate("2026-08-31"), "is_active": false, "academic_year": 2026, "sort_order": 2},
		{"id": mustUUID("ff000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000012"), "name": "Bimestre 1", "code": "B1-2026", "type": "bimester", "start_date": mustDate("2026-03-01"), "end_date": mustDate("2026-04-30"), "is_active": true, "academic_year": 2026, "sort_order": 1},
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
		// Repuntada a dd…01 (Matematicas escuela) tras colapsar la duplicada 5B
		// (ADR 0016). grades_unique=(membership,subject,period): membership bb…04
		// es único en esta materia/periodo → sin colisión.
		{"id": mustUUID("a0000000-0000-0000-0000-000000000004"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 5.0, "grade_letter": "C+", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "notes": "Debe mejorar", "finalized_at": nil},
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
		// Repuntadas a dd…01 (Matematicas escuela) tras colapsar la duplicada 5B
		// (ADR 0016). attendance_unique=(membership,subject,date): membership
		// bb…04 no tiene otra asistencia en dd…01 → sin colisión.
		{"id": mustUUID("a1000000-0000-0000-0000-000000000007"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000008"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
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
		// Repuntada a dd…01 (Matematicas escuela) tras colapsar la duplicada 5B
		// (ADR 0016). schedules no tiene unique natural → repunte directo seguro;
		// la unidad 5to B (ac…04) sigue en la propia fila de horario.
		{"id": mustUUID("a2000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "teacher_membership_id": mustUUID("bb000000-0000-0000-0000-000000000010"), "day_of_week": 2, "start_time": mustTimestamp("2026-03-04 08:00:00+00"), "end_time": mustTimestamp("2026-03-04 09:30:00+00"), "room": "Sala 103", "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "is_active": true},
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
