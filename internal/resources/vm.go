package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func GetInstances(service *compute.Service, projectID string) ([]string, error) {
	// List instances aggregated by zone in the project
	instancesList, err := service.Instances.AggregatedList(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %v", err)
	}

	var allInstances []*compute.Instance

	// Iterate over the zones and collect instances from each zone
	for _, zoneInstances := range instancesList.Items {
		allInstances = append(allInstances, zoneInstances.Instances...)
	}

	var instancelist []string
	for _, instance := range allInstances {
		jsonInstances, err := json.MarshalIndent(instance, "", " ")
		if err != nil {
			log.Fatalf("Error Marshalling instances to json: %v", err)
		}
		instancelist = append(instancelist, string(jsonInstances))
	}

	return instancelist, nil
}

func CreateInstance(service *compute.Service, projectID, instanceName, zone, machineType, username, password string, size int64) (*compute.Instance, error) {
	fmt.Printf("Creating instance %s in project %s\n", instanceName, projectID)

	// Define the startup script
	startupScript := fmt.Sprintf(`#!/bin/bash
useradd -m -s /bin/bash %s
echo "%s:%s" | chpasswd
sed -i 's/^PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config
systemctl reload sshd`, username, username, password)

	// Create an instance resource object with the instance details
	instance := &compute.Instance{
		Name:        instanceName,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType),
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				DiskSizeGb: size,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-10",
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Name: "External NAT",
						Type: "ONE_TO_ONE_NAT",
					},
				},
				Network: "global/networks/default",
			},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "startup-script",
					Value: &startupScript,
				},
			},
		},
	}

	// Call the Instances.Insert method to create the instance
	op, err := service.Instances.Insert(projectID, zone, instance).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %v", err)
	}

	waitZoneOperation(service, projectID, zone, op.Name, "creating")
	fmt.Printf("Instance %s created successfully!\n", instanceName)

	// Wait 1- seconds before removing the startup script
	time.Sleep(10 * time.Second)

	createdVM, err := service.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to get VM details: %v", err)
	}

	// Remove startup script from instance
	err2 := removeStartupScript(service, projectID, instanceName, zone)
	if err2 != nil {
		return createdVM, fmt.Errorf("failed to remove start up script from instance: %v", err)
	}

	return createdVM, nil
}

func DeleteInstance(service *compute.Service, projectID, instanceName, zone string) error {
	fmt.Printf("Deleting instance %s in project %s\n", instanceName, projectID)

	// Call the Instances.Delete method to create the instance
	op, err := service.Instances.Delete(projectID, zone, instanceName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete instance: %v", err)
	}

	waitZoneOperation(service, projectID, zone, op.Name, "deleting")
	fmt.Printf("Instance %s deleted successfully!\n", instanceName)
	return nil
}

func ChangeInstanceState(service *compute.Service, projectID, zone, instanceName, state string) error {
	fmt.Printf("Attempting to %s instance %s in project %s...\n", state, instanceName, projectID)

	// Check the current instance state and perform the desired state change
	switch state {
	case "START":
		opr, err := service.Instances.Start(projectID, zone, instanceName).Do()
		if err != nil {
			return err
		}
		return waitZoneOperation(service, projectID, zone, opr.Name, state+"ing")

	case "STOP":
		opr, err := service.Instances.Stop(projectID, zone, instanceName).Do()
		if err != nil {
			return err
		}
		return waitZoneOperation(service, projectID, zone, opr.Name, state+"ing")
	case "RESTART":
		opr, err := service.Instances.Reset(projectID, zone, instanceName).Do()
		if err != nil {
			return err
		}
		return waitZoneOperation(service, projectID, zone, opr.Name, state+"ing")

	}

	return nil
}

