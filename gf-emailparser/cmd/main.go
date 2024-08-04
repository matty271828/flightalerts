// cmd/main.go
package main

import (
	"log"

	internalGoogle "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
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
		log.Fatalf("Unable to get client: %v", err)
	}

	sheets, err := internalGoogle.NewSheetsService(client)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	gmail, err := internalGoogle.NewGmailService(client, sheets)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Example usage of GmailService
	messages, err := gmail.ListNewMessages("me")
	if err != nil {
		log.Fatalf("Unable to list new messages: %v", err)
	}

	if len(messages) == 0 {
		log.Println("No new messages found.")
		return
	}

	data, err := gmail.ExtractFlightData(messages[0])
	if err != nil {
		log.Fatalf("Failed to extract flight data from email: %v", err)
	}

	if err := sheets.AppendFlightData(*data); err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}
