package cmd

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/views"
)

// Hx is a function to render a child template wrapped with a specified layout
func Hx(w http.ResponseWriter, r *http.Request, layout string, tmpl string, data any) {
	var response bytes.Buffer
	var child bytes.Buffer

	if err := views.Tmpl.ExecuteTemplate(&child, tmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		io.WriteString(w, child.String())

	} else {
		if err := views.Tmpl.ExecuteTemplate(&response, layout, child.String()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, response.String())

	}

}
