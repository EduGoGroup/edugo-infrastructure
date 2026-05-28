package layers

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// applyL3Resources siembra el recurso "materials" de L3.
// Usa SQL crudo siguiendo el patrón de applyL0Resources para
// garantizar correctitud de columnas booleanas con default tag
// (workaround para bug GORM bool zero-value).
// Idempotente vía ON CONFLICT (id).
//
// Scope `unit`: distinto de announcements (`school`). L3 valida que
// el sistema soporta recursos con scope distintos (F5-REQ-1.1). El
// enum iam.permission_scope admite los valores 'system', 'school' y
// 'unit' (ver migrations/sql/pre_gorm.sql).
func applyL3Resources(tx *gorm.DB) error {
	const upsertSQL = `
        INSERT INTO iam.resources
            (id, key, display_name, description, icon, parent_id, sort_order, is_menu_visible, scope, is_active, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?::iam.permission_scope, ?, NOW(), NOW())
        ON CONFLICT (id) DO UPDATE SET
            key             = EXCLUDED.key,
            display_name    = EXCLUDED.display_name,
            description     = EXCLUDED.description,
            icon            = EXCLUDED.icon,
            parent_id       = EXCLUDED.parent_id,
            sort_order      = EXCLUDED.sort_order,
            is_menu_visible = EXCLUDED.is_menu_visible,
            scope           = EXCLUDED.scope,
            is_active       = EXCLUDED.is_active
    `

	id, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return fmt.Errorf("applyL3Resources: parse id: %w", err)
	}

	description := "Materiales educativos"
	icon := "book"

	if err := tx.Exec(upsertSQL,
		id,
		L3_RESOURCE_MATERIALS_KEY,
		"Materiales",
		&description,
		&icon,
		nil,    // parent_id raíz
		1,      // sort_order (después de announcements=0)
		true,   // is_menu_visible
		"unit", // scope (distinto de announcements=school)
		true,   // is_active
	).Error; err != nil {
		return fmt.Errorf("applyL3Resources: upsert materials: %w", err)
	}
	return nil
}
