package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type PurchaseItem struct {
	Name     string  `json:"name"`
	Amount   float32 `json:"amount"`
	Category string  `json:"category"`
}

type PurchaseEntry struct {
	Location      string         `json:"location"`
	Date          time.Time      `json:"date"`
	PurchaseItems []PurchaseItem `json:"purchaseItems"`
}

var formTemplate = template.Must(template.ParseFiles("src/templates/form.html"))

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	slog.Info("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		formTemplate.Execute(w, nil)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		location := r.FormValue("location")
		date := r.FormValue("date")

		parsedDate, _ := time.Parse("2006-01-02", date)

		amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		item := PurchaseItem{
			Name:     r.FormValue("name"),
			Amount:   float32(amount),
			Category: r.FormValue("category"),
		}

		entry := PurchaseEntry{
			Location:      location,
			Date:          parsedDate,
			PurchaseItems: []PurchaseItem{item},
		}

		jsonData, _ := json.Marshal(entry)

		// Here you would send jsonData to the Google Apps Script URL
		slog.Info("Data sent to Google Apps Script: ", "data", jsonData)

		w.Write([]byte("Form submitted successfully!"))
	}
}
