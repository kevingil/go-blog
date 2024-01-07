package controllers

import (
	"net/http"

	"github.com/kevingil/blog/app/models"
)

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	Hx(w, r, "main_layout", "contact", data)
}
