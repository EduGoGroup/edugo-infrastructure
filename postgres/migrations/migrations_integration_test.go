package migrations_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	_ "github.com/lib/pq"
)

// TestIntegration tests de integración con PostgreSQL real
// Solo se ejecutan si ENABLE_INTEGRATION_TESTS=true
func TestIntegration(t *testing.T) {
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set ENABLE_INTEGRATION_TESTS=true to run")
	}

	// Conectar a PostgreSQL
	connStr := os.Getenv("POSTGRES_TEST_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=edugo password=edugo_dev_2024 dbname=edugo_test sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Error verificando conexión: %v", err)
	}

	// Limpiar BD antes de empezar
	dropAllTables(t, db)

	t.Run("ApplyAll", testApplyAll(db))
	t.Run("CRUD_Users", testCRUDUsers(db))
	t.Run("CRUD_Schools", testCRUDSchools(db))
	t.Run("CRUD_Materials", testCRUDMaterials(db))
	t.Run("ApplyMockData", testApplyMockData(db))
}

func testApplyAll(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// Aplicar todas las migraciones
		if err := migrations.ApplyAll(db); err != nil {
			t.Fatalf("Error aplicando migraciones: %v", err)
		}

		// Verificar que las tablas existen
		tables := []string{
			"users", "schools", "academic_units", "memberships",
			"materials", "assessment", "assessment_attempt", "assessment_attempt_answer",
		}

		for _, table := range tables {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)`
			if err := db.QueryRow(query, table).Scan(&exists); err != nil {
				t.Errorf("Error verificando tabla %s: %v", table, err)
			}
			if !exists {
				t.Errorf("Tabla %s no fue creada", table)
			}
		}
	}
}

func testCRUDUsers(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// CREATE
		_, err := db.Exec(`
			INSERT INTO users (id, email, password_hash, role, first_name, last_name)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, "usr_test_1", "test@edugo.com", "hash123", "student", "Test", "User")
		if err != nil {
			t.Fatalf("Error insertando usuario: %v", err)
		}

		// READ
		var email, role string
		err = db.QueryRow(`SELECT email, role FROM users WHERE id = $1`, "usr_test_1").Scan(&email, &role)
		if err != nil {
			t.Fatalf("Error leyendo usuario: %v", err)
		}
		if email != "test@edugo.com" || role != "student" {
			t.Errorf("Datos incorrectos: email=%s, role=%s", email, role)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE users SET role = $1 WHERE id = $2`, "teacher", "usr_test_1")
		if err != nil {
			t.Fatalf("Error actualizando usuario: %v", err)
		}

		err = db.QueryRow(`SELECT role FROM users WHERE id = $1`, "usr_test_1").Scan(&role)
		if err != nil || role != "teacher" {
			t.Errorf("Update falló: role=%s", role)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM users WHERE id = $1`, "usr_test_1")
		if err != nil {
			t.Fatalf("Error eliminando usuario: %v", err)
		}

		var count int
		db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = $1`, "usr_test_1").Scan(&count)
		if count != 0 {
			t.Errorf("Usuario no fue eliminado")
		}
	}
}

func testCRUDSchools(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// CREATE
		_, err := db.Exec(`
			INSERT INTO schools (id, name, address, city, country)
			VALUES ($1, $2, $3, $4, $5)
		`, "sch_test_1", "Escuela Test", "Calle 123", "Buenos Aires", "Argentina")
		if err != nil {
			t.Fatalf("Error insertando escuela: %v", err)
		}

		// READ
		var name, city string
		err = db.QueryRow(`SELECT name, city FROM schools WHERE id = $1`, "sch_test_1").Scan(&name, &city)
		if err != nil {
			t.Fatalf("Error leyendo escuela: %v", err)
		}
		if name != "Escuela Test" {
			t.Errorf("Nombre incorrecto: %s", name)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE schools SET city = $1 WHERE id = $2`, "Córdoba", "sch_test_1")
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM schools WHERE id = $1`, "sch_test_1")
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}
	}
}

func testCRUDMaterials(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// Primero crear school (FK dependency)
		db.Exec(`INSERT INTO schools (id, name, address, city, country) VALUES ($1, $2, $3, $4, $5)`,
			"sch_mat_test", "School FK", "Address", "City", "Country")

		// CREATE material
		_, err := db.Exec(`
			INSERT INTO materials (id, school_id, title, content_type, storage_provider, storage_key, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, "mat_test_1", "sch_mat_test", "Material Test", "document", "s3", "key123", "usr_creator")
		if err != nil {
			t.Fatalf("Error insertando material: %v", err)
		}

		// READ
		var title string
		err = db.QueryRow(`SELECT title FROM materials WHERE id = $1`, "mat_test_1").Scan(&title)
		if err != nil {
			t.Fatalf("Error leyendo material: %v", err)
		}
		if title != "Material Test" {
			t.Errorf("Título incorrecto: %s", title)
		}

		// UPDATE
		_, err = db.Exec(`UPDATE materials SET title = $1 WHERE id = $2`, "Material Updated", "mat_test_1")
		if err != nil {
			t.Fatalf("Error actualizando: %v", err)
		}

		// DELETE
		_, err = db.Exec(`DELETE FROM materials WHERE id = $1`, "mat_test_1")
		if err != nil {
			t.Fatalf("Error eliminando: %v", err)
		}

		// Cleanup FK
		db.Exec(`DELETE FROM schools WHERE id = $1`, "sch_mat_test")
	}
}

func testApplyMockData(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		// Limpiar y aplicar estructura de nuevo
		dropAllTables(t, db)
		if err := migrations.ApplyAll(db); err != nil {
			t.Fatalf("Error aplicando migraciones: %v", err)
		}

		// Aplicar mock data
		if err := migrations.ApplyMockData(db); err != nil {
			t.Fatalf("Error aplicando mock data: %v", err)
		}

		// Verificar que se insertaron datos
		var count int
		db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
		if count == 0 {
			t.Error("Mock data no insertó usuarios")
		}

		db.QueryRow(`SELECT COUNT(*) FROM schools`).Scan(&count)
		if count == 0 {
			t.Error("Mock data no insertó escuelas")
		}

		t.Logf("✅ Mock data aplicado correctamente: %d usuarios, escuelas, etc.", count)
	}
}

// dropAllTables elimina todas las tablas para empezar limpio
func dropAllTables(t *testing.T, db *sql.DB) {
	tables := []string{
		"assessment_attempt_answer",
		"assessment_attempt",
		"assessment",
		"materials",
		"memberships",
		"academic_units",
		"schools",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(`DROP TABLE IF EXISTS ` + table + ` CASCADE`)
		if err != nil {
			t.Logf("Warning: Error eliminando tabla %s: %v", table, err)
		}
	}
}
