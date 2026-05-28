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
  "api_prefix": "academic"
}`,
	}
}

func rolesList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ROLES_LIST_ID,
		screenKey:          "roles-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Lista de Roles",
		description:        "Roles del sistema",
		scope:              "system",
		requiredPermission: "admin.roles.read",
		slotData: `{
  "title": "Roles",
  "search_placeholder": "Buscar rol...",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "code", "label": "Código"},
    {"key": "scope", "label": "Alcance"}
  ],
  "api_prefix": "identity:",
  "api_base_path": "/api/v1/roles",
  "resource": "roles",
  "form_screen_key": "roles-form"
}`,
	}
}

func rolesForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ROLES_FORM_ID,
		screenKey:          "roles-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Rol",
		description:        "Crear o editar un rol",
		scope:              "system",
		requiredPermission: "admin.roles.read",
		slotData: `{
  "title": "Rol",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "code", "label": "Código", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "scope", "label": "Alcance", "type": "select", "options": [
      {"value": "system", "label": "Sistema"},
      {"value": "school", "label": "Escuela"},
      {"value": "unit", "label": "Unidad"}
    ]}
  ],
  "api_prefix": "identity:",
  "api_base_path": "/api/v1/roles",
  "resource": "roles",
  "list_screen_key": "roles-list"
}`,
	}
}

func permissionsList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_PERMISSIONS_LIST_ID,
		screenKey:          "permissions-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Gestión de Permisos",
		description:        "Permisos del sistema",
		scope:              "system",
		requiredPermission: "admin.permissions_mgmt.read",
		slotData: `{
  "title": "Permisos",
  "search_placeholder": "Buscar permiso...",
  "filter_ready_label": "Activos",
  "filter_processing_label": "Inactivos",
  "columns": [
    {"key": "name", "label": "Permiso"},
    {"key": "resource", "label": "Recurso"},
    {"key": "action", "label": "Acción"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "identity:",
  "api_base_path": "/api/v1/permissions",
  "resource": "permissions_mgmt",
  "form_screen_key": "permissions-form"
}`,
	}
}

func permissionsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_PERMISSIONS_FORM_ID,
		screenKey:          "permissions-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Permiso",
		description:        "Crear o editar un permiso",
		scope:              "system",
		requiredPermission: "admin.permissions_mgmt.read",
		slotData: `{
  "title": "Permiso",
  "fields": [
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "resource_id", "label": "Recurso", "type": "remote_select", "required": true},
    {"key": "action", "label": "Acción", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "identity:",
  "api_base_path": "/api/v1/permissions",
  "resource": "permissions_mgmt",
  "list_screen_key": "permissions-list"
}`,
	}
}

// ===============================================================
// ADMIN: screen-config (templates + instances + endpoint screens-form)
// ===============================================================

func screenTplList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCREEN_TPL_LIST_ID,
		screenKey:          "screen-templates-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Templates de Pantalla",
		description:        "Lista de templates base disponibles",
		scope:              "system",
		requiredPermission: "admin.screen_templates.read",
		slotData: `{
  "title": "Templates",
  "search_placeholder": "Buscar template...",
  "columns": [
    {"key": "name", "label": "Nombre"},
    {"key": "pattern", "label": "Patrón"},
    {"key": "version", "label": "Versión"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "platform:"
}`,
	}
}

func screenInstList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCREEN_INST_LIST_ID,
		screenKey:          "screen-instances-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Instancias de Pantalla",
		description:        "Lista de instancias configuradas de pantalla",
		scope:              "system",
		requiredPermission: "admin.screen_instances.read",
		slotData: `{
  "title": "Instancias de Pantalla",
  "search_placeholder": "Buscar instancia...",
  "columns": [
    {"key": "screen_key", "label": "Key"},
    {"key": "name", "label": "Nombre"},
    {"key": "scope", "label": "Alcance"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "identity"
}`,
	}
}

func screenInstForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCREEN_INST_FORM_ID,
		screenKey:          "screen-instances-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Instancia",
		description:        "Crear o editar una instancia de pantalla",
		scope:              "system",
		requiredPermission: "admin.screen_instances.read",
		slotData: `{
  "title": "Instancia de Pantalla",
  "fields": [
    {"key": "screen_key", "label": "Screen Key", "type": "text", "required": true},
    {"key": "template_id", "label": "Template", "type": "remote_select", "required": true},
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "scope", "label": "Alcance", "type": "select", "options": [
      {"value": "system", "label": "Sistema"},
      {"value": "school", "label": "Escuela"},
      {"value": "unit", "label": "Unidad"}
    ]},
    {"key": "required_permission", "label": "Permiso requerido", "type": "text"},
    {"key": "is_active", "label": "Activa", "type": "toggle"}
  ],
  "api_prefix": "identity"
}`,
	}
}

// screens-form (alias legacy): mismo concepto que screen-instances-form
// pero el FE lo declara con resource=`screens` (resource API-only en
// L4 B1). api_prefix=platform documentado en TC-C como acepted.
func screensForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_SCREENS_FORM_ID,
		screenKey:          "screens-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Nueva Screen Instance",
		description:        "Formulario alias para crear screen instances (variante platform)",
		scope:              "system",
		requiredPermission: "admin.screen_instances.read",
		slotData: `{
  "title": "Nueva Screen Instance",
  "fields": [
    {"key": "screen_key", "label": "Screen Key", "type": "text", "required": true},
    {"key": "template_id", "label": "Template", "type": "remote_select", "required": true},
    {"key": "name", "label": "Nombre", "type": "text", "required": true},
    {"key": "scope", "label": "Alcance", "type": "select", "options": [
      {"value": "system", "label": "Sistema"},
      {"value": "school", "label": "Escuela"},
      {"value": "unit", "label": "Unidad"}
    ]}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "platform"
}`,
	}
}

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
  "filter_ready_label": "Info",
  "filter_processing_label": "Crítico",
  "columns": [
    {"key": "event_type", "label": "Tipo"},
    {"key": "actor", "label": "Actor"},
    {"key": "target", "label": "Recurso"},
    {"key": "created_at", "label": "Fecha"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
  "api_prefix": "identity"
}`,
	}
}

func auditDetail() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_AUDIT_DETAIL_ID,
		screenKey:          "audit-detail",
		templateID:         L0_SCREEN_TPL_DETAIL_ID_REF,
		name:               "Detalle de Auditoría",
		description:        "Detalle de un evento de auditoría",
		scope:              "system",
		requiredPermission: "admin.audit.read",
		slotData: `{
  "title": "Evento de Auditoría",
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
    {"key": "level", "label": "Nivel", "type": "text"},
    {"key": "period_id", "label": "Periodo", "type": "remote_select", "required": true}
  ],
  "api_prefix": "academic"
}`,
	}
}

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
  "api_prefix": "academic"
}`,
	}
}

// membershipsForm: form-basic-v1 con los campos del CreateMembershipRequest del
// backend. Las keys/tipos cuadran con el contrato real (academic):
//   - user_email (text): el backend acepta user_id O user_email; usamos el
//     email para no depender de un selector remoto de usuarios. Tipo `text`
//     (NO `email`) para evitar un ControlType incierto en el renderer.
//   - academic_unit_id (remote_select): el FormFieldsResolver del KMP DESCARTA
//     todo remote_select sin remote_endpoint, así que aquí SÍ lo declaramos.
//     Endpoint academic:/api/v1/units → {"units":[{id, display_name,...}]}; la
//     escuela se resuelve de la escuela activa del JWT (NUNCA por path/query/
//     body, estándar del ecosistema). El campo visible es
//     display_field=display_name (NO `name`), value_field=id.
//   - role_key (select estático): enum del backend (NO remote, NO role_id).
//     Misma forma textual que invitations-form: type "select" + options con
//     {value,label}.
// NO lleva subject_ids ni materias (retirado en F0b, no se reintroduce).
func membershipsForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MEMBERSHIPS_FORM_ID,
		screenKey:          "memberships-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Membresía",
		description:        "Asignar usuario a unidad",
		scope:              "school",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Membresía",
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
  "api_prefix": "academic"
}`,
	}
}

