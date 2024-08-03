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

	_, err = internalGoogle.NewGmailService(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	sheets, err := internalGoogle.NewSheetsService(client)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	/*
		// Example usage of GmailService
		messages, err := gmailService.ListMessages("me")
		if err != nil {
			log.Fatalf("Unable to list messages: %v", err)
		}

		for _, m := range messages {
			msg, err := gmailService.GetMessage("me", m.Id)
			if err != nil {
				log.Fatalf("Unable to retrieve message: %v", err)
			}
			log.Printf("Message snippet: %s\n", msg.Snippet)
		}
	*/
	// Example usage of SheetsService
	data := []internalGoogle.FlightData{
		{Date: "2024-08-03", Type: "OneWay", Airline: "ExampleAir", Origin: "JFK", Destination: "LAX", Duration: "1hr", URL: "https://exampleurl.com", Price: "300"},
	}
	if err := sheets.AppendFlightData(data); err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}
