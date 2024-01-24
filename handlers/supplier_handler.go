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

type supplierStats struct {
	ID       string  `json:"id" bson:"_id"`
	AvgPrice float32 `json:"AvgPrice" bson:"AvgPrice"`
	Products int     `json:"Products" bson:"Products"`
}

func GetSuppliersData(c *gin.Context) {
	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("suppliers")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier data from MongoDB"})
		return
	}
	defer cursor.Close(context.Background())

	var dbResult []app.Supplier
	if err := cursor.All(context.Background(), &dbResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
		return
	}

	result := make(map[string]supplierStats)
	for _, supplier := range dbResult {
		result[supplier.Name] = supplierStats{supplier.ID, 0, supplier.ProductsCount}
		if supplier.ProductsCount > 0 {
			currSupplier := result[supplier.Name]
			for _, productId := range supplier.Products {
				ProductReq, err := http.Get(fmt.Sprintf("http://localhost:3000/product/%s", productId))
				if err != nil {
					fmt.Println("Error making GET request:", err)
					return
				}
				defer ProductReq.Body.Close()

				if ProductReq.StatusCode != http.StatusOK {
					fmt.Printf("GET request failed with status code: %d\n", ProductReq.StatusCode)
					return
				}

				var prodRes app.ProductResponse
				err = json.NewDecoder(ProductReq.Body).Decode(&prodRes)
				if err != nil {
					fmt.Println("Error decoding ProductResponse:", err)
					return
				}

				product := prodRes.ExistingProduct
				currSupplier.AvgPrice += float32(product.Price)
			}
			currSupplier.AvgPrice /= float32(supplier.ProductsCount)
			result[supplier.Name] = currSupplier
		}
	}

	c.JSON(http.StatusOK, result)
}
