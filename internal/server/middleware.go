package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/controllers"
)

func LayoutMiddleware(c *fiber.Ctx) error {
	if c.Get("HX-Request") != "" {
		return c.Render("layout", fiber.Map{
			"TemplateChild": c.Locals("content"),
		})
	}
	return c.Next()
}

func GetCookie(c *fiber.Ctx) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = c.Cookies("session")

	return cookie
}

func Permission(c *fiber.Ctx) error {
	path := strings.Split(c.Path(), "/")[1]
	cookie := GetCookie(c)

	switch path {
	case "dashboard":
		if controllers.Sessions[cookie.Value] == nil {
			return c.Redirect("/login", fiber.StatusSeeOther)
		}
	case "login", "register":
		if controllers.Sessions[cookie.Value] != nil {
			return c.Redirect("/dashboard", fiber.StatusSeeOther)
		}
	}

	return nil
}
