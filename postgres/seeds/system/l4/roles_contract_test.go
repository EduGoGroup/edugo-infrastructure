package l4

import "testing"

// Suite de CONTRATO de los grants por rol (plan 052, Frente 4).
//
// Fija lo que el Frente 4 saneó a partir del recorrido de QA del rediseño
// de UI: qué NO debe alcanzar cada rol (hallazgos QA-08, QA-09, QA-12,
// QA-25) y —tan importante como eso— qué SÍ debe seguir alcanzando, para
// que una poda futura no se pase de frenada. Si alguien revierte una de
// estas decisiones, aquí se entera.
//
// Evalúa el permiso EFECTIVO con la semántica del matcher (deny-wins,
// glob) tras resolver la herencia de roles; NO compara literales contra la
// lista de patterns. Esa es la gracia: el test también caza que alguien
// reintroduzca un COMODÍN que vuelva a cubrir lo prohibido, que es
// exactamente como aparecieron QA-08 (`reports.*` en teacher) y QA-12
// (`reports.read` en guardian).
//
// Matcher: se reusan `evaluate`/`matchesPattern` de
// roles_inheritance_test.go, espejo 1:1 de auth.EvaluateGrants /
// auth.PermissionMatches. NO se importa `edugo-shared/auth` porque el
// módulo `postgres` no depende de edugo-shared: añadirlo metería una
// dependencia de PRODUCCIÓN en su go.mod (y arrastraría jwt/common) sólo
// para un test. El matcher autoritativo tiene sus propios golden tests.

// roleCan responde si el rol alcanza el permiso. Resuelve antes la cadena
// de herencia (parent_role_id), imprescindible para los alias: por ADR-6
// assistant_teacher/observer (→ teacher) y school_director/coordinator/
// assistant (→ school_admin) NO declaran grants propios, así que
// consultarlos en roleGrantPatterns() directamente daría vacío.
func roleCan(t *testing.T, roleID, permission string) bool {
	t.Helper()
	allow, deny := flattenRoleGrants(t, roleID)
	return evaluate(allow, deny, permission)
}

// assertDenied exige que NINGUNO de los permisos sea alcanzable por el rol.
func assertDenied(t *testing.T, roleID, roleName string, permissions ...string) {
	t.Helper()
	for _, p := range permissions {
		if roleCan(t, roleID, p) {
			t.Errorf("%s NO debería alcanzar %q, pero sus grants efectivos lo permiten", roleName, p)
		}
	}
}

// assertAllowed exige que TODOS los permisos sigan siendo alcanzables por
// el rol (guarda contra podas que se pasan de frenada).
func assertAllowed(t *testing.T, roleID, roleName string, permissions ...string) {
	t.Helper()
	for _, p := range permissions {
		if !roleCan(t, roleID, p) {
			t.Errorf("%s SÍ debería alcanzar %q, pero sus grants efectivos lo niegan", roleName, p)
		}
	}
}

// TestContratoGrants_ProfesorNoVeLasEstadisticasDeToooaLaPlataforma fija
// QA-08: el profesor tenía `reports.*`, que cubría `reports.stats.global`
// (GET /stats/global) — los totales agregados de TODOS los colegios. Un
// docente veía las cifras de la plataforma entera. El Frente 3 creó
// GET /stats/school como alternativa acotada; el profesor no recibe
// tampoco ese literal porque QA-08 pide que no vea el panel.
func TestContratoGrants_ProfesorNoVeLasEstadisticasDeTodaLaPlataforma(t *testing.T) {
	assertDenied(t, L4_ROLE_TEACHER_ID, "teacher", "reports.stats.global")

	// Los alias que heredan de teacher arrastrarían la misma fuga.
	assertDenied(t, L4_ROLE_ASSISTANT_TEACHER_ID, "assistant_teacher", "reports.stats.global")
	assertDenied(t, L4_ROLE_OBSERVER_ID, "observer", "reports.stats.global")
}

