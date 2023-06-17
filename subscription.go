package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SubscriptionRequest struct {
	Type      string `json:"type"`
	Version   string `json:"version"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
		ModeratorUserID   string `json:"moderator_user_id"`
	} `json:"condition"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
		Secret   string `json:"secret"`
	} `json:"transport"`
}

func createSubscriptions(appToken string, broadcasterUserID string) error {

	subTypes := []struct {
		Type    string
		Version string
	}{
		{"channel.follow", "2"},
		{"channel.subscribe", "1"},
		{"channel.cheer", "1"},
		{"stream.online", "1"},
		{"stream.offline", "1"},
	}

	errChan := make(chan error, len(subTypes))

	for _, sub := range subTypes {
		go func(subType string, subVersion string) {

			reqBody := SubscriptionRequest{
				Type:    sub.Type,
				Version: sub.Version,
			}

			reqBody.Transport.Method = "webhook"
			reqBody.Transport.Callback = os.Getenv("BASE_URL") + "/webhook"
			reqBody.Transport.Secret = os.Getenv("WEBHOOK_SECRET")

			if sub.Type == "channel.follow" {
				reqBody.Condition.BroadcasterUserID = broadcasterUserID
				reqBody.Condition.ModeratorUserID = broadcasterUserID
			} else {
				reqBody.Condition.BroadcasterUserID = broadcasterUserID
			}

			jsonReqBody, err := json.Marshal(reqBody)
			if err != nil {
				errChan <- fmt.Errorf("failed to create subscription request: %w", err)
				return
			}

			req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewBuffer(jsonReqBody))
			if err != nil {
				// print error
				fmt.Println(err)
				errChan <- fmt.Errorf("failed to create subscription request: %w", err)
				return
			}

			req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
			req.Header.Set("Authorization", "Bearer " + appToken)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errChan <- fmt.Errorf("failed to create %s subscription", sub.Type)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusAccepted {
				errChan <- fmt.Errorf("failed to create %s subscription", sub.Type)
				return
			}
		}(sub.Type, sub.Version)

		for range subTypes {
			if err := <-errChan; err != nil {
				return err
			}
		}

	}

	return nil
}

func deleteSubscriptions(c *gin.Context, broadcasterID string) {
	// Get the app access token
	appAccessToken, err := getAppAccessToken()
	if err != nil {
		fmt.Println("Failed to get app access token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app access token"})
		return
	}

	// Get the list of subscriptions
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/eventsub/subscriptions", nil)
	if err != nil {
		fmt.Println("Failed to create request: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer " + appAccessToken)
	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to get subscriptions: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}
	defer resp.Body.Close()

	var result map[string][]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	for _, subscription := range result["data"] {
		if subscription["condition"].(map[string]interface{})["broadcaster_user_id"].(string) == broadcasterID {
			req, err := http.NewRequest("DELETE", "https://api.twitch.tv/helix/eventsub/subscriptions?id="+subscription["id"].(string), nil)
			if err != nil {
				fmt.Println("Failed to create delete request: ", err.Error())
				continue
			}

			req.Header.Set("Authorization", "Bearer " + appAccessToken)
			req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

			_, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Failed to delete subscription: ", err.Error())
				continue
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted subscriptions"})
}

func getSubscriptions(c *gin.Context) {

	appAccessToken, err := getAppAccessToken()
	if err != nil {
		fmt.Println("Failed to get app access token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app access token"})
		return
	}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/eventsub/subscriptions", nil)
	if err != nil {
		fmt.Println("Failed to create request: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+appAccessToken)
	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to get subscriptions: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}
	defer resp.Body.Close()

	var result map[string][]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(http.StatusOK, gin.H{"subscriptions": result["data"]})
}
