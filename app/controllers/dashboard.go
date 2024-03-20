package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/kevingil/blog/app/cmd"
	"github.com/kevingil/blog/app/models"
)

// Dashboard is a controller for users to list articles.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")
	layout := "dashboard"

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
			}
		default:
			if user != nil {
				data.ArticleCount = user.CountArticles()
				data.DraftCount = user.CountDrafts()
				data.Articles = user.FindArticles()
			}

			cmd.Hx(w, r, layout, "dashboard_home", data)

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
					article := &models.Article{
						Image:     r.FormValue("image"),
						Slug:      slug.Make(r.FormValue("title")),
						Title:     r.FormValue("title"),
						Content:   r.FormValue("content"),
						Author:    *user,
						CreatedAt: time.Now(),
						IsDraft:   isDraft,
					}

					user.CreateArticle(article)
				} else {
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
					}

					// Handle tags
					rawtags := r.Form["tags"]
					tags := make([]*models.Tag, 0)
					tagNames := strings.Split(rawtags[0], ",")
					for _, tagName := range tagNames {
						trimmedTagName := strings.TrimSpace(tagName)
						if trimmedTagName != "" {
							tag := &models.Tag{
								Name: trimmedTagName,
							}
							tags = append(tags, tag)
						}
					}

					user.UpdateArticle(article)
					article.UpdateTags(tags)
				}
			}

			http.Redirect(w, r, "/dashboard/publish", http.StatusSeeOther)
		}
	}
}
