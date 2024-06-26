package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/kevingil/blog/internal/models"
)

// Dashboard is a controller for users to list articles.
func Dashboard(c *fiber.Ctx) error {
	cookie := c.Cookies("cookie_name")
	user := Sessions[cookie]
	permission(c)

	model := c.Query("edit")
	id, _ := strconv.Atoi(c.Query("id"))
	delete := c.Query("delete")

	switch c.Method() {
	case fiber.MethodGet:
		switch model {
		case "article":
			if delete != "" && id != 0 {
				article := &models.Article{
					ID: id,
				}

				user.DeleteArticle(article)
				return c.Redirect("/dashboard", fiber.StatusSeeOther)
			}
		default:

			data := map[string]interface{}{
				"ArticleCount": user.CountArticles(),
				"DraftCount":   user.CountDrafts(),
				"Articles":     user.FindArticles(),
				"User":         user,
			}

			return c.Render("dashboardPage", data)
		}
	case fiber.MethodPost:
		switch model {
		case "article":
			if user != nil {
				isDraftStr := c.FormValue("isDraft")
				isDraft, err := strconv.Atoi(isDraftStr)
				if err != nil {
					isDraft = 0
				}

				if id == 0 {
					// Create a new article
					article := &models.Article{
						Image:     c.FormValue("image"),
						Slug:      slug.Make(c.FormValue("title")),
						Title:     c.FormValue("title"),
						Content:   c.FormValue("content"),
						Author:    *user,
						CreatedAt: time.Now(),
						IsDraft:   isDraft,
						Tags:      []*models.Tag{},
					}

					user.CreateArticle(article)
				} else {
					// Update existing article
					createdAtStr := c.FormValue("createdat")
					createdAt, err := time.Parse("2006-01-02", createdAtStr)
					if err != nil {
						createdAt = time.Now()
					}
					article := &models.Article{
						ID:        id,
						Image:     c.FormValue("image"),
						Slug:      slug.Make(c.FormValue("title")),
						Title:     c.FormValue("title"),
						Content:   c.FormValue("content"),
						CreatedAt: createdAt,
						IsDraft:   isDraft,
						Tags:      []*models.Tag{},
					}

					// Convert form input to Tags and append
					rawtags := c.FormValue("tags")
					tagNames := strings.Split(rawtags, ",")
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

				return c.Redirect("/dashboard/articles", fiber.StatusSeeOther)
			}
		}
	}

	return nil
}
