package l4

// Constructores por screen_instance.
// =====================================
//
// Cada funcion devuelve un `l4ScreenInstanceRow` ya armado. Se
// separan del slice principal (`l4ScreenInstanceRows()` en
// screen_instances.go) para que los JSON crudos no inflen ese
// archivo y para que cada decision quede co-localizada con sus
// columns/fields/acciones.
//
// El template del slot_data NO es prescriptivo: cada constructor
// escoge los campos minimos para que la pantalla resuelva 200 en el
// endpoint `screen-config/resolve/key/:key` y para que el FE pueda
// renderizarla con su contrato actual. Cualquier campo que el FE
// hardcodea (ej. labels especificos) se preserva si era razonable;
// el resto se simplifica.

// ===============================================================
// ADMIN: users / schools / roles / permissions
// ===============================================================

func usersList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_USERS_LIST_ID,
		screenKey:          "users-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Lista de Usuarios",
		description:        "Listado de usuarios del sistema",
		scope:              "system",
		requiredPermission: "admin.users.read",
		slotData: `{
  "title": "Usuarios",
  "search_placeholder": "Buscar usuario...",
  "filter_all_label": "Todos",
  "filter_ready_label": "Activos",
  "filter_processing_label": "Inactivos",
  "columns": [
    {"key": "full_name", "label": "Nombre"},
    {"key": "email", "label": "Email"},
    {"key": "is_active", "label": "Activo"}
  ],
  "api_prefix": "identity"
}`,
	}
}

func usersForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_USERS_FORM_ID,
		screenKey:          "users-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Usuario",
		description:        "Crear o editar un usuario",
		scope:              "system",
		requiredPermission: "admin.users.read",
		slotData: `{
  "title": "Usuario",
  "fields": [
    {"key": "full_name", "label": "Nombre completo", "type": "text", "required": true},
    {"key": "email", "label": "Email", "type": "email", "required": true},
    {"key": "password", "label": "Contraseña", "type": "password", "required": false},
    {"key": "is_active", "label": "Activo", "type": "toggle"}
  ],
  "api_prefix": "identity"
}`,
	}
}

// schools-list (MP-08 F4, DEC-D; bug 0054): pantalla read-only. Se retiran
// create/edit/delete del header y de las filas (heredados de list-basic-v1 vía
// actions_removed). La gestión real de escuelas (alta/edición/baja) vive en el
// admin-tool de Go, no en el producto SDUI del KMP; la pantalla de Escuelas
// solo lista. Mismo patrón delta que memberships-list / concept-types-list.
func schoolsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCHOOLS_LIST_ID,
		screenKey:          "schools-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Lista de Escuelas",
		description:        "Listado de escuelas",
		scope:              "system",
		requiredPermission: "admin.schools.read",
		slotData: `{
  "title": "Escuelas",
  "search_placeholder": "Buscar escuela...",
  "filter_all_label": "Todos",
  "filter_ready_label": "Activas",
  "filter_processing_label": "Inactivas",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "code", "label": "Código"},
    {"key": "is_active", "label": "Activa"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "academic"
}`,
	}
}

func schoolsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCHOOLS_FORM_ID,
		screenKey:          "schools-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Escuela",
		description:        "Crear o editar una escuela",
		scope:              "system",
		requiredPermission: "admin.schools.read",
		slotData: `{
  "title": "Escuela",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "code", "label": "Código", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "is_active", "label": "Activa", "type": "toggle"}
  ],
  "actions_added": [
    {"id": "manage-concepts", "scope": "form-submit", "label": "Gestionar Conceptos", "icon": "settings", "permission": "admin.concept_types.read", "condition": "edit-only", "event_id": "manage-concepts", "style": "outlined", "order": 30}
  ],
  "api_prefix": "academic"
}`,
	}
}

// Poda menú (2026-05-29): los constructores rolesList(), rolesForm(),
// permissionsList() y permissionsForm() se eliminaron — el FE KMP no
// implementa esas pantallas (roles-*, permissions-*) y los recursos
// `roles`/`permissions_mgmt` fueron retirados del menú. Sus constantes
// L4_SCREEN_INST_ROLES_*/PERMISSIONS_* también se quitaron.

// Poda menú (2026-06-01): los constructores screenTplList(), screenInstList(),
// screenInstForm() y screensForm() se eliminaron — las pantallas de
// configuración SDUI (screen-templates-list, screen-instances-list/form,
// screens-form) se reimplementaron en el admin-tool de Go; los recursos
// `screen_templates`/`screen_instances` se retiraron del menú.

// ===============================================================
// ADMIN: system-settings + concept-types
// ===============================================================

func systemSettings() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SYSTEM_SETTINGS_ID,
		screenKey:          "system-settings",
		templateID:         l4TplSettingsSystemV1ID,
		name:               "Configuración del Sistema",
		description:        "Configuración global del sistema",
		scope:              "system",
		requiredPermission: "admin.system_settings.read",
		slotData: `{
  "title": "Configuración del Sistema",
  "cache_title": "Cache",
  "cache_description": "Limpia cachés locales y remotos para forzar refresh",
  "clear_cache_label": "Limpiar cache",
  "info_title": "Información",
  "app_version_label": "Versión de la app",
  "app_version_value": "1.0.0",
  "schema_version_label": "Versión del schema",
  "schema_version_value": "see /admin/version",
  "api_prefix": "platform"
}`,
	}
}

func conceptTypesList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_CONCEPT_TYPES_LIST_ID,
		screenKey:          "concept-types-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Tipos de Concepto",
		description:        "Catálogo de tipos de institución y terminología",
		scope:              "system",
		requiredPermission: "admin.concept_types.read",
		slotData: `{
  "title": "Tipos de Concepto",
  "search_placeholder": "Buscar tipo...",
  "filter_all_label": "Todos",
  "filter_ready_label": "Activos",
  "filter_processing_label": "Inactivos",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "code", "label": "Código"},
    {"key": "is_active", "label": "Activo"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

func conceptTypesForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_CONCEPT_TYPES_FORM_ID,
		screenKey:          "concept-types-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Tipo de Concepto",
		description:        "Crear o editar un tipo de concepto",
		scope:              "system",
		requiredPermission: "admin.concept_types.read",
		slotData: `{
  "title": "Tipo de Concepto",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "code", "label": "Código", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "is_active", "label": "Activo", "type": "toggle"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// ADMIN: AUDIT
// ===============================================================

func auditEventsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_AUDIT_LIST_ID,
		screenKey:          "audit-events-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Eventos de Auditoría",
		description:        "Listado de eventos de auditoría del sistema",
		scope:              "system",
		requiredPermission: "admin.audit.read",
		slotData: `{
  "title": "Auditoría",
  "search_placeholder": "Buscar evento...",
  "filter_all_label": "Todos",
  "filter_processing_label": "Solo críticos",
  "columns": [
    {"key": "action", "label": "Acción"},
    {"key": "actor_email", "label": "Actor"},
    {"key": "resource_type", "label": "Recurso"},
    {"key": "created_at", "label": "Fecha"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "identity"
}`,
	}
}

// auditDetail: detalle SOLO LECTURA de un evento de auditoría. Usa el
// template propio `audit-detail-v1` (L4) en vez del base detail-basic-v1
// (L0): este último trae slots HARDCODEADOS de material/archivo
// ("Tamaño/Subido/Estado/Descripción" + botón "Descargar") y el renderer
// de detalle del KMP está atado a las `zones` del template — el slot_data
// del instance solo sustituye labels (bind "slot:<key>"), no cambia qué
// `field` del JSON se pinta ni los slots. audit-detail-v1 declara los
// campos reales del evento (GET identity:/api/v1/audit/events/:id) con
// labels en español y sin descarga. Endpoint y permiso (admin.audit.read)
// intactos. Cada fila del template son dos slots controlType "label"
// (etiqueta estática + valor desde `field`); aquí en el instance solo va el
// título de la pantalla.
func auditDetail() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_AUDIT_DETAIL_ID,
		screenKey:          "audit-detail",
		templateID:         l4TplAuditDetailV1ID,
		name:               "Detalle de Auditoría",
		description:        "Detalle de un evento de auditoría",
		scope:              "system",
		requiredPermission: "admin.audit.read",
		slotData: `{
  "title": "Detalle de auditoría",
  "page_title": "Detalle de auditoría",
  "readonly": true,
  "api_prefix": "identity"
}`,
	}
}

// ===============================================================
// ACADEMIC: units / memberships / subjects / periods
// ===============================================================

func unitsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_UNITS_LIST_ID,
		screenKey:          "units-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Unidades Académicas",
		description:        "Listado de unidades académicas",
		scope:              "school",
		requiredPermission: "academic.units.read",
		slotData: `{
  "title": "Unidades",
  "search_placeholder": "Buscar unidad...",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "level", "label": "Nivel"},
    {"key": "period", "label": "Periodo"}
  ],
  "api_prefix": "academic"
}`,
	}
}

func unitsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_UNITS_FORM_ID,
		screenKey:          "units-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Unidad",
		description:        "Crear o editar una unidad académica",
		scope:              "school",
		requiredPermission: "academic.units.read",
		slotData: `{
  "title": "Unidad",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "type", "label": "Tipo", "type": "select", "required": true, "options": [
      {"value": "school", "label": "Colegio"},
      {"value": "grade", "label": "Grado"},
      {"value": "class", "label": "Clase"},
      {"value": "section", "label": "Sección"},
      {"value": "club", "label": "Club"},
      {"value": "department", "label": "Departamento"}
    ]},
    {"key": "parent_unit_id", "label": "Unidad Padre", "type": "entity-picker", "required": false, "remote_endpoint": "academic:/api/v1/units", "display_field": "display_name", "value_field": "id", "search_param": "search", "page_size": 20, "picker_title": "Buscar unidad padre"}
  ],
  "api_prefix": "academic"
}`,
	}
}

// membershipsList: hereda los default_actions de list-basic-v1 pero RETIRA
// "create" — la creación directa de membresías se eliminó (redundante con el
// flujo invitación→solicitud→doble-gate→aprobación, que ya crea la membresía).
// Las acciones edit/delete/expire se conservan.
func membershipsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MEMBERSHIPS_LIST_ID,
		screenKey:          "memberships-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Miembros",
		description:        "Listado de miembros por unidad",
		scope:              "school",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Miembros",
  "search_placeholder": "Buscar miembro...",
  "columns": [
    {"key": "user_name", "label": "Usuario"},
    {"key": "unit_name", "label": "Unidad"},
    {"key": "role", "label": "Rol"}
  ],
  "actions_removed": ["create"],
  "api_prefix": "academic"
}`,
	}
}

// membershipsForm: form-basic-v1 reservado para SOLO EDICIÓN de una membresía
// existente. La creación directa de membresías se retiró (redundante con el flujo
// invitación→solicitud→doble-gate→aprobación): no hay FAB de crear en
// memberships-list, no hay POST en el backend y membership-add se eliminó. Esta
// pantalla solo se alcanza desde la acción "editar" de la lista; carga por id
// (LOAD_DATA → GET /memberships/:id) y guarda con PUT.
//   - actions_removed=["save_new"]: retira el "guardar como nuevo" (action
//     save_new, condition=create-only, permission $resource$.create) heredado del
//     template form-basic-v1; queda solo `save` (condition=edit-only,
//     $resource$.update → PUT) y `delete`. Así la pantalla NUNCA puede crear.
//   - user_email (text): el usuario NO se reasigna editando; el contrato KMP lo
//     muestra read-only en edición. Las keys/tipos cuadran con el contrato real.
//   - academic_unit_id (remote_select): el FormFieldsResolver del KMP DESCARTA
//     todo remote_select sin remote_endpoint, así que aquí SÍ lo declaramos.
//     Endpoint academic:/api/v1/units → {"units":[{id, display_name,...}]}; la
//     escuela se resuelve de la escuela activa del JWT (NUNCA por path/query/
//     body, estándar del ecosistema). display_field=display_name, value_field=id.
//   - role_key (select estático): enum del backend (NO remote, NO role_id).
//
// NO lleva subject_ids ni materias (retirado en F0b, no se reintroduce).
func membershipsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MEMBERSHIPS_FORM_ID,
		screenKey:          "memberships-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Editar Membresía",
		description:        "Editar la membresía de un usuario en una unidad",
		scope:              "school",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Editar Membresía",
  "fields": [
    {"key": "user_email", "label": "Email del usuario", "type": "text", "required": true},
    {"key": "academic_unit_id", "label": "Unidad", "type": "remote_select", "required": true, "remote_endpoint": "academic:/api/v1/units", "display_field": "display_name", "value_field": "id"},
    {"key": "role_key", "label": "Rol", "type": "select", "required": true, "options": [
      {"value": "student", "label": "Estudiante"},
      {"value": "teacher", "label": "Profesor"},
      {"value": "guardian", "label": "Acudiente"},
      {"value": "assistant", "label": "Docente auxiliar"},
      {"value": "coordinator", "label": "Coordinador"},
      {"value": "admin", "label": "Administrador"}
    ]}
  ],
  "actions_removed": ["save_new"],
  "api_prefix": "academic"
}`,
	}
}

// subjectsList: hereda los default_actions de list-basic-v1
// (create/edit/delete sobre $resource$ → academic.subjects.*). Sin deltas:
// el patrón CRUD estándar es suficiente. La vista "sesiones por materia" no
// vive aquí sino embebida como pestaña "Sesiones" en subjects-form
// (master-detail con detail_configs), ver subjectsForm().
func subjectsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SUBJECTS_LIST_ID,
		screenKey:          "subjects-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Lista de Materias",
		description:        "Listado de materias",
		scope:              "school",
		requiredPermission: "academic.subjects.read",
		slotData: `{
  "title": "Materias",
  "search_placeholder": "Buscar materia...",
  "filter_all_label": "Todos",
  "filter_ready_label": "Activas",
  "filter_processing_label": "Inactivas",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "code", "label": "Código"}
  ],
  "api_prefix": "academic"
}`,
	}
}

// subjectsForm usa master-detail-v1 (plan 006, Trozo A; N1.7 F2.2/F2.3): hereda
// los 3 defaults de form (save_new/save/delete con scope=form-submit) y, vía
// detail_configs[], embebe UNA pestaña de detalle:
//   - "Sesiones" → sessions-by-subject-list (sesiones/offerings de la materia),
//     CON modal (modal_screen_key="sessions-by-subject-form", N1.7 F2.3): el
//     botón "+" crea una sesión y el click en fila la edita (asignar/cambiar
//     docente, sección, estado). El MasterDetailContainer abre el modal pasando
//     subjectId (parent) y, en edición, id.
//
// La pestaña "Alumnos" (students-by-subject-list) se RETIRÓ: la materia es
// catálogo; un alumno se inscribe en una SESIÓN, no en la materia, así que el
// roster de alumnos se gestiona dentro de cada sesión (batch-enroll/enroll-one),
// no a nivel materia. El detalle de materia queda SOLO con "Sesiones".
//
// La pestaña "Sesiones" sustituye a la antigua row-action `view-sessions` de
// subjects-list (eliminada en F2.2): ahora se llega navegando dentro del
// formulario de materia.
//
// detail_configs: la entrada lleva parent_id_param="subjectId" →
// MasterDetailContainer carga la pantalla hija pasando subjectId = id de la
// materia editada; el contrato KMP lee context.params["subjectId"].
// sessions-by-subject-list llama a
// GET /api/v1/subject-offerings?subject_id=. child_id_field="id". El frontend
// KMP interpreta detail_configs; el backend solo lo persiste.
//
// actions_removed=["detail"]: el template master-detail-v1 trae un default
// `detail` (view-detail|$resource$.read|edit-only) pensado para navegar a un
// detalle full-screen. Aquí el detalle es el panel EMBEBIDO (no hay pantalla
// destino ni handler view-detail en SubjectsFormContract), así que el botón de
// toolbar no aplica y se retira intencionalmente.
//
// Entry-points de asistencia/notas REUBICADOS (N3.5 F1, plan 014 / ADR 0018):
// las 4 acciones del docente — "Pasar lista" (take-attendance), "Historial"
// (view-attendance), "Resumen" (view-attendance-summary) y "Poner notas"
// (put-grades) — YA NO cuelgan de subjects-form. Colgaban de la materia (scope
// resource-toolbar, condition edit-only) y eso mezclaba el roster de un docente
// que dicta la misma materia en dos secciones (A/B). Ahora viven en la card de
// cada SESIÓN, como row-actions de sessions-by-subject-list (scope row): el id de
// la fila es el offering_id, así que cada acción opera sobre una sección concreta.
// Es reubicación, no convivencia: se BORRARON de aquí. Ver sessionsBySubjectList.
//
// Reintroducido en N1.7 F2 sobre el modelo de sesiones (antes de F0b dependía
// del filtro subject_id sobre membership_subjects; ahora el lector resuelve
// las sesiones de la materia).
func subjectsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SUBJECTS_FORM_ID,
		screenKey:          "subjects-form",
		templateID:         L0_SCREEN_TPL_MASTER_DETAIL_ID_REF,
		name:               "Formulario de Materia",
		description:        "Crear o editar una materia",
		scope:              "school",
		requiredPermission: "academic.subjects.read",
		slotData: `{
  "title": "Materia",
  "page_title": "Materia",
  "edit_title": "Editar materia",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "code", "label": "Código", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"}
  ],
  "detail_configs": [
    {"screen_key": "sessions-by-subject-list", "modal_screen_key": "sessions-by-subject-form", "parent_id_param": "subjectId", "child_id_field": "id", "title": "Sesiones"}
  ],
  "actions_removed": ["detail"],
  "api_prefix": "academic"
}`,
	}
}

