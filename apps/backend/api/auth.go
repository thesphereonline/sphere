package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thesphereonline/sphere/backend/database"
	"github.com/thesphereonline/sphere/backend/models"
	"github.com/thesphereonline/sphere/backend/services"
)

// Signup request struct
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login request struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup handler
func Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	user := models.User{Email: req.Email, Password: hashedPassword}
	database.DB.Create(&user)

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// Login handler
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}
