// cmd/main.go
package main

import (
	"fmt"
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

	// Retrieve the full message
	firstMessage, err := gmail.GetMessage("me", messages[0].Id)
	if err != nil {
		log.Fatalf("Unable to retrieve message: %v", err)
	}

	log.Printf("Message ID: %s", firstMessage.Id)
	log.Printf("Message snippet: %s", firstMessage.Snippet)

	// Print the entire message payload
	content := internalGoogle.GetMessageContent(firstMessage.Payload)
	if content != "" {
		fmt.Println("Message Content:")
		fmt.Println(content)
	} else {
		log.Println("No plain text content found in the message.")
	}

	// Example usage of SheetsService
	data := []internalGoogle.FlightData{
		{Date: "2024-08-03", Type: "OneWay", Airline: "ExampleAir", Origin: "JFK", Destination: "LAX", Duration: "1hr", URL: "https://exampleurl.com", Price: "300"},
	}
	if err := sheets.AppendFlightData(data); err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}