// myMembershipsList (plan 006, N1.C ETAPA 2): pantalla "Mis materias" del
// alumno. Lista readonly de las materias en las que está inscrito. El contrato
// KMP consume el lector A (GET /api/v1/me/subject-offerings) y produce las
// columnas name/code a partir de subject_name/subject_code, por eso las
// columnas declaran name y code. requiredPermission =
// academic.my_memberships.read:own (permiso ÚNICO del feature self del alumno:
// slot.permission de la pantalla, route gate del dato y visibilidad del item de
// menú). Sin acciones de mutación: actions_removed = [create, edit, delete].
// Reintroducida en N1.7 F1 sobre el modelo de sesiones.
func myMembershipsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MY_MEMBERSHIPS_LIST_ID,
		screenKey:          "my-memberships-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Mis Materias",
		description:        "Materias en las que el alumno está inscrito",
		scope:              "unit",
		requiredPermission: "academic.my_memberships.read:own",
		slotData: `{
  "title": "Mis Materias",
  "search_placeholder": "Buscar materia...",
  "columns": [
    {"key": "name", "label": "Materia"},
    {"key": "code", "label": "Código"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// myGradesList (N3 F4, consulta de notas): pantalla "Mis notas" del alumno.
// Lista readonly de sus notas por sesión de materia. Espejo de
// myMembershipsList: el contrato KMP consume el lector self
// GET /api/v1/me/grades (el seed solo declara columnas/título/permiso).
// requiredPermission = academic.my_grades.read:own (permiso ÚNICO del feature
// self del alumno: slot.permission de la pantalla, route gate del dato y
// visibilidad del item de menú). Sin acciones de mutación: actions_removed =
// [create, edit, delete]. Columnas: subject_name (materia), period_name
// (período) y grade (nota propia).
func myGradesList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MY_GRADES_LIST_ID,
		screenKey:          "my-grades-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Mis Notas",
		description:        "Notas del alumno por sesión de materia",
		scope:              "unit",
		requiredPermission: "academic.my_grades.read:own",
		slotData: `{
  "title": "Mis Notas",
  "search_placeholder": "Buscar materia...",
  "columns": [
    {"key": "subject_name", "label": "Materia"},
    {"key": "period_name", "label": "Período"},
    {"key": "grade", "label": "Nota"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "academic"
}`,
	}
}

func periodsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_PERIODS_LIST_ID,
		screenKey:          "periods-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Periodos Académicos",
		description:        "Listado de periodos académicos",
		scope:              "school",
		requiredPermission: "academic.periods.read",
		slotData: `{
  "title": "Periodos",
  "search_placeholder": "Buscar periodo...",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "start_date", "label": "Inicio"},
    {"key": "end_date", "label": "Fin"},
    {"key": "is_active", "label": "Activo"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

func periodsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_PERIODS_FORM_ID,
		screenKey:          "periods-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Periodo",
		description:        "Crear o editar un periodo académico",
		scope:              "school",
		requiredPermission: "academic.periods.read",
		slotData: `{
  "title": "Periodo Académico",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "type", "label": "Tipo", "type": "select", "required": true, "options": [
      {"value": "semester", "label": "Semestre"},
      {"value": "trimester", "label": "Trimestre"},
      {"value": "bimester", "label": "Bimestre"},
      {"value": "quarter", "label": "Cuatrimestre"}
    ]},
    {"key": "academic_year", "label": "Año académico", "type": "number", "required": true, "min": 1900, "max": 2100},
    {"key": "start_date", "label": "Inicio", "type": "date", "required": true},
    {"key": "end_date", "label": "Fin", "type": "date", "required": true},
    {"key": "is_active", "label": "Activo", "type": "toggle"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// ACADEMIC: invitations (códigos de invitación a colegio/unidad)
// ===============================================================
//
// invitations-list (N0.4-A, plan 005): listado de códigos de
// invitación que el admin reparte. Patrón delta sobre list-basic-v1:
//   - actions_removed [edit, delete]: las invitaciones NO se editan ni
//     borran como CRUD estándar; el ciclo de vida es crear → revocar.
//   - actions_added [revoke] (scope row, permission
//     academic.invitations.revoke, event_id "revoke"): desactiva el
//     código. El FE resuelve el id del item y hace POST al endpoint
//     /revoke (handler custom, no SubmitTo estándar de CRUD).
//   - create (header) se HEREDA del template: $resource$.create →
//     academic.invitations.create (derivado de required_permission).
func invitationsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_INVITATIONS_LIST_ID,
		screenKey:          "invitations-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Invitaciones",
		description:        "Listado de códigos de invitación a colegio/unidad",
		scope:              "school",
		requiredPermission: "academic.invitations.read",
		slotData: `{
  "title": "Invitaciones",
  "search_placeholder": "Buscar invitación...",
  "columns": [
    {"key": "code", "label": "Código"},
    {"key": "invitation_type_label", "label": "Tipo"},
    {"key": "label", "label": "Etiqueta"},
    {"key": "uses_count", "label": "Usos"},
    {"key": "max_uses", "label": "Máx."},
    {"key": "is_active", "label": "Activa"},
    {"key": "expires_at", "label": "Expira"}
  ],
  "actions_removed": ["edit", "delete"],
  "actions_added": [
    {"id": "copy-code", "scope": "row", "label": "Copiar código", "icon": "copy", "permission": "academic.invitations.read", "condition": "always", "event_id": "copy-code", "order": 10},
    {"id": "revoke", "scope": "row", "label": "Revocar", "icon": "ban", "permission": "academic.invitations.revoke", "condition": "always", "event_id": "revoke", "style": "destructive", "order": 20}
  ],
  "api_prefix": "academic"
}`,
	}
}

// invitations-form (N0.4-A; MP-08 F4): creación de un código de invitación.
// Solo create (no edit): patrón delta retira "save" (edit-only) y
// "delete"; conserva "save_new" → $resource$.create →
// academic.invitations.create. El campo `code` NO se incluye: lo
// genera el backend. academic_unit_id se llena vía remote_select de
// unidades del colegio (remoteSelectPrefix=academic en el contrato FE).
//
// MP-08 F4: el campo `role` (select estático con el enum legacy) murió.
// Ahora es `invitation_type` (la KEY del tipo configurado para la escuela),
// que el backend exige en CreateInvitationRequest (json:"invitation_type"
// binding:"required"). Se llena por remote_select contra el endpoint nuevo
// GET /api/v1/schools/invitation-types (JWT-only, los tipos configurados
// para la escuela activa), que responde
// {"invitation_types":[{"key","label","requires_unit"}]}. El RemoteDataLoader
// del KMP localiza el array por el fallback "primer array de objetos de nivel
// superior" (no hay envelope items/data), y el select lee value_field=key y
// display_field=label de cada objeto. value=key (lo que el body envía),
// label=label (texto legible).
func invitationsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_INVITATIONS_FORM_ID,
		screenKey:          "invitations-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Nueva Invitación",
		description:        "Generar un código de invitación a colegio/unidad",
		scope:              "school",
		requiredPermission: "academic.invitations.read",
		slotData: `{
  "title": "Nueva Invitación",
  "fields": [
    {"key": "academic_unit_id", "label": "Unidad", "type": "remote_select", "required": true, "remote_endpoint": "academic:/api/v1/units", "display_field": "display_name", "value_field": "id"},
    {"key": "invitation_type", "label": "Tipo de invitación", "type": "remote_select", "required": true, "remote_endpoint": "academic:/api/v1/schools/invitation-types", "display_field": "label", "value_field": "key"},
    {"key": "label", "label": "Etiqueta", "type": "text"},
    {"key": "expires_at", "label": "Expira", "type": "datetime"},
    {"key": "max_uses", "label": "Usos máximos", "type": "number", "min": 1}
  ],
  "actions_removed": ["save", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// invitations-detail: detalle de SOLO LECTURA de un código de invitación.
// El tap de fila en invitations-list navega aquí (two-pane en desktop / apilada
// en móvil). La invitación es inmutable (ciclo crear→revocar), así que NO es un
// form editable: reusa el template genérico form-basic-v1 en modo lectura
// (accessMode="view" lo fuerza el panel de detalle) pintando los campos como
// valores. El fetch lo hace InvitationsDetailContract (KMP) vía
// GET academic:/api/v1/schools/invitations/{id}; las KEYS de los campos abajo
// coinciden con los json tags de dto.InvitationResponse para que el prellenado
// del form mapee 1:1. NO lleva acción "copiar" aquí: copiar el código es una
// row-action de la lista (event_id "copy-code"), que funciona sin abrir el
// detalle y no pelea con el panel read-only.
func invitationsDetail() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_INVITATIONS_DETAIL_ID,
		screenKey:          "invitations-detail",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Detalle de Invitación",
		description:        "Detalle de solo lectura de un código de invitación",
		scope:              "school",
		requiredPermission: "academic.invitations.read",
		slotData: `{
  "title": "Detalle de Invitación",
  "fields": [
    {"key": "code", "label": "Código", "type": "text"},
    {"key": "invitation_type_label", "label": "Tipo", "type": "text"},
    {"key": "academic_unit_name", "label": "Unidad", "type": "text"},
    {"key": "label", "label": "Etiqueta", "type": "text"},
    {"key": "uses_count", "label": "Usos", "type": "text"},
    {"key": "max_uses", "label": "Usos máximos", "type": "text"},
    {"key": "is_active", "label": "Activa", "type": "text"},
    {"key": "expires_at", "label": "Vence", "type": "text"}
  ],
  "actions_removed": ["save", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// joinRequestsInbox (N0.4-B, plan 005): bandeja de solicitudes de
// ingreso pendientes con doble visto bueno (colegio + unidad). La
// pantalla es NATIVA en el FE (Compose, NO SDUI) porque el gating del
// botón Aprobar depende del alcance del aprobador + estado de cada gate
// + permiso de aprobación POR ROL — lógica que el motor SDUI no expresa.
//
// Esta screen_instance existe SOLO para satisfacer la FK
// resource_screens.screen_key → screen_instances.screen_key y para que el
// menú pueda resolver el screen_key `join-requests-inbox`. El slot_data
// NUNCA se renderiza por el SDUI engine: MainScreen intercepta el
// screen_key y pinta JoinRequestsInboxScreen directamente. Se conserva un
// slot_data mínimo y válido (list-basic-v1) por higiene.
func joinRequestsInbox() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_JOIN_REQUESTS_INBOX_ID,
		screenKey:          "join-requests-inbox",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Solicitudes de Ingreso",
		description:        "Bandeja de solicitudes de ingreso pendientes (pantalla nativa)",
		scope:              "school",
		requiredPermission: "academic.join_requests.read",
		slotData: `{
  "title": "Solicitudes de Ingreso",
  "columns": [
    {"key": "user", "label": "Solicitante"},
    {"key": "role", "label": "Rol"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// subjectOfferingsBatchEnroll (N1.7 F1, plan 010 / ADR 0009): pantalla de
// "inscripción por lote" de alumnos en una sesión de materia (subject_offering).
// La pantalla es NATIVA en el FE (Compose, NO SDUI): seleccionar la sesión +
// marcar/desmarcar alumnos + confirmar es un flujo de selección masiva que el
// motor SDUI list/form no expresa. MainScreen intercepta el screen_key
// `batch-enroll` y pinta la pantalla nativa.
//
// Esta screen_instance existe para satisfacer la FK
// resource_screens.screen_key → screen_instances.screen_key y para que el menú
// resuelva el screen_key. El slot_data NUNCA se renderiza por el SDUI engine;
// se conserva mínimo y válido (list-basic-v1) por higiene.
//
// Permiso de visibilidad (requiredPermission, slot.permission de la pantalla):
// academic.subject_offerings.read — ver la pantalla. El botón "Inscribir" se
// declara como action en slot_data con permission
// academic.subject_offerings.enroll (ADR 0003: única fuente del permiso del
// botón). El FE nativo lee esa action del contrato y gatea el botón con ese
// permiso. La action sigue el esquema real de actions_added[] (mismas keys que
// p.ej. attendanceList/invitationsList): id, scope, label, icon, permission,
// condition, event_id, style, order. El permiso se lee de la key `permission`.
func subjectOfferingsBatchEnroll() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SUBJECT_OFFERINGS_BATCH_ENROLL_ID,
		screenKey:          "batch-enroll",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Inscripción por Lote",
		description:        "Inscribir alumnos en una sesión de materia (pantalla nativa)",
		scope:              "school",
		requiredPermission: "academic.subject_offerings.read",
		slotData: `{
  "title": "Inscripción por Lote",
  "columns": [
    {"key": "user_name", "label": "Alumno"},
    {"key": "enrolled", "label": "Inscrito"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "actions_added": [
    {"id": "enroll", "scope": "header", "label": "Inscribir", "icon": "person_add", "permission": "academic.subject_offerings.enroll", "condition": "always", "event_id": "enroll", "style": "filled", "order": 10}
  ],
  "api_prefix": "academic"
}`,
	}
}

// enrollOne (N1.7 F2, plan 010 / ADR 0009): pantalla NATIVA de "inscripción
// individual" de UN alumno en una sesión de materia (subject_offering). Igual
// que batch-enroll, la pantalla es NATIVA en el FE (Compose, NO SDUI):
// MainScreen intercepta el screen_key `enroll-one` y pinta la pantalla nativa.
//
// Esta screen_instance existe para satisfacer la FK
// resource_screens.screen_key → screen_instances.screen_key y para que el
// handler resuelva el screen_key. El slot_data NUNCA se renderiza por el SDUI
// engine; se conserva mínimo y válido por higiene, replicando la forma de
// batch-enroll (action `enroll` IDÉNTICA: misma permission/event_id/icon/style).
//
// requiredPermission (slot.permission de la pantalla) = academic.subject_offerings.read
// para verla; el botón "Inscribir" se declara como action con permission
// academic.subject_offerings.enroll (ADR 0003: única fuente del permiso del botón).
func enrollOne() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ENROLL_ONE_ID,
		screenKey:          "enroll-one",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Inscripción Individual",
		description:        "Inscribir un alumno en una sesión de materia (pantalla nativa)",
		scope:              "school",
		requiredPermission: "academic.subject_offerings.read",
		slotData: `{
  "title": "Inscripción Individual",
  "columns": [
    {"key": "user_name", "label": "Alumno"},
    {"key": "enrolled", "label": "Inscrito"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "actions_added": [
    {"id": "enroll", "scope": "header", "label": "Inscribir", "icon": "person_add", "permission": "academic.subject_offerings.enroll", "condition": "always", "event_id": "enroll", "style": "filled", "order": 10}
  ],
  "api_prefix": "academic"
}`,
	}
}

// sessionsBySubjectList (N1.7 F2, plan 010 / ADR 0009; reubicada en F2.2): lista
// hija de "sesiones de la materia". Pantalla SDUI list estándar (no nativa). Se
// alcanza embebida como pestaña "Sesiones" del master-detail subjects-form (vía
// detail_configs[]); el contenedor le inyecta subjectId = id de la materia
// editada y consume el endpoint
// GET /api/v1/subject-offerings?subject_id={subjectId} (lo resuelve el handler
// KMP; el seed solo declara columnas/título/permiso). El id de cada fila es el
// offering_id (la sesión concreta).
//
// Columnas (reordenadas en N3.5 F1): section_label primero — es el headline que
// distingue la sección A de la B —, luego period_name y teacher_name. Se quitó
// subject_name: es redundante porque todas las filas son la MISMA materia (ya
// estamos dentro de su detalle).
//
// Row-actions de asistencia/notas (N3.5 F1, plan 014 / ADR 0018; +consulta de
// notas en N3 F4): la card de cada sesión lleva las 5 acciones del docente —
// "Pasar lista" (take-attendance), "Poner notas" (put-grades), "Historial"
// (view-attendance), "Resumen" (view-attendance-summary) y "Resumen de notas"
// (view-grades-summary) —, todas scope row (se materializan como RowAction en
// el KMP). Vinieron de subjects-form (antes scope resource-toolbar, mezclaban el
// roster de un docente con dos secciones de la misma materia); ahora operan sobre
// la sesión concreta. condition=always: la fila SIEMPRE es una sesión existente
// (no hay modo create/edit como en la toolbar del form). El id de la fila
// (offering_id) viajará como offeringId al evento (mapeo en el contrato KMP, F2).
// Cada permiso es slot.permission (ADR 0003): take-attendance →
// academic.attendance.create, put-grades → academic.grades.create, view-attendance/
// view-attendance-summary → academic.attendance.read, view-grades-summary →
// academic.grades.read (navega a grades-subject-summary; ya sembrados, cubiertos
// por el wildcard academic.* de teacher). Solo lectura del CRUD de sesiones:
// actions_removed retira create/edit/delete heredados del template.
// requiredPermission (slot.permission) = academic.subject_offerings.read.
func sessionsBySubjectList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:          L4_SCREEN_INST_SESSIONS_BY_SUBJECT_ID,
		screenKey:   "sessions-by-subject-list",
		templateID:  L0_SCREEN_TPL_LIST_ID_REF,
		name:        "Sesiones de la Materia",
		description: "Listado de sesiones (oferta) de una materia",
		// scope=unit (ADR 0016 punto 3): aunque la materia es catalogo de
		// ESCUELA, la GESTION de sus sesiones es por unidad activa — el backend
		// filtra las sesiones por la unidad del token. El scope declarado refleja
		// ese filtro real (antes decia "school", incoherente con el filtro).
		scope:              "unit",
		requiredPermission: "academic.subject_offerings.read",
		slotData: `{
  "title": "Sesiones",
  "columns": [
    {"key": "section_label", "label": "Sección"},
    {"key": "period_name", "label": "Período"},
    {"key": "teacher_name", "label": "Docente"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "actions_added": [
    {"id": "take-attendance", "scope": "row", "label": "Pasar lista", "icon": "checklist", "permission": "academic.attendance.create", "condition": "always", "event_id": "take-attendance", "style": "icon", "order": 20},
    {"id": "put-grades", "scope": "row", "label": "Poner notas", "icon": "star", "permission": "academic.grades.create", "condition": "always", "event_id": "put-grades", "style": "icon", "order": 21},
    {"id": "view-attendance", "scope": "row", "label": "Historial", "icon": "history", "permission": "academic.attendance.read", "condition": "always", "event_id": "view-attendance", "style": "icon", "order": 22},
    {"id": "view-attendance-summary", "scope": "row", "label": "Resumen", "icon": "bar_chart", "permission": "academic.attendance.read", "condition": "always", "event_id": "view-attendance-summary", "style": "icon", "order": 23},
    {"id": "view-grades-summary", "scope": "row", "label": "Resumen de notas", "icon": "pie_chart", "permission": "academic.grades.read", "condition": "always", "event_id": "view-grades-summary", "style": "icon", "order": 24}
  ],
  "api_prefix": "academic"
}`,
	}
}

// sessionsBySubjectForm (N1.7 F2.3): formulario crear/editar de "sesión de
// materia" (subject offering). Se renderiza como MODAL del master-detail
// subjects-form: la pestaña "Sesiones" lo enlaza vía detail_configs[].
// modal_screen_key. El MasterDetailContainer abre el modal con subjectId (parent)
// en create y con id+subjectId en edición.
//
// Campos:
//   - period_id (remote_select, required): catálogo GET /api/v1/periods, que
//     responde {"periods":[{id,name,...}]}; el RemoteDataLoader resuelve el
//     array por fallback (no hay envelope items/data). display_field=name.
//     En edición es identidad inmutable → el contrato KMP lo marca readonly.
//   - section_label (text, opcional): etiqueta de sección (máx 10 en el backend).
//   - teacher_membership_id (remote_select, NO required): docentes de la unidad
//     activa vía GET /api/v1/memberships/by-role?role_key=teacher, que responde
//     {"memberships":[{id,full_name,display_name,...}]}; display_field=full_name
//     (nombre real de la persona; display_name lleva el rol "Profesor").
//     Asigna o cambia el docente; el backend acepta omitirlo (deja intacto).
//   - is_active (toggle, default true): el form renderer mapea toggle→SWITCH y
//     serializa el valor como booleano JSON limpio (no string), alineado al
//     IsActive *bool del UpdateSubjectOfferingRequest.
//
// subject_id NO es un campo del form: el contrato KMP lo inyecta al body en
// create desde context.params["subjectId"]. requiredPermission =
// academic.subject_offerings.update (gate del slot de mutación de la sesión).
func sessionsBySubjectForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:          L4_SCREEN_INST_SESSIONS_BY_SUBJECT_FORM_ID,
		screenKey:   "sessions-by-subject-form",
		templateID:  L0_SCREEN_TPL_FORM_ID_REF,
		name:        "Formulario de Sesión",
		description: "Crear o editar una sesión de materia (período, sección y docente)",
		// scope=unit (ADR 0016 punto 3): el form gestiona UNA sesión, que el
		// backend filtra por la unidad activa del token, y su selector de docente
		// (memberships/by-role) requiere unidad activa. Coherente con
		// sessions-by-subject-list, ya en scope=unit (antes decía "school",
		// incoherente con el contexto que el form realmente exige).
		scope:              "unit",
		requiredPermission: "academic.subject_offerings.update",
		slotData: `{
  "title": "Sesión",
  "fields": [
    {"key": "period_id", "label": "Período", "type": "remote_select", "required": true, "remote_endpoint": "academic:/api/v1/periods", "display_field": "name", "value_field": "id"},
    {"key": "section_label", "label": "Sección", "type": "text", "max_length": 10},
    {"key": "teacher_membership_id", "label": "Docente", "type": "remote_select", "remote_endpoint": "academic:/api/v1/memberships/by-role?role_key=teacher", "display_field": "full_name", "value_field": "id"},
    {"key": "is_active", "label": "Activa", "type": "toggle", "default": "true"}
  ],
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// ACADEMIC: grades / attendance
// ===============================================================
//
// Poda F2 (plan 004-permisologia-mvp): los constructores de guardian
// (guardian-relations-list/form, guardian_relations-form alias,
// guardian-requests-list), calendar (calendar-list/form) y schedules
// (schedules-list/form) se eliminaron del MVP junto con sus constantes
// y filas en resource_screens.go. Los recursos academic.guardian_relations,
// academic.calendar y academic.schedules quedan huérfanos (prune-later,
// ver docs/plans/004-permisologia-mvp/diferido.md).

func gradesList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_GRADES_LIST_ID,
		screenKey:          "grades-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Calificaciones",
		description:        "Listado de calificaciones",
		scope:              "unit",
		requiredPermission: "academic.grades.read",
		slotData: `{
  "title": "Calificaciones",
  "search_placeholder": "Buscar calificación...",
  "filter_all_label": "Todos",
  "filter_ready_label": "Pendientes",
  "filter_processing_label": "Finalizadas",
  "columns": [
    {"key": "student_name", "label": "Estudiante"},
    {"key": "subject", "label": "Materia"},
    {"key": "score", "label": "Nota"},
    {"key": "period", "label": "Periodo"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// grades-form ELIMINADA (2026-06-09): formulario SDUI legacy huérfano,
// reemplazado por pantallas nativas (my-grade-detail para el alumno,
// grades-batch para el docente). No tenía entry-point en el FE y sus
// campos student_id/subject_id eran remote_select MUERTOS (sin endpoint).
// Su constructor, su llamada en screen_instances.go, su mapping en
// resource_screens.go y la constante L4_SCREEN_INST_GRADES_FORM_ID
// (UUID …0071) se eliminaron. UUID …0071 queda libre para reuso futuro.

func attendanceList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ATTENDANCE_LIST_ID,
		screenKey:          "attendance-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Asistencia",
		description:        "Registro de asistencia",
		scope:              "unit",
		requiredPermission: "academic.attendance.read",
		slotData: `{
  "title": "Asistencia",
  "search_placeholder": "Buscar registro...",
  "columns": [
    {"key": "student_name", "label": "Estudiante"},
    {"key": "date", "label": "Fecha"},
    {"key": "status", "label": "Estado"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "actions_added": [
    {"id": "batch", "scope": "header", "label": "Registrar día", "icon": "plus", "permission": "academic.attendance.create", "condition": "always", "event_id": "batch", "style": "icon", "order": 10}
  ],
  "api_prefix": "academic"
}`,
	}
}

// attendanceBatch (N2.S2, plan 008 D5): pantalla "pasar lista" por sesión.
// Es OVERRIDE NATIVO en el FE (Compose, NO SDUI): MainScreen intercepta el
// screen_key `attendance-batch` y pinta AttendanceBatchViewModel/Screen
// (selección masiva de presentes/ausentes + upsert), que el motor SDUI form
// no expresa. El slot_data NO lo renderiza el SDUI genérico; se conserva
// mínimo y válido (form-basic-v1) por higiene de contrato.
//
// La action `submit-batch` declara el permiso del botón "Pasar lista" del
// override nativo (ADR 0003: única fuente del permiso del botón). NO es una
// action que el SDUI genérico pinte: la pantalla nativa la lee del contrato y
// gatea su botón con `permission`. El `event_id` debe ser `submit-batch` (la
// pantalla nativa busca la action cuyo event_id/id ∈ {submit-batch, save,
// take-attendance}); el permiso `academic.attendance.create` ya está sembrado y
// lo cubre el wildcard `academic.attendance.*` del rol teacher. Espeja la
// action `enroll` de subjectOfferingsBatchEnroll (mismas keys: id, scope,
// label, icon, permission, condition, event_id, style, order).
func attendanceBatch() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ATTENDANCE_BATCH_ID,
		screenKey:          "attendance-batch",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Registrar Asistencia",
		description:        "Formulario para registrar asistencia por día",
		scope:              "unit",
		requiredPermission: "academic.attendance.read",
		slotData: `{
  "title": "Registrar Asistencia",
  "fields": [
    {"key": "date", "label": "Fecha", "type": "date", "required": true},
    {"key": "entries", "label": "Asistencias", "type": "table"}
  ],
  "actions_removed": ["save", "delete"],
  "actions_added": [
    {"id": "submit-batch", "scope": "header", "label": "Pasar lista", "icon": "checklist", "permission": "academic.attendance.create", "condition": "always", "event_id": "submit-batch", "style": "filled", "order": 10}
  ],
  "api_prefix": "academic"
}`,
	}
}

// gradesBatch (N3 F3): pantalla "poner notas" por sesión, espejo de
// attendanceBatch. Es OVERRIDE NATIVO en el FE (Compose, NO SDUI): MainScreen
// intercepta el screen_key `grades-batch` y pinta el ViewModel/Screen de
// registro masivo de calificaciones, que el motor SDUI form no expresa. El
// slot_data NO lo renderiza el SDUI genérico; se conserva mínimo y válido
// (form-basic-v1) por higiene de contrato, análogo a attendance-batch.
//
// La action `submit-batch` declara el permiso del botón "Guardar notas" del
// override nativo (ADR 0003: única fuente del permiso del botón). NO es una
// action que el SDUI genérico pinte: la pantalla nativa la lee del contrato y
// gatea su botón con `permission`. El permiso `academic.grades.create` ya está
// sembrado y lo cubre el wildcard `academic.grades.*` del rol teacher.
//
// El entry-point vive en subjects-form (action `put-grades`, event_id
// `put-grades` → NavigateTo("grades-batch", {subjectId})). requiredPermission
// = academic.grades.read abre la pantalla (espejo de attendance-batch). El
// selector de período se declara como remote_select a academic:/api/v1/periods
// (mismo endpoint que sessions-by-subject-form).
func gradesBatch() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_GRADES_BATCH_ID,
		screenKey:          "grades-batch",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Registrar Calificaciones",
		description:        "Formulario para registrar calificaciones por período",
		scope:              "unit",
		requiredPermission: "academic.grades.read",
		slotData: `{
  "title": "Registrar Calificaciones",
  "fields": [
    {"key": "period_id", "label": "Período", "type": "remote_select", "required": true, "remote_endpoint": "academic:/api/v1/periods", "display_field": "name", "value_field": "id"},
    {"key": "entries", "label": "Calificaciones", "type": "table"}
  ],
  "actions_removed": ["save", "delete"],
  "actions_added": [
    {"id": "submit-batch", "scope": "header", "label": "Guardar notas", "icon": "star", "permission": "academic.grades.create", "condition": "always", "event_id": "submit-batch", "style": "filled", "order": 10}
  ],
  "api_prefix": "academic"
}`,
	}
}

