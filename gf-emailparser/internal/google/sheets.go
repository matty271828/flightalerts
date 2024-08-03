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
	log.Println("Successfully initialised sheets service")
	return &SheetsService{Service: service}, nil
}

// FlightData is a custom type used to serve flight information
// to google sheets.
type FlightData struct {
	Date        string
	Type        string
	Airline     string
	Origin      string
	Destination string
	Duration    string
	URL         string
	Price       string
}

// prepareDataForSheets is used to transform the custom type FlightData
// into a format expected by google sheets.
func prepareDataForSheet(data []FlightData) [][]interface{} {
	var result [][]interface{}
	for _, fd := range data {
		result = append(result, []interface{}{
			fd.Date, fd.Type, fd.Airline, fd.Origin, fd.Destination, fd.Duration, fd.URL, fd.Price,
		})
	}
	return result
}

func (s *SheetsService) AppendData(sheetName string, data []FlightData) error {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	rangeToWrite := sheetName + "!A1" // Starting range, append will handle the rest
	// Prepare the data for appending
	values := prepareDataForSheet(data)

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
