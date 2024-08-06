// google/auth.go
package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Google callback")
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

	saveToken("token.json", token)
	fmt.Fprintf(w, "OAuth2 Token saved successfully")
}

// GetClient returns an HTTP client based on OAuth 2.0 configuration and saved token.
// It waits until the token file is available before proceeding.
func GetClient(config *oauth2.Config) (*http.Client, error) {
	tokFile := "token.json"
	var tok *oauth2.Token

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("Error reading token from file: %v", err)
		return nil, err
	}

	if tok == nil {
		// Generate and log the OAuth URL for user authorization
		authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	}

	for {
		tok, err := tokenFromFile(tokFile)
		if err != nil {
			log.Printf("Error reading token from file: %v", err)
			return nil, err
		}

		if tok != nil {
			break
		}

		log.Println("Token not found, waiting...")
		time.Sleep(5 * time.Second) // Wait for 5 seconds before checking again
	}

	log.Println("Token successfully read from file")
	return config.Client(context.Background(), tok), nil
}

// tokenFromFile reads a token from a file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	log.Printf("Reading token from file: %s", file)
	f, err := os.Open(file)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		// File does not exist, return nil token
		return nil, nil
	default:
		log.Printf("Error opening token file: %v", err)
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	switch {
	case err == nil:
	case err.Error() == "EOF":
		// File is empty, return nil token
		return nil, nil
	default:
		log.Printf("Error decoding token file: %v", err)
		return nil, err
	}

	log.Println("Token successfully decoded")
	return tok, nil
}

// saveToken saves the token to a file.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