// TestContratoGrants_ProfesorNoAdministraElDirectorioDeUsuarios fija
// QA-09: en el menú del profesor aparecía «Administración > Usuarios» con
// el botón de crear. El rol NUNCA tuvo `admin.users.create` — venía de un
// user_grant del fixture— pero el contrato debe dejarlo escrito para que
// nadie lo "arregle" concediéndoselo al rol.
//
// Se fijan sólo los verbos MUTATIVOS del directorio: dar de alta, editar o
// borrar usuarios del colegio es acto de school_admin. Las lecturas se
// dejan fuera a propósito (un perfil propio `admin.users.read:own` sería
// una concesión legítima a futuro y no debe romper este test).
func TestContratoGrants_ProfesorNoAdministraElDirectorioDeUsuarios(t *testing.T) {
	assertDenied(t, L4_ROLE_TEACHER_ID, "teacher",
		"admin.users.create",
		"admin.users.update",
		"admin.users.delete",
	)

	assertDenied(t, L4_ROLE_ASSISTANT_TEACHER_ID, "assistant_teacher", "admin.users.create")
	assertDenied(t, L4_ROLE_OBSERVER_ID, "observer", "admin.users.create")
}

// TestContratoGrants_ProfesorNoOperaLaConfiguracionDelSistema fija el otro
// hijo del nodo «Administración» señalado por QA-09: el profesor tenía
// `admin.system_settings.*`, que le destapaba el ítem
// «Administración > Configuración». Ninguna ruta de las 4 APIs se lo exige,
// así que quitarlo no rompe llamadas suyas; dejarlo habría vuelto cosmético
// el fix de QA-09, porque el nodo padre «Administración» seguiría pintándose.
func TestContratoGrants_ProfesorNoOperaLaConfiguracionDelSistema(t *testing.T) {
	assertDenied(t, L4_ROLE_TEACHER_ID, "teacher",
		"admin.system_settings.read",
		"admin.system_settings.update",
		"admin.system_settings.settings",
	)

	assertDenied(t, L4_ROLE_ASSISTANT_TEACHER_ID, "assistant_teacher", "admin.system_settings.read")
	assertDenied(t, L4_ROLE_OBSERVER_ID, "observer", "admin.system_settings.read")
}

// TestContratoGrants_ProfesorConservaTodoLoSuyo es la contracara de los
// tres casos anteriores y pesa lo mismo: las podas de QA-08/QA-09 recortan
// el borde administrativo del rol, NO su trabajo. Si una poda futura se
// pasa de frenada y el docente deja de poder pasar lista, poner notas,
// publicar material o escribirle a las familias, se cae aquí.
func TestContratoGrants_ProfesorConservaTodoLoSuyo(t *testing.T) {
	assertAllowed(t, L4_ROLE_TEACHER_ID, "teacher",
		// Comunicaciones a su clase.
		"academic.announcements.read",
		"academic.announcements.create",
		"academic.announcements.update",
		"academic.announcements.delete",
		// Pasar lista.
		"academic.attendance.read",
		"academic.attendance.create",
		"academic.attendance.update",
		// Calificar y cerrar notas.
		"academic.grades.read",
		"academic.grades.create",
		"academic.grades.update",
		"academic.grades.finalize",
		// Roster de su unidad: literal `.read`, no el wildcard (no crea ni
		// borra membresías).
		"academic.memberships.read",
		// «Materias que dicto» — su única vista de qué imparte tras la poda
		// de subjects/subject_offerings (plan 027 F3).
		"academic.my_teaching.read:own",
		// Biblioteca de materiales.
		"content.materials.read",
		"content.materials.create",
		"content.materials.update",
		"content.materials.delete",
		"content.materials.download",
		"content.materials.publish",
		// Evaluaciones: del diseño a la revisión.
		"content.assessments.read",
		"content.assessments.create",
		"content.assessments.update",
		"content.assessments.delete",
		"content.assessments.publish",
		"content.assessments.assign",
		"content.assessments.grade",
		"content.assessments.review",
		// WhatsApp a las familias de su clase (plan 025).
		"messaging.session.pair",
		"messaging.message.send",
		"messaging.device.link",
		// Chasis de la app: sin esto no hay menú ni pantallas.
		"dashboard.view",
		"menu.read",
		"screens.read",
		"notifications.read",
	)
}

