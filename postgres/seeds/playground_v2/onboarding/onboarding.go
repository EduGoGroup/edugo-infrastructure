// Package onboarding es un playground de la línea v2 para probar E2E el
// flujo de onboarding del plan 005 (N0.5): signup + redención de código +
// doble aprobación (gate de unidad + gate de colegio) + caso de auto-firma.
//
// A diferencia de los playgrounds v1 (que sembraban recursos+pantallas
// ad-hoc sobre L0), este paquete asume que el sistema completo (L0..L4) ya
// corrió. Reusa tal cual los roles, permisos y recursos de onboarding que
// L4 ya trae (academic.invitations.*, academic.join_requests.*,
// academic.join_request_approvals.*) — NO inventa roles ni permisos.
//
// Convive con los demás playgrounds sin colisionar: rango UUID propio
// 64000000-... y emails con sufijo @edugo.local.
//
// Roles L4 reutilizados (ver system/l4/roles_permissions.go):
//   - school_admin (scope=school): grants `academic.*` + `admin.*` +
//     `context.*` → cubre context.browse_schools, academic.invitations.*,
//     academic.join_requests.*, academic.join_request_approvals.* (firma
//     AMBOS gates y aprueba cualquier rol; ve y gestiona la bandeja).
//   - teacher (scope=unit): grants `academic.invitations.*`,
//     `academic.join_requests.*` y el literal
//     `academic.join_request_approvals.student` → firma SOLO el gate de
//     unidad de solicitudes con rol student (no profesores ni apoderados).
//   - student (scope=unit): usado por el usuario ya-miembro para tener una
//     membresía activa en el colegio.
//
// Lo que siembra:
//  1. academic.schools          — 1 colegio "Onboarding".
//  2. academic.academic_units   — 2 unidades (A y B) bajo ese colegio.
//  3. auth.users                — 3 usuarios sembrados con password "12345678":
//     colegio-admin (school_admin),
//     profesor-a    (teacher en unidad A),
//     ya-miembro    (student en unidad A).
//     El usuario NUEVO del flujo E2E NO se siembra:
//     se registra por signup en runtime.
//  4. iam.user_roles            — assignments 1×1 a los roles L4.
//  5. academic.memberships      — colegio-admin con academic_unit_id NULL
//     (alcance COLEGIO); profesor-a y ya-miembro
//     con academic_unit_id = unidad A (alcance UNIDAD).
//  6. academic.school_invitations — 3 códigos activos legibles:
//     ONB-STUDENT-A (role=student, unidad A),
//     ONB-TEACHER-A (role=teacher, unidad A),
//     ONB-STUDENT-B (role=student, unidad B).
//
// NO siembra academic.school_join_requests: se crean en runtime al redimir.
//
// Idempotente: OnConflict DoNothing en todas las inserciones (por id, o por
// la clave natural donde aplica).
package onboarding

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/catalog"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/common"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales de los usuarios SEMBRADOS (el usuario nuevo del E2E se
	// registra por signup en runtime, no aparece aquí).
	AdminEmail    = "colegio-admin@edugo.local"
	ProfesorEmail = "profesor-a@edugo.local"
	MiembroEmail  = "ya-miembro@edugo.local"
	Password      = "12345678"

	// Códigos de invitación (legibles, conocidos para el E2E).
	CodeStudentA = "ONB-STUDENT-A" // unidad A, role=student (flujo usuario nuevo)
	CodeTeacherA = "ONB-TEACHER-A" // unidad A, role=teacher (aprobar un profesor)
	CodeStudentB = "ONB-STUDENT-B" // unidad B, role=student (ya-miembro → auto-firma colegio)

	// Rango UUID 64000000-...: reservado para el playground onboarding.
	schoolID = "64000000-0000-0000-0000-000000000001"
	unitAID  = "64000000-0000-0000-0000-000000000002"
	unitBID  = "64000000-0000-0000-0000-000000000003"

	adminUserID    = "64000000-0000-0000-0000-000000000010"
	profesorUserID = "64000000-0000-0000-0000-000000000011"
	miembroUserID  = "64000000-0000-0000-0000-000000000012"

	adminMembID    = "64000000-0000-0000-0000-000000000020"
	profesorMembID = "64000000-0000-0000-0000-000000000021"
	miembroMembID  = "64000000-0000-0000-0000-000000000022"

	invStudentAID = "64000000-0000-0000-0000-000000000030"
	invTeacherAID = "64000000-0000-0000-0000-000000000031"
	invStudentBID = "64000000-0000-0000-0000-000000000032"

	schoolCode = "ONBOARDING"
	schoolName = "Colegio Onboarding"
	unitACode  = "ONB-UNIT-A"
	unitAName  = "Unidad A"
	unitBCode  = "ONB-UNIT-B"
	unitBName  = "Unidad B"

	academicYear = 2026
)

