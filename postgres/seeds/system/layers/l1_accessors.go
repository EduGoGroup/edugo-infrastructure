package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Accessors públicos de L1 — espejan el patrón de
// `system/l4/accessors.go`. Sólo se exponen las entidades que el
// cross-checker (seedaudit) necesita para validar el catálogo iam:
// roles y role_permissions. Entidades fuera de scope (schools,
// memberships, users, user_roles) NO se exponen por design (TC-5).

// L1Roles retorna el rol `announcement_viewer` sembrado por L1
// (1 fila). Mirror determinístico de applyL1Role en l1_role.go.
func L1Roles() ([]entities.Role, error) {
	id, err := uuid.Parse(L1_ROLE_ANNOUNCEMENT_VIEWER_ID)
	if err != nil {
		return nil, fmt.Errorf("L1Roles: parse id: %w", err)
	}
	desc := "Rol read-only: solo puede ver anuncios. Usado para validar gating de UI en Fase 3."
	return []entities.Role{
		{
			ID:          id,
			Name:        L1_ROLE_ANNOUNCEMENT_VIEWER_NAME,
			DisplayName: "Visualizador de Anuncios",
			Description: &desc,
			Scope:       "school",
			IsActive:    true,
			IsSystem:    true,
		},
	}, nil
}

// P4-1 (plan B): L1RolePermissions() fue eliminada. La tabla
// iam.role_permissions ya no existe; el permiso `academic.announcements.read`
// para el rol announcement_viewer se otorga vía iam.role_grants desde L4.
