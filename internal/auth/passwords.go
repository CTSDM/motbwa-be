package auth

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost_hashing, err := strconv.Atoi(os.Getenv("COST_HASHING"))
	if err != nil {
		log.Printf("Error while casting to int the cost_hashing: %v", err)
		cost_hashing = DEFAULT_COST_HASHING
	}
	if cost_hashing == 0 {
		log.Printf("Something went wrong while loading the cost_hashing: %v", err)
		cost_hashing = DEFAULT_COST_HASHING
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost_hashing)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
