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
	PAGES   = "pages"
	PARTIAL = "partial"
)

// Renders child template
// then local layout, **/layout.gohtml
// then the root layout, /layout.gohtml
// unlesss already at root
func renderPage(w http.ResponseWriter, r *http.Request, data map[string]any) {
	permission(w, r)
	cookie := getCookie(r)
	var (
		INDEX  string = "index"
		LAYOUT string = "layout"
	)

	user := Sessions[cookie.Value]
	data["User"] = user

	if template, ok := data["Template"].(string); ok {
		INDEX = template
	}

	if layout, ok := data["Layout"].(string); ok {
		LAYOUT = layout
	}

	url := r.URL.Path

	if Url, ok := data["Url"].(string); ok {
		url = Url
	}

	isHXRequest := r.Header.Get("HX-Request") == "true"

	templatePath := filepath.Join(PAGES, url, INDEX+".gohtml")
	localLayoutPath := filepath.Join(PAGES, url, LAYOUT+".gohtml")
	if url == "/" {
		templatePath = filepath.Join(PAGES, INDEX+".gohtml")
		localLayoutPath = filepath.Join(PAGES, LAYOUT+".gohtml")
	}
	rootLayoutPath := filepath.Join(PAGES, LAYOUT+".gohtml")

	isRoot := (rootLayoutPath == localLayoutPath)

	log.Println(templatePath)
	log.Println(localLayoutPath)
	log.Println(rootLayoutPath)

	log.Println("Available templates:")
	for _, tmpl := range Tmpl.Templates() {
		log.Printf("- %s", tmpl.Name())
	}

	var htmlContent bytes.Buffer

	// Render the child template
	if err := Tmpl.ExecuteTemplate(&htmlContent, templatePath, data); err != nil {
		log.Printf("Error executing template %s: %v", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Rendered main template: %s", templatePath)

	if isHXRequest {
		// If local layout exists and render and wrap child, unless Url req has same prefix
		if _, err := Fs.ReadFile(localLayoutPath); err == nil && !isRoot && !strings.HasPrefix(url, filepath.Dir(r.URL.Path)) {
			var localContent bytes.Buffer
			data["TemplateChild"] = htmlContent.String()
			if err := Tmpl.ExecuteTemplate(&localContent, localLayoutPath, data); err != nil {
				log.Printf("Error executing local layout %s: %v", localLayoutPath, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Rendered local layout template: %s", localLayoutPath)
			htmlContent = localContent
		}
	} else {
		// Apply root layout if it's not the root URL and it's not an HX-Request
		var rootContent bytes.Buffer
		data["TemplateChild"] = htmlContent.String()
		if err := Tmpl.ExecuteTemplate(&rootContent, rootLayoutPath, data); err != nil {
			log.Printf("Error executing root layout %s: %v", rootLayoutPath, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Rendered root layout template: %s", rootLayoutPath)
		htmlContent = rootContent
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlContent.String())
}

func renderTemplate(w http.ResponseWriter, r *http.Request, data map[string]any, tmpl string) {
	permission(w, r)
	cookie := getCookie(r)
	var (
		LAYOUT string = "layout"
	)

	user := Sessions[cookie.Value]
	data["User"] = user

	if layout, ok := data["Layout"].(string); ok {
		LAYOUT = layout
	}

	isHXRequest := r.Header.Get("HX-Request") == "true"
	rootLayoutPath := filepath.Join(PAGES, LAYOUT+".gohtml")

	log.Println(tmpl)
	log.Println(rootLayoutPath)

	var htmlContent bytes.Buffer

	// Render the child template
	if err := Tmpl.ExecuteTemplate(&htmlContent, tmpl, data); err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Rendered main template: %s", tmpl)

	if !isHXRequest {
		// Apply root layout if it's not the root URL and it's not an HX-Request
		var rootContent bytes.Buffer
		data["TemplateChild"] = htmlContent.String()
		if err := Tmpl.ExecuteTemplate(&rootContent, rootLayoutPath, data); err != nil {
			log.Printf("Error executing root layout %s: %v", rootLayoutPath, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Rendered root layout template: %s", rootLayoutPath)
		htmlContent = rootContent
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlContent.String())
}

func renderPartial(w http.ResponseWriter, r *http.Request, data map[string]any, tmpl string) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data["User"] = user

	log.Println("Available templates:")
	for _, tmpl := range Tmpl.Templates() {
		log.Printf("- %s", tmpl.Name())
	}

	var htmlContent bytes.Buffer

	if err := Tmpl.ExecuteTemplate(&htmlContent, tmpl, data); err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Rendered partial template: %s", tmpl)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, htmlContent.String())
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
