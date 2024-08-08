// cmd/main.go
package main

import (
	"log"
	"time"

	"github.com/matty271828/flightalerts/gf-emailparser/internal/api"
	internalGoogle "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
	"github.com/matty271828/flightalerts/gf-emailparser/internal/jobs"
	"github.com/matty271828/flightalerts/gf-emailparser/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialise OAuth config
	oauthConfig, err := internalGoogle.InitOAuth()
	if err != nil {
		log.Fatalf("error initialising oauth2 config: %v", err)
	}

	// Start the web server. A limited start is needed here in order to authenticate
	// with google before retrieving the google client. Later we will initialise the
	// full server service including our custom APIs.
	go server.Start()

	// Use the client after manual OAuth authentication
	client, err := internalGoogle.GetClient(oauthConfig)
	if err != nil {
		log.Fatalf("unable to get client: %v", err)
	}

	// Block until the token is available
	if err := internalGoogle.WaitForToken(10 * time.Minute); err != nil {
		log.Fatalf("unable to get client: %v", err)
	}

	sheets, err := internalGoogle.NewSheetsService(client)
	if err != nil {
		log.Fatalf("unable to create Sheets service: %v", err)
	}

	gmail, err := internalGoogle.NewGmailService(client, sheets)
	if err != nil {
		log.Fatalf("unable to create Gmail service: %v", err)
	}

	jobs, err := jobs.NewJobs(gmail, sheets, oauthConfig)
	if err != nil {
		log.Fatalf("unable to create Jobs service: %v", err)
	}

	api, err := api.NewAPI(jobs)
	if err != nil {
		log.Fatalf("unable to create API service: %v", err)
	}

	// Start the full server service.
	err = server.NewServer(api)
	if err != nil {
		log.Fatalf("unable to create Server service: %v", err)
	}

	// Block forever to keep the service running
	select {}
}
