package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/database"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/services/coffeeapp"
)

func main() {
	//Init blog database
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//Initialize db for every visit
	err = database.Init()
	if err != nil {
		log.Fatal(err)
	}

	//In memory logged in sessions
	controllers.Sessions = make(map[string]*models.User)

	loadRoutes()

}

func loadRoutes() {
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

	// Moderator AI
	r.HandleFunc("/projects/moderatorjs", controllers.ModeratorJS)

	// Espresso App
	r.HandleFunc("/projects/coffeeapp", coffeeapp.CoffeeApp).Methods("GET")
	r.HandleFunc("/components/completion", coffeeapp.Completion).Methods("GET")
	r.HandleFunc("/api/stream-recipe", coffeeapp.StreamRecipe).Methods("POST")
	r.HandleFunc("/api/stream-recipe", coffeeapp.StreamRecipe).Methods("GET")

	//Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
