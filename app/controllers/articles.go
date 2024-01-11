package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

// Post is the post/article controller.
func Article(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := models.FindArticle(vars["slug"])
	data.Article = article
	data.Tags = article.FindTags()

	if article == nil {
		data.Article = &models.Article{
			Image:   "",
			Title:   "Post Not Found",
			Content: "This post doesn't exists.",
		}
	}
	if data.Tags == nil {
		data.Tags = []*models.Tag{
			{
				Name: "",
			},
		}
	}
	Hx(w, r, "main_layout", "article", data)
}


