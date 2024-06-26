package server

import "github.com/gofiber/fiber/v2"

func LayoutMiddleware(c *fiber.Ctx) error {
	if c.Get("HX-Request") != "" {
		return c.Next()
	}

	return c.Render("layout", fiber.Map{
		"TemplateChild": c.Locals("content"),
	})
}
