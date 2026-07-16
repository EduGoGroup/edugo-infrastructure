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
//     / ADR 0009). Las tres sesiones (Matemáticas, Lenguaje, Ciencias) llevan
//     teacher_membership_id del docente-n1 (dicta las tres): un aula con alumnos
//     inscritos exige docente asignado (deuda 028).
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
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"gorm.io/gorm"
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
	// Ids del playground (constantes string) parseados una sola vez.
	schoolUUID := common.MustParseUUID(schoolID)
	unitUUID := common.MustParseUUID(unitID)
	subjectMathUUID := common.MustParseUUID(subjectMathID)
	subjectLangUUID := common.MustParseUUID(subjectLangID)
	subjectScienceUUID := common.MustParseUUID(subjectScienceID)
	periodUUID := common.MustParseUUID(periodID)

	teacherUserUUID := common.MustParseUUID(teacherUserID)
	student1UserUUID := common.MustParseUUID(student1UserID)
	student2UserUUID := common.MustParseUUID(student2UserID)
	student3UserUUID := common.MustParseUUID(student3UserID)
	adminUserUUID := common.MustParseUUID(adminUserID)

	teacherMembUUID := common.MustParseUUID(teacherMembID)
	student1MembUUID := common.MustParseUUID(student1MembID)
	student2MembUUID := common.MustParseUUID(student2MembID)
	student3MembUUID := common.MustParseUUID(student3MembID)
	adminMembUUID := common.MustParseUUID(adminMembID)

	offeringMathUUID := common.MustParseUUID(offeringMathID)
	offeringLangUUID := common.MustParseUUID(offeringLangID)
	offeringScienceUUID := common.MustParseUUID(offeringScienceID)

	if err := common.SeedSchool(tx, common.SchoolSpec{ID: schoolUUID, Name: schoolName, Code: schoolCode}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: school: %w", err)
	}
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{ID: unitUUID, SchoolID: schoolUUID, Name: unitName, Code: unitCode, AcademicYear: academicYear}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: academic_unit: %w", err)
	}

	if err := common.SeedSubject(tx, common.SubjectSpec{ID: subjectMathUUID, SchoolID: schoolUUID, Name: "Matemáticas", Code: "MAT"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_math: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{ID: subjectLangUUID, SchoolID: schoolUUID, Name: "Lenguaje", Code: "LEN"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_lang: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{ID: subjectScienceUUID, SchoolID: schoolUUID, Name: "Ciencias", Code: "CIE"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: subject_science: %w", err)
	}

	if err := common.SeedActivePeriod(tx, common.PeriodSpec{
		ID:             periodUUID,
		SchoolID:       schoolUUID,
		AcademicUnitID: unitUUID,
		Name:           "Semestre 1 2026",
		Code:           "N1-2026-S1",
		Type:           "semester",
		AcademicYear:   academicYear,
		SortOrder:      1,
	}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: academic_period: %w", err)
	}

	if err := common.SeedUser(tx, common.UserSpec{ID: teacherUserUUID, Email: TeacherEmail, Password: Password, FirstName: "Docente", LastName: "N1"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: student1UserUUID, Email: Student1Email, Password: Password, FirstName: "Alumno", LastName: "Uno"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: student2UserUUID, Email: Student2Email, Password: Password, FirstName: "Alumno", LastName: "Dos"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: student3UserUUID, Email: Student3Email, Password: Password, FirstName: "Alumno", LastName: "Tres"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: adminUserUUID, Email: AdminEmail, Password: Password, FirstName: "Admin", LastName: "N1"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := common.SeedUserRole(tx, teacherUserUUID, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, student1UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, student2UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, student3UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, adminUserUUID, common.MustParseUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID)); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_user_role: %w", err)
	}

	// Membresías: todos con alcance UNIDAD en la misma unidad.
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: teacherMembUUID, UserID: teacherUserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "teacher"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: teacher_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: student1MembUUID, UserID: student1UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: student2MembUUID, UserID: student2UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: student3MembUUID, UserID: student3UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student3_membership: %w", err)
	}
	// Membresía del admin con alcance COLEGIO (AcademicUnitID = NULL): el form
	// memberships-form exige contexto de colegio en el JWT del actor.
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: adminMembUUID, UserID: adminUserUUID, SchoolID: schoolUUID, AcademicUnitID: nil, Role: "admin"}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: admin_membership: %w", err)
	}

	// Sesiones de materia (subject_offerings): una por materia en la unidad.
	// Las tres llevan al mismo docente (docente-n1 dicta en esta unidad). Antes
	// Lenguaje y Ciencias quedaban sin docente (teacher_membership_id NULL), pero
	// con la regla de la deuda 028 —no puede haber alumnos inscritos en un aula sin
	// docente— eso volvía inconsistentes las inscripciones que este fixture crea
	// más abajo. Todas SIN section_label.
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringMathUUID, SchoolID: schoolUUID, SubjectID: subjectMathUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, TeacherMembershipID: &teacherMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_math: %w", err)
	}
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringLangUUID, SchoolID: schoolUUID, SubjectID: subjectLangUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, TeacherMembershipID: &teacherMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_lang: %w", err)
	}
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringScienceUUID, SchoolID: schoolUUID, SubjectID: subjectScienceUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, TeacherMembershipID: &teacherMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: offering_science: %w", err)
	}

	// Inscripciones (subject_offering_enrollments):
	//  - alumno 1: Matemáticas + Lenguaje.
	//  - alumno 2: las 3 materias.
	//  - alumno 3: SIN inscribir (no se crea ninguna fila).
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMathUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: student1MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_enroll_math: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringLangUUID, SubjectID: subjectLangUUID, PeriodID: periodUUID, StudentMembershipID: student1MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student1_enroll_lang: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMathUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: student2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_math: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringLangUUID, SubjectID: subjectLangUUID, PeriodID: periodUUID, StudentMembershipID: student2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_lang: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringScienceUUID, SubjectID: subjectScienceUUID, PeriodID: periodUUID, StudentMembershipID: student2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n1_inscripcion: student2_enroll_science: %w", err)
	}

	return nil
}
