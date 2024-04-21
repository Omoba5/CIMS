package models

import (
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

type VirtualMachine struct {
	VMName      string
	Username    string
	DateCreated string
	Status      string
	MachineType string
	ExternalIP  string
	InternalIP  string
	NetworkTags []string
	Zone        string
}

func ConvertVMToStruct(instance *compute.Instance, username string) (*VirtualMachine, error) {
	convertedVM := &VirtualMachine{
		VMName:      instance.Name,
		Username:    username, // Set username as specified
		Status:      instance.Status,
		MachineType: instance.MachineType,
		NetworkTags: instance.Tags.Items,
		Zone:        instance.Zone,
	}

	// Access the internal and external IP addresses (if available)
	internalIP := ""
	if instance.NetworkInterfaces != nil && len(instance.NetworkInterfaces) > 0 {
		internalIP = instance.NetworkInterfaces[0].NetworkIP
	}
	externalIP := ""
	if instance.NetworkInterfaces[0].AccessConfigs != nil && len(instance.NetworkInterfaces[0].AccessConfigs) > 0 {
		externalIP = instance.NetworkInterfaces[0].AccessConfigs[0].NatIP
	}
	convertedVM.ExternalIP = externalIP
	convertedVM.InternalIP = internalIP

	// Parse creation timestamp and format the date
	if instance.CreationTimestamp != "" {
		t, err := time.Parse(time.RFC3339Nano, instance.CreationTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse creation timestamp: %v", err)
		}
		convertedVM.DateCreated = t.Format(time.RFC3339) // Adjust format as needed
	}

	return convertedVM, nil
}
