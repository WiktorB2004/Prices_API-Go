package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Client model
type Client struct {
	ID      string         `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string         `json:"username" bson:"username"`
	AuthKey string         `json:"authKey" bson:"authKey"`
	ApiKeys map[string]int `json:"apiKeys" bson:"apiKeys"`
}

func IncrementAndAuthenticate(c *gin.Context, client *mongo.Client) bool {
	authKey := c.Request.Header.Get("AuthKey")
	apiKey := c.Query("apikey")
	collection := client.Database(MongoDB).Collection("clients")

	filter := bson.M{"authKey": authKey}

	var clientData Client
	err := collection.FindOne(context.Background(), filter).Decode(&clientData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client with provided authKey not found, please register using /register"})
		return false
	}

	_, exist := clientData.ApiKeys[apiKey]
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect apikey provided"})
		return false
	}
	update := bson.M{"$inc": bson.M{"apiKeys." + apiKey: 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit client data"})
		return false
	}
	return true
}
