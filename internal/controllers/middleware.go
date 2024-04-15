package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/kevingil/blog/internal/models"
)

// Tmpl is a template.
var Tmpl *template.Template

// Sessions is a user sessions.
var Sessions map[string]*models.User

// Tempalte context
type Context struct {
	User            *models.User
	Article         *models.Article
	Articles        []*models.Article
	Project         *models.Project
	Projects        []*models.Project
	Skill           *models.Project
	Skills          []*models.Skill
	About           string
	Contact         string
	ArticleCount    int
	DraftCount      int
	TemplateChild   string
	TotalArticles   int
	ArticlesPerPage int
	TotalPages      int
	CurrentPage     int
}

var data Context

// Render is a function to render a partial template if the request is an hx request
// or a partial with layout if it's a normal HTTP request
func render(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data Context) {
	var response bytes.Buffer
	var child bytes.Buffer
	permission(w, r)
	cookie := getCookie(r)
	data.User = Sessions[cookie.Value]

	if err := Tmpl.ExecuteTemplate(&child, tmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		io.WriteString(w, child.String())

	} else {
		data.TemplateChild = child.String()
		if err := Tmpl.ExecuteTemplate(&response, layout, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, response.String())

	}

}

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
