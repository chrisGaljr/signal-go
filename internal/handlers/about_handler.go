package handlers

import (
	"net/http"
	"text/template"
)

func AboutPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("template/about/index.html")

	tmpl.Execute(w, nil)
}
