package migrations_test

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
	"github.com/google/uuid"
)

// migratedPostgres starts a clean Postgres and runs the full migration on it.
func migratedPostgres(t *testing.T) *sql.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	if os.Getenv("ENABLE_INTEGRATION_TESTS") != "true" {
		t.Skip("Set ENABLE_INTEGRATION_TESTS=true to run integration tests")
	}

	db := setupPostgres(t)
	t.Cleanup(func() { _ = db.Close() })

	if _, err := migrations.Migrate(db, migrations.MigrateOptions{Force: true, DBUser: "test"}); err != nil {
		t.Fatalf("Error ejecutando migración completa: %v", err)
	}
	return db
}

// checkConstraints returns the CHECK constraint names of a table, sorted.
func checkConstraints(t *testing.T, db *sql.DB, table string) []string {
	t.Helper()
	rows, err := db.Query(
		`SELECT conname FROM pg_constraint WHERE conrelid = $1::regclass AND contype = 'c' ORDER BY conname`,
		table,
	)
	if err != nil {
		t.Fatalf("Error leyendo los CHECK de %s: %v", table, err)
	}
	defer func() { _ = rows.Close() }()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Error leyendo un CHECK de %s: %v", table, err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Error recorriendo los CHECK de %s: %v", table, err)
	}
	return names
}

