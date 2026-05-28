package l4

import (
	"sort"
	"strings"
	"testing"
)

// Esta suite valida la herencia de roles (ADR-6 / F1) a nivel de seed,
// sin BD: que tras quitar los grants propios de los alias y resolver la
// herencia por parent_role_id, los grants EFECTIVOS aplanados de cada
// alias quedan IDÉNTICOS a los de su rol canónico — que es justamente su
// baseline, porque los alias eran copias literales del canónico.
//
// Reusa la MISMA semántica del matcher (deny-wins, glob) que el runtime,
// replicada localmente en `matchesPattern`/`evaluate` para no introducir
// una dependencia del módulo postgres hacia edugo-shared/auth. El matcher
// autoritativo (edugo-shared/auth) tiene sus propios golden tests.

// matchesPattern es el espejo 1:1 de auth.PermissionMatches.
func matchesPattern(pattern, request string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == request {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := pattern[:len(pattern)-2]
		return request == prefix || strings.HasPrefix(request, prefix+".")
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:]
		return len(request) > len(suffix) && strings.HasSuffix(request, suffix)
	}
	if i := strings.Index(pattern, ".*."); i > 0 {
		head := pattern[:i+1]
		tail := pattern[i+2:]
		if strings.Contains(head, "*") || strings.Contains(tail, "*") {
			return false
		}
		if !strings.HasPrefix(request, head) || !strings.HasSuffix(request, tail) {
			return false
		}
		if len(request) <= len(head)+len(tail) {
			return false
		}
		middle := request[len(head) : len(request)-len(tail)]
		return !strings.HasPrefix(middle, ".") && !strings.HasSuffix(middle, ".")
	}
	return false
}

// evaluate es el espejo de auth.EvaluateGrants: deny-wins, default deny.
func evaluate(allow, deny []string, request string) bool {
	for _, d := range deny {
		if matchesPattern(d, request) {
			return false
		}
	}
	for _, a := range allow {
		if matchesPattern(a, request) {
			return true
		}
	}
	return false
}

// flattenRoleGrants resuelve la cadena de herencia de un rol del seed
// (vía parentByRole) y devuelve los grants efectivos aplanados (union de
// allow y deny a lo largo de la cadena). Espeja la resolución del login.
func flattenRoleGrants(t *testing.T, roleID string) (allow, deny []string) {
	t.Helper()
	allowMap := roleGrantPatterns()
	denyMap := roleGrantDenyPatterns()
	parents := parentByRole(t)

	seen := map[string]struct{}{}
	current := roleID
	for current != "" {
		if _, dup := seen[current]; dup {
			t.Fatalf("ciclo detectado en la cadena de herencia en %s", current)
		}
		seen[current] = struct{}{}
		allow = append(allow, allowMap[current]...)
		deny = append(deny, denyMap[current]...)
		current = parents[current]
	}
	return dedup(allow), dedup(deny)
}

// parentByRole materializa el mapa role_id → parent_role_id desde las
// specs declarativas del seed (misma fuente que buildL4Roles).
func parentByRole(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, s := range l4RoleSpecs() {
		if s.parentIDStr != "" {
			out[s.idStr] = s.parentIDStr
		}
	}
	return out
}