// Apply siembra el playground onboarding. Asume que L0..L4 corrieron (los
// roles school_admin/teacher/student y los recursos/permisos de onboarding
// ya existen). Orden: school → units → users → user_roles → memberships →
// invitations. Idempotente.
func Apply(tx *gorm.DB) error {
	sid := common.MustParseUUID(schoolID)
	unitA := common.MustParseUUID(unitAID)
	unitB := common.MustParseUUID(unitBID)

	if err := common.SeedSchool(tx, common.SchoolSpec{
		ID: sid, Name: schoolName, Code: schoolCode,
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: school: %w", err)
	}
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{
		ID: unitA, SchoolID: sid, Name: unitAName, Code: unitACode, AcademicYear: academicYear,
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: unit_a: %w", err)
	}
	if err := common.SeedAcademicUnit(tx, common.UnitSpec{
		ID: unitB, SchoolID: sid, Name: unitBName, Code: unitBCode, AcademicYear: academicYear,
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: unit_b: %w", err)
	}

	adminUser := common.MustParseUUID(adminUserID)
	profesorUser := common.MustParseUUID(profesorUserID)
	miembroUser := common.MustParseUUID(miembroUserID)

	if err := common.SeedUser(tx, common.UserSpec{ID: adminUser, Email: AdminEmail, Password: Password, FirstName: "Admin", LastName: "Colegio"}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: admin_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: profesorUser, Email: ProfesorEmail, Password: Password, FirstName: "Profesor", LastName: "Unidad A"}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: profesor_user: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{ID: miembroUser, Email: MiembroEmail, Password: Password, FirstName: "Ya", LastName: "Miembro"}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: miembro_user: %w", err)
	}

	// Asignación de roles L4 (no se crean roles nuevos).
	if err := common.SeedUserRole(tx, adminUser, common.MustParseUUID(l4.L4_ROLE_SCHOOL_ADMIN_ID)); err != nil {
		return fmt.Errorf("playground_v2/onboarding: admin_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, profesorUser, common.MustParseUUID(l4.L4_ROLE_TEACHER_ID)); err != nil {
		return fmt.Errorf("playground_v2/onboarding: profesor_user_role: %w", err)
	}
	if err := common.SeedUserRole(tx, miembroUser, common.MustParseUUID(l4.L4_ROLE_STUDENT_ID)); err != nil {
		return fmt.Errorf("playground_v2/onboarding: miembro_user_role: %w", err)
	}

	// Membresías: admin con alcance COLEGIO (academic_unit_id nil); profesor
	// y ya-miembro con alcance UNIDAD A.
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(adminMembID), UserID: adminUser, SchoolID: sid, AcademicUnitID: nil, Role: "admin",
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: admin_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(profesorMembID), UserID: profesorUser, SchoolID: sid, AcademicUnitID: &unitA, Role: "teacher",
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: profesor_membership: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID: common.MustParseUUID(miembroMembID), UserID: miembroUser, SchoolID: sid, AcademicUnitID: &unitA, Role: "student",
	}); err != nil {
		return fmt.Errorf("playground_v2/onboarding: miembro_membership: %w", err)
	}

	// Invitaciones activas con código conocido.
	if err := upsertInvitation(tx, invStudentAID, CodeStudentA, unitAID, "student", "Alumno — Unidad A"); err != nil {
		return fmt.Errorf("playground_v2/onboarding: inv_student_a: %w", err)
	}
	if err := upsertInvitation(tx, invTeacherAID, CodeTeacherA, unitAID, "teacher", "Profesor — Unidad A"); err != nil {
		return fmt.Errorf("playground_v2/onboarding: inv_teacher_a: %w", err)
	}
	if err := upsertInvitation(tx, invStudentBID, CodeStudentB, unitBID, "student", "Alumno — Unidad B"); err != nil {
		return fmt.Errorf("playground_v2/onboarding: inv_student_b: %w", err)
	}

	return nil
}

func upsertInvitation(tx *gorm.DB, idStr, code, unitIDStr, role, label string) error {
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
	invitationTypeID, err := catalog.ResolveInvitationTypeID(tx, role)
	if err != nil {
		return err
	}
	lbl := label
	inv := entities.SchoolInvitation{
		ID:               id,
		Code:             code,
		SchoolID:         sid,
		AcademicUnitID:   auid,
		InvitationTypeID: invitationTypeID,
		Label:            &lbl,
		UsesCount:        0,
		IsActive:         true,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&inv).Error
}
