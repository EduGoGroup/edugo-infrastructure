package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/EduGoGroup/edugo-infrastructure/tools/mock-generator/pkg/generator"
	"github.com/EduGoGroup/edugo-infrastructure/tools/mock-generator/pkg/parser"
)

var (
	testingDir string
	outputDir  string
)

var rootCmd = &cobra.Command{
	Use:   "mock-generator",
	Short: "Genera codigo Go desde scripts SQL de testing",
	Long:  "Parser de SQL que genera dataset mock para desarrollo frontend",
	Run:   runGenerator,
}

func init() {
	rootCmd.Flags().StringVar(&testingDir, "testing", "../../postgres/migrations/testing", "Directorio con SQL de testing")
	rootCmd.Flags().StringVar(&outputDir, "output", "../../mock/dataset", "Directorio de salida")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runGenerator(cmd *cobra.Command, args []string) {
	fmt.Println("Mock Generator v1.0.0")
	fmt.Printf("Testing dir: %s\n", testingDir)
	fmt.Printf("Output dir: %s\n", outputDir)
	fmt.Println("")

	// Parser
	p := parser.NewSQLParser()
	fmt.Println("Parseando archivos SQL...")
	tables, err := p.ParseDirectory(testingDir)
	if err != nil {
		fmt.Printf("Error parseando: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parseados %d tablas\n\n", len(tables))

	// Mostrar estadísticas
	fmt.Println("Estadisticas por tabla:")
	for tableName, data := range tables {
		fmt.Printf("  - %-20s: %d registros, %d columnas\n",
			tableName, len(data.Rows), len(data.Columns))
	}

	// Generador
	fmt.Println("\nGenerando dataset...")
	gen := generator.NewDatasetGenerator(outputDir, tables)
	if err := gen.Generate(); err != nil {
		fmt.Printf("Error generando: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Dataset generado exitosamente")
	fmt.Printf("Archivos en: %s\n", outputDir)
}
