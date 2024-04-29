package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kevingil/blog/internal/models"
)

// Data structure for the Publish page
type PublishData struct {
	Articles []*models.Article
}

// Refactor the Publish function
func Publish(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	req := Request{
		W:      w,
		R:      r,
		Layout: "dashboard-layout",
		Tmpl:   "publish",
		Data: PublishData{
			Articles: user.FindArticles(),
		},
	}

	permission(req)
	render(req)
}

// Data structure for the EditArticle page
type EditArticleData struct {
	Article *models.Article
}

// Refactor the EditArticle function
func EditArticle(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(r)
	user := Sessions[cookie.Value]

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "edit-article",
	}

	permission(req)

	defaultArticle := &models.Article{
		Image:   "",
		Title:   "",
		Content: "",
		IsDraft: 0,
		Tags:    []*models.Tag{},
	}

	if user != nil && id != 0 {
		article, err := user.FindArticle(id)
		if err == nil {
			req.Data = EditArticleData{
				Article: article,
			}
		} else {
			log.Print(err)
			req.Data = EditArticleData{
				Article: defaultArticle,
			}
		}
	} else {
		req.Data = EditArticleData{
			Article: defaultArticle,
		}
	}

	render(req)
}

// Data structure for the Blog page
type BlogData struct {
	Articles        []*models.Article
	TotalArticles   int
	ArticlesPerPage int
	TotalPages      int
	CurrentPage     int
}

// Refactor the Blog function
func Blog(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	articlesPerPage := 10
	result, err := models.BlogTimeline(page, articlesPerPage)
	if err != nil {
		http.Error(w, "Error fetching blog timeline", http.StatusInternalServerError)
		return
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "blog",
		Data: BlogData{
			Articles:        result.Articles,
			TotalArticles:   result.TotalArticles,
			ArticlesPerPage: articlesPerPage,
			TotalPages:      result.TotalPages,
			CurrentPage:     result.CurrentPage,
		},
	}

	render(req)
}

// Data structure for the HomeFeedService
type HomeFeedData struct {
	Articles []*models.Article
}

// Refactor the HomeFeedService function
func HomeFeedService(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	if isHTMXRequest {
		req := Request{
			W:      w,
			R:      r,
			Layout: "",
			Tmpl:   "home-feed",
			Data: HomeFeedData{
				Articles: models.LatestArticles(6),
			},
		}
		render(req)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Data structure for the Post page
type PostData struct {
	Article *models.Article
}

// Refactor the Post function
func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := models.FindArticle(vars["slug"])

	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
		Tmpl:   "post",
		Data: PostData{
			Article: &models.Article{
				Image:   "",
				Title:   "Post Not Found",
				Content: "This post doesn't exist.",
			},
		},
	}

	if article != nil {
		req.Data = PostData{
			Article: article,
		}
	}

	render(req)
}
