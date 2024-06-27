package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

// Login is a controller for users to log in.
func LoginPage(c *fiber.Ctx) error {
	user := &models.User{
		Email: c.FormValue("email"),
	}
	user = user.Find()
	data := map[string]interface{}{
		"User": "",
	}
	if user.ID != 0 {
		return c.Redirect("/dashboard", fiber.StatusSeeOther)
	} else {
		if c.Get("HX-Request") == "true" {
			return c.Render("loginPage", data, "")
		} else {
			return c.Render("loginPage", data)
		}
	}

}

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

	sessionID := uuid.New().String()
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    sessionID,
		HTTPOnly: false,
	})
	Sessions[sessionID] = user

	return c.Redirect("/dashboard", fiber.StatusSeeOther)
}

// Logout is a controller for users to log out
func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if Sessions[sessionID] != nil {
		delete(Sessions, sessionID)
	}

	c.ClearCookie("session")

	return c.Redirect("/login", fiber.StatusSeeOther)
}

// Register is a controller to register a user.
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

	sessionID := uuid.New().String()
	cookie := &fiber.Cookie{
		Name:  "session",
		Value: sessionID,
	}
	Sessions[sessionID] = user

	c.Cookie(cookie)
	return c.Redirect("/dashboard", fiber.StatusSeeOther)
}

// DashboardProfile is a controller for users to view and update their profile.
func DashboardProfile(c *fiber.Ctx) error {

	sessionID := c.Cookies("session")
	user := Sessions[sessionID]

	if user == nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

	model := c.Query("edit")
	delete := c.Query("delete")
	id, _ := strconv.Atoi(c.Query("id"))

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Projects": user.GetProjects(),
	}

	switch c.Method() {
	case fiber.MethodGet:
		switch model {
		case "user":
			if delete != "" && id != 0 {
				return c.Redirect("/dashboard/profile", fiber.StatusSeeOther)
			} else {
				return c.Render("edit-user", data, "")
			}
		case "contact":
			return c.Render("edit-contact", data, "")
		default:
			if c.Get("HX-Request") == "true" {
				return c.Render("dashboardProfilePage", data, "")
			} else {
				return c.Render("dashboardProfilePage", data)
			}
		}
	case fiber.MethodPost:
		switch model {
		case "user":
			updatedUser := &models.User{
				ID:    user.ID,
				Name:  c.FormValue("name"),
				Email: c.FormValue("email"),
				About: c.FormValue("about"),
			}
			user.UpdateUser(updatedUser)
			data["User"] = user.GetProfile()

			return c.Render("profile-user", data, "")
		case "contact":
			updatedUser := &models.User{
				ID:      user.ID,
				Contact: c.FormValue("contact"),
			}
			user.UpdateContact(updatedUser)
			data["User"] = user.GetProfile()
			return c.Render("profile-contact", data, "")
		}
	}

	return nil
}

// DashboardResume handles resume-related operations.
func DashboardResume(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	user := Sessions[sessionID]
	model := c.Query("edit")
	idStr := c.Query("id")
	delete := c.Query("delete")

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

	switch c.Method() {
	case fiber.MethodGet:
		switch model {
		case "projects":
			if delete != "" && idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).SendString("Invalid project ID")
				}
				project := &models.Project{ID: id}
				user.DeleteProject(project)
				data["Projects"] = user.GetProjects()
				return c.Redirect("/dashboard/resume", fiber.StatusSeeOther)
			} else {
				if idStr != "" {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).SendString("Invalid project ID")
					}
					project := user.FindProject(id)
					if project == nil {
						return c.Status(fiber.StatusNotFound).SendString("Project not found")
					}
					data["Project"] = project
				}
				return c.Render("edit-project", data, "")
			}
		default:
			if c.Get("HX-Request") == "true" {
				return c.Render("dashboardResumePage", data, "")
			} else {
				return c.Render("dashboardResumePage", data)
			}
		}
	case fiber.MethodPost:
		switch model {
		case "projects":
			if user != nil {
				url := c.FormValue("url")
				title := c.FormValue("title")
				classes := c.FormValue("classes")
				description := c.FormValue("description")

				if url == "" || title == "" || classes == "" || description == "" {
					return c.Status(fiber.StatusBadRequest).SendString("Missing required fields")
				}

				if idStr == "" {
					project := &models.Project{
						Url:         url,
						Title:       title,
						Classes:     classes,
						Description: description,
					}
					user.AddProject(project)
				} else {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).SendString("Invalid project ID")
					}
					project := &models.Project{
						ID:          id,
						Url:         url,
						Title:       title,
						Classes:     classes,
						Description: description,
					}
					user.UpdateProject(project)
				}
				data["Projects"] = user.GetProjects()
				return c.Redirect("/dashboard/resume", fiber.StatusSeeOther)
			}
		}
	}

	return nil
}
