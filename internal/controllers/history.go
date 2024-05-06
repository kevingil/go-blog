package controllers

import (
	"net/http"
	"strconv"

	"github.com/kevingil/blog/internal/models"
)

// Returns a partial html element with recent articles
func RecentPostsPartial(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	if isHTMXRequest {
		articles := models.LatestArticles(6) //Pass article count
		data := struct {
			Articles []*models.Article
		}{
			Articles: articles,
		}
		req := Request{
			W:      w,
			R:      r,
			Layout: "",
			Tmpl:   "home-feed",
			Data:   data,
		}
		render(req)
	} else {
		//Redirect home if trying to call the endpoint directly
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Data structure for the Blog page
type BlogData struct {
	Articles        []*models.Article
	TotalArticles   int
	ArticlesPerPage int
	TotalPages      int
	CurrentPage     int
}

// Refactor the Blog function
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

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "blog",
		Data: BlogData{
			Articles:        result.Articles,
			TotalArticles:   result.TotalArticles,
			ArticlesPerPage: articlesPerPage,
			TotalPages:      result.TotalPages,
			CurrentPage:     result.CurrentPage,
		},
	}

	render(req)
}
