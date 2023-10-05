package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/app/controllers"
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

	// HTMX components
	r.HandleFunc("/mx/", controllers.IndexMX)
	r.HandleFunc("/mx/contact", controllers.ContactMX)
	r.HandleFunc("/mx/post/{slug}", controllers.PostMX)

	// Other
	r.HandleFunc("/projects/moderatorjs", controllers.ModeratorJS)

	//Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Println(fmt.Sprintf("Your app is running on port %s.", os.Getenv("PORT")))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
