package resources

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

var ComputeService *compute.Service

func Init() {
	fmt.Println("Initializing compute service")

	if ComputeService == nil {
		var err error
		ComputeService, err = compute.NewService(context.Background()) // Use background context
		fmt.Println("Initialization complete")
		fmt.Println("")
		if err != nil {
			fmt.Printf("failed to create compute service: %v", err)
		}
	}
}
