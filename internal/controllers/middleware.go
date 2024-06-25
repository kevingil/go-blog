package controllers

import (
	"bytes"
	"embed"
	"io"
	"log"
	"net/http"
	"path/filepath"
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

var Fs embed.FS

const (
	PAGES = "internal/templates/pages"
)

// Renders child template
// then local layout, **/_layout.gohtml
// then the root layout, /_layout.gohtml
// unlesss already at root
func Handle(w http.ResponseWriter, r *http.Request, data map[string]any) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data["User"] = user

	url := r.URL.Path

	templatePath := PAGES + url + "/index.gohtml"
	localLayoutPath := PAGES + url + "/_layout.gohtml"
	if url == "/" {
		templatePath = PAGES + "/index.gohtml"
		localLayoutPath = PAGES + "/_layout.gohtml"
	}
	rootLayoutPath := PAGES + "/_layout.gohtml"

	var htmlContent bytes.Buffer

	// For HTMX requests, just render the child template
	if r.Header.Get("HX-Request") == "true" {
		if err := Tmpl.ExecuteTemplate(&htmlContent, templatePath, data); err != nil {
			log.Printf("Error executing template %s: %v", templatePath, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Rendered template: %s", filepath.Base(templatePath))
	} else {
		// Check if local layout exists
		if _, err := Fs.ReadFile(localLayoutPath); err == nil {
			if err != nil {
				log.Printf("Error parsing local layout %s: %v", localLayoutPath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := Tmpl.ExecuteTemplate(&htmlContent, localLayoutPath, data); err != nil {
				log.Printf("Error executing local layout, not present? %s: %v", localLayoutPath, err)
			}
			log.Printf("Rendered template: %s", filepath.Base(localLayoutPath))
		} else {
			if err := Tmpl.ExecuteTemplate(&htmlContent, templatePath, data); err != nil {
				log.Printf("Error executing template %s: %v", templatePath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Rendered template: %s", templatePath)
		}

		// Apply root layout if it's not the root URL
		if url != "/" {
			data["TemplateChild"] = htmlContent.String()
			var rootContent bytes.Buffer

			if err := Tmpl.ExecuteTemplate(&rootContent, rootLayoutPath, data); err != nil {
				log.Printf("Error executing root layout %s: %v", rootLayoutPath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Rendered template: %s", rootLayoutPath)
			htmlContent = rootContent
		}
	}

	log.Printf("Path: %s", url)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlContent.String())
	data["TemplateChild"] = ""
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
