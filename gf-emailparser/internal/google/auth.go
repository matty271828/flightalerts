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

const (
	tokFile = "token.json"
)

func InitOAuth() (*oauth2.Config, error) {
	// Retrieve environment variables
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")

	// Check for empty values
	if clientID == "" {
		return nil, fmt.Errorf("CLIENT_ID environment variable is not set")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("CLIENT_SECRET environment variable is not set")
	}
	if redirectURL == "" {
		return nil, fmt.Errorf("REDIRECT_URL environment variable is not set")
	}

	// Initialize OAuth
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}
	oauthState = "state-token" // This should be a random unique string in production

	return oauthConfig, nil
}

// HandleGoogleCallback handles the Google OAuth2 callback
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Google callback")
	state := r.FormValue("state")
	if state != oauthState {
		log.Printf("Invalid OAuth state, expected '%s', got '%s'\n", oauthState, state)
		http.Error(w, "State parameter doesn't match", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Println("Code parameter is missing")
		http.Error(w, "Code parameter is missing", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed: %v\n", err)
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	err = saveToken(tokFile, token)
	if err != nil {
		log.Printf("Failed to save OAuth2 token: %v\n", err)
		http.Error(w, "Failed to save OAuth2 token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OAuth2 Token saved successfully")
}

// GetClient returns an HTTP client based on OAuth 2.0 configuration and saved token.
// It waits until the token file is available before proceeding.
func GetClient(config *oauth2.Config) (*http.Client, error) {
	var tok *oauth2.Token

	// Check if there is an existing token saved.
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

	// Refresh the token if it has expired.
	if tok.Expiry.Before(time.Now()) {
		log.Println("Token expired, attempting to refresh")
		tokSource := config.TokenSource(context.Background(), tok)
		tok, err = tokSource.Token()
		if err != nil {
			log.Printf("Unable to refresh token: %v", err)
			return nil, err
		}

		// Save the refreshed token
		err = saveToken(tokFile, tok)
		if err != nil {
			log.Printf("Unable to save refreshed token: %v", err)
			return nil, err
		}
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
		log.Println("nil token file")
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
		log.Println("nil token")
		return nil, nil
	default:
		log.Printf("Error decoding token file: %v", err)
		return nil, err
	}

	log.Println("Token successfully decoded")
	return tok, nil
}

// saveToken saves the token to a file
func saveToken(path string, token *oauth2.Token) error {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create oauth token cache file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("Error closing file: %v\n", cerr)
		}
	}()

	if err := json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("unable to encode oauth token to file: %v", err)
	}

	log.Println("Token saved successfully")
	return nil
}

// WaitForToken is used to block until a token has been received from OAuth2
func WaitForToken(timeout time.Duration) error {
	start := time.Now()
	for {
		if TokenExists(tokFile) {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timed out waiting for token")
		}
		time.Sleep(5 * time.Second) // Sleep for 10 seconds before checking again
		log.Println("Waiting for token...")
	}
}

func TokenExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