func containsName(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

// TestAIFailureTable_Exists (MVP-IA HU-IA.6): assessment.ai_failure exists with the
// contract the rescue routine and the teacher screen consume: columns and types,
// its three CHECK, the unique per target, the due index and the updated_at trigger.
func TestAIFailureTable_Exists(t *testing.T) {
	db := migratedPostgres(t)

	assertTableExists(t, db, "assessment", "ai_failure")

	wantColumns := map[string][2]string{
		"id":            {"uuid", "NO"},
		"target_type":   {"character varying", "NO"},
		"target_id":     {"uuid", "NO"},
		"status":        {"character varying", "NO"},
		"cause_class":   {"character varying", "NO"},
		"last_error":    {"text", "YES"},
		"retry_count":   {"integer", "NO"},
		"next_retry_at": {"timestamp with time zone", "YES"},
		"closed_at":     {"timestamp with time zone", "YES"},
		"created_at":    {"timestamp with time zone", "NO"},
		"updated_at":    {"timestamp with time zone", "NO"},
	}
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'assessment' AND table_name = 'ai_failure'`)
	if err != nil {
		t.Fatalf("Error leyendo las columnas de assessment.ai_failure: %v", err)
	}
	gotColumns := map[string][2]string{}
	for rows.Next() {
		var name, dataType, nullable string
		if err := rows.Scan(&name, &dataType, &nullable); err != nil {
			t.Fatalf("Error leyendo una columna: %v", err)
		}
		gotColumns[name] = [2]string{dataType, nullable}
	}
	_ = rows.Close()
	if len(gotColumns) != len(wantColumns) {
		t.Errorf("assessment.ai_failure tiene %d columnas, se esperaban %d: %v", len(gotColumns), len(wantColumns), gotColumns)
	}
	for name, want := range wantColumns {
		if got, ok := gotColumns[name]; !ok || got != want {
			t.Errorf("columna %s: tipo y nulabilidad %v, se esperaba %v", name, got, want)
		}
	}

	var shortVarchars int
	if err := db.QueryRow(`
		SELECT count(*) FROM information_schema.columns
		WHERE table_schema = 'assessment' AND table_name = 'ai_failure'
		  AND column_name IN ('target_type', 'status', 'cause_class') AND character_maximum_length = 20`,
	).Scan(&shortVarchars); err != nil {
		t.Fatalf("Error leyendo el largo de los varchar: %v", err)
	}
	if shortVarchars != 3 {
		t.Errorf("target_type, status y cause_class deberían ser varchar(20); lo son %d de 3", shortVarchars)
	}

	var primaryKey string
	if err := db.QueryRow(`
		SELECT pg_get_constraintdef(oid) FROM pg_constraint
		WHERE conrelid = 'assessment.ai_failure'::regclass AND contype = 'p'`,
	).Scan(&primaryKey); err != nil {
		t.Fatalf("Error leyendo la PK: %v", err)
	}
	if primaryKey != "PRIMARY KEY (id)" {
		t.Errorf("la PK de assessment.ai_failure es %q, se esperaba PRIMARY KEY (id)", primaryKey)
	}

	checks := checkConstraints(t, db, "assessment.ai_failure")
	t.Logf("CHECK de assessment.ai_failure: %v", checks)
	wantChecks := []string{
		"ai_failure_cause_class_check",
		"ai_failure_status_check",
		"ai_failure_target_type_check",
	}
	if len(checks) != len(wantChecks) {
		t.Errorf("assessment.ai_failure tiene %d CHECK, se esperaban %d", len(checks), len(wantChecks))
	}
	for _, want := range wantChecks {
		if !containsName(checks, want) {
			t.Errorf("falta el CHECK %s", want)
		}
	}

	var indexDef string
	var isUnique bool
	err = db.QueryRow(`
		SELECT pg_get_indexdef(i.indexrelid), i.indisunique
		FROM pg_index i JOIN pg_class c ON c.oid = i.indexrelid
		WHERE i.indrelid = 'assessment.ai_failure'::regclass AND c.relname = 'uq_ai_failure_target'`,
	).Scan(&indexDef, &isUnique)
	if err != nil {
		t.Fatalf("Error leyendo uq_ai_failure_target: %v", err)
	}
	if !isUnique || !strings.Contains(indexDef, "(target_type, target_id)") {
		t.Errorf("uq_ai_failure_target no es único sobre (target_type, target_id): %s", indexDef)
	}
	err = db.QueryRow(`
		SELECT pg_get_indexdef(i.indexrelid), i.indisunique
		FROM pg_index i JOIN pg_class c ON c.oid = i.indexrelid
		WHERE i.indrelid = 'assessment.ai_failure'::regclass AND c.relname = 'idx_ai_failure_due'`,
	).Scan(&indexDef, &isUnique)
	if err != nil {
		t.Fatalf("Error leyendo idx_ai_failure_due: %v", err)
	}
	if isUnique || !strings.Contains(indexDef, "(status, next_retry_at)") {
		t.Errorf("idx_ai_failure_due no es un índice simple sobre (status, next_retry_at): %s", indexDef)
	}

	var foreignKeys int
	if err := db.QueryRow(
		`SELECT count(*) FROM pg_constraint WHERE conrelid = 'assessment.ai_failure'::regclass AND contype = 'f'`,
	).Scan(&foreignKeys); err != nil {
		t.Fatalf("Error contando las FK: %v", err)
	}
	if foreignKeys != 0 {
		t.Errorf("assessment.ai_failure tiene %d FK, el destino va sin FK", foreignKeys)
	}

	const insert = `
		INSERT INTO assessment.ai_failure (target_type, target_id, status, cause_class, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())`
	target := uuid.New()
	if _, err := db.Exec(insert, "attempt", target, "retrying", "service"); err != nil {
		t.Fatalf("Error insertando una fila válida: %v", err)
	}

	var retryCount int
	var nextRetryAt, closedAt sql.NullTime
	if err := db.QueryRow(
		`SELECT retry_count, next_retry_at, closed_at FROM assessment.ai_failure WHERE target_id = $1`, target,
	).Scan(&retryCount, &nextRetryAt, &closedAt); err != nil {
		t.Fatalf("Error leyendo la fila: %v", err)
	}
	if retryCount != 0 || nextRetryAt.Valid || closedAt.Valid {
		t.Errorf("valores por defecto: retry_count=%d next_retry_at=%v closed_at=%v", retryCount, nextRetryAt, closedAt)
	}

	// The same target under the other type is a different target.
	if _, err := db.Exec(insert, "question", target, "failed", "data"); err != nil {
		t.Errorf("el mismo id con otro target_type debería entrar: %v", err)
	}

	rejected := []struct {
		name                          string
		targetType, status, causeKind string
		targetID                      uuid.UUID
		wantInError                   string
	}{
		{"duplicate target", "attempt", "failed", "data", target, "uq_ai_failure_target"},
		{"unknown target_type", "material", "retrying", "service", uuid.New(), "ai_failure_target_type_check"},
		{"unknown status", "attempt", "pending", "service", uuid.New(), "ai_failure_status_check"},
		{"unknown cause_class", "attempt", "retrying", "quality", uuid.New(), "ai_failure_cause_class_check"},
	}
	for _, tc := range rejected {
		_, err := db.Exec(insert, tc.targetType, tc.targetID, tc.status, tc.causeKind)
		if err == nil {
			t.Errorf("%s: la fila debería rechazarse", tc.name)
			continue
		}
		if !strings.Contains(err.Error(), tc.wantInError) {
			t.Errorf("%s: el error no nombra %s: %v", tc.name, tc.wantInError, err)
		}
	}

	// The set_updated_at trigger moves updated_at on every UPDATE.
	if _, err := db.Exec(
		`UPDATE assessment.ai_failure SET updated_at = '2000-01-01T00:00:00Z' WHERE target_id = $1 AND target_type = 'attempt'`,
		target,
	); err != nil {
		t.Fatalf("Error actualizando la fila: %v", err)
	}
	var updatedYear int
	if err := db.QueryRow(
		`SELECT extract(year FROM updated_at)::int FROM assessment.ai_failure WHERE target_id = $1 AND target_type = 'attempt'`,
		target,
	).Scan(&updatedYear); err != nil {
		t.Fatalf("Error leyendo updated_at: %v", err)
	}
	if updatedYear == 2000 {
		t.Error("el trigger set_updated_at no actuó sobre assessment.ai_failure")
	}
}

// TestAnswerReviewStatus_AcceptsCancelled (MVP-IA HU-IA.6, R-062): the CHECK on
// assessment_attempt_answer.review_status is gone, so 'cancelled' can be stored,
// and it is the only CHECK that left: the ones of the attempt and of the question stay.
func TestAnswerReviewStatus_AcceptsCancelled(t *testing.T) {
	db := migratedPostgres(t)

	answerChecks := checkConstraints(t, db, "assessment.assessment_attempt_answer")
	t.Logf("CHECK de assessment.assessment_attempt_answer: %v", answerChecks)
	if len(answerChecks) != 0 {
		t.Errorf("assessment.assessment_attempt_answer no debería tener ningún CHECK, tiene %v", answerChecks)
	}

	var allConstraints []string
	rows, err := db.Query(
		`SELECT conname FROM pg_constraint WHERE conrelid = 'assessment.assessment_attempt_answer'::regclass ORDER BY conname`,
	)
	if err != nil {
		t.Fatalf("Error leyendo los constraints de la respuesta: %v", err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Error leyendo un constraint: %v", err)
		}
		allConstraints = append(allConstraints, name)
	}
	_ = rows.Close()
	t.Logf("constraints de assessment.assessment_attempt_answer: %v", allConstraints)
	wantAnswer := []string{
		"assessment_attempt_answer_attempt_fkey",
		"assessment_attempt_answer_pkey",
		"assessment_attempt_answer_question_fkey",
	}
	if strings.Join(allConstraints, ",") != strings.Join(wantAnswer, ",") {
		t.Errorf("constraints de la respuesta: %v, se esperaban %v", allConstraints, wantAnswer)
	}

	wantKept := map[string][]string{
		"assessment.assessment_attempt": {
			"assessment_attempt_status_check",
			"assessment_attempt_time_check",
			"assessment_attempt_time_spent_check",
		},
		"assessment.question": {
			"question_llm_prep_status_check",
			"question_type_check",
		},
	}
	for table, wants := range wantKept {
		got := checkConstraints(t, db, table)
		for _, want := range wants {
			if !containsName(got, want) {
				t.Errorf("%s perdió el CHECK %s (quedan %v)", table, want, got)
			}
		}
	}

	// The answer needs a parent attempt; the FKs are skipped for this session so
	// the insert only meets the column rules (NOT NULL, CHECK, unique).
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error abriendo la transacción: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`SET LOCAL session_replication_role = replica`); err != nil {
		t.Fatalf("Error desactivando las FK de la sesión: %v", err)
	}
	attemptID := uuid.New()
	for index, status := range []string{"pending", "auto_graded", "reviewed", "ai_reviewed", "cancelled"} {
		if _, err := tx.Exec(`
			INSERT INTO assessment.assessment_attempt_answer (attempt_id, question_index, review_status, created_at, updated_at)
			VALUES ($1, $2, $3, now(), now())`, attemptID, index, status); err != nil {
			t.Errorf("review_status=%s debería insertarse: %v", status, err)
			return
		}
	}
	var cancelled int
	if err := tx.QueryRow(
		`SELECT count(*) FROM assessment.assessment_attempt_answer WHERE attempt_id = $1 AND review_status = 'cancelled'`,
		attemptID,
	).Scan(&cancelled); err != nil {
		t.Fatalf("Error contando las canceladas: %v", err)
	}
	if cancelled != 1 {
		t.Errorf("se esperaba 1 respuesta cancelled, hay %d", cancelled)
	}
}
