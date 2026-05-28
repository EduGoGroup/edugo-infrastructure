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
// esperados — 4 canónicos + 6 alias — y que la herencia de
// role_permissions se cumple:
//
//  1. Los 4 canónicos (student, teacher, guardian, school_admin) están
//     en iam.roles. (PRE-4: platform_admin fue eliminado.)
//  2. Los 6 alias (school_director, school_coordinator, school_assistant,
//     assistant_teacher, observer, readonly_auditor) están en iam.roles.
//  3. Cada alias tiene AL MENOS los mismos permission_ids que su rol
//     canónico, con la excepción de readonly_auditor: éste debe ser un
//     subconjunto del conjunto de teacher SIN ninguna permission cuya
//     `Action` esté en {create, update, delete, *:create, *:update,
//     *:delete}.
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

	// (1) y (2): los 11 roles existen.
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

	// (4): readonly_auditor — subset de teacher SIN acciones de mutación.
	teacherPerms := loadRolePermissionIDs(t, gdb, l4.L4_ROLE_TEACHER_ID)
	readonlyPerms := loadRolePermissionIDs(t, gdb, l4.L4_ROLE_READONLY_AUDITOR_ID)

	if len(readonlyPerms) == 0 {
		t.Errorf("readonly_auditor has no permissions assigned (expected non-empty read-only subset of teacher)")
	}
	// Subset: cada permiso de readonly_auditor está en teacher.
	for pid := range readonlyPerms {
		if _, ok := teacherPerms[pid]; !ok {
			t.Errorf("readonly_auditor has permission %s that is NOT in teacher's set (subset violation)", pid)
		}
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
		`SELECT p.id::text AS id, p.name, p.action `+
			`FROM iam.role_permissions rp `+
			`JOIN iam.permissions p ON rp.permission_id = p.id `+
			`WHERE rp.role_id = ?::uuid`,
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
// llegan a los roles esperados (9 para :read, 5 para :update —
// post-PRE-4, platform_admin removido).
//
// Permisos esperados:
//
//	d1000000-0000-0000-0000-000000000002 → system_settings:read
//	d1000000-0000-0000-0000-000000000003 → system_settings:update
//
// Roles esperados :read (9):
//
//	super_admin (L0), school_admin, school_director, school_coordinator,
//	school_assistant, teacher, student, guardian, announcement_viewer (L1).
//
// Roles esperados :update (5):
//
//	super_admin, school_admin, school_director, school_coordinator,
//	school_assistant.
//
// El test también valida que :update NO llega a teacher/student/
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

	// (3) Grants de :read llegan a los 9 roles esperados.
	// PRE-4: platform_admin (L4_ROLE_ADMIN_ID) removido del catálogo.
	readRoles := []string{
		roleSuperAdmin,
		l4.L4_ROLE_SCHOOL_ADMIN_ID,
		l4.L4_ROLE_SCHOOL_DIRECTOR_ID,
		l4.L4_ROLE_SCHOOL_COORDINATOR_ID,
		l4.L4_ROLE_SCHOOL_ASSISTANT_ID,
		l4.L4_ROLE_TEACHER_ID,
		l4.L4_ROLE_STUDENT_ID,
		l4.L4_ROLE_GUARDIAN_ID,
		roleAnnViewer,
	}
	for _, rid := range readRoles {
		var cnt int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM iam.role_permissions `+
				`WHERE role_id = ?::uuid AND permission_id = ?::uuid`,
			rid, permReadID,
		).Scan(&cnt).Error; err != nil {
			t.Fatalf("count grant role=%s perm=read: %v", rid, err)
		}
		if cnt != 1 {
			t.Errorf("role %s missing system_settings:read grant (got %d, want 1)", rid, cnt)
		}
	}
	var totalRead int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.role_permissions WHERE permission_id = ?::uuid`,
		permReadID,
	).Scan(&totalRead).Error; err != nil {
		t.Fatalf("count total :read grants: %v", err)
	}
	if totalRead != int64(len(readRoles)) {
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
		var cnt int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM iam.role_permissions `+
				`WHERE role_id = ?::uuid AND permission_id = ?::uuid`,
			rid, permUpdateID,
		).Scan(&cnt).Error; err != nil {
			t.Fatalf("count grant role=%s perm=update: %v", rid, err)
		}
		if cnt != 1 {
			t.Errorf("role %s missing system_settings:update grant (got %d, want 1)", rid, cnt)
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
		var cnt int64
		if err := gdb.Raw(
			`SELECT COUNT(*) FROM iam.role_permissions `+
				`WHERE role_id = ?::uuid AND permission_id = ?::uuid`,
			rid, permUpdateID,
		).Scan(&cnt).Error; err != nil {
			t.Fatalf("count forbidden grant role=%s perm=update: %v", rid, err)
		}
		if cnt != 0 {
			t.Errorf("role %s unexpectedly has system_settings:update grant (got %d, want 0)", rid, cnt)
		}
	}
	var totalUpdate int64
	if err := gdb.Raw(
		`SELECT COUNT(*) FROM iam.role_permissions WHERE permission_id = ?::uuid`,
		permUpdateID,
	).Scan(&totalUpdate).Error; err != nil {
		t.Fatalf("count total :update grants: %v", err)
	}
	if totalUpdate != int64(len(updateRoles)) {
		t.Errorf("total :update grants: got %d, want %d", totalUpdate, len(updateRoles))
	}
}

// loadRolePermissionIDs retorna el set de permission_id (como string)
// asignados a un role_id.
func loadRolePermissionIDs(t *testing.T, gdb *gorm.DB, roleID string) map[string]struct{} {
	t.Helper()
	var ids []string
	if err := gdb.Raw(
		`SELECT permission_id::text FROM iam.role_permissions WHERE role_id = ?::uuid`,
		roleID,
	).Scan(&ids).Error; err != nil {
		t.Fatalf("load role_permissions for %s: %v", roleID, err)
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
