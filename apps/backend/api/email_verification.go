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
)

// Generate email verification token
func GenerateVerificationToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Send email verification link
func SendVerificationEmail(c *fiber.Ctx) error {
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

	token, err := GenerateVerificationToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", os.Getenv("FRONTEND_URL"), token)
	emailBody := fmt.Sprintf("Click <a href='%s'>here</a> to verify your email.", verificationLink)

	err = services.SendEmail(user.Email, "Email Verification", emailBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Verification email sent!"})
}

// Verify email
func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing token"})
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	userID := uint(claims["user_id"].(float64))
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_verified", true)

	return c.JSON(fiber.Map{"message": "Email verified successfully!"})
}
