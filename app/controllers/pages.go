package controllers

import (
	"net/http"

	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

func About(w http.ResponseWriter, r *http.Request) {
	data.About = models.AboutPage()
	data.Skills = models.Skills_Test()
	views.Render(w, r, "main_layout", "about", data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	views.Render(w, r, "main_layout", "contact", data)
}

// This just handles the page, Moderator is written in JS
func ModeratorJS(w http.ResponseWriter, r *http.Request) {
	views.Render(w, r, "main_layout", "moderatorjs", data)
}
