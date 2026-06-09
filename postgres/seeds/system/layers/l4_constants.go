package layers

// L4_SEED_VERSION declara la versión semántica del contenido de L4.
// Bumpear en CADA cambio de dato visible en cualquier sub-archivo
// de seeds/system/l4/ (resources, roles_permissions, etc.).
//
// Historial:
//   - 1.11.0: baseline previo al refactor SDUI por dominio.
//   - 1.12.0: migración del slot_data de 4 screen_instances de
//     assessments (assessments-form, assessments-management-list,
//     assessment-questions-list, assessment-question-form) al
//     nuevo estándar SDUI (zonas con scope + actions declarativos:
//     page_title/edit_title, save_new+save desdoblados con
//     condition create-only/edit-only, destructive flag). Cambio
//     adicional en resource_screens.go: assessments-management-list
//     pasa a ser is_default=true para el recurso `assessments`
//     (assessments-list deja de ser default). `assessments_student`
//     sigue con `assessments-list` como default.
//   - 1.13.0: corrección de bug de routing del menú docente.
//     Bundle expone `screens` indexadas por screen_type, y el
//     frontend KMP toma `screens["list"]` para navegar — ignora
//     is_default. Fix: bajo el recurso `assessments` (docente),
//     `screen_type="list"` ahora apunta a
//     `assessments-management-list` (master-detail CRUD).
//     Se eliminó la fila legacy que mapeaba `assessments-list` al
//     mismo recurso (esa pantalla pertenece a `assessments_student`,
//     no al flujo docente). Bug observado: al tocar "Evaluaciones"
//     en el menú con grants content.assessments.*, el menú abría
//     student-take en lugar de management-list.
//   - 1.14.0: (revertido en 1.15.0) — incorporó 4 actions extra al
//     slot_data de assessments-form (view_questions/publish/archive/
//     assign) más cambios visuales SDUI en KMP (icon-button para
//     scope=form, tint destructive). Resultó en parches: contrato
//     icon-name sin validar, sin overflow strategy en zonas
//     ACTION_GROUP, style hardcoded por if, mezcla semántica de
//     form-submit y resource-toolbar. Snapshot 002 anota el detalle.
//   - 1.15.0: rollback del intento 1.14.0. assessmentsForm vuelve a
//     save_new + save + delete (estado pre-fix botón faltante). Las
//     4 actions extra se reincorporarán bajo el plan arquitectónico
//     pendiente (separación form-submit vs resource-toolbar +
//     catálogo icon-name validado + tabla style→token).
//   - 1.17.0: Fase 3 (B7b) — demo CRUD data-driven sin Kotlin. Se
//     siembran 2 screen_instances nuevas (`colors-list`, `colors-form`)
//     y 1 recurso de menú (`colors` bajo admin) + 2 mappings en
//     resource_screens (list/form). slot_data declara la metadata
//     SDUI (`api_prefix`, `api_base_path`, `resource`, `*_screen_key`)
//     que el composer proyecta como bloque `contract` para el fallback
//     `GenericListContract`/`GenericFormContract` del frontend KMP.
//     Endpoint `/api/v1/colors` y permisos `platform.colors.{create,
//     read,update,delete}` ya existen en edugo-api-platform/edugo-shared
//     (Bloque 7a). Ningún Contract.kt nuevo en el frontend.
//   - 1.16.0: Fase 3a — assessmentsForm migra de form-basic-v1 a
//     master-detail-v1. slot_data pasa a modelo declarativo
//     defaults+added/removed:
//   - Templates form-basic-v1 / list-basic-v1 / master-detail-v1
//     declaran `default_actions[]` con placeholder `$resource$`
//     (resuelve a "content.assessments" en este caso) → save_new
//     / save / delete con scope=form-submit + detail con
//     scope=resource-toolbar.
//   - assessmentsForm elimina la lista legacy `actions:[...]` y
//     declara `actions_added`: detail (override del default —
//     label "Preguntas", event_id view-questions, icon
//     help_outline), publish y archive — todas con
//     scope=resource-toolbar.
//   - `detail_config` apunta a assessment-questions-list /
//     assessment-question-form con parent_id_param=assessmentId.
//     El frontend NUNCA ve actions_added/actions_removed: el composer
//     en api-platform los expande sobre defaults antes del response.
//     Pantallas con `actions:[...]` legacy (announcement-form,
//     users-form, etc.) siguen idénticas (override total).
//   - 1.17.2: TECH_DEBT_BOTONERA #19 — colorsForm() actions corregidas.
//     Se añaden scope/condition/event_id/style/order a las 3 actions
//     (save_new, save_existing, delete) para que el SlotBindingResolver
//     las expanda correctamente en la zona form-submit del template
//     form-basic-v1. La causa raíz era que el legacy "event" (mayúsculas)
//     no es leído por el resolver; los campos faltantes dejaban la zona
//     form_submit vacía y los botones no se renderizaban.
//   - 1.17.3: F1 (ADR-6) herencia de roles. Los 5 alias que sí heredan
//     (school_director/coordinator/assistant → school_admin;
//     assistant_teacher/observer → teacher) reciben parent_role_id y
//     dejan de declarar grants propios en role_grants. readonly_auditor
//     conserva su allow/deny standalone. Los grants efectivos aplanados
//     no cambian (la herencia se resuelve en el login).
//   - 1.18.0: F2 (plan 004-permisologia-mvp) — poda del seed SDUI. Se
//     retiran 13 screen_instances y sus filas en resource_screens:
//     guardian (guardian-relations-list/form, guardian_relations-form
//     alias, guardian-requests-list), horarios (schedules-list/form),
//     calendario (calendar-list/form), demo (colors-list/form) y
//     reportes detalle (stats-detail, progress-detail, report-card).
//     Se elimina además el template L4 master-detail-basic-v1 (0
//     instancias). Los recursos academic.guardian_relations,
//     academic.schedules, academic.calendar, platform.colors y reports
//     quedan huérfanos (prune-later — NO se tocan iam.resources ni
//     iam.role_grants en esta pasada). Se conservan los dashboards
//     progress-dashboard / stats-dashboard y todo el flujo de examen.
//   - 1.19.0: F3 (plan 004) — estándar de pantallas. ~34 screen_instances
//     L4 migradas al patrón delta (heredan default_actions del template +
//     override puntual con actions_added/removed). Incluye las 5 ex-legacy
//     (attendance-list/batch, assessment-assignment, assessment-questions-
//     list, user-roles) ahora como delta con override explícito. CERO
//     instancias en formato actions legacy. Sin cambio semántico: el harness
//     TestScreenActionsInvariantRoundTrip garantiza set {event_id,permission}
//     idéntico. resource_screens NO se toca (inferencia descartada: la tabla
//     es load-bearing — codifica screen_type/is_default/N:M).
//   - 1.20.0: N0.0 (plan 005, onboarding) — 3 recursos nuevos
//     (invitations + join_requests visibles bajo academic;
//     join_request_approvals API-only como namespace de permisos de
//     aprobación per-rol) + 9 permisos (invitations.{create,read,revoke},
//     join_requests.{create,read,reject},
//     join_request_approvals.{student,teacher,guardian} — la acción ES el
//     rol que se admite). Grants: teacher gana invitations.*,
//     join_requests.* y approvals.student (literal, NO el wildcard);
//     school_admin ya cubre todo vía academic.*. readonly_auditor suma
//     deny *.revoke, *.reject y approvals.* (deny-wins).
//   - 1.22.0: N0.4-B (plan 005) — bandeja de solicitudes pendientes. +1
//     fila en resource_screens: join_requests:list →
//     screen_key `join-requests-inbox` (is_default=true). El FE la
//     resuelve con una pantalla Compose NATIVA (no SDUI), por eso NO se
//     siembra screen_instance: el resolver solo necesita que el menú
//     exponga el screen_key. El item aparece para quien tenga
//     `academic.join_requests.read` (school_admin vía academic.*,
//     teacher vía join_requests.*).
//   - 1.25.0: N1.B (plan 006) — vista docente "alumnos por materia".
//     +1 screen_instance `students-by-subject-list` (scope=unit, readonly,
//     espeja unit-directory; navegación-only desde subjects-list, NO se
//     mapea en resource_screens). subjects-list suma una acción de fila
//     `view-students` (actions_added, event_id view-students, permission
//     academic.memberships.read) — aditiva, no reemplaza el tap de editar.
//     teacher gana el grant LITERAL `academic.memberships.read` (no el
//     wildcard: el docente lee membresías para ver alumnos pero no las
//     muta). Golden de screen actions actualizado: subjects-list suma
//     `view-students|academic.memberships.read`; students-by-subject-list
//     entra con set vacío.
//   - 1.26.0: N1.C (plan 006) — "mis materias" del alumno. +1 permiso
//     `academic.my_memberships.read:own` (scope=unit, resource my_memberships).
//     Grant LITERAL al rol student (NO el wildcard `academic.memberships.*`).
//     Es el permiso ÚNICO del feature self del alumno: visibilidad del item
//     de menú "Mis materias", slot.permission de my-memberships-list y route
//     gate del dato. Vive bajo path propio (academic.my_memberships.*) para
//     que el gate de menú por path-prefix NO le filtre el item admin
//     "memberships" (roster de unidad). Habilita que el alumno lea SOLO sus
//     propias membresías vía GET /users/:user_id/memberships (self-check
//     path==token en el backend); sigue sin poder listar la unidad
//     (GET /memberships exige `academic.memberships.read` amplio).
//   - 1.27.0: Trozo A (plan 006) — subjects-form pasa a master-detail-v1.
//     El tab/panel detalle embebe `students-by-subject-list` (alumnos de la
//     materia, readonly) vía detail_config (parent_id_param=subjectId,
//     child_id_field=id, modal_screen_key=null → solo lectura, sin alta/baja
//     todavía; eso es Trozo B). Se RETIRA la acción de fila standalone
//     `view-students` de subjects-list (la lista de alumnos ya no es pantalla
//     suelta, es el detalle embebido); students-by-subject-list se conserva
//     como destino del detail_config. El default `detail` del template
//     master-detail se quita vía actions_removed=["detail"] (no hay detalle
//     full-screen en Trozo A). teacher: `academic.subjects.*` → grant LITERAL
//     `academic.subjects.read` (el docente ve materias pero no las gestiona;
//     CRUD de materias es de school_admin). Golden: subjects-list pierde
//     `view-students|academic.memberships.read`; subjects-form mantiene su set
//     invariante (master-detail con `detail` removido = mismas 3 acciones de
//     form que el form-basic anterior).
//   - 1.29.0: N1.7 F0a etapa 1 (plan 010, ADR 0009) — recurso nuevo
//     `subject_offerings` (b4000000-…-23, bajo academic, IsMenuVisible=false:
//     aún sin screen_instance) + 5 permisos academic.subject_offerings.
//     {create,read,update,delete,enroll} (scope school). Grants: school_admin
//     ya cubierto por wildcard `academic.*` (sin enumerar, wildcard-first);
//     teacher gana literal `academic.subject_offerings.read` (paridad con
//     `academic.subjects.read`). Sin cambios en pantallas. Acompaña el DDL
//     aditivo de subject_offerings / subject_offering_enrollments.
//   - 1.30.0: N1.7 F0b (plan 010, ADR 0009; Opción A) — retiro del catálogo
//     ligado a la tabla `membership_subjects` (eliminada). Se quitan: la
//     screen_instance `my-memberships-list` ("Mis materias" del alumno) con su
//     recurso `my_memberships`, su mapping en resource_screens, el permiso
//     `academic.my_memberships.read:own` y su grant al student; la
//     screen_instance `students-by-subject-list` ("alumnos por materia",
//     navegación-only). `subjects-form` vuelve de master-detail-v1 a
//     form-basic-v1 (se desembebe el detail_config de alumnos por materia).
//     `memberships-form` pierde el campo `subject_ids`. El grant teacher
//     `academic.memberships.read` se CONSERVA (usos vivos roster/unit-directory).
//     Golden de screen actions: se eliminan las entradas `my-memberships-list`
//     y `students-by-subject-list`; `subjects-form` mantiene su set invariante.
//   - 1.31.0: N1.7 F1/F2 — REINTRODUCCIÓN del catálogo retirado en 1.30.0,
//     ahora apuntando al modelo de sesiones. Vuelven: el recurso de menú
//     `my_memberships` ("Mis materias" del alumno), su screen_instance
//     `my-memberships-list`, su mapping en resource_screens, el permiso
//     `academic.my_memberships.read:own` y su grant al student (el contrato KMP
//     consume el lector A GET /api/v1/me/subject-offerings); la screen_instance
//     `students-by-subject-list`, re-embebida en `subjects-form` (vuelve de
//     form-basic-v1 a master-detail-v1 con su detail_config), cuyo contrato KMP
//     consume el lector B GET /api/v1/subjects/:id/enrollments. `memberships-form`
//     NO recupera `subject_ids` (sigue retirado). El grant teacher
//     `academic.memberships.read` no se toca. Golden de screen actions: se
//     re-añaden `my-memberships-list` y `students-by-subject-list` (ambos {}).
//   - 1.32.0: fix de `memberships-form` (renderizaba vacío). El renderer KMP
//     (FormFieldsResolver) DESCARTA todo `remote_select` sin `remote_endpoint`,
//     y los campos viejos (user_id/unit_id/role_id) no lo tenían — el único con
//     endpoint era el `subject_ids` eliminado en F0b. Además las keys/tipos no
//     cuadraban con CreateMembershipRequest del backend. Se reescribe el
//     slot_data a las keys reales del contrato: `user_email` (text),
//     `academic_unit_id` (remote_select con remote_endpoint
//     academic:/api/v1/schools/{schoolId}/units, display_field=display_name,
//     value_field=id) y `role_key` (select estático con options del enum de
//     roles). Sin cambios de esquema/permisos. NO se reintroduce subject_ids.
//   - 1.33.0: `memberships-form`, campo `academic_unit_id` — el remote_endpoint
//     pasa de `academic:/api/v1/schools/{schoolId}/units` a
//     `academic:/api/v1/units`. La escuela se resuelve de la escuela activa del
//     JWT, NUNCA por path/query/body (estándar del ecosistema). Misma forma de
//     respuesta `{"units":[{id, display_name,...}]}`; display_field/value_field
//     sin cambios. Sin cambios de esquema/permisos.
//   - 1.33.1: nueva screen_instance `batch-enroll` (inscripción por lote,
//     pantalla NATIVA) bajo el recurso `subject_offerings` (N1.7 F1, plan 010 /
//     ADR 0009). requiredPermission=academic.subject_offerings.read; el botón
//     "Inscribir" se declara como action en slot_data con
//     permission=academic.subject_offerings.enroll (ADR 0003). Nuevo mapping en
//     resource_screens (subject_offerings → batch-enroll, list, default). El
//     recurso y los permisos ya se sembraron en F0a; aquí solo se consumen.
//   - 1.33.2: recurso `subject_offerings` pasa a IsMenuVisible=true (N1.7 F1):
//     ya existe screen_instance + mapping, así que el item de menú "Sesiones de
//     Materia" abre la pantalla batch-enroll (default del recurso). Sin esto la
//     pantalla quedaba inalcanzable desde el menú.
//   - 1.34.0: N1.7 F2 (plan 010 / ADR 0009) — pantallas de sesiones por materia.
//     Dos screen_instances nuevas bajo el recurso `subject_offerings`:
//     `enroll-one` (inscripción individual, pantalla NATIVA; action `enroll` con
//     permission academic.subject_offerings.enroll, ADR 0003) y
//     `sessions-by-subject-list` (lista hija SDUI; columnas
//     subject_name/section_label/period_name/teacher_name; readonly). Ambas con
//     requiredPermission=academic.subject_offerings.read y mapping en
//     resource_screens (no-default, sort 2 y 3; el default sigue siendo
//     batch-enroll). Además se añade la row-action `view-sessions` (scope=row,
//     permission academic.subject_offerings.read, icon event) a `subjects-list`,
//     que el FE enruta a `sessions-by-subject-list` con param subjectId. No se
//     tocan columnas ni otras actions de subjects-list ni de subjects-form
//     (F1 dejó su detail de alumnos intacto). Sin nuevos permisos (ya existen
//     academic.subject_offerings.*).
//   - 1.34.1: fix icono row-action view-sessions de subjects-list. El icon-name
//     "event" NO existe en IconCatalog del KMP y el renderer de row-actions hacía
//     hard-throw (crash de subjects-list). Se cambia a "list" (sí registrado en
//     IconCatalog → FormatListBulleted), semántica "ver lista de sesiones". Sin
//     cambios de permisos, columnas ni otras actions. Acompaña el fix de causa
//     raíz en ListPatternRenderer (row-actions ahora resuelven vía IconResolver
//     con fallback, no más throw).
//   - 1.35.0: N1.7 F2.2 — "Sesiones" como pestaña del master-detail subjects-form
//     y generalización del contrato a N paneles de detalle. Cambio de ESTRUCTURA
//     del contrato SDUI: la clave `detail_config` (objeto singular) se reemplaza
//     por `detail_configs` (array) en TODAS las screen_instances master-detail.
//     subjects-form ahora declara DOS detalles readonly: "Alumnos"
//     (students-by-subject-list) y "Sesiones" (sessions-by-subject-list), ambos
//     con parent_id_param=subjectId. assessments-form migra su detalle único
//     (assessment-questions-list + modal_screen_key=assessment-question-form) al
//     array de una entrada, preservando el modal. Se ELIMINA la row-action
//     `view-sessions` de subjects-list (la sesiones ahora se ven dentro del form
//     de materia, no navegando a una pantalla aparte). No hay nuevos permisos ni
//     columnas; sessions-by-subject-list y su mapping en resource_screens se
//     conservan (la pestaña carga esa pantalla). El singular `detail_config` deja
//     de existir en el seed y en el contrato KMP (sin legacy).
//   - 1.36.0: poda del menú — se eliminan 6 recursos sin pantalla KMP
//     implementada (el menú los listaba pero al abrirlos daba "screen
//     instance not found"): `roles`, `permissions_mgmt`, `colors`,
//     `calendar`, `schedules` y `guardian_relations`. Se retiran de forma
//     coherente: las 4 screen_instances de roles/permissions
//     (roles-list/form, permissions-list/form) y sus mappings en
//     resource_screens; los permisos admin.roles.*, admin.permissions_mgmt.*,
//     academic.calendar.*, academic.schedules.* y academic.guardian_relations.*;
//     y los patterns de grant que los citaban en teacher/student/guardian.
//     Los recursos colors/calendar/schedules/guardian_relations ya tenían sus
//     pantallas retiradas (poda F2 plan 004) y quedaban huérfanos. El dashboard
//     del guardian (rol) NO se toca. Sin cambios de esquema.
//   - 1.37.0 (2026-05-29): scope=unit en los recursos memberships,
//     subjects y subject_offerings. Sus endpoints exigen unidad activa
//     (RequireActiveContext → 428 NO_ACTIVE_UNIT), así que el scope del
//     recurso debe ser "unit" para que el menú/contexto del frontend
//     pida unidad antes de abrirlos. Coherente con grades/attendance.
//   - 1.38.0 (2026-06-01): poda de los recursos `screen_templates` y
//     `screen_instances` (CRUD de configuración SDUI reimplementado en el
//     admin-tool de Go): se retiran sus 2 recursos, 8 permisos, 4
//     screen_instances (…30..…33) y 4 mappings en resource_screens.
//   - 1.39.0 (2026-06-02): N1.7 F2.3 — habilitada la creación/edición de sesiones
//     de materia desde la app. Nueva screen_instance `sessions-by-subject-form`
//     (…5f) como modal del master-detail subjects-form (detail_configs[] de la
//     pestaña "Sesiones" pasa modal_screen_key null→sessions-by-subject-form) y
//     su mapping en resource_screens (…6d, recurso subject_offerings).
//   - 1.40.0 (2026-06-02): ADR 0016 punto 3 — scope de la screen_instance
//     `sessions-by-subject-list` corregido school → unit. La materia es catálogo
//     de escuela, pero la gestión de sus sesiones la filtra el backend por la
//     unidad activa del token; el scope declarado ahora refleja ese filtro real
//     (queda alineado con `students-by-subject-list`, ya en scope=unit). El
//     gating por RequiredContext del menú no cambia (deriva de resource.scope,
//     no de esta screen_instance).
//   - 1.41.0 (2026-06-02): el detalle de materia queda SOLO con la pestaña
//     "Sesiones". Cambios: (C) se quita la entrada "Alumnos"
//     (students-by-subject-list) del detail_configs de `subjects-form` — un
//     alumno se inscribe en una SESIÓN, no en la materia, así que el roster se
//     gestiona dentro de cada sesión (batch-enroll/enroll-one), no a nivel
//     materia. La screen_instance `students-by-subject-list` se ELIMINA por
//     completo (constructor, registro en el slice y constante
//     L4_SCREEN_INST_STUDENTS_BY_SUBJECT_ID, UUID …c1 libre): era SOLO ese panel
//     embebido, sin otra referencia (no estaba en menú ni en resource_screens).
//     (D) la screen_instance `sessions-by-subject-form` corrige su scope
//     school → unit: el form gestiona UNA sesión (filtrada por unidad activa) y
//     su selector de docente requiere unidad activa, alineándose con
//     `sessions-by-subject-list`.
//   - 1.42.0 (2026-06-02): se retira el camino de CREACIÓN DIRECTA de membresías
//     (redundante con el flujo invitación→solicitud→doble-gate→aprobación, que ya
//     crea la membresía). Cambios: se ELIMINAN las screen_instances
//     `memberships-form` (…53) y `membership-add` (…d2) — sus constructores,
//     registros en el slice y constantes (UUIDs …53 y …d2 libres). Se quitan sus
//     2 mappings en resource_screens (recurso memberships): quedan `memberships-list`
//     y `unit-directory`. A `memberships-list` se le agrega
//     `actions_removed:["create"]` para que el admin ya no navegue al form
//     eliminado; se conservan edit/delete/expire. Leer/editar/expirar/borrar
//     membresías sigue intacto. Sin cambios de esquema ni de permisos.
//   - 1.42.1 (2026-06-03): se reparan los campos faltantes del form
//     `periods-form`. El DTO `CreatePeriodRequest` exige `type` (string enum) y
//     `academic_year` (int) como required, pero el form sólo declaraba
//     name/start_date/end_date/is_active → el backend respondía 400 "invalid
//     request body" al crear un período desde el iPad. Se AGREGAN dos campos al
//     slot_data (después de `name`): `type` (select, required, options =
//     semester/trimester/bimester/quarter, los 4 valores válidos del CHECK de la
//     tabla academic_periods) y `academic_year` (number, required). Sin cambios de
//     esquema ni de permisos.
//   - 1.42.2 (2026-06-03): el form `sessions-by-subject-form` ahora limita la
//     entrada del campo `section_label` a 10 caracteres (atributo `max_length: 10`).
//     El backend valida `section_label max=10`; antes el usuario escribía de más y
//     el guardado fallaba con 400. El SDUI gana soporte de `max_length`: el renderer
//     KMP (SlotRenderer.applyMaxLength) trunca la entrada en `onValueChange` para
//     PREVENIR la sobreescritura en lugar de mostrar el error tras guardar. Sin
//     cambios de esquema ni de permisos.
//   - 1.42.3 (2026-06-04): se repara el campo faltante del form `units-form`. El
//     DTO `CreateUnitRequest` exige `type` (string enum) como required, pero el
//     form sólo declaraba name/level/period_id → el backend respondía 400
//     "invalid request body" al crear una unidad desde el iPad/web. Se AGREGA un
//     campo al slot_data (después de `name`): `type` (select, required, options =
//     school/grade/class/section/club/department, los 6 valores válidos de
//     domain.ValidUnitTypes). Sin cambios de esquema ni de permisos.
//   - 1.42.4 (2026-06-04, plan 011 Eje C): se SANEA el contrato del form
//     `units-form`. Se QUITAN los campos `level` y `period_id` del slot_data: el
//     DTO `CreateUnitRequest` no los acepta (solo display_name/code/type/
//     description/parent_unit_id/metadata) y el contrato KMP `UnitsFormContract`
//     los descartaba silenciosamente → el form "mentía". El form queda con
//     name + type (ambos required). Sin cambios de esquema ni de permisos.
//   - 1.42.5 (2026-06-04): se repara el campo `academic_unit_id` del form
//     `invitations-form`. Estaba declarado como `remote_select required` pero SIN
//     `remote_endpoint`, así que el FormFieldsResolver del KMP lo DESCARTA (no se
//     renderiza el selector de unidad) → el submit iba sin `academic_unit_id` y el
//     backend respondía 400 "invalid request body" (el DTO lo exige con
//     binding:"required"). Se AGREGA el endpoint espejando `memberships-form`:
//     remote_endpoint=academic:/api/v1/units, display_field=display_name,
//     value_field=id. Sin cambios de esquema ni de permisos. Ver bug 0034.
//   - 1.42.6 (2026-06-05, PRE 1a tenant→JWT de asistencia): el endpoint
//     /attendance pasa a scope=unit (RequireActiveContext) y deriva la unidad
//     del JWT. (1) Se QUITA el campo tenant `unit_id` del slot_data del form
//     `attendance-batch` (la unidad sale del token, no es campo de form); el
//     form queda con date + entries. (2) Se ELIMINA por completo el screen
//     huérfano `attendance-form` (constructor, registro en el slice y constante
//     L4_SCREEN_INST_ATTENDANCE_FORM_ID): no estaba mapeado en resource_screens
//     y solo lo respaldaba el contrato KMP, también eliminado. Cierre del
//     latente bug 0034 (attendance-form.student_id) por eliminación. Sin
//     cambios de esquema ni de permisos.
//   - 1.42.7 (2026-06-05, N2 plan 008 — feature de asistencia): (1) se corrige
//     el `api_prefix` de las 3 instancias `attendance-*` (list/batch/summary) de
//     "learning" a "academic": la asistencia vive en la API academic (:8060) y
//     el contrato KMP ya usa `academic:` (D5). (2) Entry-point "Pasar lista" en
//     el form `subjects-form`: action `take-attendance` (scope resource-toolbar,
//     condition edit-only) que navega a `attendance-batch` con subjectId, gateada
//     por `academic.attendance.create` (D2). Sin cambios de esquema ni de
//     permisos (el permiso ya estaba sembrado).
//   - 1.42.8 (2026-06-05, N2.S2 plan 008 D5 — cierre): el form
//     `attendance-batch` (override nativo "pasar lista") declara la action de
//     submit `submit-batch` (scope header, permission academic.attendance.create,
//     event_id submit-batch) en su slot_data. Es el permiso del botón del
//     override nativo (ADR 0003), espejo de la action `enroll` de batch-enroll;
//     activa el gate cliente del botón (antes el permiso quedaba null porque el
//     seed no declaraba ninguna action de submit). El permiso ya estaba sembrado
//     (cubierto por el wildcard academic.attendance.* de teacher). Sin cambios de
//     esquema ni de permisos.
//   - 1.42.9 (2026-06-05, N2.S3 plan 008 — entry-points de consulta): el form
//     `subjects-form` suma dos actions de toolbar espejo de "take-attendance":
//     `view-attendance` ("Historial", icon history, order 21) y
//     `view-attendance-summary` ("Resumen", icon bar_chart, order 22). Ambas con
//     scope resource-toolbar, condition edit-only y permission
//     academic.attendance.read; navegan a las pantallas SDUI genéricas
//     attendance-list / attendance-summary pasando subjectId. El destino del
//     evento (event_id view-attendance / view-attendance-summary) vive en
//     SubjectsFormContract del KMP. El permiso ya estaba sembrado (cubierto por el
//     wildcard academic.attendance.* de teacher). Sin cambios de esquema ni de
//     permisos.
//   - 1.42.10 (2026-06-05, F0.5 plan 013 — bug 0034): corrige el `api_prefix` de
//     las pantallas `grades-list` y `grades-form` de "learning" a "academic". El
//     endpoint de notas vive en la API academic (:8060), no learning. Sin cambios
//     de esquema ni de permisos.
//   - 1.42.11 (2026-06-05, N3 F3 plan 013): seed SDUI de la pantalla nativa
//     "Poner notas". Agrega la action `put-grades` (entry-point en subjects-form,
//     scope resource-toolbar, condition edit-only, permission
//     academic.grades.create, event_id put-grades → NavigateTo("grades-batch",
//     {subjectId})) y el screen instance `grades-batch` (override nativo Compose,
//     scope unit, requiredPermission academic.grades.read, selector de período
//     remote_select a academic:/api/v1/periods). Espejo de attendance-batch (N2).
//     Los permisos ya estaban sembrados (cubiertos por el wildcard
//     academic.grades.* de teacher). Sin cambios de esquema.
//   - 1.44.0 (2026-06-06, N3 F4 plan 013 — consulta de notas): seed SDUI de la
//     consulta de notas. (1) Action `view-grades-summary` ("Resumen de notas",
//     icon pie_chart, scope row, condition always, permission academic.grades.read,
//     event_id view-grades-summary → grades-subject-summary) como 5ª row-action de
//     la card de sesión (sessions-by-subject-list). (2) Screen instance
//     `grades-subject-summary` (resumen de notas por sesión, vista docente;
//     readonly, scope unit, requiredPermission academic.grades.read; espejo de
//     attendance-summary) + su mapping en resource_screens (recurso grades,
//     screen_type summary). (3) Feature self del alumno "Mis notas": recurso
//     `my_grades`, permiso nuevo `academic.my_grades.read:own` con grant LITERAL
//     al rol student, screen instance `my-grades-list` (readonly, requiredPermission
//     academic.my_grades.read:own; el contrato KMP consume GET /api/v1/me/grades),
//     mapping resource_screens my_grades→my-grades-list (is_default) e item de menú
//     "Mis notas" (recurso my_grades, IsMenuVisible). Espejo de my_memberships. Sin
//     cambios de esquema. Los permisos de docente (academic.grades.read) ya estaban
//     sembrados.
//   - 1.45.0 (2026-06-06, N3 F4.1 — cierre deuda de privacidad, decisión del
//     dueño): se ELIMINA el grant amplio `academic.grades.*` del rol `student`.
//     Ese wildcard era CRUD docente y dejaba al alumno ver/crear/editar notas
//     ajenas vía GET/POST /grades y ver el menú "Calificaciones" (grades-list).
//     El alumno conserva el feature self `academic.my_grades.read:own` (1.44.0),
//     que sirve solo sus propias notas vía GET /api/v1/me/grades → su única vista
//     de notas pasa a ser "Mis Notas". El grant de `guardian` (`academic.grades.*`)
//     queda intacto (deuda separada: el acudiente necesita ver notas de sus
//     acudidos). Sin cambios de esquema (cambia solo el output del seed).
//   - 1.46.0 (2026-06-06, N4 F2.6 — alineación SDUI de evaluación al contrato
//     nuevo + field option-list): (1) assessment-question-form gana el field
//     {key:options, type:option-list, correct_answer_field:correct_answer} que
//     faltaba (bug original: el editor no mostraba opciones; lo consume el
//     DynamicOptionListField del KMP con shape {option_id, option_text}); se
//     quita el field correct_answer separado (lo marca el radio de la lista) y
//     se restringe question_type a los 4 tipos del CHECK nuevo. (2) assessments-form
//     gana subject_id (remote_select a academic /subjects, FK obligatoria del
//     esquema nuevo) + acción "Asignar" (event_id=assign, permiso
//     content.assessments.assign); se quita modality. (3) assessment-assignment
//     reescrita al contrato nuevo: target = subject_offering_id (remote_select a
//     /subject-offerings) + due_date opcional, NUNCA alumnos; slot.permission
//     content.assessments.assign. (4) listas (assessments-list/-management-list/
//     -questions-list/assigned-assessments-list) alineadas a los campos del
//     esquema nuevo (subject_name, status, questions_count, question_text/_type/
//     points, due_date). (5) assessment-modality ELIMINADA (concepto muerto: el
//     esquema nuevo no tiene modalidad). take/result/review-dashboard/
//     attempt-review-detail quedan MÍNIMAS (F3, re-apuntado de UI pendiente). Sin
//     cambios de esquema (cambia solo el output del seed).
//   - 1.47.0 (2026-06-06, N4 F4.6 — catálogo del modo detallado de notas): se
//     siembran en iam.permissions los 4 permisos del recurso grades_detail
//     (academic.grades_detail.create/read/update/delete), espejando el enum de
//     edugo-shared. Cuelgan de un recurso PROPIO grades_detail (…37, NO
//     menú-visible): no comparten resource_id con `grades` porque el unique
//     (resource_id, action) prohíbe repetir create/read/update. Gestionan los
//     componentes de nota (academic.grade_item) y el desglose transparente en
//     "Mis Notas". NO se otorgan a roles vía roleGrantPatterns: el grant es
//     CONDICIONAL por perfil de escuela (academic.schools.grade_profile) y lo
//     inyecta identity en runtime (F4.5). Sin cambios de esquema (cambia solo el
//     output del seed).
//   - 1.48.0 (2026-06-08, Fase 1 — visibilidad condicional de campos SDUI en el
//     form de pregunta): assessment-question-form se vuelve REACTIVO por
//     `question_type` (campo controlador). Nuevo contrato SDUI snake_case
//     `visible_when` ({field, equals|in}; ausencia = siempre visible) en los
//     campos de respuesta correcta: `options` (option-list) solo en
//     multiple_choice, su `correct_answer_field` pasa de `correct_answer` a
//     `mc_correct_letter`; NUEVO field `correct_answer_bool` (select
//     Verdadero/Falso, required) visible en true_false; NUEVO field
//     `correct_answer_text` (text, required) visible en short_answer; open_ended
//     no muestra campo de respuesta correcta. question_text/question_type/points/
//     explanation/difficulty siguen siempre visibles. Contrato compartido con el
//     agente FRONT del KMP. Sin cambios de esquema (cambia solo el output del
//     seed).
//   - 1.49.0 (2026-06-08, Fase 2 — nuevo tipo de pregunta multiple_select):
//     assessment-question-form gana soporte para opción múltiple con VARIAS
//     respuestas correctas. (1) El dropdown `question_type` suma la opción
//     {value: multiple_select, label: "Opción múltiple (varias)"}. (2) Nuevo
//     field `options_multi` (type=option-list, selection_mode=multiple,
//     correct_answer_field=ms_correct_letters, visible_when question_type in
//     [multiple_select], NO required) — key DISTINTA de `options` (single) para
//     no colisionar el estado del componente. Contrato de datos: para este tipo
//     assessment.question.correct_answer guarda un ARRAY JSON de textos; NO se
//     añade is_correct a question_option. Acompaña el cambio de esquema (CHECK
//     question_type_check suma 'multiple_select', SchemaVersion → 3.53.0).
//     Contrato compartido con backend learning y FRONT del KMP. Solo autoría.
//   - 1.50.0: assessment-questions-list suma "actions_removed": ["edit"]. El
//     template list-basic-v1 declara `edit` (scope row) como default_action; en
//     el detalle de preguntas la edición ya la cubre el botón nativo "Editar"
//     del bottom-sheet (MasterDetailContainer, flujo N3.5), así que la row-action
//     SDUI `edit` quedaba huérfana (sin handler → "No custom handler for event:
//     edit"). Se elimina el duplicado igual que en las listas de sesiones.
//   - 1.51.0 (2026-06-08): dos ajustes de evaluación. (1) assessments-form: la
//     action "Publicar" (event_id=publish) alinea su slot.permission de
//     content.assessments.update → content.assessments.publish, para igualar el
//     gate del botón con la ruta POST /api/v1/assessments/:id/publish
//     (RequirePermission(PermissionAssessmentsPublish)). El rol teacher ya cubre
//     publish vía wildcard content.assessments.* (no cambian roles). (2) Se
//     ELIMINA la pantalla SDUI assessment-assignment (form-basic-v1): la
//     asignación a una sesión de materia pasa a un modal NATIVO ("nativa
//     prevalece, SDUI solo guía"). Se quita su screen_instance, su mapping en
//     resource_screens y su constante; se conserva el recurso assessments y el
//     permiso content.assessments.assign (lo gatean la action "Asignar" del form
//   - la ruta de assignments).
//   - 1.52.0 (2026-06-09, plan 017 F2 — picker de entidad): assessments-form
//     migra el campo `subject_id` de `remote_select` a `entity-picker` (control
//     nuevo del plan 017). El selector de materia abre un modal con búsqueda
//     server-side + paginación contra academic:/api/v1/subjects (search_param=
//     search, page_size=20) en vez de cargar todas las opciones al montar. Se
//     conservan remote_endpoint/display_field/value_field (claves legacy con
//     fallback en el resolver KMP, FormFieldsResolver). Cambia el slot_data del
//     seed L4 → bump obligatorio para invalidar la caché SDUI por contenido. Sin
//     cambios de esquema ni de permisos.
const L4_SEED_VERSION = "1.52.0"

// L4_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L4_LAYER_NAME = "L4-full"
