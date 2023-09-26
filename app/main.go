package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/router"
)

func main() {
	router.Init()
}

func init() {
	controllers.Sessions = make(map[string]*models.User)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	/*
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)
	*/
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		models.Db, models.Err = sql.Open("mysql", "root:LunaSynnax_k0990@tcp(127.0.0.1:3306)/blog")
		// models.Db, models.Err = sql.Open("mysql", dsn)
		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to MySQL server: %v\n", err)
		fmt.Printf("Retrying in %v...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		fmt.Printf("Cannot connect to MySQL server: %v\n", err)
		return
	}

	models.Err = models.Db.Ping()
	if models.Err != nil {
		log.Fatal(models.Err)
	}
}
