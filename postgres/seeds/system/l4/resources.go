package l4

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplyResources siembra los recursos restantes del menú del sistema
// (31 filas). Excluye `announcements` (sembrado en L0) y `materials`
// (sembrado en L3).
//
// Estrategia (F6-REQ-2.x, ADR-6):
//   - El inventario `[archivado pre-Fase-6] data.go::resourceSeedRows`
//     (33 filas) se usa SOLO como guía. NO se hace copy-paste literal.
//   - L4 sólo siembra recursos referenciados por al menos un permiso
//     usado por los 5 roles canónicos implementados (super_admin,
//     school_admin, teacher, student, guardian), o que son contenedores
//     padre del árbol del menú (admin, academic, content). Cualquier
//     recurso huérfano respecto a esos roles se descarta.
//     (PRE-4: el rol `platform_admin` fue eliminado del catálogo y sus
//     capacidades globales quedan cubiertas por `super_admin`.)
//   - Resultado: las 33 filas del legacy se mantienen, menos 2 que
//     viven en L0/L3 (announcements, materials). No hubo descartes
//     adicionales: TODOS los recursos legacy son referenciados por al
//     menos uno de los 6 roles implementados (verificado contra
//     `rolePermissionSeedRows` + `permissionSeedRows`).
//   - UUIDs propios bajo prefijo b4000000-* (ADR-6 §6): el FE no
//     hardcodea UUIDs de recursos (resuelve por `key`), así que se
//     puede regenerar sin coordinación.
//   - Iconos normalizados: el legacy mezclaba estilos Material 2
//     (`settings_applications`, `swap_horiz`) con Lucide (`book-open`,
//     `user-plus`). L4 normaliza a Lucide (kebab-case) por
//     consistencia con L0 (`bullhorn`) y L3 (`book`). Excepción:
//     `dashboard` que ya está en kebab/snake equivalente.
//   - Keys preservadas: `permissions_mgmt` se conserva tal cual porque
//     el FE la referencia literalmente
//     (kmp-screens/.../PermissionsListContract.kt). `assessments_student`
//     se conserva como recurso separado de `assessments` porque tiene
//     una permission distinta (`assessments_student:read`) usada por
//     student/teacher para el flujo "tomar evaluación".
//
// Idempotencia: UPSERT vía ON CONFLICT (id) — mismo patrón que
// applyL0Resources / applyL3Resources. Las columnas booleanas con
// `default` tag se setean por SQL crudo para evitar el bug GORM de
// zero-value en bool.
//
// Orden de inserción: padres antes que hijos para evitar violación FK
// `fk_resources_parent` durante la primera aplicación.
func ApplyResources(tx *gorm.DB) error {
	const upsertSQL = `
        INSERT INTO iam.resources
            (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope, is_active, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?::iam.permission_scope, ?, NOW(), NOW())
        ON CONFLICT (id) DO UPDATE SET
            key             = EXCLUDED.key,
            display_name    = EXCLUDED.display_name,
            description     = EXCLUDED.description,
            icon            = EXCLUDED.icon,
            parent_id       = EXCLUDED.parent_id,
            sort_order      = EXCLUDED.sort_order,
            is_menu_visible = EXCLUDED.is_menu_visible,
            scope           = EXCLUDED.scope,
            is_active       = EXCLUDED.is_active
    `

	for _, r := range l4Resources {
		id, err := uuid.Parse(r.ID)
		if err != nil {
			return fmt.Errorf("ApplyResources: parse id %s: %w", r.ID, err)
		}
		var parentID any
		if r.ParentID != "" {
			pid, perr := uuid.Parse(r.ParentID)
			if perr != nil {
				return fmt.Errorf("ApplyResources: parse parent_id %s for %s: %w", r.ParentID, r.Key, perr)
			}
			parentID = pid
		}

		var description any
		if r.Description != "" {
			description = r.Description
		}
		var icon any
		if r.Icon != "" {
			icon = r.Icon
		}

		if err := tx.Exec(upsertSQL,
			id,
			r.Key,
			r.DisplayName,
			description,
			icon,
			parentID,
			r.SortOrder,
			r.IsMenuVisible,
			r.Scope,
			r.IsActive,
		).Error; err != nil {
			return fmt.Errorf("ApplyResources: upsert %s: %w", r.Key, err)
		}
	}

	// ----------------------------------------------------------------
	// Re-parent recursos sembrados en capas previas (L0/L3) bajo los
	// contenedores `academic` / `content` introducidos en L4.
	//
	// Decisión (PRE-4, sub-tarea B — rediseño de permisos EduGo):
	// L0/L3 crean `announcements` y `materials` con `parent_id=NULL`
	// porque sus respectivos contenedores no existían aún. Cuando L4
	// siembra `academic` y `content` (arriba en este mismo Apply), hay
	// que enlazar los recursos huérfanos al nuevo árbol.
	//
	// Idempotente: UPDATE incondicional sobre id conocidos. Si la fila
	// no existe (caso teórico: aplicar L4 sin L0/L3 previamente, lo
	// cual no es soportado por system.Layers()) la UPDATE no afecta
	// filas y no rompe — el FK fk_resources_parent se garantiza porque
	// `academic` (b4000000-…-03) y `content` (b4000000-…-04) ya
	// fueron upserteados por el loop anterior.
	const updateParentSQL = `
        UPDATE iam.resources
        SET parent_id = ?::uuid, updated_at = NOW()
        WHERE id = ?::uuid
    `
	reparents := []struct {
		childID    string
		parentID   string
		humanLabel string
	}{
		{
			childID:    "b0000000-0000-0000-0000-000000000001", // = layers.L0_RESOURCE_ANNOUNCEMENTS_ID
			parentID:   L4_RESOURCE_ACADEMIC_ID,
			humanLabel: "announcements → academic",
		},
		{
			childID:    "b3000000-0000-0000-0000-000000000001", // = layers.L3_RESOURCE_MATERIALS_ID
			parentID:   L4_RESOURCE_CONTENT_ID,
			humanLabel: "materials → content",
		},
	}
	for _, rp := range reparents {
		if err := tx.Exec(updateParentSQL, rp.parentID, rp.childID).Error; err != nil {
			return fmt.Errorf("ApplyResources: re-parent %s: %w", rp.humanLabel, err)
		}
	}

	return nil
}

