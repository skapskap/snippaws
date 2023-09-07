package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/skapskap/snippaws/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
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

	errors := make(map[string]string)

	if expiresValue == "" {
		errors["expires"] = "Este campo não pode ficar em branco"
	}

	expires, err := strconv.Atoi(expiresValue)
	if err != nil {
	}

	expiresTime := time.Now().AddDate(0, 0, expires)

	// Checar se o título do snippet tá em branco ou ultrapassa 100 caracteres

	if strings.TrimSpace(title) == "" {
		errors["title"] = "Este campo não pode ficar em branco."
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "Limite de caracteres ultrapassado (100)"
	}

	// Checar se o conteúdo do snippet tá em branco

	if strings.TrimSpace(content) == "" {
		errors["content"] = "Este campo não pode ficar em branco"
	}

	// Checar se o campo de expiração não está em branco e se bate com um dos valores permitidos

	if expires != 365 && expires != 7 && expires != 1 {
		errors["expires"] = "Este campo é inválido"
	}

	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, created, expiresTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
