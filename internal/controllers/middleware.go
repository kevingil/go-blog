package controllers

import (
	"bytes"
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
)

// Tmpl is a template.
var Tmpl *template.Template

var app *mux.Router

// Sessions is a map for user sessions.
var Sessions map[string]*models.User

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Layout string // template layout
	Tmpl   string // template name
	User   *models.User
	Keys   map[string]any
}

// CustomRouter is a custom type that embeds mux.Router.
type CustomRouter struct {
	*mux.Router
}

// NewCustomRouter creates a new CustomRouter.
func NewCustomRouter() *CustomRouter {
	return &CustomRouter{
		Router: mux.NewRouter(),
	}
}

// Render is a function to render a partial template if the request is an HX request
// or a full template with layout if it's a normal HTTP request.
func (router *CustomRouter) Handle(route string, handler func(c Context) error) {
	router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		c := Context{
			W:    w,
			R:    r,
			Keys: make(map[string]any),
		}

		// Extract the session user
		cookie := getCookie(r)
		c.User = Sessions[cookie.Value]
		c.Keys["User"] = c.User

		// Handle the request
		if err := handler(c); err != nil {
			logging(c, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var rendered bytes.Buffer
		var child bytes.Buffer

		if err := Tmpl.ExecuteTemplate(&child, c.Tmpl, c.Keys); err != nil {
			logging(c, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Header.Get("HX-Request") == "true" {
			io.WriteString(w, child.String())
		} else {
			c.Keys["TemplateChild"] = child.String()
			if err := Tmpl.ExecuteTemplate(&rendered, c.Layout, c.Keys); err != nil {
				logging(c, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			io.WriteString(w, rendered.String())
		}
	})
}

func Serve() {
	r := mux.NewRouter()

	// Blog pages
	r.HandleFunc("/", Index)

	// User login, logout, register
	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/register", Register)

	// View posts, preview drafts
	r.HandleFunc("/blog", Blog)

	//Services
	r.HandleFunc("/blog/partial/recent", RecentPostsPartial)

	// View posts, preview drafts
	r.HandleFunc("/blog/{slug}", Post)

	// User Dashboard
	r.HandleFunc("/dashboard", Dashboard)

	// Edit articles, delete, or create new
	r.HandleFunc("/dashboard/publish", Publish)

	// View posts, preview drafts
	r.HandleFunc("/dashboard/publish/edit", EditArticle)

	// User Profile
	// Edit about me, skills, and projects
	r.HandleFunc("/dashboard/profile", Profile)

	// Resume Edit
	r.HandleFunc("/dashboard/resume", Resume)

	// Files page
	r.HandleFunc("/dashboard/files", FilesPage)
	//Files =content with pagination
	r.HandleFunc("/dashboard/files/content", FilesContent)

	// Pages
	r.HandleFunc("/about", About)
	r.HandleFunc("/contact", Contact)

	//Files
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))
	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func ParseTemplates(templateFS embed.FS) *template.Template {
	// Create a new template and add helper functions
	tmpl := template.New("").Funcs(helpers.Functions)

	// Parse the templates from the embedded file system
	parsedTemplates, err := tmpl.ParseFS(templateFS,
		"internal/templates/*.gohtml",
		"internal/templates/pages/*.gohtml",
		"internal/templates/dashboard/*.gohtml",
		"internal/templates/components/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to parse embedded templates: %v", err)
	}

	return parsedTemplates
}

// About serves the about page.
func About(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering
	data := struct {
		User   string
		About  string
		Skills []*models.Skill
	}{
		About:  models.AboutPage(),
		Skills: models.Skills_Test(),
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "about",
		Data:   data,
	}

	render(req) // render the about page with the provided data
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
