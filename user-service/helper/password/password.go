package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func GenerateRandomSalt(size int) (string, error) {
	salt := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPasswordWithSalt(password, salt string) string {
	combined := password + salt
	hash := sha256.Sum256([]byte(combined))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func ComparePasswordWithHash(password, salt, hashed string) bool {
	expected := HashPasswordWithSalt(password, salt)
	return expected == hashed
}
