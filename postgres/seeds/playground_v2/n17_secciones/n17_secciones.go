// Package n17_secciones es un playground de la línea v2 para validar el
// modelo N1.7 de SESIONES de materia con SECCIONES (plan 010 / ADR 0009).
//
// Foco de la validación:
//   - Secciones A/B: la misma materia (Matemáticas) se dicta en dos sesiones
//     distintas, una con section_label "A" y otra con "B".
//   - Un docente con DOS sesiones: el docente X dicta Mat-A y Mat-B.
//   - Un segundo docente (Y) que dicta una tercera sesión (Lenguaje-A).
//   - Un alumno inscrito en DOS sesiones (alumno A1/A2 en Mat-A y Len-A).
//   - Un alumno SIN inscribir (Alumno Libre) para ejercitar el flujo de
//     inscripción por lote en E2E.
//
// A diferencia de n1_inscripcion (que NO setea section_label), aquí cada
// subject_offering SÍ lleva section_label — ese es el punto de la validación.
//
// Como todo v2, asume que el sistema completo (L0..L4) ya corrió: reusa los
// roles L4 school_admin/teacher/student (ver system/l4/roles_permissions.go)
// para que los usuarios tengan contexto de login real — NO inventa roles ni
// permisos. El login resuelve active_context desde academic.memberships +
// iam.user_roles, así que cada usuario sembrado lleva ambas filas.
//
// Convive con los demás playgrounds sin colisionar: rango UUID propio
// 67000000-... (distinto del 66000000 de n1_inscripcion) y emails con sufijo
// @n17.edugo.local. La escuela es distinta, así que el índice único parcial
// de período activo (por school_id WHERE is_active) no colisiona con
// n1_inscripcion.
//
// Lo que siembra:
//  1. academic.schools           — 1 colegio "Colegio N1.7 Secciones".
//  2. academic.academic_units    — 1 unidad académica "Grado N1.7".
//  3. academic.subjects          — 2 materias de ESCUELA (AcademicUnitID=NULL,
//     ADR 0016): Matemáticas, Lenguaje. Nombres distintos → cumplen
//     UNIQUE(school_id, name).
//  4. academic.academic_periods  — 1 período ACTIVO (is_active=true).
//  5. auth.users                 — 1 admin + 2 docentes + 5 alumnos, password "12345678".
//  6. iam.user_roles             — admin→school_admin L4, docentes→teacher L4, alumnos→student L4.
//  7. academic.memberships       — admin con alcance COLEGIO (AcademicUnitID=NULL); los demás en la unidad.
//  8. academic.subject_offerings — 3 sesiones CON section_label:
//     - Mat-A (Matemáticas, sección "A", docente X).
//     - Mat-B (Matemáticas, sección "B", docente X).
//     - Len-A (Lenguaje, sección "A", docente Y).
//  9. academic.subject_offering_enrollments — inscripciones:
//     - Mat-A → alumno A1, alumno A2.
//     - Mat-B → alumno B1, alumno B2.
//     - Len-A → alumno A1, alumno A2 (alumno en 2 sesiones).
//     - alumno Libre → SIN filas (sin inscribir).
//
// Credenciales (todas password "12345678"):
//
//	admin-n17@n17.edugo.local        school_admin — director (alcance colegio)
//	docente-x-n17@n17.edugo.local    teacher      — dicta Mat-A y Mat-B
//	docente-y-n17@n17.edugo.local    teacher      — dicta Len-A
//	alumno-a1-n17@n17.edugo.local    student      — Mat-A + Len-A
//	alumno-a2-n17@n17.edugo.local    student      — Mat-A + Len-A
//	alumno-b1-n17@n17.edugo.local    student      — Mat-B
//	alumno-b2-n17@n17.edugo.local    student      — Mat-B
//	alumno-libre-n17@n17.edugo.local student      — SIN materias (sin inscribir)
//
// Idempotente: OnConflict DoNothing por id (o clave natural compuesta en
// subject_offering_enrollments) en todas las inserciones.
package n17_secciones

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales de los usuarios sembrados.
	TeacherXEmail     = "docente-x-n17@n17.edugo.local"
	TeacherYEmail     = "docente-y-n17@n17.edugo.local"
	StudentA1Email    = "alumno-a1-n17@n17.edugo.local"
	StudentA2Email    = "alumno-a2-n17@n17.edugo.local"
	StudentB1Email    = "alumno-b1-n17@n17.edugo.local"
	StudentB2Email    = "alumno-b2-n17@n17.edugo.local"
	StudentLibreEmail = "alumno-libre-n17@n17.edugo.local"
	AdminEmail        = "admin-n17@n17.edugo.local"
	Password          = "12345678"

	// Rango UUID 67000000-...: reservado para el playground n17_secciones.
	schoolID = "67000000-0000-0000-0000-000000000001"
	unitID   = "67000000-0000-0000-0000-000000000002"

	subjectMathID = "67000000-0000-0000-0000-000000000003"
	subjectLangID = "67000000-0000-0000-0000-000000000004"

	periodID = "67000000-0000-0000-0000-000000000006"

	teacherXUserID     = "67000000-0000-0000-0000-000000000010"
	teacherYUserID     = "67000000-0000-0000-0000-000000000011"
	studentA1UserID    = "67000000-0000-0000-0000-000000000012"
	studentA2UserID    = "67000000-0000-0000-0000-000000000013"
	studentB1UserID    = "67000000-0000-0000-0000-000000000014"
	studentB2UserID    = "67000000-0000-0000-0000-000000000015"
	studentLibreUserID = "67000000-0000-0000-0000-000000000016"
	adminUserID        = "67000000-0000-0000-0000-000000000017"

	teacherXMembID     = "67000000-0000-0000-0000-000000000020"
	teacherYMembID     = "67000000-0000-0000-0000-000000000021"
	studentA1MembID    = "67000000-0000-0000-0000-000000000022"
	studentA2MembID    = "67000000-0000-0000-0000-000000000023"
	studentB1MembID    = "67000000-0000-0000-0000-000000000024"
	studentB2MembID    = "67000000-0000-0000-0000-000000000025"
	studentLibreMembID = "67000000-0000-0000-0000-000000000026"
	adminMembID        = "67000000-0000-0000-0000-000000000027"

	// Sesiones de materia (subject_offerings) CON section_label.
	offeringMatAID = "67000000-0000-0000-0000-000000000030"
	offeringMatBID = "67000000-0000-0000-0000-000000000031"
	offeringLenAID = "67000000-0000-0000-0000-000000000032"

	schoolCode = "N17-SECCIONES"
	schoolName = "Colegio N1.7 Secciones"
	unitCode   = "N17-UNIT"
	unitName   = "Grado N1.7"

	academicYear = 2026
)

