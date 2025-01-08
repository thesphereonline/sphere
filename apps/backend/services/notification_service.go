package services

import (
	"context"
	"fmt"

	"github.com/thesphereonline/sphere/backend/api"
	"github.com/thesphereonline/sphere/backend/database"
)

// Publish notification to Redis
func PublishNotification(userID, message string) {
	database.RedisClient.Publish(context.Background(), fmt.Sprintf("notifications:%s", userID), message)
}

// Subscribe to Redis notifications
func SubscribeToNotifications() {
	pubsub := database.RedisClient.Subscribe(context.Background(), "notifications:*")

	for msg := range pubsub.Channel() {
		userID := msg.Channel[len("notifications:"):]
		api.SendNotification(userID, msg.Payload)
	}
}
