package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thesphereonline/sphere/backend/database"
	"github.com/thesphereonline/sphere/backend/models"
)

// GetUserProfile retrieves authenticated user's details
func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // Extract user ID from request context

	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"email": user.Email,
		"id":    user.ID,
	})
}
