package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "." + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// tmpl := template.Must(template.ParseFiles("dashboard.html"))
		// tmpl.Execute(w, r)
		tmpl, err := template.ParseFiles("login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tmpl.Execute(w, nil)
		// http.ServeFile(w, r, "dashboard.html")
	})

	http.HandleFunc("/bestpractices", func(w http.ResponseWriter, r *http.Request) {
		// tmpl := template.Must(template.ParseFiles("dashboard.html"))
		// tmpl.Execute(w, r)
		tmpl, err := template.ParseFiles("best_practices.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tmpl.Execute(w, nil)
		// http.ServeFile(w, r, "dashboard.html")
	})

	http.HandleFunc("/bestpractices/vm", func(w http.ResponseWriter, r *http.Request) {
		// tmpl := template.Must(template.ParseFiles("dashboard.html"))
		// tmpl.Execute(w, r)
		tmpl, err := template.ParseFiles("./cloud_blog/article_VMs.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tmpl.Execute(w, nil)
		// http.ServeFile(w, r, "dashboard.html")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
