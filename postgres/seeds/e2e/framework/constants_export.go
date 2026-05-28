package framework

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ConstantsExportSchemaVersion identifica el shape del JSON exportado
// para que los consumidores (tests Kotlin del KMP) detecten drift.
//
// Bumpear este valor cuando cambie el shape del JSON (campos nuevos,
// renombrados o quitados). Los consumidores comparan el valor leído
// con el esperado y fallan con `fixtures-constants schemaVersion
// mismatch: expected X, got Y` (C-REQ-6.5).
const ConstantsExportSchemaVersion = "1"

// DefaultExportPath es la ruta canónica del JSON dentro del repo.
// Relativa al working directory del binario seed_e2e; en el monorepo
// se resuelve a `postgres/seeds/e2e/exports/fixtures-constants.json`.
const DefaultExportPath = "seeds/e2e/exports/fixtures-constants.json"

// ConstantsExport es la representación serializable del JSON.
type ConstantsExport struct {
	SchemaVersion string                              `json:"schemaVersion"`
	GeneratedAt   time.Time                           `json:"generatedAt"`
	Scenarios     map[string]ConstantsScenarioExport  `json:"scenarios"`
}

// ConstantsScenarioExport agrupa los datos derivados de un scenario.
type ConstantsScenarioExport struct {
	TenantPrefix string            `json:"tenantPrefix"`
	SchemaPrefix string            `json:"schemaPrefix"`
	Constants    map[string]string `json:"constants"`
}

// ConstantsExporter merge-actualiza el JSON exportado con la información
// de un ApplyContext recién terminado. Es seguro entre goroutines pero
// en la práctica el binario seed_e2e lo invoca secuencialmente al final
// de cada Apply (C-REQ-6.3).
type ConstantsExporter struct {
	mu   sync.Mutex
	Path string
	Now  func() time.Time
}

// NewConstantsExporter construye un exporter que escribe al path dado.
// Si path es vacío usa DefaultExportPath.
func NewConstantsExporter(path string) *ConstantsExporter {
	if path == "" {
		path = DefaultExportPath
	}
	return &ConstantsExporter{Path: path, Now: time.Now}
}

// WriteFromContext lee el JSON existente (si lo hay), inserta/actualiza
// la entrada del scenario y reescribe el archivo de forma atómica.
//
// La operación es idempotente: invocar dos veces con el mismo
// ApplyContext produce el mismo archivo.
func (e *ConstantsExporter) WriteFromContext(ctx *ApplyContext) error {
	if ctx == nil {
		return fmt.Errorf("constants_export: nil ApplyContext")
	}
	if ctx.ScenarioName == "" {
		return fmt.Errorf("constants_export: ApplyContext.ScenarioName vacío")
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	export, err := e.read()
	if err != nil {
		return err
	}
	if export.Scenarios == nil {
		export.Scenarios = map[string]ConstantsScenarioExport{}
	}
	export.Scenarios[ctx.ScenarioName] = ConstantsScenarioExport{
		TenantPrefix: ctx.TenantPrefix,
		SchemaPrefix: ctx.SchemaPrefix,
		Constants:    cloneStringMap(ctx.Constants),
	}
	export.SchemaVersion = ConstantsExportSchemaVersion
	export.GeneratedAt = e.now()
	return e.writeAtomic(export)
}

// Read devuelve la copia del JSON tal cual está en disco. Útil para
// los tests Kotlin que quieran inspeccionarlo desde Go (no es la API
// estable de consumo — esa es leer el archivo directamente).
func (e *ConstantsExporter) Read() (ConstantsExport, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.read()
}

// read es el helper interno; asume que el lock ya está tomado.
func (e *ConstantsExporter) read() (ConstantsExport, error) {
	data, err := os.ReadFile(e.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return ConstantsExport{
				SchemaVersion: ConstantsExportSchemaVersion,
				Scenarios:     map[string]ConstantsScenarioExport{},
			}, nil
		}
		return ConstantsExport{}, fmt.Errorf("constants_export: read %s: %w", e.Path, err)
	}
	var out ConstantsExport
	if err := json.Unmarshal(data, &out); err != nil {
		return ConstantsExport{}, fmt.Errorf("constants_export: parse %s: %w", e.Path, err)
	}
	if out.Scenarios == nil {
		out.Scenarios = map[string]ConstantsScenarioExport{}
	}
	return out, nil
}

// writeAtomic serializa el export con orden estable (claves
// alfabéticas) y reescribe el archivo usando rename para evitar
// corrupción si el proceso muere a mitad de escritura.
func (e *ConstantsExporter) writeAtomic(export ConstantsExport) error {
	if err := os.MkdirAll(filepath.Dir(e.Path), 0o755); err != nil {
		return fmt.Errorf("constants_export: mkdir %s: %w", filepath.Dir(e.Path), err)
	}
	// Serializar manualmente para producir un orden estable de claves.
	buf, err := marshalStable(export)
	if err != nil {
		return fmt.Errorf("constants_export: marshal: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(e.Path), ".fixtures-constants.*.json.tmp")
	if err != nil {
		return fmt.Errorf("constants_export: temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(buf); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("constants_export: write tmp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("constants_export: close tmp: %w", err)
	}
	if err := os.Rename(tmpName, e.Path); err != nil {
		return fmt.Errorf("constants_export: rename %s -> %s: %w", tmpName, e.Path, err)
	}
	return nil
}

// marshalStable serializa el export con claves ordenadas para que el
// archivo sea estable entre ejecuciones (diff-friendly).
func marshalStable(export ConstantsExport) ([]byte, error) {
	type stableEntry struct {
		Name  string                  `json:"-"`
		Entry ConstantsScenarioExport `json:"entry"`
	}
	// Clonar el map a un slice ordenado, y luego construir un map
	// nuevo cuyo iteration order... no, el JSON encoder en Go ordena
	// keys por defecto en encoding/json desde 1.12. Para estar seguros
	// y hacer un marshal con indentación, dependemos de json.Marshal
	// que ordena claves de map[string]X alfabéticamente.
	// Sin embargo los Constants internos también se ordenan.
	cloned := ConstantsExport{
		SchemaVersion: export.SchemaVersion,
		GeneratedAt:   export.GeneratedAt.UTC(),
		Scenarios:     make(map[string]ConstantsScenarioExport, len(export.Scenarios)),
	}
	for name, sc := range export.Scenarios {
		cloned.Scenarios[name] = ConstantsScenarioExport{
			TenantPrefix: sc.TenantPrefix,
			SchemaPrefix: sc.SchemaPrefix,
			Constants:    cloneStringMap(sc.Constants),
		}
	}
	// json.MarshalIndent ordena map keys alfabéticamente desde Go 1.12+,
	// suficiente para un archivo estable.
	buf, err := json.MarshalIndent(cloned, "", "  ")
	if err != nil {
		return nil, err
	}
	buf = append(buf, '\n')
	return buf, nil
}

// cloneStringMap devuelve una copia profunda de un map[string]string.
// Si el origen es nil devuelve un map vacío para que el JSON serialice
// `{}` en vez de `null`.
func cloneStringMap(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	maps.Copy(out, src)
	return out
}

// SortedKeys helper exportado para que tests/snapshots imprimen las
// claves en orden estable.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// now devuelve el reloj configurado o time.Now por defecto.
func (e *ConstantsExporter) now() time.Time {
	if e.Now != nil {
		return e.Now()
	}
	return time.Now().UTC()
}
