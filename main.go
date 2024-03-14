package main

import (
	"cims/internal"
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

	// http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
	// 	filePath := "." + r.URL.Path
	// 	http.ServeFile(w, r, filePath)
	// })

	http.HandleFunc("/blogs/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "./assets/" + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "registration.html", nil)
	})

	http.HandleFunc("/register", internal.RegisterUser)

	http.ListenAndServe(":8080", nil)
}
