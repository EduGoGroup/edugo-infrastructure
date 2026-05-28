package fixtures

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L1ConstantsExport es una fixture pasiva que valida la presencia de
// las filas L1 (sembradas por system.ApplySystem cuando L1 está
// registrada en system.Layers()) y exporta sus identificadores al
// ApplyContext para que tests downstream y el JSON
// fixtures-constants.json puedan referenciarlos sin hardcodear.
//
// NO escribe filas: L1 vive en el namespace del production seed, y la
// regla del framework prohíbe que las fixtures lo modifiquen
// (ver framework.Fixture docstring, C-REQ-10.2). Por la misma razón
// Manifest no declara Provides, Requires ni Tables: la fixture no
// genera entidades ni participa en el cleanup selectivo por prefijo.
//
// Refs: phase-3-layer-l1/{requirements,design}.md (F3-REQ-1..6),
// ADR-7 (L1 toca academic.schools para cumplir scope=school).
type L1ConstantsExport struct{}

// Manifest implementa framework.Fixture.
func (f *L1ConstantsExport) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l1_constants_export",
		Description: "Verifica filas L1 sembradas por system.ApplySystem y exporta sus identificadores al JSON.",
		Constants: map[string]string{
			"E2EFixtureL1RoleAnnouncementViewerID":   layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
			"E2EFixtureL1RoleAnnouncementViewerName": layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME,
			"E2EFixtureL1UserViewerID":               layers.L1_USER_VIEWER_ID,
			"E2EFixtureL1UserViewerEmail":            layers.L1_VIEWER_EMAIL,
			"E2EFixtureL1UserViewerPassword":         layers.L1_VIEWER_PASSWORD,
			"E2EFixtureL1SchoolDemoID":               layers.L1_SCHOOL_DEMO_ID,
			"E2EFixtureL1SchoolDemoCode":             layers.L1_SCHOOL_DEMO_CODE,
			"E2EFixtureL1SchoolDemoName":             layers.L1_SCHOOL_DEMO_NAME,
			"E2EFixtureL1UserRoleViewerID":           layers.L1_USER_ROLE_VIEWER_ID,
		},
	}
}

// Apply verifica L1 y exporta constantes. Idempotente: sólo lee y
// llama a ctx.SetConstant.
func (f *L1ConstantsExport) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	checks := []struct {
		table string
		sql   string
		id    string
	}{
		{
			table: "academic.schools",
			sql:   `SELECT COUNT(*) FROM academic.schools WHERE id = ?`,
			id:    layers.L1_SCHOOL_DEMO_ID,
		},
		{
			table: "iam.roles",
			sql:   `SELECT COUNT(*) FROM iam.roles WHERE id = ?`,
			id:    layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID,
		},
		{
			table: "auth.users",
			sql:   `SELECT COUNT(*) FROM auth.users WHERE id = ?`,
			id:    layers.L1_USER_VIEWER_ID,
		},
		{
			table: "iam.user_roles",
			sql:   `SELECT COUNT(*) FROM iam.user_roles WHERE id = ?`,
			id:    layers.L1_USER_ROLE_VIEWER_ID,
		},
	}
	for _, c := range checks {
		var count int64
		if err := tx.Raw(c.sql, c.id).Scan(&count).Error; err != nil {
			return fmt.Errorf("l1_constants_export: verify %s: %w", c.table, err)
		}
		if count == 0 {
			return fmt.Errorf(
				"l1_constants_export: fila L1 ausente en %s id=%s — corré system.ApplySystem (con L1 registrada) antes del scenario",
				c.table, c.id,
			)
		}
	}

	ctx.SetConstant("E2EFixtureL1RoleAnnouncementViewerID", layers.L1_ROLE_ANNOUNCEMENT_VIEWER_ID)
	ctx.SetConstant("E2EFixtureL1RoleAnnouncementViewerName", layers.L1_ROLE_ANNOUNCEMENT_VIEWER_NAME)
	ctx.SetConstant("E2EFixtureL1UserViewerID", layers.L1_USER_VIEWER_ID)
	ctx.SetConstant("E2EFixtureL1UserViewerEmail", layers.L1_VIEWER_EMAIL)
	ctx.SetConstant("E2EFixtureL1UserViewerPassword", layers.L1_VIEWER_PASSWORD)
	ctx.SetConstant("E2EFixtureL1SchoolDemoID", layers.L1_SCHOOL_DEMO_ID)
	ctx.SetConstant("E2EFixtureL1SchoolDemoCode", layers.L1_SCHOOL_DEMO_CODE)
	ctx.SetConstant("E2EFixtureL1SchoolDemoName", layers.L1_SCHOOL_DEMO_NAME)
	ctx.SetConstant("E2EFixtureL1UserRoleViewerID", layers.L1_USER_ROLE_VIEWER_ID)

	return nil
}

// Cleanup es no-op: las filas L1 son del system seed, no del scenario.
func (f *L1ConstantsExport) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}
