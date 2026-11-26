package migrations_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground_v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	db := setupPostgres(t)
	defer func() { _ = db.Close() }()

	result, err := migrations.Migrate(db, migrations.MigrateOptions{
		Force:  true,
		DBUser: "test",
	})
	if err != nil {
		t.Fatalf("Error ejecutando migración completa: %v", err)
	}
	if result.Skipped {
		t.Fatalf("La migración no debería quedar en skip en entorno limpio")
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Error abriendo GORM: %v", err)
	}

	t.Run("SchemaAndIndexes", func(t *testing.T) {
		expectedSchemas := []string{
			"auth", "iam", "academic", "content", "assessment", "ui_config", "audit", "notifications",
		}
		for _, schema := range expectedSchemas {
			assertSchemaExists(t, db, schema)
		}

		assertTableExists(t, db, "auth", "users")
		assertTableExists(t, db, "auth", "identity_sync_outbox")
		assertTableExists(t, db, "iam", "resources")
		assertTableExists(t, db, "academic", "schools")
		assertTableExists(t, db, "content", "materials")
		assertTableExists(t, db, "assessment", "assessment")
		assertTableExists(t, db, "ui_config", "screen_instances")
		assertTableExists(t, db, "audit", "events")
		assertTableExists(t, db, "notifications", "notifications")

		expectedAuditIndexes := []string{
			"idx_audit_events_actor",
			"idx_audit_events_action",
			"idx_audit_events_resource",
			"idx_audit_events_school",
			"idx_audit_events_created",
			"idx_audit_events_severity",
		}
		for _, idx := range expectedAuditIndexes {
			assertIndexExists(t, db, "audit", "events", idx)
		}
	})

	t.Run("LoginAttemptsRetired", func(t *testing.T) {
		// Plan 053 T5.1: auth.login_attempts ya no se crea, y post_gorm.sql
		// lleva el DROP explícito para una BD que la arrastre. post_gorm.sql
		// entero no se puede re-aplicar sobre una BD ya migrada (no es del todo
		// idempotente), así que se ejecuta solo su sentencia de DROP.
		if tableExists(t, db, "auth", "login_attempts") {
			t.Fatalf("auth.login_attempts no debería existir tras migrar")
		}
		const dropStmt = "DROP TABLE IF EXISTS auth.login_attempts;"
		content, err := os.ReadFile("sql/post_gorm.sql")
		if err != nil {
			t.Fatalf("Error leyendo post_gorm.sql: %v", err)
		}
		if !strings.Contains(string(content), dropStmt) {
			t.Fatalf("post_gorm.sql no contiene %q", dropStmt)
		}
		if _, err := db.Exec(`CREATE TABLE auth.login_attempts (id uuid PRIMARY KEY)`); err != nil {
			t.Fatalf("Error creando la tabla vieja: %v", err)
		}
		if _, err := db.Exec(dropStmt); err != nil {
			t.Fatalf("Error ejecutando el DROP: %v", err)
		}
		if tableExists(t, db, "auth", "login_attempts") {
			t.Fatalf("el DROP no borró auth.login_attempts")
		}
	})

	t.Run("LoginPropioRetired", func(t *testing.T) {
		// Plan 053 T6.2: auth.refresh_tokens y las columnas password_hash y
		// token_version de auth.users ya no se crean, y post_gorm.sql lleva los
		// DROP explícitos para una BD que las arrastre. Se ejecutan solo esas
		// sentencias (post_gorm.sql entero no es re-aplicable), dos veces: la
		// segunda prueba que son idempotentes sobre una BD que ya no las tiene.
		if tableExists(t, db, "auth", "refresh_tokens") {
			t.Fatalf("auth.refresh_tokens no debería existir tras migrar")
		}
		for _, col := range []string{"password_hash", "token_version"} {
			if columnExists(t, db, "auth", "users", col) {
				t.Fatalf("auth.users.%s no debería existir tras migrar", col)
			}
		}
		dropStmts := []string{
			"DROP TABLE IF EXISTS auth.refresh_tokens;",
			"ALTER TABLE auth.users DROP COLUMN IF EXISTS password_hash;",
			"ALTER TABLE auth.users DROP COLUMN IF EXISTS token_version;",
		}
		content, err := os.ReadFile("sql/post_gorm.sql")
		if err != nil {
			t.Fatalf("Error leyendo post_gorm.sql: %v", err)
		}
		for _, stmt := range dropStmts {
			if !strings.Contains(string(content), stmt) {
				t.Fatalf("post_gorm.sql no contiene %q", stmt)
			}
		}
		for _, stmt := range []string{
			`CREATE TABLE auth.refresh_tokens (id uuid PRIMARY KEY, user_id uuid REFERENCES auth.users(id) ON DELETE CASCADE)`,
			`ALTER TABLE auth.users ADD COLUMN password_hash varchar(255) NOT NULL DEFAULT 'x'`,
			`ALTER TABLE auth.users ADD COLUMN token_version integer NOT NULL DEFAULT 1`,
		} {
			if _, err := db.Exec(stmt); err != nil {
				t.Fatalf("Error recreando el esquema viejo (%s): %v", stmt, err)
			}
		}
		for pass := 1; pass <= 2; pass++ {
			for _, stmt := range dropStmts {
				if _, err := db.Exec(stmt); err != nil {
					t.Fatalf("pasada %d: error ejecutando %q: %v", pass, stmt, err)
				}
			}
			if tableExists(t, db, "auth", "refresh_tokens") {
				t.Fatalf("pasada %d: el DROP no borró auth.refresh_tokens", pass)
			}
			for _, col := range []string{"password_hash", "token_version"} {
				if columnExists(t, db, "auth", "users", col) {
					t.Fatalf("pasada %d: el DROP no borró auth.users.%s", pass, col)
				}
			}
		}
	})

	t.Run("ProductionSeedsApplied", func(t *testing.T) {
		// Floors POST-Fase-6 (rebuild de seeds, ADR-6).
		//
		// Los conteos legacy (roles≥12, role_permissions≥300) ya no
		// aplican: F6 limita el sistema a los 6 roles efectivamente
		// implementados en KMP (super_admin L0 + announcement_viewer L1
		// + 5 L4: student/teacher/guardian/admin/school_admin). El resto
		// de roles del legacy se descartaron como `role_unused`
		// (decisions-log B2-D1).
		//
		// P4-1 (plan B): iam.role_permissions fue eliminada. El conteo
		// de filas se valida ahora sobre iam.role_grants (patterns
		// wildcard).
		//
		// El piso de role_grants era 130 («12 roles × 11 patterns prom»),
		// de cuando cada alias copiaba los grants de su canónico. Hoy el
		// seed L4 (roleGrantPatterns + roleGrantDenyPatterns) declara 95:
		// 59 allow (super_admin 1, announcement_viewer 2, student 11,
		// teacher 12, guardian 13, school_admin 10, readonly_auditor 10)
		// y 36 deny (school_admin 4, readonly_auditor 32). Los 5 alias
		// que heredan por parent_role_id no declaran grants propios
		// (ADR-6), PRE-4 retiró platform_admin, y los planes 027 y 052 F4
		// podaron grants de consumo.
		assertCountAtLeast(t, db, "iam.resources", 30)             // 1 L0 + 1 L3 + 31 L4
		assertCountAtLeast(t, db, "iam.roles", 7)                  // 1 L0 + 1 L1 + 5 L4
		assertCountAtLeast(t, db, "iam.permissions", 90)           // 4 L0 + 3 L3 + ~89 L4
		assertCountAtLeast(t, db, "iam.role_grants", 95)           // 59 allow + 36 deny (ver arriba)
		assertCountAtLeast(t, db, "ui_config.screen_templates", 8) // 3 L0 + 5 L4
		assertCountAtLeast(t, db, "ui_config.screen_instances", 1) // sube a ~78 al cerrar B4
		assertCountAtLeast(t, db, "ui_config.resource_screens", 1) // sube a ~64 al cerrar B5
		assertCountAtLeast(t, db, "academic.concept_types", 5)
	})

	t.Run("CRUD_Users_GORM", func(t *testing.T) {
		user := entities.User{
			ID:        uuid.New(),
			Email:     fmt.Sprintf("integration-%s@edugo.test", uuid.NewString()),
			FirstName: "Integration",
			LastName:  "User",
			IsActive:  true,
		}
		if err := gdb.Create(&user).Error; err != nil {
			t.Fatalf("Error creando user: %v", err)
		}

		var got entities.User
		if err := gdb.First(&got, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Error leyendo user: %v", err)
		}
		if got.Email != user.Email {
			t.Fatalf("Email inesperado: %s", got.Email)
		}

		if err := gdb.Model(&entities.User{}).
			Where("id = ?", user.ID).
			Update("first_name", "Updated").Error; err != nil {
			t.Fatalf("Error actualizando user: %v", err)
		}

		if err := gdb.Unscoped().Delete(&entities.User{}, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("Error eliminando user: %v", err)
		}
	})

	t.Run("CRUD_Schools_GORM", func(t *testing.T) {
		school := entities.School{
			ID:               uuid.New(),
			Name:             "Escuela Integración",
			Code:             "SCH_IT_001",
			Country:          "Chile",
			SubscriptionTier: "basic",
			MaxTeachers:      5,
			MaxStudents:      100,
			IsActive:         true,
		}
		if err := gdb.Create(&school).Error; err != nil {
			t.Fatalf("Error creando school: %v", err)
		}

		if err := gdb.Model(&entities.School{}).
			Where("id = ?", school.ID).
			Update("city", "Santiago").Error; err != nil {
			t.Fatalf("Error actualizando school: %v", err)
		}

		if err := gdb.Unscoped().Delete(&entities.School{}, "id = ?", school.ID).Error; err != nil {
			t.Fatalf("Error eliminando school: %v", err)
		}
	})

	t.Run("CRUD_Materials_GORM", func(t *testing.T) {
		seedSchoolID := mustUUID(t, "b1000000-0000-0000-0000-000000000001")
		// N4 (ADR 0019): content.materials.uploaded_by_membership_id es FK a
		// academic.memberships (era uploaded_by_teacher_id → auth.users). Se usa
		// la membresía docente bb…08 (user …05, school b1…01, unidad ac…03).
		//
		// MP-09: migrations.Migrate solo siembra system/, que es contrato
		// puro (sin escuelas, unidades ni membresías), y el paquete
		// seeds/demo que aportaba estas filas se eliminó. El mundo de datos
		// es playground_v2 `base` (el que `make docker-recreate` aplica por
		// defecto), así que el test lo siembra antes de usar sus IDs.
		if err := playground_v2.Apply(gdb, "base"); err != nil {
			t.Fatalf("Error aplicando playground_v2 base: %v", err)
		}
		seedTeacherMembershipID := mustUUID(t, "bb000000-0000-0000-0000-000000000008")
		seedUnitID := mustUUID(t, "ac000000-0000-0000-0000-000000000003")

		// F2 (plan 018): el material es el TEMA (maestro); los datos de archivo
		// viven en content.material_file (N por tema). El tema ya no lleva
		// file_url/file_type/file_size_bytes inline pero sí el `summary` (markdown).
		summary := "## Resumen\nNotas del docente."
		material := entities.Material{
			ID:                     uuid.New(),
			SchoolID:               seedSchoolID,
			UploadedByMembershipID: seedTeacherMembershipID,
			AcademicUnitID:         &seedUnitID,
			Title:                  "Material Integración",
			Summary:                &summary,
			Status:                 "ready",
			IsPublic:               true,
		}

		if err := gdb.Create(&material).Error; err != nil {
			t.Fatalf("Error creando material: %v", err)
		}

		// El hijo material_file: FK material_id → content.materials ON DELETE CASCADE.
		file := entities.MaterialFile{
			ID:            uuid.New(),
			MaterialID:    material.ID,
			FileURL:       "s3://integration/material.pdf",
			FileName:      "material.pdf",
			FileType:      "application/pdf",
			FileSizeBytes: 2048,
		}
		if err := gdb.Create(&file).Error; err != nil {
			t.Fatalf("Error creando material_file: %v", err)
		}

		if err := gdb.Model(&entities.Material{}).
			Where("id = ?", material.ID).
			Update("title", "Material Integración Updated").Error; err != nil {
			t.Fatalf("Error actualizando material: %v", err)
		}

		// Borrar el tema arrastra el hijo por la FK CASCADE; lo eliminamos
		// explícitamente primero para no depender del borrado lógico (DeletedAt)
		// del maestro, que NO cascadea a nivel BD.
		if err := gdb.Delete(&entities.MaterialFile{}, "id = ?", file.ID).Error; err != nil {
			t.Fatalf("Error eliminando material_file: %v", err)
		}

		if err := gdb.Unscoped().Delete(&entities.Material{}, "id = ?", material.ID).Error; err != nil {
			t.Fatalf("Error eliminando material: %v", err)
		}
	})

	t.Run("CRUD_AuditEvents_GORM", func(t *testing.T) {
		actorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		event := entities.AuditEvent{
			ID:           uuid.New(),
			ActorID:      actorID,
			ActorEmail:   "integration@edugo.test",
			ActorRole:    "super_admin",
			ServiceName:  "integration-test",
			Action:       "user.login",
			ResourceType: "auth",
			Severity:     "info",
			Category:     "auth",
			Changes:      map[string]any{"from": "none", "to": "login"},
			Metadata:     map[string]any{"source": "integration"},
		}

		if err := gdb.Create(&event).Error; err != nil {
			t.Fatalf("Error creando audit.event: %v", err)
		}

		if err := gdb.Model(&entities.AuditEvent{}).
			Where("id = ?", event.ID).
			Update("severity", "warning").Error; err != nil {
			t.Fatalf("Error actualizando audit.event: %v", err)
		}

		if err := gdb.Delete(&entities.AuditEvent{}, "id = ?", event.ID).Error; err != nil {
			t.Fatalf("Error eliminando audit.event: %v", err)
		}
	})
}

