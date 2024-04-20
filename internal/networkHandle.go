package internal

import (
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

	// Insert data into MongoDB using db_controller.InsertData function
	err = InsertData(networkData, networkName, "networks")
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
