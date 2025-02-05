package utils

import "golang.org/x/crypto/bcrypt"

// VerifyPassword compares the hashed password with the plain text password
func VerifyPassword(hashedPassword, plainPassword string) error {
	// bcrypt.CompareHashAndPassword returns nil if passwords match, otherwise an error
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
