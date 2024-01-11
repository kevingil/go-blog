package controllers


func Publish(w http.ResponseWriter, r *http.Request) {
	permission(w, r)
	cookie := getCookie(r)
	user := Sessions[cookie.Value]
	tmpl := "publish"
	var response bytes.Buffer

	if r.Header.Get("HX-Request") == "true" {
		if user != nil {
			data.Articles = user.FindArticles()
		}
		if err := views.Tmpl.ExecuteTemplate(&response, tmpl, data); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	io.WriteString(w, response.String())
}
