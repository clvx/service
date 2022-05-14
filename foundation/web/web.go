// Package web contains a small web framework extension
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles an http request withing our own little mini
// framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App structu.
type App struct {
	*httptreemux.ContextMux // embedding the mux
							// inner type promotion which means everything 
							// related to the inner type gets promoted to the 
							// outer type. In other words, type App is everything
							// the ContextMux is plus shutdown. It's a way to 
							// do a mux implementation without having to rewrite it
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown: shutdown, // a way to initialize a clean shutdown from inside the application
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and a path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		// PRE CODE PROCESSING
		if err := handler(r.Context(), w, r); err != nil {
			// ERROR HANDLING
			return
		}
		// POST CODE PROCESSING
	}
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.ContextMux.Handle(method, finalPath, h)

}
