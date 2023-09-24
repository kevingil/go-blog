package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/models"
)

func init() {
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
		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to MySQL server: %v\n", err)
		fmt.Printf("Retrying in %v...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		fmt.Printf("Max retries reached, cannot connect to MySQL server: %v\n", err)
		return
	}

	models.Err = models.Db.Ping()
	if models.Err != nil {
		log.Fatal(models.Err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Index)
	r.HandleFunc("/contact", controllers.Contact)
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/register", controllers.Register)
	r.HandleFunc("/dashboard", controllers.Dashboard)
	r.HandleFunc("/post/{slug}", controllers.Post)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	log.Println(fmt.Sprintf("Your app is running on port %s.", os.Getenv("PORT")))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
