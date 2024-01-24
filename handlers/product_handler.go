package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pricesAPI/app"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type productStats struct {
	ProductList []app.ExtProduct `json:"ProductList" bson:"ProductList"`
	AvgPrice    float32          `json:"AvgPrice" bson:"AvgPrice"`
}

func CreateExtProduct(p app.Product, supplierName string) app.ExtProduct {
	return app.ExtProduct{
		ID:           p.ID,
		Name:         p.Name,
		SupplierName: supplierName,
		Price:        p.Price,
	}
}

func GetProductsData(c *gin.Context) {
	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("products")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products data from MongoDB"})
		return
	}
	defer cursor.Close(context.Background())

	var dbResult []app.Product
	if err := cursor.All(context.Background(), &dbResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
		return
	}

	result := make(map[string]productStats)
	// Group products with the same name together
	for _, product := range dbResult {
		currProduct := result[product.Name]

		SupplierReq, err := http.Get(fmt.Sprintf("http://localhost:3000/supplier/%s", product.SupplierID))
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}
		defer SupplierReq.Body.Close()

		if SupplierReq.StatusCode != http.StatusOK {
			fmt.Printf("GET request failed with status code: %d\n", SupplierReq.StatusCode)
			return
		}

		var suppRes app.SupplierResponse
		err = json.NewDecoder(SupplierReq.Body).Decode(&suppRes)
		if err != nil {
			fmt.Println("Error decoding ProductResponse:", err)
			return
		}

		currProduct.ProductList = append(currProduct.ProductList, CreateExtProduct(product, suppRes.ExistingProduct.Name))
		result[product.Name] = currProduct
	}

	// Calculate average price for each group
	for prodName, stats := range result {
		if len(stats.ProductList) > 0 {
			for _, product := range stats.ProductList {
				stats.AvgPrice += float32(product.Price)
			}
			stats.AvgPrice /= float32(len(stats.ProductList))
		}
		result[prodName] = stats
	}

	c.JSON(http.StatusOK, result)
}

func GetProductsByName(c *gin.Context) {
	productName := c.Param("name")

	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("products")

	filter := bson.M{"productName": productName}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products data from MongoDB"})
		return
	}
	defer cursor.Close(context.Background())

	var dbResult []app.Product
	if err := cursor.All(context.Background(), &dbResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
		return
	}

	c.JSON(http.StatusOK, dbResult)
}
