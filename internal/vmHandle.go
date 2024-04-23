package internal

import (
	"cims/internal/resources"
	"cims/models"
	"encoding/json"

	"fmt"
	"net/http"
	"strconv"
	// "sync"
)

var zone = "us-west1-a"

// var deleteMutex sync.Mutex

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
	instance, err := resources.CreateInstance(computeService, projectID, instanceName, zone, machineType, username, password, diskSize)
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

func GetVirtualMachineList(w http.ResponseWriter, r *http.Request) {
	// Marshal the virtual machines to JSON
	vmJSON, err := resources.GetInstances(computeService, projectID)
	if err != nil {
		http.Error(w, "Failed to marshal virtual machines", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(vmJSON)
	// fmt.Println(string(vmJSON))
}

func DeleteVMsHandler(w http.ResponseWriter, r *http.Request) {
	// // Lock the mutex to prevent concurrent delete operations
	// deleteMutex.Lock()
	// defer deleteMutex.Unlock()

	// Parse request body
	var requestBody struct {
		VMs []string `json:"vms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Print the request body as a string
	requestBodyStr, _ := json.Marshal(requestBody)
	fmt.Println("Received request body:", string(requestBodyStr))

	// Delete VMs
	for _, vmName := range requestBody.VMs {
		err := resources.DeleteInstance(computeService, projectID, vmName, zone)
		if err != nil {
			http.Error(w, "Failed to delete virtual machines", http.StatusInternalServerError)
			return
		}
		fmt.Printf("VM %v deleted successfully", vmName)
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "VMs deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
