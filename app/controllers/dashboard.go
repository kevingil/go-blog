package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gosimple/slug"
	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/templates"
)

// Dashboard is a controller for users to list articles.
func Dashboard(w http.ResponseWriter, r *http.Request) {
	permission(w, r)

	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	model := r.URL.Query().Get("model")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	delete := r.URL.Query().Get("delete")
	tmpl := "page_dashboard.html"

	if r.Header.Get("HX-Request") == "true" {
		tmpl = "dashboard"
	}

	switch r.Method {
	case http.MethodGet:
		switch model {
		case "article":
			if delete != "" && id != 0 {
				article := &models.Article{
					ID: id,
				}

				user.DeleteArticle(article)

				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			} else {
				data.Article = &models.Article{
					Image:   "",
					Title:   "",
					Content: "",
				}

				if user != nil && id != 0 {
					data.Article = user.FindArticle(id)
				}

				var response bytes.Buffer
				if err := templates.Tmpl.ExecuteTemplate(&response, "article.html", data); err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				io.WriteString(w, response.String())
			}
		default:
			if user != nil {
				data.Articles = user.FindArticles()
			}

			var response bytes.Buffer
			if err := templates.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			io.WriteString(w, response.String())
		}
	case http.MethodPost:
		switch model {
		case "article":
			if user != nil {
				if id == 0 {
					article := &models.Article{
						Image:     r.FormValue("image"),
						Slug:      slug.Make(r.FormValue("title")),
						Title:     r.FormValue("title"),
						Content:   r.FormValue("content"),
						Author:    *user,
						CreatedAt: time.Now(),
					}

					user.CreateArticle(article)
				} else {
					article := &models.Article{
						ID:      id,
						Image:   r.FormValue("image"),
						Slug:    slug.Make(r.FormValue("title")),
						Title:   r.FormValue("title"),
						Content: r.FormValue("content"),
					}

					user.UpdateArticle(article)
				}
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	data.User = user.GetProfile()

	data.Skills = models.Skills_Test()

	data.Projects = models.Projects_Test()

	templateName := "profile"
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
}

func Articles(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	tmpl := "articles"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if user != nil {
			data.Articles = user.FindArticles()
		}
		if err := templates.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	io.WriteString(w, response.String())
}