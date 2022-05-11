package main

/*
the second thing we need is a router or serveMux in go terminology
This stores a mapping between the URL patterns for your application
and the corresponding handlers. Usually you have one servemux for
your application containing all your routes.
*/
import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	// mux.Handle("/", http.HandlerFunc(home))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
