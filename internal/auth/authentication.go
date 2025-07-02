package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pwd), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
