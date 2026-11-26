//go:build integration
// +build integration

package demo_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2/demo"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// accountIDs son los cuatro `sub` de identity-core.
var accountIDs = []string{demo.TeacherUserID, demo.Student1UserID, demo.Student2UserID, demo.GuardianUserID}

// TestDemoWorld aplica demo como lo hace el migrador (Force + --playground-v2=demo)
// y afirma la estructura del mundo.
func TestDemoWorld(t *testing.T) {
	gdb := startDemoDB(t)
	school := demo.SchoolID.String()
	unit := demo.UnitID.String()
	root := demo.RootUnitID.String()

	// Una sola escuela, la de demo, perfil basic, con un tipo de concepto de system/.
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.schools`)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.schools s
		JOIN academic.concept_types ct ON ct.id = s.concept_type_id
		WHERE s.id = ? AND s.grade_profile = 'basic' AND s.name = 'Colegio Los Aromos'`, school)
	if n := scalarInt(t, gdb, `SELECT count(*) FROM academic.school_invitation_roles WHERE school_id = ?`, school); n == 0 {
		t.Errorf("la escuela de demo no tiene equivalencias de invitación")
	}

	// Sus conceptos, copiados del tipo, con ids del rango propio (nacen al escribir, no en el mundo en Go).
	if n := scalarInt(t, gdb, `SELECT count(*) FROM academic.school_concepts WHERE school_id = ?`, school); n == 0 {
		t.Errorf("la escuela de demo no tiene conceptos")
	}
	expectInt(t, gdb, 0, `SELECT count(*) FROM academic.school_concepts WHERE school_id = ? AND id::text NOT LIKE ?`,
		school, demo.IDPrefix+"%")

	// La raíz y una unidad de trabajo hija de la raíz.
	expectInt(t, gdb, 2, `SELECT count(*) FROM academic.academic_units WHERE school_id = ?`, school)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.academic_units
		WHERE id = ? AND parent_unit_id IS NULL AND (metadata->>'is_root')::boolean`, root)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.academic_units WHERE id = ? AND parent_unit_id = ? AND name = '8.º Básico A'`, unit, root)

	// Un solo período activo, en la raíz, que cubre hoy.
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.academic_periods WHERE is_active`)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.academic_periods
		WHERE id = ? AND academic_unit_id = ? AND is_active AND start_date <= now() AND end_date >= now()`,
		demo.PeriodID.String(), root)

	// Tres materias y tres sesiones de la docente en la unidad de trabajo.
	expectInt(t, gdb, 3, `SELECT count(*) FROM academic.subjects WHERE school_id = ? AND academic_unit_id IS NULL`, school)
	expectInt(t, gdb, 3, `SELECT count(*) FROM academic.subject_offerings
		WHERE school_id = ? AND academic_unit_id = ? AND teacher_membership_id = ? AND period_id = ? AND is_active`,
		school, unit, demo.TeacherMembershipID.String(), demo.PeriodID.String())

	// Docente y alumnos: una sola membresía activa, en la unidad de trabajo.
	for _, id := range []string{demo.TeacherUserID, demo.Student1UserID, demo.Student2UserID} {
		expectInt(t, gdb, 1, `SELECT count(*) FROM academic.memberships WHERE user_id = ?`, id)
		expectInt(t, gdb, 1, `SELECT count(*) FROM academic.memberships
			WHERE user_id = ? AND academic_unit_id = ? AND status = 'active'`, id, unit)
	}
	expectInt(t, gdb, int64(3+demo.FillerCount), `SELECT count(*) FROM academic.memberships WHERE school_id = ?`, school)
	// Cada alumno, inscrito en las tres sesiones.
	expectInt(t, gdb, int64(3*(2+demo.FillerCount)), `SELECT count(*) FROM academic.subject_offering_enrollments`)

	// La tutora: sin membresía, con un vínculo activo y principal con Antonia.
	expectInt(t, gdb, 0, `SELECT count(*) FROM academic.memberships WHERE user_id = ?`, demo.GuardianUserID)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.guardian_relations WHERE guardian_id = ?`, demo.GuardianUserID)
	expectInt(t, gdb, 1, `SELECT count(*) FROM academic.guardian_relations
		WHERE guardian_id = ? AND student_id = ? AND school_id = ? AND is_active AND is_primary AND status = 'active'`,
		demo.GuardianUserID, demo.Student2UserID, school)

	// Los roles de las personas de demo llevan school_id; el único global es el de L0.
	demoUsers := scalarInt(t, gdb, `SELECT count(*) FROM auth.users WHERE id <> ?`, layers.L0_USER_SUPER_ADMIN_ID)
	expectInt(t, gdb, demoUsers, `SELECT count(DISTINCT user_id) FROM iam.user_roles WHERE school_id = ?`, school)
	expectInt(t, gdb, 0, `SELECT count(*) FROM iam.user_roles WHERE school_id IS NULL AND user_id <> ?`, layers.L0_USER_SUPER_ADMIN_ID)

	// Los cuatro ids de cuenta, iguales a las constantes y con su correo.
	emails := map[string]string{
		demo.TeacherUserID: demo.TeacherEmail, demo.Student1UserID: demo.Student1Email,
		demo.Student2UserID: demo.Student2Email, demo.GuardianUserID: demo.GuardianEmail,
	}
	for id, email := range emails {
		expectInt(t, gdb, 1, `SELECT count(*) FROM auth.users WHERE id = ? AND email = ? AND is_active`, id, email)
	}

	// auth.users = las personas de demo más el usuario de L0.
	expectInt(t, gdb, int64(4+demo.FillerCount), `SELECT count(*) FROM auth.users WHERE id <> ?`, layers.L0_USER_SUPER_ADMIN_ID)
	expectInt(t, gdb, int64(5+demo.FillerCount), `SELECT count(*) FROM auth.users`)

	// Lo fallado (HU-IA.7, decisión 10): tres filas, cada una sobre un intento de E-R.
	expectInt(t, gdb, demo.AIFailureCount, `SELECT count(*) FROM assessment.ai_failure`)
	expectInt(t, gdb, demo.AIFailureCount, `SELECT count(*) FROM assessment.ai_failure f
		JOIN assessment.assessment_attempt a ON a.id = f.target_id
		WHERE f.target_type = 'attempt' AND a.assessment_id = ?`, demo.AssessmentReviewID.String())
	for n, want := range map[int]string{1: "failed/service/5", 2: "failed/data/0", 3: "retrying/service/2"} {
		expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.ai_failure
			WHERE id = ? AND target_id = ? AND status || '/' || cause_class || '/' || retry_count = ?`,
			demo.AIFailureID(n), demo.ReviewAttemptID(n), want)
	}
	expectInt(t, gdb, 1, `SELECT count(*) FROM assessment.ai_failure WHERE next_retry_at IS NOT NULL`)

	// Sincronización apagada: nada en el outbox.
	expectInt(t, gdb, 0, `SELECT count(*) FROM auth.identity_sync_outbox`)

	// Re-aplicar sobre su propio mundo no cambia nada.
	before := scalarInt(t, gdb, `SELECT count(*) FROM academic.attendance`)
	if err := playground_v2.Apply(gdb, "demo"); err != nil {
		t.Fatalf("segunda aplicación de demo: %v", err)
	}
	expectInt(t, gdb, before, `SELECT count(*) FROM academic.attendance`)
}

