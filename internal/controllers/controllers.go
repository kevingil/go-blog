package controllers

import (
	"net/http"

	"github.com/kevingil/blog/internal/models"
)

func Index(w http.ResponseWriter, r *http.Request) {

	data.About = models.About()
	//data.Projects = models.HomeProjects()
	data.Projects = models.GetProjects()

	// Render the template using the utility function
	render(w, r, "layout", "index", data)
}

func About(w http.ResponseWriter, r *http.Request) {
	data.About = models.AboutPage()
	data.Skills = models.Skills_Test()
	render(w, r, "layout", "about", data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	render(w, r, "layout", "contact", data)
}
