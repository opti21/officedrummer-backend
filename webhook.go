package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleWebhook(c *gin.Context) {
	// Get the type of the notification
	messageType := c.GetHeader("Twitch-Eventsub-Message-Type")

	switch messageType {
	case "notification":
		handleEventNotification(c)
	case "webhook_callback_verification":
		handleVerification(c)
	case "revocation":
		handleRevocation(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message type"})
	}
}

func handleEventNotification(c *gin.Context) {
	fmt.Println("Handling notification")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	// may need to recreate the request body stream since it has been read already
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	c.Status(http.StatusOK)

}

func handleVerification(c *gin.Context) {
	fmt.Println("Handling verification")

	// Parse the request body
	var body map[string]interface{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	challenge, ok := body["challenge"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No challenge found in request body"})
		return
	}

	c.String(http.StatusOK, challenge)
}

func handleRevocation(c *gin.Context) {
	c.Status(http.StatusOK)
}