func attendanceSummary() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ATTENDANCE_SUMMARY_ID,
		screenKey:          "attendance-summary",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Resumen de Asistencia",
		description:        "Resumen estadístico de asistencia",
		scope:              "unit",
		requiredPermission: "academic.attendance.read",
		slotData: `{
  "title": "Resumen",
  "columns": [
    {"key": "student_name", "label": "Estudiante"},
    {"key": "total_classes", "label": "Clases"},
    {"key": "present", "label": "Presentes"},
    {"key": "absent", "label": "Ausentes"},
    {"key": "rate", "label": "% Asistencia"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "academic"
}`,
	}
}

// gradesSubjectSummary (N3 F4, consulta de notas): resumen de notas por sesión
// de materia (vista del docente). Lista readonly espejo de attendanceSummary:
// template list, actions_removed=[create,edit,delete] + readonly, api_prefix
// academic, scope unit. El contrato KMP consume el endpoint de resumen ya
// existente GET /api/v1/grades/subject-summary (el seed solo declara
// columnas/título/permiso); el destino del evento view-grades-summary de la
// card de sesión (sessions-by-subject-list) es esta pantalla. Columnas:
// student_name + nota (grade_value numérico y grade_letter literal) + graded
// (indicador "sin nota"/ungraded cuando el alumno aún no tiene calificación).
// requiredPermission = academic.grades.read (slot.permission, ADR 0003; ya
// sembrado, cubierto por el wildcard academic.grades.* de teacher).
func gradesSubjectSummary() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_GRADES_SUBJECT_SUMMARY_ID,
		screenKey:          "grades-subject-summary",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Resumen de Notas",
		description:        "Resumen de notas por sesión de materia",
		scope:              "unit",
		requiredPermission: "academic.grades.read",
		slotData: `{
  "title": "Resumen de notas",
  "columns": [
    {"key": "student_name", "label": "Estudiante"},
    {"key": "grade_value", "label": "Nota"},
    {"key": "grade_letter", "label": "Letra"},
    {"key": "graded", "label": "Calificado"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// CONTENT: assessments
// ===============================================================

func assessmentsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_LIST_ID,
		screenKey:          "assessments-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Evaluaciones",
		description:        "Listado de evaluaciones",
		scope:              "unit",
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Evaluaciones",
  "search_placeholder": "Buscar evaluación...",
  "columns": [
    {"key": "title", "label": "Título"},
    {"key": "subject_name", "label": "Materia"},
    {"key": "status", "label": "Estado"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "learning"
}`,
	}
}

