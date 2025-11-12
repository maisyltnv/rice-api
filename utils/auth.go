package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key (ໃນໂປຣເຈັກຈິງຄວນເກັບໃນ environment variable)
var JWTSecret = []byte(getenv("JWT_SECRET", "my-secret-key-change-in-production"))

// Hash password ດ້ວຍ bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ກວດສອບວ່າ password ຖືກຕ້ອງບໍ່
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ສ້າງ JWT token
func GenerateToken(userID uint, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role, // "admin" or "customer"
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Token ໝົດອາຍຸໃນ 7 ວັນ
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
