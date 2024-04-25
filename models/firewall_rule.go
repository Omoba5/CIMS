package models

import (
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

type FirewallRule struct {
	FwName        string
	Username      string
	DateCreated   string
	Type          string
	Targets       []string
	SourceFilter  []string
	ProtocolPorts []ProtocolPorts
	Action        string
	Priority      int64
	Network       string
}

type ProtocolPorts struct {
	Protocols string
	Ports     []string
}

func ConvertFirewallToStruct(firewall *compute.Firewall, username, projectID string) (*FirewallRule, error) {
	convertedFirewall := &FirewallRule{
		FwName:       firewall.Name,
		Username:     username,
		Type:         firewall.Direction,
		Targets:      firewall.TargetTags,
		SourceFilter: firewall.SourceRanges,
		Action:       "Allowed",
		Priority:     firewall.Priority,
		Network:      firewall.Network,
	}

	// Parse creation timestamp and format the date
	if firewall.CreationTimestamp != "" {
		t, err := time.Parse(time.RFC3339Nano, firewall.CreationTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse creation timestamp: %v", err)
		}
		convertedFirewall.DateCreated = t.Format(time.RFC3339) // Adjust format as needed
	}
	// Convert the firewall's allowed protocols and ports
	for _, allowed := range firewall.Allowed {
		protocolPorts := ProtocolPorts{
			Protocols: allowed.IPProtocol,
			Ports:     allowed.Ports,
		}
		convertedFirewall.ProtocolPorts = append(convertedFirewall.ProtocolPorts, protocolPorts)
	}

	return convertedFirewall, nil
}
