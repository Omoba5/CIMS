package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username    string
	Password    string
	Email       string
	CompanyName string
	Address     string
}

type VirtualMachine struct {
	VMName      string
	Username    string
	DateCreated string
	Status      string
	MachineType string
	ExternalIP  string
	InternalIP  string
	NetworkTags string
	Zone        string
}

type Network struct {
	NetworkName   string
	Username      string
	DateCreated   string
	Subnets       string
	MTU           string
	FirewallRules string
}

type FirewallRule struct {
	FwName        string
	Username      string
	DateCreated   string
	Type          string
	Targets       string
	SourceFilter  string
	ProtocolPorts string
	Actions       string
	Priority      string
	Network       string
}

type Subnet struct {
	SubnetName       string
	NetworkName      string
	DateCreated      string
	Region           string
	InternalIPRanges string
	ExternalIPRanges string
	Gateway          string
}

func main() {
	tmpl, err := template.ParseGlob("*.html")
	if err != nil {
		fmt.Println("Parsing Templates Error:")
		panic(err.Error)
	}

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "." + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/cloud_blog/", func(w http.ResponseWriter, r *http.Request) {
		filePath := "." + r.URL.Path
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("***Login Handler running***")
		tmpl.ExecuteTemplate(w, "login.html", nil)
	})

	http.HandleFunc("/register", registerUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("*.html")
	if err != nil {
		fmt.Println("Parsing Templates Error:")
		panic(err.Error)
	}

	fmt.Println("***Register Users on the CIMS platform***")
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

	user := User{
		Username:    username,
		Password:    string(hash),
		Email:       emailAddress,
		CompanyName: companyName,
		Address:     address,
	}

	// check if username already exists for availability
	insertData(user, "cims", username)
}

func insertData(record any, table string, primarykey string) {

	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + url.QueryEscape(dbPass) + dbHost).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testting").Collection(table)

	// Check username availability
	var existingUser User
	err = collection.FindOne(context.Background(), bson.M{"username": primarykey}).Decode(&existingUser)
	if err == nil {
		// User name already exists, handle the case accordingly
		fmt.Println("User name already exists:", existingUser.Username)
		return
	} else if err != mongo.ErrNoDocuments {
		// Error occurred during the query
		panic(err)
	} else {
		// Insert Document
		_, err = collection.InsertOne(context.TODO(), record)
		if err != nil {
			panic(err)
		}
	}

}
