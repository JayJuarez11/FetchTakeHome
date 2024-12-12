package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price string `json:"price" binding:"required"`
}
type Receipt struct {
	Retailer string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items []Item `json:"items" binding:"required,dive"`
	Total string `json:"total" binding:"required"`
}

var (
	idMap = make(map[string]int)
	mutex = &sync.Mutex{}
)

func main() {

	// Creating Gin Engine, HTTP Router
	r := gin.Default()

	// POST: Process Receipts
	r.POST("/receipts/process", func(c *gin.Context) {
		var receipt Receipt
		// Validating the receipt
		if err := c.ShouldBindJSON(&receipt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		// Receipt passed validation
		var newId = uuid.New().String()
		c.JSON(http.StatusOK, gin.H{"id": newId})
		// Updating our local db
		pointsAwarded := determinePointsAwarded(receipt)
		mutex.Lock()
		idMap[newId] = pointsAwarded
		mutex.Unlock()
	})

	// GET: Points awarded for the receipt
	r.GET("/receipts/:id/points", func(c *gin.Context) {
		id := c.Param("id")
		if val, ok := idMap[id]; ok {
			c.JSON(http.StatusOK, gin.H{"points": val})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Receipt has not been processed before. No points."})
			return
		}
	})

	r.Run(":8080")
}

func determinePointsAwarded(receipt Receipt) int{
	pointsAwarded := 0
	// Rule 1: One point for every alphanumeric character in the retailer name
	for _, val := range receipt.Retailer {
		if unicode.IsLetter(val) || unicode.IsNumber(val) {
			pointsAwarded++
		}
	}
	fmt.Println("Points Awarded After Rule 1: ", pointsAwarded)
	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalChange := receipt.Total[len(receipt.Total) - 2:]
	if totalChange == "00" {
		pointsAwarded += 50
	}
	fmt.Println("Points Awarded After Rule 2: ", pointsAwarded)
	// Rule 3: 25 points if the total is a multiple of 0.25.
	if totalChange == "00" || totalChange == "25" || totalChange == "50" || totalChange == "75" {
		pointsAwarded += 25
	}
	fmt.Println("Points Awarded After Rule 3: ", pointsAwarded)
	// Rule 4: 5 points for every two items on the receipt.
	pointsAwarded = pointsAwarded + (len(receipt.Items) / 2) * 5
	fmt.Println("Points Awarded After Rule 4: ", pointsAwarded)
	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up
	//to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedShortDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedShortDescription) % 3 == 0 {
			itemPrice, _ := strconv.ParseFloat(item.Price, 64)
			itemPointsAwarded:= int(math.Ceil(itemPrice * 0.2))
			fmt.Println("Item Points Awarded: ", itemPointsAwarded)
			pointsAwarded += itemPointsAwarded
		}
	}
	fmt.Println("Points Awarded After Rule 5: ", pointsAwarded)
	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDay := receipt.PurchaseDate[len(receipt.PurchaseDate) - 2:]
	if purchaseDay[1] % 2 != 0 {
		pointsAwarded += 6
	}
	fmt.Println("Points Awarded After Rule 6: ", pointsAwarded)
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	afterTime, _ := time.Parse("15:04", "14:00")
	beforeTime, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(afterTime) && purchaseTime.Before(beforeTime) {
		pointsAwarded += 10
	}
	fmt.Println("Points Awarded After Rule 7: ", pointsAwarded)

	return pointsAwarded
}