// Apply siembra el playground n17_secciones. Asume que L0..L4 corrieron (los
// roles teacher/student y el esquema academic completo ya existen). Orden:
// school → unit → subjects → period → users → user_roles → memberships →
// subject_offerings → subject_offering_enrollments. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: academic_unit: %w", err)
	}

	if err := upsertSubject(tx, subjectMathID, "Matemáticas", "MAT"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: subject_math: %w", err)
	}
	if err := upsertSubject(tx, subjectLangID, "Lenguaje", "LEN"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: subject_lang: %w", err)
	}

	if err := upsertActivePeriod(tx); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: academic_period: %w", err)
	}

	if err := upsertUser(tx, teacherXUserID, TeacherXEmail, "Docente", "X"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_user: %w", err)
	}
	if err := upsertUser(tx, teacherYUserID, TeacherYEmail, "Docente", "Y"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_user: %w", err)
	}
	if err := upsertUser(tx, studentA1UserID, StudentA1Email, "Alumno", "A1"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_user: %w", err)
	}
	if err := upsertUser(tx, studentA2UserID, StudentA2Email, "Alumno", "A2"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_user: %w", err)
	}
	if err := upsertUser(tx, studentB1UserID, StudentB1Email, "Alumno", "B1"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_user: %w", err)
	}
	if err := upsertUser(tx, studentB2UserID, StudentB2Email, "Alumno", "B2"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_user: %w", err)
	}
	if err := upsertUser(tx, studentLibreUserID, StudentLibreEmail, "Alumno", "Libre"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_user: %w", err)
	}
	if err := upsertUser(tx, adminUserID, AdminEmail, "Admin", "N17"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := upsertUserRole(tx, teacherXUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_user_role: %w", err)
	}
	if err := upsertUserRole(tx, teacherYUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentA1UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentA2UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentB1UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentB2UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentLibreUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_user_role: %w", err)
	}
	if err := upsertUserRole(tx, adminUserID, l4.L4_ROLE_SCHOOL_ADMIN_ID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_user_role: %w", err)
	}

	// Membresías: docentes y alumnos con alcance UNIDAD en la misma unidad.
	if err := upsertMembership(tx, teacherXMembID, teacherXUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_membership: %w", err)
	}
	if err := upsertMembership(tx, teacherYMembID, teacherYUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_membership: %w", err)
	}
	if err := upsertMembership(tx, studentA1MembID, studentA1UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_membership: %w", err)
	}
	if err := upsertMembership(tx, studentA2MembID, studentA2UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_membership: %w", err)
	}
	if err := upsertMembership(tx, studentB1MembID, studentB1UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_membership: %w", err)
	}
	if err := upsertMembership(tx, studentB2MembID, studentB2UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_membership: %w", err)
	}
	if err := upsertMembership(tx, studentLibreMembID, studentLibreUserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_membership: %w", err)
	}
	// Membresía del admin con alcance COLEGIO (AcademicUnitID = NULL): el form
	// memberships-form exige contexto de colegio en el JWT del actor.
	if err := upsertSchoolMembership(tx, adminMembID, adminUserID, "admin"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_membership: %w", err)
	}

	// Sesiones de materia (subject_offerings) CON section_label:
	//  - Mat-A: Matemáticas, sección "A", docente X.
	//  - Mat-B: Matemáticas, sección "B", docente X (mismo docente, 2 sesiones).
	//  - Len-A: Lenguaje, sección "A", docente Y.
	if err := upsertOffering(tx, offeringMatAID, subjectMathID, teacherXMembID, "A"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_mat_a: %w", err)
	}
	if err := upsertOffering(tx, offeringMatBID, subjectMathID, teacherXMembID, "B"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_mat_b: %w", err)
	}
	if err := upsertOffering(tx, offeringLenAID, subjectLangID, teacherYMembID, "A"); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_len_a: %w", err)
	}

	// Inscripciones (subject_offering_enrollments):
	//  - Mat-A: alumno A1, alumno A2.
	//  - Mat-B: alumno B1, alumno B2.
	//  - Len-A: alumno A1, alumno A2 (alumno en 2 sesiones).
	//  - alumno Libre: SIN inscribir (no se crea ninguna fila).
	if err := upsertEnrollment(tx, offeringMatAID, studentA1MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_a_a1: %w", err)
	}
	if err := upsertEnrollment(tx, offeringMatAID, studentA2MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_a_a2: %w", err)
	}
	if err := upsertEnrollment(tx, offeringMatBID, studentB1MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_b_b1: %w", err)
	}
	if err := upsertEnrollment(tx, offeringMatBID, studentB2MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_b_b2: %w", err)
	}
	if err := upsertEnrollment(tx, offeringLenAID, studentA1MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_len_a_a1: %w", err)
	}
	if err := upsertEnrollment(tx, offeringLenAID, studentA2MembID); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_len_a_a2: %w", err)
	}

	return nil
}

