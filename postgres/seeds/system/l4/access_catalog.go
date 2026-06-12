package l4

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MP-08 — acceso por sistema + tipos de invitación (capa L4).
//
// Este archivo siembra los 4 catálogos de MP-08:
//   1. iam.systems                  — apps del ecosistema (kmp, admin-tool).
//   2. iam.system_roles             — qué roles entran a cada app (puente).
//   3. academic.invitation_types    — catálogo global de tipos de invitación.
//   4. academic.school_invitation_roles — equivalencia tipo→rol IAM por escuela
//      (sólo la escuela demo L1 desde aquí; las escuelas de playground reciben
//      las equivalencias por defecto vía common.SeedSchool, que invoca el helper
//      compartido SeedDefaultSchoolInvitationRoles).
//
// Todas las funciones son idempotentes (ON CONFLICT DO NOTHING) y referencian
// roles/tipos por id (FK por id, nunca por nombre). Los ids hardcodeados viven
// en access_catalog_constants.go.

// l1SchoolDemoID es el id de la escuela demo sembrada por la capa L1. Re-pegado
// como literal igual que los role ids de L0/L1: `l4` no puede importar `layers`
// (ciclo, porque layers.l4Layer importa l4).
const l1SchoolDemoID = "b1000000-0000-0000-0000-000000000003" // layers.L1_SCHOOL_DEMO_ID

// ApplySystems siembra iam.systems (2 apps del ecosistema).
//
// Sin deps de FK. Debe correr antes de ApplySystemRoles (FK system_id).
func ApplySystems(tx *gorm.DB) error {
	systems := buildL4Systems()
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&systems).Error
}

// buildL4Systems construye las 2 filas de iam.systems. Helper compartido por
// ApplySystems y por el accessor público Systems().
func buildL4Systems() []entities.System {
	kmpDesc := "App principal del ecosistema EduGo para estudiantes, docentes, representantes y staff."
	adminDesc := "Herramienta de escritorio para la administración del ecosistema EduGo."
	return []entities.System{
		{
			ID:          uuid.MustParse(L4_SYSTEM_KMP_ID),
			Key:         L4_SYSTEM_KMP_KEY,
			Name:        "EduGo App",
			Description: &kmpDesc,
		},
		{
			ID:          uuid.MustParse(L4_SYSTEM_ADMIN_TOOL_ID),
			Key:         L4_SYSTEM_ADMIN_TOOL_KEY,
			Name:        "Herramienta de Administración",
			Description: &adminDesc,
		},
	}
}

// ApplySystemRoles siembra iam.system_roles: qué roles IAM entran a cada app.
//
//   - kmp:        los 12 roles del ecosistema (acceso completo).
//   - admin-tool: SOLO staff/admin (super_admin, school_admin,
//     school_coordinator, school_director, readonly_auditor). REGLA DURA
//     (DEC-C, MP-08): student/teacher/guardian NO entran a admin-tool.
//
// FK: requiere iam.systems (ApplySystems) e iam.roles (ApplyRolesPermissions /
// L0 / L1) sembrados antes. Idempotente: id derivado SHA1(system:role) para que
// reaplicar no duplique (mismo patrón que buildUserRole).
func ApplySystemRoles(tx *gorm.DB) error {
	// Los 12 roles del ecosistema (FK por id, vía constantes).
	allRoles := []string{
		L4_ROLE_STUDENT_ID,
		L4_ROLE_TEACHER_ID,
		L4_ROLE_GUARDIAN_ID,
		L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_SCHOOL_DIRECTOR_ID,
		L4_ROLE_SCHOOL_COORDINATOR_ID,
		L4_ROLE_SCHOOL_ASSISTANT_ID,
		L4_ROLE_ASSISTANT_TEACHER_ID,
		L4_ROLE_OBSERVER_ID,
		L4_ROLE_READONLY_AUDITOR_ID,
		l0RoleSuperAdminID,
		l1RoleAnnouncementViewerID,
	}

	// admin-tool: SOLO staff/admin. student/teacher/guardian excluidos (DEC-C).
	adminToolRoles := []string{
		l0RoleSuperAdminID,
		L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_SCHOOL_COORDINATOR_ID,
		L4_ROLE_SCHOOL_DIRECTOR_ID,
		L4_ROLE_READONLY_AUDITOR_ID,
	}

	rows := make([]entities.SystemRole, 0, len(allRoles)+len(adminToolRoles))

	kmpSystemID := uuid.MustParse(L4_SYSTEM_KMP_ID)
	for _, roleStr := range allRoles {
		roleID, err := uuid.Parse(roleStr)
		if err != nil {
			return fmt.Errorf("ApplySystemRoles: parse kmp role id %s: %w", roleStr, err)
		}
		rows = append(rows, buildSystemRole(kmpSystemID, roleID))
	}

	adminToolSystemID := uuid.MustParse(L4_SYSTEM_ADMIN_TOOL_ID)
	for _, roleStr := range adminToolRoles {
		roleID, err := uuid.Parse(roleStr)
		if err != nil {
			return fmt.Errorf("ApplySystemRoles: parse admin-tool role id %s: %w", roleStr, err)
		}
		rows = append(rows, buildSystemRole(adminToolSystemID, roleID))
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(&rows, 20).Error
}

// buildSystemRole construye una fila de iam.system_roles con id determinístico
// SHA1(system:role), de modo que reaplicar el seed no produzca duplicados aun
// cuando el unique compuesto (system_id, iam_role_id) ya exista.
func buildSystemRole(systemID, roleID uuid.UUID) entities.SystemRole {
	derived := uuid.NewSHA1(uuid.NameSpaceOID, []byte("system_role:"+systemID.String()+":"+roleID.String()))
	return entities.SystemRole{
		ID:        derived,
		SystemID:  systemID,
		IAMRoleID: roleID,
	}
}

// ApplyInvitationTypes siembra academic.invitation_types (catálogo global de 6
// tipos). Sin deps de FK. Debe correr antes de las equivalencias por escuela
// (school_invitation_roles → invitation_type_id).
func ApplyInvitationTypes(tx *gorm.DB) error {
	types := buildL4InvitationTypes()
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&types).Error
}

