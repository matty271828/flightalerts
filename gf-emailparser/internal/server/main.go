// server/main.go
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	internalGoogle "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
	oauthState  string
)

func InitOAuth() *oauth2.Config {
	// Initialize OAuth
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}
	oauthState = "state-token" // This should be a random unique string in production
	return oauthConfig
}

func Start() {
	http.HandleFunc("/callback", handleGoogleCallback)
	log.Println("Starting google callback server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling google callback")
	if r.FormValue("state") != oauthState {
		http.Error(w, "State parameter doesn't match", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	internalGoogle.SaveToken("token.json", token)

	fmt.Fprintf(w, "OAuth2 Token saved successfully")
}
