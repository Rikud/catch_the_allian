package encoder

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return string(hash)
}

func ComparePasswords(dbHash string, password string) bool {
	hashedPassword := []byte(dbHash)
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}