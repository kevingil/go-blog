package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/models"
)

func Publish(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	if user != nil {
		data.Articles = user.FindArticles()
	}

	render(w, r, "dashboard-layout", "publish", data)
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
	}

	render(w, r, "layout", "edit-article", data)

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

	render(w, r, "layout", "blog", ctx)
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

	if err := Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
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

	if article == nil {
		data.Article = &models.Article{
			Image:   "",
			Title:   "Post Not Found",
			Content: "This post doesn't exists.",
		}
	}
	render(w, r, "layout", "post", data)
}
