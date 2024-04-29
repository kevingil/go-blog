package controllers

import (
	"net/http"
	"strconv"

	"github.com/kevingil/blog/internal/models"
)

// Data structure for the HomeFeedService
type HomeFeedData struct {
	Articles []*models.Article
}

// Refactor the HomeFeedService function
func HomeFeedService(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	if isHTMXRequest {
		req := Request{
			W:      w,
			R:      r,
			Layout: "",
			Tmpl:   "home-feed",
			Data: HomeFeedData{
				Articles: models.LatestArticles(6),
			},
		}
		render(req)
	} else {
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
