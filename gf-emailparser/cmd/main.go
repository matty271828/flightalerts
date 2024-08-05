// cmd/main.go
package main

import (
	"log"

	internalGoogle "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
	"github.com/matty271828/flightalerts/gf-emailparser/internal/jobs"
	"github.com/matty271828/flightalerts/gf-emailparser/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Start OAuth server
	oauthConfig := server.InitOAuth()
	go server.Start()

	// Use the client after manual OAuth authentication
	client, err := internalGoogle.GetClient(oauthConfig)
	if err != nil {
		log.Fatalf("unable to get client: %v", err)
	}

	sheets, err := internalGoogle.NewSheetsService(client, internalGoogle.NewQueuer(100))
	if err != nil {
		log.Fatalf("unable to create Sheets service: %v", err)
	}

	gmail, err := internalGoogle.NewGmailService(client, sheets)
	if err != nil {
		log.Fatalf("unable to create Gmail service: %v", err)
	}

	jobs := jobs.NewJobs(gmail, sheets)

	// run job to read emails
	err = jobs.ReadEmailsJob()
	if err != nil {
		log.Fatalf("failed to run read emails job: %v", err)
	}
}
