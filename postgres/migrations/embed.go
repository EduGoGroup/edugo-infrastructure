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

// Files contiene todos los archivos SQL embebidos (estructura, constraints, seeds, testing)
//
//go:embed structure/*.sql constraints/*.sql seeds/*.sql testing/*.sql
var Files embed.FS

// ApplyAll ejecuta structure + constraints (base de datos limpia lista para usar)
// Equivalente a: ApplyStructure() + ApplyConstraints()
//
// Uso típico: Inicializar base de datos en ambiente de desarrollo o testing
//
// Ejemplo:
//   import "github.com/EduGoGroup/edugo-infrastructure/postgres/migrations"
//   if err := migrations.ApplyAll(db); err != nil {
//       log.Fatal(err)
//   }
func ApplyAll(db *sql.DB) error {
	if err := ApplyStructure(db); err != nil {
		return fmt.Errorf("error aplicando structure: %w", err)
	}
	if err := ApplyConstraints(db); err != nil {
		return fmt.Errorf("error aplicando constraints: %w", err)
	}
	return nil
}

// ApplyStructure ejecuta solo los scripts de structure/ (CREATE TABLE sin FK)
// Crea las tablas base sin foreign keys, índices adicionales, ni triggers
//
// Uso típico: Cuando necesitas crear tablas en orden específico sin dependencias
//
// Ejemplo:
//   migrations.ApplyStructure(db)
func ApplyStructure(db *sql.DB) error {
	return applyLayer(db, "structure")
}

// ApplyConstraints ejecuta solo los scripts de constraints/ (FK, índices, triggers, views)
// DEBE ejecutarse DESPUÉS de ApplyStructure()
//
// Uso típico: Agregar constraints después de haber creado las tablas
//
// Ejemplo:
//   migrations.ApplyStructure(db)
//   migrations.ApplyConstraints(db)
func ApplyConstraints(db *sql.DB) error {
	return applyLayer(db, "constraints")
}

// ApplySeeds ejecuta scripts de seeds/ (datos iniciales del ecosistema)
// Datos básicos necesarios para que el sistema funcione (ej: regiones, configuración)
//
// Uso típico: Inicializar datos mínimos en ambiente de producción/staging
//
// Ejemplo:
//   migrations.ApplyAll(db)
//   migrations.ApplySeeds(db)  // Datos iniciales
func ApplySeeds(db *sql.DB) error {
	return applyLayer(db, "seeds")
}

// ApplyMockData ejecuta scripts de testing/ (datos mock para tests)
// Datos de prueba para desarrollo y testing
//
// Uso típico: Tests de integración, ambiente de desarrollo
//
// Ejemplo:
//   migrations.ApplyAll(db)
//   migrations.ApplyMockData(db)  // Datos de prueba
func ApplyMockData(db *sql.DB) error {
	return applyLayer(db, "testing")
}

// GetScript obtiene el contenido de un script específico como string
// Permite al cliente ejecutar scripts individuales con flexibilidad
//
// Parámetros:
//   - name: Ruta relativa al script, ej: "structure/001_users.sql"
//
// Retorna:
//   - string: Contenido del script
//   - error: Si el script no existe
//
// Ejemplo:
//   script, err := migrations.GetScript("structure/001_users.sql")
//   if err != nil {
//       log.Fatal(err)
//   }
//   db.Exec(script)
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
//   {
//     "structure": ["001_users.sql", "002_schools.sql", ...],
//     "constraints": ["001_users.sql", ...],
//     "seeds": [...],
//     "testing": [...]
//   }
//
// Ejemplo:
//   scripts := migrations.ListScripts()
//   for layer, files := range scripts {
//       fmt.Printf("%s: %v\n", layer, files)
//   }
func ListScripts() map[string][]string {
	result := make(map[string][]string)
	layers := []string{"structure", "constraints", "seeds", "testing"}

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
//   - layer: "structure", "constraints", "seeds", o "testing"
//
// Retorna:
//   - map[string]string: Mapa de nombre_archivo -> contenido
//
// Ejemplo:
//   scripts, err := migrations.GetScriptsByLayer("structure")
//   for name, content := range scripts {
//       fmt.Printf("Script: %s\n", name)
//       db.Exec(content)
//   }
func GetScriptsByLayer(layer string) (map[string]string, error) {
	validLayers := map[string]bool{
		"structure":   true,
		"constraints": true,
		"seeds":       true,
		"testing":     true,
	}

	if !validLayers[layer] {
		return nil, fmt.Errorf("capa inválida: %s (debe ser: structure, constraints, seeds, testing)", layer)
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
