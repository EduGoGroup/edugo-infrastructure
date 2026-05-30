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
const L4_SEED_VERSION = "1.36.0"

// L4_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L4_LAYER_NAME = "L4-full"
