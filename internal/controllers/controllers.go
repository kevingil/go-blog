package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

// Index serves the homepage.
func Index(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"About":    models.About(),
		"Projects": models.GetProjects(),
	}

	return c.Render("indexPage", data)
}

// About serves the about page.
func About(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"About":  models.AboutPage(),
		"Skills": models.Skills_Test(),
	}

	return c.Render("aboutPage", data)
}

// Contact serves the contact page.
func Contact(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"Contact": models.ContactPage(),
	}

	return c.Render("contactPage", data)
}
