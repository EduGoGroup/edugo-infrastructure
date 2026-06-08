package layers

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Accessors públicos de L3 — espejan el patrón de
// `system/l4/accessors.go`. L3 introduce el segundo recurso del
// sistema (`materials`) con CRUD parcial y dos pantallas propias.

// L3Resources retorna el recurso `materials` sembrado por L3 (1 fila).
// Mirror determinístico de applyL3Resources en l3_resources.go.
func L3Resources() ([]entities.Resource, error) {
	id, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return nil, fmt.Errorf("L3Resources: parse id: %w", err)
	}
	description := "Materiales educativos"
	icon := "book"
	return []entities.Resource{
		{
			ID:            id,
			Key:           L3_RESOURCE_MATERIALS_KEY,
			DisplayName:   "Materiales",
			Description:   &description,
			Icon:          &icon,
			ParentID:      nil,
			SortOrder:     1,
			IsMenuVisible: true,
			Scope:         "unit",
			IsActive:      true,
		},
	}, nil
}

// L3Permissions retorna los 3 permisos CRUD parcial sobre materials
// (read, create, update — sin delete) sembrados por L3.
// Mirror determinístico de applyL3Permissions en l3_permissions.go.
func L3Permissions() ([]entities.Permission, error) {
	resourceID, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return nil, fmt.Errorf("L3Permissions: parse resource_id: %w", err)
	}
	specs := []struct {
		idStr       string
		name        string
		displayName string
		action      string
	}{
		{L3_PERM_MATERIALS_READ_ID, "content.materials.read", "Ver Materiales", "read"},
		{L3_PERM_MATERIALS_CREATE_ID, "content.materials.create", "Crear Material", "create"},
		{L3_PERM_MATERIALS_UPDATE_ID, "content.materials.update", "Editar Material", "update"},
	}
	out := make([]entities.Permission, 0, len(specs))
	for _, s := range specs {
		id, err := uuid.Parse(s.idStr)
		if err != nil {
			return nil, fmt.Errorf("L3Permissions: parse id %s: %w", s.idStr, err)
		}
		out = append(out, entities.Permission{
			ID:          id,
			Name:        s.name,
			DisplayName: s.displayName,
			ResourceID:  resourceID,
			Action:      s.action,
			Scope:       "unit",
			IsActive:    true,
		})
	}
	return out, nil
}

// P4-1 (plan B): L3RolePermissions() fue eliminada. La tabla
// iam.role_permissions ya no existe; los permisos materials.* del
// super_admin se otorgan vía iam.role_grants (pattern `*`) desde L4.

// L3ScreenInstances retorna las ScreenInstances sembradas por L3.
// Mirror determinístico de l3_screens.go.
//
// Poda SDUI material (2026-06-07): L3 ya NO siembra screen_instances.
// Las pantallas `materials-list` / `material-form` eran código muerto
// (la app las renderiza nativas, no por SDUI) y fueron eliminadas. El
// recurso materials sigue vivo en el menú vía el mapping resource_screen
// `materials:list` (SIN ScreenInstance) — ver L3ResourceScreens. Por eso
// este mirror devuelve un slice vacío.
func L3ScreenInstances() ([]entities.ScreenInstance, error) {
	return []entities.ScreenInstance{}, nil
}

// L3ResourceScreens retorna el mapping materials ↔ materials-list
// sembrado por L3. Mirror determinístico de applyL3ResourceScreens en
// l3_resource_screens.go.
//
// El screen_key `materials-list` apunta a una pantalla NATIVA (sin
// ScreenInstance). El mapping `material-form` fue podado junto con su
// ScreenInstance (poda SDUI material 2026-06-07).
func L3ResourceScreens() ([]entities.ResourceScreen, error) {
	materialsID, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ResourceScreens: parse resource_id: %w", err)
	}
	idList := uuid.NewSHA1(uuid.NameSpaceOID, []byte(materialsID.String()+":list"))
	return []entities.ResourceScreen{
		{
			ID:          idList,
			ResourceID:  materialsID,
			ResourceKey: L3_RESOURCE_MATERIALS_KEY,
			ScreenKey:   L3_SCREEN_KEY_MATERIALS_LIST,
			ScreenType:  "list",
			IsDefault:   true,
			SortOrder:   0,
			IsActive:    true,
		},
	}, nil
}