// TestDemoNotInAll: una base migrada con "all" no trae ninguna fila de la escuela
// de demo, y demo se niega a aplicarse encima de ella.
//
// Se migra con "all" desde cero en vez de aplicar "all" sobre el mundo de demo:
// base no trunca notifications.notifications y esa tabla no tiene FK a
// auth.users, así que los avisos de demo quedarían huérfanos (es de base, no de
// demo; en el ambiente de demo el migrador siempre fuerza la recreación).
func TestDemoNotInAll(t *testing.T) {
	gdb := startDB(t, "all")
	school := demo.SchoolID.String()
	for _, q := range []string{
		`SELECT count(*) FROM academic.schools WHERE id = ?`,
		`SELECT count(*) FROM academic.memberships WHERE school_id = ?`,
		`SELECT count(*) FROM assessment.assessment WHERE school_id = ?`,
		`SELECT count(*) FROM content.materials WHERE school_id = ?`,
		`SELECT count(*) FROM notifications.notifications WHERE school_id = ?`,
		`SELECT count(*) FROM academic.announcements WHERE school_id = ?`,
	} {
		expectInt(t, gdb, 0, q, school)
	}
	expectInt(t, gdb, 0, `SELECT count(*) FROM auth.users WHERE id IN ?`, accountIDs)

	err := playground_v2.Apply(gdb, "demo")
	if err == nil || !strings.Contains(err.Error(), "recién recreada") {
		t.Fatalf("demo sobre base debía negarse, devolvió: %v", err)
	}
	expectInt(t, gdb, 0, `SELECT count(*) FROM academic.schools WHERE id = ?`, school)
}

// startDemoDB arranca un Postgres efímero y lo migra como el migrador con
// --playground-v2=demo: Force, system y demo.
func startDemoDB(tb testing.TB) *gorm.DB { return startDB(tb, "demo") }

// startDB arranca un Postgres efímero y lo migra con Force, system y el
// playground indicado.
func startDB(tb testing.TB, playground string) *gorm.DB {
	tb.Helper()
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		tb.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		},
		Started: true,
	})
	if err != nil {
		tb.Fatalf("demo integration: GenericContainer: %v", err)
	}
	tb.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		tb.Fatalf("demo integration: Host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		tb.Fatalf("demo integration: MappedPort: %v", err)
	}
	sqlDB, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port()))
	if err != nil {
		tb.Fatalf("demo integration: sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		tb.Fatalf("demo integration: ping: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	if _, err := migrations.Migrate(sqlDB, migrations.MigrateOptions{
		Force:        true,
		DBUser:       "test",
		PlaygroundV2: playground,
	}); err != nil {
		tb.Fatalf("demo integration: Migrate: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		tb.Fatalf("demo integration: gorm.Open: %v", err)
	}
	return gdb
}

func scalarInt(t *testing.T, gdb *gorm.DB, query string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := gdb.Raw(query, args...).Scan(&n).Error; err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

func expectInt(t *testing.T, gdb *gorm.DB, want int64, query string, args ...any) {
	t.Helper()
	if got := scalarInt(t, gdb, query, args...); got != want {
		t.Errorf("%d, esperado %d: %s %v", got, want, strings.Join(strings.Fields(query), " "), args)
	}
}
