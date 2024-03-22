package controllers

import (
	"html/template"
	"net/http"

	"github.com/kevingil/blog/app/models"
	"github.com/kevingil/blog/app/views"
)

type Services struct {
	MonthlyViews    int
	MonthlyVisitors int
	TopArticles     []int
}

type Context struct {
	User            *models.User
	Article         *models.Article
	Articles        []*models.Article
	Project         *models.Project
	Projects        []*models.Project
	Skill           *models.Project
	Skills          []*models.Skill
	Tags            []*models.Tag
	About           string
	Contact         string
	ArticleCount    int
	DraftCount      int
	View            template.HTML
	TotalArticles   int
	ArticlesPerPage int
	TotalPages      int
	CurrentPage     int
}

var data Context

func Index(w http.ResponseWriter, r *http.Request) {

	data.About = models.About()
	//data.Projects = models.HomeProjects()
	data.Projects = models.GetProjects()

	// Render the template using the utility function
	views.Render(w, r, "layout", "index", data)
}

func About(w http.ResponseWriter, r *http.Request) {
	data.About = models.AboutPage()
	data.Skills = models.Skills_Test()
	views.Render(w, r, "layout", "about", data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data.Contact = models.ContactPage()
	views.Render(w, r, "layout", "contact", data)
}

// This just handles the page, Moderator is written in JS
func ModeratorJS(w http.ResponseWriter, r *http.Request) {
	views.Render(w, r, "layout", "moderatorjs", data)
}
