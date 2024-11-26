package models

import "time"

type PurchaseItem struct {
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Quantity  int     `json:"quantity"`
	Category  string  `json:"category"`
	Recipient string  `json:"recipient"`
}

type PurchaseEntry struct {
	Merchant      string         `json:"merchant"`
	Date          time.Time      `json:"date"`
	PurchaseItems []PurchaseItem `json:"purchaseItems"`
}
