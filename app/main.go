package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/router"
)

func main() {

	//Init router
	router.Init()
}

func init() {
	//Init blog database
	controllers.Sessions = make(map[string]*models.User)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	maxRetries := 5
	retryInterval := 5 * time.Second
	for i := 0; i < maxRetries; i++ {
		models.Db, models.Err = sql.Open("mysql", dsn)
		fmt.Printf("Connecting to MySQL server\n")
		models.Err = models.Db.Ping()

		//Try prod db
		if models.Err == nil {
			fmt.Printf("Connected to MySQL database in container cluster\n")
			break
		} else {
			// Try test db
			fmt.Printf("Trying test database\n")
			models.Db, models.Err = sql.Open("mysql", os.Getenv("TEST_MYSQL"))
			models.Err = models.Db.Ping()
			if models.Err == nil {
				fmt.Printf("Connected to test MySQL database\n")
				break
			}
		}

		fmt.Printf("Failed to connect to any MySQL server: %v\n", models.Err)
		fmt.Printf("Retrying ( %v )\n", retryInterval)
		time.Sleep(retryInterval)
	}
}
