package models

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
)

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
	Tags      []*Tag
}

type Tag struct {
	ID   int
	Name string
}

var (
	// Db is a database connection.
	Db *sql.DB

	// Err is an error returned.
	Err error
)

// FindArticle is to print an article.
func FindArticle(slug string) *Article {
	rows, err := Db.Query(`SELECT articles.id, articles.image, articles.title, articles.content, users.name, articles.created_at FROM articles JOIN users ON users.id = articles.author WHERE slug = ?`, slug)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var createdAt []byte
	user := &User{}
	article := &Article{}

	for rows.Next() {
		err = rows.Scan(&article.ID, &article.Image, &article.Title, &article.Content, &user.Name, &createdAt)
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
		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt, 0, nil})
	}

	return articles
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

func (article Article) FindTags() []*Tag {
	var tags []*Tag

	rows, err := Db.Query(`
	SELECT tags.tag_id, tags.tag_name
	FROM tags
	JOIN article_tags ON article_tags.tag_id = tags.tag_id
	WHERE article_tags.article_id = ?
`, article.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)
		err = rows.Scan(&id, &name)
		if err != nil {
			print("Error finding tags")
			log.Fatal(err)
		}
		tags = append(tags, &Tag{id, name})
	}
	return tags
}

/*

CREATE TABLE `articles` (
	`id` int NOT NULL AUTO_INCREMENT,
	`image` varchar(255),
	`slug` varchar(255) NOT NULL,
	`title` varchar(60) NOT NULL,
	`content` text NOT NULL,
	`author` int NOT NULL,
	`created_at` datetime NOT NULL,
	`is_draft` tinyint(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	UNIQUE KEY `slug` (`slug`),
	KEY `author` (`author`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE `tags` (
	`tag_id` int NOT NULL AUTO_INCREMENT,
	`tag_name` varchar(50),
	PRIMARY KEY (`tag_id`),
	UNIQUE KEY `tag_name` (`tag_name`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE `article_tags` (
	`article_id` int NOT NULL,
	`tag_id` int NOT NULL,
	PRIMARY KEY (`article_id`, `tag_id`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;
*/

// Check for tags linked to article, add new tags if not exists, add relationship between article and tag.
// if removed tags no longer have matching articles, delete unused tags

// UpdateTags updates the tags associated with an article.
func (article Article) UpdateTags(tags []*Tag) {
	// Ensure there is a valid database connection
	if Db == nil {
		fmt.Println("Database connection is not initialized.")
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}

	// Clear existing tags for the article
	_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", article.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting existing tags:", err)
		return
	}

	// Iterate through the provided tags and update the relationships
	for _, tag := range tags {
		// Check if the tag exists
		var tagID int64
		err := tx.QueryRow("SELECT tag_id FROM tags WHERE tag_name = ?", tag.Name).Scan(&tagID)
		if err == sql.ErrNoRows {
			// If the tag doesn't exist, create it
			result, err := tx.Exec("INSERT INTO tags (tag_name) VALUES (?)", tag.Name)
			if err != nil {
				tx.Rollback()
				fmt.Println("Error creating tag:", err)
				return
			}
			tagID, _ = result.LastInsertId()
		} else if err != nil {
			// If there is an error querying the database, rollback the transaction
			tx.Rollback()
			fmt.Println("Error checking tag existence:", err)
			return
		}

		// Create the relationship between the article and the tag
		_, err = tx.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", article.ID, tagID)
		if err != nil {
			// If there is an error creating the relationship, rollback the transaction
			tx.Rollback()
			fmt.Println("Error creating article-tag relationship:", err)
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
	}
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

		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt, isDraft, nil})
	}

	// Sort the articles in descending order by created_at
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	return articles
}

func (user User) CountArticles() int {
	var count int

	err := Db.QueryRow(`SELECT COUNT(*) FROM articles WHERE author = ? AND is_draft = 0`, user.ID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

func (user User) CountDrafts() int {
	var count int

	err := Db.QueryRow(`SELECT COUNT(*) FROM articles WHERE author = ? AND is_draft = 1`, user.ID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
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
		"UPDATE articles SET image = ?, slug = ?, title = ?, content = ?, created_at = ?, is_draft = ? WHERE id = ? AND author = ?",
		article.Image,
		article.Slug,
		article.Title,
		article.Content,
		article.CreatedAt,
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
