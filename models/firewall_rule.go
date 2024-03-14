package models

type FirewallRule struct {
	FwName        string
	Username      string
	DateCreated   string
	Type          string
	Targets       string
	SourceFilter  string
	ProtocolPorts string
	Actions       string
	Priority      string
	Network       string
}