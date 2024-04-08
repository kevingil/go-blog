package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kevingil/blog/app/helpers"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
	"golang.org/x/crypto/bcrypt"
)

// Sessions is a user sessions.
var Sessions map[string]*models.User

func getCookie(r *http.Request) *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: "",
	}

	for _, c := range r.Cookies() {
		if c.Name == "session" {
			cookie.Value = c.Value
			break
		}
	}

	return cookie
}

func permission(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")[1]
	cookie := getCookie(r)

	switch path {
	case "dashboard":
		if Sessions[cookie.Value] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	case "login":
	case "register":
		if Sessions[cookie.Value] != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}

// Login is a controller for users to log in.
func Login(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "login.gohtml", nil); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
	case http.MethodPost:
		user := &models.User{
			Email: r.FormValue("email"),
		}
		user = user.Find()

		if user.ID == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password")))
		if err != nil {
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
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := views.Tmpl.ExecuteTemplate(&response, "register.gohtml", nil); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		io.WriteString(w, response.String())
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

func Profile(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	delete := r.URL.Query().Get("delete")
	//new := r.URL.Query().Get("new")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	data.User = user.GetProfile()

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

	tmpl := "page_profile.gohtml"

	if r.Header.Get("HX-Request") == "true" {
		tmpl = "profile"
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
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit-user", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			}
		case "contact":
			if user != nil {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit-contact", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		default:
			if user != nil {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}
	case http.MethodPost:
		switch model {
		case "user":
			if user != nil {
				updatedUser := &models.User{
					ID:    user.ID,
					Name:  r.FormValue("name"),
					Email: r.FormValue("email"),
					About: r.FormValue("about"),
				}
				user.UpdateUser(updatedUser)
				data.User = user.GetProfile()
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "profile-user", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		case "contact":
			if user != nil {
				updatedUser := &models.User{
					ID:      user.ID,
					Contact: r.FormValue("contact"),
				}
				user.UpdateContact(updatedUser)
				data.User = user.GetProfile()
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "profile-contact", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}

		}

		//model := r.URL.Query().Get("model")
	}

}

func Resume(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")
	data.User = user.GetProfile()
	data.Skills = user.GetSkills()
	data.Projects = user.GetProjects()

	tmpl := "dashboard-resume"

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "skills":
			//skills edit
			if delete != "" && id != 0 {
				//Not allowed
				//TODO: Add delete functionality as set to blank
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
			} else {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit-skill", data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				io.WriteString(w, response.String())
			}
		case "projects":
			if delete != "" && id != 0 {
				project := &models.Project{
					ID: id,
				}

				user.DeleteProject(project)
				data.Projects = user.GetProjects()
				var response bytes.Buffer
				if r.Header.Get("HX-Request") == "true" {

					if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
						log.Printf("Delete Project: %v", project.ID)

						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}

					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					io.WriteString(w, response.String())
					permission(w, r)

				} else {
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}

			} else {
				project := &models.Project{
					ID:          0,
					Title:       "",
					Url:         "",
					Description: "",
				}

				if user != nil && id != 0 {
					project = user.FindProject(id)
				}
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit-project", project); err != nil {
					log.Fatal("Template Error:", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())

			}
		default:
			//default template
			var response bytes.Buffer
			if r.Header.Get("HX-Request") == "true" {
				if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.WriteString(w, response.String())
				permission(w, r)

			} else {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}

		}
	case http.MethodPost:
		switch model {
		case "skills":
			//exec skills edit
			log.Println("skills post test")
		case "projects":
			//default template
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
					if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}

					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					io.WriteString(w, response.String())
					permission(w, r)

				} else {
					http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				}

			}
		}

		//model := r.URL.Query().Get("model")
	}

}

func Skills(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

	tmpl := "edit-skills"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
	permission(w, r)
	//model := r.URL.Query().Get("model")
}

func Projects(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

	tmpl := "edit-projects"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
	permission(w, r)
	//model := r.URL.Query().Get("model")
}
