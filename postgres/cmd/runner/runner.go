package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

type Layer struct {
	Name      string
	Directory string
	Color     string
}

var layers = []Layer{
	{Name: "STRUCTURE", Directory: "structure", Color: colorBlue},
	{Name: "CONSTRAINTS", Directory: "constraints", Color: colorPurple},
	{Name: "SEEDS", Directory: "seeds", Color: colorGreen},
	{Name: "TESTING", Directory: "testing", Color: colorCyan},
}

func main() {
	// Configuración de conexión a PostgreSQL
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "edugo")
	dbPassword := getEnv("POSTGRES_PASSWORD", "edugo_dev_2024")
	dbName := getEnv("POSTGRES_DB", "edugo_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("%s✗ Error conectando a PostgreSQL: %v%s\n", colorRed, err, colorReset)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("%s✗ Error verificando conexión: %v%s\n", colorRed, err, colorReset)
	}

	fmt.Printf("%s✓ Conectado a PostgreSQL: %s@%s:%s/%s%s\n\n",
		colorGreen, dbUser, dbHost, dbPort, dbName, colorReset)

	// Obtener directorio base
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s✗ Error obteniendo directorio actual: %v%s\n", colorRed, err, colorReset)
	}

	// Ejecutar cada capa en orden
	totalExecuted := 0
	totalSkipped := 0

	for _, layer := range layers {
		layerDir := filepath.Join(baseDir, layer.Directory)

		// Verificar si el directorio existe
		if _, err := os.Stat(layerDir); os.IsNotExist(err) {
			fmt.Printf("%s⊘ Capa %s: directorio no existe, omitiendo...%s\n\n",
				colorYellow, layer.Name, colorReset)
			continue
		}

		fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", layer.Color, colorReset)
		fmt.Printf("%s  CAPA: %s%s\n", layer.Color, layer.Name, colorReset)
		fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n\n", layer.Color, colorReset)

		executed, skipped := executeLayer(db, layerDir, layer.Color)
		totalExecuted += executed
		totalSkipped += skipped

		fmt.Println()
	}

	// Resumen final
	fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", colorGreen, colorReset)
	fmt.Printf("%s  RESUMEN FINAL%s\n", colorGreen, colorReset)
	fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", colorGreen, colorReset)
	fmt.Printf("%s✓ Archivos ejecutados: %d%s\n", colorGreen, totalExecuted, colorReset)
	fmt.Printf("%s⊘ Archivos omitidos: %d%s\n", colorYellow, totalSkipped, colorReset)
	fmt.Printf("%s✓ Todas las capas procesadas exitosamente%s\n", colorGreen, colorReset)
}

func executeLayer(db *sql.DB, layerDir string, color string) (executed, skipped int) {
	// Leer archivos del directorio
	files, err := os.ReadDir(layerDir)
	if err != nil {
		log.Printf("%s✗ Error leyendo directorio %s: %v%s\n", colorRed, layerDir, err, colorReset)
		return 0, 0
	}

	// Filtrar y ordenar archivos .sql
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		fmt.Printf("%s⊘ No se encontraron archivos SQL en %s%s\n", colorYellow, layerDir, colorReset)
		return 0, 0
	}

	// Ejecutar cada archivo SQL
	for _, filename := range sqlFiles {
		filePath := filepath.Join(layerDir, filename)

		// Leer contenido del archivo
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("%s✗ Error leyendo %s: %v%s\n", colorRed, filename, err, colorReset)
			continue
		}

		sqlContent := string(content)

		// Verificar si el archivo tiene contenido ejecutable
		if isEmptyOrComment(sqlContent) {
			fmt.Printf("%s⊘ %s (vacío/comentarios)%s\n", colorYellow, filename, colorReset)
			skipped++
			continue
		}

		// Ejecutar el SQL
		fmt.Printf("%s▸ Ejecutando: %s%s\n", color, filename, colorReset)

		_, err = db.Exec(sqlContent)
		if err != nil {
			log.Printf("%s  ✗ Error: %v%s\n", colorRed, err, colorReset)
			continue
		}

		fmt.Printf("%s  ✓ Éxito%s\n", colorGreen, colorReset)
		executed++
	}

	return executed, skipped
}

func isEmptyOrComment(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Si hay una línea que no es vacía ni comentario, el archivo tiene contenido
		if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
			return false
		}
	}
	return true
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
