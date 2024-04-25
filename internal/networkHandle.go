package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cims/internal/resources"
	"cims/models"
)

func CreateNetworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the username from the cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		// If the cookie is not found, redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get the username from the cookie
	username := cookie.Value
	fmt.Println("Cookie Value from CreateNetworkHandler: ", username)
	fmt.Println("")

	// Parse the form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	networkName := username + "-" + r.FormValue("networkName")

	// Create the network using networks.CreateNetwork function
	network, err := resources.CreateNetwork(computeService, projectID, networkName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating network: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare data for the model.Network struct
	networkData, err := models.ConvertNetworkToStruct(computeService, network, username, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating network struct: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert data into MongoDB using db_controller.InsertDocument function
	err = InsertDocument(networkData, networkName, "networks")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting network data: %v", err), http.StatusInternalServerError)
		return
	}

	// Set a temporary success message attribute on the body element
	w.Header().Set("HX-Trigger", "body.showMessage afterbegin")
	w.Header().Set("HX-Show", "#network-message 'Network created successfully!'")

	// Redirect to networks_subnets.html for refresh (implement as needed)
	http.Redirect(w, r, "/networks_subnets.html", http.StatusSeeOther)
}

func GetNetworkList(w http.ResponseWriter, r *http.Request) {
	// Marshal the networks to JSON
	networkJSON, err := resources.GetNetworks(computeService, projectID, "")
	if err != nil {
		http.Error(w, "Failed to marshal networks", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(networkJSON)
	fmt.Println(string(networkJSON))
}

func DeleteNetworksHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		Networks []string `json:"networks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Print the request body as a string
	requestBodyStr, _ := json.Marshal(requestBody)
	fmt.Println("Received request body:", string(requestBodyStr))

	// Delete VMs
	for _, networkName := range requestBody.Networks {
		err := resources.DeleteNetwork(computeService, projectID, networkName)
		if err != nil {
			http.Error(w, "Failed to delete networks", http.StatusInternalServerError)
			return
		}
		fmt.Printf("VM %v deleted successfully", networkName)
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Networks deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
