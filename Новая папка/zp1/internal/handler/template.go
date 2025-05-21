package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}, useLayout bool) {
	var templates []string
	if useLayout {
		templates = append(templates, "templates/layout.html")
	}
	templates = append(templates, filepath.Join("templates", tmpl))

	t, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
