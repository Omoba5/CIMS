package internal

import (
	"cims/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDocument(record interface{}, primaryKey string, table string) error {
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

func GetNetwork(username, networkName, table string) ([]*models.Network, error) {
	// Access the MongoDB collection
	collection := client.Database("testting").Collection(table)

	// Define filter for finding user by username
	filter := bson.M{"username": username, "networkname": networkName}

	// Find network documents matching the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error finding networks:", err)
		return nil, fmt.Errorf("error finding networks: %s", err)
	}
	defer cursor.Close(ctx)

	var networks []*models.Network

	// Iterate over the cursor and decode each document
	for cursor.Next(ctx) {
		var network models.Network
		if err := cursor.Decode(&network); err != nil {
			fmt.Println("Error decoding network:", err)
			return nil, fmt.Errorf("error decoding network: %s", err)
		}
		networks = append(networks, &network)
	}

	// Check if any networks were found
	if len(networks) == 0 {
		// No networks found
		return nil, nil
	}

	return networks, nil
}

func GetFirewall(username, firewallName, table string) ([]*models.FirewallRule, error) {
	// Access the MongoDB collection
	collection := client.Database("testting").Collection(table)

	// Define filter for finding user by username
	filter := bson.M{"username": username, "fwname": firewallName}

	// Find network documents matching the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error finding firewall rules:", err)
		return nil, fmt.Errorf("error finding firewall rules: %s", err)
	}
	defer cursor.Close(ctx)

	var firewalls []*models.FirewallRule

	// Iterate over the cursor and decode each document
	for cursor.Next(ctx) {
		var firewall models.FirewallRule
		if err := cursor.Decode(&firewall); err != nil {
			fmt.Println("Error decoding firewall rules:", err)
			return nil, fmt.Errorf("error decoding firewall rules: %s", err)
		}
		firewalls = append(firewalls, &firewall)
	}

	// Check if any networks were found
	if len(firewalls) == 0 {
		// No networks found
		return nil, nil
	}

	return firewalls, nil
}

func GetVirtualMachine(username, instanceName, table string) ([]*models.VirtualMachine, error) {
	// Access the MongoDB collection
	collection := client.Database("testting").Collection(table)

	// Define filter for finding user by username
	filter := bson.M{"username": username, "instanceName": instanceName}

	// Find network documents matching the filter
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error finding virtual_machines:", err)
		return nil, fmt.Errorf("error finding virtual_machines: %s", err)
	}
	defer cursor.Close(ctx)

	var virtual_machines []*models.VirtualMachine

	// Iterate over the cursor and decode each document
	for cursor.Next(ctx) {
		var virtual_machine models.VirtualMachine
		if err := cursor.Decode(&virtual_machine); err != nil {
			fmt.Println("Error decoding virtual machines:", err)
			return nil, fmt.Errorf("error decoding virtual machines: %s", err)
		}
		virtual_machines = append(virtual_machines, &virtual_machine)
	}

	// Check if any networks were found
	if len(virtual_machines) == 0 {
		// No networks found
		return nil, nil
	}

	return virtual_machines, nil
}
