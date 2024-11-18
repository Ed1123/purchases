package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load(".env")
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	sheet_id, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		log.Fatalf("SHEET_ID not found")
	}
	readSheet(sheetsService, sheet_id, "A1:H4")
}
