package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/services/coffeeapp"
)

func Init() {
	r := mux.NewRouter()

	// Blog pages
	r.HandleFunc("/", controllers.Index)

	//Services
	r.HandleFunc("/service/feed", controllers.HomeFeed)

	// User login, logout, register
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/register", controllers.Register)

	// View posts, preview drafts
	r.HandleFunc("/articles", controllers.Articles)

	// View posts, preview drafts
	r.HandleFunc("/article/{slug}", controllers.Article)

	// User Dashboard
	r.HandleFunc("/dashboard", controllers.Dashboard)

	// Edit articles, delete, or create new
	r.HandleFunc("/dashboard/publish", controllers.Publish)

	// View posts, preview drafts
	r.HandleFunc("/dashboard/publish/edit", controllers.Editor)

	// User Profile
	// Edit about me, skills, and projects
	r.HandleFunc("/dashboard/profile", controllers.Profile)

	// Resume Edit
	r.HandleFunc("/dashboard/resume", controllers.Resume)

	// Pages
	//r.HandleFunc("/about", controllers.About)
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
