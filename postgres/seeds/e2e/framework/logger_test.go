package framework

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestJSONLogger_EmitOneLine(t *testing.T) {
	var buf bytes.Buffer
	log := NewJSONLoggerTo(&buf)
	log.now = func() time.Time { return time.Unix(1700000000, 0).UTC() }
	log.Emit(LogEntry{Event: EventFixtureApply, Fixture: "role_only"})
	out := buf.String()
	if !strings.HasSuffix(out, "\n") {
		t.Fatalf("salida sin newline: %q", out)
	}
	if strings.Count(out, "\n") != 1 {
		t.Fatalf("se esperaba 1 línea; got=%q", out)
	}
	var parsed LogEntry
	if err := json.Unmarshal([]byte(strings.TrimRight(out, "\n")), &parsed); err != nil {
		t.Fatalf("no parsea como JSON: %v", err)
	}
	if parsed.Event != EventFixtureApply {
		t.Errorf("Event=%v", parsed.Event)
	}
	if parsed.Fixture != "role_only" {
		t.Errorf("Fixture=%q", parsed.Fixture)
	}
}

func TestMemoryLogger_CaptureAndReset(t *testing.T) {
	log := NewMemoryLogger()
	log.Emit(LogEntry{Event: EventFixtureApply, Fixture: "a"})
	log.Emit(LogEntry{Event: EventFixtureCleanup, Fixture: "a"})
	got := log.Captured()
	if len(got) != 2 {
		t.Fatalf("Captured len=%d, want 2", len(got))
	}
	if got[0].Time.IsZero() {
		t.Error("Time debería poblarse automáticamente")
	}
	log.Reset()
	if len(log.Captured()) != 0 {
		t.Error("Reset no limpió")
	}
}

func TestNopLogger(t *testing.T) {
	log := NewNopLogger()
	// No debe panic.
	log.Emit(LogEntry{Event: EventFixtureApply})
}

func TestJSONLogger_Concurrent(t *testing.T) {
	var buf bytes.Buffer
	log := NewJSONLoggerTo(&buf)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 10; j++ {
				log.Emit(LogEntry{Event: EventFixtureApply, Fixture: "x"})
			}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if got := strings.Count(buf.String(), "\n"); got != 100 {
		t.Errorf("got=%d líneas, want 100", got)
	}
}
