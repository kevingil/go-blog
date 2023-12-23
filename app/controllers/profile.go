package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/kevingil/blog/app/views"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")

	data.User = user.GetProfile()

	data.Skills = user.FindSkills()

	data.Projects = user.FindProjects()

	tmpl := "page_profile.gohtml"

	if r.Header.Get("HX-Request") == "true" {
		tmpl = "profile"
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "about":
			if delete != "" && id != 0 {
				//Not allowed
				//TODO: Add delete functionality as set to blank
				http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
			} else {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit_about", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			}
		default:
			if user != nil {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}
	case http.MethodPost:
		switch model {
		case "about":
			if user != nil {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "profile_about", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}
		//model := r.URL.Query().Get("model")
	}
}

func Skills(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = user.FindSkills()

	data.Projects = user.FindProjects()

	templateName := "edit_skills"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if err := views.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
	permission(w, r)
	//model := r.URL.Query().Get("model")
}

func Projects(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = user.FindSkills()

	data.Projects = user.FindProjects()

	templateName := "edit_projects"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if err := views.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
	permission(w, r)
	//model := r.URL.Query().Get("model")
}
