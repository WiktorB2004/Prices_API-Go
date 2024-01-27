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

	router.GET("/product-data", handlers.GetProductsData)
	router.GET("/product-data/:name", handlers.GetProductsByName)
	router.GET("/supplier-data", handlers.GetSuppliersData)
	router.GET("/", handlers.GetIndex)

	serverAddr := ":3001"
	fmt.Printf("Server is running on http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
