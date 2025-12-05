package common

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// The functions below are merely wrapper functions of bcrypt

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), fmt.Errorf("failed hashing the password %w", err)
}

func CheckHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
