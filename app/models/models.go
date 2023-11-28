package models

import (
	"database/sql"
	"log"
	"sort"
	"time"
)

// User is a model for users.
type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	About    string
	Content  string
}

// Article is a model for articles.
type Article struct {
	ID        int
	Image     string
	Slug      string
	Title     string
	Content   string
	Author    User
	CreatedAt time.Time
	IsDraft   int
}

var (
	// Db is a database connection.
	Db *sql.DB

	// Err is an error returned.
	Err error
)

// FindArticle is to print an article.
func FindArticle(slug string) *Article {
	rows, err := Db.Query(`SELECT articles.image, articles.title, articles.content, users.name, articles.created_at articles.is_draft FROM articles JOIN users ON users.id = articles.author WHERE slug = ?`, slug)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var createdAt []byte
	user := &User{}
	article := &Article{}

	for rows.Next() {
		err = rows.Scan(&article.Image, &article.Title, &article.Content, &user.Name, &createdAt, &article.IsDraft)
		if err != nil {
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}
		article.CreatedAt = parsedCreatedAt
		article.Author = *user
	}

	return article
}

// Articles is a list of all articles.
func Articles() []*Article {
	var articles []*Article

	rows, err := Db.Query(`
    SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at
    FROM articles
    JOIN users ON users.id = articles.author
    WHERE articles.is_draft = 0
    ORDER BY articles.created_at DESC 
`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			image     string
			slug      string
			title     string
			content   string
			author    string
			createdAt []byte
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &author, &createdAt)
		if err != nil {
			print("Error finding articles")
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}
		user := User{
			Name: author,
		}
		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt, 0})
	}

	return articles
}

// Find finds a user by email.
func (user User) Find() *User {
	rows, err := Db.Query(`SELECT * FROM users WHERE email = ?`, user.Email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var about, content sql.NullString
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &about, &content)
		if err != nil {
			log.Fatal(err)
		}

		// Check for NULL values
		if about.Valid {
			user.About = about.String
		} else {
			user.About = ""
		}

		if content.Valid {
			user.Content = content.String
		} else {
			user.Content = ""
		}
	}

	return &user
}

// GetProfile finds a user by email and returns a user profile.
func (user User) GetProfile() *User {
	rows, err := Db.Query(`SELECT * FROM users WHERE email = ?`, user.Email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	profile := &User{}

	for rows.Next() {
		var about, content sql.NullString
		err = rows.Scan(&profile.ID, &profile.Name, &profile.Email, &profile.Password, &about, &content)
		if err != nil {
			log.Fatal(err)
		}

		// Check for NULL values
		if about.Valid {
			profile.About = about.String
		} else {
			profile.About = ""
		}

		if content.Valid {
			profile.Content = content.String
		} else {
			profile.Content = ""
		}
	}

	return profile
}

// FindArticle finds an user article by ID.
func (user User) FindArticle(id int) *Article {
	rows, err := Db.Query(`SELECT image, slug, title, content, created_at, is_draft FROM articles WHERE id = ? AND author = ?`, id, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var createdAt []byte
	article := &Article{
		ID:     id,
		Author: user,
	}

	for rows.Next() {
		err = rows.Scan(&article.Image, &article.Slug, &article.Title, &article.Content, &createdAt, &article.IsDraft)
		if err != nil {
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}
		article.CreatedAt = parsedCreatedAt
	}

	return article
}

// FindArticles finds user articles.
// https://go.dev/doc/database/querying

func (user User) FindArticles() []*Article {
	var articles []*Article

	rows, err := Db.Query(`SELECT id, image, slug, title, content, created_at, is_draft FROM articles WHERE author = ?`, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			image     string
			slug      string
			title     string
			content   string
			createdAt []byte
			isDraft   int
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &createdAt, &isDraft)
		if err != nil {
			print("Error finding article")
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt, isDraft})
	}

	// Sort the articles in descending order by created_at
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	return articles
}

// Create creates a user.
func (user User) Create() *User {
	result, err := Db.Exec("INSERT INTO users(name, email, password, about, content) VALUES(?, ?, ?, NULL, NULL)",
		user.Name, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	if id != 0 {
		user.ID = int(id)
	}

	return &user
}

// CreateArticle creates an article.
func (user User) CreateArticle(article *Article) {
	_, err := Db.Exec(
		"INSERT INTO articles(image, slug, title, content, author, created_at, is_draft) VALUES(?, ?, ?, ?, ?, ?, ?)",
		article.Image,
		article.Slug,
		article.Title,
		article.Content,
		article.Author.ID,
		article.CreatedAt,
		article.IsDraft,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateArticle updates an article.
func (user User) UpdateArticle(article *Article) {
	_, err := Db.Exec(
		"UPDATE articles SET image = ?, slug = ?, title = ?, content = ?, is_draft = ? WHERE id = ? AND author = ?",
		article.Image,
		article.Slug,
		article.Title,
		article.Content,
		article.IsDraft,
		article.ID,
		user.ID,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// DeleteArticle deletes an article.
func (user User) DeleteArticle(article *Article) {
	_, err := Db.Exec(
		"DELETE FROM articles WHERE id = ? AND author = ?",
		article.ID,
		user.ID,
	)
	if err != nil {
		log.Fatal(err)
	}
}
