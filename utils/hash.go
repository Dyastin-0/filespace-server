package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = bcrypt.DefaultCost
)

func GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func CompareHash(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("hashed password and password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
