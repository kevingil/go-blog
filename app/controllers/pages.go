package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "contact"
	} else {
		templateName = "page_contact.gohtml"
	}

	var response bytes.Buffer

	if err := views.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}
