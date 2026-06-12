// Package multi_unidad es un playground de la línea v2 para validar el
// SELECTOR DE UNIDAD en el login web.
//
// Foco de la validación:
//   - Una escuela con DOS unidades académicas activas ("Sede Norte" y
//     "Sede Sur"). El login de identity auto-selecciona unidad SOLO si la
//     escuela tiene exactamente una unidad activa
//     (FindSchoolAcademicUnits → len(units)==1, ver login.go). Con dos
//     unidades, el active_context del JWT queda con academic_unit_id SIN
//     setear → la web debe pintar el selector de unidad al abrir una
//     pantalla unit-scoped (memberships / subjects).
//   - Un school_admin con membership de ALCANCE COLEGIO
//     (AcademicUnitID = NULL): tiene contexto de colegio en el JWT pero
//     ninguna unidad preseleccionada.
//   - Datos mínimos en AMBAS unidades (materias + miembros) para que, tras
//     elegir una en el selector, las pantallas unit-scoped muestren algo.
//
// Como todo v2, asume que el sistema completo (L0..L4) ya corrió: reusa los
// roles L4 school_admin/teacher/student (ver system/l4/roles_permissions.go)
// para que los usuarios tengan contexto de login real — NO inventa roles ni
// permisos. El login resuelve active_context desde academic.memberships +
// iam.user_roles, así que cada usuario sembrado lleva ambas filas.
//
// Convive con los demás playgrounds sin colisionar: rango UUID propio
// 68000000-... (distinto del 67000000 de n17_secciones y del 66000000 de
// n1_inscripcion) y emails con sufijo @multi.edugo.local. La escuela es
// distinta, así que el índice único parcial de período activo (por school_id
// WHERE is_active) no colisiona con los otros playgrounds.
//
// Lo que siembra:
//  1. academic.schools           — 1 colegio "Colegio Multi-Unidad".
//  2. academic.academic_units    — 2 unidades: "Sede Norte", "Sede Sur".
//  3. academic.subjects          — 4 materias de ESCUELA (AcademicUnitID=NULL,
//     ADR 0016), nombres distintos "… Norte" / "… Sur" → UNIQUE(school_id, name).
//  4. academic.academic_periods  — 1 período ACTIVO (is_active=true).
//  5. auth.users                 — 1 admin + 2 docentes + 2 alumnos, password "12345678".
//  6. iam.user_roles             — admin→school_admin L4, docentes→teacher L4, alumnos→student L4.
//  7. academic.memberships       — admin con alcance COLEGIO (AcademicUnitID=NULL);
//     docente+alumno Norte en "Sede Norte"; docente+alumno Sur en "Sede Sur".
//
// Credenciales (todas password "12345678"):
//
//	admin-multi@multi.edugo.local        school_admin — director (alcance colegio, 2 unidades)
//	docente-norte@multi.edugo.local      teacher      — Sede Norte
//	docente-sur@multi.edugo.local        teacher      — Sede Sur
//	alumno-norte@multi.edugo.local       student      — Sede Norte
//	alumno-sur@multi.edugo.local         student      — Sede Sur
//
// Idempotente: OnConflict DoNothing por id en todas las inserciones.
package multi_unidad

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"gorm.io/gorm"
)

const (
	// Credenciales de los usuarios sembrados.
	AdminEmail        = "admin-multi@multi.edugo.local"
	TeacherNorteEmail = "docente-norte@multi.edugo.local"
	TeacherSurEmail   = "docente-sur@multi.edugo.local"
	StudentNorteEmail = "alumno-norte@multi.edugo.local"
	StudentSurEmail   = "alumno-sur@multi.edugo.local"
	Password          = "12345678"

	// Rango UUID 68000000-...: reservado para el playground multi_unidad.
	schoolID    = "68000000-0000-0000-0000-000000000001"
	unitNorteID = "68000000-0000-0000-0000-000000000002"
	unitSurID   = "68000000-0000-0000-0000-000000000003"

	// Materias (2 por unidad).
	subjectNorteMathID = "68000000-0000-0000-0000-000000000004"
	subjectNorteLangID = "68000000-0000-0000-0000-000000000005"
	subjectSurMathID   = "68000000-0000-0000-0000-000000000006"
	subjectSurLangID   = "68000000-0000-0000-0000-000000000007"

	// Un período ACTIVO por unidad: el índice único parcial de período activo
	// es por (school_id, academic_unit_id), así que cada sede lleva el suyo.
	periodNorteID = "68000000-0000-0000-0000-000000000008"
	periodSurID   = "68000000-0000-0000-0000-000000000009"

	adminUserID        = "68000000-0000-0000-0000-000000000010"
	teacherNorteUserID = "68000000-0000-0000-0000-000000000011"
	teacherSurUserID   = "68000000-0000-0000-0000-000000000012"
	studentNorteUserID = "68000000-0000-0000-0000-000000000013"
	studentSurUserID   = "68000000-0000-0000-0000-000000000014"

	adminMembID        = "68000000-0000-0000-0000-000000000020"
	teacherNorteMembID = "68000000-0000-0000-0000-000000000021"
	teacherSurMembID   = "68000000-0000-0000-0000-000000000022"
	studentNorteMembID = "68000000-0000-0000-0000-000000000023"
	studentSurMembID   = "68000000-0000-0000-0000-000000000024"

	schoolCode = "MULTI-UNIDAD"
	schoolName = "Colegio Multi-Unidad"

	unitNorteCode = "MULTI-NORTE"
	unitNorteName = "Sede Norte"
	unitSurCode   = "MULTI-SUR"
	unitSurName   = "Sede Sur"

	academicYear = 2026
)