func dedup(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// requestCorpus reúne todos los permisos del catálogo (path-based) que
// sirven de requests para comparar decisiones de autorización. Incluye
// los nombres L4 (l4Permissions) más los nombres conocidos de L0
// (announcements) y L3 (materials) que viven fuera de este paquete.
func requestCorpus(t *testing.T) []string {
	t.Helper()
	set := map[string]struct{}{}
	for _, p := range l4Permissions() {
		set[p.name] = struct{}{}
	}
	for _, extra := range []string{
		"academic.announcements.read",
		"academic.announcements.create",
		"academic.announcements.update",
		"academic.announcements.delete",
		"content.materials.read",
		"content.materials.create",
		"content.materials.update",
	} {
		set[extra] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for r := range set {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

// TestRoleInheritance_AliasEffectiveEqualsCanonical valida la NO
// regresión: para cada alias que hereda, la decisión de autorización
// (allow/deny) sobre TODO el corpus de permisos debe coincidir con la de
// su rol canónico — el baseline antes del dedup, porque los alias eran
// copias literales del canónico.
func TestRoleInheritance_AliasEffectiveEqualsCanonical(t *testing.T) {
	aliasParent := map[string]string{
		L4_ROLE_SCHOOL_DIRECTOR_ID:    L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_SCHOOL_COORDINATOR_ID: L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_SCHOOL_ASSISTANT_ID:   L4_ROLE_SCHOOL_ADMIN_ID,
		L4_ROLE_ASSISTANT_TEACHER_ID:  L4_ROLE_TEACHER_ID,
		L4_ROLE_OBSERVER_ID:           L4_ROLE_TEACHER_ID,
	}
	corpus := requestCorpus(t)

	for alias, canonical := range aliasParent {
		aAllow, aDeny := flattenRoleGrants(t, alias)
		cAllow, cDeny := flattenRoleGrants(t, canonical)
		for _, req := range corpus {
			gotAlias := evaluate(aAllow, aDeny, req)
			gotCanon := evaluate(cAllow, cDeny, req)
			if gotAlias != gotCanon {
				t.Errorf("alias %s difiere del canónico %s en %q: alias=%v canónico=%v",
					alias, canonical, req, gotAlias, gotCanon)
			}
		}
	}
}

// TestRoleInheritance_AliasesHaveNoOwnGrants confirma que los alias que
// heredan ya NO declaran patterns propios en el seed (toda su
// autorización proviene del canónico vía parent_role_id).
func TestRoleInheritance_AliasesHaveNoOwnGrants(t *testing.T) {
	allowMap := roleGrantPatterns()
	denyMap := roleGrantDenyPatterns()
	inheriting := []string{
		L4_ROLE_SCHOOL_DIRECTOR_ID,
		L4_ROLE_SCHOOL_COORDINATOR_ID,
		L4_ROLE_SCHOOL_ASSISTANT_ID,
		L4_ROLE_ASSISTANT_TEACHER_ID,
		L4_ROLE_OBSERVER_ID,
	}
	for _, rid := range inheriting {
		if len(allowMap[rid]) != 0 {
			t.Errorf("alias %s aún declara %d allow propios (debería heredar)", rid, len(allowMap[rid]))
		}
		if len(denyMap[rid]) != 0 {
			t.Errorf("alias %s aún declara %d deny propios (debería heredar)", rid, len(denyMap[rid]))
		}
	}
}

// TestRoleInheritance_ReadonlyAuditorStandalone documenta que
// readonly_auditor NO hereda (no es superset exacto de teacher) y por
// tanto conserva sus grants propios. Su set efectivo debe seguir
// excluyendo toda mutación.
func TestRoleInheritance_ReadonlyAuditorStandalone(t *testing.T) {
	parents := parentByRole(t)
	if p, ok := parents[L4_ROLE_READONLY_AUDITOR_ID]; ok {
		t.Fatalf("readonly_auditor no debería tener parent, tiene %s", p)
	}
	allow := roleGrantPatterns()[L4_ROLE_READONLY_AUDITOR_ID]
	deny := roleGrantDenyPatterns()[L4_ROLE_READONLY_AUDITOR_ID]
	if len(allow) == 0 || len(deny) == 0 {
		t.Fatalf("readonly_auditor debe conservar allow y deny propios (allow=%d deny=%d)", len(allow), len(deny))
	}
	// Ninguna mutación debe quedar permitida.
	for _, req := range requestCorpus(t) {
		if !evaluate(allow, deny, req) {
			continue
		}
		// permitida → debe ser de lectura/consulta (no mutación).
		for _, verb := range []string{".create", ".update", ".delete", ".publish",
			".finalize", ".activate", ".approve", ".grade", ".attempt",
			".assign", ".review", ".manage", ".request"} {
			if strings.HasSuffix(req, verb) {
				t.Errorf("readonly_auditor permite mutación %q", req)
			}
		}
	}
}