func upsertSchool(tx *gorm.DB) error {
	id, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	s := entities.School{
		ID:               id,
		Name:             schoolName,
		Code:             schoolCode,
		Country:          "Chile",
		SubscriptionTier: "basic",
		MaxTeachers:      0,
		MaxStudents:      0,
		IsActive:         true,
		Metadata:         json.RawMessage(`{}`),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&s).Error
}

func upsertAcademicUnit(tx *gorm.DB) error {
	id, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	u := entities.AcademicUnit{
		ID:           id,
		SchoolID:     sid,
		Name:         unitName,
		Code:         unitCode,
		Type:         "class",
		AcademicYear: academicYear,
		Metadata:     json.RawMessage(`{}`),
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

// upsertSubject siembra una materia con scope de ESCUELA (ADR 0016):
// AcademicUnitID = nil. La materia es catálogo de la escuela; la sección y la
// unidad las aporta la sesión (subject_offering). Los 2 nombres del playground
// (Matemáticas, Lenguaje) son distintos → cumplen UNIQUE(school_id, name); la
// misma materia "Matemáticas" se dicta en 2 secciones vía 2 offerings (Mat-A/
// Mat-B), NO vía 2 filas de subjects.
func upsertSubject(tx *gorm.DB, idStr, name, code string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	c := code
	subj := entities.Subject{
		ID:             id,
		SchoolID:       sid,
		AcademicUnitID: nil,
		Name:           name,
		Code:           &c,
		IsActive:       true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&subj).Error
}

// upsertActivePeriod siembra un período académico ACTIVO (is_active=true)
// para el colegio. Hay un índice único parcial por school_id WHERE is_active,
// así que sólo puede haber uno activo por colegio (cumplido: único período).
// Como es OTRA escuela distinta a n1_inscripcion, no colisiona.
func upsertActivePeriod(tx *gorm.DB) error {
	id, err := uuid.Parse(periodID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	code := "N17-2026-S1"
	p := entities.AcademicPeriod{
		ID:           id,
		SchoolID:     sid,
		Name:         "Semestre 1 2026",
		Code:         &code,
		Type:         "semester",
		StartDate:    time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
		IsActive:     true,
		AcademicYear: academicYear,
		SortOrder:    1,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&p).Error
}

func upsertUser(tx *gorm.DB, idStr, email, first, last string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}
	u := entities.User{
		ID:           id,
		Email:        email,
		PasswordHash: string(hash),
		FirstName:    first,
		LastName:     last,
		IsActive:     true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&u).Error
}

func upsertUserRole(tx *gorm.DB, userIDStr, roleIDStr string) error {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	rid, err := uuid.Parse(roleIDStr)
	if err != nil {
		return err
	}
	derived := uuid.NewSHA1(uuid.NameSpaceOID, []byte(uid.String()+":"+rid.String()))
	ur := entities.UserRole{
		ID:        derived,
		UserID:    uid,
		RoleID:    rid,
		IsActive:  true,
		GrantedAt: time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&ur).Error
}

func upsertMembership(tx *gorm.DB, idStr, userIDStr, roleKind string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:             id,
		UserID:         uid,
		SchoolID:       sid,
		AcademicUnitID: &auid,
		Role:           roleKind,
		Metadata:       json.RawMessage(`{}`),
		IsActive:       true,
		EnrolledAt:     time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}

// upsertSchoolMembership crea una membresía con alcance COLEGIO (no UNIDAD):
// AcademicUnitID = nil. Es lo que necesita el school_admin para que su JWT
// lleve contexto de colegio y pueda abrir memberships-form. Idempotente por id.
func upsertSchoolMembership(tx *gorm.DB, idStr, userIDStr, roleKind string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	m := entities.Membership{
		ID:             id,
		UserID:         uid,
		SchoolID:       sid,
		AcademicUnitID: nil,
		Role:           roleKind,
		Metadata:       json.RawMessage(`{}`),
		IsActive:       true,
		EnrolledAt:     time.Now().UTC(),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&m).Error
}

// upsertOffering crea una sesión de materia (subject_offering) para la materia
// dada en la unidad y período del playground, CON section_label (a diferencia
// de n1_inscripcion). teacherMembIDStr puede ser "" (sesión sin docente →
// teacher_membership_id NULL). sectionLabel se setea siempre que no sea ""
// (forma parte del índice único natural uq_subject_offerings_natural).
// Idempotente por id.
func upsertOffering(tx *gorm.DB, idStr, subjectIDStr, teacherMembIDStr, sectionLabel string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	subjID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitID)
	if err != nil {
		return err
	}
	pid, err := uuid.Parse(periodID)
	if err != nil {
		return err
	}

	var teacherMembID *uuid.UUID
	if teacherMembIDStr != "" {
		tmid, err := uuid.Parse(teacherMembIDStr)
		if err != nil {
			return err
		}
		teacherMembID = &tmid
	}

	var section *string
	if sectionLabel != "" {
		s := sectionLabel
		section = &s
	}

	o := entities.SubjectOffering{
		ID:                  id,
		SchoolID:            sid,
		SubjectID:           subjID,
		AcademicUnitID:      auid,
		SectionLabel:        section,
		PeriodID:            pid,
		TeacherMembershipID: teacherMembID,
		IsActive:            true,
		Metadata:            json.RawMessage(`{}`),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&o).Error
}

// upsertEnrollment inscribe al alumno (membership) en una sesión de materia
// (subject_offering_enrollment). La PK es compuesta (offering_id,
// student_membership_id); OnConflict sobre ambas → idempotente.
func upsertEnrollment(tx *gorm.DB, offeringIDStr, studentMembIDStr string) error {
	oid, err := uuid.Parse(offeringIDStr)
	if err != nil {
		return err
	}
	smid, err := uuid.Parse(studentMembIDStr)
	if err != nil {
		return err
	}
	e := entities.SubjectOfferingEnrollment{
		OfferingID:          oid,
		StudentMembershipID: smid,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "offering_id"}, {Name: "student_membership_id"}},
		DoNothing: true,
	}).Create(&e).Error
}
