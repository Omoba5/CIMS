package internal

import (
	"cims/models"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template

func init() {
	// Parse templates
	var err error
	tmpl, err = template.ParseGlob("./views/*.html")
	if err != nil {
		panic(err)
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		tmpl.ExecuteTemplate(w, "registration.html", "Failed to Parse Form")
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
	err2 := InsertData(user, username)
	if err2 == nil {
		// You can redirect the user to a dashboard or any other page upon successful login
		tmpl.ExecuteTemplate(w, "login.html", "Registration Successful")
		return
	}

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		tmpl.ExecuteTemplate(w, "login.html", "Failed to parse form")
		// http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Access form values
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println("Username:", username)
	fmt.Println("Password:", password)

	// Retrieve user from the database based on the provided username
	user, err := GetUser(username)
	fmt.Println(user)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	// Check if user exists and verify password
	if user != nil && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		// User exists and password matches, authentication successful
		fmt.Println("Login successful for user:", username)
		// You can redirect the user to a dashboard or any other page upon successful login
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	} else {
		// Invalid username or password
		fmt.Println("Invalid username or password")
		// You can render an error message on the login page
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
}
