package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kevingil/blog/app/views"
)

func Resume(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")
	data.User = user.GetProfile()
	data.Skills = user.GetSkills()
	data.Projects = user.GetProjects()

	templateName := "dash_contact"

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "skills":
			//skills edit
			if delete != "" && id != 0 {
				//Not allowed
				//TODO: Add delete functionality as set to blank
				http.Redirect(w, r, "/dashboard/resume", http.StatusSeeOther)
			} else {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit_skill", data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				io.WriteString(w, response.String())
			}
		case "projects":
			var response bytes.Buffer
			if err := views.Tmpl.ExecuteTemplate(&response, "edit_project", data); err != nil {
				log.Fatal("Template Error:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			io.WriteString(w, response.String())
		default:
			//default template
			var response bytes.Buffer
			if r.Header.Get("HX-Request") == "true" {
				if err := views.Tmpl.ExecuteTemplate(&response, templateName, data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				io.WriteString(w, response.String())
				permission(w, r)

			} else {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}

		}
	case http.MethodPost:
		switch model {
		case "skills":
			//exec skills edit
			log.Println("skills post test")
		case "projects":
			//exec projects edit
			log.Println("proj post test")
		}

		//model := r.URL.Query().Get("model")
	}

}

func Skills(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

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

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

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
