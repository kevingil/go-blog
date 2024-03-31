package database

import (
	"database/sql"
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
		Cache()
		return nil
	}
}
