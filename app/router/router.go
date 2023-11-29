package router

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
	r.HandleFunc("/contact", controllers.Contact)
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/register", controllers.Register)
	r.HandleFunc("/dashboard", controllers.Dashboard)
	r.HandleFunc("/dashboard/articles", controllers.Articles)
	r.HandleFunc("/dashboard/profile", controllers.Profile)
	r.HandleFunc("/post/{slug}", controllers.Post)

	// API Requests
	r.HandleFunc("/api/r2", controllers.R2)
	//r.HandleFunc("/api/profile", controllers.ProfileEdit).Methods("POST")
	//r.HandleFunc("/api/contact", controllers.ContactEdit).Methods("POST")
	r.HandleFunc("/api/skills", controllers.Skills).Methods("POST")
	r.HandleFunc("/api/skills/delete", controllers.Skills).Methods("POST")
	r.HandleFunc("/api/projects", controllers.Projects).Methods("POST")
	r.HandleFunc("/api/projects/delete", controllers.Projects).Methods("POST")

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
