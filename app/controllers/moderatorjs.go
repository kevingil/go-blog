package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/templates"
)

// This just handles the page, Moderator is written in JS

func ModeratorJS(w http.ResponseWriter, r *http.Request) {

	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "moderatorjs"
	} else {
		templateName = "page_moderatorjs.html"
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}
