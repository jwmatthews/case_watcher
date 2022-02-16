package spreadsheet

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/api"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

func CreateIfSheetDoesNotExist(srv *sheets.Service, spreadsheetId, sheetName string) {
	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&req},
	}

	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetId, rbb).Context(context.Background()).Do()
	if err != nil {
		log.Printf("Ignoring error from creation of sheet '%s', error message: %s", sheetName, err)
	}
}

func Update(spreadsheetId, email, private_key, private_key_id, originalSearchQuery string, data *api.ResponseCasesQueryBody) error {

	// See: https://stackoverflow.com/questions/58874474/unable-to-access-google-spreadsheets-with-a-service-account-credentials-using-go
	// Create a JWT configurations object for the Google service account
	conf := &jwt.Config{
		Email:        email,
		PrivateKey:   []byte(private_key),
		PrivateKeyID: private_key_id,
		TokenURL:     "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	client := conf.Client(oauth2.NoContext)

	// Create a service object for Google sheets
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	sheetName := fmt.Sprintf("OpenCases - %s", currentDate)
	CreateIfSheetDoesNotExist(srv, spreadsheetId, sheetName)

	sheetRange := fmt.Sprintf("%s!A1:Z9999", sheetName)
	log.Printf("Will update %s next\n", sheetRange)

	values := [][]interface{}{{"sample_A1", "sample_B1"}, {"sample_A2", "sample_B2"}, {"sample_A3", "sample_A3"}}
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  sheetRange,
		Values: values,
	})

	_, err = srv.Spreadsheets.Values.Clear(spreadsheetId, sheetRange, &sheets.ClearValuesRequest{}).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Error, failed to clear spreadsheet: %s %s, received error: %s\n", spreadsheetId, sheetRange, err)
	}

	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Error, failed to write to spreadsheet: %s %s, received error: %s\n", spreadsheetId, sheetRange, err)
	}

	return nil
}
