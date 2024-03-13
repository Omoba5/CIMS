package internal

import (
	"cims/models"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func InsertData(record any, table string, primarykey string) {

	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testting").Collection(table)

	// Check username availability
	var existingUser models.User
	err = collection.FindOne(context.Background(), bson.M{"username": primarykey}).Decode(&existingUser)
	if err == nil {
		// User name already exists, handle the case accordingly
		fmt.Println("User name already exists:", existingUser.Username)
		return
	} else if err != mongo.ErrNoDocuments {
		// Error occurred during the query
		panic(err)
	} else {
		// Insert Document
		_, err = collection.InsertOne(context.TODO(), record)
		if err != nil {
			panic(err)
		}
	}

}
