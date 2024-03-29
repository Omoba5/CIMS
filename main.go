package main

import (
	"cims/internal/ssh"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
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
		tmpl.ExecuteTemplate(w, "registration.html", nil)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "registration.html", nil)
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "dashboard.html", nil)
	})

	http.HandleFunc("/virtualMachines", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "virtual_machines.html", nil)
	})

	http.HandleFunc("/infastructureOverview", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "infastructure_overview.html", nil)
	})

	http.HandleFunc("/network&subnets", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "network_subnets.html", nil)
	})

	http.HandleFunc("/network&firewallRules", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "firewall_rules.html", nil)
	})

	http.HandleFunc("/bestPractices", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "best_practices.html", nil)
	})

	http.HandleFunc("/cloudEdu", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "cloud_edu.html", nil)
	})

	http.ListenAndServe(":8080", nil)
}
