package google

import (
	"context"
	"log"
	"net/http"
	"os"

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

func (s *SheetsService) AppendData(sheetName string, values [][]interface{}) error {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	rangeToWrite := sheetName + "!A1" // Starting range, append will handle the rest
	vr := &sheets.ValueRange{
		Values: values,
	}
	_, err := s.Service.Spreadsheets.Values.Append(spreadsheetId, rangeToWrite, vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
	if err != nil {
		return err
	}
	log.Println("Data appended to sheet")
	return nil
}
