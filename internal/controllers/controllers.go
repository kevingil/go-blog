package controllers

import (
	"net/http"

	"github.com/kevingil/blog/internal/models"
)

// Index serves the homepage.
func Index(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering
	data := struct {
		About    string
		Projects []*models.Project
	}{
		About:    models.About(),
		Projects: models.GetProjects(),
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "index",
		Data:   data,
	}

	render(req)
}

// About serves the about page.
func About(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering
	data := struct {
		About  string
		Skills []*models.Skill
	}{
		About:  models.AboutPage(),
		Skills: models.Skills_Test(),
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "about",
		Data:   data,
	}

	render(req) // render the about page with the provided data
}

// Contact serves the contact page.
func Contact(w http.ResponseWriter, r *http.Request) {
	// Prepare the data for rendering
	data := struct {
		Contact string
	}{
		Contact: models.ContactPage(),
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "contact",
		Data:   data,
	}

	render(req)
}
