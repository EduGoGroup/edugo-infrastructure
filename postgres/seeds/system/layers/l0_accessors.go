package layers

import (
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Accessors públicos de L0 — espejan el patrón de
// `system/l4/accessors.go` para que herramientas estáticas
// (cross-checker, seedaudit) puedan materializar los datos sembrados
// por L0 sin tocar la base.
//
// Política: cada accessor reconstruye in-memory las MISMAS entidades
// que el `applyL0*` correspondiente persiste en Postgres (mismas IDs,
// keys, scopes, etc.). NO se exponen entidades fuera del scope del
// cross-checker — `users` y `user_roles` quedan deliberadamente sin
// accessor porque el seedaudit no las consume.
//
// Resuelve TC-5: previamente el loader sólo conocía L4 y reportaba
// las entidades de L0..L3 como phantom/unused (39 falsos positivos).

// L0Resources retorna el recurso `announcements` sembrado por L0
// (1 fila). Mirror determinístico de applyL0Resources en l0_resources.go.
func L0Resources() ([]entities.Resource, error) {
	id, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return nil, fmt.Errorf("L0Resources: parse id: %w", err)
	}
	description := "Comunicaciones y anuncios institucionales"
	icon := "bullhorn"
	return []entities.Resource{
		{
			ID:            id,
			Key:           L0_RESOURCE_ANNOUNCEMENTS_KEY,
			DisplayName:   "Anuncios",
			Description:   &description,
			Icon:          &icon,
			ParentID:      nil,
			SortOrder:     0,
			IsMenuVisible: true,
			Scope:         "school",
			IsActive:      true,
		},
	}, nil
}

// L0Roles retorna el rol `super_admin` sembrado por L0 (1 fila).
// Mirror determinístico de upsertL0Role en l0_roles.go.
func L0Roles() ([]entities.Role, error) {
	id, err := uuid.Parse(L0_ROLE_SUPER_ADMIN_ID)
	if err != nil {
		return nil, fmt.Errorf("L0Roles: parse id: %w", err)
	}
	desc := "Rol con acceso total al sistema. Usuario de bootstrapping."
	return []entities.Role{
		{
			ID:          id,
			Name:        L0_ROLE_SUPER_ADMIN_NAME,
			DisplayName: "Super Administrador",
			Description: &desc,
			Scope:       "system",
			IsActive:    true,
		},
	}, nil
}

// L0Permissions retorna los 4 permisos CRUD sobre announcements
// sembrados por L0. Mirror determinístico de upsertL0Permissions.
func L0Permissions() ([]entities.Permission, error) {
	resourceID, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return nil, fmt.Errorf("L0Permissions: parse resource_id: %w", err)
	}
	specs := []struct {
		idStr       string
		name        string
		displayName string
		action      string
	}{
		{L0_PERM_ANNOUNCEMENTS_READ, "academic.announcements.read", "Ver Anuncios", "read"},
		{L0_PERM_ANNOUNCEMENTS_CREATE, "academic.announcements.create", "Crear Anuncio", "create"},
		{L0_PERM_ANNOUNCEMENTS_UPDATE, "academic.announcements.update", "Editar Anuncio", "update"},
		{L0_PERM_ANNOUNCEMENTS_DELETE, "academic.announcements.delete", "Eliminar Anuncio", "delete"},
	}
	out := make([]entities.Permission, 0, len(specs))
	for _, s := range specs {
		id, err := uuid.Parse(s.idStr)
		if err != nil {
			return nil, fmt.Errorf("L0Permissions: parse id %s: %w", s.idStr, err)
		}
		out = append(out, entities.Permission{
			ID:          id,
			Name:        s.name,
			DisplayName: s.displayName,
			ResourceID:  resourceID,
			Action:      s.action,
			Scope:       "school",
			IsActive:    true,
		})
	}
	return out, nil
}

// P4-1 (plan B): L0RolePermissions() fue eliminada. La tabla
// iam.role_permissions ya no existe; los permisos efectivos del
// super_admin se otorgan vía iam.role_grants con pattern `*` desde L4.

// L0ScreenTemplates retorna las 4 ScreenTemplates compartidas (list,
// detail, form, master-detail) v1 sembradas por L0. Mirror
// determinístico de upsertL0ScreenTemplates en l0_screens.go.
func L0ScreenTemplates() ([]entities.ScreenTemplate, error) {
	listID, err := uuid.Parse(L0_SCREEN_TPL_LIST_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenTemplates: parse list id: %w", err)
	}
	detailID, err := uuid.Parse(L0_SCREEN_TPL_DETAIL_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenTemplates: parse detail id: %w", err)
	}
	formID, err := uuid.Parse(L0_SCREEN_TPL_FORM_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenTemplates: parse form id: %w", err)
	}
	masterDetailID, err := uuid.Parse(L0_SCREEN_TPL_MASTER_DETAIL_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenTemplates: parse master-detail id: %w", err)
	}
	return []entities.ScreenTemplate{
		{
			ID:         listID,
			Pattern:    "list",
			Name:       "list-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(listBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         detailID,
			Pattern:    "detail",
			Name:       "detail-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(detailBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         formID,
			Pattern:    "form",
			Name:       "form-basic-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(formBasicV1Definition)),
			IsActive:   true,
		},
		{
			ID:         masterDetailID,
			Pattern:    "master-detail",
			Name:       "master-detail-v1",
			Version:    1,
			Definition: json.RawMessage([]byte(masterDetailV1Definition)),
			IsActive:   true,
		},
	}, nil
}

// L0ScreenInstances retorna la instancia `announcements-list`
// sembrada por L0 (1 fila). Mirror de upsertL0ScreenInstances.
func L0ScreenInstances() ([]entities.ScreenInstance, error) {
	id, err := uuid.Parse(L0_SCREEN_INST_ANNOUNCEMENTS_LIST_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenInstances: parse id: %w", err)
	}
	templateID, err := uuid.Parse(L0_SCREEN_TPL_LIST_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ScreenInstances: parse template id: %w", err)
	}
	description := "Listado de anuncios institucionales"
	requiredPermission := "academic.announcements.read"
	return []entities.ScreenInstance{
		{
			ID:                 id,
			ScreenKey:          L0_SCREEN_KEY_ANNOUNCEMENTS_LIST,
			TemplateID:         templateID,
			Name:               "Anuncios — Listado",
			Description:        &description,
			SlotData:           json.RawMessage([]byte(announcementsListSlotData)),
			Scope:              "school",
			RequiredPermission: &requiredPermission,
			HandlerKey:         nil,
			IsActive:           true,
		},
	}, nil
}

// L0ResourceScreens retorna el mapping announcements ↔ announcements-list
// sembrado por L0 (1 fila). Mirror de upsertL0ResourceScreens.
func L0ResourceScreens() ([]entities.ResourceScreen, error) {
	resourceID, err := uuid.Parse(L0_RESOURCE_ANNOUNCEMENTS_ID)
	if err != nil {
		return nil, fmt.Errorf("L0ResourceScreens: parse resource_id: %w", err)
	}
	screenType := "list"
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(resourceID.String()+":"+screenType))
	return []entities.ResourceScreen{
		{
			ID:          id,
			ResourceID:  resourceID,
			ResourceKey: L0_RESOURCE_ANNOUNCEMENTS_KEY,
			ScreenKey:   L0_SCREEN_KEY_ANNOUNCEMENTS_LIST,
			ScreenType:  screenType,
			IsDefault:   true,
			SortOrder:   0,
			IsActive:    true,
		},
	}, nil
}
