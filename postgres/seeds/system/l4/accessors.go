package l4

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Accessors públicos de L4 — espejan los de `[archivado pre-Fase-6] accessors.go`
// para que herramientas estáticas (cross-checker, seedaudit) puedan
// materializar los datos sembrados por L4 sin tocar la base.
//
// Patrón: cada función reusa el mismo helper de construcción que usa la
// función `Apply*` del dominio. NO se duplica la lógica — si Apply
// inserta N filas, la lista que retorna el accessor son las mismas N
// filas (mismas IDs, mismos valores), modulo timestamps generados en
// runtime.
//
// Las funciones devuelven `([]T, error)` para alinearse con la firma
// pública del paquete `legacy`; en la práctica L4 sólo puede fallar si
// alguna constante UUID está corrupta (mismo panic-equivalente que ya
// haría el Apply en runtime).

// Resources materializa l4Resources como entities.Resource.
func Resources() ([]entities.Resource, error) {
	out := make([]entities.Resource, 0, len(l4Resources))
	for _, r := range l4Resources {
		id, err := uuid.Parse(r.ID)
		if err != nil {
			return nil, fmt.Errorf("l4.Resources: parse id %s: %w", r.ID, err)
		}
		var parentID *uuid.UUID
		if r.ParentID != "" {
			pid, perr := uuid.Parse(r.ParentID)
			if perr != nil {
				return nil, fmt.Errorf("l4.Resources: parse parent_id %s for %s: %w", r.ParentID, r.Key, perr)
			}
			parentID = &pid
		}
		var description *string
		if r.Description != "" {
			d := r.Description
			description = &d
		}
		var icon *string
		if r.Icon != "" {
			i := r.Icon
			icon = &i
		}
		out = append(out, entities.Resource{
			ID:            id,
			Key:           r.Key,
			DisplayName:   r.DisplayName,
			Description:   description,
			Icon:          icon,
			ParentID:      parentID,
			SortOrder:     r.SortOrder,
			IsMenuVisible: r.IsMenuVisible,
			Scope:         r.Scope,
			IsActive:      r.IsActive,
		})
	}
	return out, nil
}

// Roles retorna los 5 roles sembrados por L4 (student, teacher,
// guardian, admin, school_admin).
func Roles() ([]entities.Role, error) {
	return buildL4Roles()
}

// Permissions retorna las permissions sembradas por L4 (catálogo
// completo definido en l4Permissions()).
func Permissions() ([]entities.Permission, error) {
	return buildL4Permissions()
}

// P4-1 (plan B): RolePermissions() fue eliminado. La tabla
// iam.role_permissions ya no existe; los permisos efectivos por rol se
// resuelven via iam.role_grants (patterns wildcard).

// RoleGrants retorna los role_grants generados por applyL4RoleGrants.
// A diferencia de los demás accessors, NO se construye en memoria desde
// specs declarativas — la fuente de verdad es la BD. Útil para
// cross-checker.
func RoleGrants(tx *gorm.DB) ([]entities.RoleGrant, error) {
	var rgs []entities.RoleGrant
	if err := tx.Find(&rgs).Error; err != nil {
		return nil, err
	}
	return rgs, nil
}

// ScreenTemplates retorna los 6 templates adicionales sembrados por L4.
func ScreenTemplates() ([]entities.ScreenTemplate, error) {
	return buildL4ScreenTemplates(), nil
}

// ScreenInstances retorna las screen instances sembradas por L4.
func ScreenInstances() ([]entities.ScreenInstance, error) {
	return buildL4ScreenInstances(), nil
}

// ResourceScreens retorna los mappings recurso↔pantalla sembrados por L4.
func ResourceScreens() ([]entities.ResourceScreen, error) {
	return buildL4ResourceScreens(), nil
}

// ConceptTypes retorna los 5 tipos de concepto sembrados por L4.
func ConceptTypes() ([]entities.ConceptType, error) {
	return buildL4ConceptTypes(), nil
}

// ConceptDefinitions retorna las 50 definiciones (10 × 5 tipos)
// sembradas por L4. Como el schema asigna el ID con
// `gen_random_uuid()`, el accessor sintetiza un UUID determinístico
// derivado de (concept_type_id, term_key) para que llamadas repetidas
// devuelvan slices estables.
func ConceptDefinitions() ([]entities.ConceptDefinition, error) {
	rows, err := buildL4ConceptDefinitions()
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].ID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(rows[i].ConceptTypeID.String()+":"+rows[i].TermKey))
	}
	return rows, nil
}
