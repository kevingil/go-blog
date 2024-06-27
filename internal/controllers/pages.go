package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Index serves the homepage.
func Index(c *fiber.Ctx) error {
	data := fiber.Map{
		"About":    models.About(),
		"Projects": models.GetProjects(),
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("indexPage", data, "")
	} else {
		return c.Render("indexPage", data)
	}
}

// About serves the about page.
func About(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"About":  models.AboutPage(),
		"Skills": models.Skills_Test(),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("aboutPage", data, "")
	} else {
		return c.Render("aboutPage", data)
	}
}

// Contact serves the contact page.
func Contact(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"Contact": models.ContactPage(),
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("contactPage", data, "")
	} else {
		return c.Render("contactPage", data)
	}
}