// subjectsList: hereda los default_actions de list-basic-v1
// (create/edit/delete sobre $resource$ → academic.subjects.*). Sin deltas:
// el patrón CRUD estándar es suficiente. La vista "alumnos por materia" no vive
// aquí sino embebida en subjects-form (master-detail), ver subjectsForm().
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

// subjectsForm usa master-detail-v1 (plan 006, Trozo A): hereda los 3
// defaults de form (save_new/save/delete con scope=form-submit) y, vía
// detail_config, embebe `students-by-subject-list` como tab/panel detalle
// (lista readonly de alumnos de la materia). Sin modal: modal_screen_key=null
// porque la lista es solo lectura.
//
// detail_config: parent_id_param="subjectId" → MasterDetailContainer carga
// students-by-subject-list pasando subjectId = id de la materia editada; el
// contrato KMP lee context.params["subjectId"] y llama al lector B
// (GET /api/v1/subjects/:id/enrollments). child_id_field="id". El frontend KMP
// interpreta detail_config; el backend solo lo persiste.
//
// actions_removed=["detail"]: el template master-detail-v1 trae un default
// `detail` (view-detail|$resource$.read|edit-only) pensado para navegar a un
// detalle full-screen. Aquí el detalle es el panel EMBEBIDO (no hay pantalla
// destino ni handler view-detail en SubjectsFormContract), así que el botón de
// toolbar no aplica y se retira intencionalmente.
//
// Reintroducido en N1.7 F2 sobre el modelo de sesiones (antes de F0b dependía
// del filtro subject_id sobre membership_subjects; ahora el lector B resuelve
// las inscripciones por sesión).
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
  "detail_config": {
    "screen_key": "students-by-subject-list",
    "modal_screen_key": null,
    "parent_id_param": "subjectId",
    "child_id_field": "id",
    "title": "Alumnos"
  },
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

