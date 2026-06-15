// Package base es el MUNDO DE DATOS POR DEFECTO de EduGo (MP-09 / F0, 2026-06-14).
// Reubicado verbatim desde seeds/demo: mismos UUIDs/emails @edugo.test, para que dev
// y los tests de integracion compartan UNA sola fuente de la verdad. El paquete
// seeds/demo se elimina en F2 (repoint demo.ApplyDemo -> base.Apply).
package base

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
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

func Apply(gdb *gorm.DB) error {
	return gdb.Transaction(func(tx *gorm.DB) error {
		if err := truncateDevelopmentData(tx); err != nil {
			return err
		}
		if err := seedSchools(tx); err != nil {
			return err
		}
		if err := seedSchoolInvitationRoles(tx); err != nil {
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
		if err := seedAnnouncements(tx); err != nil {
			return err
		}
		// Plan 024 F1 (ADR 0026): tejido de representante (guardian). Va al final
		// porque depende de users/schools/units/memberships ya sembrados. DEC-R-B:
		// el guardián NO lleva membership, solo user_role guardián scoped +
		// guardian_relations.
		if err := seedGuardianTejido(tx); err != nil {
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
		"academic.announcements",
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

// seedSchoolInvitationRoles siembra las equivalencias tipo→rol IAM por defecto
// (academic.school_invitation_roles) para las 2 escuelas de base.
//
// base siembra sus escuelas con upsertMaps raw (sin pasar por common.SeedSchool),
// así que debe invocar explícitamente el helper compartido — de lo contrario sus
// escuelas nacerían sin equivalencias y el flujo de aprobación de invitaciones
// (que admin-go edita por UI) quedaría sin mapeo. MP-09 F4: el sistema (L4) dejó
// de sembrar la escuela demo, así que este es el único punto donde las escuelas
// de base obtienen sus equivalencias.
//
// PRECONDICIÓN: academic.invitation_types ya está sembrado por el system seed
// (l4.ApplyInvitationTypes, que corre ANTES que cualquier playground). El helper
// es idempotente (id derivado SHA1(school:type) + ON CONFLICT DO NOTHING).
func seedSchoolInvitationRoles(tx *gorm.DB) error {
	schools := []uuid.UUID{
		mustUUID("b1000000-0000-0000-0000-000000000001"), // Colegio San Ignacio
		mustUUID("b3000000-0000-0000-0000-000000000003"), // Academia Global English
	}
	for _, schoolID := range schools {
		if err := l4.SeedDefaultSchoolInvitationRoles(tx, schoolID); err != nil {
			return fmt.Errorf("seedSchoolInvitationRoles: %w", err)
		}
	}
	return nil
}

func seedAcademicUnits(tx *gorm.DB) error {
	rows := []map[string]any{
		{"id": mustUUID("ac000000-0000-0000-0000-000000000001"), "parent_unit_id": nil, "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "Colegio San Ignacio", "code": "CSI-ROOT", "type": "school", "description": "Unidad raiz del Colegio San Ignacio", "level": "secondary", "academic_year": 0, "metadata": mustJSON(`{"is_root":true}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000012"), "parent_unit_id": nil, "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Academia Global English", "code": "AGE-ROOT", "type": "school", "description": "Unidad raiz de la Academia Global English", "level": "language", "academic_year": 0, "metadata": mustJSON(`{"is_root":true}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000002"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000001"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to Basico", "code": "GRADE-05", "type": "grade", "description": "Quinto ano de educacion basica, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000013"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000012"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Level A2", "code": "LVL-A2", "type": "grade", "description": "Elementary level A2", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"cefr_level":"A2"}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000003"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to A", "code": "5A", "type": "class", "description": "Seccion A del 5to Basico, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"section":"A","grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000004"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "name": "5to B", "code": "5B", "type": "class", "description": "Seccion B del 5to Basico, 2026", "level": "secondary", "academic_year": 2026, "metadata": mustJSON(`{"section":"B","grade_number":5}`), "is_active": true},
		{"id": mustUUID("ac000000-0000-0000-0000-000000000014"), "parent_unit_id": mustUUID("ac000000-0000-0000-0000-000000000013"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "name": "Class Monday", "code": "CLS-MON", "type": "class", "description": "Monday class - Level A2", "level": "language", "academic_year": 2026, "metadata": mustJSON(`{"day":"monday"}`), "is_active": true},
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
		{"id": mustUUID("00000000-0000-0000-0000-000000000004"), "email": "coord.academico@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Lucia", "last_name": "Fernandez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000005"), "email": "prof.martinez@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Maria", "last_name": "Martinez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000006"), "email": "prof.gonzalez@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Pedro", "last_name": "Gonzalez", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000008"), "email": "est.carlos@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Carlos", "last_name": "Mendoza", "is_active": true},
		// est.sofia (…0009): alumna de S1/5A, espejo de est.diego pero en la
		// sección A (compañera de carlos). MP-09-base-logica la incluye; base la
		// sembraba sin ella, plan 024 F1 la reconcilia para colgar un representante.
		{"id": mustUUID("00000000-0000-0000-0000-000000000009"), "email": "est.sofia@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Sofia", "last_name": "Rojas", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000010"), "email": "est.diego@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Diego", "last_name": "Vargas", "is_active": true},
		{"id": mustUUID("00000000-0000-0000-0000-000000000018"), "email": "asist.prof@edugo.test", "password_hash": defaultPasswordHash, "first_name": "Andres", "last_name": "Gomez", "is_active": true},
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

	// MP-08 (esquema 3.64.0): academic.memberships ya no tiene columna `role`;
	// la membresía referencia su tipo vía invitation_type_id (FK a
	// academic.invitation_types). Resolvemos cada key del catálogo a su id UNA
	// vez (data-driven, sin hardcodear UUIDs) antes de construir las filas.
	// PRECONDICIÓN: el catálogo invitation_types lo siembra L4
	// (l4.ApplyInvitationTypes) y el runner aplica system (L0..L4) ANTES de demo
	// (ver migrate.go / cmd/seed / cmd/runner), así que aquí ya existe.
	invType := func(key string) uuid.UUID {
		id, err := catalog.ResolveInvitationTypeID(tx, key)
		if err != nil {
			panic(fmt.Sprintf("seedMemberships: %v", err))
		}
		return id
	}
	student := invType("student")
	teacher := invType("teacher")
	admin := invType("admin")
	coordinator := invType("coordinator")
	assistant := invType("assistant")

	rows := []map[string]any{
		{"id": mustUUID("bb000000-0000-0000-0000-000000000001"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "invitation_type_id": student, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000002"), "user_id": mustUUID("00000000-0000-0000-0000-000000000008"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "invitation_type_id": student, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		// Sofia (…0009) — membership de alumna en S1/5A (ac…03), misma sección que
		// carlos. Espejo de diego (bb…04, 5B) pero en la sección A.
		{"id": mustUUID("bb000000-0000-0000-0000-000000000003"), "user_id": mustUUID("00000000-0000-0000-0000-000000000009"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "invitation_type_id": student, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000004"), "user_id": mustUUID("00000000-0000-0000-0000-000000000010"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "invitation_type_id": student, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-01 08:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000008"), "user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "invitation_type_id": teacher, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000009"), "user_id": mustUUID("00000000-0000-0000-0000-000000000005"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000014"), "invitation_type_id": teacher, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000010"), "user_id": mustUUID("00000000-0000-0000-0000-000000000006"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000004"), "invitation_type_id": teacher, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000014"), "user_id": mustUUID("00000000-0000-0000-0000-000000000002"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "invitation_type_id": admin, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-01-15 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000016"), "user_id": mustUUID("00000000-0000-0000-0000-000000000004"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "invitation_type_id": coordinator, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-10 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000024"), "user_id": mustUUID("00000000-0000-0000-0000-000000000018"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": mustUUID("ac000000-0000-0000-0000-000000000003"), "invitation_type_id": assistant, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-02-15 09:00:00+00")},
		{"id": mustUUID("bb000000-0000-0000-0000-000000000028"), "user_id": mustUUID("00000000-0000-0000-0000-000000000021"), "school_id": mustUUID("b1000000-0000-0000-0000-000000000001"), "academic_unit_id": nil, "invitation_type_id": admin, "metadata": mustJSON(`{}`), "is_active": true, "enrolled_at": parse("2026-03-25 08:00:00+00")},
	}

	return upsertMaps(
		tx,
		"academic.memberships",
		rows,
		[]string{"id"},
		[]string{"user_id", "school_id", "academic_unit_id", "invitation_type_id", "metadata", "is_active", "enrolled_at", "withdrawn_at"},
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
		{ID: mustUUID("cc000000-0000-0000-0000-000000000004"), UserID: mustUUID("00000000-0000-0000-0000-000000000004"), RoleID: mustUUID(l4.L4_ROLE_SCHOOL_COORDINATOR_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-01 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000006"), UserID: mustUUID("00000000-0000-0000-0000-000000000005"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000007"), UserID: mustUUID("00000000-0000-0000-0000-000000000005"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b3000000-0000-0000-0000-000000000003"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000008"), UserID: mustUUID("00000000-0000-0000-0000-000000000006"), RoleID: mustUUID(l4.L4_ROLE_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-10 09:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000010"), UserID: mustUUID("00000000-0000-0000-0000-000000000008"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000011"), UserID: mustUUID("00000000-0000-0000-0000-000000000008"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b3000000-0000-0000-0000-000000000003"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000001"), GrantedAt: parse("2026-03-01 08:00:00")},
		// Sofia (…0009) — rol student scoped a S1 (espejo de diego cc…13).
		{ID: mustUUID("cc000000-0000-0000-0000-000000000012"), UserID: mustUUID("00000000-0000-0000-0000-000000000009"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000013"), UserID: mustUUID("00000000-0000-0000-0000-000000000010"), RoleID: mustUUID(l4.L4_ROLE_STUDENT_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-03-01 08:00:00")},
		// PRE-4: usuario admin.plataforma@edugo.test re-mapeado de
		// platform_admin (eliminado) a super_admin (L0).
		{ID: mustUUID("cc000000-0000-0000-0000-000000000023"), UserID: mustUUID("00000000-0000-0000-0000-000000000018"), RoleID: mustUUID(l4.L4_ROLE_ASSISTANT_TEACHER_ID), SchoolID: uuidPtr("b1000000-0000-0000-0000-000000000001"), AcademicUnitID: nil, IsActive: true, GrantedBy: uuidPtr("00000000-0000-0000-0000-000000000002"), GrantedAt: parse("2026-02-15 09:00:00")},
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
// un teacher). El expires_at es relativo a la fecha de aplicación (un año en el
// futuro) para que el grant siga ACTIVO en pruebas sin importar cuándo se
// siembre. Idempotente vía OnConflict.DoNothing sobre id.
func seedUserGrants(tx *gorm.DB) error {
	grantedBy := mustUUID("00000000-0000-0000-0000-000000000001")
	expiresInOneYear := time.Now().UTC().AddDate(1, 0, 0)
	rows := []entities.UserGrant{
		{
			ID:                mustUUID("ee000000-0000-0000-0000-000000000001"),
			UserID:            mustUUID("00000000-0000-0000-0000-000000000008"),
			PermissionPattern: "academic.grades.read",
			Effect:            "deny",
			GrantedBy:         &grantedBy,
		},
		{
			ID:                mustUUID("ee000000-0000-0000-0000-000000000002"),
			UserID:            mustUUID("00000000-0000-0000-0000-000000000005"),
			PermissionPattern: "admin.users.create",
			Effect:            "allow",
			ExpiresAt:         &expiresInOneYear,
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
		{"id": mustUUID("dd000000-0000-0000-0000-000000000007"), "school_id": mustUUID("b3000000-0000-0000-0000-000000000003"), "academic_unit_id": nil, "name": "English Basics A2", "code": "ENG-A2", "description": "English course for level A2", "is_active": true},
	}
	return upsertMaps(tx, "academic.subjects", rows, []string{"id"}, nil, false)
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
	}

	return upsertMaps(tx, "ui_config.screen_user_preferences", rows, []string{"screen_key", "user_id"}, nil, false)
}

func seedSchoolConcepts(tx *gorm.DB) error {
	mappings := []struct {
		SchoolID      uuid.UUID
		ConceptTypeID uuid.UUID
	}{
		{SchoolID: mustUUID("b1000000-0000-0000-0000-000000000001"), ConceptTypeID: mustUUID("c1000000-0000-0000-0000-000000000002")},
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
	// Periodos: offMat5A/offMat5B/offSci5B → ff…01; offEngA2 → ff…05.
	enrollments := []map[string]any{
		{"offering_id": mustUUID(offMat5A), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Carlos (5to A, Matematicas)
		{"offering_id": mustUUID(offMat5A), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Sofia (5to A, Matematicas)
		{"offering_id": mustUUID(offMat5B), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B, Matematicas)
		{"offering_id": mustUUID(offSci5B), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000002"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Diego (5to B, Ciencias)
		{"offering_id": mustUUID(offEngA2), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000007"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000005"), "student_membership_id": mustUUID("bb000000-0000-0000-0000-000000000002"), "enrolled_at": mustTimestamp("2026-03-01 08:00:00+00")}, // Carlos (Global English, multi-escuela)
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
		// Sofia (bb…03) — nota de Matematicas 5A (dd…01, period ff…01), teacher
		// bb…08 (martinez, docente de Mate5A). Espejo de la nota de diego (a0…04)
		// pero en la sección A; finalized como nota real (paralelo a a0…01).
		{"id": mustUUID("a0000000-0000-0000-0000-000000000003"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "period_id": mustUUID("ff000000-0000-0000-0000-000000000001"), "grade_value": 6.0, "grade_letter": "B+", "teacher_id": mustUUID("bb000000-0000-0000-0000-000000000008"), "notes": "Progreso constante", "finalized_at": mustTimestamp("2026-03-20 10:00:00+00")},
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
		{"id": mustUUID("a1000000-0000-0000-0000-000000000003"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "late", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000005"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000001"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-19"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		// Sofia (bb…03) — asistencia de Mate5A (dd…01), misma sección/materia que
		// carlos (recorded_by martinez …05). Usa los IDs libres del pool a1…02/04/06
		// (carlos toma 01/03/05, diego 07/08). Mismas 3 fechas que carlos.
		{"id": mustUUID("a1000000-0000-0000-0000-000000000002"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000004"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000006"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000003"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-19"), "status": "late", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000005")},
		// Repuntadas a dd…01 (Matematicas escuela) tras colapsar la duplicada 5B
		// (ADR 0016). attendance_unique=(membership,subject,date): membership
		// bb…04 no tiene otra asistencia en dd…01 → sin colisión.
		{"id": mustUUID("a1000000-0000-0000-0000-000000000007"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-17"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
		{"id": mustUUID("a1000000-0000-0000-0000-000000000008"), "membership_id": mustUUID("bb000000-0000-0000-0000-000000000004"), "subject_id": mustUUID("dd000000-0000-0000-0000-000000000001"), "date": mustDate("2026-03-18"), "status": "present", "recorded_by": mustUUID("00000000-0000-0000-0000-000000000006")},
	}

	return upsertMaps(tx, "academic.attendance", rows, []string{"id"}, nil, false)
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

// seedGuardianTejido siembra el modelo de representante (guardian) del plan 024
// F1 / ADR 0026 sobre el dataset base ya existente. DEC-R-B: el guardián es un
// usuario con rol guardián SCOPED a la(s) escuela(s) de sus acudidos, SIN
// membership (no es alumno ni staff de la escuela). El vínculo con cada alumno
// vive en academic.guardian_relations, con school_id en el índice único: el
// mismo guardián puede colgar del mismo alumno en dos escuelas distintas (caso
// repB↔carlos en S1 y S3).
//
// Se invoca al final de Apply: depende de users/schools/units/memberships ya
// sembrados. Idempotente: auth.users / iam.user_roles / academic.guardian_relations
// se truncan en truncateDevelopmentData; school_guardian_policy usa OnConflict
// DoNothing por id.
func seedGuardianTejido(tx *gorm.DB) error {
	const (
		// Escuelas de los acudidos.
		schoolS1 = "b1000000-0000-0000-0000-000000000001" // Colegio San Ignacio
		schoolS3 = "b3000000-0000-0000-0000-000000000003" // Academia Global English

		// Alumnos (user IDs de auth.users sembrados arriba).
		studSofia  = "00000000-0000-0000-0000-000000000009"
		studCarlos = "00000000-0000-0000-0000-000000000008"
		studDiego  = "00000000-0000-0000-0000-000000000010"

		// Representantes (user IDs NUEVOS, libres en base).
		guardianA = "00000000-0000-0000-0000-000000000011" // Laura Mendoza
		guardianB = "00000000-0000-0000-0000-000000000012" // Miguel Castro
	)

	// 1) Usuarios representante (password 12345678 vía helper común con bcrypt).
	guardians := []common.UserSpec{
		{ID: mustUUID(guardianA), Email: "tutor.mendoza@edugo.test", Password: "12345678", FirstName: "Laura", LastName: "Mendoza"},
		{ID: mustUUID(guardianB), Email: "tutor.castro@edugo.test", Password: "12345678", FirstName: "Miguel", LastName: "Castro"},
	}
	for _, spec := range guardians {
		if err := common.SeedUser(tx, spec); err != nil {
			return fmt.Errorf("seedGuardianTejido (user): %w", err)
		}
	}

	// 2) user_roles guardián SCOPED por escuela (sin membership; DEC-R-B). IDs
	// cc… nuevos y libres. El BeforeSave de UserRole calcula scope_pattern desde
	// school_id (school:<id>). repB lleva DOS roles (S1 y S3): tiene acudidos en
	// ambas escuelas.
	uuidPtr := func(v string) *uuid.UUID { u := mustUUID(v); return &u }
	grantedBy := mustUUID("00000000-0000-0000-0000-000000000001") // super
	grantedAt := mustTimestamp("2026-03-01 08:00:00")
	guardianRoleID := mustUUID(l4.L4_ROLE_GUARDIAN_ID)
	guardianRoles := []entities.UserRole{
		{ID: mustUUID("cc000000-0000-0000-0000-000000000017"), UserID: mustUUID(guardianA), RoleID: guardianRoleID, SchoolID: uuidPtr(schoolS1), AcademicUnitID: nil, IsActive: true, GrantedBy: &grantedBy, GrantedAt: grantedAt},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000018"), UserID: mustUUID(guardianB), RoleID: guardianRoleID, SchoolID: uuidPtr(schoolS1), AcademicUnitID: nil, IsActive: true, GrantedBy: &grantedBy, GrantedAt: grantedAt},
		{ID: mustUUID("cc000000-0000-0000-0000-000000000019"), UserID: mustUUID(guardianB), RoleID: guardianRoleID, SchoolID: uuidPtr(schoolS3), AcademicUnitID: nil, IsActive: true, GrantedBy: &grantedBy, GrantedAt: grantedAt},
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&guardianRoles).Error; err != nil {
		return fmt.Errorf("seedGuardianTejido (user_roles): %w", err)
	}

	// 3) guardian_relations (link a nivel escuela: AcademicUnitID nil porque la
	// política por defecto es link_scope=school). IsPrimary=true, Status="active".
	// repB↔carlos aparece DOS veces (S1 y S3): mismo guardián+alumno, otra
	// escuela — habilitado por el índice único ampliado a (guardian,student,school).
	relations := []common.GuardianRelationSpec{
		{ID: mustUUID("9a000000-0000-0000-0000-000000000001"), GuardianID: mustUUID(guardianA), StudentID: mustUUID(studSofia), SchoolID: mustUUID(schoolS1), IsPrimary: true},  // Laura ↔ Sofia (S1)
		{ID: mustUUID("9a000000-0000-0000-0000-000000000002"), GuardianID: mustUUID(guardianB), StudentID: mustUUID(studCarlos), SchoolID: mustUUID(schoolS1), IsPrimary: true}, // Miguel ↔ Carlos (S1)
		{ID: mustUUID("9a000000-0000-0000-0000-000000000003"), GuardianID: mustUUID(guardianB), StudentID: mustUUID(studCarlos), SchoolID: mustUUID(schoolS3), IsPrimary: true}, // Miguel ↔ Carlos (S3, multi-escuela)
		{ID: mustUUID("9a000000-0000-0000-0000-000000000004"), GuardianID: mustUUID(guardianB), StudentID: mustUUID(studDiego), SchoolID: mustUUID(schoolS1), IsPrimary: true},  // Miguel ↔ Diego (S1)
	}
	for _, spec := range relations {
		if err := common.SeedGuardianRelation(tx, spec); err != nil {
			return fmt.Errorf("seedGuardianTejido (relation): %w", err)
		}
	}

	// 4) school_guardian_policy: solo S3 aparta del default (S1 no lleva fila =
	// usa los defaults del esquema). S3 invita en la inscripción y gatea
	// activación con aprobación de cualquier representante.
	if err := common.SeedSchoolGuardianPolicy(tx, common.SchoolGuardianPolicySpec{
		ID:              mustUUID("9b000000-0000-0000-0000-000000000001"),
		SchoolID:        mustUUID(schoolS3),
		AcademicUnitID:  nil,
		InvitationMode:  "on_enrollment",
		GatesActivation: true,
		GatingApprover:  "any",
		LinkScope:       "school",
	}); err != nil {
		return fmt.Errorf("seedGuardianTejido (policy): %w", err)
	}

	return nil
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
