package google

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Ed1123/purchases/src/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
)

func TestAddPurchaseToSheet(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load environment variables: %v", err)
	}
	srv, err := sheets.NewService(context.Background())
	if err != nil {
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		t.Fatalf("SHEET_ID not found")
	}
	purchase := models.PurchaseEntry{
		Location: "Test Location",
		Date:     time.Now(),
		PurchaseItems: []models.PurchaseItem{
			{
				Name:      "Test Item",
				Price:     1.99,
				Quantity:  1,
				Category:  "Test Category",
				Recipient: "Test Recipient",
			},
			{
				Name:      "Test Item",
				Price:     1.99,
				Quantity:  1,
				Category:  "Test Category",
				Recipient: "Test Recipient",
			},
		},
	}

	err = AddPurchaseToSheet(srv, spreadsheetId, purchase)
	if err != nil {
		t.Fatalf("Failed to add purchase to sheet: %v", err)
	}

}
func TestGetCategories(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Failed to load environment variables: %v", err)
	}
	srv, err := sheets.NewService(context.Background())
	if err != nil {
		t.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		t.Fatalf("SHEET_ID not found")
	}

	categories, err := GetCategories(srv, spreadsheetId)
	if err != nil {
		t.Fatalf("Failed to get categories from sheet: %v", err)
	}
	t.Log(categories)

	if len(categories) == 0 {
		t.Fatalf("Expected to find categories, but found none")
	}
}
