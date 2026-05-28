package l4

// UUIDs hardcodeados de los recursos sembrados por L4 (B1).
// Prefijo `b4000000-*` por convención de fase (L0 usa `b0xxx`,
// L3 usa `b3xxx`). Ver phase-6-layer-l4/design.md §6.
//
// El FE no hardcodea estos UUIDs (resuelve por `key`), así que se
// generaron de cero respecto al legacy. Los sufijos de 2 dígitos
// siguen el patrón mnemotécnico del legacy (10s=admin, 20s=academic,
// 30s=content, 40s=reports, 50s=screens, 60s=guardian, 70s=audit,
// 80s=concept, 90s=settings, A0+=hidden) sólo para legibilidad
// humana al cross-referenciar el inventario; no hay dependencia
// semántica con el legacy.
//
// Las constantes Key existen para que sub-archivos posteriores
// (`roles_permissions.go`, `resource_screens.go`) puedan referenciar
// los recursos por nombre simbólico sin duplicar strings literales.
const (
	// Raíces del menú
	L4_RESOURCE_DASHBOARD_ID  = "b4000000-0000-0000-0000-000000000001"
	L4_RESOURCE_ADMIN_ID      = "b4000000-0000-0000-0000-000000000002"
	L4_RESOURCE_ACADEMIC_ID   = "b4000000-0000-0000-0000-000000000003"
	L4_RESOURCE_CONTENT_ID    = "b4000000-0000-0000-0000-000000000004"
	L4_RESOURCE_REPORTS_ID    = "b4000000-0000-0000-0000-000000000005"
	L4_RESOURCE_DASHBOARD_KEY = "dashboard"
	L4_RESOURCE_ADMIN_KEY     = "admin"
	L4_RESOURCE_ACADEMIC_KEY  = "academic"
	L4_RESOURCE_CONTENT_KEY   = "content"
	L4_RESOURCE_REPORTS_KEY   = "reports"

	// Hijos de admin
	L4_RESOURCE_USERS_ID             = "b4000000-0000-0000-0000-000000000010"
	L4_RESOURCE_SCHOOLS_ID           = "b4000000-0000-0000-0000-000000000011"
	L4_RESOURCE_ROLES_ID             = "b4000000-0000-0000-0000-000000000012"
	L4_RESOURCE_PERMISSIONS_MGMT_ID  = "b4000000-0000-0000-0000-000000000013"
	L4_RESOURCE_SCREEN_TEMPLATES_ID  = "b4000000-0000-0000-0000-000000000050"
	L4_RESOURCE_SCREEN_INSTANCES_ID  = "b4000000-0000-0000-0000-000000000051"
	L4_RESOURCE_AUDIT_ID             = "b4000000-0000-0000-0000-000000000070"
	L4_RESOURCE_CONCEPT_TYPES_ID     = "b4000000-0000-0000-0000-000000000080"
	L4_RESOURCE_SYSTEM_SETTINGS_ID   = "b4000000-0000-0000-0000-000000000090"
	L4_RESOURCE_USERS_KEY            = "users"
	L4_RESOURCE_SCHOOLS_KEY          = "schools"
	L4_RESOURCE_ROLES_KEY            = "roles"
	L4_RESOURCE_PERMISSIONS_MGMT_KEY = "permissions_mgmt"
	L4_RESOURCE_SCREEN_TEMPLATES_KEY = "screen_templates"
	L4_RESOURCE_SCREEN_INSTANCES_KEY = "screen_instances"
	L4_RESOURCE_AUDIT_KEY            = "audit"
	L4_RESOURCE_CONCEPT_TYPES_KEY    = "concept_types"
	L4_RESOURCE_SYSTEM_SETTINGS_KEY  = "system_settings"

	// Hijos de academic
	L4_RESOURCE_UNITS_ID               = "b4000000-0000-0000-0000-000000000020"
	L4_RESOURCE_MEMBERSHIPS_ID         = "b4000000-0000-0000-0000-000000000021"
	L4_RESOURCE_SUBJECTS_ID            = "b4000000-0000-0000-0000-000000000032"
	L4_RESOURCE_GUARDIAN_RELATIONS_ID  = "b4000000-0000-0000-0000-000000000060"
	L4_RESOURCE_PERIODS_ID             = "b4000000-0000-0000-0000-000000000034"
	L4_RESOURCE_GRADES_ID              = "b4000000-0000-0000-0000-000000000035"
	L4_RESOURCE_ATTENDANCE_ID          = "b4000000-0000-0000-0000-000000000036"
	L4_RESOURCE_SCHEDULES_ID           = "b4000000-0000-0000-0000-000000000037"
	L4_RESOURCE_CALENDAR_ID            = "b4000000-0000-0000-0000-000000000039"
	L4_RESOURCE_UNITS_KEY              = "units"
	L4_RESOURCE_MEMBERSHIPS_KEY        = "memberships"
	// Recurso de menú "Mis materias" del alumno (plan 006, N1.C). Es un recurso
	// de MENÚ separado de `memberships` (roster de unidad para admin/teacher):
	// la pantalla default difiere y el path-prefix del gate de menú no
	// distingue por rol sobre el mismo path. Path propio academic.my_memberships
	// aísla el gate para que el alumno NO vea el item admin "memberships".
	// Reintroducido en N1.7 F1 sobre sesiones. Sufijo …22 (adyacente a
	// memberships …21).
	L4_RESOURCE_MY_MEMBERSHIPS_ID  = "b4000000-0000-0000-0000-000000000022"
	L4_RESOURCE_MY_MEMBERSHIPS_KEY = "my_memberships"
	// Plan 010 (N1.7, ADR 0009): "sesiones de materia" (oferta = materia +
	// seccion + periodo + docente como unidad de inscripcion). Recurso bajo
	// `academic`, scope school (coherente con subjects/memberships). Sufijo
	// …23 (adyacente a memberships …21 / my_memberships …22).
	L4_RESOURCE_SUBJECT_OFFERINGS_ID  = "b4000000-0000-0000-0000-000000000023"
	L4_RESOURCE_SUBJECT_OFFERINGS_KEY = "subject_offerings"
	L4_RESOURCE_SUBJECTS_KEY           = "subjects"
	L4_RESOURCE_GUARDIAN_RELATIONS_KEY = "guardian_relations"
	L4_RESOURCE_PERIODS_KEY            = "periods"
	L4_RESOURCE_GRADES_KEY             = "grades"
	L4_RESOURCE_ATTENDANCE_KEY         = "attendance"
	L4_RESOURCE_SCHEDULES_KEY          = "schedules"
	L4_RESOURCE_CALENDAR_KEY           = "calendar"

	// Onboarding (plan 005, N0.0): invitaciones + solicitudes de ingreso.
	// IDs libres adyacentes a guardian (…60). join_request_approvals es
	// solo un namespace de permisos (API-only, sin pantalla propia): la
	// acción del permiso codifica el rol que se admite.
	L4_RESOURCE_INVITATIONS_ID            = "b4000000-0000-0000-0000-000000000061"
	L4_RESOURCE_JOIN_REQUESTS_ID          = "b4000000-0000-0000-0000-000000000062"
	L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID = "b4000000-0000-0000-0000-000000000063"
	L4_RESOURCE_INVITATIONS_KEY            = "invitations"
	L4_RESOURCE_JOIN_REQUESTS_KEY          = "join_requests"
	L4_RESOURCE_JOIN_REQUEST_APPROVALS_KEY = "join_request_approvals"

	// Hijos de content (materials vive en L3, no en L4).
	L4_RESOURCE_ASSESSMENTS_ID          = "b4000000-0000-0000-0000-000000000031"
	L4_RESOURCE_ASSESSMENTS_STUDENT_ID  = "b4000000-0000-0000-0000-000000000033"
	L4_RESOURCE_ASSESSMENTS_KEY         = "assessments"
	L4_RESOURCE_ASSESSMENTS_STUDENT_KEY = "assessments_student"

	// Hijos de reports
	L4_RESOURCE_PROGRESS_ID  = "b4000000-0000-0000-0000-000000000040"
	L4_RESOURCE_STATS_ID     = "b4000000-0000-0000-0000-000000000041"
	L4_RESOURCE_PROGRESS_KEY = "progress"
	L4_RESOURCE_STATS_KEY    = "stats"

	// Recursos "API-only" (IsMenuVisible=false). No raíces visibles.
	L4_RESOURCE_SCREENS_ID        = "b4000000-0000-0000-0000-000000000052"
	L4_RESOURCE_CONTEXT_ID        = "b4000000-0000-0000-0000-0000000000a0"
	L4_RESOURCE_NOTIFICATIONS_ID  = "b4000000-0000-0000-0000-0000000000b0"
	L4_RESOURCE_MENU_ID           = "b4000000-0000-0000-0000-0000000000c0"
	L4_RESOURCE_SCREENS_KEY       = "screens"
	L4_RESOURCE_CONTEXT_KEY       = "context"
	L4_RESOURCE_NOTIFICATIONS_KEY = "notifications"
	L4_RESOURCE_MENU_KEY          = "menu"

	// Demo Fase 3 (SDUI B7b): recurso CRUD plano "colors" para pilotar
	// `GenericListContract`/`GenericFormContract` sin código Kotlin.
	// Visible en el menú bajo "admin" para que el playground respectivo
	// pueda abrirlo desde el navegador del menú dinámico.
	L4_RESOURCE_COLORS_ID  = "b4000000-0000-0000-0000-0000000000e0"
	L4_RESOURCE_COLORS_KEY = "colors"
)
