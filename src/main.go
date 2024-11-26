package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ed1123/purchases/src/google"
	"github.com/Ed1123/purchases/src/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
)

var form = template.Must(template.ParseFiles("src/templates/form.html"))
var submitted = template.Must(template.ParseFiles("src/templates/submitted.html"))

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		form.Execute(w, nil)
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
	merchant := r.FormValue("merchant")
	date := r.FormValue("date")
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	formTotal, err := strconv.ParseFloat(r.FormValue("total"), 32)
	if err != nil {
		http.Error(w, "Invalid total", http.StatusBadRequest)
		return
	}

	var items []models.PurchaseItem
	formItems := r.Form["name"]
	for i := range formItems {
		price, err := strconv.ParseFloat(r.Form["price"][i], 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		quantity, err := strconv.Atoi(r.Form["quantity"][i])
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}
		item := models.PurchaseItem{
			Name:      r.Form["name"][i],
			Price:     float32(price),
			Quantity:  quantity,
			Category:  r.Form["category"][i],
			Recipient: r.Form["recipient"][i],
		}
		items = append(items, item)
	}

	correctTax(items, float32(formTotal))

	entry := models.PurchaseEntry{
		Merchant:      merchant,
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
	submitted.Execute(w, entry)
}

// Applies taxes proportionally if the total amount is different
// from the sum of the items
func correctTax(items []models.PurchaseItem, taxedTotal float32) {
	var calcTotal float32
	for _, item := range items {
		calcTotal += item.Price * float32(item.Quantity)
	}
	if taxedTotal == calcTotal {
		return
	}
	tax := taxedTotal - calcTotal
	for i, item := range items {
		item.Price += tax / calcTotal * item.Price
		items[i] = item
	}
}

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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

	data, err := google.GetAutocompleteData(srv, spreadsheetId)
	if err != nil {
		slog.Error("Failed to get autocomplete data", "error", err)
		http.Error(w, "Failed to get autocomplete data", http.StatusInternalServerError)
		return
	}
	option := "<option value=\"%s\"></option>"
	options := ""
	for _, merchant := range data.Merchants {
		options += fmt.Sprintf(option, merchant)
	}
	merchants := fmt.Sprintf("<datalist id=\"merchants\" >%s</datalist>", options)

	options = ""
	for _, itemName := range data.ItemNames {
		options += fmt.Sprintf(option, itemName)
	}
	itemNames := fmt.Sprintf("<datalist id=\"item-names\" >%s</datalist>", options)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(merchants + itemNames))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn(".env file not found")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT not found")
		os.Exit(1)
	}

	http.Handle("/static/", http.FileServer(http.Dir("src")))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/autocomplete", autocompleteHandler)

	slog.Info("Server started", "port", port)
	http.ListenAndServe(":"+port, nil)
}
