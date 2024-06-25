package controllers

import (
	"net/http"

	"github.com/kevingil/blog/internal/models"
)

// Index serves the homepage.
func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"About":    models.About(),
		"Projects": models.GetProjects(),
	}

	Handle(w, r, data)

}

// About serves the about page.
func About(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"About":  models.AboutPage(),
		"Skills": models.Skills_Test(),
	}

	Handle(w, r, data)
}

// Contact serves the contact page.
func Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Contact": models.ContactPage(),
	}

	Handle(w, r, data)
}
