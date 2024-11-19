package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Ed1123/purchases/src/google"
	"github.com/Ed1123/purchases/src/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
)

var formTemplate = template.Must(template.ParseFiles("src/templates/form.html"))

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		formTemplate.Execute(w, nil)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	location := r.FormValue("location")
	date := r.FormValue("date")
	parsedDate, _ := time.Parse("2006-01-02", date)

	var items []models.PurchaseItem
	names := r.Form["name"]
	for i, name := range names {
		amount, err := strconv.ParseFloat(r.Form["amount"][i], 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		item := models.PurchaseItem{
			Name:     name,
			Amount:   float32(amount),
			Category: r.Form["category"][i],
		}
		items = append(items, item)
	}

	entry := models.PurchaseEntry{
		Location:      location,
		Date:          parsedDate,
		PurchaseItems: items,
	}

	srv, err := sheets.NewService(context.Background())
	if err != nil {
		slog.Error("Unable to create Google Sheets service", "error", err)
		http.Error(w, "Failed to connect to Google Sheets", http.StatusInternalServerError)
		return
	}

	spreadsheetId, ok := os.LookupEnv("SHEET_ID")
	if !ok {
		slog.Error("SHEET_ID not found")
		http.Error(w, "Failed to connect to Google Sheets", http.StatusInternalServerError)
		return
	}

	err = google.AddPurchaseToSheet(srv, spreadsheetId, entry)
	if err != nil {
		slog.Error("Failed to add purchase to sheet", "error", err)
		http.Error(w, "Failed to add purchase to sheet", http.StatusInternalServerError)
		return
	}

	slog.Info("Data sent to Google Apps Script: ", "data", entry)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data sent: %v", entry)))
}

func newItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	slog.Info("New item row requested")
	html := `
		<div class="item">
        	<label for="name">Name:</label>
        	<input id="name" type="text" name="name" required>
        	<label for="amount">Amount:</label>
        	<input id="amount" type="float" name="amount" required>
        	<label for="category">Category:</label>
        	<input id="category" type="text" name="category" required>
    	</div>
	`
	w.Write([]byte(html))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Failed to load environment variables", "error", err)
		os.Exit(1)
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT not found")
		os.Exit(1)
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/item", newItemHandler)

	slog.Info("Server started", "port", port)
	http.ListenAndServe(":"+port, nil)
}