// TestContratoGrants_RepresentanteNoTieneNodoDeReportes fija QA-12: el
// guardián tenía `reports.read`. Ninguna ruta de ninguna API lo exige, pero
// SÍ «tocaba» el recurso de menú `reports` (el gate es por prefijo de path),
// así que le pintaba un ítem raíz «Reportes» cuyo único hijo —`stats`— ese
// mismo permiso no habilitaba: un nodo que no navegaba a ningún sitio.
func TestContratoGrants_RepresentanteNoTieneNodoDeReportes(t *testing.T) {
	assertDenied(t, L4_ROLE_GUARDIAN_ID, "guardian",
		"reports.read",
		"reports.stats.school",
		// Por si acaso: los totales de la plataforma tampoco, obviamente.
		"reports.stats.global",
	)
}

// TestContratoGrants_AuditorAuditaSuColegioNoLaPlataforma fija la decisión
// del dueño (2026-08-01): readonly_auditor tenía `reports.*`, que incluía
// los totales de TODOS los colegios, mientras su contexto tiene school_id
// fijado — la amplitud contradecía su propio alcance. Queda el literal de
// colegio, que además MANTIENE visible el ítem «Estadísticas» (el
// resourcePath del recurso es `reports.stats` y el literal lo toca por
// prefijo).
func TestContratoGrants_AuditorAuditaSuColegioNoLaPlataforma(t *testing.T) {
	assertDenied(t, L4_ROLE_READONLY_AUDITOR_ID, "readonly_auditor", "reports.stats.global")
	assertAllowed(t, L4_ROLE_READONLY_AUDITOR_ID, "readonly_auditor", "reports.stats.school")
}

// TestContratoGrants_AuditorNoVeLasPantallasPersonalesDeOtros fija QA-25:
// el allow `academic.*` del auditor le arrastraba los cuatro recursos
// "self" (my_teaching / my_memberships / my_grades / my_attendance), que
// para él devuelven listas vacías porque no es ni profesor ni alumno —y
// además le dejaban DOS ítems «Mis Materias» idénticos en el menú—. Los
// tapa el deny `academic.*.read:own`, el mismo que ya llevaba school_admin
// desde el plan 027 F4.8.
func TestContratoGrants_AuditorNoVeLasPantallasPersonalesDeOtros(t *testing.T) {
	assertDenied(t, L4_ROLE_READONLY_AUDITOR_ID, "readonly_auditor",
		"academic.my_teaching.read:own",
		"academic.my_memberships.read:own",
		"academic.my_grades.read:own",
		"academic.my_attendance.read:own",
	)
}

// TestContratoGrants_AuditorSigueSiendoDeSoloLectura protege lo que
// garantizan sus deny `*.<verbo>`: por amplio que sea su allow
// (`academic.*`, `content.*`), ninguna mutación debe quedar alcanzable.
// Nota sobre `academic.subjects.manage`: hoy NO existe en el catálogo — se
// evalúa a propósito para demostrar que un permiso mutativo NUEVO cae bajo
// el deny sin tocar el seed, que es justo el diseño del rol.
func TestContratoGrants_AuditorSigueSiendoDeSoloLectura(t *testing.T) {
	assertDenied(t, L4_ROLE_READONLY_AUDITOR_ID, "readonly_auditor",
		"academic.units.create",
		"academic.grades.update",
		"content.materials.delete",
		"academic.subjects.manage",
	)
}

