package l4

// Constantes UUID del bloque MP-08 (acceso por sistema + tipos de invitación).
//
// MP-08 modela en datos dos catálogos globales que antes vivían como
// enumeraciones hardcodeadas:
//
//   - iam.systems / iam.system_roles: qué roles IAM entran a cada app del
//     ecosistema (kmp, admin-tool). El reparto se resuelve por la tabla puente
//     system_roles, NO por nombres hardcodeados en código.
//   - academic.invitation_types: catálogo de tipos de invitación
//     (teacher/student/guardian/coordinator/admin/assistant). La equivalencia
//     tipo→rol IAM por escuela vive en academic.school_invitation_roles.
//
// Política de prefijos (igual que el resto de L4): UUIDs propios estables,
// nunca reutilizados aunque una fila se elimine. El segundo bloque codifica
// MP-08 (`0008`):
//   - systems          → f8000000-0008-0001-...
//   - invitation_types → f9000000-0008-0001-...
//
// Prefijos f8/f9 elegidos por estar libres en el catálogo de seeds del sistema
// (b0..b4 = capas, c2..c5 = concept types, e0..f2 = permisos L4 por dominio).

// -----------------------------------------------------------------
// iam.systems — 2 apps del ecosistema.
// -----------------------------------------------------------------
const (
	L4_SYSTEM_KMP_ID         = "f8000000-0008-0001-0000-000000000001"
	L4_SYSTEM_KMP_KEY        = "kmp"
	L4_SYSTEM_ADMIN_TOOL_ID  = "f8000000-0008-0001-0000-000000000002"
	L4_SYSTEM_ADMIN_TOOL_KEY = "admin-tool"
	// Plan 025 (mensajería WhatsApp): app de mensajería del ecosistema. La API
	// edugo-api-messaging autoriza por los grants del JWT (no consulta IAM); esta
	// fila de iam.systems existe para que la web pública/admin reconozcan el
	// system vía iam.system_roles (puente sistema↔rol), igual que kmp/admin-tool.
	L4_SYSTEM_MESSAGING_ID  = "f8000000-0008-0001-0000-000000000003"
	L4_SYSTEM_MESSAGING_KEY = "messaging"
)

// -----------------------------------------------------------------
// academic.invitation_types — catálogo global de 6 tipos.
//
// requires_unit = true para los tipos que operan a nivel de unidad académica
// (teacher/student/guardian/assistant); false para los school-scoped
// (coordinator/admin). El label en español es el texto visible al invitar.
// guardian.label = "Representante" (explícito de negocio, NO "Acudiente").
// -----------------------------------------------------------------
const (
	L4_INVITATION_TYPE_TEACHER_ID      = "f9000000-0008-0001-0000-000000000001"
	L4_INVITATION_TYPE_TEACHER_KEY     = "teacher"
	L4_INVITATION_TYPE_STUDENT_ID      = "f9000000-0008-0001-0000-000000000002"
	L4_INVITATION_TYPE_STUDENT_KEY     = "student"
	L4_INVITATION_TYPE_GUARDIAN_ID     = "f9000000-0008-0001-0000-000000000003"
	L4_INVITATION_TYPE_GUARDIAN_KEY    = "guardian"
	L4_INVITATION_TYPE_COORDINATOR_ID  = "f9000000-0008-0001-0000-000000000004"
	L4_INVITATION_TYPE_COORDINATOR_KEY = "coordinator"
	L4_INVITATION_TYPE_ADMIN_ID        = "f9000000-0008-0001-0000-000000000005"
	L4_INVITATION_TYPE_ADMIN_KEY       = "admin"
	L4_INVITATION_TYPE_ASSISTANT_ID    = "f9000000-0008-0001-0000-000000000006"
	L4_INVITATION_TYPE_ASSISTANT_KEY   = "assistant"
)

// Espejo de role IDs de capas previas (L0 super_admin, L1 announcement_viewer)
// que MP-08 necesita referenciar. Re-pegados como literal + comentario igual
// que en roles_permissions.go:529-530 (l4 NO puede importar `layers`: ciclo).
const (
	l0RoleSuperAdminID         = "10000000-0000-0000-0000-000000000001" // layers.L0_ROLE_SUPER_ADMIN_ID
	l1RoleAnnouncementViewerID = "b1000000-0000-0000-0000-000000000001" // layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID
)
