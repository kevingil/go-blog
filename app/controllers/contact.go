package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/templates"
)

func ContactEdit(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	templateName := "dash_contact"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if err := templates.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
	permission(w, r)

}