// TestContratoRoles_CadaRolAterrizaEnElDashboardDeSuArquetipo fija el
// landing_screen_key de LOS 10 ROLES de L4 (4 canónicos + 6 alias).
//
// Por qué merece un golden propio: el landing NO se hereda. La cascada del
// backend es (landing del rol ?? default de la escuela ?? "dashboard-home") y
// mira SOLO el campo propio del rol —no resuelve parent_role_id (ADR 0024,
// nota en l4RoleSpecs)—, así que un alias con el campo vacío no aterriza donde
// su canónico: cae al home genérico. Es un fallo silencioso: nadie ve un 403,
// el usuario simplemente entra a otra pantalla. Los grants no lo detectan
// (assertAllowed/assertDenied no miran esta columna) y el Frente 4 acaba de
// mover el del auditor de `dashboard-teacher` a `dashboard-schooladmin`
// (QA-11) sin nada que lo sujete.
//
// Además de la tabla, el test exige dos invariantes:
//   - ningún rol se queda sin landing (vacío → NULL → home genérico);
//   - el landing apunta a una screen_instance que EXISTE en el seed (un typo
//     o el borrado de una pantalla dejaría al rol aterrizando en el vacío).
//
// Alcance: los roles de L4. `super_admin` (L0) y `announcement_viewer` (L1)
// declaran el suyo en el paquete `layers` y quedan fuera de este paquete.
func TestContratoRoles_CadaRolAterrizaEnElDashboardDeSuArquetipo(t *testing.T) {
	landingEsperado := map[string]string{
		// Canónicos: cada arquetipo a su panel.
		L4_ROLE_STUDENT_ID:      "dashboard-student",
		L4_ROLE_TEACHER_ID:      "dashboard-teacher",
		L4_ROLE_GUARDIAN_ID:     "dashboard-guardian",
		L4_ROLE_SCHOOL_ADMIN_ID: "dashboard-schooladmin",
		// Alias de school_admin: lo reciben EXPLÍCITO (no se hereda).
		L4_ROLE_SCHOOL_DIRECTOR_ID:    "dashboard-schooladmin",
		L4_ROLE_SCHOOL_COORDINATOR_ID: "dashboard-schooladmin",
		L4_ROLE_SCHOOL_ASSISTANT_ID:   "dashboard-schooladmin",
		// Alias de teacher: idem.
		L4_ROLE_ASSISTANT_TEACHER_ID: "dashboard-teacher",
		L4_ROLE_OBSERVER_ID:          "dashboard-teacher",
		// readonly_auditor no hereda de nadie y el Frente 4 lo movió aquí
		// (QA-11): en `dashboard-teacher` el panel le pedía «sus» sesiones —
		// GET /me/teaching, GET /me/subject-offerings— que él no tiene, y
		// devolvían 428 apilados. `dashboard-schooladmin` sí encaja: los
		// indicadores del colegio vía GET /stats/school, que le responde 200
		// desde que este mismo frente le dio `reports.stats.school`.
		L4_ROLE_READONLY_AUDITOR_ID: "dashboard-schooladmin",
	}

	specs := l4RoleSpecs()
	if len(specs) != len(landingEsperado) {
		t.Fatalf("el seed declara %d roles y la tabla del contrato fija %d: si agregaste o quitaste un rol, decláralo aquí con su landing",
			len(specs), len(landingEsperado))
	}

	// Las pantallas que L4 siembra, para verificar que el landing existe.
	pantallas := make(map[string]struct{})
	for _, inst := range buildL4ScreenInstances() {
		pantallas[inst.ScreenKey] = struct{}{}
	}

	for _, s := range specs {
		esperado, fijado := landingEsperado[s.idStr]
		if !fijado {
			t.Errorf("el rol %s (%s) no está en la tabla del contrato: agrégalo con el landing que le corresponde", s.name, s.idStr)
			continue
		}
		if s.landingScreenKey == "" {
			t.Errorf("el rol %s se quedó SIN landing: el campo vacío se siembra como NULL y la cascada lo manda al home genérico «dashboard-home», no a %q (el landing NO se hereda del rol padre)",
				s.name, esperado)
			continue
		}
		if s.landingScreenKey != esperado {
			t.Errorf("el rol %s aterriza en %q y el contrato fija %q", s.name, s.landingScreenKey, esperado)
			continue
		}
		if _, existe := pantallas[s.landingScreenKey]; !existe {
			t.Errorf("el rol %s aterriza en %q, que NO es ninguna screen_instance sembrada: el rol quedaría sin pantalla de inicio",
				s.name, s.landingScreenKey)
		}
	}

	// La spec es declarativa; lo que llega a la BD lo escribe buildL4Roles.
	// Si el builder dejara de propagar el campo, la tabla de arriba seguiría
	// verde y el producto estaría roto igual.
	roles, err := buildL4Roles()
	if err != nil {
		t.Fatalf("buildL4Roles: %v", err)
	}
	for _, r := range roles {
		if r.LandingScreenKey == nil {
			t.Errorf("buildL4Roles no propagó el landing del rol %s: llegaría NULL a iam.roles", r.Name)
			continue
		}
		if got, want := *r.LandingScreenKey, landingEsperado[r.ID.String()]; got != want {
			t.Errorf("buildL4Roles materializa el landing del rol %s como %q; el contrato fija %q", r.Name, got, want)
		}
	}
}

