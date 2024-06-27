package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/controllers"
)

func Permission(c *fiber.Ctx) error {
	path := strings.Split(c.Path(), "/")[1]

	switch path {
	case "dashboard":
		if controllers.Sessions[c.Cookies("session")] == nil {
			return c.Redirect("/login", fiber.StatusSeeOther)
		}
	case "login", "register":
		if controllers.Sessions[c.Cookies("session")] != nil {
			return c.Redirect("/dashboard", fiber.StatusSeeOther)
		}
	}

	return nil
}
