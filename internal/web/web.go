package web

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func Render(w http.ResponseWriter, t string) {
	partials := []string{
		"/data/templates/base.layout.gohtml",
		"/data/templates/header.partial.gohtml",
		"/data/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, filepath.Join("/data/templates", t))

	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
