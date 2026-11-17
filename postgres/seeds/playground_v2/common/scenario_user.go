package common

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScenarioUser describe un usuario que un escenario on-top necesita (plan 053
// T5.2). El escenario NO lo crea: lo crea la API de identidad (con su hash de
// contraseña) y el escenario solo lo busca por email en auth.users.
type ScenarioUser struct {
	Email     string
	FirstName string
	LastName  string
}

// UserIDs mapea email → user_id resuelto en auth.users.
type UserIDs map[string]uuid.UUID

// Require comprueba que ids trae un id para cada usuario de users. Devuelve un
// error que lista, ordenados, TODOS los emails ausentes.
func (ids UserIDs) Require(users []ScenarioUser) error {
	var missing []string
	for _, u := range users {
		if _, ok := ids[u.Email]; !ok {
			missing = append(missing, u.Email)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("faltan %d usuarios en auth.users: %s", len(missing), strings.Join(missing, ", "))
	}
	return nil
}

// ResolveUserIDs resuelve los emails de users a su id en auth.users (solo
// filas vivas: deleted_at IS NULL) en una sola consulta. Devuelve los ids
// encontrados y, ordenados, los emails que no existen. No escribe nada.
//
// Mismo cuidado que catalog.ResolveInvitationTypeID: el id se escanea como
// texto y se parsea, porque uuid.UUID ([16]byte) no se escanea bien directo
// desde el driver.
func ResolveUserIDs(tx *gorm.DB, users []ScenarioUser) (UserIDs, []string, error) {
	emails := make([]string, 0, len(users))
	for _, u := range users {
		emails = append(emails, u.Email)
	}
	ids := make(UserIDs, len(users))
	if len(emails) == 0 {
		return ids, nil, nil
	}

	var rows []struct {
		ID    string
		Email string
	}
	if err := tx.Raw(
		`SELECT id::text AS id, email FROM auth.users WHERE email IN ? AND deleted_at IS NULL`,
		emails,
	).Scan(&rows).Error; err != nil {
		return nil, nil, fmt.Errorf("ResolveUserIDs: consultando auth.users: %w", err)
	}
	for _, r := range rows {
		id, err := uuid.Parse(r.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("ResolveUserIDs: id inválido para %q: %w", r.Email, err)
		}
		ids[r.Email] = id
	}

	var missing []string
	for _, e := range emails {
		if _, ok := ids[e]; !ok {
			missing = append(missing, e)
		}
	}
	sort.Strings(missing)
	return ids, missing, nil
}
