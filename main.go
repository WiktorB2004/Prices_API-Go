package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB configuration
var (
	mongoDBURL string
	mongoDB    string
)

// Product model
type Product struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"productName" bson:"productName"`
	SupplierID string `json:"supplierId" bson:"supplierId"`
	Price      int    `json:"price" bson:"price"`
}

// Supplier model
type Supplier struct {
	ID            string   `json:"id" bson:"_id"`
	Name          string   `json:"supplierName" bson:"supplierName"`
	Phone         string   `json:"phoneNumber" bson:"phoneNumber"`
	Email         string   `json:"email" bson:"email"`
	Products      []string `json:"products" bson:"products"`
	ProductsCount int      `json:"productsCount" bson:"productsCount"`
}

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoDBURL = os.Getenv("MONGODB_URL")
	mongoDB = os.Getenv("MONGODB_DB")
}

func main() {
	// Set up MongoDB client with authentication options
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Set up Gin router
	router := gin.Default()

	// Endpoint to retrieve product data from MongoDB
	router.GET("/get-product-data", func(c *gin.Context) {
		// Retrieve Product from MongoDB
		collection := client.Database(mongoDB).Collection("products")
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products data from MongoDB"})
			return
		}
		defer cursor.Close(context.Background())

		// Decode MongoDB documents into a slice
		var result []Product
		if err := cursor.All(context.Background(), &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Endpoint to retrieve supplier data from MongoDB
	router.GET("/get-supplier-data", func(c *gin.Context) {
		// Retrieve Suppliers from MongoDB
		collection := client.Database(mongoDB).Collection("suppliers")
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier data from MongoDB"})
			return
		}
		defer cursor.Close(context.Background())

		// Decode MongoDB documents into a slice
		var result []Supplier
		if err := cursor.All(context.Background(), &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Run the server
	serverAddr := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
