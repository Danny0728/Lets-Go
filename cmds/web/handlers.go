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
	if r.URL.Path != "/" {
		app.notFound(w) //use of notFound helper
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: snippets})
}

//add a showsnippet handler function
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: snippet})
}

//add a createSnippet handler function
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.Header()["Content-Length"] = nil
		app.clientError(w, http.StatusMethodNotAllowed) //this sends the response to the user for us
		return
	}
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
