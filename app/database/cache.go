package database

import (
	"log"

	"github.com/kevingil/blog/app/pkg/redix"
)

var testTable *redix.Table

func testCache() {
	// Testing cache in memory database
	redixInstance := redix.New()

	var err error
	testTable, err = redixInstance.CreateTable("test")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	err = testTable.Set("hello", []byte("world"))
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	value, err := testTable.Get("hello")
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}

	log.Printf("Database name: %s\n", testTable.Name)
	log.Printf("value test, hello: %s\n", value)
}
