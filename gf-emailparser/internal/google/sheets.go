package google

import (
	"context"
	"log"
	"net/http"

	sheets "google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	Service *sheets.Service
}

func NewSheetsService(client *http.Client) (*SheetsService, error) {
	service, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return &SheetsService{Service: service}, nil
}

func (s *SheetsService) WriteData(spreadsheetId string, writeRange string, values [][]interface{}) error {
	vr := &sheets.ValueRange{
		Values: values,
	}
	_, err := s.Service.Spreadsheets.Values.Update(spreadsheetId, writeRange, vr).ValueInputOption("RAW").Context(context.Background()).Do()
	if err != nil {
		return err
	}
	log.Println("Data written to sheet")
	return nil
}
