package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold prepared statements
type PreparedQueries struct {
	GetUserByID *sql.Stmt
}

func userhandle(pq *PreparedQueries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the query parameter
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		// Execute the prepared statement
		var name string
		err := pq.GetUserByID.QueryRow(id).Scan(&id, &name)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Write the response
		fmt.Fprintf(w, "User found: ID=%s, Name=%s\n", id, name)
	}
}

// Function to prepare SQL statements
func PrepareQueries(db *sql.DB) (*PreparedQueries, error) {
	pq := &PreparedQueries{}

	var err error
	pq.GetUserByID, err = db.Prepare("SELECT id, name FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}

	return pq, nil
}

// Function to create the necessary tables
func CreateTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255)
    )`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Function to import fixtures into the users table
func ImportFixtures(db *sql.DB) error {
	fixtures := []string{"Alice", "Bob", "Charlie"}

	for _, name := range fixtures {
		_, err := db.Exec("INSERT INTO users (name) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Connect to the database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(20 * time.Second)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	// db.SetConnMaxIdleTime(time.Second * 10)

	if !true {
		go func() {
			ticker := time.NewTicker(time.Second * 1)
			for {
				spew.Dump(db.Stats())
				<-ticker.C
			}
		}()
	}

	// Create the necessary tables
	if err := CreateTables(db); err != nil {
		log.Fatal(err)
	}

	// Import fixtures into the users table
	if err := ImportFixtures(db); err != nil {
		log.Fatal(err)
	}

	// Prepare SQL statements
	pq, err := PrepareQueries(db)
	if err != nil {
		log.Fatal(err)
	}
	defer pq.GetUserByID.Close()

	// HTTP handler function
	http.HandleFunc("/user", userhandle(pq))

	// Start the HTTP server
	log.Println("Server listening on port 7070...")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
