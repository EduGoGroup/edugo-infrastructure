package layers

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// applyL0Resources siembra el recurso "announcements" de L0.
// Usa SQL crudo para garantizar correctitud de columnas booleanas
// con default tag (workaround para bug GORM bool zero-value).
// Idempotente vía ON CONFLICT (id).
func applyL0Resources(tx *gorm.DB) error {
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

	id, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return fmt.Errorf("applyL0Resources: parse id: %w", err)
	}

	description := "Comunicaciones y anuncios institucionales"
	icon := "bullhorn"

	if err := tx.Exec(upsertSQL,
		id,
		L0_RESOURCE_ANNOUNCEMENTS_KEY,
		"Anuncios",
		&description,
		&icon,
		nil,  // parent_id raíz
		0,    // sort_order
		true, // is_menu_visible
		"school",
		true, // is_active
	).Error; err != nil {
		return fmt.Errorf("applyL0Resources: upsert announcements: %w", err)
	}
	return nil
}
