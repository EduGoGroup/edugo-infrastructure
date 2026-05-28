package framework

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Cleaner ejecuta el cleanup selectivo de un scenario: borra
// exclusivamente las filas con TenantPrefix/SchemaPrefix del scenario
// y respeta el orden topológico inverso al de creación.
type Cleaner struct {
	Registry *Registry
	Logger   Logger
	Now      func() time.Time
}

// NewCleaner construye un cleaner con valores por defecto.
func NewCleaner(reg *Registry, log Logger) *Cleaner {
	if reg == nil {
		reg = DefaultRegistry
	}
	if log == nil {
		log = NewJSONLogger()
	}
	return &Cleaner{Registry: reg, Logger: log, Now: time.Now}
}

// Cleanup elimina los datos del scenario `scenarioName` siguiendo el
// orden inverso de aplicación. Si el scenario nunca corrió contra esta
// BD el resultado es un no-op con un warning estructurado (C-REQ-3.3).
//
// La operación se envuelve en una transacción: si el borrado falla
// por una FK no contemplada, todo se revierte y se emite un
// `fixture.error` con el último table tocado.
func (c *Cleaner) Cleanup(db *gorm.DB, scenarioName string) error {
	scenario, err := c.Registry.LookupScenario(scenarioName)
	if err != nil {
		return err
	}
	release, err := c.Registry.AcquireApplyLock(scenarioName)
	if err != nil {
		return fmt.Errorf("cleanup already in progress: %s", scenarioName)
	}
	defer release()

	tenantPrefix, schemaPrefix := Derive(scenarioName)
	ctx := NewApplyContext(scenarioName, tenantPrefix, schemaPrefix)

	fixtures := scenario.BuildFixtures(ctx)
	if len(fixtures) == 0 {
		c.Logger.Emit(LogEntry{
			Event:        EventFixtureCleanup,
			Scenario:     scenarioName,
			TenantPrefix: tenantPrefix,
		})
		return nil
	}

	resolved, err := resolve(fixtures)
	if err != nil {
		return err
	}
	// Cleanup en orden inverso (C-REQ-3.5).
	reverse(resolved)

	if db == nil {
		return fmt.Errorf("cleaner.Cleanup: nil db")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, f := range resolved {
			manifest := f.Manifest()
			start := c.Now()
			if err := f.Cleanup(tx, ctx); err != nil {
				c.Logger.Emit(LogEntry{
					Event:        EventFixtureError,
					Scenario:     scenarioName,
					Fixture:      manifest.Name,
					TenantPrefix: tenantPrefix,
					Stage:        "cleanup",
					Error:        err.Error(),
				})
				return fmt.Errorf("cleanup failed: fixture %q: %w", manifest.Name, err)
			}
			c.Logger.Emit(LogEntry{
				Event:        EventFixtureCleanup,
				Scenario:     scenarioName,
				Fixture:      manifest.Name,
				TenantPrefix: tenantPrefix,
				Tables:       manifest.Tables,
				DurationMs:   c.Now().Sub(start).Milliseconds(),
			})
		}
		return nil
	})
}

// DeleteByPrefix ejecuta un DELETE selectivo basado en LIKE sobre una
// columna textual (UUID o code). Devuelve filas afectadas para que la
// fixture pueda alimentar el log estructurado.
//
// idColumn debe ser una columna UUID/text directamente comparable con
// `<col>::text LIKE 'prefix%'`. En PostgreSQL esto evita falsos
// positivos cuando la columna es UUID.
func DeleteByPrefix(tx *gorm.DB, table, idColumn, prefix string) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf("delete_by_prefix: nil transaction")
	}
	if table == "" || idColumn == "" || prefix == "" {
		return 0, fmt.Errorf("delete_by_prefix: empty table/idColumn/prefix")
	}
	// Validación defensiva contra SQL injection en identificadores: el
	// caller no debe poder pasar nada raro. Aceptamos `[a-zA-Z0-9_."]`.
	if !isSafeIdentifier(table) || !isSafeIdentifier(idColumn) {
		return 0, fmt.Errorf("delete_by_prefix: unsafe identifier (table=%q col=%q)", table, idColumn)
	}
	stmt := fmt.Sprintf("DELETE FROM %s WHERE %s::text LIKE ?", table, idColumn)
	res := tx.Exec(stmt, prefix+"%")
	if res.Error != nil {
		return 0, fmt.Errorf("delete_by_prefix: %s where %s LIKE %s%%: %w", table, idColumn, prefix, res.Error)
	}
	return res.RowsAffected, nil
}

func isSafeIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_' || r == '.' || r == '"':
		default:
			return false
		}
	}
	return true
}

// reverse invierte un slice de fixtures en sitio.
func reverse(fs []Fixture) {
	for i, j := 0, len(fs)-1; i < j; i, j = i+1, j-1 {
		fs[i], fs[j] = fs[j], fs[i]
	}
}

// FormatPrefixedClause produce una cláusula útil para logs de errores
// FK detallados (C-REQ-3.4).
func FormatPrefixedClause(table, column, prefix string) string {
	return fmt.Sprintf("%s.%s LIKE %q", table, column, strings.TrimSuffix(prefix, "%")+"%")
}
