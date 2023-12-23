package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/controllers"

	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/routes"
)

func main() {

	//Init routes
	routes.Init()
}

func init() {
	//Init blog database
	controllers.Sessions = make(map[string]*models.User)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	models.Db, models.Err = sql.Open("mysql", os.Getenv("PROD_DSN"))
	if models.Err != nil {
		log.Fatal(models.Err)
	} else {
		testDb(models.Db)
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
	rows, err := db.Query("SELECT * FROM users LIMIT 3")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into variables
		var id int
		var name, email, password string
		var about, contact sql.NullString

		err := rows.Scan(&id, &name, &email, &password, &about, &contact)
		if err != nil {
			log.Fatal(err)
		}

		// Check for NULL values
		var aboutValue, contentValue string
		if about.Valid {
			aboutValue = about.String
		} else {
			aboutValue = "Null"
		}

		if contact.Valid {
			contentValue = contact.String
		} else {
			contentValue = "Null"
		}

		// Print the results
		fmt.Printf("User ID: %d\nName: %s\nEmail: %s\nPassword: %s\nAbout: %s\nContact: %s\n\n", id, name, email, password, aboutValue, contentValue)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	articles := testTable(db, "articles")
	if articles != nil {
		log.Print(articles)
	}
	skills := testTable(db, "skills")
	if skills != nil {
		log.Print(skills)
	}
	projects := testTable(db, "projects")
	if projects != nil {
		log.Print(projects)
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
