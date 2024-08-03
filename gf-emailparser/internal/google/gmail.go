// google/gmail.go
package google

import (
	"net/http"

	"google.golang.org/api/gmail/v1"
)

type GmailService struct {
	Service *gmail.Service
}

func NewGmailService(client *http.Client) (*GmailService, error) {
	service, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return &GmailService{Service: service}, nil
}

func (g *GmailService) ListMessages(user string) ([]*gmail.Message, error) {
	r, err := g.Service.Users.Messages.List(user).Q("label:inbox").MaxResults(10).Do()
	if err != nil {
		return nil, err
	}
	return r.Messages, nil
}

func (g *GmailService) GetMessage(user, messageId string) (*gmail.Message, error) {
	msg, err := g.Service.Users.Messages.Get(user, messageId).Do()
	if err != nil {
		return nil, err
	}
	return msg, nil
}
