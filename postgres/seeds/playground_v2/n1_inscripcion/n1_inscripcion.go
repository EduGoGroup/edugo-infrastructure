// Package n1_inscripcion es un playground de la línea v2 para validar el
// flujo N1 de inscripción de alumnos a materias (plan 006) y diagnosticar
// el landing del student (bug 0006).
//
// Como todo v2, asume que el sistema completo (L0..L4) ya corrió: reusa los
// roles L4 school_admin/teacher/student (ver system/l4/roles_permissions.go)
// para que los usuarios tengan contexto de login real — NO inventa roles ni
// permisos. El login resuelve active_context desde academic.memberships +
// iam.user_roles, así que cada usuario sembrado lleva ambas filas.
//
// Convive con los demás playgrounds sin colisionar: rango UUID propio
// 66000000-... (62000000=v2_screens_catalog, 64000000=onboarding) y emails
// con sufijo @n1.edugo.local.
//
// Lo que siembra:
//  1. academic.schools           — 1 colegio "Colegio N1 Inscripción".
//  2. academic.academic_units    — 1 unidad académica bajo ese colegio.
//  3. academic.subjects          — 3 materias de ESCUELA (AcademicUnitID=NULL,
//     ADR 0016): Matemáticas, Lenguaje, Ciencias. Nombres distintos → cumplen
//     UNIQUE(school_id, name).
//  4. academic.academic_periods  — 1 período ACTIVO (is_active=true).
//  5. auth.users                 — 1 admin + 1 docente + 3 alumnos, password "12345678".
//  6. iam.user_roles             — admin→school_admin L4, docente→teacher L4, alumnos→student L4.
//  7. academic.memberships       — admin con alcance COLEGIO (AcademicUnitID=NULL); los demás en la unidad.
//  8. academic.subject_offerings — 1 sesión por materia (modelo N1.7, plan 010
//     / ADR 0009). La sesión de Matemáticas lleva teacher_membership_id del
//     docente (lo dicta); Lenguaje y Ciencias quedan sin docente asignado
//     (teacher_membership_id NULL).
//  9. academic.subject_offering_enrollments — inscripción del alumno a la sesión:
//     - alumno 1  → Matemáticas + Lenguaje (inscrito).
//     - alumno 2  → Matemáticas + Lenguaje + Ciencias (inscrito).
//     - alumno 3  → SIN filas (sin inscribir, para ejercitar el flujo).
//
// Credenciales (todas password "12345678"):
//
//	admin-n1@n1.edugo.local      school_admin — director del colegio (alcance colegio)
//	docente-n1@n1.edugo.local    teacher      — dicta Matemáticas
//	alumno1-n1@n1.edugo.local    student      — inscrito en Matemáticas + Lenguaje
//	alumno2-n1@n1.edugo.local    student      — inscrito en las 3 materias
//	alumno3-n1@n1.edugo.local    student      — SIN materias (sin inscribir)
//
// Idempotente: OnConflict DoNothing por id (o clave natural compuesta en
// subject_offering_enrollments) en todas las inserciones.
package n1_inscripcion

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
	TeacherEmail  = "docente-n1@n1.edugo.local"
	Student1Email = "alumno1-n1@n1.edugo.local"
	Student2Email = "alumno2-n1@n1.edugo.local"
	Student3Email = "alumno3-n1@n1.edugo.local"
	AdminEmail    = "admin-n1@n1.edugo.local"
	Password      = "12345678"

	// Rango UUID 66000000-...: reservado para el playground n1_inscripcion.
	schoolID = "66000000-0000-0000-0000-000000000001"
	unitID   = "66000000-0000-0000-0000-000000000002"

	subjectMathID    = "66000000-0000-0000-0000-000000000003"
	subjectLangID    = "66000000-0000-0000-0000-000000000004"
	subjectScienceID = "66000000-0000-0000-0000-000000000005"

	periodID = "66000000-0000-0000-0000-000000000006"

	teacherUserID  = "66000000-0000-0000-0000-000000000010"
	student1UserID = "66000000-0000-0000-0000-000000000011"
	student2UserID = "66000000-0000-0000-0000-000000000012"
	student3UserID = "66000000-0000-0000-0000-000000000013"
	adminUserID    = "66000000-0000-0000-0000-000000000014"

	teacherMembID  = "66000000-0000-0000-0000-000000000020"
	student1MembID = "66000000-0000-0000-0000-000000000021"
	student2MembID = "66000000-0000-0000-0000-000000000022"
	student3MembID = "66000000-0000-0000-0000-000000000023"
	adminMembID    = "66000000-0000-0000-0000-000000000024"

	// Sesiones de materia (subject_offerings): una por materia en la unidad.
	offeringMathID    = "66000000-0000-0000-0000-000000000030"
	offeringLangID    = "66000000-0000-0000-0000-000000000031"
	offeringScienceID = "66000000-0000-0000-0000-000000000032"

	schoolCode = "N1-INSCRIPCION"
	schoolName = "Colegio N1 Inscripción"
	unitCode   = "N1-UNIT-MAIN"
	unitName   = "Sede Única N1"

	academicYear = 2026
)

