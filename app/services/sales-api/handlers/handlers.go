// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/clvx/service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/clvx/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/clvx/service/foundation/web"
	"go.uber.org/zap"
)

// We want packages that provides not contains.
// Package that contains like utils, models, helpers flip the piramid upside down
// creating a cascading dependency effect.
// A way to identify a package that provides is having a file named after the
// package, e.g. handlers/handlers.go. It's not the case with utils/ which might
// different files for different purposes.

// DebugStandardLibrary registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the DefaultServerMux
// would be a security risk since a dependency could inject a handler into our
// service without us kowing it.
func DebugStandardLibraryMux() *http.ServeMux{
	mux := http.NewServeMux()

	//Register all the standard library debug endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom 
// debug application routes for the service. This bypassing the use of the 
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints
	cgh := checkgrp.Handlers{
		Build: build,
		Log: log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}


// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log *zap.SugaredLogger
}

// APIMux constructs an http.Handler with all application routes defined
// Guideline: when it comes to an api you have data coming in (input) you can use a concrete type if you wanna use data based on what it is, or use an interface if you wanna use data based on what it can do. When the api returns data, use concrete data unless it's an error (error interface) or if you need an empty interface (discouraged)
func APIMux(cfg APIMuxConfig) *web.App{

	//construct the web.App which holds all routes as well.
	app := web.NewApp(
		cfg.Shutdown,
	)

	v1(app, cfg)
	return app
}

// v1 binds all the version 1 routes
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, "v1", "/test", tgh.Test)
}
