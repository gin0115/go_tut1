package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {

	// Attempt to setup the database.
	setupDB()

	// Check for CLI arguments.
	if len(os.Args) < 2 {
		log.Fatal("no arguments provided")
	}

	// Get the command and handle it with a switch.
	switch command := os.Args[1]; command {
	case "create":
		createUserHandler(os.Args[2], os.Args[3])
	case "list":
		listUsers()
	default:
		log.Fatalf("unknown command: %s", command)
	}

}

// setupDB sets up the database.
func setupDB() {
	// Open the database.
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Close any open connections.
	err = db.Close()

	if err != nil {
		log.Fatal(err)
	}
}

// Get the current database.
func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// createUserHandler creates a user.
func createUserHandler(first, last string) {

	// Get the database.
	db := getDB()

	// Insert the user.
	result, err := db.Exec("INSERT INTO users (first_name, last_name) VALUES (?, ?)", first, last)
	if err != nil {
		log.Fatal(err)
	}

	// Get the ID of the last inserted row
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to retrieve last insert ID: %v\n", err)
	}
	fmt.Printf("Row created with ID: %d\n", lastInsertID)

}

// listUsers lists all users in a table.
func listUsers() {

	// Get the database.
	db := getDB()

	// Query the database.
	rows, err := db.Query("SELECT id, first_name, last_name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows.
	for rows.Next() {
		var id int
		var first, last string
		err = rows.Scan(&id, &first, &last)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, First Name: %s, Last Name: %s\n", id, first, last)
	}

}
