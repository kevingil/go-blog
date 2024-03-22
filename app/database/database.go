package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Db is a database connection.
var Db *sql.DB

// Err is an error returned.
var Err error

func Init() {
	Db, Err = sql.Open("mysql", os.Getenv("PROD_DSN"))
	if Err != nil {
		log.Fatal(Err)
	} else {
		testDb(Db)
	}
}

// inintDb will check for necessary tables and create them if not exists
func testDb(db *sql.DB) {
	testSetup(db)
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

func testSetup(db *sql.DB) {

	// Check if the time zone is already set
	rowsTimezone, err := db.Query("SELECT @@session.time_zone")
	if err != nil {
		log.Fatal(err)
	}
	defer rowsTimezone.Close()

	var dbTimeZone string
	for rowsTimezone.Next() {
		err := rowsTimezone.Scan(&dbTimeZone)
		if err != nil {
			log.Fatal(err)
		}
	}

	// If not, set to PST
	if dbTimeZone == "" {
		_, err := db.Exec("SET time_zone = '-08:00';")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Time zone initialized.")
	} else {
		log.Printf("Time zone is set to: %s\n", dbTimeZone)
	}

	// Check if the SQL_MODE is already set
	rowsMode, err := db.Query("SELECT @@session.sql_mode")
	if err != nil {
		log.Fatal(err)
	}
	defer rowsMode.Close()

	var sqlMode string
	for rowsMode.Next() {
		err := rowsMode.Scan(&sqlMode)
		if err != nil {
			log.Fatal(err)
		}
	}

	// If not, set to NO_AUTO_VALUE_ON_ZERO
	if sqlMode == "" {
		_, err := db.Exec("SET SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO';")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("SQL_MODE initialized.")
	} else {
		log.Printf("SQL_MODE is set to: %s\n", sqlMode)
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
