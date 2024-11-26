package google

import (
	"fmt"

	"github.com/Ed1123/purchases/src/models"
	sheets "google.golang.org/api/sheets/v4"
)

func AddPurchaseToSheet(srv *sheets.Service, spreadsheetId string, purchase models.PurchaseEntry) error {
	var vr sheets.ValueRange
	for _, item := range purchase.PurchaseItems {
		dataRow := []interface{}{purchase.Merchant, purchase.Date.Format("2006-01-02"),
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
		return nil, fmt.Errorf("no categories found in spreadsheet")
	}

	var categories []string
	for _, row := range resp.Values {
		categories = append(categories, row[0].(string))
	}
	return categories, nil
}

type AutocompleteData struct {
	Merchants  []string
	Names      []string
	Categories []string
}

func GetAutocompleteData(srv *sheets.Service, spreadsheetId string) (AutocompleteData, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, "autocomplete!A:C").Do()
	if err != nil {
		return AutocompleteData{}, fmt.Errorf("unable to retrieve autocomplete data from sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return AutocompleteData{}, fmt.Errorf("no autocomplete data found in spreadsheet")
	}

	var data AutocompleteData
	for _, row := range resp.Values {
		var name, category string
		merchant := row[0].(string)
		if len(row) >= 2 {
			name = row[1].(string)
		}
		if len(row) >= 3 {
			category = row[2].(string)
		}

		if merchant != "" {
			data.Merchants = append(data.Merchants, merchant)
		}
		if name != "" {
			data.Names = append(data.Names, name)
		}
		if category != "" {
			data.Categories = append(data.Categories, category)
		}
	}
	return data, nil
}
