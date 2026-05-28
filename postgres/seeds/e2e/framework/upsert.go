package framework

import (
	"fmt"

	"gorm.io/gorm"
)

// UpsertBool actualiza una columna booleana de forma segura, evitando
// la trampa documentada en F2·H5: el tag `gorm:"default:..."` provoca
// que GORM ignore los valores `false` al hacer UPSERT (los considera
// "zero value" y no los envía). Para booleanos críticos como
// IsActive, IsMenuVisible, IsDefault, IsPinned, IsPublic, IsTimed o
// memberships.is_active, hay que tocar la columna por SQL crudo.
//
// La función ejecuta un UPDATE atómico:
//
//	UPDATE <table> SET <col> = <value> WHERE <idCol> = <id>
//
// Si la fila no existe, devuelve un error explícito (la fixture debe
// haberla creado antes de llamar a UpsertBool).
func UpsertBool(tx *gorm.DB, table, idCol string, id any, col string, value bool) error {
	if tx == nil {
		return fmt.Errorf("upsert_bool: nil transaction")
	}
	if table == "" || idCol == "" || col == "" {
		return fmt.Errorf("upsert_bool: empty table/idCol/col (got %q/%q/%q)", table, idCol, col)
	}
	stmt := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, col, idCol)
	res := tx.Exec(stmt, value, id)
	if res.Error != nil {
		return fmt.Errorf("upsert_bool: %s.%s where %s=%v: %w", table, col, idCol, id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("upsert_bool: row not found in %s where %s=%v (call after creating the row)", table, idCol, id)
	}
	return nil
}

// UpsertString es el equivalente para columnas string que también
// pueden colisionar con la trampa de `default:` cuando el valor que
// queremos asignar es la cadena vacía.
func UpsertString(tx *gorm.DB, table, idCol string, id any, col string, value string) error {
	if tx == nil {
		return fmt.Errorf("upsert_string: nil transaction")
	}
	if table == "" || idCol == "" || col == "" {
		return fmt.Errorf("upsert_string: empty table/idCol/col")
	}
	stmt := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, col, idCol)
	res := tx.Exec(stmt, value, id)
	if res.Error != nil {
		return fmt.Errorf("upsert_string: %s.%s where %s=%v: %w", table, col, idCol, id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("upsert_string: row not found in %s where %s=%v", table, idCol, id)
	}
	return nil
}

// UpsertJSON aplica el mismo patrón pero para columnas JSON/JSONB. El
// caller debe pasar el JSON ya serializado.
func UpsertJSON(tx *gorm.DB, table, idCol string, id any, col string, jsonValue []byte) error {
	if tx == nil {
		return fmt.Errorf("upsert_json: nil transaction")
	}
	if table == "" || idCol == "" || col == "" {
		return fmt.Errorf("upsert_json: empty table/idCol/col")
	}
	stmt := fmt.Sprintf("UPDATE %s SET %s = ?::jsonb WHERE %s = ?", table, col, idCol)
	res := tx.Exec(stmt, string(jsonValue), id)
	if res.Error != nil {
		return fmt.Errorf("upsert_json: %s.%s where %s=%v: %w", table, col, idCol, id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("upsert_json: row not found in %s where %s=%v", table, idCol, id)
	}
	return nil
}
