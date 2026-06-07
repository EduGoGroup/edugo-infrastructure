// Package n4_evaluacion es un playground de la línea v2 para validar el cierre
// de N4: evaluación + notas con procedencia, en una escuela con perfil de notas
// DETALLADO (plan 015 / ADR 0019-0020).
//
// Foco de la validación:
//   - La escuela siembra GradeProfile = "detailed" (a diferencia de los demás
//     playgrounds, que quedan en 'basic'). Esto es lo que habilita que el alumno
//     vea el desglose por componentes (grade_item) en "Mis notas". Sin "detailed"
//     no se ejercita el modo de notas detallado del cierre N4.
//   - Una sesión de "Ciencias Naturales" con un docente y 3 alumnos inscritos,
//     sobre la cual el E2E arma un examen (tipo "sistema solar"), lo toma cada
//     alumno y verifica nota con procedencia (auto_scored/manual/auto_llm) y
//     componentes.
//   - Una segunda materia "Matemáticas" (sin oferta por defecto) queda como
//     catálogo para no acoplar el foco; el foco es Ciencias.
//
// Como todo v2, asume que el sistema completo (L0..L4) ya corrió: reusa los
// roles L4 school_admin/teacher/student (ver system/l4/roles_permissions.go)
// para que los usuarios tengan contexto de login real — NO inventa roles ni
// permisos. El login resuelve active_context desde academic.memberships +
// iam.user_roles, así que cada usuario sembrado lleva ambas filas.
//
// Convive con los demás playgrounds sin colisionar: rango UUID propio
// 6a000000-... (LIBRE; tomados: 10000000 n0n1, 64000000 onboarding, 66000000
// n1_inscripcion, 67000000 n17_secciones, 68000000 multi_unidad, 69000000 y
// c4000000 n0n1_escuelas) y emails con sufijo @n4.edugo.local. La escuela es
// distinta, así que el índice único parcial de período activo (por school_id
// WHERE is_active) no colisiona con los demás playgrounds.
//
// Lo que siembra:
//  1. academic.schools           — 1 colegio "Colegio N4 Evaluación", GradeProfile "detailed".
//  2. academic.academic_units    — 1 unidad académica "Grado N4".
//  3. academic.subjects          — 2 materias de ESCUELA (AcademicUnitID=NULL,
//     ADR 0016): Ciencias Naturales (CIE), Matemáticas (MAT).
//  4. academic.academic_periods  — 1 período ACTIVO (is_active=true).
//  5. auth.users                 — 1 admin + 1 docente + 3 alumnos, password "12345678".
//  6. iam.user_roles             — admin→school_admin L4, docente→teacher L4, alumnos→student L4.
//  7. academic.memberships       — admin con alcance COLEGIO (AcademicUnitID=NULL); los demás en la unidad.
//  8. academic.subject_offerings — 1 sesión: Ciencias Naturales, sección "A", docente.
//  9. academic.subject_offering_enrollments — Ana, Bruno, Caro inscritos en Ciencias.
//
// Credenciales (todas password "12345678"):
//
//	admin-n4@n4.edugo.local       school_admin — director (alcance colegio)
//	profe-ciencias@n4.edugo.local teacher      — dicta Ciencias Naturales (sección A)
//	ana@n4.edugo.local            student      — Ciencias
//	bruno@n4.edugo.local          student      — Ciencias
//	caro@n4.edugo.local           student      — Ciencias
//
// Idempotente: OnConflict DoNothing por id (o clave natural compuesta en
// subject_offering_enrollments) en todas las inserciones.
package n4_evaluacion

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
	AdminEmail        = "admin-n4@n4.edugo.local"
	TeacherEmail      = "profe-ciencias@n4.edugo.local"
	StudentAnaEmail   = "ana@n4.edugo.local"
	StudentBrunoEmail = "bruno@n4.edugo.local"
	StudentCaroEmail  = "caro@n4.edugo.local"
	Password          = "12345678"

	// Rango UUID 6a000000-...: reservado para el playground n4_evaluacion.
	schoolID = "6a000000-0000-0000-0000-000000000001"
	unitID   = "6a000000-0000-0000-0000-000000000002"

	subjectScienceID = "6a000000-0000-0000-0000-000000000003"
	subjectMathID    = "6a000000-0000-0000-0000-000000000004"

	periodID = "6a000000-0000-0000-0000-000000000006"

	adminUserID        = "6a000000-0000-0000-0000-000000000010"
	teacherUserID      = "6a000000-0000-0000-0000-000000000011"
	studentAnaUserID   = "6a000000-0000-0000-0000-000000000012"
	studentBrunoUserID = "6a000000-0000-0000-0000-000000000013"
	studentCaroUserID  = "6a000000-0000-0000-0000-000000000014"

	adminMembID        = "6a000000-0000-0000-0000-000000000020"
	teacherMembID      = "6a000000-0000-0000-0000-000000000021"
	studentAnaMembID   = "6a000000-0000-0000-0000-000000000022"
	studentBrunoMembID = "6a000000-0000-0000-0000-000000000023"
	studentCaroMembID  = "6a000000-0000-0000-0000-000000000024"

	// Sesión de materia (subject_offering) de Ciencias, sección "A".
	offeringScienceID = "6a000000-0000-0000-0000-000000000030"

	schoolCode = "N4-EVALUACION"
	schoolName = "Colegio N4 Evaluación"
	unitCode   = "N4-UNIT"
	unitName   = "Grado N4"

	academicYear = 2026
)

