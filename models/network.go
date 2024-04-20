package models

import (
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

type Network struct {
	NetworkName   string
	Username      string
	DateCreated   string
	Subnet        int
	MTU           int64
	FirewallRules int
}

func ConvertNetworkToStruct(service *compute.Service, network *compute.Network, username, projectID string) (*Network, error) {
	convertedNetwork := &Network{
		NetworkName: network.Name,
		Username:    username, // Set username as specified
		// ...
	}

	// Parse creation timestamp and format the date
	if network.CreationTimestamp != "" {
		t, err := time.Parse(time.RFC3339Nano, network.CreationTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse creation timestamp: %v", err)
		}
		convertedNetwork.DateCreated = t.Format(time.RFC3339) // Adjust format as needed
	}

	// Count subnetworks
	subnetCount := 0
	if network.Subnetworks != nil {
		subnetCount = len(network.Subnetworks)
	}
	convertedNetwork.Subnet = subnetCount

	// Count firewall rules (assuming Firewall field exists)
	firewallCount := 0
	filter := fmt.Sprintf("network eq '.*%s'", network.Name)
	list, err := service.Firewalls.List(projectID).Filter(filter).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list firewall rules: %v", err)
	}
	firewallCount = len(list.Items)

	convertedNetwork.FirewallRules = firewallCount

	// ... (set other fields based on compute.Network properties)

	return convertedNetwork, nil
}
