// google/gmail.go
package google

import (
	"fmt"
	"log"
	"net/http"

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
