package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
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

// GetScript obtiene el contenido de un script específico como string
// Permite al cliente ejecutar scripts individuales con flexibilidad
//
// Parámetros:
//   - name: Ruta relativa al script, ej: "structure/010_auth_users.sql"
//
// Retorna:
//   - string: Contenido del script
//   - error: Si el script no existe
//
// Ejemplo:
//
//	script, err := migrations.GetScript("structure/010_auth_users.sql")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	db.Exec(script)
func GetScript(name string) (string, error) {
	content, err := Files.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("script no encontrado: %s", name)
	}
	return string(content), nil
}

// ListScripts lista todos los scripts disponibles por capa
// Útil para inspeccionar qué scripts están disponibles
//
// Retorna map con estructura:
//
//	{
//	  "structure": ["000_schemas_and_extensions.sql", "001_shared_functions.sql", ...],
//	}
//
// Ejemplo:
//
//	scripts := migrations.ListScripts()
//	for layer, files := range scripts {
//	    fmt.Printf("%s: %v\n", layer, files)
//	}
func ListScripts() map[string][]string {
	result := make(map[string][]string)
	layers := []string{"structure"}

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

// GetScriptsByLayer obtiene todos los scripts de una capa específica como map[nombre]contenido
// Útil cuando necesitas acceso a todos los scripts de una capa
//
// Parámetros:
//   - layer: "structure"
//
// Retorna:
//   - map[string]string: Mapa de nombre_archivo -> contenido
//
// Ejemplo:
//
//	scripts, err := migrations.GetScriptsByLayer("structure")
//	for name, content := range scripts {
//	    fmt.Printf("Script: %s\n", name)
//	    db.Exec(content)
//	}
func GetScriptsByLayer(layer string) (map[string]string, error) {
	validLayers := map[string]bool{
		"structure": true,
	}

	if !validLayers[layer] {
		return nil, fmt.Errorf("capa inválida: %s (debe ser: structure)", layer)
	}

	files, err := fs.ReadDir(Files, layer)
	if err != nil {
		return nil, fmt.Errorf("error leyendo capa %s: %w", layer, err)
	}

	scripts := make(map[string]string)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		path := filepath.Join(layer, file.Name())
		content, err := Files.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error leyendo %s: %w", path, err)
		}

		scripts[file.Name()] = string(content)
	}

	return scripts, nil
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
		if isEmptyOrComment(sqlContent) {
			continue
		}

		// Ejecutar el SQL
		if _, err := db.Exec(sqlContent); err != nil {
			return fmt.Errorf("error ejecutando %s: %w", path, err)
		}
	}

	return nil
}

// isEmptyOrComment verifica si un archivo SQL solo contiene comentarios o está vacío
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
