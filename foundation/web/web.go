// Package web contains a small web framework extension
package web

import (
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

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
