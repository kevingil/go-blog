package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/database"
	"github.com/kevingil/blog/internal/models"
	"github.com/kevingil/blog/internal/views"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//Initialize db conn
	err = database.Init()
	if err != nil {
		log.Fatal(err)
	}

	//Sessions
	controllers.Sessions = make(map[string]*models.User)
	controllers.Tmpl = parseTemplates()
	serve()
}

// Parse templates, checks for errors
func parseTemplates() *template.Template {
	dirs := []string{
		"./internal/views/*.gohtml",
		"./internal/views/pages/*.gohtml",
		"./internal/views/forms/*.gohtml",
		"./internal/views/components/*.gohtml",
	}

	//Parse templates with helper functions
	t := template.New("").Funcs(views.Functions)
	for _, dir := range dirs {
		files, err := filepath.Glob(dir)
		if err != nil {
			log.Fatalf("Failed to find template files: %v", err)
		}

		for _, file := range files {
			_, err = t.ParseFiles(file)
			if err != nil {
				log.Fatalf("Failed to parse template file (%s): %v", file, err)
			}
		}
	}
	return t
}

func serve() {
	r := mux.NewRouter()

	// Blog pages
	r.HandleFunc("/", controllers.Index)

	//Services
	r.HandleFunc("/service/feed", controllers.HomeFeedService)

	// User login, logout, register
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/register", controllers.Register)

	// View posts, preview drafts
	r.HandleFunc("/blog", controllers.Blog)

	// View posts, preview drafts
	r.HandleFunc("/blog/{slug}", controllers.Post)

	// User Dashboard
	r.HandleFunc("/dashboard", controllers.Dashboard)

	// Edit articles, delete, or create new
	r.HandleFunc("/dashboard/publish", controllers.Publish)

	// View posts, preview drafts
	r.HandleFunc("/dashboard/publish/edit", controllers.EditArticle)

	// User Profile
	// Edit about me, skills, and projects
	r.HandleFunc("/dashboard/profile", controllers.Profile)

	// Resume Edit
	r.HandleFunc("/dashboard/resume", controllers.Resume)

	// Files page
	r.HandleFunc("/dashboard/files", controllers.Files)

	// Pages
	r.HandleFunc("/about", controllers.About)
	r.HandleFunc("/contact", controllers.Contact)

	//Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}