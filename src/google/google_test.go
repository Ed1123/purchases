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
			{Name: "Item1", Amount: 10.0, Category: "Category1"},
			{Name: "Item2", Amount: 20.0, Category: "Category2"},
		},
	}

	err = AddPurchaseToSheet(srv, spreadsheetId, purchase)
	if err != nil {
		t.Fatalf("Failed to add purchase to sheet: %v", err)
	}

}
