package google

import (
	"fmt"
	"log"

	"github.com/Ed1123/purchases/src/models"
	sheets "google.golang.org/api/sheets/v4"
)

func readSheet(srv *sheets.Service, spreadsheetId, readRange string) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from spreadsheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("data:")
		for _, row := range resp.Values {
			fmt.Printf("%s, %s\n", row[0], row[1])
		}
	}
}

func AddPurchaseToSheet(srv *sheets.Service, spreadsheetId string, purchase models.PurchaseEntry) error {
	var vr sheets.ValueRange
	for _, item := range purchase.PurchaseItems {
		dataRow := []interface{}{purchase.Location, purchase.Date.Format("2006-01-02"),
			item.Name, item.Amount, item.Category}
		vr.Values = append(vr.Values,
			dataRow,
		)
	}
	_, err := srv.Spreadsheets.Values.
		Append(spreadsheetId, "A1", &vr).
		ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("Unable to write data to sheet: %w", err)
	}
	return nil
}
