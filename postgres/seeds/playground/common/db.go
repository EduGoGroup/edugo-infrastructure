package common

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OnConflictIgnore inserta value ignorando conflictos por las columnas del
// índice único primario detectado por GORM. Atajo del patrón repetido
// tx.Clauses(clause.OnConflict{DoNothing: true}).Create(value).
//
// Si necesitas especificar columnas concretas (no la PK), usa la cláusula
// completa directamente en el caller.
func OnConflictIgnore(tx *gorm.DB, value any) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(value).Error
}
