package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/kevingil/blog/internal/models"
)

// Dashboard is a controller for users to list articles.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	req := Request{
		W:      w,
		R:      r,
		Layout: "dashboard-layout",
		User:   user,
	}

	permission(w, r)

	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "article":
			if delete != "" && id != 0 {
				article := &models.Article{
					ID: id,
				}

				user.DeleteArticle(article)
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
		default:
			// Dashboard data for the default view
			dashboardData := struct {
				ArticleCount int
				DraftCount   int
				Articles     []*models.Article
				User         *models.User
			}{
				ArticleCount: user.CountArticles(),
				DraftCount:   user.CountDrafts(),
				Articles:     user.FindArticles(),
				User:         user,
			}

			req.Tmpl = "dashboard-home"
			req.Data = dashboardData
			render(req)
		}
	case http.MethodPost:
		switch model {
		case "article":
			if user != nil {
				isDraftStr := r.FormValue("isDraft")
				isDraft, err := strconv.Atoi(isDraftStr)
				if err != nil {
					isDraft = 0
				}

				if id == 0 {
					// Create a new article
					article := &models.Article{
						Image:     r.FormValue("image"),
						Slug:      slug.Make(r.FormValue("title")),
						Title:     r.FormValue("title"),
						Content:   r.FormValue("content"),
						Author:    *user,
						CreatedAt: time.Now(),
						IsDraft:   isDraft,
						Tags:      []*models.Tag{},
					}

					user.CreateArticle(article)
				} else {
					// Update existing article
					createdAtStr := r.FormValue("createdat")
					createdAt, err := time.Parse("2006-01-02", createdAtStr)
					if err != nil {
						createdAt = time.Now()
					}
					article := &models.Article{
						ID:        id,
						Image:     r.FormValue("image"),
						Slug:      slug.Make(r.FormValue("title")),
						Title:     r.FormValue("title"),
						Content:   r.FormValue("content"),
						CreatedAt: createdAt,
						IsDraft:   isDraft,
						Tags:      []*models.Tag{},
					}

					// Convert form input to Tags and append
					rawtags := r.Form["tags"]
					tagNames := strings.Split(rawtags[0], ",")
					for _, tagName := range tagNames {
						trimmedTagName := strings.TrimSpace(tagName)
						if trimmedTagName != "" {
							tag := &models.Tag{
								Name: trimmedTagName,
							}
							article.Tags = append(article.Tags, tag)
						}
					}

					user.UpdateArticle(article)
				}

				http.Redirect(w, r, "/dashboard/publish", http.StatusSeeOther)
				return
			}
		}
	}
}