// studentsBySubjectList (plan 006, N1.B): vista de "alumnos por materia".
// Espeja unit-directory (scope=unit, readonly). El permiso del slot es
// academic.memberships.read (única fuente de gateo, ADR-0003). Se alcanza como
// panel/tab detalle EMBEBIDO de subjects-form (master-detail-v1): el
// detail_config de subjects-form la carga pasando subjectId = id de la materia;
// el contrato KMP lee context.params["subjectId"] y llama al lector B
// (GET /api/v1/subjects/:id/enrollments), que devuelve la misma forma que
// GET /memberships. actions_removed retira create/edit/delete heredados del
// template (la pantalla es de solo lectura, igual que unit-directory).
// Reintroducida en N1.7 F2 sobre el modelo de sesiones.
func studentsBySubjectList() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_STUDENTS_BY_SUBJECT_ID,
		screenKey:          "students-by-subject-list",
		templateID:         L0_SCREEN_TPL_LIST_ID_REF,
		name:               "Alumnos por Materia",
		description:        "Listado de alumnos inscritos en una materia",
		scope:              "unit",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Alumnos",
  "search_placeholder": "Buscar alumno...",
  "columns": [
    {"key": "user_name", "label": "Nombre"},
    {"key": "role", "label": "Rol"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "readonly": true,
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
    {"key": "role", "label": "Rol"},
    {"key": "label", "label": "Etiqueta"},
    {"key": "uses_count", "label": "Usos"},
    {"key": "max_uses", "label": "Máx."},
    {"key": "is_active", "label": "Activa"},
    {"key": "expires_at", "label": "Expira"}
  ],
  "actions_removed": ["edit", "delete"],
  "actions_added": [
    {"id": "revoke", "scope": "row", "label": "Revocar", "icon": "ban", "permission": "academic.invitations.revoke", "condition": "always", "event_id": "revoke", "style": "destructive", "order": 20}
  ],
  "api_prefix": "academic"
}`,
	}
}

// invitations-form (N0.4-A): creación de un código de invitación.
// Solo create (no edit): patrón delta retira "save" (edit-only) y
// "delete"; conserva "save_new" → $resource$.create →
// academic.invitations.create. El campo `code` NO se incluye: lo
// genera el backend. academic_unit_id se llena vía remote_select de
// unidades del colegio (remoteSelectPrefix=academic en el contrato FE).
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
    {"key": "academic_unit_id", "label": "Unidad", "type": "remote_select", "required": true},
    {"key": "role", "label": "Rol", "type": "select", "required": true, "options": [
      {"value": "student", "label": "Estudiante"},
      {"value": "teacher", "label": "Profesor"},
      {"value": "guardian", "label": "Acudiente"}
    ]},
    {"key": "label", "label": "Etiqueta", "type": "text"},
    {"key": "expires_at", "label": "Expira", "type": "datetime"},
    {"key": "max_uses", "label": "Usos máximos", "type": "number"}
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
  "actions_removed": ["delete"],
  "api_prefix": "learning"
}`,
	}
}

func gradesForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_GRADES_FORM_ID,
		screenKey:          "grades-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Calificación",
		description:        "Registrar o editar una calificación",
		scope:              "unit",
		requiredPermission: "academic.grades.read",
		slotData: `{
  "title": "Calificación",
  "fields": [
    {"key": "student_id", "label": "Estudiante", "type": "remote_select", "required": true},
    {"key": "subject_id", "label": "Materia", "type": "remote_select", "required": true},
    {"key": "score", "label": "Nota", "type": "number", "required": true},
    {"key": "comment", "label": "Comentario", "type": "textarea"}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "learning"
}`,
	}
}

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
  "api_prefix": "learning"
}`,
	}
}

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
    {"key": "unit_id", "label": "Unidad", "type": "remote_select", "required": true},
    {"key": "date", "label": "Fecha", "type": "date", "required": true},
    {"key": "entries", "label": "Asistencias", "type": "table"}
  ],
  "actions_removed": ["save", "delete"],
  "api_prefix": "learning"
}`,
	}
}

// attendance-form (phantom legítimo): el FE
// (AttendanceFormContract.kt) declara este screenKey y el legacy lo
// tenia con MISMO slot que attendance-batch. Conservado aqui como
// variante simple (un solo estudiante).
func attendanceForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ATTENDANCE_FORM_ID,
		screenKey:          "attendance-form",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Formulario de Asistencia",
		description:        "Formulario individual de asistencia",
		scope:              "unit",
		requiredPermission: "academic.attendance.read",
		slotData: `{
  "title": "Asistencia",
  "fields": [
    {"key": "student_id", "label": "Estudiante", "type": "remote_select", "required": true},
    {"key": "date", "label": "Fecha", "type": "date", "required": true},
    {"key": "status", "label": "Estado", "type": "select", "options": [
      {"value": "present", "label": "Presente"},
      {"value": "absent", "label": "Ausente"},
      {"value": "late", "label": "Tarde"},
      {"value": "excused", "label": "Justificado"}
    ]}
  ],
  "actions_removed": ["delete"],
  "api_prefix": "learning"
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
  "api_prefix": "learning"
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
    {"key": "subject", "label": "Materia"},
    {"key": "scheduled_at", "label": "Fecha"}
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
// detail_config[] describe la navegación al panel detalle. El frontend
// KMP es quien lo interpreta; el backend solo lo persiste como blob.
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
  "fields": [
    {"key": "title", "label": "Título", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "pass_threshold", "label": "Umbral de aprobación (%)", "type": "number"},
    {"key": "max_attempts", "label": "Intentos máximos", "type": "number"},
    {"key": "time_limit_minutes", "label": "Tiempo límite (min)", "type": "number"},
    {"key": "is_timed", "label": "Cronometrada", "type": "boolean"},
    {"key": "shuffle_questions", "label": "Mezclar preguntas", "type": "boolean"},
    {"key": "show_correct_answers", "label": "Mostrar respuestas correctas", "type": "boolean"},
    {"key": "available_from", "label": "Disponible desde", "type": "datetime"},
    {"key": "available_until", "label": "Disponible hasta", "type": "datetime"}
  ],
  "detail_config": {
    "screen_key": "assessment-questions-list",
    "modal_screen_key": "assessment-question-form",
    "parent_id_param": "assessmentId",
    "child_id_field": "id"
  },
  "actions_added": [
    {"id": "detail",  "scope": "resource-toolbar", "icon": "help_outline", "label": "Preguntas", "permission": "content.assessments.read",   "condition": "edit-only", "event_id": "view-questions", "style": "icon", "order": 15},
    {"id": "publish", "scope": "resource-toolbar", "icon": "check_circle", "label": "Publicar",  "permission": "content.assessments.update", "condition": "edit-only", "event_id": "publish",        "style": "icon", "order": 30},
    {"id": "archive", "scope": "resource-toolbar", "icon": "archive",      "label": "Archivar",  "permission": "content.assessments.update", "condition": "edit-only", "event_id": "archive",        "style": "icon", "order": 40}
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
    {"key": "subject", "label": "Materia"},
    {"key": "modality", "label": "Modalidad"},
    {"key": "scheduled_at", "label": "Fecha"},
    {"key": "status", "label": "Estado"}
  ],
  "api_prefix": "learning"
}`,
	}
}

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
    {"key": "statement", "label": "Pregunta"},
    {"key": "kind", "label": "Tipo"},
    {"key": "score", "label": "Puntaje"}
  ],
  "actions_added": [
    {"id": "create", "scope": "header", "label": "Nuevo",    "icon": "plus",  "permission": "content.assessments.update", "condition": "always", "event_id": "create", "style": "icon",        "order": 10},
    {"id": "delete", "scope": "row",    "label": "Eliminar", "icon": "trash", "permission": "content.assessments.update", "condition": "always", "event_id": "delete", "style": "destructive", "order": 20}
  ],
  "api_prefix": "learning"
}`,
	}
}

func assessmentQuestionForm() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:          L4_SCREEN_INST_ASSESS_QUESTION_FORM_ID,
		screenKey:   "assessment-question-form",
		templateID:  L0_SCREEN_TPL_FORM_ID_REF,
		name:        "Formulario de Pregunta",
		description: "Crear o editar una pregunta",
		scope:       "unit",
		// TC-A del baseline.
		requiredPermission: "content.assessments.read",
		slotData: `{
  "title": "Pregunta",
  "page_title": "Pregunta",
  "edit_title": "Editar pregunta",
  "fields": [
    {"key": "question_text", "label": "Enunciado", "type": "textarea", "required": true},
    {"key": "question_type", "label": "Tipo", "type": "select", "required": true, "options": [
      {"value": "multiple_choice", "label": "Opción múltiple"},
      {"value": "true_false", "label": "Verdadero/Falso"},
      {"value": "short_answer", "label": "Respuesta corta"},
      {"value": "open_ended", "label": "Respuesta abierta"}
    ]},
    {"key": "points", "label": "Puntaje", "type": "number", "required": true},
    {"key": "correct_answer", "label": "Respuesta correcta", "type": "text"},
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

// assessmentAssignment: phantom-nueva. Pantalla para asignar una
// evaluación creada a las unidades destino.
func assessmentAssignment() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_ASSIGNMENT_ID,
		screenKey:          "assessment-assignment",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Asignación de Evaluación",
		description:        "Asignar una evaluación a unidades destino",
		scope:              "unit",
		requiredPermission: "content.assessments.update",
		slotData: `{
  "title": "Asignar Evaluación",
  "fields": [
    {"key": "assessment_id", "label": "Evaluación", "type": "remote_select", "required": true},
    {"key": "units", "label": "Unidades", "type": "multi_select", "required": true},
    {"key": "starts_at", "label": "Inicio", "type": "datetime"},
    {"key": "ends_at", "label": "Fin", "type": "datetime"}
  ],
  "actions_removed": ["save", "delete"],
  "actions_added": [
    {"id": "save_new", "scope": "form-submit", "label": "Asignar", "icon": "save", "permission": "content.assessments.update", "condition": "create-only", "event_id": "submit-form", "style": "filled", "order": 10}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessmentModality: phantom-nueva. Selector previo al form de
// creación de evaluación (modalidad: quiz, examen, tarea).
func assessmentModality() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_ASSESS_MODALITY_ID,
		screenKey:          "assessment-modality",
		templateID:         L0_SCREEN_TPL_DETAIL_ID_REF,
		name:               "Modalidad de Evaluación",
		description:        "Selección de modalidad antes de crear una evaluación",
		scope:              "unit",
		requiredPermission: "content.assessments.create",
		slotData: `{
  "title": "Modalidad de Evaluación",
  "options": [
    {"value": "quiz", "label": "Quiz", "icon": "zap"},
    {"value": "exam", "label": "Examen", "icon": "clipboard"},
    {"value": "assignment", "label": "Tarea", "icon": "file-text"}
  ],
  "api_prefix": "learning"
}`,
	}
}

// assessmentReviewDashboard: phantom-nueva. Dashboard de revisión
// de intentos por evaluación (docente).
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

// assignedAssessmentsList: phantom-nueva. Lista para el estudiante
// de las evaluaciones asignadas.
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
    {"key": "subject", "label": "Materia"},
    {"key": "due_at", "label": "Vence"},
    {"key": "status", "label": "Estado"}
  ],
  "actions_removed": ["create", "edit", "delete"],
  "api_prefix": "learning"
}`,
	}
}

