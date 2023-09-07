package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/skapskap/snippaws/pkg/models"
	"net/http"
	"strconv"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return

	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	return
	//}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	created := time.Now()
	expiresValue := r.PostForm.Get("expires")

	expires, err := strconv.Atoi(expiresValue)
	if err != nil {
	}

	expiresTime := time.Now().AddDate(0, 0, expires)

	id, err := app.snippets.Insert(title, content, created, expiresTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
