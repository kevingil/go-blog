package controllers

import (
	"bytes"
	"io"
	"log"
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

const (
	PAGES = "internal/templates/pages"
)

// Render is a function to render a partial template if the request is an HX request
// or a full template with layout if it's a normal HTTP request.
// Handle serves the templates based on the URL.
func Handle(w http.ResponseWriter, r *http.Request, data map[string]any) {

	// Renders child template
	// then local layout, **/_layout.gohtml
	// then the root layout, /_layout.gohtml
	// unlesss already at root

	log.Printf("Request URL: %s", r.URL.Path)
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data["User"] = user

	url := r.URL.Path

	templatePath := PAGES + url + "index.gohtml"
	localLayoutPath := PAGES + url + "_layout.gohtml"
	rootLayoutPath := PAGES + "/_layout.gohtml"

	// Render the child template
	var htmlChildContent bytes.Buffer

	if err := Tmpl.ExecuteTemplate(&htmlChildContent, templatePath, data); err != nil {
		log.Printf("Error executing template %s: %v", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["TemplateChild"] = htmlChildContent.String()
	log.Printf("Rendered Child, Parsed URL path: %s", url)

	// Function to check if a template exists
	templateExists := func(name string) bool {
		_, err := Tmpl.ParseFiles(name)
		return err == nil
	}

	var htmlWrappedContent bytes.Buffer

	// For HTMX requests, just return the child template
	if r.Header.Get("HX-Request") == "true" {
		htmlWrappedContent = htmlChildContent
	} else {

		// Apply local layout if it exists
		if templateExists(localLayoutPath) {
			data["TemplateChild"] = htmlChildContent.String()
			if err := Tmpl.ExecuteTemplate(&htmlWrappedContent, localLayoutPath, data); err != nil {
				log.Printf("Error executing local layout %s: %v", localLayoutPath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Rendered with local layout")

		}

		// Apply root layout
		if url != "/" {
			data["TemplateChild"] = htmlWrappedContent.String()
			if err := Tmpl.ExecuteTemplate(&htmlWrappedContent, rootLayoutPath, data); err != nil {
				log.Printf("Error executing root layout %s: %v", rootLayoutPath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlWrappedContent.String())

}

// Index serves the homepage.
func Index(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering

	data := map[string]interface{}{
		"About":    models.About(),
		"Projects": models.GetProjects(),
	}

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
