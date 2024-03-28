package main

import (
	"cims/internal"
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
