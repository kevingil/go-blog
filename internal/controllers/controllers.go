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
		//TODO cache this
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
