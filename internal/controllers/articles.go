package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/kevingil/blog/internal/models"
)

// Articles page, shows edit actions
func EditArticlesPage(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	data := map[string]interface{}{
		"User":     user,
		"Articles": user.FindArticles(),
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("adminArticlesPage", data, "")
	} else {
		return c.Render("adminArticlesPage", data)
	}
}

// Edit article form page
func EditArticlePage(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	data := map[string]interface{}{}
	id, _ := strconv.Atoi(c.Query("id"))

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
			data["Article"] = article
		} else {
			data["Article"] = defaultArticle
		}
	} else {
		data["Article"] = defaultArticle
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("edit-article", data, "")
	} else {
		return c.Render("edit-article", data)
	}
}

func DeleteArticle(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	id, _ := strconv.Atoi(c.Query("id"))
	delete := c.Query("delete")

	if delete != "" && id != 0 {
		article := &models.Article{
			ID: id,
		}

		user.DeleteArticle(article)
		return c.Redirect("/admin/articles", fiber.StatusSeeOther)
	}

	return c.Redirect("/login", fiber.StatusSeeOther)
}

func EditArticle(c *fiber.Ctx) error {

	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	id, _ := strconv.Atoi(c.Query("id"))

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

		return c.Redirect("/admin/articles", fiber.StatusSeeOther)
	}
	return c.Redirect("/login", fiber.StatusSeeOther)

}
