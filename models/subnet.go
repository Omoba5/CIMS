package models

type Subnet struct {
	SubnetName       string
	NetworkName      string
	DateCreated      string
	Region           string
	InternalIPRanges string
	ExternalIPRanges string
	Gateway          string
}