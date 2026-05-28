// Package focal_colors_demo es el playground 6 (Fase 3 — SDUI demo
// data-driven sin Kotlin). Valida que agregar un CRUD plano nuevo a
// EduGo NO requiere código Kotlin: las pantallas `colors-list` /
// `colors-form` se resuelven via `ScreenContractRegistry.findOrCreate`
// con fallback a `GenericListContract` / `GenericFormContract`,
// parametrizadas con la metadata SDUI declarada en los `slot_data` de
// las dos screen_instances sembradas por L4 (B7b).
//
// **Composicion (autosuficiente)**: depende de [admin] para tener
// `super_admin` con grant `*` y un tenant minimo (1 escuela + 1 unidad
// + admin@edugo.local). Apply() encadena admin.Apply() al inicio para
// soportar `P=focal-colors-demo` standalone. Con `P=all` la dependencia
// ya corrio antes en el registry — la doble aplicacion es no-op gracias
// a OnConflict DoNothing. Idempotente.
//
// **Roles agregados** (wildcard-first segun feedback_wildcard_first):
//   - colors-author -> `platform.colors.*` (CRUD completo)
//
// Se evita reutilizar el rol super_admin (que ya tiene `*` y por tanto
// cubre platform.colors.*); el rol `colors-author` existe explicitamente
// para que el smoke test pueda demostrar que UN GRANT ESPECIFICO basta
// para abrir el flow generico — sin depender de super_admin.
//
// **Usuarios agregados** (al tenant heredado de `admin`):
//   - colors-author@edugo.local / 12345678
//
// **Datos de seed** (academic.colors):
//   6 colores de muestra en la escuela Playground Admin para que la
//   lista no se vea vacia al abrirla. Los nombres son humanos (paleta
//   editorial) y los hex se validan contra el CHECK constraint
//   `^#[0-9A-Fa-f]{6}$` (post_gorm.sql).
//
// Rango UUID 65000000-... para usuarios/memberships/roles/colors
// (no colisiona con 60..., 61..., 62..., 63... de otros playgrounds).
//
// Idempotente: INSERTs usan OnConflict DoNothing.
package focal_colors_demo

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/admin"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales del playground.
	AuthorEmail = "colors-author@edugo.local"
	Password    = "12345678"

	// Tenant heredado de `admin` (school + unit ya sembrados ahi).
	tenantSchoolID = "60000000-0000-0000-0000-000000000002"

	// Rango 65000000-... — usuario + membership + rol del playground.
	authorUserID = "65000000-0000-0000-0000-000000000001"
	authorMembID = "65000000-0000-0000-0000-000000000011"

	authorRoleID   = "15000000-0000-0000-0000-000000000001"
	authorRoleName = "colors-author"
)

// Apply siembra el playground focal_colors_demo sobre admin. Encadena
// admin.Apply al inicio para ser autosuficiente cuando se invoca con
// `P=focal-colors-demo` standalone. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := admin.Apply(tx); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: dependencia admin: %w", err)
	}
	// Rol del playground (wildcard-first: un pattern cubre cualquier permiso futuro).
	if err := common.SeedRole(tx, common.RoleSpec{
		ID:          common.MustParseUUID(authorRoleID),
		Name:        authorRoleName,
		DisplayName: "Colors — Author",
		Description: "CRUD completo sobre platform.colors. Playground demo SDUI data-driven (Fase 3).",
	}); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: role: %w", err)
	}
	if err := common.SeedRoleGrant(tx, common.MustParseUUID(authorRoleID), "platform.colors.*"); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: role_grants: %w", err)
	}
	if err := common.SeedUser(tx, common.UserSpec{
		ID:        common.MustParseUUID(authorUserID),
		Email:     AuthorEmail,
		Password:  Password,
		FirstName: "Colors",
		LastName:  "Author",
	}); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: user: %w", err)
	}
	if err := common.SeedUserRole(tx, common.MustParseUUID(authorUserID), common.MustParseUUID(authorRoleID)); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: user_role: %w", err)
	}
	if err := common.SeedMembership(tx, common.MembershipSpec{
		ID:       common.MustParseUUID(authorMembID),
		UserID:   common.MustParseUUID(authorUserID),
		SchoolID: common.MustParseUUID(tenantSchoolID),
		Role:     "admin",
	}); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: membership: %w", err)
	}
	if err := upsertSampleColors(tx); err != nil {
		return fmt.Errorf("playground/focal_colors_demo: colors: %w", err)
	}
	return nil
}

// upsertSampleColors siembra 6 colores en la escuela del playground
// admin para que `colors-list` no se vea vacia al abrirla por primera
// vez. Hex en formato `#RRGGBB` (mayusculas) para validar el CHECK
// constraint regex de migrations/sql/post_gorm.sql.
func upsertSampleColors(tx *gorm.DB) error {
	sid := common.MustParseUUID(tenantSchoolID)
	items := []entities.Color{
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000021"), SchoolID: sid, Name: "Rojo carmin", Hex: "#C0392B"},
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000022"), SchoolID: sid, Name: "Naranja calido", Hex: "#E67E22"},
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000023"), SchoolID: sid, Name: "Amarillo mostaza", Hex: "#F1C40F"},
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000024"), SchoolID: sid, Name: "Verde esmeralda", Hex: "#2ECC71"},
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000025"), SchoolID: sid, Name: "Azul oceano", Hex: "#2980B9"},
		{ID: common.MustParseUUID("65000000-0000-0000-0000-000000000026"), SchoolID: sid, Name: "Violeta profundo", Hex: "#8E44AD"},
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&items).Error
}
