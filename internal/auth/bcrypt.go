package auth

import "golang.org/x/crypto/bcrypt"

func HashingPass(password string) string {
	hashpass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashpass)
}

func CheckHashPass(password, hashpass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password)) == nil
}
