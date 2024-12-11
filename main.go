package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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

type ReceiptID struct {
	ID string `json:"id" binding:"required"`
}

func main() {
	fmt.Println("Creating Gin Engine.")
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
		c.JSON(http.StatusOK, gin.H{"id": uuid.New().String()})
	})

	// GET: Points awarded for the receipt
	r.GET("/receipts/{id}/points", func(c *gin.Context) {
		var id ReceiptID
		// Validating the ID
		if err := c.ShouldBindJSON(&id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
			return
		}
		// ID passed validation
		// TODO: Determine actual points awarded
		c.JSON(http.StatusOK, gin.H{"points": 100})
	})


	fmt.Println("Starting the server on port 8080.")
	r.Run(":8080")
}
