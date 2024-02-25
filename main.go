package main

import (
	"log"
	"net/http"
	// "text/template"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// tmpl := template.Must(template.ParseFiles("dashboard.html"))
		// tmpl.Execute(w, r)
		http.ServeFile(w, r, "dashboard.html")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
