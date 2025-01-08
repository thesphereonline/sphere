package services

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thesphereonline/sphere/backend/database"
	"github.com/thesphereonline/sphere/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// JWT claims struct
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with plain text
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a new JWT token
func GenerateJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthenticateUser checks email and password, returning a JWT if successful
func AuthenticateUser(email, password string) (string, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", errors.New("user not found")
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return GenerateJWT(user.ID)
}
