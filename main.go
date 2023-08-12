package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}

var users = []User{
	{Name: "Nick", ID: "0"},
	{Name: "John", ID: "1"},
	{Name: "Sam", ID: "2"},
}

func main() {

	fmt.Println("Starting the server")
	time.Sleep(time.Second * 1)

	connectDB()
	time.Sleep(time.Second * 1)

	router := mux.NewRouter()

	// CRUD handlers
	fmt.Println("Setting up routes")
	router.HandleFunc("/create", createUser).Methods("POST")
	router.HandleFunc("/read/{name}", readUser).Methods("GET")
	router.HandleFunc("/update", updateUser).Methods("PUT")
	router.HandleFunc("/delete/{name}", deleteUser).Methods("DELETE")

	// To test after performing CRUD operations
	router.HandleFunc("/users", getUsers).Methods("GET")
	time.Sleep(time.Second * 1)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", router)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("There was an error decoding the request")
	}

	users = append(users, user)

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		fmt.Println("There was en error encoding the response")
	}

}

func readUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	name := vars["name"]

	for _, user := range users {
		if user.Name == name {
			err := json.NewEncoder(w).Encode(&user)
			if err != nil {
				fmt.Println("There was an error encoding the user")
			}
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("There was an error decoding the user")
	}

	for index, oldUser := range users {
		if user.Name == oldUser.Name {
			users = append(users[:index], users[index+1:]...)
		}
	}

	users = append(users, user)

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		fmt.Println("There was an error encoding the user")
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	name := mux.Vars(r)["name"]

	for index, oldUser := range users {
		if oldUser.Name == name {
			users = append(users[:index], users[index+1:]...)
		}
	}

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for _, user := range users {
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			fmt.Println("There was an error encoding the user")
		}
	}
}

// Enviroment helper function

func getEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading enviroment variables")
	}

	return os.Getenv(key)

}

// Connecting to DB

var db *sql.DB

func connectDB() {
	var err error
	dsn := getEnv("DSN")

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to PlanetScale!")

	addTestData()

}

func addTestData() {
	for _, user := range users {

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE ID = ?", user.ID).Scan(&count)
		if err != nil {
			log.Println("Failed to check user existence:", err)
			continue
		}

		if count == 0 {
			_, err := db.Exec("INSERT INTO Users (name, ID) VALUES (?, ?)", user.Name, user.ID)
			if err != nil {
				log.Println("Failed to insert user:", err)
			} else {
				fmt.Printf("User %s inserted\n", user.Name)
			}
		} else {
			fmt.Printf("User %s already exists\n", user.Name)
		}
	}

	fmt.Println("Test data synchronization complete")
}
