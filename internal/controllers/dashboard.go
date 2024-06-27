package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/internal/models"
)

// Dashboard
func AdminPage(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		log.Println("Admin store not found")
	}

	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if sess.Get("userID") != nil {

		data := map[string]interface{}{
			"ArticleCount": user.CountArticles(),
			"DraftCount":   user.CountDrafts(),
			"Articles":     user.FindArticles(),
			"User":         user,
		}
		if c.Get("HX-Request") == "true" {
			return c.Render("adminPage", data, "")
		} else {
			return c.Render("adminPage", data)
		}
	} else {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

}

// AdminProfile is a controller for users to view and update their profile.
func EditProfilePage(c *fiber.Ctx) error {

	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

	model := c.Query("edit")
	delete := c.Query("delete")
	id, _ := strconv.Atoi(c.Query("id"))

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Projects": user.GetProjects(),
	}

	switch model {
	case "user":
		if delete != "" && id != 0 {
			return c.Redirect("/admin/profile", fiber.StatusSeeOther)
		} else {
			return c.Render("edit-user", data, "")
		}
	case "contact":
		return c.Render("edit-contact", data, "")
	default:
		if c.Get("HX-Request") == "true" {
			return c.Render("adminProfilePage", data, "")
		} else {
			return c.Render("adminProfilePage", data)
		}
	}

}

// AdminProfile is a controller for users to view and update their profile.
func EditProfile(c *fiber.Ctx) error {

	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

	model := c.Query("edit")
	//delete := c.Query("delete")
	//id, _ := strconv.Atoi(c.Query("id"))

	switch model {
	case "user":
		data := map[string]interface{}{
			"User":     user.GetProfile(),
			"Projects": user.GetProjects(),
		}

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
		data := map[string]interface{}{
			"User": user.GetProfile(),
		}
		return c.Render("profile-contact", data, "")
	default:
		return c.Redirect("/admin/profile", fiber.StatusSeeOther)
	}

}

// AdminResume handles resume-related operations.
func EditResumePage(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("adminResumePage", data, "")
	} else {
		return c.Render("adminResumePage", data)
	}

}

func EditResumeProject(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	model := c.Query("edit")
	idStr := c.Query("id")
	delete := c.Query("delete")

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

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
			return c.Redirect("/admin/resume", fiber.StatusSeeOther)
		}
		if delete != "" && idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid project ID")
			}
			project := &models.Project{ID: id}
			user.DeleteProject(project)
			data["Projects"] = user.GetProjects()
			return c.Redirect("/admin/resume", fiber.StatusSeeOther)
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
}
