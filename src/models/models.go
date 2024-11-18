package models

import "time"

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
