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

		data := map[string]interface{}{
			"Articles": models.LatestArticles(6), //Pass article count
		}

		Partial(w, r, data, "homeFeed")

	} else {
		//Redirect home if trying to call the endpoint directly
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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

	data := map[string]interface{}{
		"Articles":        result.Articles,
		"TotalArticles":   result.TotalArticles,
		"ArticlesPerPage": articlesPerPage,
		"TotalPages":      result.TotalPages,
		"CurrentPage":     result.CurrentPage,
	}
	Handle(w, r, data)

}
