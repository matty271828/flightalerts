package jobs

import (
	"fmt"
	"log"
	"sync"
	"time"

	google "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
	"golang.org/x/oauth2"
)

type Jobs struct {
	Config *oauth2.Config
	Gmail  *google.GmailService
	Sheets *google.SheetsService
	Wg     *sync.WaitGroup
}

func NewJobs(gmail *google.GmailService, sheets *google.SheetsService, config *oauth2.Config) (*Jobs, error) {
	wg := sync.WaitGroup{}

	return &Jobs{
		Config: config,
		Gmail:  gmail,
		Sheets: sheets,
		Wg:     &wg,
	}, nil
}

// RefreshOauthTokenJob is a used to check if the oauth token has expired
// and refresh it if it has.
func (j *Jobs) RefreshTokenJob() error {
	log.Println("RefreshOauthTokenJob requested.")

	_, err := google.GetClient(j.Config)
	if err != nil {
		log.Println("RefreshOauthTokenJob failed.")
	}

	return nil
}

// ReadEmailsJob is used to access the gmail inbox, read all unread emails
// and extract and store found flight data.
func (j *Jobs) ReadEmailsJob() error {
	// Add a count to the WaitGroup for the goroutine
	j.Wg.Add(1)

	// Spawn a goroutine to handle email reading and processing
	go func() {
		defer j.Wg.Done() // Mark the goroutine as done when it finishes

		if err := j.ReadEmailsSubJob(); err != nil {
			log.Printf("Error in ReadEmailsSubJob: %v", err)
		}
	}()

	return nil
}

// ReadEmailSubJob is a subjob used to read and extract the contents
// for an individual email.
func (j *Jobs) ReadEmailsSubJob() error {
	messages, err := j.Gmail.ListNewMessages("me")
	if err != nil {
		return fmt.Errorf("unable to list new messages: %v", err)
	}

	if len(messages) == 0 {
		log.Println("no new messages found.")
		return nil
	}

	for i, message := range messages {
		// Add a 2 second delay between processing each message
		time.Sleep(2 * time.Second)

		metaData, data, err := j.Gmail.ExtractFlightData(message)
		if err != nil {
			return fmt.Errorf("failed to extract flight data from email: %v", err)
		}

		var (
			id           = metaData.ID
			internalDate = metaData.InternalDate
		)

		if err := j.Sheets.AppendFlightData(*data); err != nil {
			return fmt.Errorf("unable to write data to sheet: %v", err)
		}

		err = j.Gmail.Sheets.MarkMessageAsRead(id, internalDate)
		if err != nil {
			return fmt.Errorf("failed to mark message as read: %v", err)
		}

		if i == 0 {
			time.Sleep(2 * time.Second)
			err = j.Gmail.Sheets.MarkMessageAsCutoff(id, internalDate)
			if err != nil {
				return fmt.Errorf("failed to mark message as cutoff: %v", err)
			}
		}
	}

	return nil
}
