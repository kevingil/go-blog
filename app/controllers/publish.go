package controllers

import (
	"net/http"
	"strconv"

	"github.com/kevingil/blog/app/models"
)

func Publish(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	if user != nil {
		data.Articles = user.FindArticles()
	}

	Hx(w, r, "dashboard", "publish", data)
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	data.Article = &models.Article{
		Image:   "",
		Title:   "",
		Content: "",
		IsDraft: 0,
	}

	if user != nil && id != 0 {
		data.Article = user.FindArticle(id)
		data.Tags = data.Article.FindTags()
	}

	Hx(w, r, "main_layout", "edit_article", data)

}
