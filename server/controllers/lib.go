package controllers

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func generateApiKey() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		log.Error().AnErr("error", err).Msg("While filling randomBytes varible with random bytes")
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func hashApiKey(apiKey string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(apiKey), 10)

	return string(bytes), err
}

func checkHashedApiKey(hashedApiKey, apiKey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedApiKey), []byte(apiKey))

	return err == nil
}
