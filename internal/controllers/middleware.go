package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

func getCookie(c *fiber.Ctx) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = c.Cookies("session")

	return cookie
}

func permission(c *fiber.Ctx) error {
	path := strings.Split(c.Path(), "/")[1]
	cookie := getCookie(c)

	switch path {
	case "dashboard":
		if Sessions[cookie.Value] == nil {
			return c.Redirect("/login", fiber.StatusSeeOther)
		}
	case "login", "register":
		if Sessions[cookie.Value] != nil {
			return c.Redirect("/dashboard", fiber.StatusSeeOther)
		}
	}

	return nil
}
