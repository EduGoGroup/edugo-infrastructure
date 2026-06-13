// Package catalog ofrece resolvers data-driven sobre catálogos del esquema:
// dado un valor de negocio (p.ej. la key de un tipo de invitación) devuelve su
// id, sin hardcodear UUIDs en los seeds que lo consumen. Espeja cómo los seeds
// resuelven roles/recursos por su nombre semántico en vez de por literal.
package catalog

import (
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResolveInvitationTypeID resuelve la key de un tipo de invitación
// (teacher|student|guardian|coordinator|admin|assistant) a su id en
// academic.invitation_types.
//
// Data-driven: consulta la tabla por key, NO hardcodea el UUID. Los seeds que
// siembran membresías/invitaciones lo usan para poblar invitation_type_id (FK
// por id) sin depender del literal del catálogo.
//
// PRECONDICIÓN: academic.invitation_types debe estar sembrado antes de invocar
// (lo siembra la capa L4 vía l4.ApplyInvitationTypes, y L1 lo adelanta para su
// propia membresía). Si la key no existe, devuelve un error claro en vez de un
// uuid cero silencioso.
func ResolveInvitationTypeID(tx *gorm.DB, key string) (uuid.UUID, error) {
	// Pluck a string (no a uuid.UUID): el driver entrega el uuid como texto y
	// uuid.UUID es [16]byte, así que escanear directo lo descompone byte-a-byte
	// sin pasar por su Scanner ("converting driver.Value type string to uint8").
	// Resolvemos a string y parseamos explícitamente.
	var idStr string
	err := tx.Model(&entities.InvitationType{}).
		Where("key = ?", key).
		Limit(1).
		Pluck("id", &idStr).Error
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog.ResolveInvitationTypeID: query key %q: %w", key, err)
	}
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("catalog.ResolveInvitationTypeID: invitation_type key %q no existe en academic.invitation_types (¿se sembró ApplyInvitationTypes antes?)", key)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("catalog.ResolveInvitationTypeID: id inválido para key %q: %w", key, err)
	}
	return id, nil
}
