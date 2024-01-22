package controllers

import (
	"net/http"

	"github.com/kevingil/blog/app/models"
)

func About(w http.ResponseWriter, r *http.Request) {
	data.About = models.AboutPage()
	Hx(w, r, "main_layout", "about", data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	Hx(w, r, "main_layout", "contact", data)
}

// This just handles the page, Moderator is written in JS
func ModeratorJS(w http.ResponseWriter, r *http.Request) {
	Hx(w, r, "main_layout", "moderatorjs", data)
}
