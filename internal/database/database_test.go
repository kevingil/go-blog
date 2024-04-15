package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabase(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("PROD_DSN"))
	if err != nil {
		t.Fatalf("Error opening database connection: %v", err)
	} else {
		t.Log("Connected to database.")
	}
	defer db.Close()
}
