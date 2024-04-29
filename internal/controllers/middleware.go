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

// Tmpl is a template.
var Tmpl *template.Template

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

var response Response

type Request struct {
	W      http.ResponseWriter
	R      *http.Request
	Layout string //template layout
	Tmpl   string //template name
	User   *models.User
	Data   interface{}
}

type Response struct {
	User          *models.User
	TemplateChild string //Rendered HTML template child
	Data          interface{}
}

// Render is a function to render a partial template if the request is an HX request
// or a full template with layout if it's a normal HTTP request.
func render(req Request) {
	var rendered bytes.Buffer
	var child bytes.Buffer

	permission(req)
	cookie := getCookie(req.R)
	response = Response{
		User:          Sessions[cookie.Value],
		TemplateChild: "",
		Data:          req.Data,
	}

	if err := Tmpl.ExecuteTemplate(&child, req.Tmpl, response); err != nil {
		logging(req, err)
		http.Error(req.W, err.Error(), http.StatusInternalServerError)
		return
	}

	req.W.Header().Set("Content-Type", "text/html; charset=utf-8")

	if req.R.Header.Get("HX-Request") == "true" {
		io.WriteString(req.W, child.String())
	} else {
		response.TemplateChild = child.String()
		if err := Tmpl.ExecuteTemplate(&rendered, req.Layout, response); err != nil {
			logging(req, err)
			http.Error(req.W, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(req.W, rendered.String())
	}

	logging(req, nil)
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

func permission(req Request) {
	path := strings.Split(req.R.URL.Path, "/")[1]
	cookie := getCookie(req.R)

	switch path {
	case "dashboard":
		if Sessions[cookie.Value] == nil {
			http.Redirect(req.W, req.R, "/login", http.StatusSeeOther)
		}
	case "login", "register":
		if Sessions[cookie.Value] != nil {
			http.Redirect(req.W, req.R, "/dashboard", http.StatusSeeOther)
		}
	}
}

func logging(req Request, err error) {
	// Log request method
	log.Printf("Request: %s %s", req.R.Method, req.R.URL.Path)

	// Log form parameters
	req.R.ParseForm()
	if req.R.Form != nil {
		log.Printf("Form: %v", req.R.Form)
	}

	// Log session user
	if req.User != nil {
		log.Printf("User: %s", req.User.Name)
	}

	if err != nil {
		log.Println("Error:", err.Error())
	}
}
