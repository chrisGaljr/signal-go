package handlers

import (
	"net/http"
	"text/template"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	// not the most elegant way to handle 404s, but it gets the job done
	if r.URL.Path != "/" {
		tmpl, _ := template.ParseFiles("template/fourOhFour/index.html")
		tmpl.Execute(w, nil)
		return
	}

	tmpl, _ := template.ParseFiles("template/index.html")
	tmpl.Execute(w, nil)
}
