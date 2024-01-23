package main

import (
	"fmt"
	"log"

	"pricesAPI/app"
	"pricesAPI/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app.InitMongoDB()
	router := gin.Default()

	router.GET("/get-product-data", handlers.GetProductData)
	router.GET("/get-supplier-data", handlers.GetSupplierData)

	serverAddr := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
