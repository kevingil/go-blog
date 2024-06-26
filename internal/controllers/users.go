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
		return
	}

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	switch r.Method {
	case http.MethodGet:
		model := r.URL.Query().Get("edit")
		delete := r.URL.Query().Get("delete")
		if model == "user" && delete != "" && id != 0 {
			http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
			return
		}

		switch model {
		case "user":
			editUser(w, r)
		case "contact":
			editContact(w, r)
		default:
			renderPage(w, r, data)
		}

	case http.MethodPost:
		// renderPage POST updates to user profile
		updatedUser := &models.User{
			ID:    user.ID,
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		user.UpdateUser(updatedUser)

		renderPartial(w, r, data, "profile-user")
	}
}

// editUser renderPages editing user profile.
func editUser(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	permission(w, r)
	renderPage(w, r, data)
}

// editContact renderPages editing contact information.
func editContact(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	permission(w, r)
	renderPage(w, r, data)

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
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")

	data := map[string]interface{}{
		"User":     user.GetProfile(),
		"Skills":   user.GetSkills(),
		"Projects": user.GetProjects(),
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "skills":
			if delete != "" && id != 0 {
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
			} else {

			}
		case "projects":
			if delete != "" && id != 0 {
				project := &models.Project{ID: id}
				user.DeleteProject(project)
				data["project"] = user.GetProjects()

			} else {
				data["project"] = &models.Project{}
				if user != nil && id != 0 {
					data["project"] = user.FindProject(id)
				}
			}
		default:
			renderPage(w, r, data)
		}
	case http.MethodPost:
		switch model {
		case "skills":
			log.Println("skills post test")
		case "projects":
			if user != nil {
				if id == 0 {
					project := &models.Project{
						Url:         r.FormValue("url"),
						Title:       r.FormValue("title"),
						Classes:     r.FormValue("classes"),
						Description: r.FormValue("description"),
					}
					user.AddProject(project)
				} else {
					project := &models.Project{
						ID:          id,
						Url:         r.FormValue("url"),
						Title:       r.FormValue("title"),
						Classes:     r.FormValue("classes"),
						Description: r.FormValue("description"),
					}
					user.UpdateProject(project)
				}
				data["project"] = user.GetProjects()
			}
		}
	}
}
