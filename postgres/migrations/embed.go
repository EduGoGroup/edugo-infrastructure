package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/internal/sqlutil"
)

// Files contiene todos los archivos SQL embebidos (estructura completa con domain schemas)
//
//go:embed structure/*.sql
var Files embed.FS

// ApplyAll ejecuta todos los scripts de structure/ en orden alfabético.
// Crea schemas, funciones, tablas, foreign keys, funciones IAM y vistas.
//
// Uso típico: Inicializar base de datos desde cero
//
// Ejemplo:
//
//	import "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
//	if err := migrations.ApplyAll(db); err != nil {
//	    log.Fatal(err)
//	}
func ApplyAll(db *sql.DB) error {
	return applyLayer(db, "structure")
}

// applyLayer es una función interna que ejecuta todos los scripts de una capa en orden
func applyLayer(db *sql.DB, layer string) error {
	files, err := fs.ReadDir(Files, layer)
	if err != nil {
		// Si el directorio no existe, no es error (puede estar vacío)
		return nil
	}

	// Filtrar y ordenar archivos .sql
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Ejecutar cada archivo
	for _, filename := range sqlFiles {
		path := filepath.Join(layer, filename)
		content, err := Files.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %w", path, err)
		}

		sqlContent := string(content)

		// Verificar si el archivo tiene contenido ejecutable
		if sqlutil.IsEmptyOrComment(sqlContent) {
			continue
		}

		// Ejecutar el SQL
		if _, err := db.Exec(sqlContent); err != nil {
			return fmt.Errorf("error ejecutando %s: %w", path, err)
		}
	}

	return nil
}
