package utils

import (
	"html/template"
	"net/http"
)

func HxRender(w http.ResponseWriter, r *http.Request, htmxTemplate, normalTemplate string, data interface{}) error {
	if r.Header.Get("X-Hx-Request") == "true" {
		tmpl, err := template.ParseFiles(htmxTemplate)
		if err != nil {
			return err
		}
		return tmpl.Execute(w, data)
	}

	tmpl, err := template.ParseFiles(normalTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
