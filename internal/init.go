package internal

import (
	"cims/internal/resources"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/compute/v1"
)

var (
	ctx            context.Context // Global context (optional)
	projectID      string
	computeService *compute.Service
	client         *mongo.Client // Optional: MongoDB client
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Read environment variables
	projectID = os.Getenv("projectID")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Initialize Compute Service
	resources.Init()
	computeService = resources.ComputeService

	// Initialize MongoDB connection (optional)
	if dbHost != "" {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Error connecting to MongoDB: %s", err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Error pinging MongoDB: %s", err)
		}
		fmt.Println("Connected to MongoDB!")
	} else {
		fmt.Println("MongoDB connection not configured (skipping)")
	}
}