func UpdateBootDiskSize(service *compute.Service, projectID, zone, instanceName string, newDiskSizeGb int64) error {
	fmt.Printf("Updating boot disk size for instance %s in project %s\n", instanceName, projectID)

	// Resize the boot disk
	op, err := service.Disks.Resize(projectID, zone, instanceName, &compute.DisksResizeRequest{
		SizeGb: newDiskSizeGb,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to resize boot disk: %v", err)
	}

	if err := waitZoneOperation(service, projectID, zone, op.Name, "resizing boot disk"); err != nil {
		return err
	}

	fmt.Printf("Boot disk for instance %s updated successfully!\n", instanceName)
	return nil
}

func UpdateMachineType(service *compute.Service, projectID, zone, instanceName, machineType string) error {
	fmt.Printf("Updating machine type for instance %s in project %s\n", instanceName, projectID)

	// Update machine type
	op, err := service.Instances.SetMachineType(projectID, zone, instanceName, &compute.InstancesSetMachineTypeRequest{
		MachineType: fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", projectID, zone, machineType),
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to update machine type: %v", err)
	}

	// Wait for the operation to complete
	if err := waitZoneOperation(service, projectID, zone, op.Name, "updating machine type"); err != nil {
		return err
	}

	fmt.Printf("Machine type for instance %s updated successfully!\n", instanceName)
	return nil
}

// UpdateNetworkTags updates the network tags for the specified instance in a project and zone.
func UpdateNetworkTags(service *compute.Service, projectID, zone, instanceName string, networkTags []string) error {
	fmt.Printf("Updating network tags for instance %s in project %s\n", instanceName, projectID)

	// Get the current instance details to obtain the current fingerprint of the network tags.
	instance, err := service.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve instance details: %v", err)
	}

	// Retrieve the current fingerprint of the network tags.
	currentFingerprint := instance.Tags.Fingerprint

	// Update network tags with the latest fingerprint.
	update := &compute.Tags{
		Items:       networkTags,
		Fingerprint: currentFingerprint,
	}

	// Perform the update.
	op, err := service.Instances.SetTags(projectID, zone, instanceName, update).Do()
	if err != nil {
		// If the error is due to the fingerprint mismatch, retry the operation.
		if isFingerprintMismatch(err) {
			return retryUpdateNetworkTags(service, projectID, zone, instanceName, networkTags, currentFingerprint)
		}
		return fmt.Errorf("failed to update network tags: %v", err)
	}

	// Wait for the operation to complete.
	if err := waitZoneOperation(service, projectID, zone, op.Name, "updating network tags"); err != nil {
		return err
	}

	fmt.Printf("Network tags for instance %s updated successfully!\n", instanceName)
	return nil
}

// isFingerprintMismatch checks whether the error is a fingerprint mismatch error.
func isFingerprintMismatch(err error) bool {
	apiErr, ok := err.(*googleapi.Error)
	return ok && apiErr.Code == 412 && apiErr.Errors[0].Reason == "conditionNotMet"
}

// retryUpdateNetworkTags retries the update operation with the latest fingerprint.
func retryUpdateNetworkTags(service *compute.Service, projectID, zone, instanceName string, networkTags []string, currentFingerprint string) error {
	// Retry the update operation with the latest fingerprint.
	update := &compute.Tags{
		Items:       networkTags,
		Fingerprint: currentFingerprint,
	}

	op, err := service.Instances.SetTags(projectID, zone, instanceName, update).Do()
	if err != nil {
		return fmt.Errorf("failed to retry update network tags: %v", err)
	}

	// Wait for the operation to complete.
	if err := waitZoneOperation(service, projectID, zone, op.Name, "retrying update of network tags"); err != nil {
		return err
	}

	fmt.Printf("Network tags for instance %s updated successfully after retry!\n", instanceName)
	return nil
}

// RemoveStartupScript removes the startup script from the instance's metadata.
func removeStartupScript(service *compute.Service, projectID, instanceName, zone string) error {
	// Get the current instance metadata
	instance, err := service.Instances.Get(projectID, zone, instanceName).Do()
	if err != nil {
		return fmt.Errorf("failed to get instance metadata: %v", err)
	}

	// Remove the startup-script item from metadata
	var updatedItems []*compute.MetadataItems
	for _, item := range instance.Metadata.Items {
		if item.Key != "startup-script" {
			updatedItems = append(updatedItems, item)
		}
	}

	// Update the instance metadata to remove the startup script
	instance.Metadata.Items = updatedItems
	_, err = service.Instances.SetMetadata(projectID, zone, instanceName, instance.Metadata).Do()
	if err != nil {
		return fmt.Errorf("failed to update instance metadata: %v", err)
	}

	fmt.Println("Startup script removed successfully!")
	return nil
}
