package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

func Publish(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	if user != nil {
		data.Articles = user.FindArticles()
	}

	views.Hx(w, r, "dashboard", "publish", data)
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

	views.Hx(w, r, "main_layout", "edit_article", data)

}

// Blog post
func Blog(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	articlesPerPage := 10
	result, err := models.BlogTimeline(page, articlesPerPage)
	if err != nil {
		http.Error(w, "Error fetching blog timeline", http.StatusInternalServerError)
		return
	}

	var ctx Context
	ctx.Articles = result.Articles
	ctx.TotalArticles = result.TotalArticles
	ctx.ArticlesPerPage = result.ArticlesPerPage
	ctx.TotalPages = result.TotalPages
	ctx.CurrentPage = result.CurrentPage

	views.Hx(w, r, "main_layout", "blog", ctx)
}

func HomeFeedService(w http.ResponseWriter, r *http.Request) {
	data.Articles = models.LatestArticles(6) // 6 latest articles
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var tmpl string

	if isHTMXRequest {
		tmpl = "home-feed"
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	var response bytes.Buffer

	if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}

// Post is the post/article controller.
func Post(w http.ResponseWriter, r *http.Request) {
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
	views.Hx(w, r, "main_layout", "post", data)
}
