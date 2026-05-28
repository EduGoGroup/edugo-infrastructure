package common

import "golang.org/x/crypto/bcrypt"

// BcryptHash devuelve el hash bcrypt de password. Panic si bcrypt falla
// (este es código de seed, no runtime de producción).
func BcryptHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("seeds/playground/common: bcrypt failed: " + err.Error())
	}
	return string(hash)
}
