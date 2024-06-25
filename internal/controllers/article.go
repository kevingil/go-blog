package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	Handle(w, r, data)
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

// Data structure for the Post page
type PostData struct {
	Article *models.Article
}

// Refactor the Post function
func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := models.FindArticle(vars["slug"])
	data := map[string]interface{}{
		"Article":  article,
		"Template": "slug",
		"Url":      "/blog",
	}

	Handle(w, r, data)
}
