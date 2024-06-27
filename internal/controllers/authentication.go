package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Sessions is a map for user sessions.
var Store *session.Store

// GetUser is a helper function to get the current user from the session
func GetUser(c *fiber.Ctx) (*models.User, error) {
	sess, err := Store.Get(c)
	if err != nil {
		log.Println("store not found")
		return nil, err
	}

	userID := sess.Get("userID")
	if userID == nil {
		log.Println("userID nil")
		return nil, nil
	}

	userEmail := sess.Get("userEmail")
	if userEmail == nil {
		log.Println("userEmail nil")
		return nil, nil
	}

	user := &models.User{
		ID:    userID.(int),
		Email: userEmail.(string)}
	return user.Find(), nil
}

// LoginPage is a controller for the login page
func LoginPage(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if sess.Get("userID") != nil {
		return c.Redirect("/admin", fiber.StatusSeeOther)
	}

	data := fiber.Map{
		"User": "",
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("loginPage", data, "")
	}
	return c.Render("loginPage", data)
}

// AuthenticateUser is a controller to authenticate users
func AuthenticateUser(c *fiber.Ctx) error {
	user := &models.User{
		Email: c.FormValue("email"),
	}
	user = user.Find()

	if user.ID == 0 {
		log.Println("User not found for email:", c.FormValue("email"))
		return c.Status(fiber.StatusBadRequest).SendString("User not found")
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(c.FormValue("password")))
	if err != nil {
		log.Println("Password doesn't match:", user.Email, "Error:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	sess.Set("userID", user.ID)
	sess.Set("userEmail", user.Email)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/admin", fiber.StatusSeeOther)
}

// Logout is a controller for users to log out
func Logout(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/login", fiber.StatusSeeOther)
}

// Register is a controller to register a user
func Register(c *fiber.Ctx) error {
	user := &models.User{
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
	}
	password, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.MinCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := helpers.ValidateEmail(user.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user.Password = password
	user = user.Find()

	if user.ID == 0 {
		user = user.Create()
	}

	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	sess.Set("userID", user.ID)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/admin", fiber.StatusSeeOther)
}
