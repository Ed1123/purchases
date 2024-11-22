package google

import (
	"fmt"

	"github.com/Ed1123/purchases/src/models"
	sheets "google.golang.org/api/sheets/v4"
)

func AddPurchaseToSheet(srv *sheets.Service, spreadsheetId string, purchase models.PurchaseEntry) error {
	var vr sheets.ValueRange
	for _, item := range purchase.PurchaseItems {
		dataRow := []interface{}{purchase.Location, purchase.Date.Format("2006-01-02"),
			item.Name, item.Price, item.Quantity, item.Category, item.Recipient}
		vr.Values = append(vr.Values,
			dataRow,
		)
	}
	_, err := srv.Spreadsheets.Values.
		Append(spreadsheetId, "A1", &vr).
		ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("unable to write data to sheet: %w", err)
	}
	return nil
}

func GetCategories(srv *sheets.Service, spreadsheetId string) ([]string, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, "categories!A:A").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve categories from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("no cateogories found in spreadsheet")
	}

	var categories []string
	for _, row := range resp.Values {
		categories = append(categories, row[0].(string))
	}
	return categories, nil
}
