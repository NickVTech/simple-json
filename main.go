package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	// CRUD handlers
	router.HandleFunc("/create", createUser).Methods("POST")
	router.HandleFunc("/read/{name}", readUser).Methods("GET")
	router.HandleFunc("/update", updateUser).Methods("PUT")
	router.HandleFunc("/delete/{name}", deleteUser).Methods("DELETE")

	// To test after performing CRUD operations
	router.HandleFunc("/users", getUsers).Methods("GET")

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

	fmt.Fprintln(w, user.Name)
	fmt.Fprintln(w, user.ID)

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