// buildL4InvitationTypes construye las 6 filas de academic.invitation_types.
// Helper compartido por ApplyInvitationTypes y por el accessor InvitationTypes().
func buildL4InvitationTypes() []entities.InvitationType {
	return []entities.InvitationType{
		{ID: uuid.MustParse(L4_INVITATION_TYPE_TEACHER_ID), Key: L4_INVITATION_TYPE_TEACHER_KEY, Label: "Profesor", RequiresUnit: true},
		{ID: uuid.MustParse(L4_INVITATION_TYPE_STUDENT_ID), Key: L4_INVITATION_TYPE_STUDENT_KEY, Label: "Estudiante", RequiresUnit: true},
		{ID: uuid.MustParse(L4_INVITATION_TYPE_GUARDIAN_ID), Key: L4_INVITATION_TYPE_GUARDIAN_KEY, Label: "Representante", RequiresUnit: true},
		{ID: uuid.MustParse(L4_INVITATION_TYPE_COORDINATOR_ID), Key: L4_INVITATION_TYPE_COORDINATOR_KEY, Label: "Coordinador", RequiresUnit: false},
		{ID: uuid.MustParse(L4_INVITATION_TYPE_ADMIN_ID), Key: L4_INVITATION_TYPE_ADMIN_KEY, Label: "Administrador", RequiresUnit: false},
		{ID: uuid.MustParse(L4_INVITATION_TYPE_ASSISTANT_ID), Key: L4_INVITATION_TYPE_ASSISTANT_KEY, Label: "Asistente", RequiresUnit: true},
	}
}

// invitationTypeToRole mapea cada tipo de invitación a su rol IAM por defecto.
// Origen: edugo-api-academic/internal/core/domain/membership.go (switch tipo→rol).
// Es el reparto default que recibe TODA escuela; admin-tool podrá ajustarlo por
// escuela en una fase posterior (school_invitation_roles es por escuela).
var invitationTypeToRole = []struct {
	typeIDStr string
	roleIDStr string
}{
	{L4_INVITATION_TYPE_TEACHER_ID, L4_ROLE_TEACHER_ID},
	{L4_INVITATION_TYPE_STUDENT_ID, L4_ROLE_STUDENT_ID},
	{L4_INVITATION_TYPE_GUARDIAN_ID, L4_ROLE_GUARDIAN_ID},
	{L4_INVITATION_TYPE_COORDINATOR_ID, L4_ROLE_SCHOOL_COORDINATOR_ID},
	{L4_INVITATION_TYPE_ADMIN_ID, L4_ROLE_SCHOOL_ADMIN_ID},
	{L4_INVITATION_TYPE_ASSISTANT_ID, L4_ROLE_ASSISTANT_TEACHER_ID},
}

// SeedDefaultSchoolInvitationRoles siembra las 6 equivalencias tipo→rol por
// defecto para schoolID en academic.school_invitation_roles.
//
// Es el PUNTO ÚNICO de la equivalencia por escuela: lo invoca tanto la capa L4
// (para la escuela demo L1) como common.SeedSchool (para toda escuela de
// playground), de modo que ninguna escuela nazca sin sus equivalencias y sin
// duplicar el mapeo (shared over inline, una fuente un punto).
//
// PRECONDICIÓN: academic.invitation_types debe estar sembrado (ApplyInvitationTypes,
// que corre en la capa L4 del system seed, ANTES que cualquier playground). Si la
// FK invitation_type_id no existe, la inserción falla — el error NO se silencia.
//
// Idempotente: id derivado SHA1(school:type) para que reaplicar no duplique aun
// cuando el unique compuesto (school_id, invitation_type_id) ya exista.
func SeedDefaultSchoolInvitationRoles(tx *gorm.DB, schoolID uuid.UUID) error {
	rows := make([]entities.SchoolInvitationRole, 0, len(invitationTypeToRole))
	for _, m := range invitationTypeToRole {
		typeID, err := uuid.Parse(m.typeIDStr)
		if err != nil {
			return fmt.Errorf("SeedDefaultSchoolInvitationRoles: parse invitation_type id %s: %w", m.typeIDStr, err)
		}
		roleID, err := uuid.Parse(m.roleIDStr)
		if err != nil {
			return fmt.Errorf("SeedDefaultSchoolInvitationRoles: parse role id %s: %w", m.roleIDStr, err)
		}
		derived := uuid.NewSHA1(uuid.NameSpaceOID, []byte("school_invitation_role:"+schoolID.String()+":"+typeID.String()))
		rows = append(rows, entities.SchoolInvitationRole{
			ID:               derived,
			SchoolID:         schoolID,
			InvitationTypeID: typeID,
			IAMRoleID:        roleID,
		})
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&rows).Error
}

// ApplyDemoSchoolInvitationRoles aplica las equivalencias por defecto a la
// escuela demo L1. Se invoca desde la capa L4 DESPUÉS de ApplyInvitationTypes.
func ApplyDemoSchoolInvitationRoles(tx *gorm.DB) error {
	schoolID, err := uuid.Parse(l1SchoolDemoID)
	if err != nil {
		return fmt.Errorf("ApplyDemoSchoolInvitationRoles: parse school id %s: %w", l1SchoolDemoID, err)
	}
	return SeedDefaultSchoolInvitationRoles(tx, schoolID)
}