// attemptReviewDetail: phantom-nueva. Detalle de un intento de
// evaluación (vista revisión).
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
// REPORTS: progress / stats / report-card
// ===============================================================
//
// Poda F2 (plan 004-permisologia-mvp): progress-detail, stats-detail y
// report-card se eliminaron del MVP junto con sus constantes y filas en
// resource_screens.go. Los dashboards progress-dashboard /
// stats-dashboard SÍ se conservan (definidos en screen_instances.go).

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

// membership-add: form simplificado para vincular un usuario a una
// unidad sin pasar por el form completo. El FE lo usa como flow
// abreviado (`MembershipAddContract.kt`).
func membershipAdd() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_MEMBERSHIP_ADD_ID,
		screenKey:          "membership-add",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Vincular Miembro",
		description:        "Form simplificado de vinculación usuario-unidad",
		scope:              "school",
		requiredPermission: "academic.memberships.read",
		slotData: `{
  "title": "Vincular Miembro",
  "fields": [
    {"key": "email", "label": "Email del usuario", "type": "email", "required": true},
    {"key": "role_id", "label": "Rol", "type": "remote_select", "required": true}
  ],
  "actions_removed": ["save", "delete"],
  "api_prefix": "academic"
}`,
	}
}

// ===============================================================
// Fase 3 (B7b) — Demo CRUD data-driven sin Kotlin (colors)
// ===============================================================
//
// Poda F2 (plan 004-permisologia-mvp): los constructores de la pareja
// demo `colors-list` / `colors-form` se eliminaron del MVP junto con
// sus constantes y filas en resource_screens.go. El recurso
// platform.colors queda huérfano (prune-later, ver
// docs/plans/004-permisologia-mvp/diferido.md).