// assessmentsForm usa master-detail-v1: hereda los 3 defaults
// (save_new/save/delete con scope=form-submit) y declara via
// actions_added[] las acciones de recurso específicas de evaluación
// (detail=Preguntas, publish, archive), todas con scope=resource-toolbar.
// El default "detail" del template se overridea por id — el composer
// reemplaza el default genérico ("Detalle") por la versión específica
// ("Preguntas", event_id=view-questions, icon=help_outline).
//
// detail_configs[] describe los paneles detalle (aquí uno solo: "Preguntas"
// con modal de crear/editar). El frontend KMP es quien lo interpreta; el
// backend solo lo persiste como blob.
//
// Contrato N4 (plan 015): POST/GET /api/v1/assessments, GET/PUT/DELETE
// /assessments/:assessment_id (read/update/delete), POST .../publish y
// .../archive. El cuerpo de crear NO lleva school_id ni autor (del JWT); el
// `subject_id` SÍ va en el cuerpo (FK al catálogo de materias). `modality`
// se eliminó (no existe en el esquema nuevo). `subject_id` se llena con un
// remote_select al catálogo de materias de academic (GET /api/v1/subjects,
// display=name, value=id).
func assessmentsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_FORM_ID,
		screenKey:          "assessments-form",
		templateID:         L0_SCREEN_TPL_MASTER_DETAIL_ID_REF,
		name:               "Formulario de Evaluación",
		description:        "Crear o editar una evaluación",
		scope:              "unit",
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Evaluación",
  "page_title": "Evaluación",
  "edit_title": "Editar evaluación",
  "view_when": {"field": "status", "in": ["published", "archived"]},
  "fields": [
    {"key": "title", "label": "Título", "type": "text", "required": true},
    {"key": "subject_id", "label": "Materia", "type": "entity-picker", "required": true, "remote_endpoint": "academic:/api/v1/subjects", "display_field": "name", "value_field": "id", "search_param": "search", "page_size": 20, "picker_title": "Buscar materia"},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "pass_threshold", "label": "Umbral de aprobación (%)", "type": "number", "min": 0, "max": 100},
    {"key": "max_attempts", "label": "Intentos máximos", "type": "number", "min": 1},
    {"key": "time_limit_minutes", "label": "Tiempo límite (min)", "type": "number", "min": 1},
    {"key": "is_timed", "label": "Cronometrada", "type": "boolean"},
    {"key": "shuffle_questions", "label": "Mezclar preguntas", "type": "boolean"},
    {"key": "show_correct_answers", "label": "Mostrar respuestas correctas", "type": "boolean"},
    {"key": "available_from", "label": "Disponible desde", "type": "datetime"},
    {"key": "available_until", "label": "Disponible hasta", "type": "datetime"}
  ],
  "detail_configs": [
    {"screen_key": "assessment-questions-list", "modal_screen_key": "assessment-question-form", "parent_id_param": "assessmentId", "child_id_field": "id"}
  ],
  "actions_removed": ["delete"],
  "actions_added": [
    {"id": "detail",  "scope": "resource-toolbar", "icon": "help_outline", "label": "Preguntas", "permission": "content.assessments.read",   "condition": "edit-only", "event_id": "view-questions", "style": "icon", "order": 15},
    {"id": "assign",  "scope": "resource-toolbar", "icon": "assignment",   "label": "Asignar",   "permission": "content.assessments.assign", "condition": "edit-only", "event_id": "assign",         "style": "icon", "order": 20, "visible_when": {"field": "status", "equals": "published"}},
    {"id": "publish", "scope": "resource-toolbar", "icon": "check_circle", "label": "Publicar",  "permission": "content.assessments.publish", "condition": "edit-only", "event_id": "publish",        "style": "icon", "order": 30, "visible_when": {"field": "status", "equals": "draft"}},
    {"id": "archive", "scope": "resource-toolbar", "icon": "archive",      "label": "Archivar",  "permission": "content.assessments.update", "condition": "edit-only", "event_id": "archive",        "style": "icon", "order": 40, "visible_when": {"field": "status", "equals": "published"}},
    {"id": "delete",  "scope": "form-submit",      "icon": "trash",        "label": "Eliminar",  "permission": "content.assessments.delete", "condition": "edit-only", "event_id": "delete",         "style": "destructive", "order": 50, "visible_when": {"field": "status", "equals": "draft"}}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessmentsMgmtList — F3.1 (plan 004): migrada al patrón delta.
// Hereda los 3 default_actions de list-basic-v1 (create/edit/delete con
// $resource$ → "content.assessments"). El legacy ya declaraba
// scope=header/row y los mismos permisos, así que el conjunto invariante
// {event_id, permission, scope} no cambia (verificado por el harness).
func assessmentsMgmtList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_MGMT_LIST_ID,
		screenKey:          "assessments-management-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Gestión de Evaluaciones",
		description:        "Vista de gestión para docentes",
		scope:              "unit",
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Gestión de Evaluaciones",
  "page_title": "Evaluaciones",
  "search_placeholder": "Buscar...",
  "columns": [
    {"key": "title", "label": "Título"},
    {"key": "subject_name", "label": "Materia"},
    {"key": "questions_count", "label": "Preguntas"},
    {"key": "status", "label": "Estado"}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessmentTake: F3 (el alumno presenta la evaluación). Pendiente de
// re-apuntar a los endpoints de intento del backend nuevo (start/save/submit
// por student_membership_id) en F3.1/F3.2 — aquí queda MÍNIMO, no se inventa
// el contrato. Permiso del alumno: content.assessments_student.read.
func assessmentTake() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_TAKE_ID,
		screenKey:          "assessment-take",
		templateID:         L0_SCREEN_TPL_DETAIL_ID_REF,
		name:               "Tomar Evaluación",
		description:        "Pantalla para rendir una evaluación",
		scope:              "unit",
		requiredPermission: "content.assessments_student.read",
		slotData: `{
  "title": "Tomar Evaluación",
  "submit_label": "Enviar respuestas",
  "api_prefix": "learning"
}`,
	}
}

