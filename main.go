package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	conf  *oauth2.Config
	store *sessions.CookieStore
)

type AppAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}


func getAppAccessToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("TWITCH_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", data)
	if err != nil {
		return "", fmt.Errorf("failed to get app access token: %w", err)
	}
	defer resp.Body.Close()

	var result AppAccessTokenResponse
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println("App Token: ", result.AccessToken)

	return result.AccessToken, nil
}


func getUserInfo(accessToken string) (string, string, error) {

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string][]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	userID := result["data"][0]["id"].(string)
	username := result["data"][0]["login"].(string)

	return userID, username, nil
}


func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	conf = &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		Scopes:       []string{"openid", "channel:read:subscriptions", "channel:read:redemptions", "bits:read", "moderator:read:followers", "channel:manage:polls"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
		RedirectURL: os.Getenv("BASE_URL") + "/auth",
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	r := gin.Default()

	r.GET("/login", handleLogin)
	r.GET("/logout", handleLogout)
	r.GET("/auth", handleAuth)
	r.GET("/delete", func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session-name")
		accessToken := session.Values["access_token"].(string)

		userID, _, err := getUserInfo(accessToken)
		if err != nil {
			fmt.Println("Failed to get user info: ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}

		deleteSubscriptions(c, userID)
	})
	r.GET("/current", getSubscriptions)
	r.POST("/webhook", handleWebhook)

	r.Run()
}