package main

import (
	"fmt"
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

	var items []PurchaseItem
	names := r.Form["name"]
	for i, name := range names {
		amount, err := strconv.ParseFloat(r.Form["amount"][i], 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		item := PurchaseItem{
			Name:     name,
			Amount:   float32(amount),
			Category: r.Form["category"][i],
		}
		items = append(items, item)
	}

	entry := PurchaseEntry{
		Location:      location,
		Date:          parsedDate,
		PurchaseItems: items,
	}

	// jsonData, _ := json.Marshal(entry)

	// Here you would send jsonData to the Google Apps Script URL
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
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/item", newItemHandler)

	slog.Info("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
