package handlers

import (
	"context"
	"log"
	"net/http"
	"pricesAPI/app"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetClientData(c *gin.Context) {
	authKey := c.Request.Header.Get("AuthKey")

	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("clients")

	filter := bson.M{"authKey": authKey}

	var dbResult app.Client
	err := collection.FindOne(context.Background(), filter).Decode(&dbResult)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client with provided authKey not found"})
		return
	}

	c.JSON(http.StatusOK, dbResult)
}

func PostClientRegister(c *gin.Context) {
	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("clients")

	var req struct {
		Name string `form:"username"`
	}

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "username cant be empty"})
		return
	}

	var newClient app.Client
	authKey := uuid.New().String()
	apiKey := uuid.New().String()
	newClient.Name = req.Name
	newClient.ApiKeys = make(map[string]int)
	newClient.ApiKeys[apiKey] = 0

	var existingUser app.Client
	err := collection.FindOne(context.Background(), bson.M{"authKey": authKey}).Decode(&existingUser)
	if err == nil {
		// Auth key is not unique, generate a new one (the probability of generating existing one for 2nd time is close to 0)
		authKey = uuid.New().String()
	}

	err = collection.FindOne(context.Background(), bson.M{"username": req.Name}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Found user with provided username, change your usernamename": req.Name})
		return
	}

	// Insert the client into MongoDB
	newClient.AuthKey = authKey
	_, err = collection.InsertOne(context.Background(), newClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to register client": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Client registered successfully", "authKey": authKey})
}

func PostClientKey(c *gin.Context) {
	authKey := c.Request.Header.Get("AuthKey")

	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("clients")
	filter := bson.M{"authKey": authKey}

	// Find a client with authkey provided in request header
	var dbResult app.Client
	err := collection.FindOne(context.Background(), filter).Decode(&dbResult)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client with provided authKey not found"})
		return
	}

	// Add new ApiKey
	apiKey := uuid.New().String()
	update := bson.M{"$set": bson.M{"apiKeys." + apiKey: 0}}

	// Update database document
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Client document update failed"})
		return
	}

	// Check if document got modified successflly and provide user with feedback
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusOK, "No client edited")
	} else {
		c.JSON(http.StatusOK, gin.H{"Client:": dbResult.Name, "Deleted apikey": apiKey})
	}
}

func DeleteClient(c *gin.Context) {
	authKey := c.Request.Header.Get("AuthKey")
	apiKey := c.Query("apikey")

	client := app.GetMongoClient()
	collection := client.Database(app.MongoDB).Collection("clients")
	filter := bson.M{"authKey": authKey}

	if apiKey == "" {
		result, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Fatal(err)
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusOK, "No client deleted")
		} else {
			c.JSON(http.StatusOK, gin.H{"Message": "Deleted client successfully", "Deleted user AuthKey": authKey})
		}
	} else {
		var dbResult app.Client
		if err := collection.FindOne(context.Background(), filter).Decode(&dbResult); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Client with provided authKey not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client data from MongoDB"})
			return
		}

		delete(dbResult.ApiKeys, apiKey)
		update := bson.M{"$unset": bson.M{"apiKeys." + apiKey: ""}}
		result, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MongoDB document"})
			return
		}
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusOK, "No client edited")
		} else {
			c.JSON(http.StatusOK, gin.H{"Client:": dbResult.Name, "Deleted apikey": apiKey})
		}
	}

}
