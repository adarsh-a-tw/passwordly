package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(pwd string) string
	ComparePassword(rawPwd string, hashedPwd string) bool
}

type PasswordHasherImpl struct{}

func (ph *PasswordHasherImpl) HashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Hashing failed")
	}
	return string(hash)
}

func (ph *PasswordHasherImpl) ComparePassword(rawPwd string, hashedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(rawPwd))
	return err == nil
}
