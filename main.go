package main

import (
	"cims/internal"
	// "cims/internal/resources"
	"cims/internal/ssh"

	"fmt"
	"html/template"
	"net/http"
)

// var tmpl *template.Template

func main() {
	// // Initialize the compute.Service variable from the resource package
	// resources.Init()

	tmpl, err := template.ParseGlob("./views/*.html")
	if err != nil {
		fmt.Println("Parsing Templates Error:")
		panic(err.Error)
	}

	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./internal/ssh/static/css/"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./internal/ssh/static/js/"))))
	http.HandleFunc("/ssh", ssh.Home)
	http.HandleFunc("/ws/v1", ssh.WsHandle)

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "." + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	})
	http.HandleFunc("/login", internal.LoginUser)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "registration.html", nil)
	})
	http.HandleFunc("/registerauth", internal.RegisterUser)

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the username from the cookie
		_, err := r.Cookie("username")
		if err != nil {
			// If the cookie is not found, redirect to the login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tmpl.ExecuteTemplate(w, "dashboard.html", nil)
	})

	http.HandleFunc("/infastructureOverview", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "infastructure_overview.html", nil)
	})
	http.HandleFunc("/bestPractices", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "best_practices.html", nil)
	})
	http.HandleFunc("/cloudEdu", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "cloud_edu.html", nil)
	})

	// Virtual Machine Handlers
	http.HandleFunc("/virtualMachines", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "virtual_machines.html", nil)
	})
	http.HandleFunc("/createVM", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "createVM.html", nil)
	})
	http.HandleFunc("/vmauth", internal.CreateVMHandler)

	// Firewall Rules Handlers
	http.HandleFunc("/firewallRules", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "firewall_rules.html", nil)
	})
	http.HandleFunc("/createFirewall", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "createFirewall.html", nil)
	})

	// Network Handlers
	http.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "networks_subnets.html", nil)
	})
	http.HandleFunc("/createNetwork", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "createNetwork.html", nil)
	})
	http.HandleFunc("/networkauth", internal.CreateNetworkHandler)

	http.ListenAndServe(":8080", nil)
}
