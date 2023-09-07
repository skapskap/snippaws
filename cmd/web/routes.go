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

	r.With(app.session.Enable).Get("/", http.HandlerFunc(app.home))
	r.With(app.session.Enable).Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	r.With(app.session.Enable).Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	r.With(app.session.Enable).Get("/snippet/{id}", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", fileServer).ServeHTTP(w, r)
	})

	return r
}
