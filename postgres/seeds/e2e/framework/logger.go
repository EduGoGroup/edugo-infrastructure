package framework

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogEvent es el conjunto de eventos estructurados que el framework
// emite (C-REQ-8). Las fixtures NO los emiten directamente: el
// composer los emite alrededor de cada Apply/Cleanup.
type LogEvent string

const (
	EventFixtureApply   LogEvent = "fixture.apply"
	EventFixtureCleanup LogEvent = "fixture.cleanup"
	EventFixtureError   LogEvent = "fixture.error"
	EventScenarioApply  LogEvent = "scenario.apply"
	EventScenarioDone   LogEvent = "scenario.done"
)

// LogEntry es el contenido estructurado de un log line. Los campos
// vacíos se omiten al serializar.
type LogEntry struct {
	Time         time.Time `json:"time"`
	Event        LogEvent  `json:"event"`
	Scenario     string    `json:"scenario,omitempty"`
	Fixture      string    `json:"fixture,omitempty"`
	TenantPrefix string    `json:"tenant_prefix,omitempty"`
	Tables       []string  `json:"tables,omitempty"`
	RowsInserted int64     `json:"rows_inserted,omitempty"`
	RowsUpdated  int64     `json:"rows_updated,omitempty"`
	RowsDeleted  int64     `json:"rows_deleted,omitempty"`
	DurationMs   int64     `json:"duration_ms,omitempty"`
	Stage        string    `json:"stage,omitempty"`
	LastTable    string    `json:"last_table,omitempty"`
	Error        string    `json:"error,omitempty"`
}

// Logger abstrae la salida del framework para que los tests puedan
// capturar eventos sin tocar stdout (C-REQ-8.4).
type Logger interface {
	Emit(entry LogEntry)
}

// JSONLogger serializa cada entry como una línea JSON sobre el writer
// configurado (por defecto stdout).
type JSONLogger struct {
	mu  sync.Mutex
	w   io.Writer
	now func() time.Time
}

// NewJSONLogger devuelve un logger que escribe a stdout.
func NewJSONLogger() *JSONLogger {
	return &JSONLogger{w: os.Stdout, now: time.Now}
}

// NewJSONLoggerTo construye un logger con writer arbitrario (útil
// cuando el binario quiere redirigir a un fichero o syslog).
func NewJSONLoggerTo(w io.Writer) *JSONLogger {
	return &JSONLogger{w: w, now: time.Now}
}

// Emit serializa la entry como JSON-line.
func (l *JSONLogger) Emit(entry LogEntry) {
	if entry.Time.IsZero() {
		entry.Time = l.now()
	}
	buf, err := json.Marshal(entry)
	if err != nil {
		// Fallback: nunca debería ocurrir con los tipos planos de LogEntry.
		_, _ = fmt.Fprintf(l.w, `{"time":%q,"event":%q,"error":"json marshal failed: %v"}`+"\n",
			entry.Time.Format(time.RFC3339Nano), entry.Event, err)
		return
	}
	l.mu.Lock()
	_, _ = l.w.Write(buf)
	_, _ = l.w.Write([]byte{'\n'})
	l.mu.Unlock()
}

// MemoryLogger acumula los eventos en memoria; los tests lo consultan
// con Captured() para hacer asserts.
type MemoryLogger struct {
	mu      sync.Mutex
	entries []LogEntry
	now     func() time.Time
}

// NewMemoryLogger construye un logger en memoria con reloj real.
func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{now: time.Now}
}

// Emit agrega la entry a la lista interna.
func (l *MemoryLogger) Emit(entry LogEntry) {
	if entry.Time.IsZero() {
		entry.Time = l.now()
	}
	l.mu.Lock()
	l.entries = append(l.entries, entry)
	l.mu.Unlock()
}

// Captured devuelve una copia de los eventos capturados hasta el
// momento.
func (l *MemoryLogger) Captured() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LogEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Reset descarta los eventos capturados (útil entre subtests).
func (l *MemoryLogger) Reset() {
	l.mu.Lock()
	l.entries = nil
	l.mu.Unlock()
}

// nopLogger descarta los eventos. Útil para tests que no inspeccionan
// la salida pero necesitan un Logger no-nil.
type nopLogger struct{}

// NewNopLogger devuelve un logger que descarta todos los eventos.
func NewNopLogger() Logger { return nopLogger{} }

// Emit es un no-op.
func (nopLogger) Emit(LogEntry) {}
