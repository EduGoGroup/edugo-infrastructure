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
	// Poda menú (2026-05-29): se eliminaron los recursos `roles` (…12) y
	// `permissions_mgmt` (…13); sus UUIDs quedan libres.
	// Poda menú (2026-06-01): se eliminaron los recursos `screen_templates`
	// (…50) y `screen_instances` (…51) — su CRUD de configuración SDUI se
	// reimplementó en el admin-tool de Go; sus UUIDs quedan libres.
	L4_RESOURCE_USERS_ID            = "b4000000-0000-0000-0000-000000000010"
	L4_RESOURCE_SCHOOLS_ID          = "b4000000-0000-0000-0000-000000000011"
	L4_RESOURCE_AUDIT_ID            = "b4000000-0000-0000-0000-000000000070"
	L4_RESOURCE_CONCEPT_TYPES_ID    = "b4000000-0000-0000-0000-000000000080"
	L4_RESOURCE_SYSTEM_SETTINGS_ID  = "b4000000-0000-0000-0000-000000000090"
	L4_RESOURCE_USERS_KEY           = "users"
	L4_RESOURCE_SCHOOLS_KEY         = "schools"
	L4_RESOURCE_AUDIT_KEY           = "audit"
	L4_RESOURCE_CONCEPT_TYPES_KEY   = "concept_types"
	L4_RESOURCE_SYSTEM_SETTINGS_KEY = "system_settings"

	// Hijos de academic
	L4_RESOURCE_UNITS_ID       = "b4000000-0000-0000-0000-000000000020"
	L4_RESOURCE_MEMBERSHIPS_ID = "b4000000-0000-0000-0000-000000000021"
	L4_RESOURCE_SUBJECTS_ID    = "b4000000-0000-0000-0000-000000000032"
	// Poda menú (2026-05-29): se eliminaron los recursos `guardian_relations`
	// (…60), `schedules` (…37) y `calendar` (…39); sus UUIDs quedan libres.
	L4_RESOURCE_PERIODS_ID    = "b4000000-0000-0000-0000-000000000034"
	L4_RESOURCE_GRADES_ID     = "b4000000-0000-0000-0000-000000000035"
	L4_RESOURCE_ATTENDANCE_ID = "b4000000-0000-0000-0000-000000000036"
	// grades_detail (N4 / ADR 0020 — MODO DETALLADO). Recurso del desglose por
	// componente de nota (academic.grade_item). NO es menú-visible (no tiene
	// pantalla propia: el desglose vive embebido en grades / "Mis Notas"); existe
	// para colgar los permisos academic.grades_detail.* (que NO pueden compartir
	// resource_id con `grades` por el unique (resource_id, action)). El grant de
	// estos permisos es CONDICIONAL por perfil de escuela y lo inyecta identity en
	// runtime (F4.5). Sufijo …37 (adyacente a grades …35 / attendance …36; UUID
	// liberado por la poda de `schedules`).
	L4_RESOURCE_GRADES_DETAIL_ID  = "b4000000-0000-0000-0000-000000000037"
	L4_RESOURCE_GRADES_DETAIL_KEY = "grades_detail"
	L4_RESOURCE_UNITS_KEY         = "units"
	L4_RESOURCE_MEMBERSHIPS_KEY   = "memberships"
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
	// `academic`, scope unit (los endpoints exigen unidad activa vía
	// RequireActiveContext → 428; coherente con subjects/memberships). Sufijo
	// …23 (adyacente a memberships …21 / my_memberships …22).
	L4_RESOURCE_SUBJECT_OFFERINGS_ID  = "b4000000-0000-0000-0000-000000000023"
	L4_RESOURCE_SUBJECT_OFFERINGS_KEY = "subject_offerings"
	// Recurso de menú "Mis notas" del alumno (N3 F4, consulta de notas). Espejo
	// de my_memberships: recurso de MENÚ separado de `grades` (CRUD docente) con
	// path propio academic.my_grades para que el gate de menú por path-prefix NO
	// le filtre el item admin "grades" ni dependa del wildcard academic.grades.*.
	// Sufijo …24 (adyacente a subject_offerings …23 / my_memberships …22).
	L4_RESOURCE_MY_GRADES_ID   = "b4000000-0000-0000-0000-000000000024"
	L4_RESOURCE_MY_GRADES_KEY  = "my_grades"
	L4_RESOURCE_SUBJECTS_KEY   = "subjects"
	L4_RESOURCE_PERIODS_KEY    = "periods"
	L4_RESOURCE_GRADES_KEY     = "grades"
	L4_RESOURCE_ATTENDANCE_KEY = "attendance"

	// Onboarding (plan 005, N0.0): invitaciones + solicitudes de ingreso.
	// IDs libres adyacentes a guardian (…60). join_request_approvals es
	// solo un namespace de permisos (API-only, sin pantalla propia): la
	// acción del permiso codifica el rol que se admite.
	L4_RESOURCE_INVITATIONS_ID             = "b4000000-0000-0000-0000-000000000061"
	L4_RESOURCE_JOIN_REQUESTS_ID           = "b4000000-0000-0000-0000-000000000062"
	L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID  = "b4000000-0000-0000-0000-000000000063"
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

	// Poda menú (2026-05-29): se eliminó el recurso demo `colors` (…e0); su
	// UUID queda libre. La pareja colors-list/form ya estaba retirada.
)
