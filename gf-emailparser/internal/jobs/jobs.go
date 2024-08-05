package jobs

import (
	"fmt"
	"log"

	google "github.com/matty271828/flightalerts/gf-emailparser/internal/google"
)

type Jobs struct {
	Gmail  *google.GmailService
	Sheets *google.SheetsService
}

func NewJobs(gmail *google.GmailService, sheets *google.SheetsService) *Jobs {
	return &Jobs{
		Gmail:  gmail,
		Sheets: sheets,
	}
}

// ReadEmailsJob is used to access the gmail inbox, read all unread emails
// and extract and store found flight data.
func (j *Jobs) ReadEmailsJob() error {
	messages, err := j.Gmail.ListNewMessages("me")
	if err != nil {
		return fmt.Errorf("unable to list new messages: %v", err)
	}

	if len(messages) == 0 {
		log.Println("no new messages found.")
		return nil
	}

	for _, message := range messages {
		metaData, data, err := j.Gmail.ExtractFlightData(message)
		if err != nil {
			return fmt.Errorf("failed to extract flight data from email: %v", err)
		}

		if err := j.Sheets.AppendFlightData(*data); err != nil {
			return fmt.Errorf("unable to write data to sheet: %v", err)
		}

		err = j.Gmail.Sheets.MarkMessageAsRead(metaData.ID, metaData.InternalDate)
		if err != nil {
			return fmt.Errorf("failed to mark message as read: %v", err)
		}
	}
	return nil
}
