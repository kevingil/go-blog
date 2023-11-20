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

	maxRetries := 3
	retryInterval := 3 * time.Second
	for i := 0; i < maxRetries; i++ {
		models.Db, models.Err = sql.Open("mysql", dsn)
		fmt.Printf("Connecting to MySQL server\n")

		//Try prod db
		models.Err = models.Db.Ping()
		if models.Err == nil {
			fmt.Printf("Connected to cloud database\n")
			break
		} else {
			// Try test db
			fmt.Printf("Trying test database\n")
			models.Db, models.Err = sql.Open("mysql", os.Getenv("TEST_MYSQL"))
			models.Err = models.Db.Ping()
			if models.Err == nil {
				fmt.Printf("Connected to test database\n")
				break
			}
		}

		fmt.Printf("Failed to connect to any MySQL server: %v\n", models.Err)
		fmt.Printf("Retrying ( %v )\n", retryInterval)
		time.Sleep(retryInterval)
	}

	err = check_db(models.Db)
	if err != nil {
		log.Fatal(err)
	}

}

func check_db(db *sql.DB) error {
	// Check users table
	tableName := "users"
	query := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil // Table exists
	}

	return fmt.Errorf("table '%s' does not exist", tableName)

	/*
		_, err = db.Exec("SET time_zone = '+00:00';")
		if err != nil {
			return err
		}

		// Create articles table
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS articles (
				id int(11) NOT NULL,
				image varchar(255) DEFAULT NULL,
				slug varchar(255) NOT NULL,
				title varchar(60) NOT NULL,
				content text NOT NULL,
				author int(11) NOT NULL,
				created_at datetime NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY slug (slug),
				KEY author (author)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}

		// Create users table
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
				id int(11) NOT NULL,
				name varchar(64) NOT NULL,
				email varchar(320) NOT NULL,
				password varchar(255) NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY email (email)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`)
		if err != nil {
			return err
		}

		// Add indexes for articles table
		_, err = db.Exec(`
			ALTER TABLE articles
			ADD PRIMARY KEY (id),
			ADD UNIQUE KEY slug (slug),
			ADD KEY author (author);
		`)
		if err != nil {
			return err
		}

		// Add indexes for users table
		_, err = db.Exec(`
			ALTER TABLE users
			ADD PRIMARY KEY (id),
			ADD UNIQUE KEY email (email);
		`)
		if err != nil {
			return err
		}

		// Set AUTO_INCREMENT for articles table
		_, err = db.Exec("ALTER TABLE articles MODIFY id int(11) NOT NULL AUTO_INCREMENT;")
		if err != nil {
			return err
		}

		// Set AUTO_INCREMENT for users table
		_, err = db.Exec("ALTER TABLE users MODIFY id int(11) NOT NULL AUTO_INCREMENT;")
		if err != nil {
			return err
		}
	*/

}
