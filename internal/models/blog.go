package models

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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

type Timeline struct {
	Articles        []*Article
	TotalArticles   int
	ArticlesPerPage int
	TotalPages      int
	CurrentPage     int
}

type Homepage struct {
	LatestArticles []*Article
	TopArticles    []*Article
}

// LatestArticles returns the latest articles with a limit.
func LatestArticles(limit int) []*Article {
	var articles []*Article

	rows, err := database.Db.Query(`
		SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at, 
			group_concat(tags.tag_name) as tags
		FROM articles
		JOIN users ON users.id = articles.author
		LEFT JOIN article_tags ON article_tags.article_id = articles.id
		LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
		WHERE articles.is_draft = 0
		GROUP BY articles.id
		ORDER BY articles.created_at DESC
		LIMIT ?
	`, limit)
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
			tagsStr   string
		)
		err := rows.Scan(&id, &image, &slug, &title, &content, &author, &createdAt, &tagsStr)
		if err != nil {
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}

		tags := []*Tag{}
		if tagsStr != "" {
			for _, tagName := range strings.Split(tagsStr, ",") {
				tags = append(tags, &Tag{Name: tagName})
			}
		}

		article := &Article{
			ID:        id,
			Image:     image,
			Slug:      slug,
			Title:     title,
			Content:   content,
			Author:    User{Name: author},
			CreatedAt: parsedCreatedAt,
			IsDraft:   0,
			Tags:      tags,
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return articles
}

func BlogTimeline(page int, articlesPerPage int) (*Timeline, error) {

	var result Timeline

	offset := (page - 1) * articlesPerPage

	// Count total articles
	var totalArticles int
	err := database.Db.QueryRow("SELECT COUNT(id) FROM articles WHERE is_draft = 0").Scan(&totalArticles)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalArticles) / float64(articlesPerPage)))

	// Query for articles and tags
	rows, err := database.Db.Query(`
    SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at, group_concat(tags.tag_name) as tags
    FROM articles
    JOIN users ON users.id = articles.author
    LEFT JOIN article_tags ON article_tags.article_id = articles.id
    LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
    WHERE articles.is_draft = 0
    GROUP BY articles.id
    ORDER BY articles.created_at DESC
    LIMIT ? OFFSET ?
    `, articlesPerPage, offset)
	if err != nil {
		return nil, err
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
			author    string
			createdAt []byte
			tagsStr   string
		)
		err := rows.Scan(&id, &image, &slug, &title, &content, &author, &createdAt, &tagsStr)
		if err != nil {
			return nil, err
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}

		tags := []*Tag{}
		if tagsStr != "" {
			for _, tagName := range strings.Split(tagsStr, ",") {
				tags = append(tags, &Tag{Name: tagName})
			}
		}

		article := &Article{
			ID:        id,
			Image:     image,
			Slug:      slug,
			Title:     title,
			Content:   content,
			Author:    User{Name: author},
			CreatedAt: parsedCreatedAt,
			IsDraft:   0,
			Tags:      tags,
		}
		articles = append(articles, article)
	}

	result = Timeline{
		Articles:        articles,
		TotalArticles:   totalArticles,
		ArticlesPerPage: articlesPerPage,
		TotalPages:      totalPages,
		CurrentPage:     page,
	}

	return &result, nil
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
func (user User) FindArticle(id int) *Article {
	rows, err := database.Db.Query(`SELECT image, slug, title, content, created_at, is_draft FROM articles WHERE id = ? AND author = ?`, id, user.ID)
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

func (article Article) UpdateTags(tags []*Tag) {
	if database.Db == nil {
		fmt.Println("Database connection is not initialized.")
		return
	}
	tx, err := database.Db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}
	_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", article.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting existing tags:", err)
		return
	}

	for _, tag := range tags {
		var tagID int64
		err := tx.QueryRow("SELECT tag_id FROM tags WHERE tag_name = ?", tag.Name).Scan(&tagID)
		if err == sql.ErrNoRows {
			result, err := tx.Exec("INSERT INTO tags (tag_name) VALUES (?)", tag.Name)
			if err != nil {
				tx.Rollback()
				fmt.Println("Error creating tag:", err)
				return
			}
			tagID, _ = result.LastInsertId()
		} else if err != nil {
			tx.Rollback()
			fmt.Println("Error checking tag existence:", err)
			return
		}

		_, err = tx.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", article.ID, tagID)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error creating article-tag relationship:", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
	}
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

// UpdateArticle updates an article.
func (user User) UpdateArticle(article *Article) {
	_, err := database.Db.Exec(
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