// Apply siembra el playground n1_inscripcion. Asume que L0..L4 corrieron (los
// roles teacher/student y el esquema academic completo ya existen). Orden:
// school → unit → subjects → period → users → user_roles → memberships →
// subject_offerings → subject_offering_enrollments. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: school: %w", err)
	}
	if err := upsertAcademicUnit(tx); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: academic_unit: %w", err)
	}

	if err := upsertSubject(tx, subjectMathID, "Matemáticas", "MAT"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_math: %w", err)
	}
	if err := upsertSubject(tx, subjectLangID, "Lenguaje", "LEN"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_lang: %w", err)
	}
	if err := upsertSubject(tx, subjectScienceID, "Ciencias", "CIE"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_science: %w", err)
	}

	if err := upsertActivePeriod(tx); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: academic_period: %w", err)
	}

	if err := upsertUser(tx, teacherUserID, TeacherEmail, "Docente", "N1"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_user: %w", err)
	}
	if err := upsertUser(tx, student1UserID, Student1Email, "Alumno", "Uno"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_user: %w", err)
	}
	if err := upsertUser(tx, student2UserID, Student2Email, "Alumno", "Dos"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_user: %w", err)
	}
	if err := upsertUser(tx, student3UserID, Student3Email, "Alumno", "Tres"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_user: %w", err)
	}
	if err := upsertUser(tx, adminUserID, AdminEmail, "Admin", "N1"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := upsertUserRole(tx, teacherUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_user_role: %w", err)
	}
	if err := upsertUserRole(tx, student1UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_user_role: %w", err)
	}
	if err := upsertUserRole(tx, student2UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_user_role: %w", err)
	}
	if err := upsertUserRole(tx, student3UserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_user_role: %w", err)
	}
	if err := upsertUserRole(tx, adminUserID, l4.L4_ROLE_SCHOOL_ADMIN_ID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_user_role: %w", err)
	}

	// Membresías: todos con alcance UNIDAD en la misma unidad.
	if err := upsertMembership(tx, teacherMembID, teacherUserID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_membership: %w", err)
	}
	if err := upsertMembership(tx, student1MembID, student1UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_membership: %w", err)
	}
	if err := upsertMembership(tx, student2MembID, student2UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_membership: %w", err)
	}
	if err := upsertMembership(tx, student3MembID, student3UserID, "student"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_membership: %w", err)
	}
	// Membresía del admin con alcance COLEGIO (AcademicUnitID = NULL): el form
	// memberships-form exige contexto de colegio en el JWT del actor.
	if err := upsertSchoolMembership(tx, adminMembID, adminUserID, "admin"); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_membership: %w", err)
	}

	// Sesiones de materia (subject_offerings): una por materia en la unidad.
	// La de Matemáticas lleva al docente (la dicta); Lenguaje y Ciencias quedan
	// sin docente asignado (teacher_membership_id NULL).
	if err := upsertOffering(tx, offeringMathID, subjectMathID, teacherMembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_math: %w", err)
	}
	if err := upsertOffering(tx, offeringLangID, subjectLangID, ""); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_lang: %w", err)
	}
	if err := upsertOffering(tx, offeringScienceID, subjectScienceID, ""); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_science: %w", err)
	}

	// Inscripciones (subject_offering_enrollments):
	//  - alumno 1: Matemáticas + Lenguaje.
	//  - alumno 2: las 3 materias.
	//  - alumno 3: SIN inscribir (no se crea ninguna fila).
	if err := upsertEnrollment(tx, offeringMathID, student1MembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_enroll_math: %w", err)
	}
	if err := upsertEnrollment(tx, offeringLangID, student1MembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_enroll_lang: %w", err)
	}
	if err := upsertEnrollment(tx, offeringMathID, student2MembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_math: %w", err)
	}
	if err := upsertEnrollment(tx, offeringLangID, student2MembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_lang: %w", err)
	}
	if err := upsertEnrollment(tx, offeringScienceID, student2MembID); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_science: %w", err)
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
// AcademicUnitID = nil. La materia es catálogo de la escuela; su ubicación en
// la unidad la da la sesión (subject_offering). Los 3 nombres del playground
// son distintos, así que cumplen UNIQUE(school_id, name) sin deduplicar.
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
	code := "N1-2026-S1"
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
// dada en la unidad y período del playground. teacherMembIDStr puede ser ""
// (sesión sin docente asignado → teacher_membership_id NULL). Idempotente por id.
func upsertOffering(tx *gorm.DB, idStr, subjectIDStr, teacherMembIDStr string) error {
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

	o := entities.SubjectOffering{
		ID:                  id,
		SchoolID:            sid,
		SubjectID:           subjID,
		AcademicUnitID:      auid,
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
