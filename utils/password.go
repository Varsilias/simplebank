package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
)

var (
	lowercase = regexp.MustCompile(`[a-z]`)
	uppercase = regexp.MustCompile(`[A-Z]`)
	digit     = regexp.MustCompile(`[0-9]`)
	special   = regexp.MustCompile(`[@$!%*?&]`)
)

const (
	SaltSize = 32 // 32 bytes = 256 bits
)

// PasswordData contains the hashed password and its salt
type PasswordData struct {
	HashedPassword string
	Salt           string
}

// GenerateRandomSalt generates a random salt of specified size
func GenerateRandomSalt() (string, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("error generating salt: %w", err)
	}

	// Convert to base64 for storage
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword hashes the password with the given salt
func HashPassword(password string) (*PasswordData, error) {
	// Generate new salt
	salt, err := GenerateRandomSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash password with salt
	hashedPassword, err := hashPasswordWithSalt(password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &PasswordData{
		HashedPassword: hashedPassword,
		Salt:           salt,
	}, nil
}

// hashPasswordWithSalt hashes the password with a given salt
func hashPasswordWithSalt(password, salt string) (string, error) {
	// Decode salt from base64
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	// Create hasher
	hasher := sha512.New()

	// First write password
	hasher.Write([]byte(password))

	// Then write salt
	hasher.Write(saltBytes)

	// Get the final hash
	hashedBytes := hasher.Sum(nil)

	// Convert to hex string for storage
	return hex.EncodeToString(hashedBytes), nil
}

// VerifyPassword checks if an input password matches the stored hash
func VerifyPassword(inputPassword, storedHash, storedSalt string) (bool, error) {
	calculatedHash, err := hashPasswordWithSalt(inputPassword, storedSalt)
	if err != nil {
		return false, fmt.Errorf("failed to hash input password: %w", err)
	}

	return calculatedHash == storedHash, nil
}

// IsValidPassword verifies that the input password is at least 8 characters long,
// has a lowercase letter, an uppercase letter, a number and a special character
func IsValidPassword(password string) bool {
	return len(password) >= 8 &&
		lowercase.MatchString(password) &&
		uppercase.MatchString(password) &&
		digit.MatchString(password) &&
		special.MatchString(password)
}