// user-roles: pantalla para asignar/revocar roles de un usuario.
// TC-A del baseline: resource=users y permisos users:read/update (no
// `user_roles:*` que no existen en el seed).
func userRoles() l4ScreenInstanceRow {
	return l4ScreenInstanceRow{
		id:                 L4_SCREEN_INST_USER_ROLES_ID,
		screenKey:          "user-roles",
		templateID:         L0_SCREEN_TPL_FORM_ID_REF,
		name:               "Roles del Usuario",
		description:        "Asignación de roles a un usuario",
		scope:              "system",
		requiredPermission: "admin.users.read",
		slotData: `{
  "title": "Roles del Usuario",
  "fields": [
    {"key": "user_id", "label": "Usuario", "type": "remote_select", "required": true},
    {"key": "roles", "label": "Roles", "type": "multi_select", "required": true}
  ],
  "actions_removed": ["save_new", "save", "delete"],
  "actions_added": [
    {"id": "assign-role", "scope": "form-submit", "label": "Asignar", "icon": "plus",  "permission": "admin.users.update", "condition": "always", "event_id": "assign-role", "style": "filled",      "order": 10},
    {"id": "revoke-role", "scope": "form-submit", "label": "Revocar", "icon": "minus", "permission": "admin.users.update", "condition": "always", "event_id": "revoke-role", "style": "destructive", "order": 20}
  ],
  "api_prefix": "identity"
}`,
	}
}
