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

// Blog post
func Blog(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	Hx(w, r, "main_layout", "blog", data)
}

func TimelineService(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	var ctx Context
	ctx.Articles = models.BlogTimeline(page)
	Hx(w, r, "main_layout", "blog-feed", ctx)

}

func HomeFeedService(w http.ResponseWriter, r *http.Request) {
	data.Articles = models.HomeFeed()
	for _, post := range data.Articles {
		post.Tags = post.FindTags()
	}
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
	Hx(w, r, "main_layout", "post", data)
}
