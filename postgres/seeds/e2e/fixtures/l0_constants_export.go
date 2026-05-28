package fixtures

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
)

// L0ConstantsExport es una fixture pasiva que valida la presencia de
// las filas L0 (sembradas por system.ApplySystem) y exporta sus
// identificadores al ApplyContext para que tests downstream y el JSON
// fixtures-constants.json puedan referenciarlos sin hardcodear.
//
// NO escribe filas: L0 vive en el namespace del production seed, y la
// regla del framework prohíbe que las fixtures lo modifiquen
// (ver framework.Fixture docstring, C-REQ-10.2). Por la misma razón
// Manifest no declara Provides, Requires ni Tables: la fixture no
// genera entidades ni participa en el cleanup selectivo por prefijo.
type L0ConstantsExport struct{}

// Manifest implementa framework.Fixture.
func (f *L0ConstantsExport) Manifest() framework.FixtureManifest {
	return framework.FixtureManifest{
		Name:        "l0_constants_export",
		Description: "Verifica filas L0 sembradas por system.ApplySystem y exporta sus identificadores al JSON.",
		Constants: map[string]string{
			"E2EFixtureL0ResourceAnnouncementsID":       layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
			"E2EFixtureL0ResourceAnnouncementsKey":      layers.L0_RESOURCE_ANNOUNCEMENTS_KEY,
			"E2EFixtureL0RoleSuperAdminID":              layers.L0_ROLE_SUPER_ADMIN_ID,
			"E2EFixtureL0RoleSuperAdminName":            layers.L0_ROLE_SUPER_ADMIN_NAME,
			"E2EFixtureL0PermAnnouncementsRead":         layers.L0_PERM_ANNOUNCEMENTS_READ,
			"E2EFixtureL0PermAnnouncementsCreate":       layers.L0_PERM_ANNOUNCEMENTS_CREATE,
			"E2EFixtureL0PermAnnouncementsUpdate":       layers.L0_PERM_ANNOUNCEMENTS_UPDATE,
			"E2EFixtureL0PermAnnouncementsDelete":       layers.L0_PERM_ANNOUNCEMENTS_DELETE,
			"E2EFixtureL0ScreenInstAnnouncementsListID": layers.L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID,
			"E2EFixtureL0ScreenKeyAnnouncementsList":    layers.L0_SCREEN_KEY_ANNOUNCEMENTS_LIST,
			"E2EFixtureL0UserSuperAdminID":              layers.L0_USER_SUPER_ADMIN_ID,
			"E2EFixtureL0UserSuperAdminEmail":           layers.L0_SUPER_ADMIN_EMAIL,
			"E2EFixtureL0UserSuperAdminPassword":        layers.L0_SUPER_ADMIN_PASSWORD,
		},
	}
}

// Apply verifica L0 y exporta constantes. Idempotente: sólo lee y
// llama a ctx.SetConstant.
func (f *L0ConstantsExport) Apply(tx *gorm.DB, ctx *framework.ApplyContext) error {
	checks := []struct {
		table string
		sql   string
		id    string
	}{
		{
			table: "iam.resources",
			sql:   `SELECT COUNT(*) FROM iam.resources WHERE id = ?`,
			id:    layers.L0_RESOURCE_ANNOUNCEMENTS_ID,
		},
		{
			table: "iam.roles",
			sql:   `SELECT COUNT(*) FROM iam.roles WHERE id = ?`,
			id:    layers.L0_ROLE_SUPER_ADMIN_ID,
		},
		{
			table: "ui_config.screen_instances",
			sql:   `SELECT COUNT(*) FROM ui_config.screen_instances WHERE id = ?`,
			id:    layers.L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID,
		},
		{
			table: "auth.users",
			sql:   `SELECT COUNT(*) FROM auth.users WHERE id = ?`,
			id:    layers.L0_USER_SUPER_ADMIN_ID,
		},
	}
	for _, c := range checks {
		var count int64
		if err := tx.Raw(c.sql, c.id).Scan(&count).Error; err != nil {
			return fmt.Errorf("l0_constants_export: verify %s: %w", c.table, err)
		}
		if count == 0 {
			return fmt.Errorf(
				"l0_constants_export: fila L0 ausente en %s id=%s — corré system.ApplySystem antes del scenario",
				c.table, c.id,
			)
		}
	}

	ctx.SetConstant("E2EFixtureL0ResourceAnnouncementsID", layers.L0_RESOURCE_ANNOUNCEMENTS_ID)
	ctx.SetConstant("E2EFixtureL0ResourceAnnouncementsKey", layers.L0_RESOURCE_ANNOUNCEMENTS_KEY)
	ctx.SetConstant("E2EFixtureL0RoleSuperAdminID", layers.L0_ROLE_SUPER_ADMIN_ID)
	ctx.SetConstant("E2EFixtureL0RoleSuperAdminName", layers.L0_ROLE_SUPER_ADMIN_NAME)
	ctx.SetConstant("E2EFixtureL0PermAnnouncementsRead", layers.L0_PERM_ANNOUNCEMENTS_READ)
	ctx.SetConstant("E2EFixtureL0PermAnnouncementsCreate", layers.L0_PERM_ANNOUNCEMENTS_CREATE)
	ctx.SetConstant("E2EFixtureL0PermAnnouncementsUpdate", layers.L0_PERM_ANNOUNCEMENTS_UPDATE)
	ctx.SetConstant("E2EFixtureL0PermAnnouncementsDelete", layers.L0_PERM_ANNOUNCEMENTS_DELETE)
	ctx.SetConstant("E2EFixtureL0ScreenInstAnnouncementsListID", layers.L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID)
	ctx.SetConstant("E2EFixtureL0ScreenKeyAnnouncementsList", layers.L0_SCREEN_KEY_ANNOUNCEMENTS_LIST)
	ctx.SetConstant("E2EFixtureL0UserSuperAdminID", layers.L0_USER_SUPER_ADMIN_ID)
	ctx.SetConstant("E2EFixtureL0UserSuperAdminEmail", layers.L0_SUPER_ADMIN_EMAIL)
	ctx.SetConstant("E2EFixtureL0UserSuperAdminPassword", layers.L0_SUPER_ADMIN_PASSWORD)

	return nil
}

// Cleanup es no-op: las filas L0 son del system seed, no del scenario.
func (f *L0ConstantsExport) Cleanup(tx *gorm.DB, ctx *framework.ApplyContext) error {
	return nil
}
