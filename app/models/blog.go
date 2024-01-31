package models

import (
	"log"
	"time"
)

// HomeFeed is a list of all articles.
func HomeFeed() []*Article {
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

// BlogTimeline fetches a list of blog articles for a specific page.
// ToDo: tags
func BlogTimeline(page int) []*Article {
	var articles []*Article

	articlesPerPage := 10
	offset := (page - 1) * articlesPerPage

	rows, err := Db.Query(`
        SELECT articles.id, articles.image, articles.slug, articles.title, articles.content, users.name, articles.created_at
        FROM articles
        JOIN users ON users.id = articles.author
        WHERE articles.is_draft = 0
        ORDER BY articles.created_at DESC
        LIMIT ? OFFSET ?
    `, articlesPerPage, offset)
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
