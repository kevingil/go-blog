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
}

// Projects is a model for home page projects.
type Project struct {
	ID       int
	Title    string
	Abstract string
	Url      string
	Addon    string
}

var (
	// Db is a database connection.
	Db *sql.DB

	// Err is an error returned.
	Err error
)

// Projects_Test returns test data for projects.
func Projects_Test() []*Project {
	var projects []*Project

	// Create dummy Project objects
	project0 := &Project{
		Title:    "Interior Designer AI",
		Abstract: "Home design image renders with DALL·E 2. Python backend, React frontend.",
		Url:      "http://147.182.233.135:3000/",
		Addon:    "col-span-2",
	}

	project1 := &Project{
		Title:    "Blog",
		Abstract: "Minimalist Go blog with mysql and htmx frontend",
		Url:      "/post/minimalist-blog-with-go-mysql-htmx-and-tailwind",
	}

	project2 := &Project{
		Title:    "CoffeeGPT",
		Abstract: "Use OpenAI to dial in your morning specialty coffee.",
		Url:      "/projects/coffeeapp",
	}

	project4 := &Project{
		Title:    "Client Side Moderation",
		Abstract: "Demo of TensorflowJS toxicity AI model for social media.",
		Url:      "/projects/moderatorjs",
	}

	project3 := &Project{
		Title:    "Document Viewer",
		Abstract: "Pure JS, drag, zoom, and resize for iframe content.",
		Url:      "/post/document-viewer-for-embedded-html-pages",
	}

	// Append the dummy projects to the list
	projects = append(projects, project0, project1, project2, project3, project4)

	return projects
}

// FindArticle is to print an article.
func FindArticle(slug string) *Article {
	rows, err := Db.Query(`SELECT articles.image, articles.title, articles.content, users.name, articles.created_at FROM articles JOIN users ON users.id = articles.author WHERE slug = ?`, slug)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var createdAt []byte
	user := &User{}
	article := &Article{}

	for rows.Next() {
		err = rows.Scan(&article.Image, &article.Title, &article.Content, &user.Name, &createdAt)
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

	rows, err := Db.Query(`SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at FROM articles JOIN users ON users.id = articles.author ORDER BY articles.created_at DESC`)
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
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}
		user := User{
			Name: author,
		}
		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt})
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
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &user
}

// FindArticle finds an user article by ID.
func (user User) FindArticle(id int) *Article {
	rows, err := Db.Query(`SELECT image, slug, title, content, created_at FROM articles WHERE id = ? AND author = ?`, id, user.ID)
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
		err = rows.Scan(&article.Image, &article.Slug, &article.Title, &article.Content, &createdAt)
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

	rows, err := Db.Query(`SELECT id, image, slug, title, content, created_at FROM articles WHERE author = ?`, user.ID)
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
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt})
	}

	// Sort the articles in descending order by created_at
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	return articles
}

// Create creates a user.
func (user User) Create() *User {
	result, err := Db.Exec("INSERT INTO users(name, email, password) VALUES(?, ?, ?)", user.Name, user.Email, user.Password)
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
		"INSERT INTO articles(image, slug, title, content, author, created_at) VALUES(?, ?, ?, ?, ?, ?)",
		article.Image,
		article.Slug,
		article.Title,
		article.Content,
		article.Author.ID,
		article.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateArticle updates an article.
func (user User) UpdateArticle(article *Article) {
	_, err := Db.Exec(
		"UPDATE articles SET image = ?, slug = ?, title = ?, content = ? WHERE id = ? AND author = ?",
		article.Image,
		article.Slug,
		article.Title,
		article.Content,
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
