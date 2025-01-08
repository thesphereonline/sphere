package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/thesphereonline/sphere/backend/api"
	"github.com/thesphereonline/sphere/backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")

	apiGroup.Get("/health", api.HealthCheck)
	apiGroup.Post("/signup", api.Signup)
	apiGroup.Post("/login", api.Login)
	apiGroup.Post("/send-verification-email", api.SendVerificationEmail)
	apiGroup.Get("/verify-email", api.VerifyEmail)
	apiGroup.Post("/send-password-reset", api.SendPasswordResetEmail)
	apiGroup.Post("/reset-password", api.ResetPassword)

	// Protected routes
	protected := apiGroup.Group("/user")
	protected.Use(middleware.AuthRequired)
	protected.Get("/profile", api.GetUserProfile)
	protected.Post("/upload-video", api.UploadVideo)
	protected.Get("/activity-feed", api.GetActivityFeed)

	// WebSocket route for notifications
	app.Get("/ws/notifications", websocket.New(api.NotificationHandler))
}
