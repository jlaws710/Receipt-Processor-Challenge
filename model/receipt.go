package model

import (
	"fmt"

	"github.com/google/uuid"
)

type shortDescription string

type price string

type Item struct {
	ShortDescription shortDescription `json:"shortDescription"`
	Price            price            `json:"price"`
}

type itemArray []Item

type Receipt struct {
	ID           string    `json: "id"`
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        itemArray `json:"items"`
	Total        string    `json:"total"`
	Points       int       `json:"points"`
}

var Receipts = []Receipt{}

func (receipt *Receipt) GenerateID() {
	if receipt.ID == "" {
		receipt.ID = uuid.New().String()
	}
}

func (receipt *Receipt) SetPoints(points int) {
	receipt.Points = points
}

func GetReceiptByID(id string) (*Receipt, error) {
	if id == "" {
		return nil, fmt.Errorf("No ID passed")
	}

	for _, receipt := range Receipts {
		if receipt.ID == id {
			return &receipt, nil
		}
	}
	return nil, fmt.Errorf("receipt not found")
}
