package gaia

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/pbkdf2"
)

type PasswdEncode struct{}

// validPassword checks if the password matches the required pattern.
func (PasswdEncode) validPassword(password string) (string, error) {
	re := regexp.MustCompile(`^(?=.*[a-zA-Z])(?=.*\d).{8,}$`)
	if re.MatchString(password) {
		return password, nil
	}
	return "", errors.New("密码必须包含字母和数字，且长度必须大于8位")
}

// hashPassword hashes the password with the given salt using PBKDF2 and SHA-256.
func hashPassword(passwordStr string, salt []byte) string {
	dk := pbkdf2.Key([]byte(passwordStr), salt, 10000, sha256.Size, sha256.New)
	return hex.EncodeToString(dk)
}

// ComparePassword compares the given password with the stored hashed password.
func (PasswdEncode) ComparePassword(passwordStr, passwordHashedBase64, saltBase64 string) (bool, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false, err
	}
	hashedPassword := hashPassword(passwordStr, salt)

	expectedHash, err := base64.StdEncoding.DecodeString(passwordHashedBase64)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString([]byte(hashedPassword)) == hex.EncodeToString(expectedHash), nil
}

// EncodePassword generates a salt, hashes the password, and encodes both using Base64.
func (PasswdEncode) EncodePassword(password string) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	base64Salt := base64.StdEncoding.EncodeToString(salt)
	passwordHashed := hashPassword(password, salt)
	base64PasswordHashed := base64.StdEncoding.EncodeToString([]byte(passwordHashed))

	return base64PasswordHashed, base64Salt, nil
}
