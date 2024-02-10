package main

import (
	"fmt"
	"log"
	"os"

	"pricesAPI/app"
	"pricesAPI/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	// Load environment variables from the .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found or unable to load.")
	}

	// Check if required environment variables are set
	requiredEnvVars := []string{"MONGODB_URL", "MONGODB_DB", "MODE"}
	for _, envVar := range requiredEnvVars {
		if _, exists := os.LookupEnv(envVar); !exists {
			log.Fatalf("Error: Required environment variable %s is not set.", envVar)
		}
	}

	if os.Getenv("MODE") == "release" {
		requiredEnvVars := []string{"SHOP_URL", "PRICE_DB", "MARKET_URL"}
		for _, envVar := range requiredEnvVars {
			if _, exists := os.LookupEnv(envVar); !exists {
				log.Fatalf("Error: Required environment variable %s is not set.", envVar)
			}
		}
	}

}

func main() {
	Init()
	fmt.Printf("%s", os.Getenv("SHOP_PORT"))
	app.InitMongoDB()
	router := gin.Default()
	// /product* routes
	router.GET("/product-data", handlers.GetProductsData)
	router.GET("/product-data/:name", handlers.GetProductsByName)

	// /supplier* routes
	router.GET("/supplier-data", handlers.GetSuppliersData)

	// /register* routes
	router.POST("/register", handlers.PostClientRegister)

	// /client* routes
	router.GET("/client", handlers.GetClientData)
	router.POST("/client/keys", handlers.PostClientKey)
	router.DELETE("/client", handlers.DeleteClient)

	// Index route
	router.GET("/", handlers.GetIndex)

	serverAddr := ":" + os.Getenv("PRICE_PORT")
	fmt.Printf("Server is running on http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
