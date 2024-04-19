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

// func InsertData(record any, table string, primarykey string) {

// 	// Find .env file
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %s", err)
// 	}

// 	// Read environment variables
// 	dbHost := os.Getenv("DB_HOST")
// 	dbUser := os.Getenv("DB_USER")
// 	dbPass := os.Getenv("DB_PASS")

// 	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	collection := client.Database("testting").Collection(table)

// 	// Check username availability
// 	var existingUser models.User
// 	err = collection.FindOne(context.Background(), bson.M{"username": primarykey}).Decode(&existingUser)
// 	if err == nil {
// 		// User name already exists, handle the case accordingly
// 		fmt.Println("User name already exists:", existingUser.Username)
// 		return
// 	} else if err != mongo.ErrNoDocuments {
// 		// Error occurred during the query
// 		panic(err)
// 	} else {
// 		// Insert Document
// 		_, err = collection.InsertOne(context.TODO(), record)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// }

// func GetUserByUsername(username string) (*models.User, error) {
// 	// Find .env file
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		return nil, fmt.Errorf("error loading .env file: %s", err)
// 	}

// 	// Read environment variables
// 	dbHost := os.Getenv("DB_HOST")
// 	dbUser := os.Getenv("DB_USER")
// 	dbPass := os.Getenv("DB_PASS")

// 	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("error connecting to MongoDB: %s", err)
// 	}
// 	defer client.Disconnect(context.TODO())

// 	// Access the "users" collection
// 	collection := client.Database("testting").Collection("users")

// 	// Define a filter for finding the user by username
// 	filter := bson.M{"username": username}

// 	// Find the user document matching the filter
// 	var user models.User
// 	err = collection.FindOne(context.Background(), filter).Decode(&user)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			// User not found
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("error finding user: %s", err)
// 	}

// 	return &user, nil
// }

var (
	client *mongo.Client
	ctx    = context.Background()
	table  = "cims"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Set MongoDB connection options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

	// Create a new MongoDB client and connect to the server
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %s", err)
	}
	fmt.Println("Connected to MongoDB!")
}

func InsertData(record interface{}, primaryKey string) error {
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

func GetUser(username string) (*models.User, error) {
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
