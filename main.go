package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price string `json:"price" binding:"required"`
}
type Receipt struct {
	Retailer string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"retailer" binding:"required"`
	PurchaseTime string `json:"purchaseTime" binding:"required"`
	Items []Item `json:"items" binding:"required,dive"`
	Total string `json:"total" binding:"required"`
}

func main() {
	fmt.Println("Hello, Go with IntelliJ!")
	// Creating Gin Engine, HTTP Router
	r := gin.Default()

	// POST: Process Receipts
	r.POST("/receipts/process", func(c *gin.Context) {
		var receipt Receipt
		// Validating the receipt
		if err := c.ShouldBindJSON(&receipt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Receipt passed validation
		c.JSON(http.StatusOK, gin.H{"id": uuid.New().String()})
	})


	fmt.Println("Starting the server on port 8080")
	r.Run(":8080")
}
