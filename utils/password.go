package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword Generate Hash Password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("field to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassowrd(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
