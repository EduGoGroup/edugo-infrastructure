//go:build integration
// +build integration

package l4_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l4"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// TestL4_RolesAlias_Seeded verifica que tras aplicar la capa L4 (que
// internamente invoca ApplyRolesPermissions) existen los 10 roles
// esperados — 4 canónicos + 6 alias — y que la herencia de permisos
// efectivos se cumple:
//
//  1. Los 4 canónicos (student, teacher, guardian, school_admin) están
//     en iam.roles. (PRE-4: platform_admin fue eliminado.)
//  2. Los 6 alias (school_director, school_coordinator, school_assistant,
//     assistant_teacher, observer, readonly_auditor) están en iam.roles.
//  3. Cada alias que hereda (parent_role_id) tiene AL MENOS los mismos
//     permisos efectivos que su rol canónico.
//  4. readonly_auditor NO hereda de nadie (standalone; ver la nota en
//     l4RoleSpecs y TestRoleInheritance_ReadonlyAuditorStandalone): su
//     set efectivo es no vacío y no contiene ninguna permission cuya
//     `Action` sea de mutación.
//
// P4-1 (plan B, SchemaVersion 3.21.0): iam.role_permissions (asignación
// 1:1 rol×permiso) ya no existe. El permiso efectivo de un rol se
// resuelve con los patterns de iam.role_grants (allow + deny, deny-wins)
// del rol y de su cadena parent_role_id (ADR-6). Este test lo computa en
// SQL con iam.permission_matches (ver effectivePermissionsSQL), que es
// el espejo en BD del matcher del backend.
//
// Ejecución:
//
//	ENABLE_INTEGRATION_TESTS=true \
//	    go test -tags=integration -run TestL4_RolesAlias -count=1 \
//	        ./postgres/seeds/system/l4/...
//
// Requiere docker corriendo (testcontainers).
func TestL4_RolesAlias_Seeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL4AliasTest(t)

	canonicalRoleIDs := []string{
		l4.L4_ROLE_STUDENT_ID,
		l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_GUARDIAN_ID,
		// PRE-4: l4.L4_ROLE_ADMIN_ID removido (platform_admin eliminado).
		l4.L4_ROLE_SCHOOL_ADMIN_ID,
	}
	aliasRoleIDs := []string{
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
		l4.L4_ROLE_ASSISTANT_TEACHER_ID,
		l4.L4_ROLE_OBSERVER_ID,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
	}

	// (1) y (2): los 10 roles existen.
	for _, rid := range append(append([]string(nil), canonicalRoleIDs...), aliasRoleIDs...) {
		var cnt int64
		if err := gdb.Raw(`SELECT COUNT(*) FROM iam.roles WHERE id = ?::uuid`, rid).Scan(&cnt).Error; err != nil {
			t.Fatalf("count role %s: %v", rid, err)
		}
		if cnt != 1 {
			t.Errorf("role %s missing in iam.roles (got %d, want 1)", rid, cnt)
		}
	}

	// (3): herencia de permission_ids — aliases que heredan de school_admin.
	aliasParent := map[string]string{
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID:    l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID: l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID:   l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_ASSISTANT_TEACHER_ID:  l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_OBSERVER_ID:           l4.L4_ROLE_TEACHER_ID,
	}

	for alias, canonical := range aliasParent {
		canonicalPerms := loadRolePermissionIDs(t, gdb, canonical)
		aliasPerms := loadRolePermissionIDs(t, gdb, alias)
		if len(aliasPerms) < len(canonicalPerms) {
			t.Errorf(
				"alias %s has fewer permissions (%d) than canonical %s (%d)",
				alias, len(aliasPerms), canonical, len(canonicalPerms),
			)
		}
		for pid := range canonicalPerms {
			if _, ok := aliasPerms[pid]; !ok {
				t.Errorf(
					"alias %s missing permission %s inherited from canonical %s",
					alias, pid, canonical,
				)
			}
		}
	}

	// (4): readonly_auditor — set propio, no vacío, SIN acciones de mutación.
	//
	// El check "subset de teacher" se retiró: el seed lo declara standalone
	// a propósito. Su allow read-only (academic.*, content.*,
	// reports.stats.school, admin.users.read, admin.system_settings.read)
	// cubre permisos que teacher no tiene, y desde el plan 052 F4 teacher
	// tampoco tiene admin.system_settings.*. La invariante que el seed sí
	// garantiza es la de solo lectura.
	readonlyPerms := loadRolePermissionIDs(t, gdb, l4.L4_ROLE_READONLY_AUDITOR_ID)

	if len(readonlyPerms) == 0 {
		t.Errorf("readonly_auditor has no effective permissions (expected non-empty read-only set)")
	}

	var readonlyParents int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.roles WHERE id = ?::uuid AND parent_role_id IS NOT NULL`,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
	).Scan(&readonlyParents).Error; err != nil {
		t.Fatalf("query readonly_auditor parent_role_id: %v", err)
	}
	if readonlyParents != 0 {
		t.Errorf("readonly_auditor must be standalone (parent_role_id NULL)")
	}

	// Sin acciones de mutación: ninguna permission asignada a
	// readonly_auditor debe tener Action en {create, update, delete}
	// ni sufijo :create/:update/:delete.
	type perm struct {
		ID     string `gorm:"column:id"`
		Name   string `gorm:"column:name"`
		Action string `gorm:"column:action"`
	}
	var rows []perm
	if err := gdb.Raw(
		`SELECT p.id::text AS id, p.name, p.action FROM (`+effectivePermissionsSQL+`) p`,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
	).Scan(&rows).Error; err != nil {
		t.Fatalf("load readonly_auditor perms with actions: %v", err)
	}
	for _, r := range rows {
		if isMutation(r.Action) {
			t.Errorf(
				"readonly_auditor has mutation permission: id=%s name=%s action=%s",
				r.ID, r.Name, r.Action,
			)
		}
	}
}

// TestL4_SystemSettings_ReadUpdate_Seeded verifica el bloque R-F7-1:
// los dos permisos system_settings:read y system_settings:update están
// sembrados sobre el recurso system_settings (b4...90) y sus grants
// llegan a los roles esperados (7 para :read, 5 para :update).
//
// Los grants se evalúan como permiso EFECTIVO sobre iam.role_grants
// (patterns + herencia por parent_role_id), no como filas de la tabla
// iam.role_permissions, eliminada en P4-1.
//
// Permisos esperados:
//
//	d1000000-0000-0000-0000-000000000002 → system_settings:read
//	d1000000-0000-0000-0000-000000000003 → system_settings:update
//
// Roles esperados :read (7):
//
//	super_admin (L0, pattern `*`), school_admin (`admin.*`),
//	school_director, school_coordinator, school_assistant (heredan de
//	school_admin), readonly_auditor y announcement_viewer (L1) (literal
//	`admin.system_settings.read`).
//
// teacher, student y guardian ya NO lo tienen: el plan 027 podó los
// grants de consumo y el plan 052 F4 eliminó `admin.system_settings.*`
// de teacher y guardian (destapaba el ítem de menú Configuración). Con
// teacher lo pierden también sus alias assistant_teacher y observer.
//
// Roles esperados :update (5):
//
//	super_admin, school_admin, school_director, school_coordinator,
//	school_assistant.
//
// El test también valida que :read NO llega a teacher/student/guardian/
// assistant_teacher/observer, y que :update NO llega a teacher/student/
// guardian/announcement_viewer/assistant_teacher/observer/
// readonly_auditor.
func TestL4_SystemSettings_ReadUpdate_Seeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	gdb := startPostgresForL4AliasTest(t)

	const (
		permReadID       = "d1000000-0000-0000-0000-000000000002"
		permUpdateID     = "d1000000-0000-0000-0000-000000000003"
		resourceSettings = "b4000000-0000-0000-0000-000000000090"
		roleSuperAdmin   = "10000000-0000-0000-0000-000000000001"
		roleAnnViewer    = "b1000000-0000-0000-0000-000000000001"
	)

	// (1) Ambos permisos existen con el resource_id correcto.
	type permRow struct {
		ID         string `gorm:"column:id"`
		Name       string `gorm:"column:name"`
		Action     string `gorm:"column:action"`
		ResourceID string `gorm:"column:resource_id"`
	}
	for _, expected := range []permRow{
		{ID: permReadID, Name: "admin.system_settings.read", Action: "read", ResourceID: resourceSettings},
		{ID: permUpdateID, Name: "admin.system_settings.update", Action: "update", ResourceID: resourceSettings},
	} {
		var got permRow
		if err := gdb.Raw(
			`SELECT id::text AS id, name, action, resource_id::text AS resource_id `+
				`FROM iam.permissions WHERE id = ?::uuid`,
			expected.ID,
		).Scan(&got).Error; err != nil {
			t.Fatalf("query permission %s: %v", expected.ID, err)
		}
		if got.ID == "" {
			t.Fatalf("permission %s (%s) not seeded", expected.ID, expected.Name)
		}
		if got.Name != expected.Name {
			t.Errorf("permission %s name: got %q, want %q", expected.ID, got.Name, expected.Name)
		}
		if got.Action != expected.Action {
			t.Errorf("permission %s action: got %q, want %q", expected.ID, got.Action, expected.Action)
		}
		if got.ResourceID != expected.ResourceID {
			t.Errorf("permission %s resource_id: got %s, want %s", expected.ID, got.ResourceID, expected.ResourceID)
		}
	}

	// (2) Accessor Permissions() incluye ambos.
	perms, err := l4.Permissions()
	if err != nil {
		t.Fatalf("l4.Permissions(): %v", err)
	}
	foundRead, foundUpdate := false, false
	for _, p := range perms {
		switch p.ID.String() {
		case permReadID:
			foundRead = true
		case permUpdateID:
			foundUpdate = true
		}
	}
	if !foundRead {
		t.Errorf("l4.Permissions() missing system_settings:read (%s)", permReadID)
	}
	if !foundUpdate {
		t.Errorf("l4.Permissions() missing system_settings:update (%s)", permUpdateID)
	}

	// (3) Grants de :read llegan a los 7 roles esperados.
	// PRE-4: platform_admin (L4_ROLE_ADMIN_ID) removido del catálogo.
	// Plan 027 / plan 052 F4: teacher, student y guardian ya no lo tienen.
	readRoles := []string{
		roleSuperAdmin,
		l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
		roleAnnViewer,
	}
	for _, rid := range readRoles {
		if !roleHasPermission(t, gdb, rid, permReadID) {
			t.Errorf("role %s missing system_settings:read grant", rid)
		}
	}
	forbiddenRead := []string{
		l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_STUDENT_ID,
		l4.L4_ROLE_GUARDIAN_ID,
		l4.L4_ROLE_ASSISTANT_TEACHER_ID,
		l4.L4_ROLE_OBSERVER_ID,
	}
	for _, rid := range forbiddenRead {
		if roleHasPermission(t, gdb, rid, permReadID) {
			t.Errorf("role %s unexpectedly has system_settings:read grant", rid)
		}
	}
	if totalRead := countRolesWithPermission(t, gdb, permReadID); totalRead != int64(len(readRoles)) {
		t.Errorf("total :read grants: got %d, want %d", totalRead, len(readRoles))
	}

	// (4) Grants de :update llegan a los 5 roles esperados.
	// PRE-4: platform_admin (L4_ROLE_ADMIN_ID) removido del catálogo.
	updateRoles := []string{
		roleSuperAdmin,
		l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
	}
	for _, rid := range updateRoles {
		if !roleHasPermission(t, gdb, rid, permUpdateID) {
			t.Errorf("role %s missing system_settings:update grant", rid)
		}
	}

	// (5) :update NO llega a teacher/student/guardian/announcement_viewer
	// ni a los 3 alias de teacher (assistant_teacher, observer, readonly_auditor).
	forbiddenUpdate := []string{
		l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_STUDENT_ID,
		l4.L4_ROLE_GUARDIAN_ID,
		roleAnnViewer,
		l4.L4_ROLE_ASSISTANT_TEACHER_ID,
		l4.L4_ROLE_OBSERVER_ID,
		l4.L4_ROLE_READONLY_AUDITOR_ID,
	}
	for _, rid := range forbiddenUpdate {
		if roleHasPermission(t, gdb, rid, permUpdateID) {
			t.Errorf("role %s unexpectedly has system_settings:update grant", rid)
		}
	}
	if totalUpdate := countRolesWithPermission(t, gdb, permUpdateID); totalUpdate != int64(len(updateRoles)) {
		t.Errorf("total :update grants: got %d, want %d", totalUpdate, len(updateRoles))
	}
}

// effectivePermissionsSQL selecciona las filas de iam.permissions que un
// rol (único parámetro: role_id) tiene como permiso EFECTIVO en el modelo
// de iam.role_grants:
//
//   - `chain` recorre el rol y sus ancestros por parent_role_id (ADR-6:
//     el alias hereda los grants del canónico; el backend los aplana en
//     el login).
//   - Un permiso es efectivo si algún grant allow de la cadena lo cubre
//     (iam.permission_matches(pattern, permissions.name)) y ningún grant
//     deny de la cadena lo cubre (deny-wins, default deny).
//
// Reemplaza a la lectura directa de iam.role_permissions, tabla
// eliminada en P4-1 (plan B).
const effectivePermissionsSQL = `
	WITH RECURSIVE chain AS (
		SELECT id, parent_role_id FROM iam.roles WHERE id = ?::uuid
		UNION ALL
		SELECT r.id, r.parent_role_id
		FROM iam.roles r
		JOIN chain c ON r.id = c.parent_role_id
	)
	SELECT p.id, p.name, p.action
	FROM iam.permissions p
	WHERE EXISTS (
		SELECT 1 FROM iam.role_grants g
		JOIN chain c ON c.id = g.role_id
		WHERE g.effect = 'allow' AND iam.permission_matches(g.pattern, p.name)
	)
	AND NOT EXISTS (
		SELECT 1 FROM iam.role_grants g
		JOIN chain c ON c.id = g.role_id
		WHERE g.effect = 'deny' AND iam.permission_matches(g.pattern, p.name)
	)`

// roleHasPermission indica si permissionID está en el set efectivo del rol.
func roleHasPermission(t *testing.T, gdb *gorm.DB, roleID, permissionID string) bool {
	t.Helper()
	_, ok := loadRolePermissionIDs(t, gdb, roleID)[permissionID]
	return ok
}

// countRolesWithPermission cuenta los roles de iam.roles cuyo set
// efectivo contiene permissionID.
func countRolesWithPermission(t *testing.T, gdb *gorm.DB, permissionID string) int64 {
	t.Helper()
	var roleIDs []string
	if err := gdb.Raw(`SELECT id::text FROM iam.roles`).Scan(&roleIDs).Error; err != nil {
		t.Fatalf("load iam.roles: %v", err)
	}
	var total int64
	for _, rid := range roleIDs {
		if roleHasPermission(t, gdb, rid, permissionID) {
			total++
		}
	}
	return total
}

// loadRolePermissionIDs retorna el set de permission_id (como string)
// efectivos de un role_id (ver effectivePermissionsSQL).
func loadRolePermissionIDs(t *testing.T, gdb *gorm.DB, roleID string) map[string]struct{} {
	t.Helper()
	var ids []string
	if err := gdb.Raw(
		`SELECT p.id::text FROM (`+effectivePermissionsSQL+`) p`,
		roleID,
	).Scan(&ids).Error; err != nil {
		t.Fatalf("load effective permissions for %s: %v", roleID, err)
	}
	out := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		out[id] = struct{}{}
	}
	return out
}

// isMutation refleja la heurística usada por filterMutationGrants en
// roles_permissions.go. Tras SUB-2 (Pass 1) cubre el conjunto completo
// de verbos mutativos del catálogo: CRUD + workflow/lifecycle
// (publish/finalize/activate/approve/grade/attempt/assign/review/
// manage/request) y sus variantes con sufijo `:own`.
func isMutation(action string) bool {
	mutationVerbs := map[string]struct{}{
		"create": {}, "update": {}, "delete": {},
		"publish": {}, "finalize": {}, "activate": {}, "approve": {},
		"grade": {}, "attempt": {}, "assign": {}, "review": {},
		"manage": {}, "request": {},
	}
	for _, tok := range strings.Split(action, ":") {
		if _, ok := mutationVerbs[tok]; ok {
			return true
		}
	}
	return false
}

// startPostgresForL4AliasTest levanta postgres:15-alpine y ejecuta
// migrations.Migrate(Force=true, SeedUpToLayer=L4_LAYER_NAME) — siembra
// L0+L1+L2+L3+L4 una vez. Replica el patrón de startPostgresForL3Test.
func startPostgresForL4AliasTest(tb testing.TB) *gorm.DB {
	tb.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		tb.Fatalf("l4 alias integration: testcontainers GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("l4 alias integration: container.Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("l4 alias integration: container.MappedPort: %v", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Fatalf("l4 alias integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("l4 alias integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:         true,
		DBUser:        "test",
		SeedUpToLayer: layers.L4_LAYER_NAME,
	}); err != nil {
		tb.Fatalf("l4 alias integration: migrations.Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("l4 alias integration: gorm.Open: %v", err)
	}
	return gdb
}
