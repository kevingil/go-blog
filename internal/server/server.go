package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kevingil/blog/internal/controllers"
)

func Boot() {

	var app = http.NewServeMux()

	// User login, logout, register
	app.HandleFunc("/login", controllers.Login)
	app.HandleFunc("/logout", controllers.Logout)
	app.HandleFunc("/register", controllers.Register)

	// View posts, preview drafts
	app.HandleFunc("/blog", controllers.Blog)

	//Services
	app.HandleFunc("/blog/partial/recent", controllers.RecentPostsPartial)

	// View posts, preview drafts
	app.HandleFunc("/blog/{slug}", controllers.Post)

	// User Dashboard
	app.HandleFunc("/dashboard", controllers.Dashboard)

	// Edit articles, delete, or create new
	app.HandleFunc("/dashboard/publish", controllers.Publish)

	// View posts, preview drafts
	app.HandleFunc("/dashboard/publish/edit", controllers.EditArticle)

	// User Profile
	// Edit about me, skills, and projects
	app.HandleFunc("/dashboard/profile", controllers.Profile)

	// Resume Edit
	app.HandleFunc("/dashboard/resume", controllers.Resume)

	// Files page
	app.HandleFunc("/dashboard/files", controllers.FilesPage)
	//Files =content with pagination
	app.HandleFunc("/dashboard/files/content", controllers.FilesContent)

	// Pages
	app.HandleFunc("/about", controllers.About)
	app.HandleFunc("/contact", controllers.Contact)

	// Combine file server and index handler
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index for root path
		if r.URL.Path == "/" {
			controllers.Index(w, r)
			return
		}

		// Check if the requested file exists
		path := filepath.Join("web", r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// If the file doesn't exist, serve the index page
			controllers.Index(w, r)
			return
		}

		// If the file exists, serve it
		http.FileServer(http.Dir("web")).ServeHTTP(w, r)
	})

	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), app))
}
