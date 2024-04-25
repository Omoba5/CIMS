package resources

import (
	"cims/models"
	"encoding/json"
	"fmt"

	"google.golang.org/api/compute/v1"
)

func GetNetworks(service *compute.Service, projectID string, username string) ([]byte, error) {
	// List networks in the project
	networkList, err := service.Networks.List(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %v", err)
	}

	var allNetworks []*models.Network

	// Iterate over the list of networks
	for _, network := range networkList.Items {
		// Convert network to Network struct
		convertedNetwork, err := models.ConvertNetworkToStruct(service, network, username, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert network: %v", err)
		}

		// Append converted network to slice
		allNetworks = append(allNetworks, convertedNetwork)
	}

	// Marshal the slice of JSON objects with indentation
	finalJSON, err := json.MarshalIndent(allNetworks, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling networks to JSON: %v", err)
	}

	return finalJSON, nil
}

func CreateNetwork(service *compute.Service, projectID, networkName string) (*compute.Network, error) {
	fmt.Printf("Creating network %s in project %s\n", networkName, projectID)

	subnet := &compute.Subnetwork{
		Name:        networkName,
		Network:     fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName),
		IpCidrRange: "10.138.0.0/20",
		Region:      "us-west1",
	}

	// Define the network resource
	network := &compute.Network{
		Name:                  networkName,
		AutoCreateSubnetworks: false,
		ForceSendFields:       []string{"AutoCreateSubnetworks"},
	}

	// Perform the network creation
	op, err := service.Networks.Insert(projectID, network).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %v", err)
	}

	// Wait for the operation to complete
	if err := waitGlobalOperation(service, projectID, op.Name, "creating network"); err != nil {
		return nil, err
	}

	// Get the network details of the created network
	createdNetwork, err := service.Networks.Get(projectID, networkName).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get network details: %v", err)
	}

	// Create a subnetwork within the created network
	op, err2 := service.Subnetworks.Insert(projectID, "us-west1", subnet).Do()
	if err2 != nil {
		return nil, fmt.Errorf("failed to create subnetwork: %v", err)
	}

	// Wait for the operation to complete
	if err := waitRegionOperation(service, projectID, "us-west1", op.Name, "creating Subnetwork"); err != nil {
		return nil, err
	}

	fmt.Printf("Network %s created successfully!\n", networkName)
	return createdNetwork, nil
}

// DeleteNetwork deletes the network with the specified name from the specified project.
func DeleteNetwork(service *compute.Service, projectID, networkName string) error {
	fmt.Printf("Deleting network %s from project %s\n", networkName, projectID)

	// List all subnets in the network
	subnets, err := service.Subnetworks.List(projectID, "us-west1").Filter(fmt.Sprintf("network eq '.*%s'", networkName)).Do()
	if err != nil {
		return fmt.Errorf("failed to list subnets for network %s: %v", networkName, err)
	}

	// Print the list of subnets for debugging
	fmt.Printf("Subnets associated with network %s:\n", networkName)
	for _, subnet := range subnets.Items {
		fmt.Printf("  %s\n", subnet.Name)
	}

	// Delete each subnet
	for _, subnet := range subnets.Items {
		fmt.Println("Does this loop even run at all???")
		if err := DeleteSubnet(service, projectID, "us-west1", subnet.Name); err != nil {
			return err
		}
	}

	// Perform the network deletion
	op, err := service.Networks.Delete(projectID, networkName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete network: %v", err)
	}

	// Wait for the operation to complete
	if err := waitGlobalOperation(service, projectID, op.Name, "deleting network"); err != nil {
		return err
	}

	fmt.Printf("Network %s deleted successfully!\n", networkName)
	return nil
}

// DeleteSubnet deletes the subnet with the specified name from the specified project and region.
func DeleteSubnet(service *compute.Service, projectID, region, subnetName string) error {
	fmt.Printf("Deleting subnet %s from project %s, region %s\n", subnetName, projectID, region)

	// Perform the subnet deletion
	op, err := service.Subnetworks.Delete(projectID, region, subnetName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete subnet: %v", err)
	}

	// Wait for the operation to complete
	if err := waitRegionOperation(service, projectID, region, op.Name, "deleting subnet"); err != nil {
		return err
	}

	fmt.Printf("Subnet %s deleted successfully!\n", subnetName)
	return nil
}
