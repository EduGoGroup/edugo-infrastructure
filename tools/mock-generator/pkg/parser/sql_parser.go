package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type SQLParser struct{}

func NewSQLParser() *SQLParser {
	return &SQLParser{}
}

// ParseDirectory lee todos los archivos SQL de un directorio y extrae datos INSERT
func (sp *SQLParser) ParseDirectory(dir string) (map[string]*TableData, error) {
	result := make(map[string]*TableData)

	// Listar archivos .sql
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("error listando archivos: %w", err)
	}

	// Parsear cada archivo
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue // Skip archivos que no se pueden leer
		}

		// Dividir en statements individuales
		statements := strings.Split(string(content), ";")

		for _, stmtStr := range statements {
			stmtStr = strings.TrimSpace(stmtStr)
			if stmtStr == "" {
				continue
			}

			// Parsear statement
			stmt, err := sqlparser.Parse(stmtStr)
			if err != nil {
				continue // Skip statements con errores de sintaxis
			}

			// Extraer INSERT statements
			if insert, ok := stmt.(*sqlparser.Insert); ok {
				data := sp.extractInsertData(insert)
				if existing, exists := result[data.Table]; exists {
					// Agregar filas a tabla existente
					existing.Rows = append(existing.Rows, data.Rows...)
				} else {
					result[data.Table] = data
				}
			}
		}
	}

	return result, nil
}

// extractInsertData extrae datos de un INSERT statement
func (sp *SQLParser) extractInsertData(insert *sqlparser.Insert) *TableData {
	// Extraer nombre de tabla
	tableName := insert.Table.Name.String()

	// Extraer columnas
	var columns []string
	for _, col := range insert.Columns {
		columns = append(columns, col.String())
	}

	// Extraer valores
	var rows [][]interface{}
	switch values := insert.Rows.(type) {
	case sqlparser.Values:
		for _, valTuple := range values {
			var row []interface{}
			for _, expr := range valTuple {
				val := sp.evalExpr(expr)
				row = append(row, val)
			}
			rows = append(rows, row)
		}
	}

	return &TableData{
		Table:   tableName,
		Columns: columns,
		Rows:    rows,
	}
}

// evalExpr evalua una expresion SQL y retorna su valor
func (sp *SQLParser) evalExpr(expr sqlparser.Expr) interface{} {
	switch e := expr.(type) {
	case *sqlparser.SQLVal:
		// Valores literales (strings, numeros)
		val := string(e.Val)
		switch e.Type {
		case sqlparser.StrVal:
			return val
		case sqlparser.IntVal:
			return val
		case sqlparser.FloatVal:
			return val
		case sqlparser.HexVal:
			return val
		default:
			return val
		}

	case *sqlparser.FuncExpr:
		// Funciones SQL
		funcName := strings.ToLower(e.Name.String())
		switch funcName {
		case "now":
			return "NOW()"
		case "gen_random_uuid":
			return "UUID()"
		case "current_date":
			return "CURRENT_DATE()"
		case "uuid":
			return "UUID()"
		default:
			return fmt.Sprintf("FUNC(%s)", funcName)
		}

	case *sqlparser.NullVal:
		return nil

	case sqlparser.BoolVal:
		return bool(e)

	case *sqlparser.ColName:
		// Referencia a columna
		return e.Name.String()

	default:
		// Para otros tipos, retornar representación string
		return sqlparser.String(expr)
	}
}