// Apply siembra el playground multi_unidad. Asume que L0..L4 corrieron (los
// roles school_admin/teacher/student y el esquema academic completo ya
// existen). Orden: school → units → subjects → period → users → user_roles →
// memberships. Idempotente.
func Apply(tx *gorm.DB) error {
	sid := common.MustParseUUID(schoolID)
	unitNorte := common.MustParseUUID(unitNorteID)
	unitSur := common.MustParseUUID(unitSurID)

	if err := common.SeedSchool(tx, common.SchoolSpec{
		ID:   sid,
		Name: schoolName,
		Code: schoolCode,
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: school: %w", err)
	}

	// DOS unidades académicas en la MISMA escuela: esto es lo que impide la
	// auto-selección de unidad en el login (len(units) != 1).
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{
		ID: unitNorte, SchoolID: sid, Name: unitNorteName, Code: unitNorteCode, AcademicYear: academicYear,
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: unit_norte: %w", err)
	}
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{
		ID: unitSur, SchoolID: sid, Name: unitSurName, Code: unitSurCode, AcademicYear: academicYear,
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: unit_sur: %w", err)
	}

	// Materias de ESCUELA (ADR 0016): catálogo único de la escuela, sin anclar
	// a unidad. Los 4 nombres son distintos → cumplen UNIQUE(school_id, name).
	if err := common.SeedSubject(tx, common.SubjectSpec{
		ID: common.MustParseUUID(subjectNorteMathID), SchoolID: sid, Name: "Matemáticas Norte", Code: "MAT-N",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_norte_math: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{
		ID: common.MustParseUUID(subjectNorteLangID), SchoolID: sid, Name: "Lenguaje Norte", Code: "LEN-N",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_norte_lang: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{
		ID: common.MustParseUUID(subjectSurMathID), SchoolID: sid, Name: "Matemáticas Sur", Code: "MAT-S",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_sur_math: %w", err)
	}
	if err := common.SeedSubject(tx, common.SubjectSpec{
		ID: common.MustParseUUID(subjectSurLangID), SchoolID: sid, Name: "Lenguaje Sur", Code: "LEN-S",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_sur_lang: %w", err)
	}

	// Un período ACTIVO por unidad (la exclusividad del activo es por unidad).
	// Nombre/tipo/fechas usan los defaults del común (Semestre 1 2026, semester,
	// 2026-03-01 → 2026-07-31); SortOrder=1 como en el seed original.
	if err := common.SeedActivePeriod(tx, common.PeriodSpec{
		ID: common.MustParseUUID(periodNorteID), SchoolID: sid, AcademicUnitID: unitNorte,
		Code: "MULTI-N-2026-S1", AcademicYear: academicYear, SortOrder: 1,
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: academic_period_norte: %w", err)
	}
	if err := common.SeedActivePeriod(tx, common.PeriodSpec{
		ID: common.MustParseUUID(periodSurID), SchoolID: sid, AcademicUnitID: unitSur,
		Code: "MULTI-S-2026-S1", AcademicYear: academicYear, SortOrder: 1,
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: academic_period_sur: %w", err)
	}

	// Usuarios.
	adminUser := common.MustParseUUID(adminUserID)
	teacherNorteUser := common.MustParseUUID(teacherNorteUserID)
	teacherSurUser := common.MustParseUUID(teacherSurUserID)
	studentNorteUser := common.MustParseUUID(studentNorteUserID)
	studentSurUser := common.MustParseUUID(studentSurUserID)

	if err := common.SeedUser(tx, common.UserSpec{ID: adminUser, Email: AdminEmail, Password: Password, FirstName: "Admin", LastName: "Multi"}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: teacherNorteUser, Email: TeacherNorteEmail, Password: Password, FirstName: "Docente", LastName: "Norte"}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: teacherSurUser, Email: TeacherSurEmail, Password: Password, FirstName: "Docente", LastName: "Sur"}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentNorteUser, Email: StudentNorteEmail, Password: Password, FirstName: "Alumno", LastName: "Norte"}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: studentSurUser, Email: StudentSurEmail, Password: Password, FirstName: "Alumno", LastName: "Sur"}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := common.SeedUserRole(tx, adminUser, common.MustParseUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID)); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, teacherNorteUser, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, teacherSurUser, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentNorteUser, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, studentSurUser, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_user_role: %w", err)
	}

	// Membresías:
	//  - admin con alcance COLEGIO (AcademicUnitID = nil): contexto de
	//    colegio en el JWT, SIN unidad preseleccionada → dispara el selector.
	//  - docente/alumno Norte con alcance UNIDAD en "Sede Norte".
	//  - docente/alumno Sur con alcance UNIDAD en "Sede Sur".
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(adminMembID), UserID: adminUser, SchoolID: sid, AcademicUnitID: nil, Role: "admin",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(teacherNorteMembID), UserID: teacherNorteUser, SchoolID: sid, AcademicUnitID: &unitNorte, Role: "teacher",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(teacherSurMembID), UserID: teacherSurUser, SchoolID: sid, AcademicUnitID: &unitSur, Role: "teacher",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(studentNorteMembID), UserID: studentNorteUser, SchoolID: sid, AcademicUnitID: &unitNorte, Role: "student",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(studentSurMembID), UserID: studentSurUser, SchoolID: sid, AcademicUnitID: &unitSur, Role: "student",
	}); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_membership: %w", err)
	}

	return nil
}
