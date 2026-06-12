package common

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// BcryptHash devuelve el hash bcrypt de password. Panic si bcrypt falla (código
// de seed, no runtime de producción).
func BcryptHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("seeds/playground_v2/common: bcrypt failed: " + err.Error())
	}
	return string(hash)
}
