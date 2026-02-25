package migrations_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestIntegration tests de integración con PostgreSQL en testcontainer
// Solo se ejecutan si ENABLE_INTEGRATION_TESTS=true
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set ENABLE_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()

	// Crear testcontainer PostgreSQL
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
	defer func() { _ = container.Terminate(ctx) }()

	// Obtener host y puerto
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Error obteniendo puerto: %v", err)
	}

	// Conectar a PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port())

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Fatalf("Error verificando conexión: %v", err)
	}

	t.Log("Container PostgreSQL creado y conectado")

	// Ejecutar tests
	t.Run("ApplyAll", testApplyAll(db))
	t.Run("CRUD_Users", testCRUDUsers(db))
	t.Run("CRUD_Schools", testCRUDSchools(db))
	t.Run("CRUD_Materials", testCRUDMaterials(db))
}

func testApplyAll(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// Aplicar todas las migraciones
		if err := migrations.ApplyAll(db); err != nil {
			t.Fatalf("Error aplicando migraciones: %v", err)
		}

		// Verificar que los schemas existen
		schemas := []string{"auth", "iam", "academic", "content", "assessment", "ui_config"}
		for _, schema := range schemas {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.schemata
				WHERE schema_name = $1
			)`
			if err := db.QueryRow(query, schema).Scan(&exists); err != nil {
				t.Errorf("Error verificando schema %s: %v", schema, err)
			}
			if !exists {
				t.Errorf("Schema %s no fue creado", schema)
			}
		}

		// Verificar que las tablas existen en sus schemas correctos
		tables := map[string][]string{
			"auth":       {"users", "refresh_tokens", "login_attempts"},
			"iam":        {"resources", "roles", "permissions", "role_permissions", "user_roles"},
			"academic":   {"schools", "academic_units", "memberships", "subjects", "guardian_relations"},
			"content":    {"materials", "material_versions", "progress"},
			"assessment": {"assessment", "assessment_attempt", "assessment_attempt_answer"},
			"ui_config":  {"screen_templates", "screen_instances", "resource_screens", "screen_user_preferences"},
		}

		totalTables := 0
		for schema, tableList := range tables {
			for _, table := range tableList {
				var exists bool
				query := `SELECT EXISTS (
					SELECT FROM information_schema.tables
					WHERE table_schema = $1
					AND table_name = $2
				)`
				if err := db.QueryRow(query, schema, table).Scan(&exists); err != nil {
					t.Errorf("Error verificando tabla %s.%s: %v", schema, table, err)
				}
				if !exists {
					t.Errorf("Tabla %s.%s no fue creada", schema, table)
				}
				totalTables++
			}
		}

		t.Logf("Todos los %d schemas y %d tablas creados correctamente", len(schemas), totalTables)
	}
}

func testCRUDUsers(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// CREATE
		userID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
		_, err := db.Exec(`
			INSERT INTO auth.users (id, email, password_hash, first_name, last_name)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, "test@edugo.com", "hash123", "Test", "User")
		if err != nil {
			t.Fatalf("Error insertando usuario: %v", err)
		}

		// READ
		var email string
		err = db.QueryRow(`SELECT email FROM auth.users WHERE id = $1`, userID).Scan(&email)
		if err != nil {
			t.Fatalf("Error leyendo usuario: %v", err)
		}
		if email != "test@edugo.com" {
			t.Errorf("Datos incorrectos: email=%s", email)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE auth.users SET first_name = $1 WHERE id = $2`, "Updated", userID)
		if err != nil {
			t.Fatalf("Error actualizando usuario: %v", err)
		}

		var firstName string
		err = db.QueryRow(`SELECT first_name FROM auth.users WHERE id = $1`, userID).Scan(&firstName)
		if err != nil || firstName != "Updated" {
			t.Errorf("Update falló: first_name=%s", firstName)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM auth.users WHERE id = $1`, userID)
		if err != nil {
			t.Fatalf("Error eliminando usuario: %v", err)
		}

		var count int
		_ = db.QueryRow(`SELECT COUNT(*) FROM auth.users WHERE id = $1`, userID).Scan(&count)
		if count != 0 {
			t.Errorf("Usuario no fue eliminado")
		}

		t.Log("CRUD users OK")
	}
}

func testCRUDSchools(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// CREATE
		schoolID := "b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"
		_, err := db.Exec(`
			INSERT INTO academic.schools (id, name, code, address, city, country)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, schoolID, "Escuela Test", "SCH_TEST_001", "Calle 123", "Buenos Aires", "Argentina")
		if err != nil {
			t.Fatalf("Error insertando escuela: %v", err)
		}

		// READ
		var name, city string
		err = db.QueryRow(`SELECT name, city FROM academic.schools WHERE id = $1`, schoolID).Scan(&name, &city)
		if err != nil {
			t.Fatalf("Error leyendo escuela: %v", err)
		}
		if name != "Escuela Test" {
			t.Errorf("Nombre incorrecto: %s", name)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE academic.schools SET city = $1 WHERE id = $2`, "Córdoba", schoolID)
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM academic.schools WHERE id = $1`, schoolID)
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}

		t.Log("CRUD schools OK")
	}
}

func testCRUDMaterials(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// Primero crear user y school (FK dependencies)
		userID := "c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33"
		schoolID := "d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44"
		materialID := "e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55"

		_, err := db.Exec(`INSERT INTO auth.users (id, email, password_hash, first_name, last_name) VALUES ($1, $2, $3, $4, $5)`,
			userID, "teacher.mat@test.com", "hash", "Teacher", "Material")
		if err != nil {
			t.Fatalf("Error creando user FK: %v", err)
		}

		_, err = db.Exec(`INSERT INTO academic.schools (id, name, code, address, city, country) VALUES ($1, $2, $3, $4, $5, $6)`,
			schoolID, "School FK", "SCH002", "Address", "City", "Country")
		if err != nil {
			t.Fatalf("Error creando school FK: %v", err)
		}

		// CREATE material
		_, err = db.Exec(`
			INSERT INTO content.materials (id, school_id, uploaded_by_teacher_id, title, file_url, file_type, file_size_bytes)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, materialID, schoolID, userID, "Material Test", "https://example.com/file.pdf", "application/pdf", 1024)
		if err != nil {
			t.Fatalf("Error insertando material: %v", err)
		}

		// READ
		var title string
		err = db.QueryRow(`SELECT title FROM content.materials WHERE id = $1`, materialID).Scan(&title)
		if err != nil {
			t.Fatalf("Error leyendo material: %v", err)
		}
		if title != "Material Test" {
			t.Errorf("Título incorrecto: %s", title)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE content.materials SET title = $1 WHERE id = $2`, "Material Updated", materialID)
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM content.materials WHERE id = $1`, materialID)
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}

		// Cleanup FK
		_, _ = db.Exec(`DELETE FROM academic.schools WHERE id = $1`, schoolID)
		_, _ = db.Exec(`DELETE FROM auth.users WHERE id = $1`, userID)

		t.Log("CRUD materials OK")
	}
}
