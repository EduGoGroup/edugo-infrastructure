package fixtures

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L1ConstantsExport es una fixture pasiva que valida la presencia del
// rol de contrato L1 (announcement_viewer, sembrado por
// system.ApplySystem cuando L1 está registrada en system.Layers()) y
// exporta su identificador al ApplyContext para que tests downstream y
// el JSON fixtures-constants.json puedan referenciarlo sin hardcodear.
//
// NO escribe filas: L1 vive en el namespace del production seed, y la
// regla del framework prohíbe que las fixtures lo modifiquen
// (ver framework.Fixture docstring, C-REQ-10.2). Por la misma razón
// Manifest no declara Provides, Requires ni Tables: la fixture no
// genera entidades ni participa en el cleanup selectivo por prefijo.
//
// MP-09 F4: L1 dejó de sembrar DATO DE TENANT (escuela demo, usuario
// viewer, user_role, membership). system/ es CONTRATO PURO: la fixture
// sólo valida/exporta el rol de contrato; el dato vivo equivalente vive
// en playground_v2/base.
//
// Refs: phase-3-layer-l1/{requirements,design}.md (F3-REQ-1..6),
// ADR-7.
type L1ConstantsExport struct{}

// Manifest implementa framework.Fixture.
func (f *L1ConstantsExport) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l1_constants_export",
		Description: "Verifica el rol de contrato L1 (announcement_viewer) sembrado por system.ApplySystem y exporta su identificador al JSON.",
		Constants: map[string]string{
			"E2EFixtureL1RoleAnnouncementViewerID":   layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
			"E2EFixtureL1RoleAnnouncementViewerName": layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME,
		},
	}
}

// Apply verifica el rol de contrato L1 y exporta sus constantes.
// Idempotente: sólo lee y llama a ctx.SetConstant.
func (f *L1ConstantsExport) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	var count int64
	if err := tx.Raw(
		`SELECT COUNT(*) FROM iam.roles WHERE id = ?`,
		layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
	).Scan(&count).Error; err != nil {
		return fmt.Errorf("l1_constants_export: verify iam.roles: %w", err)
	}
	if count == 0 {
		return fmt.Errorf(
			"l1_constants_export: rol L1 ausente en iam.roles id=%s — corré system.ApplySystem (con L1 registrada) antes del scenario",
			layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
		)
	}

	ctx.SetConstant("E2EFixtureL1RoleAnnouncementViewerID", layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID)
	ctx.SetConstant("E2EFixtureL1RoleAnnouncementViewerName", layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME)

	return nil
}

// Cleanup es no-op: el rol L1 es del system seed, no del scenario.
func (f *L1ConstantsExport) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}
