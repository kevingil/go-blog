package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kevingil/blog/internal/models"
)

// Refactor the Publish function
func DashboardArticles(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	permission(w, r)
	user := Sessions[cookie.Value]
	data := map[string]interface{}{
		"User":     user,
		"Articles": user.FindArticles(),
	}
	renderPage(w, r, data)
}

// Data structure for the EditArticle page
type EditArticleData struct {
	Article *models.Article
}

// Refactor the EditArticle function
func EditArticle(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "edit-article",
	}

	permission(w, r)

	defaultArticle := &models.Article{
		Image:   "",
		Title:   "",
		Content: "",
		IsDraft: 0,
		Tags:    []*models.Tag{},
	}

	if user != nil && id != 0 {
		article, err := user.FindArticle(id)
		if err == nil {
			req.Data = EditArticleData{
				Article: article,
			}
		} else {
			log.Print(err)
			req.Data = EditArticleData{
				Article: defaultArticle,
			}
		}
	} else {
		req.Data = EditArticleData{
			Article: defaultArticle,
		}
	}

	render(req)
}

// Refactor the Post function
func Post(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	article := models.FindArticle(slug)
	data := map[string]interface{}{
		"Article": article,
	}
	renderTemplate(w, r, data, "blogSlug")

}
