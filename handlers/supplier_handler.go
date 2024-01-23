package handlers

import (
	"context"
	"net/http"

	"pricesAPI/app"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSupplierData(c *gin.Context) {
	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("suppliers")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve supplier data from MongoDB"})
		return
	}
	defer cursor.Close(context.Background())

	var result []app.Supplier
	if err := cursor.All(context.Background(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB documents"})
		return
	}

	c.JSON(http.StatusOK, result)
}