// Apply siembra el playground n4_evaluacion. Asume que L0..L4 corrieron (los
// roles teacher/student y el esquema academic completo ya existen). Orden:
// school → unit → subjects → period → users → user_roles → memberships →
// subject_offerings → subject_offering_enrollments. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: academic_unit: %w", err)
	}

	if err := upsertSubject(tx, subjectScienceID, "Ciencias Naturales", "CIE"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: subject_science: %w", err)
	}
	if err := upsertSubject(tx, subjectMathID, "Matemáticas", "MAT"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: subject_math: %w", err)
	}

	if err := upsertActivePeriod(tx); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: academic_period: %w", err)
	}

	if err := upsertUser(tx, adminUserID, AdminEmail, "Admin", "N4"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: admin_user: %w", err)
	}
	if err := upsertUser(tx, teacherUserID, TeacherEmail, "Profe", "Ciencias"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: teacher_user: %w", err)
	}
	if err := upsertUser(tx, studentAnaUserID, StudentAnaEmail, "Ana", "Alumna"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_ana_user: %w", err)
	}
	if err := upsertUser(tx, studentBrunoUserID, StudentBrunoEmail, "Bruno", "Alumno"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_bruno_user: %w", err)
	}
	if err := upsertUser(tx, studentCaroUserID, StudentCaroEmail, "Caro", "Alumna"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_caro_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := upsertUserRole(tx, adminUserID, l4.L4_ROLE_SCHOOL_ADMIN_ID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: admin_user_role: %w", err)
	}
	if err := upsertUserRole(tx, teacherUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: teacher_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentAnaUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_ana_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentBrunoUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_bruno_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentCaroUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_caro_user_role: %w", err)
	}

	// Membresías: docente y alumnos con alcance UNIDAD en la misma unidad.
	if err := upsertMembership(tx, teacherMembID, teacherUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: teacher_membership: %w", err)
	}
	if err := upsertMembership(tx, studentAnaMembID, studentAnaUserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_ana_membership: %w", err)
	}
	if err := upsertMembership(tx, studentBrunoMembID, studentBrunoUserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_bruno_membership: %w", err)
	}
	if err := upsertMembership(tx, studentCaroMembID, studentCaroUserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: student_caro_membership: %w", err)
	}
	// Membresía del admin con alcance COLEGIO (AcademicUnitID = NULL): el form
	// memberships-form exige contexto de colegio en el JWT del actor.
	if err := upsertSchoolMembership(tx, adminMembID, adminUserID, "admin"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: admin_membership: %w", err)
	}

	// Sesión de materia (subject_offering): Ciencias Naturales, sección "A",
	// dictada por el docente.
	if err := upsertOffering(tx, offeringScienceID, subjectScienceID, teacherMembID, "A"); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: offering_science: %w", err)
	}

	// Inscripciones (subject_offering_enrollments): Ana, Bruno, Caro en Ciencias.
	if err := upsertEnrollment(tx, offeringScienceID, subjectScienceID, studentAnaMembID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: enroll_science_ana: %w", err)
	}
	if err := upsertEnrollment(tx, offeringScienceID, subjectScienceID, studentBrunoMembID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: enroll_science_bruno: %w", err)
	}
	if err := upsertEnrollment(tx, offeringScienceID, subjectScienceID, studentCaroMembID); err != nil {
		return fmt.Errorf("playground_v2/n4_evaluacion: enroll_science_caro: %w", err)
	}

	return nil
}

// upsertSchool siembra el colegio del playground con GradeProfile "detailed"
// (N4 / ADR 0020): es lo que habilita el desglose por componentes en "Mis notas"
// del alumno. A diferencia de los demás playgrounds (que dejan el default
// 'basic'), aquí DEBE ser "detailed" para ejercitar el cierre N4.
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
		GradeProfile:     "detailed",
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
// (Ciencias Naturales, Matemáticas) son distintos → cumplen UNIQUE(school_id, name).
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
// Como es OTRA escuela distinta a los demás playgrounds, no colisiona.
func upsertActivePeriod(tx *gorm.DB) error {
	id, err := uuid.Parse(periodID)
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
	code := "N4-2026-S1"
	p := entities.AcademicPeriod{
		ID:             id,
		SchoolID:       sid,
		AcademicUnitID: auid,
		Name:           "Semestre 1 2026",
		Code:           &code,
		Type:           "semester",
		StartDate:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
		IsActive:       true,
		AcademicYear:   academicYear,
		SortOrder:      1,
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
// dada en la unidad y período del playground, CON section_label.
// teacherMembIDStr puede ser "" (sesión sin docente → teacher_membership_id
// NULL). sectionLabel se setea siempre que no sea "" (forma parte del índice
// único natural uq_subject_offerings_natural). Idempotente por id.
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
// student_membership_id); OnConflict sobre ambas → idempotente. subjectIDStr y
// el periodID del playground son el subject_id y period_id de la oferta (copias
// denormalizadas e inmutables que respaldan el invariante una-oferta-por-materia-
// por-período, bug 0036). El playground tiene un único período activo.
func upsertEnrollment(tx *gorm.DB, offeringIDStr, subjectIDStr, studentMembIDStr string) error {
	oid, err := uuid.Parse(offeringIDStr)
	if err != nil {
		return err
	}
	subjID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		return err
	}
	pid, err := uuid.Parse(periodID)
	if err != nil {
		return err
	}
	smid, err := uuid.Parse(studentMembIDStr)
	if err != nil {
		return err
	}
	e := entities.SubjectOfferingEnrollment{
		OfferingID:          oid,
		SubjectID:           subjID,
		PeriodID:            pid,
		StudentMembershipID: smid,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "offering_id"}, {Name: "student_membership_id"}},
		DoNothing: true,
	}).Create(&e).Error
}
