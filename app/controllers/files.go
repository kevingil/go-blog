package controllers

import (
	"net/http"
)

func Files(w http.ResponseWriter, r *http.Request) {

	var data Context

	Hx(w, r, "dashboard", "dashboard-files", data)
}
