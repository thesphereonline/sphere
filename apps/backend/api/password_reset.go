package api

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/thesphereonline/sphere/backend/database"
	"github.com/thesphereonline/sphere/backend/models"
	"github.com/thesphereonline/sphere/backend/services"
	"golang.org/x/crypto/bcrypt"
)

// Generate password reset token
func GeneratePasswordResetToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Send password reset email
func SendPasswordResetEmail(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	token, err := GeneratePasswordResetToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)
	emailBody := fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetLink)

	err = services.SendEmail(user.Email, "Password Reset", emailBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Password reset email sent!"})
}

// Reset password
func ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	userID := uint(claims["user_id"].(float64))
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)

	return c.JSON(fiber.Map{"message": "Password reset successful!"})
}
