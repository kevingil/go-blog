package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Refactor the DashboardArticles function
func DashboardArticles(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	user := Sessions[sessionID]
	data := map[string]interface{}{
		"User":     user,
		"Articles": user.FindArticles(),
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("dashboardArticlesPage", data, "")
	} else {
		return c.Render("dashboardArticlesPage", data)
	}
}

// Refactor the EditArticle function
func EditArticle(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	user := Sessions[sessionID]
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

// View blog post
func BlogPostPage(c *fiber.Ctx) error {
	slug := c.Params("slug")
	article := models.FindArticle(slug)
	data := map[string]interface{}{
		"Article": article,
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("blogPostPage", data, "")
	} else {
		return c.Render("blogPostPage", data)
	}
}
