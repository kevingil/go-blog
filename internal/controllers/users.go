package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Login is a controller for users to log in.
func Login(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	data := map[string]interface{}{}
	switch r.Method {
	case http.MethodGet:
		renderPage(w, r, data)
	case http.MethodPost:
		user := &models.User{
			Email: r.FormValue("email"),
		}
		user = user.Find()

		if user.ID == 0 {
			log.Println("User not found for email:", r.FormValue("email"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password")))
		if err != nil {
			log.Println("Password doesn't match:", user.Email, "Error:", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Logout is a controller for users to log out.
func Logout(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	if Sessions[cookie.Value] != nil {
		delete(Sessions, cookie.Value)
	}

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register is a controller to register a user.
func Register(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		renderPage(w, r, data)
	case http.MethodPost:
		user := &models.User{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := helpers.ValidateEmail(user.Email); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user.Password = password
		user = user.Find()

		if user.ID == 0 {
			user = user.Create()
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Profile is a controller for users to view and update their profile.
func DashboardProfile(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return // Important to return here to prevent further execution
	}

	model := r.URL.Query().Get("edit")
	delete := r.URL.Query().Get("delete")
	//new := r.URL.Query().Get("new")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Projects": user.GetProjects(),
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "user":
			if delete != "" && id != 0 {
				//Not allowed
				//TODO: Add delete functionality as set to blank
				http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
			} else {
				renderTemplate(w, r, data, "edit-user")
			}
		case "contact":
			renderTemplate(w, r, data, "edit-contact")
		default:
			renderPage(w, r, data)

		}
	case http.MethodPost:
		switch model {
		case "user":
			updatedUser := &models.User{
				ID:    user.ID,
				Name:  r.FormValue("name"),
				Email: r.FormValue("email"),
				About: r.FormValue("about"),
			}
			user.UpdateUser(updatedUser)
			data["User"] = user.GetProfile()
			renderTemplate(w, r, data, "profile-user")
		case "contact":
			updatedUser := &models.User{
				ID:      user.ID,
				Contact: r.FormValue("contact"),
			}
			user.UpdateContact(updatedUser)
			data["User"] = user.GetProfile()
			renderTemplate(w, r, data, "profile-contact")

		}

	}

}

// Skills renderPages skill-related operations.
func DashboardProfileSkills(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}
	permission(w, r)
	renderPage(w, r, data)
}
func DashboardResume(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	idStr := r.URL.Query().Get("id")
	delete := r.URL.Query().Get("delete")

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "projects":
			if delete != "" && idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					http.Error(w, "Invalid project ID", http.StatusBadRequest)
					return
				}
				project := &models.Project{ID: id}
				user.DeleteProject(project)
				data["Projects"] = user.GetProjects()
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
				return
			} else {
				if idStr != "" {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						http.Error(w, "Invalid project ID", http.StatusBadRequest)
						return
					}
					project := user.FindProject(id)
					if project == nil {
						http.Error(w, "Project not found", http.StatusNotFound)
						return
					}
					data["Project"] = project
				}
				renderTemplate(w, r, data, "edit-project")
			}
		default:
			renderPage(w, r, data)
		}
	case http.MethodPost:
		switch model {
		case "projects":
			if user != nil {
				url := r.FormValue("url")
				title := r.FormValue("title")
				classes := r.FormValue("classes")
				description := r.FormValue("description")

				// Validate form data
				if url == "" || title == "" || classes == "" || description == "" {
					http.Error(w, "Missing required fields", http.StatusBadRequest)
					return
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
						http.Error(w, "Invalid project ID", http.StatusBadRequest)
						return
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
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
				return
			}
		}
	}
}
