package controllers

import (
	"bytes"
	"io"
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
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	switch r.Method {
	case http.MethodGet:
		req.Layout = "layout"
		req.Tmpl = "login"
		req.Data = nil
		render(req)
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
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
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
	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
	}
	permission(req)

	switch r.Method {
	case http.MethodGet:
		req.Tmpl = "register"
		req.Data = nil
		render(req)
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
func Profile(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
	user := Sessions[cookie.Value]

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		User     *models.User
		Skills   []*models.Skill
		Projects []*models.Project
	}{
		User:     user.GetProfile(),
		Skills:   user.GetSkills(),
		Projects: user.GetProjects(),
	}

	req.Data = data
	req.Layout = "dashboard-layout"

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
			req.Tmpl = "page_profile"
			if r.Header.Get("HX-Request") == "true" {
				req.Tmpl = "profile"
			}
			render(req)
		}

	case http.MethodPost:
		// Handle POST updates to user profile
		updatedUser := &models.User{
			ID:    user.ID,
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		user.UpdateUser(updatedUser)

		req.Data = data
		req.Tmpl = "profile-user"
		render(req)
	}
}

// editUser handles editing user profile.
func editUser(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	// Add logic specific to editing user profile here
	req.Tmpl = "edit-user"
	render(req)
}

// editContact handles editing contact information.
func editContact(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	// Add logic specific to editing contact information here
	req.Tmpl = "edit-contact"
	render(req)
}

// Skills handles skill-related operations.
func Skills(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
	user := Sessions[cookie.Value]

	data := struct {
		User     *models.User
		Skills   []*models.Skill
		Projects []*models.Project
	}{
		User:     user.GetProfile(),
		Skills:   user.GetSkills(),
		Projects: user.GetProjects(),
	}

	req.Data = data
	req.Layout = "dashboard-layout"
	req.Tmpl = "edit-skills"

	render(req)
}

type ResumeData struct {
	User     *models.User
	Skills   []*models.Skill
	Projects []*models.Project
}

func Resume(w http.ResponseWriter, r *http.Request) {

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")

	data := ResumeData{
		User:     user.GetProfile(),
		Skills:   user.GetSkills(),
		Projects: user.GetProjects(),
	}

	for _, skill := range data.Skills {
		println(skill.Name)
	}

	req := Request{
		W:    w,
		R:    r,
		Tmpl: "dashboard-resume",
	}

	permission(req)

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "skills":
			if delete != "" && id != 0 {
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
			} else {
				var response bytes.Buffer
				if err := Tmpl.ExecuteTemplate(&response, "edit-skill", data); err != nil {
					log.Printf("Template error for 'edit-skill': %v", err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				io.WriteString(w, response.String())
			}
		case "projects":
			if delete != "" && id != 0 {
				project := &models.Project{ID: id}
				user.DeleteProject(project)
				data.Projects = user.GetProjects()
				var response bytes.Buffer
				if r.Header.Get("HX-Request") == "true" {
					if err := Tmpl.ExecuteTemplate(&response, req.Tmpl, data); err != nil {
						log.Printf("Template error for '%s': %v", req.Tmpl, err)
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					io.WriteString(w, response.String())
					permission(req)
				} else {
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}
			} else {
				project := &models.Project{}
				if user != nil && id != 0 {
					project = user.FindProject(id)
				}
				var response bytes.Buffer
				if err := Tmpl.ExecuteTemplate(&response, "edit-project", project); err != nil {
					log.Printf("Template error for 'edit-project': %v", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			}
		default:
			var response bytes.Buffer
			if r.Header.Get("HX-Request") == "true" {
				if err := Tmpl.ExecuteTemplate(&response, req.Tmpl, data); err != nil {
					log.Printf("Template error for '%s': %v", req.Tmpl, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.WriteString(w, response.String())
				permission(req)
			} else {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}
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
				data.Projects = user.GetProjects()
				var response bytes.Buffer
				if r.Header.Get("HX-Request") == "true" {
					if err := Tmpl.ExecuteTemplate(&response, req.Tmpl, data); err != nil {
						log.Printf("Template error for '%s': %v", req.Tmpl, err)
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					io.WriteString(w, response.String())
					permission(req)
				} else {
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}
			}
		}
	}
}
