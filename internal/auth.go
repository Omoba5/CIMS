package internal

import (
	"cims/models"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./views/*.html")
	if err != nil {
		fmt.Println("Parsing Templates Error:")
		panic(err.Error)
	}

	tmpl.ExecuteTemplate(w, "registration.html", nil)

	// Parse the form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Access form values
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordAgain := r.FormValue("passwordagain")
	companyName := r.FormValue("name")
	emailAddress := r.FormValue("email")
	address := r.FormValue("address")

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	fmt.Println("Password Again:", passwordAgain)
	fmt.Println("Company Name:", companyName)
	fmt.Println("Email Address:", emailAddress)
	fmt.Println("Address:", address)

	if password != passwordAgain {
		tmpl.ExecuteTemplate(w, "registration.html", "Please confirm that the passwords match")
		return
	}

	// create hash from password
	var hash []byte

	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tmpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))
	// func (db *DB) Prepare(query string) (*Stmt, error)

	user := models.User{
		Username:    username,
		Password:    string(hash),
		Email:       emailAddress,
		CompanyName: companyName,
		Address:     address,
	}

	// check if username already exists for availability
	InsertData(user, "cims", username)
}