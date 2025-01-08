package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thesphereonline/sphere/backend/database"
)

// SaveActivity logs an event in Redis
func SaveActivity(userID, action string) {
	key := "activity:" + userID
	database.RedisClient.LPush(context.Background(), key, action)
	database.RedisClient.Expire(context.Background(), key, 24*time.Hour) // Expire in 24 hours
}

// GetActivityFeed retrieves recent activities
func GetActivityFeed(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	key := "activity:" + string(userID)

	activities, err := database.RedisClient.LRange(context.Background(), key, 0, 10).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch activity feed"})
	}

	return c.JSON(fiber.Map{"activities": activities})
}
