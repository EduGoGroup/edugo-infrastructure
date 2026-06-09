package l4

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestL4ScreenInstancesSlotDataIsValidJSON valida en compile/test
// time que cada slot_data literal embebido en los row builders es
// JSON parseable. Si el test falla, el seed romperia el migrator
// porque `entities.ScreenInstance.SlotData` es jsonb y postgres
// rechazara la fila al UPSERT.
//
// B4: este test queda como artefacto del bloque para que regresiones
// futuras (alguien tipea mal una `,` en el JSON literal) se detecten
// antes de correr el migrador.
func TestL4ScreenInstancesSlotDataIsValidJSON(t *testing.T) {
	rows := l4ScreenInstanceRows()
	if len(rows) == 0 {
		t.Fatal("l4ScreenInstanceRows() retornó vacío — B4 no sembraría nada")
	}
	seen := make(map[string]bool, len(rows))
	for _, r := range rows {
		if seen[r.screenKey] {
			t.Errorf("screen_key duplicado en B4: %q", r.screenKey)
		}
		seen[r.screenKey] = true

		var v any
		if err := json.Unmarshal([]byte(r.slotData), &v); err != nil {
			t.Errorf("slot_data inválido para %s: %v", r.screenKey, err)
		}
		if r.id == "" || r.templateID == "" || r.name == "" || r.scope == "" {
			t.Errorf("campos obligatorios vacíos en %s", r.screenKey)
		}
	}
}

// TestL4EntityPickerConformanceFixture fija (golden in-code) el contrato
// EXACTO del control `entity-picker` que hoy expone el seed productivo en
// `assessments-form.subject_id` (screen_instances_rows.go:assessmentsForm).
//
// Por qué un golden:
//
//	El contrato SDUI es la frontera estable entre el backend y CUALQUIER
//	front (KMP hoy, `apple_new` SwiftUI mañana — ADR 0007/0021, plan 016 §7
//	"Compatibilidad hacia adelante"). Un segundo front logra paridad sin
//	leer el código KMP validando contra fixtures JSON golden de las
//	respuestas reales del contrato. Este test versiona ese golden para el
//	control `entity-picker` y rompe si el seed cambia los hints sin que el
//	contrato (y los fronts que lo consumen) se actualicen a conciencia.
//
// Alcance (plan 019, WI-5 + WI-4): SOLO `entity-picker`. La "zona de
// búsqueda" NO se modela (WI-4 2026-06-09: hints `searchCollapsible` /
// `searchPersistentOn` se dejan implícitos, no entran al contrato).
//
// El golden refleja los hints REALES del seed: `search_param`, `page_size`,
// `picker_title` (y `remote_endpoint`/`display_field`/`value_field`). NO hay
// `min_chars` en este control — el test lo afirma explícitamente para
// documentar la ausencia y detectar si alguien lo añade sin actualizar la
// conformidad.
func TestL4EntityPickerConformanceFixture(t *testing.T) {
	const (
		screenKey = "assessments-form"
		fieldKey  = "subject_id"
	)

	// goldenEntityPicker: el shape EXACTO esperado del campo `subject_id`.
	// Si el seed productivo cambia un hint, este map debe actualizarse en
	// el MISMO cambio (regla "no muerto") y los fronts re-validarse.
	goldenEntityPicker := map[string]any{
		"key":             "subject_id",
		"label":           "Materia",
		"type":            "entity-picker",
		"required":        true,
		"remote_endpoint": "academic:/api/v1/subjects",
		"display_field":   "name",
		"value_field":     "id",
		"search_param":    "search",
		"page_size":       float64(20), // JSON numbers → float64
		"picker_title":    "Buscar materia",
	}

	// Localiza el row del seed productivo y parsea su slot_data.
	var row *l4ScreenInstanceRow
	for _, r := range l4ScreenInstanceRows() {
		if r.screenKey == screenKey {
			rr := r
			row = &rr
			break
		}
	}
	if row == nil {
		t.Fatalf("no se encontró el screen_instance %q en el seed productivo", screenKey)
	}

	var slot struct {
		Fields []map[string]any `json:"fields"`
	}
	if err := json.Unmarshal([]byte(row.slotData), &slot); err != nil {
		t.Fatalf("slot_data de %q no parsea: %v", screenKey, err)
	}

	// Encuentra el campo entity-picker bajo prueba.
	var field map[string]any
	for _, f := range slot.Fields {
		if k, _ := f["key"].(string); k == fieldKey {
			field = f
			break
		}
	}
	if field == nil {
		t.Fatalf("no se encontró el field %q en %q", fieldKey, screenKey)
	}

	// Pre-condición: debe seguir siendo un entity-picker (si migró de tipo,
	// este fixture ya no aplica y debe revisarse el alcance de WI-5).
	if got, _ := field["type"].(string); got != "entity-picker" {
		t.Fatalf("%s.%s ya no es entity-picker (type=%q); el fixture de conformidad quedó obsoleto", screenKey, fieldKey, got)
	}

	// Golden assertion: el campo del seed === el contrato esperado, exacto.
	if !reflect.DeepEqual(field, goldenEntityPicker) {
		t.Errorf("contrato entity-picker de %s.%s difiere del golden de conformidad.\n  esperado: %#v\n  obtenido: %#v",
			screenKey, fieldKey, goldenEntityPicker, field)
	}

	// Afirma la AUSENCIA de `min_chars` en este control (documenta el shape
	// real: el typeahead no declara umbral mínimo de caracteres en el seed).
	if _, present := field["min_chars"]; present {
		t.Errorf("%s.%s ahora declara `min_chars`; actualiza el golden de conformidad y los fronts que validan contra él", screenKey, fieldKey)
	}
}
