package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thesphereonline/sphere/backend/services"
)

// UploadVideo handles video file upload
func UploadVideo(c *fiber.Ctx) error {
	// Get file from request
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file"})
	}

	// Upload file to S3
	videoURL, err := services.UploadVideo(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload video"})
	}

	// Return video URL
	return c.JSON(fiber.Map{"video_url": videoURL})
}
