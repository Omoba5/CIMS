package resources

import (
	"fmt"

	"google.golang.org/api/compute/v1"
)

func CreateFirewallRule(service *compute.Service, projectID, firewallName string, portMap map[string][]string, targets []string) error {
	fmt.Printf("Creating firewall rule %s in project %s\n", firewallName, projectID)

	// Create a firewall resource object with the firewall details
	firewall := &compute.Firewall{
		Name:         firewallName,
		Allowed:      convertPortMapToAllowed(portMap),
		SourceRanges: []string{"0.0.0.0/0"},
		TargetTags:   targets,
	}

	// Call the Firewalls.Insert method to create the firewall rule
	op, err := service.Firewalls.Insert(projectID, firewall).Do()
	if err != nil {
		return fmt.Errorf("failed to create firewall rule: %v", err)
	}

	waitGlobalOperation(service, projectID, op.Name, "creating")
	fmt.Printf("Firewall rule %s created successfully!\n", firewallName)
	return nil
}

// UpdateFirewallPorts updates the ports allowed by a firewall rule.
func UpdateFirewallPorts(service *compute.Service, projectID, firewallName string, portMap map[string][]string) error {
	fmt.Printf("Updating firewall rule %s ports in project %s\n", firewallName, projectID)

	// Retrieve the current firewall rule details
	firewall, err := service.Firewalls.Get(projectID, firewallName).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve firewall details: %v", err)
	}

	// Update the allowed ports in the firewall rule based on the port map
	firewall.Allowed = convertPortMapToAllowed(portMap)

	// Perform the update
	op, err := service.Firewalls.Update(projectID, firewallName, firewall).Do()
	if err != nil {
		return fmt.Errorf("failed to update firewall rule: %v", err)
	}

	// Wait for the operation to complete
	if err := waitGlobalOperation(service, projectID, op.Name, "updating firewall rule"); err != nil {
		return err
	}

	fmt.Printf("Firewall rule %s ports updated successfully!\n", firewallName)
	return nil
}

func UpdateFirewallTargetTags(service *compute.Service, projectID, firewallName string, targetTags []string) error {
	fmt.Printf("Updating target tags for firewall rule %s in project %s\n", firewallName, projectID)

	// Retrieve the current firewall rule details
	firewall, err := service.Firewalls.Get(projectID, firewallName).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve firewall details: %v", err)
	}

	// Update the target tags in the firewall rule
	firewall.TargetTags = targetTags

	// Perform the update
	op, err := service.Firewalls.Update(projectID, firewallName, firewall).Do()
	if err != nil {
		return fmt.Errorf("failed to update firewall rule: %v", err)
	}

	// Wait for the operation to complete
	if err := waitGlobalOperation(service, projectID, op.Name, "updating firewall rule"); err != nil {
		return err
	}

	fmt.Printf("Target tags for firewall rule %s updated successfully!\n", firewallName)
	return nil
}

func UpdateFirewallSourceRanges(service *compute.Service, projectID, firewallName string, soruceRanges []string) error {
	fmt.Printf("Updating Source Ranges for firewall rule %s in project %s\n", firewallName, projectID)

	// Retrieve the current firewall rule details
	firewall, err := service.Firewalls.Get(projectID, firewallName).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve firewall details: %v", err)
	}

	// Update the Source Ranges in the firewall rule
	firewall.SourceRanges = soruceRanges

	// Perform the update
	op, err := service.Firewalls.Update(projectID, firewallName, firewall).Do()
	if err != nil {
		return fmt.Errorf("failed to update firewall rule: %v", err)
	}

	// Wait for the operation to complete
	if err := waitGlobalOperation(service, projectID, op.Name, "updating firewall rule"); err != nil {
		return err
	}

	fmt.Printf("Source Ranges for firewall rule %s updated successfully!\n", firewallName)
	return nil
}

// convertPortMapToAllowed converts a map of ports to compute.Allowed types.
func convertPortMapToAllowed(portMap map[string][]string) []*compute.FirewallAllowed {
	var allowed []*compute.FirewallAllowed
	for protocol, ports := range portMap {
		for _, port := range ports {
			allowed = append(allowed, &compute.FirewallAllowed{
				IPProtocol: protocol,
				Ports:      []string{port},
			})
		}
	}
	return allowed
}

func DeleteFirewallRule(service *compute.Service, projectID, firewallName string) error {
	fmt.Printf("Deleting firewall rule %s in project %s\n", firewallName, projectID)

	// Call the Firewalls.Delete method to delete the firewall rule
	op, err := service.Firewalls.Delete(projectID, firewallName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete firewall rule: %v", err)
	}

	waitGlobalOperation(service, projectID, op.Name, "deleting")
	fmt.Printf("Firewall rule %s deleted successfully!\n", firewallName)
	return nil
}
