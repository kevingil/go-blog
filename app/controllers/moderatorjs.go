package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/templates"
)

// This just handles the page, ModAI is written in JS
func ModeratorJS(w http.ResponseWriter, r *http.Request) {
	var response bytes.Buffer
	if err := templates.Tmpl.ExecuteTemplate(&response, "moderatorai.htmx", data); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	io.WriteString(w, response.String())

}