// assessmentResult: F3 (resultado/revisión del intento, vista alumno).
// Pendiente de re-apuntar al backend nuevo en F3.1 — queda MÍNIMO.
func assessmentResult() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_RESULT_ID,
		screenKey:          "assessment-result",
		templateID:         L0_SCREEN_TPL_DETAIL_ID_REF,
		name:               "Resultado de Evaluación",
		description:        "Resultado y revisión de una evaluación rendida",
		scope:              "unit",
		requiredPermission: "content.assessments_student.read",
		slotData: `{
  "title": "Resultado",
  "readonly": true,
  "api_prefix": "learning"
}`,
	}
}

// assessmentQuestionsList — detalle (lista) de preguntas de una evaluación.
// `actions_removed: ["edit"]` poda la row-action `edit` que list-basic-v1 trae
// por default: en este detalle la edición la abre el botón nativo "Editar" del
// bottom-sheet (MasterDetailContainer, flujo N3.5), por lo que la acción SDUI
// quedaba huérfana (sin handler). Mismo criterio que las listas de sesiones.
func assessmentQuestionsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:          L4_SCREEN_INST_ASSESS_QUESTIONS_LIST_ID,
		screenKey:   "assessment-questions-list",
		templateID:  L0_SCREEN_TPL_LIST_ID_REF,
		name:        "Preguntas de Evaluación",
		description: "Listado de preguntas de una evaluación",
		scope:       "unit",
		// TC-A del baseline: resource=assessments (no assessment-questions).
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Preguntas",
  "page_title": "Preguntas",
  "columns": [
    {"key": "question_text", "label": "Pregunta"},
    {"key": "question_type", "label": "Tipo"},
    {"key": "points", "label": "Puntaje"}
  ],
  "actions_removed": ["edit"],
  "actions_added": [
    {"id": "create", "scope": "header", "label": "Nuevo",    "icon": "plus",  "permission": "content.assessments.update", "condition": "always", "event_id": "create", "style": "icon",        "order": 10},
    {"id": "delete", "scope": "row",    "label": "Eliminar", "icon": "trash", "permission": "content.assessments.update", "condition": "always", "event_id": "delete", "style": "destructive", "order": 20}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessmentQuestionForm — F2.6 (plan 015) + Fase 1 visibilidad condicional +
// Fase 2 multiple_select: editor de preguntas REACTIVO por tipo.
// `question_type` es el campo CONTROLADOR (5 tipos del CHECK del esquema
// nuevo); el resto de campos de respuesta correcta declaran `visible_when`
// (formato {field, equals|in}) para mostrarse solo cuando aplican:
//   - multiple_choice → field `options` (type=option-list); lo consume el
//     componente KMP DynamicOptionListField (shape {option_id, option_text} por
//     opción). Su `correct_answer_field` apunta a `mc_correct_letter`: el
//     radio-button de la lista marca la opción correcta y escribe ese valor,
//     por eso NO hay un field separado para la opción correcta (sería un control
//     duplicado sobre el mismo dato).
//   - multiple_select → field `options_multi` (type=option-list,
//     selection_mode=multiple): opción múltiple con VARIAS respuestas
//     correctas. Key DISTINTA de `options` (single) para no colisionar el
//     estado del componente. `correct_answer_field` apunta a
//     `ms_correct_letters` (checkboxes de la lista marcan N correctas). En BD,
//     correct_answer guarda un ARRAY JSON de textos (["Texto A","Texto C"]).
//   - true_false → field `correct_answer_bool` (select Verdadero/Falso).
//   - short_answer → field `correct_answer_text` (text).
//   - open_ended → ningún campo de respuesta correcta (no se evalúa de forma
//     automática); queda implícito porque ninguno lo incluye en su in/equals.
//
// Los campos sin `visible_when` (question_text, question_type, points,
// explanation, difficulty) son siempre visibles. Endpoints (resueltos por el
// AssessmentQuestionFormContract del KMP): GET/POST
// /api/v1/assessments/:assessment_id/questions y PUT/DELETE .../:question_id,
// bajo api_prefix=learning.
func assessmentQuestionForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:          L4_SCREEN_INST_ASSESS_QUESTION_FORM_ID,
		screenKey:   "assessment-question-form",
		templateID:  L0_SCREEN_TPL_FORM_ID_REF,
		name:        "Formulario de Pregunta",
		description: "Crear o editar una pregunta",
		scope:       "unit",
		// slot.permission del editor: ver la pantalla con read; las mutaciones
		// (POST/PUT/DELETE de preguntas) son content.assessments.update y se
		// gatean en las acciones de la lista de preguntas + el submit del form.
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Pregunta",
  "page_title": "Pregunta",
  "edit_title": "Editar pregunta",
  "fields": [
    {"key": "question_text", "label": "Enunciado", "type": "textarea", "required": true},
    {"key": "question_type", "label": "Tipo", "type": "select", "required": true, "options": [
      {"value": "multiple_choice", "label": "Opción múltiple"},
      {"value": "multiple_select", "label": "Opción múltiple (varias)"},
      {"value": "true_false", "label": "Verdadero/Falso"},
      {"value": "short_answer", "label": "Respuesta corta"},
      {"value": "open_ended", "label": "Respuesta abierta"}
    ]},
    {"key": "options", "type": "option-list", "correct_answer_field": "mc_correct_letter", "visible_when": {"field": "question_type", "in": ["multiple_choice"]}},
    {"key": "options_multi", "type": "option-list", "selection_mode": "multiple", "correct_answer_field": "ms_correct_letters", "visible_when": {"field": "question_type", "in": ["multiple_select"]}},
    {"key": "correct_answer_bool", "label": "Respuesta correcta", "type": "select", "required": true, "visible_when": {"field": "question_type", "equals": "true_false"}, "options": [
      {"value": "true", "label": "Verdadero"},
      {"value": "false", "label": "Falso"}
    ]},
    {"key": "correct_answer_text", "label": "Respuesta correcta", "type": "text", "required": true, "visible_when": {"field": "question_type", "equals": "short_answer"}},
    {"key": "points", "label": "Puntaje", "type": "number", "required": true, "min": 0},
    {"key": "explanation", "label": "Explicación", "type": "textarea"},
    {"key": "difficulty", "label": "Dificultad", "type": "select", "options": [
      {"value": "easy", "label": "Fácil"},
      {"value": "medium", "label": "Media"},
      {"value": "hard", "label": "Difícil"}
    ]}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessment-assignment ELIMINADA: la asignación de una evaluación a una sesión
// de materia (subject_offering) se reemplaza por un modal NATIVO ("nativa
// prevalece, SDUI solo guía"). El SDUI form-basic-v1 quedaba muerto. Se conserva
// el recurso L4_RESOURCE_ASSESSMENTS y su permiso content.assessments.assign
// (lo sigue gateando la action "Asignar" del form de evaluación + la ruta
// POST /api/v1/assessments/:assessment_id/assignments). Concepto vivo, pantalla
// SDUI muerta: no se deprecó, se eliminó.

// assessmentModality ELIMINADA (plan 015, F2.6): la "modalidad"
// (quiz/examen/tarea) no existe en el esquema nuevo de N4 — el assessment solo
// tiene `source_type` (manual/ai_generated). El flujo de creación va directo al
// form de evaluación, sin selector previo. Concepto muerto, no deprecado.
// Deuda front: AssessmentModalityContract.kt + su test quedan inertes (sin
// screen_instance que los resuelva); limpiar en el re-apuntado de UI de F3.1.

// assessmentReviewDashboard: F3 (revisión docente). Dashboard de revisión de
// intentos por evaluación. Pendiente de re-apuntar a los endpoints de revisión
// del backend nuevo en F3.1 — aquí queda MÍNIMO (no se inventa el contrato).
func assessmentReviewDashboard() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_REVIEW_DASH_ID,
		screenKey:          "assessment-review-dashboard",
		templateID:         l4TplDashboardV1ID,
		name:               "Revisión de Evaluación",
		description:        "Dashboard de revisión de intentos por evaluación",
		scope:              "unit",
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Revisión",
  "greeting_text": "Evaluación",
  "kpi_students_label": "Intentos",
  "kpi_materials_label": "Promedio",
  "kpi_avg_score_label": "Aprobados",
  "kpi_completion_label": "Pendientes",
  "activity_title": "Intentos recientes",
  "api_prefix": "learning"
}`,
	}
}

// assignedAssessmentsList: lista para el estudiante de las evaluaciones
// asignadas. Contrato N4 (plan 015): GET /api/v1/me/assigned-assessments
// (resuelto por oferta→inscritos), permiso content.assessments_student.read.
// Solo lectura (sin create/edit/delete). Las columnas se alinean a los campos
// del esquema nuevo (subject_name, due_date).
func assignedAssessmentsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_ASSIGNED_LIST_ID,
		screenKey:          "assigned-assessments-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Mis Evaluaciones",
		description:        "Evaluaciones asignadas al estudiante",
		scope:              "unit",
		requiredPermission: "content.assessments_student.read",
		slotData: `{
  "title": "Mis Evaluaciones",
  "search_placeholder": "Buscar...",
  "columns": [
    {"key": "title", "label": "Título"},
    {"key": "subject_name", "label": "Materia"},
    {"key": "due_date", "label": "Vence"},
    {"key": "status", "label": "Estado"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "learning"
}`,
	}
}

// attemptReviewDetail: F3 (detalle de un intento, vista revisión docente).
// Pendiente de re-apuntar al backend nuevo en F3.1 — queda MÍNIMO.
func attemptReviewDetail() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ATTEMPT_REVIEW_DETAIL_ID,
		screenKey:          "attempt-review-detail",
		templateID:         L0_SCREEN_TPL_DETAIL_ID_REF,
		name:               "Detalle de Intento",
		description:        "Detalle de revisión de un intento de evaluación",
		scope:              "unit",
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Revisión de Intento",
  "readonly": true,
  "api_prefix": "learning"
}`,
	}
}

// ===============================================================
// REPORTS: stats / report-card
// ===============================================================
//
// Poda F2 (plan 004-permisologia-mvp): progress-detail, stats-detail y
// report-card se eliminaron del MVP junto con sus constantes y filas en
// resource_screens.go. El dashboard stats-dashboard SÍ se conserva
// (definido en screen_instances.go). progress-dashboard se eliminó el
// 2026-06-15 (apuntaba a un endpoint inexistente).

// ===============================================================
// DIRECTORIES & MISC
// ===============================================================

func unitDirectory() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_UNIT_DIRECTORY_ID,
		screenKey:          "unit-directory",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Directorio de Unidad",
		description:        "Directorio de miembros de la unidad",
		scope:              "unit",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Directorio",
  "search_placeholder": "Buscar miembro...",
  "columns": [
    {"key": "full_name", "label": "Nombre"},
    {"key": "role", "label": "Rol"},
    {"key": "subjects", "label": "Materias"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// PHANTOM-NUEVAS NO-ASSESSMENT
// ===============================================================

// school-concepts-list / school-concepts-form: variante scope=school
// del CRUD de concept_types. Permite que el admin de cada escuela
// haga overrides locales de la terminologia.
func schoolConceptsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCHOOL_CONCEPTS_LIST_ID,
		screenKey:          "school-concepts-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Conceptos de la Escuela",
		description:        "Terminología configurada a nivel escuela",
		scope:              "school",
		requiredPermission: "admin.concept_types.read",
		slotData: `{
  "title": "Conceptos (Escuela)",
  "columns": [
    {"key": "term_key", "label": "Clave"},
    {"key": "term_value", "label": "Valor"},
    {"key": "category", "label": "Categoría"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

func schoolConceptsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCHOOL_CONCEPTS_FORM_ID,
		screenKey:          "school-concepts-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Concepto Escolar",
		description:        "Crear o editar un término local",
		scope:              "school",
		requiredPermission: "admin.concept_types.read",
		slotData: `{
  "title": "Concepto (Escuela)",
  "fields": [
    {"key": "term_key", "label": "Clave", "type": "text", "required": true},
    {"key": "term_value", "label": "Valor", "type": "text", "required": true},
    {"key": "category", "label": "Categoría", "type": "select", "options": [
      {"value": "org", "label": "Organización"},
      {"value": "unit", "label": "Unidad"},
      {"value": "member", "label": "Miembro"},
      {"value": "content", "label": "Contenido"}
    ]}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "academic"
}`,
	}
}

// user-roles ELIMINADA (2026-06-09): pantalla SDUI legacy huérfana, sin
// reemplazo y sin navegación que la abriera (ningún entry-point en el FE).
// Su campo user_id era un remote_select MUERTO (sin endpoint). Su
// constructor, su llamada en screen_instances.go, su mapping en
// resource_screens.go y la constante L4_SCREEN_INST_USER_ROLES_ID
// (UUID …00d3) se eliminaron. UUID …00d3 queda libre para reuso futuro.

// messaging (plan 025 F5): pantalla NATIVA de mensajería del staff hacia las
// familias (Compose, NO SDUI). MainScreen intercepta el screen_key `messaging`
// y pinta Route.Messaging directamente; el slot_data NUNCA se renderiza por el
// SDUI engine. Esta screen_instance existe SOLO para satisfacer la FK
// resource_screens.screen_key → screen_instances.screen_key y para que el menú
// resuelva el screen_key. Se conserva un slot_data mínimo y válido
// (list-basic-v1) por higiene; SIN `api_prefix` (la pantalla nativa habla con
// la API messaging por su propio cliente, no por el contrato SDUI genérico, así
// que no se proyecta bloque "contract").
//
// requiredPermission (slot.permission de la pantalla) = messaging.view: gatea
// tanto el item de menú como el acceso a la pantalla. El wildcard `messaging.*`
// (school_admin/teacher) lo cubre — no se enumera por rol (wildcard-first).
// scope=system, coherente con el recurso messaging (la capability no se ata a
// una escuela; el alcance lo da el JWT).
func messaging() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MESSAGING_ID,
		screenKey:          "messaging",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Mensajería",
		description:        "Mensajería del staff hacia las familias (pantalla nativa)",
		scope:              "system",
		requiredPermission: "messaging.view",
		slotData: `{
  "title": "Mensajería",
  "columns": [
    {"key": "recipient", "label": "Destinatario"},
    {"key": "status", "label": "Estado"}
  ],
  "actions_removed": ["create", "edit", "delete"]
}`,
	}
}
