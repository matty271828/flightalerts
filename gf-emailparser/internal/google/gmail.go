// google/gmail.go
package google

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/api/gmail/v1"
)

type GmailService struct {
	Service *gmail.Service
	Sheets  *SheetsService
}

func NewGmailService(client *http.Client, sheets *SheetsService) (*GmailService, error) {
	service, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully initialised gmail service")
	return &GmailService{Service: service, Sheets: sheets}, nil
}

// ListNewMessages is used to return all currently unread flight alert emails in the inbox.
func (g *GmailService) ListNewMessages(user string) ([]*gmail.Message, error) {
	// Retrieve the latest processed message metadata
	processedMetadata, err := g.Sheets.GetLatestProcessedMessageMetadata()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve latest processed message metadata: %v", err)
	}

	// Prepare the query to fetch messages after the latest processed timestamp
	var query string
	if processedMetadata != nil {
		query = fmt.Sprintf("after:%s", processedMetadata.Timestamp)
	} else {
		query = "" // No filter if no previous metadata is available
	}

	// Retrieve messages from Gmail
	r, err := g.Service.Users.Messages.List(user).Q(query).MaxResults(100).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}

	// Create a map of processed message IDs for quick lookup
	processedMessageMap := make(map[string]struct{})
	if processedMetadata != nil {
		processedMessageMap[processedMetadata.ID] = struct{}{}
	}

	// Filter out messages that have already been processed
	var newMessages []*gmail.Message
	for _, msg := range r.Messages {
		if _, processed := processedMessageMap[msg.Id]; !processed {
			newMessages = append(newMessages, msg)
		}
	}

	return newMessages, nil
}

func (g *GmailService) GetMessage(user, messageId string) (*gmail.Message, error) {
	msg, err := g.Service.Users.Messages.Get(user, messageId).Do()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// decodeMessagePart decodes the message part and returns its plain text content
func (g *GmailService) DecodeMessagePart(part *gmail.MessagePart) (string, error) {
	data, err := base64.URLEncoding.DecodeString(part.Body.Data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getMessageContent traverses the message payload and returns the plain text content
func (g *GmailService) GetMessageContent(payload *gmail.MessagePart) string {
	if payload.MimeType == "text/plain" {
		content, err := g.DecodeMessagePart(payload)
		if err == nil {
			return content
		}
	}

	for _, part := range payload.Parts {
		content := g.GetMessageContent(part)
		if content != "" {
			return content
		}
	}

	return ""
}

// ExtractFlightData is used to return the flight data contained
// within a gmail flight alert.
func (g *GmailService) ExtractFlightData(message *gmail.Message) (*[]FlightData, error) {
	// Retrieve the full message
	fullMessage, err := g.GetMessage("me", message.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve message: %v", err)
	}

	// Print the entire message payload
	content := g.GetMessageContent(fullMessage.Payload)
	if content == "" {
		return nil, fmt.Errorf("no plain text content found in the message")
	}

	return extractFlightData(content), nil
}

// ExtractFlightData parses the message content and extracts flight data
func extractFlightData(content string) *[]FlightData {
	var (
		flights       []FlightData
		flight        FlightData
		combinedLines string
	)

	// Regular expressions to extract required fields
	reDate := regexp.MustCompile(`\b(\w{3}, \w{3} \d{1,2})\b`)
	reAirline := regexp.MustCompile(`\b(\w+(?: \w+)?) · Nonstop`)
	reOriginDestination := regexp.MustCompile(`· (\w{3})–(\w{3}) ·`)
	reDuration := regexp.MustCompile(`· (\d+ hr)`)
	reURL := regexp.MustCompile(`\((https?://[^\s)]+)\)`)
	rePrice := regexp.MustCompile(`From £(\d+)`)
	reDiscount := regexp.MustCompile(`SAVE (\d+%)`)

	// Split the content by new lines
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if match := reDate.FindString(line); match != "" {
			// Initialize a new FlightData instance here since the date
			// acts as a header for a new flight in gmail flight alerts.
			flight = FlightData{Date: match}
		}

		// TODO: Extract either "one-way" or "return" when we expand to
		// set up flight alerts for return flights.
		flight.Type = "OneWay"

		// Extract and set the Airline
		if match := reAirline.FindStringSubmatch(line); len(match) > 1 {
			flight.Airline = match[1]
		}

		// Extract and set the Origin and Destination
		if match := reOriginDestination.FindStringSubmatch(line); len(match) > 2 {
			origins := strings.Split(match[1], ", ")
			flight.Origin = match[1]
			if len(origins) == 2 {
				flight.Origin = origins[0]
			}
			flight.Destination = match[2]
		}

		// Extract and set the Duration
		if match := reDuration.FindStringSubmatch(line); len(match) > 1 {
			flight.Duration = match[1]
		}

		// Start combining lines after detecting "View"
		if strings.Contains(line, "View") {
			combinedLines = line
			continue
		}

		if combinedLines != "" {
			combinedLines += line
			// Extract the URL, ensuring to remove unwanted parts
			if match := reURL.FindString(combinedLines); match != "" {
				// Clean the URL from "View" and brackets
				cleanURL := strings.Trim(match, "()")
				flight.URL = cleanURL
				combinedLines = ""
			}
		}

		// Extract and set the Discount
		if match := reDiscount.FindStringSubmatch(line); len(match) > 1 {
			flight.Discount = match[1]
		}

		// Extract and set the Price
		if match := rePrice.FindStringSubmatch(line); len(match) > 1 {
			flight.Price = "£" + match[1]
		}

		if flight.validateFlightData() {
			flights = append(flights, flight)
			flight = FlightData{}
		}

	}
	return &flights
}
