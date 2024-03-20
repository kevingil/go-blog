package controllers

import (
	"net/http"

	"github.com/kevingil/blog/app/cmd"
)

func Files(w http.ResponseWriter, r *http.Request) {
	cmd.Hx(w, r, "dashboard", "dashboard-files", data)
}
