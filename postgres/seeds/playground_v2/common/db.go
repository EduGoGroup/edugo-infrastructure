package common

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// onConflictIgnore inserta value ignorando conflictos por la PK detectada por
// GORM. Atajo del patrón repetido tx.Clauses(clause.OnConflict{DoNothing:true}).
func onConflictIgnore(tx *gorm.DB, value any) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(value).Error
}

// MustParseUUID convierte s a uuid.UUID; panic si está malformado. Pensado para
// constantes UUID de seeds (donde un id malformado es un bug del seed, no un
// error recuperable en runtime).
func MustParseUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}
