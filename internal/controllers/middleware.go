package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/kevingil/blog/internal/models"
)

type Request struct {
	W      http.ResponseWriter
	R      *http.Request
	Layout string
	Tmpl   string
	User   *models.User
	Data   interface{}
}

func render(r Request) {

}

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

var Tmpl *template.Template

// Render is a function to render a partial template if the request is an HX request
// or a full template with layout if it's a normal HTTP request.
func Handle(w http.ResponseWriter, r *http.Request, data map[string]any) {

	var template string
	var layout string

	// Extract the session user
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data["User"] = user

	var rendered bytes.Buffer
	var child bytes.Buffer

	if err := Tmpl.ExecuteTemplate(&child, template, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		io.WriteString(w, child.String())
	} else {
		data["TemplateChild"] = child.String()
		if err := Tmpl.ExecuteTemplate(&rendered, layout, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, rendered.String())
	}
}

// Index serves the homepage.
func Index(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering
	var data map[string]any
	data["About"] = models.About()
	data["Projects"] = models.GetProjects()

	Handle(w, r, data)

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
	case "login", "register":
		if Sessions[cookie.Value] != nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}

/*
func logging(w http.ResponseWriter, r *http.Request, err error) {
	// Log request method
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	// Log form parameters
	r.ParseForm()
	if r.Form != nil {
		log.Printf("Form: %v", r.Form)
	}

	if err != nil {
		log.Println("Error:", err.Error())
	}
}
*/
