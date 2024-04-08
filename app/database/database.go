package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Db is a database connection.
var Db *sql.DB

// Err is an error returned.
var Err error

func Init() error {

	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)

	Db, Err = sql.Open("mysql", dsn)
	// Must import "github.com/go-sql-driver/mysql" for mysql driver
	if Err != nil {
		return Err
	} else {
		Cache()
		return nil
	}
}
