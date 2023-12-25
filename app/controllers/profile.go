package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("edit")
	delete := r.URL.Query().Get("delete")
	//new := r.URL.Query().Get("new")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	data.User = user.GetProfile()

	data.Skills = user.GetSkills()

	data.Projects = user.GetProjects()

	tmpl := "page_profile.gohtml"

	if r.Header.Get("HX-Request") == "true" {
		tmpl = "profile"
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "user":
			if delete != "" && id != 0 {
				//Not allowed
				//TODO: Add delete functionality as set to blank
				http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
			} else {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit_user", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			}
		case "contact":
			if user != nil {
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "edit_contact", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
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
		case "user":
			if user != nil {
				updatedUser := &models.User{
					ID:    user.ID,
					Name:  r.FormValue("name"),
					Email: r.FormValue("email"),
					About: r.FormValue("about"),
				}
				user.UpdateUser(updatedUser)
				data.User = user.GetProfile()
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "profile_user", data); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				io.WriteString(w, response.String())
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		case "contact":
			if user != nil {
				updatedUser := &models.User{
					ID:      user.ID,
					Contact: r.FormValue("contact"),
				}
				user.UpdateContact(updatedUser)
				data.User = user.GetProfile()
				var response bytes.Buffer
				if err := views.Tmpl.ExecuteTemplate(&response, "profile_contact", data); err != nil {
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
