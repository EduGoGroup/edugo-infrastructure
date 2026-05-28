package layers

import (
	"encoding/json"
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

// L3ScreenInstances retorna las 2 instancias `materials-list` y
// `material-form` sembradas por L3. Mirror determinístico de
// applyL3Screens en l3_screens.go.
func L3ScreenInstances() ([]entities.ScreenInstance, error) {
	listID, err := uuid.Parse(L3_SCREEN_INSTANCE_MATERIALS_LIST_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ScreenInstances: parse list id: %w", err)
	}
	formID, err := uuid.Parse(L3_SCREEN_INSTANCE_MATERIAL_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ScreenInstances: parse form id: %w", err)
	}
	listTemplateID, err := uuid.Parse(L0_SCREEN_TPL_LIST_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ScreenInstances: parse list template id: %w", err)
	}
	formTemplateID, err := uuid.Parse(L0_SCREEN_TPL_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ScreenInstances: parse form template id: %w", err)
	}
	descList := "Listado de materiales educativos"
	descForm := "Formulario de creación/edición de materiales"
	requiredPermissionList := "content.materials.read"
	requiredPermissionForm := "content.materials.read"
	return []entities.ScreenInstance{
		{
			ID:                 listID,
			ScreenKey:          L3_SCREEN_KEY_MATERIALS_LIST,
			TemplateID:         listTemplateID,
			Name:               "Listado de materiales",
			Description:        &descList,
			SlotData:           json.RawMessage([]byte(materialsListSlotData)),
			Scope:              "unit",
			RequiredPermission: &requiredPermissionList,
			HandlerKey:         nil,
			IsActive:           true,
		},
		{
			ID:                 formID,
			ScreenKey:          L3_SCREEN_KEY_MATERIAL_FORM,
			TemplateID:         formTemplateID,
			Name:               "Formulario de material",
			Description:        &descForm,
			SlotData:           json.RawMessage([]byte(materialFormSlotData)),
			Scope:              "unit",
			RequiredPermission: &requiredPermissionForm,
			HandlerKey:         nil,
			IsActive:           true,
		},
	}, nil
}

// L3ResourceScreens retorna los 2 mappings materials ↔ {materials-list,
// material-form} sembrados por L3. Mirror determinístico de
// applyL3ResourceScreens en l3_resource_screens.go.
func L3ResourceScreens() ([]entities.ResourceScreen, error) {
	materialsID, err := uuid.Parse(L3_RESOURCE_MATERIALS_ID)
	if err != nil {
		return nil, fmt.Errorf("L3ResourceScreens: parse resource_id: %w", err)
	}
	idList := uuid.NewSHA1(uuid.NameSpaceOID, []byte(materialsID.String()+":list"))
	idForm := uuid.NewSHA1(uuid.NameSpaceOID, []byte(materialsID.String()+":form"))
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
		{
			ID:          idForm,
			ResourceID:  materialsID,
			ResourceKey: L3_RESOURCE_MATERIALS_KEY,
			ScreenKey:   L3_SCREEN_KEY_MATERIAL_FORM,
			ScreenType:  "form",
			IsDefault:   false,
			SortOrder:   1,
			IsActive:    true,
		},
	}, nil
}
