package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)
	dynamicMiddleware := alice.New()

	mux := pat.New()

	// main page
	// mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// mux.Post("/", dynamicMiddleware.ThenFunc(app.homeGetFiles))

	// mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	// mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))

	mux.Get("/", dynamicMiddleware.ThenFunc(app.redirectHome))
	mux.Get("/upload/:code", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/upload/:code", dynamicMiddleware.ThenFunc(app.homeGetFiles))
	mux.Get("/archive/:code", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/archive/download/:code", dynamicMiddleware.ThenFunc(app.getSnippet))

	mux.Get("/download", dynamicMiddleware.ThenFunc(app.createDownloadForm))

	// ping
	mux.Get("/ping", http.HandlerFunc(ping))

	// static
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}
