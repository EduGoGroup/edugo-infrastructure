package l4

// Constantes UUID propias del bloque B2 (roles + permisos + role_permissions).
//
// Política ADR-6/7 (design §6): L4 genera UUIDs propios. Los IDs del
// legacy (`[archivado pre-Fase-6] data.go`) se usaron como inventario
// semántico pero NO se reutilizan binariamente. Excepción: para los
// role_permissions cuyos `permission_id` pertenecen a permisos ya
// sembrados por L0 (announcements:*) o L3 (materials:read/create/update),
// el FK apunta a los IDs definidos en los `layers.L0_PERM_*` /
// `layers.L3_PERM_*`.
//
// Prefijo de UUIDs L4 (igual que L1/L3): `b4xxxxxx-...`.
// El segundo bloque codifica la categoría:
//   b4000000-0001-... → roles
//   b4000000-0002-... → permissions (resource block dictates 4to byte)
//   b4000000-0003-... → role_permissions (derivados con NewSHA1 — ver
//                       applyL4RolePermissions)

// -----------------------------------------------------------------
// Roles L4 — 10 roles totales (super_admin en L0, announcement_viewer en L1).
//
// PRE-4 (rediseño de permisos EduGo): el rol `platform_admin`
// (L4_ROLE_ADMIN_*) fue ELIMINADO porque solapaba con `super_admin`
// (L0) sin aportar diferencia semántica clara. Cualquier flujo que
// antes requería `platform_admin` ahora se cubre con `super_admin`.
//
//   - 4 canónicos: student, teacher, guardian, school_admin
//   - 6 alias: school_director, school_coordinator, school_assistant,
//     assistant_teacher, observer, readonly_auditor
//
// Los roles alias son consultados literalmente por el front KMP vía
// hasRole(...) en DynamicDashboardScreen.kt:33-44. Heredan permisos
// del rol canónico correspondiente (ver tabla en
// rolePermissionGrants()). `readonly_auditor` aplica un filtro
// adicional: excluye toda permission cuya Action sea create/update/
// delete (o termine en :create/:update/:delete).
// -----------------------------------------------------------------
const (
	L4_ROLE_STUDENT_ID  = "b4000000-0001-0000-0000-000000000001"
	L4_ROLE_TEACHER_ID  = "b4000000-0001-0000-0000-000000000002"
	L4_ROLE_GUARDIAN_ID = "b4000000-0001-0000-0000-000000000003"
	// L4_ROLE_ADMIN_ID / L4_ROLE_ADMIN_NAME eliminados en PRE-4
	// (slot UUID b4000000-0001-0000-0000-000000000004 queda libre,
	// no se reutiliza para preservar la lectura del UUID legacy en
	// audits / dumps de BD ya creadas).
	L4_ROLE_SCHOOL_ADMIN_ID = "b4000000-0001-0000-0000-000000000005"

	L4_ROLE_STUDENT_NAME      = "student"
	L4_ROLE_TEACHER_NAME      = "teacher"
	L4_ROLE_GUARDIAN_NAME     = "guardian"
	L4_ROLE_SCHOOL_ADMIN_NAME = "school_admin"

	// Alias roles — heredan de school_admin
	L4_ROLE_SCHOOL_DIRECTOR_ID      = "b4000000-0001-0000-0000-000000000006"
	L4_ROLE_SCHOOL_DIRECTOR_NAME    = "school_director"
	L4_ROLE_SCHOOL_COORDINATOR_ID   = "b4000000-0001-0000-0000-000000000007"
	L4_ROLE_SCHOOL_COORDINATOR_NAME = "school_coordinator"
	L4_ROLE_SCHOOL_ASSISTANT_ID     = "b4000000-0001-0000-0000-000000000008"
	L4_ROLE_SCHOOL_ASSISTANT_NAME   = "school_assistant"

	// Alias roles — heredan de teacher
	L4_ROLE_ASSISTANT_TEACHER_ID   = "b4000000-0001-0000-0000-000000000009"
	L4_ROLE_ASSISTANT_TEACHER_NAME = "assistant_teacher"
	L4_ROLE_OBSERVER_ID            = "b4000000-0001-0000-0000-00000000000a"
	L4_ROLE_OBSERVER_NAME          = "observer"
	L4_ROLE_READONLY_AUDITOR_ID    = "b4000000-0001-0000-0000-00000000000b"
	L4_ROLE_READONLY_AUDITOR_NAME  = "readonly_auditor"
)

