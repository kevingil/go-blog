package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/services/espresso"
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
	r.HandleFunc("/post/{slug}", controllers.Post)

	// Moderator AI
	r.HandleFunc("/projects/moderatorjs", controllers.ModeratorJS)

	// Espresso App
	r.HandleFunc("/projects/espressoai", espresso.EspressoApp)
	r.HandleFunc("/api/request-coffee-recipe/", espresso.GenerateResponse).Methods("POST")
	r.HandleFunc("/api/load-coffee-recipe/", espresso.LoadAIResponse).Methods("GET")
	r.HandleFunc("/api/request-coffee-recipe-followup/", espresso.GenerateResponse).Methods("POST")
	r.HandleFunc("/api/load-coffee-recipe/", espresso.LoadAIResponse).Methods("GET")

	//Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
