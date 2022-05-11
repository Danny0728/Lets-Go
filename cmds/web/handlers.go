package main

/*The first thing we need is a handler.
handler's are like controllers from MVC
They are responsible for executing ur application logic and for writing HTTP response header and bodies*/
import (
	"fmt"
	"github.com/yash/snippetbox/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

//add a showsnippet handler function
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	//fmt.Fprintf(w, "Display a specific snippet with id = %d...\n", id)
	snippet, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a Snippet Form"))
}

//add a createSnippet handler function
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//creating some dummy builds we are going to delete this later
	title := "hello"
	content := "Hello is a nice word"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d \n", id), http.StatusSeeOther)
}
