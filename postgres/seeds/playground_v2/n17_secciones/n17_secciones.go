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
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"gorm.io/gorm"
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
	// Ids del playground (constantes string) parseados una sola vez.
	schoolUUID := common.MustParseUUID(schoolID)
	unitUUID := common.MustParseUUID(unitID)
	subjectMathUUID := common.MustParseUUID(subjectMathID)
	subjectLangUUID := common.MustParseUUID(subjectLangID)
	periodUUID := common.MustParseUUID(periodID)

	teacherXUserUUID := common.MustParseUUID(teacherXUserID)
	teacherYUserUUID := common.MustParseUUID(teacherYUserID)
	studentA1UserUUID := common.MustParseUUID(studentA1UserID)
	studentA2UserUUID := common.MustParseUUID(studentA2UserID)
	studentB1UserUUID := common.MustParseUUID(studentB1UserID)
	studentB2UserUUID := common.MustParseUUID(studentB2UserID)
	studentLibreUserUUID := common.MustParseUUID(studentLibreUserID)
	adminUserUUID := common.MustParseUUID(adminUserID)

	teacherXMembUUID := common.MustParseUUID(teacherXMembID)
	teacherYMembUUID := common.MustParseUUID(teacherYMembID)
	studentA1MembUUID := common.MustParseUUID(studentA1MembID)
	studentA2MembUUID := common.MustParseUUID(studentA2MembID)
	studentB1MembUUID := common.MustParseUUID(studentB1MembID)
	studentB2MembUUID := common.MustParseUUID(studentB2MembID)
	studentLibreMembUUID := common.MustParseUUID(studentLibreMembID)
	adminMembUUID := common.MustParseUUID(adminMembID)

	offeringMatAUUID := common.MustParseUUID(offeringMatAID)
	offeringMatBUUID := common.MustParseUUID(offeringMatBID)
	offeringLenAUUID := common.MustParseUUID(offeringLenAID)

	if err := common.SeedSchool(tx, common.SchoolSpec{ID: schoolUUID, Name: schoolName, Code: schoolCode}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: school: %w", err)
	}
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{ID: unitUUID, SchoolID: schoolUUID, Name: unitName, Code: unitCode, AcademicYear: academicYear}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: academic_unit: %w", err)
	}

	if err := common.SeedSubject(tx, common.SubjectSpec{ID: subjectMathUUID, SchoolID: schoolUUID, Name: "Matemáticas", Code: "MAT"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: subject_math: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{ID: subjectLangUUID, SchoolID: schoolUUID, Name: "Lenguaje", Code: "LEN"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: subject_lang: %w", err)
	}

	if err := common.SeedActivePeriod(tx, common.PeriodSpec{
		ID:             periodUUID,
		SchoolID:       schoolUUID,
		AcademicUnitID: unitUUID,
		Name:           "Semestre 1 2026",
		Code:           "N17-2026-S1",
		Type:           "semester",
		AcademicYear:   academicYear,
		SortOrder:      1,
	}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: academic_period: %w", err)
	}

	if err := common.SeedUser(tx, common.UserSpec{ID: teacherXUserUUID, Email: TeacherXEmail, Password: Password, FirstName: "Docente", LastName: "X"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: teacherYUserUUID, Email: TeacherYEmail, Password: Password, FirstName: "Docente", LastName: "Y"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentA1UserUUID, Email: StudentA1Email, Password: Password, FirstName: "Alumno", LastName: "A1"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentA2UserUUID, Email: StudentA2Email, Password: Password, FirstName: "Alumno", LastName: "A2"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentB1UserUUID, Email: StudentB1Email, Password: Password, FirstName: "Alumno", LastName: "B1"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentB2UserUUID, Email: StudentB2Email, Password: Password, FirstName: "Alumno", LastName: "B2"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentLibreUserUUID, Email: StudentLibreEmail, Password: Password, FirstName: "Alumno", LastName: "Libre"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: adminUserUUID, Email: AdminEmail, Password: Password, FirstName: "Admin", LastName: "N17"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := common.SeedUserRole(tx, teacherXUserUUID, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, teacherYUserUUID, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentA1UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentA2UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentB1UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentB2UserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentLibreUserUUID, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, adminUserUUID, common.MustParseUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID)); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_user_role: %w", err)
	}

	// Membresías: docentes y alumnos con alcance UNIDAD en la misma unidad.
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: teacherXMembUUID, UserID: teacherXUserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "teacher"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_x_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: teacherYMembUUID, UserID: teacherYUserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "teacher"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: teacher_y_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: studentA1MembUUID, UserID: studentA1UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a1_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: studentA2MembUUID, UserID: studentA2UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_a2_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: studentB1MembUUID, UserID: studentB1UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b1_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: studentB2MembUUID, UserID: studentB2UserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_b2_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: studentLibreMembUUID, UserID: studentLibreUserUUID, SchoolID: schoolUUID, AcademicUnitID: &unitUUID, Role: "student"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: student_libre_membership: %w", err)
	}
	// Membresía del admin con alcance COLEGIO (AcademicUnitID = NULL): el form
	// memberships-form exige contexto de colegio en el JWT del actor.
	if err := common.SeedMembership(tx, common.MembershipSpec{ID: adminMembUUID, UserID: adminUserUUID, SchoolID: schoolUUID, AcademicUnitID: nil, Role: "admin"}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: admin_membership: %w", err)
	}

	// Sesiones de materia (subject_offerings) CON section_label:
	//  - Mat-A: Matemáticas, sección "A", docente X.
	//  - Mat-B: Matemáticas, sección "B", docente X (mismo docente, 2 sesiones).
	//  - Len-A: Lenguaje, sección "A", docente Y.
	sectionA := "A"
	sectionB := "B"
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringMatAUUID, SchoolID: schoolUUID, SubjectID: subjectMathUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, SectionLabel: &sectionA, TeacherMembershipID: &teacherXMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_mat_a: %w", err)
	}
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringMatBUUID, SchoolID: schoolUUID, SubjectID: subjectMathUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, SectionLabel: &sectionB, TeacherMembershipID: &teacherXMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_mat_b: %w", err)
	}
	if err := common.SeedOffering(tx, common.OfferingSpec{ID: offeringLenAUUID, SchoolID: schoolUUID, SubjectID: subjectLangUUID, AcademicUnitID: unitUUID, PeriodID: periodUUID, SectionLabel: &sectionA, TeacherMembershipID: &teacherYMembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: offering_len_a: %w", err)
	}

	// Inscripciones (subject_offering_enrollments):
	//  - Mat-A: alumno A1, alumno A2.
	//  - Mat-B: alumno B1, alumno B2.
	//  - Len-A: alumno A1, alumno A2 (alumno en 2 sesiones).
	//  - alumno Libre: SIN inscribir (no se crea ninguna fila).
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMatAUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: studentA1MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_a_a1: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMatAUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: studentA2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_a_a2: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMatBUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: studentB1MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_b_b1: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringMatBUUID, SubjectID: subjectMathUUID, PeriodID: periodUUID, StudentMembershipID: studentB2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_mat_b_b2: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringLenAUUID, SubjectID: subjectLangUUID, PeriodID: periodUUID, StudentMembershipID: studentA1MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_len_a_a1: %w", err)
	}
	if err := common.SeedEnrollment(tx, common.EnrollmentSpec{OfferingID: offeringLenAUUID, SubjectID: subjectLangUUID, PeriodID: periodUUID, StudentMembershipID: studentA2MembUUID}); err != nil {
		return fmt.Errorf("playground_v2/n17_secciones: enroll_len_a_a2: %w", err)
	}

	return nil
}
