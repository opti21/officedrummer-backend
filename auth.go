package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func handleLogin(c *gin.Context) {

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleAuth(c *gin.Context) {
	code := c.Query("code")

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Failed to exchange token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Save the access token in the user's session
	session, _ := store.Get(c.Request, "session-name")
	session.Values["access_token"] = token.AccessToken
	session.Save(c.Request, c.Writer)

	// Get the user's Twitch username and ID
	userID, username, err := getUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println("Failed to get user info: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}	

	appToken, err := getAppAccessToken()
	if err != nil {
		fmt.Println("Failed to get app access token: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app access token"})
		return
	}

	// Create the subscriptions
	createSubscriptions(appToken, userID)

	fmt.Println("Successfully authenticated user: ", username, " userID: ", userID)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "username": username, "userID": userID})
}

func handleLogout(c *gin.Context) {
	// Get the user's session
	session, _ := store.Get(c.Request, "session-name")

	// Clear the session
	session.Options = &sessions.Options{
		MaxAge:   -1,
		HttpOnly: true,
	}

	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}