// TestContratoRecursos_NoHayDosEtiquetasDeMenuIgualesBajoElMismoPadre fija
// la otra mitad de QA-25: `my_memberships` («Mis Materias» del alumno) y
// `my_teaching` («Mis Materias» del profesor) tenían DisplayName idéntico y
// el mismo icono, así que quien alcanzaba ambos —el auditor, vía
// `academic.*`— veía dos entradas indistinguibles en el menú. La regla se
// escribe genérica: dos recursos de menú hermanos nunca deben compartir
// etiqueta, así el test caza también la próxima colisión.
func TestContratoRecursos_NoHayDosEtiquetasDeMenuIgualesBajoElMismoPadre(t *testing.T) {
	byKey := make(map[string]l4ResourceRow, len(l4Resources))
	for _, r := range l4Resources {
		byKey[r.Key] = r
	}

	// El caso concreto que originó la regla.
	memberships, ok := byKey["my_memberships"]
	if !ok {
		t.Fatal("no se encontró el recurso my_memberships en l4Resources")
	}
	teaching, ok := byKey["my_teaching"]
	if !ok {
		t.Fatal("no se encontró el recurso my_teaching en l4Resources")
	}
	if memberships.DisplayName == teaching.DisplayName {
		t.Errorf("my_memberships y my_teaching comparten la etiqueta %q: el menú vuelve a mostrar dos «Mis Materias» idénticas",
			memberships.DisplayName)
	}

	// La regla general.
	type siblingLabel struct {
		parentID    string
		displayName string
	}
	seen := make(map[siblingLabel]string, len(l4Resources))
	for _, r := range l4Resources {
		if !r.IsMenuVisible || !r.IsActive {
			continue
		}
		key := siblingLabel{parentID: r.ParentID, displayName: r.DisplayName}
		if previous, dup := seen[key]; dup {
			t.Errorf("los recursos de menú %q y %q comparten la etiqueta %q bajo el mismo padre (%s): quien alcance ambos verá dos ítems indistinguibles",
				previous, r.Key, r.DisplayName, parentLabel(r.ParentID, byKey))
			continue
		}
		seen[key] = r.Key
	}
}

// parentLabel traduce un parent_id a algo legible en el mensaje de error
// (la key del padre, o "raíz del menú" si no tiene).
func parentLabel(parentID string, byKey map[string]l4ResourceRow) string {
	if parentID == "" {
		return "raíz del menú"
	}
	for _, r := range byKey {
		if r.ID == parentID {
			return r.Key
		}
	}
	return parentID
}
