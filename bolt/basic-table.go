package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"

	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	// Attempt to setup the database.
	setupDB()

	// Look for cli parameters.
	arguments := os.Args

	// If we have no arguments, return an error.
	if len(arguments) == 1 {
		log.Fatal("no arguments provided")
	}

	// Get the first argument.
	command := arguments[1]

	// If the command is create, create a user.
	if command == "create" {
		createUserHandler()
	}
	if command == "list" {
		listUsers()
	}

}

// list all users.
func listUsers() {
	// Get DB
	db := getDB()

	// Get all users.
	err := db.View(func(tx *bolt.Tx) error {
		// Get the users bucket.
		b := tx.Bucket([]byte("users"))

		// Create a cursor.
		c := b.Cursor()

		// Iterate over the users.
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// Create a new user.
			var user User

			// Unmarshal the user.
			err := json.Unmarshal(v, &user)
			if err != nil {
				return fmt.Errorf("could not unmarshal user: %w", err)
			}

			// Print the user.
			fmt.Println(user)
		}

		// Return nil to indicate that the users were listed successfully.
		return nil
	})

	// Log any errors.
	if err != nil {
		log.Fatal(err)
	}

	// Close the database.
	defer db.Close()

	// Log a message to indicate that the users were listed successfully.
	log.Println("users listed successfully")
}

// Create a user handler.
func createUserHandler() {
	// Get the 2nd cli argument.
	firstName := os.Args[2]

	// if we dont have a first name, return an error.
	if firstName == "" {
		log.Fatal("first name is required")
	}

	// Get the 3rd cli argument.
	lastName := os.Args[3]

	// if we dont have a last name, return an error.
	if lastName == "" {
		log.Fatal("last name is required")
	}

	// Create a new user.
	user := User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
	}

	// Marshal the user struct into json
	buf, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("could not marshal user: %v", err)
	}

	// Print the user.
	fmt.Println(string(buf))

	// Add the user to the database.
	err = addUser(user)
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
}
func getDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	return db
}

func setupDB() {
	db := getDB()

	// Create a bucket.
	err := db.Update(func(tx *bolt.Tx) error {
		// Create the users bucket, if it doesn't exist.
		_, err := tx.CreateBucketIfNotExists([]byte("users"))

		// Log any errors.
		if err != nil {
			return fmt.Errorf("could not create users bucket: %w", err)
		}

		// Return nil to indicate that the bucket was created successfully.
		return nil
	})

	// Log any errors.
	if err != nil {
		log.Fatal(err)
	}

	// Close the database.
	defer db.Close()

	// Log a message to indicate that the database was created successfully.
	log.Println("database created successfully")
}

// Add a user to the database.
func addUser(user User) error {
	// Get new database connection.
	db := getDB()

	err := db.Update(func(tx *bolt.Tx) error {
		// Get the users bucket.
		b := tx.Bucket([]byte("users"))

		// Marshal the user struct into a byte slice.
		buf, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("could not marshal user: %w", err)
		}

		// Put the user in the bucket.
		err = b.Put([]byte(fmt.Sprint(user.ID)), buf)
		if err != nil {
			return fmt.Errorf("could not put user in bucket: %w", err)
		}

		// Return nil to indicate that the user was added successfully.
		return nil
	})

	// Log any errors.
	if err != nil {
		log.Fatal(err)
	}

	// Close the database.
	defer db.Close()

	// Log a message to indicate that the user was added successfully.
	log.Println("user added successfully")

	// Return nil to indicate that the user was added successfully.
	return nil

}
