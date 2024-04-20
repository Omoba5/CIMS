package internal

import (
	"cims/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertData(record interface{}, primaryKey string, table string) error {
	// Access the MongoDB collection
	collection := client.Database("testting").Collection(table)

	// Check if the primary key already exists
	filter := bson.M{"username": primaryKey}
	var existingUser models.User
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		// Primary key already exists
		fmt.Println("User name already exists:", existingUser.Username)
		return err
	} else if err != mongo.ErrNoDocuments {
		// Error occurred during the query
		log.Fatalf("Error checking primary key: %s", err)
		return err
	}

	// Insert document
	_, err = collection.InsertOne(ctx, record)
	if err != nil {
		log.Fatalf("Error inserting document: %s", err)
		return err
	}
	return nil
}

func GetUser(username string, table string) (*models.User, error) {
	// Access the MongoDB collection
	collection := client.Database("testting").Collection(table)

	// Define filter for finding user by username
	filter := bson.M{"username": username}

	// Find user document matching the filter
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found
			return nil, nil
		}
		fmt.Println("Error finding user:", err)
		return nil, fmt.Errorf("error finding user: %s", err)
	}

	return &user, nil
}
