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
//  3. academic.subjects          — 4 materias (2 por unidad).
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

	periodID = "68000000-0000-0000-0000-000000000008"

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
	if err := upsertSchool(tx); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: school: %w", err)
	}

	// DOS unidades académicas en la MISMA escuela: esto es lo que impide la
	// auto-selección de unidad en el login (len(units) != 1).
	if err := upsertAcademicUnit(tx, unitNorteID, unitNorteName, unitNorteCode); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: unit_norte: %w", err)
	}
	if err := upsertAcademicUnit(tx, unitSurID, unitSurName, unitSurCode); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: unit_sur: %w", err)
	}

	// Materias: 2 por unidad (datos para las pantallas unit-scoped).
	if err := upsertSubject(tx, subjectNorteMathID, unitNorteID, "Matemáticas Norte", "MAT-N"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_norte_math: %w", err)
	}
	if err := upsertSubject(tx, subjectNorteLangID, unitNorteID, "Lenguaje Norte", "LEN-N"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_norte_lang: %w", err)
	}
	if err := upsertSubject(tx, subjectSurMathID, unitSurID, "Matemáticas Sur", "MAT-S"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_sur_math: %w", err)
	}
	if err := upsertSubject(tx, subjectSurLangID, unitSurID, "Lenguaje Sur", "LEN-S"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: subject_sur_lang: %w", err)
	}

	if err := upsertActivePeriod(tx); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: academic_period: %w", err)
	}

	// Usuarios.
	if err := upsertUser(tx, adminUserID, AdminEmail, "Admin", "Multi"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_user: %w", err)
	}
	if err := upsertUser(tx, teacherNorteUserID, TeacherNorteEmail, "Docente", "Norte"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_user: %w", err)
	}
	if err := upsertUser(tx, teacherSurUserID, TeacherSurEmail, "Docente", "Sur"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_user: %w", err)
	}
	if err := upsertUser(tx, studentNorteUserID, StudentNorteEmail, "Alumno", "Norte"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_user: %w", err)
	}
	if err := upsertUser(tx, studentSurUserID, StudentSurEmail, "Alumno", "Sur"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_user: %w", err)
	}

	// Roles L4 para contexto de login (no se crean roles nuevos).
	if err := upsertUserRole(tx, adminUserID, l4.L4_ROLE_SCHOOL_ADMIN_ID); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_user_role: %w", err)
	}
	if err := upsertUserRole(tx, teacherNorteUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_user_role: %w", err)
	}
	if err := upsertUserRole(tx, teacherSurUserID, l4.L4_ROLE_TEACHER_ID); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentNorteUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_user_role: %w", err)
	}
	if err := upsertUserRole(tx, studentSurUserID, l4.L4_ROLE_STUDENT_ID); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_user_role: %w", err)
	}

	// Membresías:
	//  - admin con alcance COLEGIO (AcademicUnitID = NULL): contexto de
	//    colegio en el JWT, SIN unidad preseleccionada → dispara el selector.
	//  - docente/alumno Norte con alcance UNIDAD en "Sede Norte".
	//  - docente/alumno Sur con alcance UNIDAD en "Sede Sur".
	if err := upsertSchoolMembership(tx, adminMembID, adminUserID, "admin"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: admin_membership: %w", err)
	}
	if err := upsertUnitMembership(tx, teacherNorteMembID, teacherNorteUserID, unitNorteID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_norte_membership: %w", err)
	}
	if err := upsertUnitMembership(tx, teacherSurMembID, teacherSurUserID, unitSurID, "teacher"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: teacher_sur_membership: %w", err)
	}
	if err := upsertUnitMembership(tx, studentNorteMembID, studentNorteUserID, unitNorteID, "student"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_norte_membership: %w", err)
	}
	if err := upsertUnitMembership(tx, studentSurMembID, studentSurUserID, unitSurID, "student"); err != nil {
		return fmt.Errorf("playground_v2/multi_unidad: student_sur_membership: %w", err)
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

func upsertAcademicUnit(tx *gorm.DB, idStr, name, code string) error {
	id, err := uuid.Parse(idStr)
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
		Name:         name,
		Code:         code,
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

func upsertSubject(tx *gorm.DB, idStr, unitIDStr, name, code string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return err
	}
	auid, err := uuid.Parse(unitIDStr)
	if err != nil {
		return err
	}
	c := code
	subj := entities.Subject{
		ID:             id,
		SchoolID:       sid,
		AcademicUnitID: &auid,
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
	code := "MULTI-2026-S1"
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

// upsertUnitMembership crea una membresía con alcance UNIDAD (AcademicUnitID
// seteado). Idempotente por id.
func upsertUnitMembership(tx *gorm.DB, idStr, userIDStr, unitIDStr, roleKind string) error {
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
	auid, err := uuid.Parse(unitIDStr)
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
// lleve contexto de colegio SIN unidad preseleccionada. Idempotente por id.
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
