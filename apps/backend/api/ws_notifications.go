package api

import (
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
)

// Client structure to track connected users
type Client struct {
	ID   string
	Conn *websocket.Conn
}

var clients = make(map[string]*Client) // Map of connected users

// WebSocket connection handler
func NotificationHandler(c *websocket.Conn) {
	userID := c.Query("user_id") // Extract user ID from query parameters
	if userID == "" {
		c.Close()
		return
	}

	clients[userID] = &Client{ID: userID, Conn: c}
	defer func() {
		delete(clients, userID)
		c.Close()
	}()

	log.Println("User connected:", userID)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}

		fmt.Println("Received message:", string(msg))
	}
}

// Send notification to a specific user
func SendNotification(userID, message string) {
	client, exists := clients[userID]
	if exists {
		client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
