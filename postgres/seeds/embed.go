package seeds

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
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

// GetScript obtiene el contenido de un script específico
func GetScript(name string) (string, error) {
	content, err := Files.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("script no encontrado: %s", name)
	}
	return string(content), nil
}

// ListScripts lista todos los scripts disponibles por capa
func ListScripts() map[string][]string {
	result := make(map[string][]string)
	layers := []string{"production", "development"}
	for _, layer := range layers {
		files, err := fs.ReadDir(Files, layer)
		if err != nil {
			continue
		}
		var sqlFiles []string
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
				sqlFiles = append(sqlFiles, file.Name())
			}
		}
		sort.Strings(sqlFiles)
		result[layer] = sqlFiles
	}
	return result
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
		if isEmptyOrComment(sqlContent) {
			continue
		}
		if _, err := db.Exec(sqlContent); err != nil {
			return fmt.Errorf("error ejecutando %s: %w", path, err)
		}
	}
	return nil
}

func isEmptyOrComment(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
			return false
		}
	}
	return true
}