func setupPostgres(t *testing.T) *sql.DB {
	t.Helper()
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
		t.Fatalf("Error creando container: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Error obteniendo puerto: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host,
		port.Port(),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error abriendo conexión: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Error conectando PostgreSQL: %v", err)
	}
	return db
}

func assertSchemaExists(t *testing.T, db *sql.DB, schema string) {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.schemata
			WHERE schema_name = $1
		)
	`, schema).Scan(&exists)
	if err != nil {
		t.Fatalf("Error verificando schema %s: %v", schema, err)
	}
	if !exists {
		t.Fatalf("Schema %s no existe", schema)
	}
}

func assertTableExists(t *testing.T, db *sql.DB, schema, table string) {
	t.Helper()
	if !tableExists(t, db, schema, table) {
		t.Fatalf("Tabla %s.%s no existe", schema, table)
	}
}

func tableExists(t *testing.T, db *sql.DB, schema, table string) bool {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = $1
			  AND table_name = $2
		)
	`, schema, table).Scan(&exists)
	if err != nil {
		t.Fatalf("Error verificando tabla %s.%s: %v", schema, table, err)
	}
	return exists
}

func columnExists(t *testing.T, db *sql.DB, schema, table, column string) bool {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_schema = $1 AND table_name = $2 AND column_name = $3
		)
	`, schema, table, column).Scan(&exists)
	if err != nil {
		t.Fatalf("Error verificando columna %s.%s.%s: %v", schema, table, column, err)
	}
	return exists
}

func assertIndexExists(t *testing.T, db *sql.DB, schema, table, index string) {
	t.Helper()
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = $1
			  AND tablename = $2
			  AND indexname = $3
		)
	`, schema, table, index).Scan(&exists)
	if err != nil {
		t.Fatalf("Error verificando índice %s: %v", index, err)
	}
	if !exists {
		t.Fatalf("Índice %s no existe en %s.%s", index, schema, table)
	}
}

func assertCountAtLeast(t *testing.T, db *sql.DB, table string, minCount int) {
	t.Helper()
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	if err != nil {
		t.Fatalf("Error contando filas en %s: %v", table, err)
	}
	if count < minCount {
		t.Fatalf("Conteo insuficiente en %s: %d < %d", table, count, minCount)
	}
}

func mustUUID(t *testing.T, value string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(value)
	if err != nil {
		t.Fatalf("UUID inválido %s: %v", value, err)
	}
	return id
}
