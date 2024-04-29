package controllers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Login is a controller for users to log in.
func Login(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	switch r.Method {
	case http.MethodGet:
		req.Layout = "layout"
		req.Tmpl = "login"
		req.Data = nil
		render(req)
	case http.MethodPost:
		user := &models.User{
			Email: r.FormValue("email"),
		}
		user = user.Find()

		if user.ID == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password")))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Logout is a controller for users to log out.
func Logout(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
	if Sessions[cookie.Value] != nil {
		delete(Sessions, cookie.Value)
	}

	cookie = &http.Cookie{
		Name:  "session",
		Value: "",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register is a controller to register a user.
func Register(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W:      w,
		R:      r,
		Layout: "layout",
	}
	permission(req)

	switch r.Method {
	case http.MethodGet:
		req.Tmpl = "register"
		req.Data = nil
		render(req)
	case http.MethodPost:
		user := &models.User{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := helpers.ValidateEmail(user.Email); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user.Password = password
		user = user.Find()

		if user.ID == 0 {
			user = user.Create()
		}

		sessionID := uuid.New().String()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID,
		}
		Sessions[sessionID] = user

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Profile is a controller for users to view and update their profile.
func Profile(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
	user := Sessions[cookie.Value]

	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		User     *models.User
		Skills   []*models.Skill
		Projects []*models.Project
	}{
		User:     user.GetProfile(),
		Skills:   user.GetSkills(),
		Projects: user.GetProjects(),
	}

	req.Data = data
	req.Layout = "dashboard-layout"

	//new := r.URL.Query().Get("new")
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	switch r.Method {
	case http.MethodGet:
		model := r.URL.Query().Get("edit")
		delete := r.URL.Query().Get("delete")
		if model == "user" && delete != "" && id != 0 {
			http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
			return
		}

		req.Tmpl = "page_profile"
		if r.Header.Get("HX-Request") == "true" {
			req.Tmpl = "profile"
		}

		render(req)

	case http.MethodPost:
		// Handle POST updates to user profile
		updatedUser := &models.User{
			ID:    user.ID,
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
		}
		user.UpdateUser(updatedUser)

		req.Data = data
		req.Tmpl = "profile-user"
		render(req)
	}
}

// Skills handles skill-related operations.
func Skills(w http.ResponseWriter, r *http.Request) {
	req := Request{
		W: w,
		R: r,
	}
	permission(req)

	cookie := getCookie(req.R)
	user := Sessions[cookie.Value]

	data := struct {
		User     *models.User
		Skills   []*models.Skill
		Projects []*models.Project
	}{
		User:     user.GetProfile(),
		Skills:   user.GetSkills(),
		Projects: user.GetProjects(),
	}

	req.Data = data
	req.Layout = "dashboard-layout"
	req.Tmpl = "edit-skills"

	render(req)
}
