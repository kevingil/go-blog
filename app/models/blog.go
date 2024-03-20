package models

import (
	"log"
	"math"
	"strings"
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

// LatestArticles returns the latest articles with a limit.
func LatestArticles(limit int) []*Article {
	var articles []*Article

	rows, err := Db.Query(`
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
	err := Db.QueryRow("SELECT COUNT(id) FROM articles WHERE is_draft = 0").Scan(&totalArticles)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalArticles) / float64(articlesPerPage)))

	// Query for articles and tags
	rows, err := Db.Query(`
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
