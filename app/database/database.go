package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Db is a database connection.
var Db *sql.DB

// Err is an error returned.
var Err error

func Init() error {
	Db, Err = sql.Open("mysql", os.Getenv("PROD_DSN"))
	// Must import "github.com/go-sql-driver/mysql" for mysql driver
	if Err != nil {
		return Err
	} else {
		testDb(Db)
		testCache()
		return nil
	}
}

func testDb(db *sql.DB) {

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