// Las constantes `L4_RESOURCE_*_ID` que B2 necesita para construir
// `permissions.resource_id` viven en `resources_constants.go` (owner:
// B1). El placeholder original con UUIDs `20000000-*` fue removido al
// completarse B1 (que regenera UUIDs propios `b4000000-*` según
// ADR-6 §6).

// -----------------------------------------------------------------
// Permisos L4 — UUIDs deterministicos para los permisos que el
// cross-checker (`make seed-audit-strict`) reportó como
// SLOT_REF_MISSING en TC-5: 7 permisos nuevos referenciados por
// slot_data de pantallas L4 que no estaban sembrados.
//
// Patron: `b4000000-0002-<resource-byte>-0000-<sequence>` donde
// `<resource-byte>` es el sufijo hexa del L4_RESOURCE_*_ID asociado
// (12=roles, 13=permissions_mgmt, 36=attendance, 80=concept_types).
// -----------------------------------------------------------------
const (
	// Poda menú (2026-05-29): L4_PERM_ROLES_* y L4_PERM_PERMISSIONS_MGMT_CREATE_ID
	// eliminados junto con los recursos `roles` y `permissions_mgmt`.
	L4_PERM_CONCEPT_TYPES_CREATE_ID = "b4000000-0002-0080-0000-000000000001"
	L4_PERM_CONCEPT_TYPES_UPDATE_ID = "b4000000-0002-0080-0000-000000000002"
	L4_PERM_ATTENDANCE_UPDATE_ID    = "b4000000-0002-0036-0000-000000000001"
	// Permiso ÚNICO del feature "mis materias" del alumno (resource 22,
	// my_memberships). Cubre visibilidad de menú, slot.permission de la
	// pantalla my-memberships-list y route gate del dato. Vive bajo path
	// propio (academic.my_memberships.*) para no filtrar el item de menú admin
	// "memberships" por el gate path-prefix. Reintroducido en N1.7 F1 sobre
	// sesiones.
	L4_PERM_MY_MEMBERSHIPS_READ_OWN_ID = "b4000000-0002-0022-0000-000000000001"
	// Permiso ÚNICO del feature "mis notas" del alumno (resource 24, my_grades).
	// Cubre visibilidad de menú, slot.permission de la pantalla my-grades-list y
	// route gate del dato. Vive bajo path propio (academic.my_grades.*) para no
	// filtrar el item de menú admin "grades" por el gate path-prefix ni depender
	// del wildcard academic.grades.*. N3 F4 (consulta de notas).
	L4_PERM_MY_GRADES_READ_OWN_ID = "b4000000-0002-0024-0000-000000000001"
	// Permisos del representante (plan 024 F1): vistas `:own` del acudido.
	L4_PERM_MY_WARDS_GRADES_READ_OWN_ID        = "b4000000-0002-0025-0000-000000000001"
	L4_PERM_MY_WARDS_ATTENDANCE_READ_OWN_ID    = "b4000000-0002-0026-0000-000000000001"
	L4_PERM_MY_WARDS_ANNOUNCEMENTS_READ_OWN_ID = "b4000000-0002-0027-0000-000000000001"
	L4_PERM_MY_WARDS_MATERIALS_READ_OWN_ID     = "b4000000-0002-0028-0000-000000000001"
	L4_PERM_MY_WARDS_ASSESSMENTS_READ_OWN_ID   = "b4000000-0002-0029-0000-000000000001"
)
