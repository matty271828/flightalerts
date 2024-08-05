package google

import (
	"context"
	"fmt"
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

const (
	allFlights   = "all_flights"
	readMessages = "read_messages"
)

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
	Discount    string
}

// validateFlightData is used to validate that all fields in a struct of
// FlightData have been populated. The exception is discounts, which is optional.
func (f *FlightData) validateFlightData() bool {
	switch {
	case f.Date == "":
		return false
	case f.Type == "":
		return false
	case f.Airline == "":
		return false
	case f.Origin == "":
		return false
	case f.Destination == "":
		return false
	case f.Duration == "":
		return false
	case f.URL == "":
		return false
	case f.Price == "":
		return false
	}
	return true
}

// prepareDataForSheets is used to transform the custom type FlightData
// into a format expected by google sheets.
func prepareFlightDataForSheet(data []FlightData) [][]interface{} {
	var result [][]interface{}
	for _, fd := range data {
		result = append(result, []interface{}{
			fd.Date, fd.Type, fd.Airline, fd.Origin, fd.Destination, fd.Duration, fd.URL, fd.Price,
		})
	}
	return result
}

// AppendFlightData is used to add row(s) of FlightData to the
// flight_data sheet in Google Sheets.
func (s *SheetsService) AppendFlightData(data []FlightData) error {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	rangeToWrite := allFlights + "!A1" // Starting range, append will handle the rest
	// Prepare the data for appending
	values := prepareFlightDataForSheet(data)

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

// MessageMetaData is a custom type used to serve metadata
// on already read emails to Google Sheets. This metadata is then
// used to ensure FlightData is not duplicated when parsing emails.
type MessageMetaData struct {
	ID        string
	Timestamp string
}

// prepareMessageMetaDataForSheet is used to transform the custom type MessageMetaData
// into a format expected by Google Sheets.
func prepareMessageMetaDataForSheet(data []MessageMetaData) [][]interface{} {
	var result [][]interface{}
	for _, md := range data {
		result = append(result, []interface{}{
			md.ID, md.Timestamp,
		})
	}
	return result
}

// AppendMessageMetaData is used to add row(s) of MessageMetaData to the
// specified sheet in Google Sheets.
func (s *SheetsService) AppendMessageMetaData(data []MessageMetaData, sheetName string) error {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	rangeToWrite := sheetName + "!A1" // Starting range, append will handle the rest
	// Prepare the data for appending
	values := prepareMessageMetaDataForSheet(data)

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

// GetLatestProcessedMessage is used to find the last previously read message in order
// to act as a cutoff for reading emails.
func (s *SheetsService) GetLatestProcessedMessageMetadata() (*MessageMetaData, error) {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")

	// Get the sheet metadata to find the total number of rows
	sheetMetadata, err := s.Service.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title,gridProperties(rowCount)))").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve sheet metadata: %v", err)
	}

	var totalRows int64
	for _, sheet := range sheetMetadata.Sheets {
		if sheet.Properties.Title == readMessages {
			totalRows = sheet.Properties.GridProperties.RowCount
			break
		}
	}

	if totalRows == 0 {
		return nil, fmt.Errorf("no data found in sheet")
	}

	// Construct the range to read the last row
	rangeToRead := fmt.Sprintf("%s!A%d:B%d", readMessages, totalRows, totalRows)

	resp, err := s.Service.Spreadsheets.Values.Get(spreadsheetId, rangeToRead).MajorDimension("ROWS").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	// No error since this indicates that no messages have been read yet.
	if len(resp.Values) == 0 {
		return nil, nil
	}

	latestRow := resp.Values[0]
	if len(latestRow) < 2 {
		return nil, fmt.Errorf("incomplete data in the last row")
	}

	messageMetaData := &MessageMetaData{
		ID:        fmt.Sprintf("%v", latestRow[0]),
		Timestamp: fmt.Sprintf("%v", latestRow[1]),
	}

	return messageMetaData, nil
}
