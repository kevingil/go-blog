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

	renderPage(w, r, data)

}

// About serves the about page.
func About(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"About":  models.AboutPage(),
		"Skills": models.Skills_Test(),
	}

	renderPage(w, r, data)
}

// Contact serves the contact page.
func Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Contact": models.ContactPage(),
	}

	renderPage(w, r, data)
}
