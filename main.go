package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/database"
	"github.com/kevingil/blog/internal/server"
)

// Entrypoint

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Every app instance
	// initializes a database connection
	err = database.Init()
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}

	testDb(database.Db)

	// Start HTTP server
	controllers.Store = session.New()
	server.Boot()

}

// inintDb will check for necessary tables and create them if not exists
func testDb(db *sql.DB) {
	users := testTable(db, "users")
	if users != nil {
		log.Print(users)
	}

	// Query the first three rows from the users table
	rows, err := db.Query("SELECT id, name, email FROM users LIMIT 3")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into variables
		var id int
		var name, email string

		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal(err)
		}
		// Print the results
		fmt.Printf("User ID: %d\nName: %s\nEmail: %s\n\n", id, name, email)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// testTable will test if a table exists
func testTable(db *sql.DB, name string) error {
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", name)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Printf("'%s' table OK \n", name)
		return nil // Table exists
	}
	error := "'" + name + "' table ERROR \n"
	return fmt.Errorf(error)
}
