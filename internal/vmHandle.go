package internal

import (
	"cims/internal/resources"
	"cims/models"
	"fmt"
	"net/http"
	"strconv"
)

func CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now Running the CreateVMHandler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	instanceName := r.FormValue("instanceName")
	machineType := r.FormValue("machine-type")
	diskSizeStr := r.FormValue("size")
	password := r.FormValue("password")

	// Convert diskSizeStr to int64
	diskSize, err := strconv.ParseInt(diskSizeStr, 10, 64)
	if err != nil {
		// Handle the error if conversion fails
		http.Error(w, "Invalid disk size", http.StatusBadRequest)
		return
	}

	fmt.Println(instanceName)
	fmt.Println(machineType)
	fmt.Println(diskSize)
	fmt.Println(username)
	fmt.Println(password)

	// Create the network using networks.CreateNetwork function
	instance, err := resources.CreateInstance(computeService, projectID, instanceName, "us-west1-a", machineType, username, password, diskSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating network: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare data for the model.Network struct
	instanceData, err := models.ConvertVMToStruct(instance, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating network struct: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert data into MongoDB using db_controller.InsertDocument function
	err = InsertDocument(instanceData, instanceName, "virtual_machines")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting Virtual Machine data: %v", err), http.StatusInternalServerError)
		return
	}

	// // Set a temporary success message attribute on the body element
	// w.Header().Set("HX-Trigger", "body.showMessage afterbegin")
	// w.Header().Set("HX-Show", "#vm-message 'VM created successfully!'")

	// Redirect to networks_subnets.html for refresh (implement as needed)
	http.Redirect(w, r, "/virtual_machines.html", http.StatusSeeOther)
}
