package controllers

import (
	"net/http"

	"github.com/kevingil/blog/app/views"
)

func Files(w http.ResponseWriter, r *http.Request) {
	views.Hx(w, r, "dashboard", "dashboard-files", data)
}
