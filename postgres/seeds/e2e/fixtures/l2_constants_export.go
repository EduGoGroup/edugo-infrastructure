package fixtures

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L2ConstantsExport es una fixture pasiva que valida la presencia de
// las filas L2 (sembradas por system.ApplySystem cuando L2 está
// registrada en system.Layers()) y exporta sus identificadores al
// ApplyContext para que tests downstream y el JSON
// fixtures-constants.json puedan referenciarlos sin hardcodear.
//
// NO escribe filas: L2 vive en el namespace del production seed, y la
// regla del framework prohíbe que las fixtures lo modifiquen
// (ver framework.Fixture docstring, C-REQ-10.2). Por la misma razón
// Manifest no declara Provides, Requires ni Tables: la fixture no
// genera entidades ni participa en el cleanup selectivo por prefijo.
//
// Además de la presencia de filas, la fixture realiza las assertions
// SQL focales de la Fase 4 (cubren F4-REQ-1.1, F4-REQ-1.2, F4-REQ-2.1,
// F4-REQ-3.1, F4-REQ-3.2). Las assertions HTTP/UI (F4-REQ-1.3,
// F4-REQ-3.3, F4-REQ-5.2, F4-REQ-5.3) quedan diferidas — ver docstring
// del scenario L2Form.
//
// MP-09 F4: la no-regresión sobre la cadena L1 viewer→permisos se
// retiró. El usuario viewer era DATO DE TENANT que L1 ya no siembra
// (system/ es contrato puro); el dato vivo equivalente vive en
// playground_v2/base, no en el contrato que estas fixtures validan.
//
// Refs: phase-4-layer-l2/{requirements,design}.md.
type L2ConstantsExport struct{}

// Manifest implementa framework.Fixture.
func (f *L2ConstantsExport) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l2_constants_export",
		Description: "Verifica filas L2 sembradas por system.ApplySystem (ScreenInstance announcement-form + ResourceScreen form) y exporta sus identificadores al JSON.",
		Constants: map[string]string{
			"E2EFixtureL2ScreenInstanceAnnouncementFormID":  layers.L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID,
			"E2EFixtureL2ResourceScreenAnnouncementsFormID": layers.L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID,
			"E2EFixtureL2ScreenKeyAnnouncementForm":         layers.L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
		},
	}
}

// Apply verifica L2 y exporta constantes. Idempotente: sólo lee y
// llama a ctx.SetConstant.
//
// Cubre por SQL:
//   - F4-REQ-1.1: ScreenInstance announcement-form existe con
//     template_id = L0_SCREEN_TPL_FORM_ID.
//   - F4-REQ-1.2: slot_data es JSON válido con 4 fields
//     (title, body, scope, published_at), 3 actions y api_prefix=platform.
//   - F4-REQ-2.1: ResourceScreen (resource=announcements,
//     screen_type=form, is_default=false) existe.
//   - F4-REQ-3.1: action SAVE_NEW lleva permission=announcements:create.
//   - F4-REQ-3.2: action SAVE_EXISTING lleva permission=announcements:update.
func (f *L2ConstantsExport) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	if err := f.verifyScreenInstance(tx); err != nil {
		return err
	}
	if err := f.verifyResourceScreen(tx); err != nil {
		return err
	}

	ctx.SetConstant("E2EFixtureL2ScreenInstanceAnnouncementFormID", layers.L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID)
	ctx.SetConstant("E2EFixtureL2ResourceScreenAnnouncementsFormID", layers.L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID)
	ctx.SetConstant("E2EFixtureL2ScreenKeyAnnouncementForm", layers.L2_SCREEN_KEY_ANNOUNCEMENT_FORM)

	return nil
}

// Cleanup es no-op: las filas L2 son del system seed, no del scenario.
func (f *L2ConstantsExport) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}

