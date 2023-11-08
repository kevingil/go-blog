package models

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDbInit(t *testing.T) {
	// Test Initializing the database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/blog")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Db = db
}

func TestFindArticle(t *testing.T) {
	TestDbInit(t)
	// Test get article by slug
	slug := "test-article"
	article := FindArticle(slug)
	if article == nil {
		t.Fatal("Expected Test Article, got none")
	}

}

func TestArticles(t *testing.T) {
	TestDbInit(t)
	// Test find all articles
	articles := Articles()
	if len(articles) == 0 {
		t.Fatal("Expected articles, got none")
	}
}

func TestUserMethods(t *testing.T) {
	TestDbInit(t)
	// TODO
	// Test other methods
}
