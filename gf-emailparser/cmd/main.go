package main

import (
	"log"

	google "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
)

func main() {
	client, err := google.GetClient()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	gmailService, err := google.NewGmailService(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	sheetsService, err := google.NewSheetsService(client)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

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

	// Example usage of SheetsService
	sheetName := "all_flights"
	values := [][]interface{}{
		{"Date", "Type", "Airline", "Origin", "Destination", "Duration", "URL", "Price"},
		{"2024-08-03", "OneWay", "ExampleAir", "JFK", "LAX", "1hr", "https://exampleurl.com", "300"},
	}
	if err := sheetsService.AppendData(sheetName, values); err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}
