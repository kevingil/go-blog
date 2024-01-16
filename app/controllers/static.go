package controllers

import (
	"net/http"
)

// This just handles the page, Moderator is written in JS

func ModeratorJS(w http.ResponseWriter, r *http.Request) {
	Hx(w, r, "main_layout", "moderatorjs", data)
}
