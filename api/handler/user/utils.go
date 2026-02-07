package ehandleruser

import "golang.org/x/crypto/bcrypt"

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}