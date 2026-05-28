package common

import "github.com/google/uuid"

// MustParseUUID convierte s a uuid.UUID. Panic si el string es malformado.
// Pensado para constantes UUID de seeds.
func MustParseUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}
