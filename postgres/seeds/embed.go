package seeds

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

//go:embed production/*.sql development/*.sql
var Files embed.FS

// ApplyProduction ejecuta los seeds de producción (datos del sistema: roles, permisos, resources, ui_config)
func ApplyProduction(db *sql.DB) error {
	return applyLayer(db, "production")
}

// ApplyDevelopment ejecuta los seeds de desarrollo (datos de prueba: escuelas, usuarios, materiales, etc.)
func ApplyDevelopment(db *sql.DB) error {
	return applyLayer(db, "development")
}

func applyLayer(db *sql.DB, layer string) error {
	files, err := fs.ReadDir(Files, layer)
	if err != nil {
		return nil
	}
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)
	for _, filename := range sqlFiles {
		path := filepath.Join(layer, filename)
		content, err := Files.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %w", path, err)
		}
		sqlContent := string(content)
		if sqlutil.IsEmptyOrComment(sqlContent) {
			continue
		}
		if _, err := db.Exec(sqlContent); err != nil {
			return fmt.Errorf("error ejecutando %s: %w", path, err)
		}
	}
	return nil
}
