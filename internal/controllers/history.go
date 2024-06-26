package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Returns a partial html element with recent articles
func RecentPostsPartial(c *fiber.Ctx) error {
	isHTMXRequest := c.Get("HX-Request") == "true"
	if isHTMXRequest {
		data := map[string]interface{}{
			"Articles": models.LatestArticles(6), //Pass article count
		}

		return c.Render("homeFeed", data)
	} else {
		//Redirect home if trying to call the endpoint directly
		return c.Redirect("/", fiber.StatusSeeOther)
	}
}

func BlogPage(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	articlesPerPage := 10
	result, err := models.BlogTimeline(page, articlesPerPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching blog timeline")
	}

	data := map[string]interface{}{
		"Articles":        result.Articles,
		"TotalArticles":   result.TotalArticles,
		"ArticlesPerPage": articlesPerPage,
		"TotalPages":      result.TotalPages,
		"CurrentPage":     result.CurrentPage,
	}

	return c.Render("blogPostsPage", data)
}