// verifyScreenInstance cubre F4-REQ-1.1 y F4-REQ-1.2.
//
// Verifica que existe la fila en ui_config.screen_instances con
// screen_key="announcement-form" y template_id=L0_SCREEN_TPL_FORM_ID,
// y que su slot_data es JSON válido con la forma declarada en
// design.md §3 (4 fields con keys title/body/scope/published_at,
// 3 actions, api_prefix=platform). Las assertions sobre las actions
// SAVE_NEW/SAVE_EXISTING (F4-REQ-3.1/3.2) se hacen en el mismo barrido
// para no reparsear el JSON dos veces.
func (f *L2ConstantsExport) verifyScreenInstance(tx *gorm.DB) error {
	type row struct {
		ID         string
		TemplateID string
		ScreenKey  string
		SlotData   []byte
	}
	const q = `
SELECT id::text AS id,
       template_id::text AS template_id,
       screen_key,
       slot_data
FROM ui_config.screen_instances
WHERE id = ?::uuid
`
	var r row
	if err := tx.Raw(q, layers.L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID).Scan(&r).Error; err != nil {
		return fmt.Errorf("l2_constants_export: query screen_instances: %w", err)
	}
	if r.ID == "" {
		return fmt.Errorf(
			"l2_constants_export: ScreenInstance L2 ausente id=%s — corré system.ApplySystem (con L2 registrada) antes del scenario",
			layers.L2_SCREEN_INSTANCE_ANNOUNCEMENT_FORM_ID,
		)
	}

	// F4-REQ-1.1: template_id = L0_SCREEN_TPL_FORM_ID, screen_key =
	// "announcement-form".
	if r.ScreenKey != layers.L2_SCREEN_KEY_ANNOUNCEMENT_FORM {
		return fmt.Errorf(
			"l2_constants_export: F4-REQ-1.1 violado — screen_key=%q, want %q",
			r.ScreenKey, layers.L2_SCREEN_KEY_ANNOUNCEMENT_FORM,
		)
	}
	if r.TemplateID != layers.L0_SCREEN_TPL_FORM_ID {
		return fmt.Errorf(
			"l2_constants_export: F4-REQ-1.1 violado — template_id=%q, want %q (L0_SCREEN_TPL_FORM_ID)",
			r.TemplateID, layers.L0_SCREEN_TPL_FORM_ID,
		)
	}

	// F4-REQ-1.2: slot_data es JSON válido con la forma esperada.
	type field struct {
		Key string `json:"key"`
	}
	type action struct {
		ID         string `json:"id"`
		Event      string `json:"event"`
		Permission string `json:"permission"`
	}
	var slot struct {
		Fields    []field  `json:"fields"`
		Actions   []action `json:"actions"`
		APIPrefix string   `json:"api_prefix"`
	}
	if err := json.Unmarshal(r.SlotData, &slot); err != nil {
		return fmt.Errorf("l2_constants_export: F4-REQ-1.2 violado — slot_data no es JSON válido: %w", err)
	}

	if got := len(slot.Fields); got != 4 {
		return fmt.Errorf("l2_constants_export: F4-REQ-1.2 violado — fields[]=%d, want 4", got)
	}
	wantFieldKeys := map[string]bool{"title": false, "body": false, "scope": false, "published_at": false}
	for _, fd := range slot.Fields {
		if _, ok := wantFieldKeys[fd.Key]; ok {
			wantFieldKeys[fd.Key] = true
		}
	}
	for k, seen := range wantFieldKeys {
		if !seen {
			return fmt.Errorf("l2_constants_export: F4-REQ-1.2 violado — field key %q ausente en slot_data.fields", k)
		}
	}

	if got := len(slot.Actions); got != 3 {
		return fmt.Errorf("l2_constants_export: F4-REQ-1.2 violado — actions[]=%d, want 3", got)
	}
	if slot.APIPrefix != "platform" {
		return fmt.Errorf("l2_constants_export: F4-REQ-1.2 violado — api_prefix=%q, want %q", slot.APIPrefix, "platform")
	}

	// F4-REQ-3.1 / F4-REQ-3.2: las actions con event SAVE_NEW /
	// SAVE_EXISTING llevan los permisos correctos.
	wantPerms := map[string]string{
		"SAVE_NEW":      "academic.announcements.create",
		"SAVE_EXISTING": "academic.announcements.update",
	}
	seenPerms := map[string]string{}
	for _, a := range slot.Actions {
		if _, expected := wantPerms[a.Event]; expected {
			seenPerms[a.Event] = a.Permission
		}
	}
	for event, want := range wantPerms {
		got, ok := seenPerms[event]
		if !ok {
			return fmt.Errorf("l2_constants_export: F4-REQ-3 violado — action con event=%q ausente en slot_data.actions", event)
		}
		if got != want {
			return fmt.Errorf(
				"l2_constants_export: F4-REQ-3 violado — action event=%q permission=%q, want %q",
				event, got, want,
			)
		}
	}

	return nil
}

// verifyResourceScreen cubre F4-REQ-2.1: existe el mapping
// (resource_id=L0_RESOURCE_ANNOUNCEMENTS_ID, screen_type=form,
// is_default=false) en ui_config.resource_screens.
func (f *L2ConstantsExport) verifyResourceScreen(tx *gorm.DB) error {
	type row struct {
		ID         string
		ResourceID string
		ScreenType string
		IsDefault  bool
	}
	const q = `
SELECT id::text AS id,
       resource_id::text AS resource_id,
       screen_type,
       is_default
FROM ui_config.resource_screens
WHERE resource_id = ?::uuid AND screen_type = ?
`
	var r row
	if err := tx.Raw(q, layers.L0_RESOURCE_ANNOUNCEMENTS_ID, "form").Scan(&r).Error; err != nil {
		return fmt.Errorf("l2_constants_export: query resource_screens: %w", err)
	}
	if r.ID == "" {
		return fmt.Errorf(
			"l2_constants_export: F4-REQ-2.1 violado — resource_screens (resource=%s, type=form) ausente; corré system.ApplySystem con L2 registrada",
			layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
		)
	}
	if r.ID != layers.L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID {
		return fmt.Errorf(
			"l2_constants_export: F4-REQ-2.1 violado — resource_screens.id=%q, want %q (L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID)",
			r.ID, layers.L2_RESOURCE_SCREEN_ANNOUNCEMENTS_FORM_ID,
		)
	}
	if r.IsDefault {
		return fmt.Errorf("l2_constants_export: F4-REQ-2.1 violado — resource_screens.is_default=true, want false (form NO es default)")
	}
	return nil
}
