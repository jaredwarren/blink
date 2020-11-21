package service

import (
	"html/template"
	"net/http"

	"github.com/jaredwarren/app"
)

// Register ...
func Register(a *app.Service) {
	m := a.Mux
	m.HandleFunc("/", home).Methods("GET")
}

// Close handler.
func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("layout.html", "base.html"))
	tmpl.ExecuteTemplate(w, "base", struct {
		PageTitle string
	}{
		"Something",
	})
}
