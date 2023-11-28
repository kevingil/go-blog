package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/helpers"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/templates"
	"golang.org/x/crypto/bcrypt"
)

var data struct {
	User     *models.User
	Article  *models.Article
	Articles []*models.Article
	Projects []*models.Project
	Skills   []*models.Skill
}

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
func Index(w http.ResponseWriter, r *http.Request) {

	data.Skills = models.HomeSkills()
	data.Articles = models.Articles()
	data.Projects = models.HomeProjects()
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "home"
	} else {
		templateName = "index.html"
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}

func Contact(w http.ResponseWriter, r *http.Request) {

	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "contact"
	} else {
		templateName = "page_contact.html"
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}

// Post is the post/article controller.
func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := models.FindArticle(vars["slug"])

	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "post-content.html"
	} else {
		templateName = "single.html"
	}

	if article == nil {
		data.Article = &models.Article{
			Image:   "",
			Title:   "",
			Content: "Post Not Found",
		}
	} else {
		data.Article = article
	}

	if err := templates.Tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Register is a controller to register a user.
func Register(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := templates.Tmpl.ExecuteTemplate(&response, "register.html", nil); err != nil {
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

// Login is a controller for users to log in.
func Login(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	switch r.Method {
	case http.MethodGet:
		var response bytes.Buffer
		if err := templates.Tmpl.ExecuteTemplate(&response, "login.html", nil); err != nil {
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
