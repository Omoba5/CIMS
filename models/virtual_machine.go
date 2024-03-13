package models

type VirtualMachine struct {
	VMName      string
	Username    string
	DateCreated string
	Status      string
	MachineType string
	ExternalIP  string
	InternalIP  string
	NetworkTags string
	Zone        string
}