package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/kevingil/blog/app/templates"
)

func R2(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "r2"
	} else {
		templateName = "page_r2.html"
	}

	numRecentFiles := 10

	fileLinks := []string{}

	for i := 1; i <= numRecentFiles; i++ {
		fileName := fmt.Sprintf("file%d", i)
		fileLink := fmt.Sprintf("https://cdn.kevingil.com/%s", fileName)
		fileLinks = append(fileLinks, fileLink)
	}

	data := struct {
		FileLinks []string
	}{
		FileLinks: fileLinks,
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}
