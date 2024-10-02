package main

import (
	"fmt"
	"math"
	"net/http"
	"receipt-processor-challenge/model"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var newReceipt model.Receipt

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", GetPointsForReceipt)

	router.Run("localhost:8080")
}

func processReceipt(ctx *gin.Context) {

	newReceipt.Points = 0
	newReceipt.ID = ""

	newReceipt.GenerateID()
	newReceipt.SetPoints(points(newReceipt))

	model.Receipts = append(model.Receipts, newReceipt)

	ctx.Bind(&newReceipt)
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"id:": newReceipt.ID,
	})
}

func GetPointsForReceipt(ctx *gin.Context) {
	receiptID := ctx.Param("id")

	receipt, err := model.GetReceiptByID(receiptID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for this id"})

		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"points": receipt.Points})
}

func countAlphanumeric(str string) int {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	count := 0

	for _, v := range str {
		if strings.ContainsRune(alphanumeric, v) {
			count++
		}
	}
	return count
}

func points(reciept model.Receipt) int {
	points := 0
	// Add a point for every alphanumeric character in retailer name
	points += countAlphanumeric(reciept.Retailer)
	// Add 50 points if total is integer
	total, _ := strconv.ParseFloat(reciept.Total, 64)

	if total == math.Trunc(total) {
		points += 50
	}
	// Add 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	// Add 5 points for every two items on the receipt
	totalItems := len(reciept.Items)

	if totalItems > 2 {
		points += totalItems / 2 * 5
	}

	for _, item := range reciept.Items {
		/**
		If the trimmed length of the item description is a multiple of 3,
		multiply the price by 0.2 and round up to the nearest integer.
		*/
		descrLen := utf8.RuneCountInString(strings.TrimSpace(string(item.ShortDescription)))

		if descrLen%3 == 0 {
			itemPrice, _ := strconv.ParseFloat(string(item.Price), 64)
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}
	// Combine date and time into a parsed time.Time
	purchasedAt, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v %v", reciept.PurchaseDate, reciept.PurchaseTime))
	// error check
	if err != nil {
		return points
	}
	// Add 6 points if purchase date is odd.
	if purchasedAt.Day()%2 != 0 {
		points += 6
	}
	// Add 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if purchasedAt.Hour() >= 14 && purchasedAt.Hour() <= 16 {
		points += 10
	}
	return points
}
