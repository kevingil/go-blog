package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kevingil/blog/internal/database"
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

// FindArticle is to print an article.
func FindArticle(slug string) *Article {
	// Adjusted SQL query to include tags
	rows, err := database.Db.Query(`
        SELECT articles.id, articles.image, articles.title, articles.content, users.name, articles.created_at, 
        GROUP_CONCAT(tags.tag_name) AS tags
        FROM articles
        JOIN users ON users.id = articles.author
        LEFT JOIN article_tags ON article_tags.article_id = articles.id
        LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
        WHERE slug = ?
        GROUP BY articles.id
    `, slug)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var createdAt []byte
	var tagsStr string
	user := &User{}
	article := &Article{}

	if rows.Next() {
		err = rows.Scan(&article.ID, &article.Image, &article.Title, &article.Content, &user.Name, &createdAt, &tagsStr)
		if err != nil {
			log.Fatal(err)
		}
		article.CreatedAt = parseTime(createdAt)
		article.Author = *user
		article.Tags = parseTags(tagsStr)
	}

	return article
}

// FindArticle finds an user article by ID.
func (user User) FindArticle(id int) (*Article, error) {
	var createdAt []byte
	var tagsStr sql.NullString
	var authorID int // Use an integer to capture the author ID only

	article := &Article{
		ID:     id,
		Author: user,
	}

	rows, err := database.Db.Query(`
	SELECT articles.image, articles.slug, articles.title, articles.content, articles.author, articles.created_at, articles.is_draft,
	GROUP_CONCAT(tags.tag_name) AS tags
	FROM articles
	JOIN users ON users.id = articles.author
	LEFT JOIN article_tags ON article_tags.article_id = articles.id
	LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
	WHERE articles.id = ? AND articles.author = ?
	GROUP BY articles.id
    `, id, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&article.Image, &article.Slug, &article.Title, &article.Content, &authorID, &createdAt, &article.IsDraft, &tagsStr)
		if err != nil {
			log.Fatal(err)
		}
		article.CreatedAt = parseTime(createdAt)

		if tagsStr.Valid {
			article.Tags = parseTags(tagsStr.String)
		} else {
			article.Tags = []*Tag{}
		}
	}
	return article, nil
}

// FindArticles finds user articles
func (user User) FindArticles() []*Article {
	// Adjusted SQL query to include tags
	rows, err := database.Db.Query(`
        SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, articles.created_at, articles.is_draft, 
        GROUP_CONCAT(tags.tag_name) AS tags
        FROM articles
        LEFT JOIN article_tags ON article_tags.article_id = articles.id
        LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
        WHERE author = ?
        GROUP BY articles.id
        ORDER BY articles.created_at DESC
    `, user.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var articles []*Article

	for rows.Next() {
		var (
			id        int
			image     string
			slug      string
			title     string
			content   string
			createdAt []byte
			isDraft   int
			tagsStr   sql.NullString
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &createdAt, &isDraft, &tagsStr)
		if err != nil {
			log.Fatal(err)
		}
		tags := parseTags("")
		if tagsStr.Valid {
			tags = parseTags(tagsStr.String)
		}
		article := &Article{
			ID:        id,
			Image:     image,
			Slug:      slug,
			Title:     title,
			Content:   content,
			Author:    user,
			CreatedAt: parseTime(createdAt),
			IsDraft:   isDraft,
			Tags:      tags,
		}
		articles = append(articles, article)
	}

	return articles
}

func (user User) CountArticles() int {
	var count int

	err := database.Db.QueryRow(`SELECT COUNT(*) FROM articles WHERE author = ? AND is_draft = 0`, user.ID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

func (user User) CountDrafts() int {
	var count int

	err := database.Db.QueryRow(`SELECT COUNT(*) FROM articles WHERE author = ? AND is_draft = 1`, user.ID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

// CreateArticle creates an article.
func (user User) CreateArticle(article *Article) {
	_, err := database.Db.Exec(
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

func (user User) UpdateArticle(article *Article) error {

	tx, err := database.Db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Update the article details
	_, err = tx.Exec(
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
		tx.Rollback() // Roll back the transaction on error
		log.Println("Error updating article:", err)
		return err
	}

	// Delete existing tag associations
	_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", article.ID)
	if err != nil {
		tx.Rollback()
		log.Println("Error deleting existing tags:", err)
		return err
	}

	// Update tags
	for _, tag := range article.Tags {
		var tagID int64
		err := tx.QueryRow("SELECT tag_id FROM tags WHERE tag_name = ?", tag.Name).Scan(&tagID)
		if err == sql.ErrNoRows {
			// Tag does not exist, create it
			result, err := tx.Exec("INSERT INTO tags (tag_name) VALUES (?)", tag.Name)
			if err != nil {
				tx.Rollback()
				log.Println("Error creating tag:", err)
				return err
			}
			tagID, _ = result.LastInsertId()
		} else if err != nil {
			tx.Rollback()
			log.Println("Error checking tag existence:", err)
			return err
		}

		// Create new article-tag relationship
		_, err = tx.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", article.ID, tagID)
		if err != nil {
			tx.Rollback()
			log.Println("Error creating article-tag relationship:", err)
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		return err
	}

	return nil
}

// DeleteArticle deletes an article.
func (user User) DeleteArticle(article *Article) {
	_, err := database.Db.Exec(
		"DELETE FROM articles WHERE id = ? AND author = ?",
		article.ID,
		user.ID,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// Will return time from bytes or current time
func parseTime(datetimeBytes []byte) time.Time {
	const layout = "2006-01-02 15:04:05" // This is the Go time layout format
	datetimeStr := string(datetimeBytes)
	parsedTime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		return time.Now()
	}
	return parsedTime
}

func parseTags(tagsStr string) []*Tag {
	var tags []*Tag
	if tagsStr == "" {
		return tags
	}
	tagNames := strings.Split(tagsStr, ",")
	for _, tagName := range tagNames {
		tags = append(tags, &Tag{Name: tagName})
	}
	return tags
}
