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

	initDb(models.Db)
}

// inintDb will check for necessary tables and create them if not exists
func initDb(db *sql.DB) {
	testSetup(db)
	users := testTable(db, "users")
	if users != nil {
		// Create users table
		_, err := db.Exec(`
					CREATE TABLE IF NOT EXISTS users (
					id int(11) NOT NULL AUTO_INCREMENT,
					name varchar(64) NOT NULL,
					email varchar(320) NOT NULL,
					password varchar(255) NOT NULL,
					about varchar(64) DEFAULT NULL,
					content text DEFAULT NULL,
					PRIMARY KEY (id),
					UNIQUE KEY email (email)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
			`)
		if err != nil {
			log.Fatal(err)
		}
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
		var name, email, password, about, content string

		err := rows.Scan(&id, &name, &email, &password, &about, &content)
		if err != nil {
			log.Fatal(err)
		}

		// Print the results
		fmt.Printf("User ID: %d\nName: %s\nEmail: %s\nPassword: %s\nAbout: %s\nContent: %s\n\n", id, name, email, password, about, content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	articles := testTable(db, "articles")
	if articles != nil {

		// Create articles table
		_, err := db.Exec(`
					CREATE TABLE IF NOT EXISTS articles (
					id int(11) NOT NULL AUTO_INCREMENT,
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
			log.Fatal(err)
		}
		log.Print(articles)
	}
	skills := testTable(db, "skills")
	if skills != nil {

		// Create SKILLS table
		_, err := db.Exec(`
					CREATE TABLE IF NOT EXISTS skills (
					id int(11) NOT NULL AUTO_INCREMENT,
					name varchar(60) NOT NULL,
					logo text NOT NULL,
					textcolor varchar(255) NOT NULL,
					fillcolor varchar(255) NOT NULL,
					bgcolor varchar(255) NOT NULL,
					author int(11) NOT NULL,
					PRIMARY KEY (id),
					KEY author (author)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
			`)
		skills := "skills table created"
		if err != nil {
			log.Fatal(err)
		}
		log.Print(skills)
	}
	projects := testTable(db, "projects")
	if projects != nil {

		// Create PROJECTS table
		_, err := db.Exec(`
					CREATE TABLE IF NOT EXISTS projects (
					id int(11) NOT NULL AUTO_INCREMENT,
					title varchar(255) NOT NULL,
					description varchar(255) NOT NULL,
					url varchar(255) NOT NULL,
					image varchar(255) DEFAULT NULL,
					classes varchar(255) DEFAULT NULL,
					author int(11) NOT NULL,
					PRIMARY KEY (id),
					KEY author (author)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
			`)
		projects := "projects table created"
		if err != nil {
			log.Fatal(err)
		}

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
