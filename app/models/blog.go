package models

import (
	"database/sql"
	"log"
	"math"
	"time"
)

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

// HomeFeed is a list of all articles.
func LatestArticles() []*Article {
	var articles []*Article

	rows, err := Db.Query(`
    SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at, tags.tag_id, tags.tag_name
    FROM articles
    JOIN users ON users.id = articles.author
    LEFT JOIN article_tags ON article_tags.article_id = articles.id
    LEFT JOIN tags ON tags.tag_id = article_tags.tag_id
    WHERE articles.is_draft = 0
    ORDER BY articles.created_at DESC 
    LIMIT 6
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	articleMap := make(map[int]*Article)

	for rows.Next() {
		var (
			id        int
			image     string
			slug      string
			title     string
			content   string
			author    string
			createdAt []byte
			tagID     sql.NullInt64
			tagName   sql.NullString
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &author, &createdAt, &tagID, &tagName)
		if err != nil {
			print("Error finding articles")
			log.Fatal(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Fatal(err)
		}

		// Check if the article already exists in the map
		article, ok := articleMap[id]
		if !ok {
			article = &Article{
				ID:        id,
				Image:     image,
				Slug:      slug,
				Title:     title,
				Content:   content,
				Author:    User{Name: author},
				CreatedAt: parsedCreatedAt,
				IsDraft:   0,
				Tags:      make([]*Tag, 0),
			}
			// Add the article to the map
			articleMap[id] = article
			articles = append(articles, article)
		}

		// Add tags to the article
		if tagID.Valid && tagName.Valid {
			tag := &Tag{
				ID:   int(tagID.Int64),
				Name: tagName.String,
			}
			article.Tags = append(article.Tags, tag)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return articles
}

// BlogTimeline fetches a list of blog articles for a specific page.
// ToDo: tags
func BlogTimeline(page int, articlesPerPage int) (*Timeline, error) {
	var result Timeline

	offset := (page - 1) * articlesPerPage

	// Count total articles
	var totalArticles int
	err := Db.QueryRow("SELECT COUNT(id) FROM articles WHERE is_draft = 0").Scan(&totalArticles)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalArticles) / float64(articlesPerPage)))

	rows, err := Db.Query(`
        SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at
        FROM articles
        JOIN users ON users.id = articles.author
        WHERE articles.is_draft = 0
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
		)
		err = rows.Scan(&id, &image, &slug, &title, &content, &author, &createdAt)
		if err != nil {
			return nil, err
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}
		user := User{
			Name: author,
		}
		articles = append(articles, &Article{id, image, slug, title, content, user, parsedCreatedAt, 0, nil})
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
