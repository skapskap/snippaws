package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()

	r.Use(app.recoverPanic)
	r.Use(app.logRequest)
	r.Use(secureHeaders)

	r.Get("/", http.HandlerFunc(app.home))
	r.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	r.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	r.Get("/snippet/{id}", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
	})

	return r
}
