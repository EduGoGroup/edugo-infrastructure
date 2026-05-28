package l4

import (
	"encoding/json"
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