// l4ResourceRow es la representación local de una fila de
// iam.resources para L4. Usar string vacío en lugar de *string para
// Description/Icon/ParentID simplifica la declaración del slice; la
// conversión a NULL ocurre en ApplyResources.
type l4ResourceRow struct {
	ID            string
	Key           string
	DisplayName   string
	Description   string
	Icon          string
	ParentID      string
	SortOrder     int
	IsMenuVisible bool
	Scope         string
	IsActive      bool
}

// l4Resources es la definición final de los 31 recursos sembrados por
// L4. Orden: padres antes que hijos para preservar FK al aplicar la
// primera vez. Decisiones de criterio documentadas inline.
var l4Resources = []l4ResourceRow{
	// -------------------------------------------------------------
	// Raíces del menú (sin parent_id, visibles)
	// -------------------------------------------------------------
	{ID: L4_RESOURCE_DASHBOARD_ID, Key: "dashboard", DisplayName: "Dashboard", Description: "Panel principal", Icon: "dashboard", Scope: "system", ParentID: "", SortOrder: 1, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_ADMIN_ID, Key: "admin", DisplayName: "Administración", Description: "Módulo de administración", Icon: "settings", Scope: "system", ParentID: "", SortOrder: 2, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_ACADEMIC_ID, Key: "academic", DisplayName: "Académico", Description: "Módulo académico", Icon: "graduation-cap", Scope: "school", ParentID: "", SortOrder: 3, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_CONTENT_ID, Key: "content", DisplayName: "Contenido", Description: "Contenido educativo", Icon: "book-open", Scope: "unit", ParentID: "", SortOrder: 4, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_REPORTS_ID, Key: "reports", DisplayName: "Reportes", Description: "Reportes y estadísticas", Icon: "bar-chart", Scope: "school", ParentID: "", SortOrder: 5, IsMenuVisible: true, IsActive: true},

	// -------------------------------------------------------------
	// Hijos de "admin" (gestión y operación del sistema)
	// -------------------------------------------------------------
	{ID: L4_RESOURCE_USERS_ID, Key: "users", DisplayName: "Usuarios", Description: "Gestión de usuarios", Icon: "users", Scope: "school", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 1, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_SCHOOLS_ID, Key: "schools", DisplayName: "Escuelas", Description: "Gestión de escuelas", Icon: "school", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 2, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_ROLES_ID, Key: "roles", DisplayName: "Roles", Description: "Gestión de roles", Icon: "shield", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 3, IsMenuVisible: true, IsActive: true},
	// Key `permissions_mgmt` preservada (no renombrar): el FE la
	// referencia literalmente en `PermissionsListContract.kt` y
	// `PermissionsFormContract.kt`. Renombrarla rompería el FE.
	{ID: L4_RESOURCE_PERMISSIONS_MGMT_ID, Key: "permissions_mgmt", DisplayName: "Permisos", Description: "Gestión de permisos", Icon: "key", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 4, IsMenuVisible: true, IsActive: true},
	// Icono normalizado de "settings_applications" (Material 2) a
	// "settings-2" (Lucide) por consistencia con L0/L3.
	{ID: L4_RESOURCE_SCREEN_TEMPLATES_ID, Key: "screen_templates", DisplayName: "Templates de Pantalla", Description: "Templates base para configuración de pantallas", Icon: "settings-2", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 5, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_SCREEN_INSTANCES_ID, Key: "screen_instances", DisplayName: "Instancias de Pantalla", Description: "Instancias configuradas de pantalla por escuela", Icon: "monitor-smartphone", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 6, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_AUDIT_ID, Key: "audit", DisplayName: "Auditoría", Description: "Registro de auditoría del sistema", Icon: "file-search", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 7, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_CONCEPT_TYPES_ID, Key: "concept_types", DisplayName: "Tipos de Concepto", Description: "Tipos de institución y terminología", Icon: "tag", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 8, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_SYSTEM_SETTINGS_ID, Key: "system_settings", DisplayName: "Configuración", Description: "Configuración y mantenimiento del sistema", Icon: "settings", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 9, IsMenuVisible: true, IsActive: true},
	// Demo Fase 3 (SDUI B7b): recurso CRUD plano para validar el
	// fallback genérico data-driven sin escribir Kotlin nuevo.
	// Permisos `platform.colors.*` (edugo-shared). Endpoint
	// `/api/v1/colors` en edugo-api-platform. Icono Lucide `palette`.
	{ID: L4_RESOURCE_COLORS_ID, Key: "colors", DisplayName: "Colores", Description: "Demo CRUD data-driven (Fase 3)", Icon: "palette", Scope: "system", ParentID: L4_RESOURCE_ADMIN_ID, SortOrder: 10, IsMenuVisible: true, IsActive: true},

	// -------------------------------------------------------------
	// Hijos de "academic" (gestión académica y comunicaciones)
	// -------------------------------------------------------------
	{ID: L4_RESOURCE_UNITS_ID, Key: "units", DisplayName: "Unidades Académicas", Description: "Gestión de clases", Icon: "layers", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 1, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_MEMBERSHIPS_ID, Key: "memberships", DisplayName: "Miembros", Description: "Asignación de miembros", Icon: "user-plus", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 2, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_SUBJECTS_ID, Key: "subjects", DisplayName: "Materias", Description: "Gestión de materias", Icon: "book", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 3, IsMenuVisible: true, IsActive: true},
	// Plan 010 (N1.7, ADR 0009): sesiones de materia. Recurso de permisos
	// bajo "academic". IsMenuVisible=false en esta etapa de esquema (F0a):
	// aún no hay screen_instance ni mapping en resource_screens, así que no
	// debe aparecer en el menú como item muerto; el item de menú/pantalla
	// se siembra en una etapa posterior del plan 010.
	{ID: L4_RESOURCE_SUBJECT_OFFERINGS_ID, Key: "subject_offerings", DisplayName: "Sesiones de Materia", Description: "Oferta de materia: sección, período y docente", Icon: "book", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 14, IsMenuVisible: false, IsActive: true},
	{ID: L4_RESOURCE_GUARDIAN_RELATIONS_ID, Key: "guardian_relations", DisplayName: "Vínculos Guardian", Description: "Gestión de vínculos guardian-estudiante", Icon: "user-check", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 4, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_PERIODS_ID, Key: "periods", DisplayName: "Periodos Académicos", Description: "Gestión de periodos académicos", Icon: "calendar-range", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 5, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_GRADES_ID, Key: "grades", DisplayName: "Calificaciones", Description: "Gestión de calificaciones", Icon: "award", Scope: "unit", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 6, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_ATTENDANCE_ID, Key: "attendance", DisplayName: "Asistencia", Description: "Registro de asistencia", Icon: "check-square", Scope: "unit", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 7, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_SCHEDULES_ID, Key: "schedules", DisplayName: "Horarios", Description: "Horarios semanales por unidad", Icon: "clock", Scope: "unit", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 8, IsMenuVisible: true, IsActive: true},
	// SortOrder=10 en el legacy: el 9 lo ocupaba "announcements" que
	// vive en L0. Se preserva el gap para que el orden visual del
	// menú no cambie cuando L0+L4 conviven.
	{ID: L4_RESOURCE_CALENDAR_ID, Key: "calendar", DisplayName: "Calendario", Description: "Calendario escolar", Icon: "calendar", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 10, IsMenuVisible: true, IsActive: true},
	// Onboarding (plan 005, N0.0): invitaciones y solicitudes de ingreso.
	// Menu-visibles bajo "academic" como guardian_relations.
	{ID: L4_RESOURCE_INVITATIONS_ID, Key: "invitations", DisplayName: "Invitaciones", Description: "Códigos de invitación a colegio/unidad", Icon: "ticket", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 11, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_JOIN_REQUESTS_ID, Key: "join_requests", DisplayName: "Solicitudes de Ingreso", Description: "Bandeja de solicitudes de ingreso pendientes", Icon: "user-plus", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 12, IsMenuVisible: true, IsActive: true},
	// "Mis materias" (recurso my_memberships, plan 006) — item de menú del
	// alumno que abre la lista de materias en las que está inscrito. Recurso
	// separado de `memberships` (roster de unidad para admin/teacher): la
	// pantalla default y el gate de visibilidad difieren. Path
	// `academic.my_memberships`; solo lo tocan grants que matcheen ese path →
	// student (grant dedicado academic.my_memberships.read:own) y school_admin
	// (academic.*). El teacher tiene el literal `academic.memberships.read`
	// que NO toca este path, así que no lo ve. Scope=unit (el alumno lee dentro
	// de su unidad activa). Reintroducido en N1.7 F1 sobre el modelo de sesiones
	// (lector A: GET /api/v1/me/subject-offerings).
	{ID: L4_RESOURCE_MY_MEMBERSHIPS_ID, Key: "my_memberships", DisplayName: "Mis Materias", Description: "Materias en las que el alumno está inscrito", Icon: "book", Scope: "unit", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 13, IsMenuVisible: true, IsActive: true},

	// -------------------------------------------------------------
	// Hijos de "content" (evaluaciones y materiales)
	// -------------------------------------------------------------
	// NOTA: `materials` (sort_order=1 en legacy) vive en L3. Se
	// preserva el gap dejando assessments en sort_order=2.
	{ID: L4_RESOURCE_ASSESSMENTS_ID, Key: "assessments", DisplayName: "Evaluaciones", Description: "Evaluaciones y exámenes", Icon: "clipboard", Scope: "unit", ParentID: L4_RESOURCE_CONTENT_ID, SortOrder: 2, IsMenuVisible: true, IsActive: true},
	// Recurso separado de `assessments` (no renombrar): la
	// permission `assessments_student:read` distingue la vista de
	// estudiante (rendir) de la del docente (configurar/calificar).
	// Renombrar rompería esta semántica.
	{ID: L4_RESOURCE_ASSESSMENTS_STUDENT_ID, Key: "assessments_student", DisplayName: "Tomar Evaluación", Description: "Vista de evaluaciones desde perspectiva del estudiante", Icon: "play-circle", Scope: "unit", ParentID: L4_RESOURCE_CONTENT_ID, SortOrder: 3, IsMenuVisible: true, IsActive: true},

	// -------------------------------------------------------------
	// Hijos de "reports"
	// -------------------------------------------------------------
	{ID: L4_RESOURCE_PROGRESS_ID, Key: "progress", DisplayName: "Progreso", Description: "Seguimiento de progreso", Icon: "trending-up", Scope: "unit", ParentID: L4_RESOURCE_REPORTS_ID, SortOrder: 1, IsMenuVisible: true, IsActive: true},
	{ID: L4_RESOURCE_STATS_ID, Key: "stats", DisplayName: "Estadísticas", Description: "Estadísticas del sistema", Icon: "pie-chart", Scope: "school", ParentID: L4_RESOURCE_REPORTS_ID, SortOrder: 2, IsMenuVisible: true, IsActive: true},

	// -------------------------------------------------------------
	// Recursos "API-only" (IsMenuVisible=false): no aparecen en el
	// menú pero existen como targets de permisos del backend
	// (resolución de pantallas, contexto activo, notificaciones,
	// menú dinámico).
	// -------------------------------------------------------------
	{ID: L4_RESOURCE_SCREENS_ID, Key: "screens", DisplayName: "Pantallas (Mobile)", Description: "Lectura de pantallas desde aplicación mobile", Icon: "smartphone", Scope: "system", ParentID: "", SortOrder: 0, IsMenuVisible: false, IsActive: true},
	// Icono normalizado de "swap_horiz" (Material 2) a "arrow-left-right"
	// (Lucide) por consistencia.
	{ID: L4_RESOURCE_CONTEXT_ID, Key: "context", DisplayName: "Contexto", Description: "Exploración de escuelas y unidades para selección de contexto", Icon: "arrow-left-right", Scope: "system", ParentID: "", SortOrder: 99, IsMenuVisible: false, IsActive: true},
	{ID: L4_RESOURCE_NOTIFICATIONS_ID, Key: "notifications", DisplayName: "Notificaciones", Description: "Centro de notificaciones del usuario", Icon: "bell", Scope: "system", ParentID: "", SortOrder: 100, IsMenuVisible: false, IsActive: true},
	{ID: L4_RESOURCE_MENU_ID, Key: "menu", DisplayName: "Menu", Description: "Navegación y menu de la aplicación", Icon: "menu", Scope: "system", ParentID: "", SortOrder: 101, IsMenuVisible: false, IsActive: true},
	// Onboarding (plan 005, N0.0): namespace de permisos de aprobación
	// per-rol. API-only (sin pantalla): la acción del permiso ES el rol que
	// se admite (academic.join_request_approvals.{student,teacher,guardian}).
	// Padre `academic` para coherencia de dominio.
	{ID: L4_RESOURCE_JOIN_REQUEST_APPROVALS_ID, Key: "join_request_approvals", DisplayName: "Aprobación de Solicitudes", Description: "Permiso de aprobación de solicitudes de ingreso por rol", Icon: "user-check", Scope: "school", ParentID: L4_RESOURCE_ACADEMIC_ID, SortOrder: 102, IsMenuVisible: false, IsActive: true},
}
