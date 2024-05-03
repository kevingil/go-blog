package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/database"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
)

// Tempate files for embedded file system
//
//go:embed internal/templates/*.gohtml internal/templates/pages/*.gohtml internal/templates/forms/*.gohtml internal/templates/components/*.gohtml
var templateFS embed.FS

// Entrypoint

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Every app instance
	// initializes a database connection
	err = database.Init()
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}

	//Current session is stored in memory
	//users are authenticated via JWT cookies
	controllers.Sessions = make(map[string]*models.User)

	// Parse templates
	controllers.Tmpl = parseTemplates()

	// Start HTTP server
	serve()
}

func parseTemplates() *template.Template {
	// Create a new template and add helper functions
	tmpl := template.New("").Funcs(helpers.Functions)

	// Parse the templates from the embedded file system
	parsedTemplates, err := tmpl.ParseFS(templateFS,
		"internal/templates/*.gohtml",
		"internal/templates/pages/*.gohtml",
		"internal/templates/forms/*.gohtml",
		"internal/templates/components/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to parse embedded templates: %v", err)
	}

	return parsedTemplates
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
	r.HandleFunc("/dashboard/files", controllers.FilesPage)
	//Files =content with pagination
	r.HandleFunc("/dashboard/files/content", controllers.FilesContent)

	// Pages
	r.HandleFunc("/about", controllers.About)
	r.HandleFunc("/contact", controllers.Contact)

	//Files
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))
	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
